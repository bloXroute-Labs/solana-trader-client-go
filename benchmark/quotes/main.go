package main

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/logger"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/stream"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/utils"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "benchmark-quotes",
		Usage: "Compares Solana Trader API AMM quotes with Jupiter API",
		Flags: []cli.Flag{
			utils.APIWSEndpoint,
			// utils.OutputFileFlag,
			InputMintFlag,
			OutputMintFlag,
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
	inputMint := c.String(InputMintFlag.Name)
	outputMint := c.String(OutputMintFlag.Name)
	// iterations := c.Int(IterationsFlag.Name)

	jupiterApi, err := stream.NewJupiterAPI(ctx, stream.WithJupiterTokenPair(inputMint, outputMint))
	if err != nil {
		return err
	}

	var (
		// tradeUpdates   []stream.RawUpdate[[]byte]
		jupiterUpdates []stream.RawUpdate[stream.JupiterResponse]
		errCh          = make(chan error, 2)
	)

	go func() {
		var err error

		jupiterUpdates, err = jupiterApi.Run(ctx)
		if err != nil {
			errCh <- errors.Wrap(err, "could not collect results from Solana")
			return
		}

		errCh <- nil
	}()

	completionCount := 0

	for completionCount < 1 {
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

	return nil
}

var (
	InputMintFlag = &cli.StringFlag{
		Name:  "input-mint",
		Usage: "input mint to fetch quote for (inactive pairs are best)",
		Value: "",
	}

	OutputMintFlag = &cli.StringFlag{
		Name:  "output-mint",
		Usage: "output mint to fetch quote for (inactive pairs are best)",
		Value: "",
	}

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
