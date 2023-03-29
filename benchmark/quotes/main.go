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
)

// requires AUTH_HEADER and PRIVATE_KEY to work.

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

	// traderAPIEndpoint := c.String(utils.APIWSEndpoint.Name)
	mint := c.String(MintFlag.Name)
	iterations := c.Int(IterationsFlag.Name)

	jupiterAPI, err := stream.NewJupiterAPI(stream.WithJupiterToken(mint))
	if err != nil {
		return err
	}

	traderAPI, err := stream.NewTraderPriceStream(stream.WithTraderWSMint(mint))
	if err != nil {
		return err
	}

	jupiterActor, err := actor.NewJupiterSwap()
	if err != nil {
		return err
	}

	var (
		tradeUpdates      []stream.RawUpdate[*pb.GetPricesStreamResponse]
		jupiterUpdates    []stream.RawUpdate[stream.JupiterPriceResponse]
		errCh             = make(chan error, 2)
		runCtx, runCancel = context.WithCancel(ctx)
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

		tradeUpdates, err = traderAPI.Run(runCtx)
		if err != nil {
			errCh <- errors.Wrap(err, "could not collect results from Trader API")
			return
		}

		errCh <- nil
	}()

	// do swaps
	err = jupiterActor.Swap(runCtx, iterations)
	if err != nil {
		return err
	}

	// wait for routines to exit
	runCancel()
	completionCount := 0
	for completionCount < 2 {
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
	fmt.Println(tradeUpdates)

	return nil
}

var (
	MintFlag = &cli.StringFlag{
		Name:  "mint",
		Usage: "mint to fetch price for (inactive token is best)",
		Value: "zebeczgi5fSEtbpfQKVZKCJ3WgYXxjkMUkNNx7fLKAF", // zbc
	}

	// InputMintFlag = &cli.StringFlag{
	// 	Name:  "input-mint",
	// 	Usage: "input mint to fetch quote for (inactive pairs are best)",
	// 	Value: "So11111111111111111111111111111111111111112",
	// }
	//
	// OutputMintFlag = &cli.StringFlag{
	// 	Name:  "output-mint",
	// 	Usage: "output mint to fetch quote for (inactive pairs are best)",
	// 	// Value: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
	// 	Value: "zebeczgi5fSEtbpfQKVZKCJ3WgYXxjkMUkNNx7fLKAF", // zbc
	// }

	TriggerActivityFlag = &cli.BoolFlag{
		Name:  "trigger-activity",
		Usage: "if true, send trigger transaction to force quote updates",
		Value: true,
	}

	IterationsFlag = &cli.IntFlag{
		Name:  "iterations",
		Usage: "number of quotes to compare",
		Value: 10,
	}
)
