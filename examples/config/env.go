package config

import (
	"fmt"
	"os"
	"strings"
)

type Example struct {
	Env            Env
	RunTradeStream bool
	RunTrades      bool
}

func Load() (Example, error) {
	env, err := loadEnv()
	if err != nil {
		return Example{}, err
	}

	runTradesStream := true
	rtsV := os.Getenv("RUN_TRADE_STREAM")
	if rtsV == "false" {
		runTradesStream = false
	}

	runTrades := true
	rtV := os.Getenv("RUN_TRADES")
	if rtV == "false" {
		runTrades = false
	}

	return Example{
		Env:            env,
		RunTrades:      runTrades,
		RunTradeStream: runTradesStream,
	}, nil
}

type Env string

const (
	EnvMainnet Env = "mainnet"
	EnvTestnet Env = "testnet"
)

func loadEnv() (Env, error) {
	v, ok := os.LookupEnv("API_ENV")
	if !ok {
		return EnvMainnet, nil
	}

	switch Env(strings.ToLower(v)) {
	case EnvTestnet:
		return EnvTestnet, nil
	case EnvMainnet:
		return EnvMainnet, nil
	default:
		return EnvMainnet, fmt.Errorf("API_ENV %v not supported", v)
	}
}
