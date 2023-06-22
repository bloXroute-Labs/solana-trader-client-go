package main

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/actor"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/logger"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/stream"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/utils"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/gagliardetto/solana-go"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

// requires AUTH_HEADER and PRIVATE_KEY to work.

func main() {
	utils.SolanaHTTPRPCEndpointFlag.Required = true
	utils.SolanaHTTPRPCEndpointFlag.Value = "https://api.mainnet-beta.solana.com"

	app := &cli.App{
		Name:  "benchmark-quotes",
		Usage: "Compares Solana Trader API AMM quotes with Jupiter API",
		Flags: []cli.Flag{
			utils.APIWSEndpoint,
			utils.OutputFileFlag,
			utils.SolanaHTTPRPCEndpointFlag,
			MintFlag,
			MintDecimalsFlag,
			TriggerActivityFlag,
			IterationsFlag,
			PublicKeyFlag,
			MaxRuntimeFlag,
			SwapAmountFlag,
			SwapMintFlag,
			SwapIntervalFlag,
			SwapInitialWaitFlag,
			SwapAfterWaitFlag,
			SwapAlternateFlag,
			QueryIntervalFlag,
			EnvFlag,
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
		env             = c.String(EnvFlag.Name)
		mint            = c.String(MintFlag.Name)
		mintDecimals    = c.Int(MintDecimalsFlag.Name)
		iterations      = c.Int(IterationsFlag.Name)
		triggerActivity = c.Bool(TriggerActivityFlag.Name)
		publicKey       = c.String(PublicKeyFlag.Name)

		maxRuntime      = c.Duration(MaxRuntimeFlag.Name)
		swapAmount      = c.Float64(SwapAmountFlag.Name)
		swapMint        = c.String(SwapMintFlag.Name)
		swapInterval    = c.Duration(SwapIntervalFlag.Name)
		swapInitialWait = c.Duration(SwapInitialWaitFlag.Name)
		swapAlternate   = c.Bool(SwapAlternateFlag.Name)
		swapAfterWait   = c.Duration(SwapAfterWaitFlag.Name)
		queryInterval   = c.Duration(QueryIntervalFlag.Name)

		rpcEndpoint = c.String(utils.SolanaHTTPRPCEndpointFlag.Name)
		outputFile  = c.String(utils.OutputFileFlag.Name)
	)

	if triggerActivity {
		_, ok := os.LookupEnv("PRIVATE_KEY")
		if !ok {
			return errors.New("PRIVATE_KEY not set in environment when --trigger-activity set")
		}
	}

	httpClient, wsClient, err := traderClients(env)
	if err != nil {
		return err
	}

	logger.Log().Infow("trader API clients connected", "env", env)

	syncedTicker := time.NewTicker(queryInterval)
	defer syncedTicker.Stop()

	jupiterAPI, err := stream.NewJupiterAPI(stream.WithJupiterToken(mint, mintDecimals), stream.WithJupiterTicker(syncedTicker))
	if err != nil {
		return err
	}

	traderAPIWS, err := stream.NewTraderWSPrice(stream.WithTraderWSMint(mint), stream.WithTraderWSClient(wsClient))
	if err != nil {
		return err
	}

	traderAPIHTTP, err := stream.NewTraderHTTPPriceStream(stream.WithTraderHTTPMint(mint), stream.WithTraderHTTPTicker(syncedTicker), stream.WithTraderHTTPClient(httpClient))
	if err != nil {
		return err
	}

	jupiterActor, err := actor.NewJupiterSwap(actor.WithJupiterTokenPair(swapMint, mint), actor.WithJupiterPublicKey(publicKey), actor.WithJupiterInitialTimeout(swapInitialWait), actor.WithJupiterAfterTimeout(swapAfterWait), actor.WithJupiterInterval(swapInterval), actor.WithJupiterAmount(swapAmount), actor.WithJupiterClient(httpClient), actor.WithJupiterAlternate(swapAlternate))
	if err != nil {
		return err
	}

	var (
		tradeWSUpdates    []stream.RawUpdate[*pb.GetPricesStreamResponse]
		tradeHTTPUpdates  []stream.RawUpdate[stream.DurationUpdate[*pb.GetPriceResponse]]
		jupiterUpdates    []stream.RawUpdate[stream.DurationUpdate[*stream.JupiterPriceResponse]]
		errCh             = make(chan error, 2)
		runCtx, runCancel = context.WithTimeout(ctx, maxRuntime)
	)
	defer runCancel()

	logger.Log().Infow("starting all routines", "duration", maxRuntime)

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

	jupiterProcessedUpdate, jupiterDuplicates, _ := jupiterAPI.Process(jupiterUpdates, true)
	logger.Log().Infow("ignoring jupiter duplicates", "count", len(jupiterDuplicates))

	tradeWSProcessedUpdate, tradeWSDuplicates, _ := traderAPIWS.Process(tradeWSUpdates, true)
	logger.Log().Infow("ignoring tradeWS duplicates", "count", len(tradeWSDuplicates))

	tradeHTTPProcessedUpdate, tradeHTTPDuplicates, _ := traderAPIHTTP.Process(tradeHTTPUpdates, true)
	logger.Log().Infow("ignoring tradeWS duplicates", "count", len(tradeHTTPDuplicates))

	swapUpdates, err := fetchTransactionInfo(ctx, swaps, rpcEndpoint)
	if err != nil {
		return err
	}

	fmt.Println()
	result := benchmarkResult{
		mint:                      mint,
		swaps:                     swaps,
		jupiterRawUpdates:         jupiterUpdates,
		tradeWSRawUpdates:         tradeWSUpdates,
		tradeHTTPRawUpdates:       tradeHTTPUpdates,
		jupiterProcessedUpdates:   jupiterProcessedUpdate,
		tradeWSProcessedUpdates:   tradeWSProcessedUpdate,
		tradeHTTPProcessedUpdates: tradeHTTPProcessedUpdate,
		swapProcessedUpdates:      swapUpdates,
	}

	result.PrintSummary()
	result.PrintRaw()
	return result.WriteCSV(outputFile)
}

func fetchTransactionInfo(ctx context.Context, swaps []actor.SwapEvent, rpcEndpoint string) (map[int][]stream.ProcessedUpdate[stream.QuoteResult], error) {
	rpc := solanarpc.New(rpcEndpoint)
	transactionEvents, err := utils.AsyncGather(ctx, swaps, func(i int, ctx context.Context, t actor.SwapEvent) (r stream.ProcessedUpdate[stream.QuoteResult], err error) {
		maxVersion := uint64(0)
		swap := swaps[i]
		result, err := rpc.GetTransaction(ctx, solana.MustSignatureFromBase58(swap.Signature), &solanarpc.GetTransactionOpts{
			Encoding:                       solana.EncodingBase64,
			Commitment:                     solanarpc.CommitmentConfirmed,
			MaxSupportedTransactionVersion: &maxVersion,
		})
		if err != nil {
			return
		}

		r = stream.ProcessedUpdate[stream.QuoteResult]{
			Timestamp: result.BlockTime.Time(),
			Slot:      int(result.Slot),
			Data: stream.QuoteResult{
				Elapsed:   result.BlockTime.Time().Sub(swap.Timestamp),
				BuyPrice:  0,
				SellPrice: 0,
				Source:    fmt.Sprintf("transaction-%v-%v", swap.Signature, swap.Info),
			},
		}
		return
	})

	if err != nil {
		return nil, err
	}

	result := make(map[int][]stream.ProcessedUpdate[stream.QuoteResult])
	for _, event := range transactionEvents {
		result[event.Slot] = append(result[event.Slot], event)
	}
	return result, nil
}
