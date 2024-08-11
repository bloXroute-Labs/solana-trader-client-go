package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Example struct {
	Env           Env
	RunSlowStream bool
	RunTrades     bool
	RunPumpFun    bool
}

func BoolPtr(val bool) *bool {
	return &val
}

func Load() (Example, error) {
	env, err := loadEnv()
	if err != nil {
		return Example{}, err
	}

	runSlowStream := true
	rtsV := os.Getenv("RUN_SLOW_STREAM")
	if rtsV == "false" {
		runSlowStream = false
	}

	runTrades := true
	rtV := os.Getenv("RUN_TRADES")
	if rtV == "false" {
		runTrades = false
	}

	runPumpFun := true
	runPumpValue := os.Getenv("RUN_PUMP_FUN")
	if runPumpValue == "false" {
		runPumpFun = false
	}

	return Example{
		Env:           env,
		RunTrades:     runTrades,
		RunSlowStream: runSlowStream,
		RunPumpFun:    runPumpFun,
	}, nil
}

type Env string

const (
	EnvMainnet Env = "mainnet"
	EnvTestnet Env = "testnet"
	EnvLocal   Env = "local"
)

func loadEnv() (Env, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	v, ok := os.LookupEnv("API_ENV")
	if !ok {
		return EnvMainnet, nil
	}

	switch Env(strings.ToLower(v)) {
	case EnvLocal:
		return EnvLocal, nil
	case EnvTestnet:
		return EnvTestnet, nil
	case EnvMainnet:
		return EnvMainnet, nil
	default:
		return EnvMainnet, fmt.Errorf("API_ENV %v not supported", v)
	}
}
