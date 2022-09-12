package main

import (
	"github.com/bloXroute-Labs/serum-client-go/benchmark/internal/arrival"
	"github.com/bloXroute-Labs/serum-client-go/benchmark/internal/csv"
	"github.com/bloXroute-Labs/serum-client-go/benchmark/internal/logger"
	"github.com/bloXroute-Labs/serum-client-go/benchmark/internal/utils"
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
		Name:  "benchmark-serumapi",
		Usage: "Compares Serum API orderbook stream with a direct connection to a Solana node to determine the efficacy of using the Serum API stream",
		Flags: []cli.Flag{
			utils.SolanaHTTPRPCEndpointFlag,
			utils.SolanaWSRPCEndpointFlag,
			SerumWSEndpointFlag,
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
	serumEndpoint := c.String(SerumWSEndpointFlag.Name)
	solanaRPCEndpoint := c.String(utils.SolanaHTTPRPCEndpointFlag.Name)
	solanaWSEndpoint := c.String(utils.SolanaWSRPCEndpointFlag.Name)

	authHeader, ok := os.LookupEnv("AUTH_HEADER")
	if !ok {
		return errors.New("AUTH_HEADER not set in environment")
	}

	serumOS, err := arrival.NewSerumOrderbookStream(serumEndpoint, marketAddr, authHeader)
	if err != nil {
		return err
	}
	solanaOS, err := arrival.NewSolanaOrderbookStream(ctx, solanaRPCEndpoint, solanaWSEndpoint, marketAddr)
	if err != nil {
		return err
	}

	startTime := time.Now()
	duration := c.Duration(DurationFlag.Name)
	runCtx, runCancel := context.WithTimeout(ctx, duration)
	defer runCancel()

	var (
		serumUpdates  []arrival.StreamUpdate[[]byte]
		solanaUpdates []arrival.StreamUpdate[arrival.SolanaRawUpdate]
	)
	errCh := make(chan error, 2)

	go func() {
		var err error

		serumUpdates, err = serumOS.Run(runCtx)
		if err != nil {
			errCh <- errors.Wrap(err, "could not collect results from Serum")
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

	logger.Log().Infow("finished collecting data points", "serumcount", len(serumUpdates), "solanacount", len(solanaUpdates))

	removeDuplicates := c.Bool(RemoveDuplicatesFlag.Name)
	serumResults, serumDuplicates, err := serumOS.Process(serumUpdates, removeDuplicates)
	if err != nil {
		return errors.Wrap(err, "could not process serum updates")
	}
	logger.Log().Debugw("processed serum results", "range", FormatSortRange(serumResults), "count", len(serumResults), "duplicaterange", FormatSortRange(serumDuplicates), "duplicatecount", len(serumDuplicates))

	solanaResults, solanaDuplicates, err := solanaOS.Process(solanaUpdates, removeDuplicates)
	if err != nil {
		return errors.Wrap(err, "could not process solana results")
	}

	logger.Log().Debugw("processed solana results", "range", FormatSortRange(solanaResults), "count", len(solanaResults), "duplicaterange", FormatSortRange(solanaDuplicates), "duplicatecount", len(solanaDuplicates))

	slots := SlotRange(serumResults, solanaResults)
	logger.Log().Debugw("finished processing data points", "startSlot", slots[0], "endSlot", slots[len(slots)-1], "count", len(slots))

	datapoints, _, _, err := Merge(slots, serumResults, solanaResults)
	if err != nil {
		return err
	}

	logger.Log().Infow("completed merging: outputting data...")

	// dump results to stdout
	removeUnmatched := c.Bool(RemoveUnmatchedFlag.Name)
	PrintSummary(duration, serumEndpoint, solanaWSEndpoint, datapoints)

	// write results to csv
	outputFile := c.String(utils.OutputFileFlag.Name)
	header := []string{"slot", "diff", "seq", "serum-time", "solana-side", "solana-time"}
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

	arrival.OrderbookEqualIndex(nil, nil)

	return nil
}

var (
	SerumWSEndpointFlag = &cli.StringFlag{
		Name:  "serum-ws-endpoint",
		Usage: "Serum API websocket connection endpoint",
		Value: "wss://virginia.solana.dex.blxrbdn.com/ws",
	}
	MarketAddrFlag = &cli.StringFlag{
		Name:  "market",
		Usage: "market to run analysis for",
		Value: "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT", // SOL/USDC market
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
