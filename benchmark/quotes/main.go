package main

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/actor"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/logger"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/stream"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/utils"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

// requires AUTH_HEADER and PRIVATE_KEY to work.

const (
	maxRuntime      = 10 * time.Second
	swapMint        = "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
	pollingInterval = time.Second
)

func main() {
	app := &cli.App{
		Name:  "benchmark-quotes",
		Usage: "Compares Solana Trader API AMM quotes with Jupiter API",
		Flags: []cli.Flag{
			utils.APIWSEndpoint,
			// utils.OutputFileFlag,
			MintFlag,
			TriggerActivityFlag,
			IterationsFlag,
		},
		Action: run,
	}

	err := app.Run(os.Args)
	defer func() {
		if logger.Log() != nil {
			_ = logger.Log().Sync()
		}
	}()

	if err != nil {
		panic(err)
	}
}

func run(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mint := c.String(MintFlag.Name)
	iterations := c.Int(IterationsFlag.Name)
	triggerActivity := c.Bool(TriggerActivityFlag.Name)

	syncedTicker := time.NewTicker(pollingInterval)
	defer syncedTicker.Stop()

	jupiterAPI, err := stream.NewJupiterAPI(stream.WithJupiterToken(mint), stream.WithJupiterTicker(syncedTicker))
	if err != nil {
		return err
	}

	traderAPIWS, err := stream.NewTraderWSPrice(stream.WithTraderWSMint(mint))
	if err != nil {
		return err
	}

	traderAPIHTTP, err := stream.NewTraderHTTPPriceStream(stream.WithTraderHTTPMint(mint), stream.WithTraderHTTPTicker(syncedTicker))
	if err != nil {
		return err
	}

	jupiterActor, err := actor.NewJupiterSwap(actor.WithJupiterTokenPair(swapMint, mint))
	if err != nil {
		return err
	}

	var (
		tradeWSUpdates    []stream.RawUpdate[*pb.GetPricesStreamResponse]
		tradeHTTPUpdates  []stream.RawUpdate[stream.DurationUpdate[*pb.GetPriceResponse]]
		jupiterUpdates    []stream.RawUpdate[stream.DurationUpdate[stream.JupiterPriceResponse]]
		errCh             = make(chan error, 2)
		runCtx, runCancel = context.WithTimeout(ctx, maxRuntime)
	)
	defer runCancel()

	go func() {
		var err error

		jupiterUpdates, err = jupiterAPI.Run(runCtx)
		if err != nil {
			errCh <- errors.Wrap(err, "could not collect results from Solana")
			return
		}

		errCh <- nil
	}()

	go func() {
		var err error

		tradeWSUpdates, err = traderAPIWS.Run(runCtx)
		if err != nil {
			errCh <- errors.Wrap(err, "could not collect results from Trader API")
			return
		}

		errCh <- nil
	}()

	go func() {
		var err error

		tradeHTTPUpdates, err = traderAPIHTTP.Run(runCtx)
		if err != nil {
			errCh <- errors.Wrap(err, "could not collect results from Trader API")
			return
		}

		errCh <- nil
	}()

	if triggerActivity {
		err = jupiterActor.Swap(runCtx, iterations)
		if err != nil {
			return err
		}

		runCancel()
	}

	// wait for routines to exit
	completionCount := 0
	for completionCount < 3 {
		select {
		case runErr := <-errCh:
			completionCount++
			if runErr != nil {
				logger.Log().Errorw("fatal error during runtime: exiting", "err", err)
				return runErr
			}
		}
	}

	fmt.Println(jupiterUpdates)
	fmt.Println(tradeWSUpdates)
	fmt.Println(tradeHTTPUpdates)

	return nil
}

var (
	MintFlag = &cli.StringFlag{
		Name:  "mint",
		Usage: "mint to fetch price for (inactive token is best)",
		Value: "zebeczgi5fSEtbpfQKVZKCJ3WgYXxjkMUkNNx7fLKAF", // zbc
	}

	TriggerActivityFlag = &cli.BoolFlag{
		Name:  "trigger-activity",
		Usage: "if true, send trigger transactions to force quote updates (requires PRIVATE_KEY environment variable_",
		// Value: true,
	}

	IterationsFlag = &cli.IntFlag{
		Name:  "iterations",
		Usage: "number of quotes to compare",
		Value: 10,
	}
)
