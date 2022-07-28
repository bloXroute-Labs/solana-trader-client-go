package main

import (
	"context"
	"errors"
	"github.com/bloXroute-Labs/serum-client-go/benchmark/internal/csv"
	"github.com/bloXroute-Labs/serum-client-go/benchmark/internal/logger"
	"github.com/bloXroute-Labs/serum-client-go/benchmark/internal/transaction"
	"github.com/bloXroute-Labs/serum-client-go/benchmark/internal/utils"
	"github.com/bloXroute-Labs/serum-client-go/bxserum/provider"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "benchmark-txcompare",
		Usage: "Compares submitting transactions to multiple Solana nodes to determine if one is consistently faster",
		Flags: []cli.Flag{
			IterationCountFlag,
			SolanaHTTPEndpointsFlag,
			SolanaQueryEndpointsFlag,
			utils.OutputFileFlag,
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

	opts := provider.DefaultRPCOpts(provider.MainnetSerumAPIGRPC)
	if opts.PrivateKey == nil {
		return errors.New("PRIVATE_KEY environment variable must be set")
	}

	iterations := c.Int(IterationCountFlag.Name)
	endpoints := c.StringSlice(SolanaHTTPEndpointsFlag.Name)
	queryEndpoint := c.String(SolanaQueryEndpointsFlag.Name)

	querier := transaction.NewStatusQuerier(queryEndpoint)

	recentBlockHash, err := querier.RecentBlockHash(ctx)
	if err != nil {
		return err
	}
	submitter := transaction.NewSubmitter(endpoints, transaction.MemoBuilder(*opts.PrivateKey, recentBlockHash))

	signatures, creationTimes, err := submitter.SubmitIterations(ctx, iterations)
	if err != nil {
		return err
	}

	datapoints := make([]Datapoint, 0)
	best := make([]int, len(endpoints))
	lost := make([]int, len(endpoints))
	for i, iterationSignatures := range signatures {
		summary, statuses, err := querier.FetchBatch(ctx, iterationSignatures)
		if err != nil {
			return err
		}

		if summary.Best >= 0 {
			logger.Log().Debugw("iteration results found", "iteration", i, "winner", endpoints[summary.Best])
			best[summary.Best]++
		} else {
			logger.Log().Debugw("iteration no transactions confirmed", "iteration", i)
		}
		for j, status := range statuses {
			dp := Datapoint{
				Iteration:     i,
				CreationTime:  creationTimes[i],
				Signature:     iterationSignatures[j].String(),
				Endpoint:      endpoints[j],
				Executed:      status.Found,
				ExecutionTime: status.ExecutionTime,
				Slot:          status.Slot,
				Position:      status.Position,
			}
			datapoints = append(datapoints, dp)

			if !status.Found {
				lost[j]++
			}

			logger.Log().Debugw("iteration transaction result", "iteration", i, "endpoint", dp.Endpoint, "slot", dp.Slot, "position", dp.Position, "signature", dp.Signature)
		}
	}

	Print(iterations, endpoints, best, lost)

	outputFile := c.String(utils.OutputFileFlag.Name)
	header := []string{"iteration", "creation-time", "signature", "endpoint", "executed", "execution-time", "slot", "position"}
	err = csv.Write(outputFile, header, datapoints, func(line []string) bool {
		return false
	})
	if err != nil {
		return err
	}

	return nil
}

var (
	IterationCountFlag = &cli.IntFlag{
		Name:     "iterations",
		Aliases:  []string{"n"},
		Usage:    "number of transaction pairs to submit",
		Required: true,
	}
	SolanaHTTPEndpointsFlag = &cli.StringSliceFlag{
		Name:     "endpoint",
		Aliases:  []string{"e"},
		Usage:    "solana endpoints to submit transactions to (multiple allowed)",
		Required: true,
	}
	SolanaQueryEndpointsFlag = &cli.StringFlag{
		Name:     "query-endpoint",
		Usage:    "solana endpoints to query for transaction inclusion (useful when submission endpoint doesn't index transactions)",
		Required: true,
	}
)
