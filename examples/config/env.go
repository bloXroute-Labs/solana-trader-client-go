package config

import (
	"fmt"
	"os"
	"strings"
)

type Env string

const (
	EnvMainnet Env = "mainnet"
	EnvTestnet Env = "testnet"
)

func LoadEnv() (Env, error) {
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
