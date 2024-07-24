package main

import (
	"context"
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

const updateInterval = 10 * time.Second

const APIGRPCEndpoint = "localhost:9000"
const GeyserGRPCEndpoint = "185.209.178.215:6677"

// SOL-USDC
var pools = []string{
	"58oQChx4yWmvKdwLLZzBi4ChoCc2fqCUWBkwMihLYQo2",
	"DQyrAcCrDXQ7NeoqGgDCZwBvWDcYmFCjSb9JtteuvPpz",
	"HLmqeL62xR1QoZ1HKKbXRrdN1p3phKpxRMb2VVopvBBz",
	"HmiHHzq4Fym9e1D4qzLS6LDDM3tNsCTBPDWHTLZ763jY",
}

func main() {
	app := &cli.App{
		Name:  "benchmark-poolreserves",
		Usage: "Compares Solana Trader API pool reserves stream with a direct Geyser connection to determine the efficacy of using the Solana Trader API stream",
		Flags: []cli.Flag{
			DurationFlag,
			utils.OutputFileFlag,
			RemoveUnmatchedFlag,
			RemoveDuplicatesFlag,
			APIGRPCEndpointFlag,
			GeyserGRPCEndpointFlag,
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

	apiGRPCEndpoint := c.String(APIGRPCEndpointFlag.Name)
	if apiGRPCEndpoint == "" {
		apiGRPCEndpoint = APIGRPCEndpoint
	}

	geyserGRPCEndpoint := c.String(GeyserGRPCEndpointFlag.Name)
	if geyserGRPCEndpoint == "" {
		geyserGRPCEndpoint = GeyserGRPCEndpoint
	}

	authHeader, ok := os.LookupEnv("AUTH_HEADER")
	if !ok {
		return errors.New("AUTH_HEADER not set in environment")
	}

	logger.Log().Infow("Attempting connections with endpoints: ", "api endpoint", apiGRPCEndpoint, "geyser endpoint", geyserGRPCEndpoint)
	logger.Log().Infow("Requested pool keys: ", "keys", pools)

	traderOS, err := stream.NewTraderAPIGRPCStream(apiGRPCEndpoint, authHeader, pools)
	if err != nil {
		return err
	}

	geyserOS, err := stream.NewGeyserGRPCStream(ctx, geyserGRPCEndpoint, pools)
	if err != nil {
		return err
	}

	duration := c.Duration(DurationFlag.Name)
	runCtx, runCancel := context.WithTimeout(ctx, duration)
	defer runCancel()

	traderUpdatesChan := make(chan stream.RawUpdate[stream.TraderAPIRawUpdateGRPC], 100)
	geyserUpdatesChan := make(chan stream.RawUpdate[stream.GeyserRawUpdateGRPC], 100)
	errCh := make(chan error, 2)

	go func() {
		errCh <- traderOS.Run(runCtx, traderUpdatesChan, pools)
	}()

	go func() {
		errCh <- geyserOS.Run(runCtx, geyserUpdatesChan)
	}()

	var traderUpdates []stream.RawUpdate[stream.TraderAPIRawUpdateGRPC]
	var geyserUpdates []stream.RawUpdate[stream.GeyserRawUpdateGRPC]

	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	for {
		select {
		case update := <-traderUpdatesChan:
			traderUpdates = append(traderUpdates, update)
		case update := <-geyserUpdatesChan:
			geyserUpdates = append(geyserUpdates, update)
		case err := <-errCh:
			logger.Log().Errorw("Stream error", "err", err)
		case <-ticker.C:
			logger.Log().Infof("Collected %d trader updates and %d geyser updates", len(traderUpdates), len(geyserUpdates))
		case <-runCtx.Done():
			logger.Log().Info("Run time completed")
			goto ProcessResults
		}
	}

ProcessResults:
	logger.Log().Infow("finished collecting data points", "tradercount", len(traderUpdates), "geysercount", len(geyserUpdates))

	removeDuplicates := c.Bool(RemoveDuplicatesFlag.Name)
	traderResults, traderDuplicates, err := traderOS.Process(traderUpdates, removeDuplicates)
	if err != nil {
		return errors.Wrap(err, "could not process trader API updates")
	}
	logger.Log().Debugw("processed trader API results", "range", output.FormatSortRange(traderResults), "count", len(traderResults), "duplicaterange", output.FormatSortRange(traderDuplicates), "duplicatecount", len(traderDuplicates))

	geyserResults, geyserDuplicates, err := geyserOS.Process(geyserUpdates, removeDuplicates)
	if err != nil {
		return errors.Wrap(err, "could not process geyser results")
	}
	logger.Log().Debugw("processed geyser results", "range", output.FormatSortRange(geyserResults), "count", len(geyserResults), "duplicaterange", output.FormatSortRange(geyserDuplicates), "duplicatecount", len(geyserDuplicates))

	slots := SlotRange(traderResults, geyserResults)
	logger.Log().Debugw("finished processing data points", "startSlot", slots[0], "endSlot", slots[len(slots)-1], "count", len(slots))

	datapoints, _, _, err := Merge(slots, traderResults, geyserResults)
	if err != nil {
		return err
	}

	logger.Log().Infow("completed merging: outputting data...")

	// dump results to stdout
	removeUnmatched := c.Bool(RemoveUnmatchedFlag.Name)
	PrintSummary(duration, apiGRPCEndpoint, geyserGRPCEndpoint, datapoints)

	// write results to csv
	outputFile := c.String(utils.OutputFileFlag.Name)
	header := []string{"slot", "diff", "trader-api-time", "geyser-time"}
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
	APIGRPCEndpointFlag = &cli.StringFlag{
		Name:  "api-endpoint",
		Usage: "Override the default API GRPC endpoint",
	}
	GeyserGRPCEndpointFlag = &cli.StringFlag{
		Name:  "geyser-endpoint",
		Usage: "Override the default Geyser GRPC endpoint",
	}
)
