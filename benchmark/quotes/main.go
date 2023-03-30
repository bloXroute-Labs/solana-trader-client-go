package main

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/actor"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/logger"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/stream"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/utils"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"os"
	"time"
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
			PublicKeyFlag,
			MaxRuntimeFlag,
			SwapAmountFlag,
			SwapMintFlag,
			SwapIntervalFlag,
			SwapInitialWaitFlag,
			SwapAfterWaitFlag,
			QueryIntervalFlag,
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

	_, ok := os.LookupEnv("AUTH_HEADER")
	if !ok {
		return errors.New("AUTH_HEADER not set in environment")
	}

	var (
		mint            = c.String(MintFlag.Name)
		iterations      = c.Int(IterationsFlag.Name)
		triggerActivity = c.Bool(TriggerActivityFlag.Name)
		publicKey       = c.String(PublicKeyFlag.Name)

		maxRuntime      = c.Duration(MaxRuntimeFlag.Name)
		swapAmount      = c.Float64(SwapAmountFlag.Name)
		swapMint        = c.String(SwapMintFlag.Name)
		swapInterval    = c.Duration(SwapIntervalFlag.Name)
		swapInitialWait = c.Duration(SwapInitialWaitFlag.Name)
		swapAfterWait   = c.Duration(SwapAfterWaitFlag.Name)
		queryInterval   = c.Duration(QueryIntervalFlag.Name)
	)

	if triggerActivity {
		_, ok := os.LookupEnv("PRIVATE_KEY")
		if !ok {
			return errors.New("PRIVATE_KEY not set in environment when --trigger-activity set")
		}
	}

	syncedTicker := time.NewTicker(queryInterval)
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

	privateKey, err := transaction.LoadPrivateKeyFromEnv()
	if err != nil {
		return err
	}
	rpcOpts := provider.RPCOpts{
		Endpoint:   "http://18.208.115.90:1809",
		PrivateKey: &privateKey,
		AuthHeader: "ZDJhYjkzYmEtMWE4Yi00MTg3LTk5NGUtYzYzODk2YzkzNmUzOmE2MTY4MWE5NDU2Y2EzMTlhOTAwMzZlODM2MWRiYzcz",
	}
	client := provider.NewHTTPClientWithOpts(nil, rpcOpts)
	jupiterActor, err := actor.NewJupiterSwap(actor.WithJupiterTokenPair(swapMint, mint), actor.WithJupiterPublicKey(publicKey), actor.WithJupiterInitialTimeout(swapInitialWait), actor.WithJupiterAfterTimeout(swapAfterWait), actor.WithJupiterInterval(swapInterval), actor.WithJupiterAmount(swapAmount), actor.WithJupiterClient(client))
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

	var swaps []actor.SwapEvent
	if triggerActivity {
		swaps, err = jupiterActor.Swap(runCtx, iterations)
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

	fmt.Println()
	result := benchmarkResult{
		mint:             mint,
		swaps:            swaps,
		jupiterUpdates:   jupiterUpdates,
		tradeWSUpdates:   tradeWSUpdates,
		tradeHTTPUpdates: tradeHTTPUpdates,
	}

	result.PrintSummary()
	result.PrintRaw()
	return nil
}
