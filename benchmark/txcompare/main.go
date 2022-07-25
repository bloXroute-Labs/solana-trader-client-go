package main

import (
	"context"
	"fmt"
	transaction2 "github.com/bloXroute-Labs/serum-client-go/benchmark/internal/transaction"
	"github.com/bloXroute-Labs/serum-client-go/bxserum/provider"
	"github.com/bloXroute-Labs/serum-client-go/examples/benchmark/internal/logger"
	"github.com/bloXroute-Labs/serum-client-go/examples/benchmark/internal/utils"
	"github.com/gagliardetto/solana-go"
	"github.com/urfave/cli/v2"
	"os"
)

/*

benchmark/txcompare/main.go

PRIVATE_KEY=...
OPEN_ORDERS=...

*/

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
			err := logger.Log().Sync()
			if err != nil {
				fmt.Println("error flushing log, may not have logged all messages: ", err)
			}
		}
	}()

	if err != nil {
		panic(err)
	}
}

func run(c *cli.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ooAddress := os.Getenv("OPEN_ORDERS")
	ooPk, err := solana.PublicKeyFromBase58(ooAddress)
	if err != nil {
		return err
	}
	opts := provider.DefaultRPCOpts(provider.MainnetSerumAPIGRPC)
	publicKey := opts.PrivateKey.PublicKey()
	g, err := provider.NewGRPCClientWithOpts(opts)

	iterations := c.Int(IterationCountFlag.Name)
	endpoints := c.StringSlice(SolanaHTTPEndpointsFlag.Name)
	queryEndpoint := c.String(SolanaQueryEndpointsFlag.Name)

	submitter := transaction2.NewSubmitter(endpoints, transaction2.SerumBuilder(ctx, g, publicKey, ooPk, *opts.PrivateKey))
	querier := transaction2.NewStatusQuerier(queryEndpoint)

	signatures, err := submitter.SubmitIterations(ctx, iterations)
	if err != nil {
		return err
	}

	for i, iterationSignatures := range signatures {
		summary, statuses, err := querier.FetchBatch(ctx, iterationSignatures)
		if err != nil {
			return err
		}

		logger.Log().Debugw("iteration results found", "iteration", i, "winner", endpoints[summary.Best])
		for j, status := range statuses {
			logger.Log().Debugw("iteration transaction result", "iteration", i, "slot", status.Slot, "position", status.Position, "signature", iterationSignatures[j])
		}
	}

	for i := 0; i < iterations; i++ {
		signatures, err := submitter.SubmitIteration(ctx)
		if err != nil {
			return err
		}

		logger.Log().Debugw("submitted iteration of transactions", "iteration", i, "count", len(signatures))

		for j, signature := range signatures {
			status, err := querier.Fetch(ctx, signature)
			if err != nil {
				return err
			}

			logger.Log().Infow("found submitted transaction", "iteration", i, "endpoint", endpoints[j], "status", status, "signature", signature)
		}
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
