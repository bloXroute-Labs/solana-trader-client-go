package main

import (
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
)

func traderClients(env string) (*provider.HTTPClient, *provider.WSClient, error) {
	var (
		httpClient *provider.HTTPClient
		wsClient   *provider.WSClient
		err        error
	)
	switch env {
	case "testnet":
		httpClient = provider.NewHTTPTestnet()
		wsClient, err = provider.NewWSClientTestnet()
		if err != nil {
			return nil, nil, err
		}
	case "devnet":
		httpClient = provider.NewHTTPDevnet()
		wsClient, err = provider.NewWSClientDevnet()
		if err != nil {
			return nil, nil, err
		}
	case "devnet1":
		httpClient = provider.NewHTTPClientWithOpts(nil, provider.DefaultRPCOpts("http://3.239.217.218:1809"))
		wsClient, err = provider.NewWSClientWithOpts(provider.DefaultRPCOpts("ws://3.239.217.218:1809/ws"))
		if err != nil {
			return nil, nil, err
		}
	case "mainnet":
		httpClient = provider.NewHTTPClient()
		wsClient, err = provider.NewWSClient()
		if err != nil {
			return nil, nil, err
		}
	case "mainnet1":
		httpClient = provider.NewHTTPClientWithOpts(nil, provider.DefaultRPCOpts("http://54.161.46.25:1809"))
		wsClient, err = provider.NewWSClientWithOpts(provider.DefaultRPCOpts("ws://54.161.46.25:1809/ws"))
		if err != nil {
			return nil, nil, err
		}
	default:
		return nil, nil, fmt.Errorf("unknown environment: %v", env)
	}

	return httpClient, wsClient, nil
}
