package main

import (
	"github.com/urfave/cli/v2"
	"time"
)

var (
	MintFlag = &cli.StringFlag{
		Name:  "mint",
		Usage: "mint to fetch price for (inactive token is best)",
		Value: "6D7nXHAhsRbwj8KFZR2agB6GEjMLg4BM7MAqZzRT8F1j", // gosu
	}

	TriggerActivityFlag = &cli.BoolFlag{
		Name:  "trigger-activity",
		Usage: "if true, send trigger transactions to force quote updates (requires PRIVATE_KEY environment variable_",
	}

	IterationsFlag = &cli.IntFlag{
		Name:  "iterations",
		Usage: "number of quotes to compare",
		Value: 5,
	}

	MaxRuntimeFlag = &cli.DurationFlag{
		Name:  "runtime",
		Usage: "max time to run benchmark for",
		Value: 5 * time.Minute,
	}

	SwapAmountFlag = &cli.Float64Flag{
		Name:  "swap-amount",
		Usage: "amount to swap for each trigger transaction (for unit, see --swap-mint)",
		Value: 1,
	}

	SwapMintFlag = &cli.StringFlag{
		Name:  "swap-mint",
		Usage: "corresponding token to swap from to --mint token for triggers",
		Value: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
	}

	SwapIntervalFlag = &cli.DurationFlag{
		Name:  "swap-interval",
		Usage: "time to wait between each swap",
		Value: 3 * time.Second,
	}

	SwapInitialWaitFlag = &cli.DurationFlag{
		Name:  "swap-initial-wait",
		Usage: "initial wait before beginning swaps",
		Value: 10 * time.Second,
	}

	SwapAfterWaitFlag = &cli.DurationFlag{
		Name:  "swap-after-wait",
		Usage: "initial wait after finishing swaps",
		Value: 10 * time.Second,
	}

	QueryIntervalFlag = &cli.DurationFlag{
		Name:  "query-interval",
		Usage: "time to wait between each poll for poll-based streams",
		Value: 500 * time.Millisecond,
	}

	PublicKeyFlag = &cli.StringFlag{
		Name:  "public-key",
		Usage: "public key to place swaps over (requires PRIVATE_KEY environment variable)",
		Value: "AFT8VayE7qr8MoQsW3wHsDS83HhEvhGWdbNSHRKeUDfQ",
	}

	EnvFlag = &cli.StringFlag{
		Name:  "env",
		Usage: "trader API environment (options: mainnet, testnet, devnet)",
		Value: "mainnet",
	}
)
