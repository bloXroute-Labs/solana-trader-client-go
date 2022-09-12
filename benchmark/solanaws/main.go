package main

import (
	"github.com/bloXroute-Labs/serum-client-go/benchmark/internal/arrival"
	"github.com/bloXroute-Labs/serum-client-go/benchmark/internal/logger"
	"github.com/bloXroute-Labs/serum-client-go/benchmark/internal/utils"
	gserum "github.com/gagliardetto/solana-go/programs/serum"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/maps"
	"os"
	"sort"
	"time"
)

import (
	"context"
)

const updateInterval = 10 * time.Second

func main() {
	app := &cli.App{
		Name:  "benchmark-solanaws",
		Usage: "Compares two Solana websocket streams for any differences",
		Flags: []cli.Flag{
			utils.SolanaHTTPRPCEndpointFlag,
			SolanaWSEndpoint1Flag,
			SolanaWSEndpoint2Flag,
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
	solanaRPCEndpoint := c.String(utils.SolanaHTTPRPCEndpointFlag.Name)
	solanaWSEndpoint1 := c.String(SolanaWSEndpoint1Flag.Name)
	solanaWSEndpoint2 := c.String(SolanaWSEndpoint2Flag.Name)

	solanaOS1, err := arrival.NewSolanaOrderbookStream(ctx, solanaRPCEndpoint, solanaWSEndpoint1, marketAddr)
	if err != nil {
		return err
	}
	solanaOS2, err := arrival.NewSolanaOrderbookStream(ctx, solanaRPCEndpoint, solanaWSEndpoint2, marketAddr)
	if err != nil {
		return err
	}

	startTime := time.Now()
	duration := c.Duration(DurationFlag.Name)
	runCtx, runCancel := context.WithTimeout(ctx, duration)
	defer runCancel()

	var (
		solanaUpdates1 []arrival.StreamUpdate[arrival.SolanaRawUpdate]
		solanaUpdates2 []arrival.StreamUpdate[arrival.SolanaRawUpdate]
	)
	errCh := make(chan error, 2)

	go func() {
		var err error

		solanaUpdates1, err = solanaOS1.Run(runCtx)
		if err != nil {
			errCh <- errors.Wrap(err, "could not collect results from solana 1")
			return
		}

		errCh <- nil
	}()

	go func() {
		var err error

		solanaUpdates2, err = solanaOS2.Run(runCtx)
		if err != nil {
			errCh <- errors.Wrap(err, "could not collect results from solana 2")
			return
		}

		errCh <- nil
	}()

	completionCount := 0
	ticker := time.NewTicker(updateInterval)

Loop:
	for {
		select {
		case runErr := <-errCh:
			completionCount++
			if runErr != nil {
				logger.Log().Errorw("fatal error during runtime: exiting", "err", err)
				return runErr
			}
			if completionCount == 2 {
				break Loop
			}
		case <-ticker.C:
			elapsedTime := time.Now().Sub(startTime).Round(time.Second)
			logger.Log().Infof("waited %v out of %v...", elapsedTime, duration)
		}

	}

	logger.Log().Infow("finished collecting data points", "solana1count", len(solanaUpdates1), "solana2count", len(solanaUpdates2))

	updateMap1 := make(map[gserum.Side]map[int]bool)
	for _, r := range solanaUpdates1 {
		slot := r.Data.Data.Context.Slot
		slots, ok := updateMap1[r.Data.Side]
		if !ok {
			slots = make(map[int]bool)
		}
		slots[int(slot)] = true

		updateMap1[r.Data.Side] = slots
	}

	bidSlots1 := maps.Keys(updateMap1[gserum.SideBid])
	sort.Ints(bidSlots1)

	askSlots1 := maps.Keys(updateMap1[gserum.SideAsk])
	sort.Ints(askSlots1)

	logger.Log().Infow("solana1 bid slot range", "low", bidSlots1[0], "high", bidSlots1[len(bidSlots1)-1])
	logger.Log().Infow("solana1 ask slot range", "low", askSlots1[0], "high", askSlots1[len(askSlots1)-1])

	updateMap2 := make(map[gserum.Side]map[int]bool)
	for _, r := range solanaUpdates2 {
		slot := r.Data.Data.Context.Slot
		slots, ok := updateMap2[r.Data.Side]
		if !ok {
			slots = make(map[int]bool)
		}
		slots[int(slot)] = true

		updateMap2[r.Data.Side] = slots
	}

	bidSlots2 := maps.Keys(updateMap2[gserum.SideBid])
	sort.Ints(bidSlots2)

	askSlots2 := maps.Keys(updateMap2[gserum.SideAsk])
	sort.Ints(askSlots2)

	logger.Log().Infow("solana2 bid slot range", "low", bidSlots2[0], "high", bidSlots2[len(bidSlots2)-1])
	logger.Log().Infow("solana2 ask slot range", "low", askSlots2[0], "high", askSlots2[len(askSlots2)-1])

	for slot := range updateMap1[gserum.SideBid] {
		_, ok := updateMap2[gserum.SideBid][slot]
		if !ok {
			logger.Log().Warnw("missing in solana2 bid", "slot", slot)
		}
	}

	for slot := range updateMap1[gserum.SideAsk] {
		_, ok := updateMap2[gserum.SideAsk][slot]
		if !ok {
			logger.Log().Warnw("missing in solana2 ask", "slot", slot)
		}
	}

	//removeDuplicates := c.Bool(RemoveDuplicatesFlag.Name)

	//serumResults, serumDuplicates, err := serumOS.Process(serumUpdates, removeDuplicates)
	//if err != nil {
	//	return errors.Wrap(err, "could not process serum updates")
	//}
	//logger.Log().Infow("serum duplicate updates", "count", len(serumDuplicates))
	//
	//solanaResults, solanaDuplicates, err := solanaOS.Process(solanaUpdates, removeDuplicates)
	//if err != nil {
	//	return errors.Wrap(err, "could not process solana results")
	//}
	//logger.Log().Infow("solana duplicate updates", "count", len(solanaDuplicates))
	//
	//slots := SlotRange(serumResults, solanaResults)
	//logger.Log().Infow("finished processing data points", "startSlot", slots[0], "endSlot", slots[len(slots)-1])
	//
	//datapoints, _, _, err := Merge(slots, serumResults, solanaResults)
	//if err != nil {
	//	return err
	//}
	//
	//logger.Log().Infow("completed merging: outputting data...")
	//
	//// dump results to stdout
	//removeUnmatched := c.Bool(RemoveUnmatchedFlag.Name)
	//PrintSummary(duration, serumEndpoint, solanaWSEndpoint, datapoints)
	//
	//// write results to csv
	//outputFile := c.String(utils.OutputFileFlag.Name)
	//header := []string{"slot", "diff", "seq", "serum-time", "solana-side", "solana-time"}
	//err = csv.Write(outputFile, header, datapoints, func(line []string) bool {
	//	if removeUnmatched {
	//		for _, col := range line {
	//			return col == "n/a"
	//		}
	//	}
	//	return false
	//})
	//if err != nil {
	//	return err
	//}

	return nil
}

var (
	SolanaWSEndpoint1Flag = &cli.StringFlag{
		Name:  "solana-ws-endpoint-1",
		Usage: "first Solana websocket connection endpoint",
	}
	SolanaWSEndpoint2Flag = &cli.StringFlag{
		Name:  "solana-ws-endpoint-2",
		Usage: "second Solana websocket connection endpoint",
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
