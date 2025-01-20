package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Example struct {
	Env           Env
	RunSlowStream bool
	RunTrades     bool
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

	return Example{
		Env:           env,
		RunTrades:     runTrades,
		RunSlowStream: runSlowStream,
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

type EnvironmentVariables struct {
	PrivateKey        string
	PublicKey         string
	OpenOrdersAddress string
	Payer             string
}

func InitializeEnvironmentVariables() EnvironmentVariables {
	// Load .env file if it exists
	godotenv.Load()

	// Check required AUTH_HEADER
	if os.Getenv("AUTH_HEADER") == "" {
		log.Fatal("must specify bloXroute authorization header!")
	}

	// Get private key
	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		log.Error("PRIVATE_KEY environment variable not set, cannot run any examples that require tx submission")
	}

	// Get public key
	publicKey := os.Getenv("PUBLIC_KEY")
	if publicKey == "" {
		log.Warn("PUBLIC_KEY environment variable not set: will skip place/cancel/settle examples")
	}

	// Get open orders
	openOrders := os.Getenv("OPEN_ORDERS")
	if openOrders == "" {
		log.Error("OPEN_ORDERS environment variable not set: requests will be slower")
	}

	// Get payer or default to public key
	payer := os.Getenv("PAYER")
	if payer == "" {
		log.Warn("PAYER environment variable not set: will be set to owner address")
		payer = publicKey
	}

	return EnvironmentVariables{
		PrivateKey:        privateKey,
		PublicKey:         publicKey,
		OpenOrdersAddress: openOrders,
		Payer:             payer,
	}
}
