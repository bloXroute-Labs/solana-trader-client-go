package main

import (
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/csv"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/logger"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/output"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/stream"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/utils"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

import (
	"context"
)

const updateInterval = 10 * time.Second

func main() {
	app := &cli.App{
		Name:  "benchmark-traderapi",
		Usage: "Compares Solana Trader API orderbook stream with a direct connection to a Solana node to determine the efficacy of using the Solana Trader API stream",
		Flags: []cli.Flag{
			utils.SolanaHTTPRPCEndpointFlag,
			utils.SolanaWSRPCEndpointFlag,
			utils.APIWSEndpoint,
			MarketAddrFlag,
			DurationFlag,
			utils.OutputFileFlag,
			RemoveUnmatchedFlag,
			RemoveDuplicatesFlag,
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

	marketAddr := c.String(MarketAddrFlag.Name)
	traderAPIEndpoint := c.String(utils.APIWSEndpoint.Name)
	solanaRPCEndpoint := c.String(utils.SolanaHTTPRPCEndpointFlag.Name)
	solanaWSEndpoint := c.String(utils.SolanaWSRPCEndpointFlag.Name)

	authHeader, ok := os.LookupEnv("AUTH_HEADER")
	if !ok {
		return errors.New("AUTH_HEADER not set in environment")
	}

	connectCtx, connectCancel := context.WithTimeout(ctx, 5*time.Second)
	defer connectCancel()

	traderOS, err := stream.NewAPIOrderbookStream(traderAPIEndpoint, marketAddr, authHeader)
	if err != nil {
		return err
	}

	solanaOS, err := stream.NewSolanaOrderbookStream(connectCtx, solanaRPCEndpoint, solanaWSEndpoint, marketAddr)
	if err != nil {
		return err
	}

	startTime := time.Now()
	duration := c.Duration(DurationFlag.Name)
	runCtx, runCancel := context.WithTimeout(ctx, duration)
	defer runCancel()

	var (
		tradeUpdates  []stream.RawUpdate[[]byte]
		solanaUpdates []stream.RawUpdate[stream.SolanaRawUpdate]
	)
	errCh := make(chan error, 2)

	go func() {
		var err error

		tradeUpdates, err = traderOS.Run(runCtx)
		if err != nil {
			errCh <- errors.Wrap(err, "could not collect results from trader API")
			return
		}

		errCh <- nil
	}()

	go func() {
		var err error

		solanaUpdates, err = solanaOS.Run(runCtx)
		if err != nil {
			errCh <- errors.Wrap(err, "could not collect results from Solana")
			return
		}

		errCh <- nil
	}()

	completionCount := 0
	ticker := time.NewTicker(updateInterval)

	for completionCount < 2 {
		select {
		case runErr := <-errCh:
			completionCount++
			if runErr != nil {
				logger.Log().Errorw("fatal error during runtime: exiting", "err", err)
				return runErr
			}
		case <-ticker.C:
			elapsedTime := time.Now().Sub(startTime).Round(time.Second)
			logger.Log().Infof("waited %v out of %v...", elapsedTime, duration)
		}

	}

	logger.Log().Infow("finished collecting data points", "tradercount", len(tradeUpdates), "solanacount", len(solanaUpdates))

	removeDuplicates := c.Bool(RemoveDuplicatesFlag.Name)
	traderResults, traderDuplicates, err := traderOS.Process(tradeUpdates, removeDuplicates)
	if err != nil {
		return errors.Wrap(err, "could not process trader API updates")
	}
	logger.Log().Debugw("processed trader API results", "range", output.FormatSortRange(traderResults), "count", len(traderResults), "duplicaterange", output.FormatSortRange(traderDuplicates), "duplicatecount", len(traderDuplicates))

	solanaResults, solanaDuplicates, err := solanaOS.Process(solanaUpdates, removeDuplicates)
	if err != nil {
		return errors.Wrap(err, "could not process solana results")
	}

	logger.Log().Debugw("processed solana results", "range", output.FormatSortRange(solanaResults), "count", len(solanaResults), "duplicaterange", output.FormatSortRange(solanaDuplicates), "duplicatecount", len(solanaDuplicates))

	slots := SlotRange(traderResults, solanaResults)
	logger.Log().Debugw("finished processing data points", "startSlot", slots[0], "endSlot", slots[len(slots)-1], "count", len(slots))

	datapoints, _, _, err := Merge(slots, traderResults, solanaResults)
	if err != nil {
		return err
	}

	logger.Log().Infow("completed merging: outputting data...")

	// dump results to stdout
	removeUnmatched := c.Bool(RemoveUnmatchedFlag.Name)
	PrintSummary(duration, traderAPIEndpoint, solanaWSEndpoint, datapoints)

	// write results to csv
	outputFile := c.String(utils.OutputFileFlag.Name)
	header := []string{"slot", "diff", "seq", "trader-api-time", "solana-side", "solana-time"}
	err = csv.Write(outputFile, header, datapoints, func(line []string) bool {
		if removeUnmatched {
			for _, col := range line {
				return col == "n/a"
			}
		}
		return false
	})
	if err != nil {
		return err
	}

	return nil
}

var (
	MarketAddrFlag = &cli.StringFlag{
		Name:  "market",
		Usage: "market to run analysis for",
		Value: "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh6", // SOL/USDC market
	}
	DurationFlag = &cli.DurationFlag{
		Name:     "run-time",
		Usage:    "amount of time to run script for (seconds)",
		Required: true,
	}
	RemoveUnmatchedFlag = &cli.BoolFlag{
		Name:  "remove-unmatched",
		Usage: "skip events without a match from other source",
	}
	RemoveDuplicatesFlag = &cli.BoolFlag{
		Name:  "remove-duplicates",
		Usage: "skip events that are identical to the previous",
	}
)
