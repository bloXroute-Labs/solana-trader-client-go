package provider

import (
	"github.com/bloXroute-Labs/serum-api/bxserum/transaction"
	"github.com/gagliardetto/solana-go"
	"time"
)

const (
	MainnetSerumAPIHTTP = "https://virginia.solana.dex.blxrbdn.com"
	MainnetSerumAPIWS   = "wss://virginia.solana.dex.blxrbdn.com/ws"
	MainnetSerumAPIGRPC = "107.22.159.156:1810" // TODO: setup load balancer
	TestnetSerumAPIHTTP = "http://44.200.244.204:80"
	TestnetSerumAPIWS   = "ws://44.200.244.204:80/ws"
	TestnetSerumAPIGRPC = "44.200.244.204:443"
	defaultRPCTimeout   = 7 * time.Second
)

type PostOrderOpts struct {
	OpenOrdersAddress string
	ClientOrderID     uint64
}

type RPCOpts struct {
	Endpoint   string
	Timeout    time.Duration
	PrivateKey solana.PrivateKey
}

func DefaultRPCOpts(endpoint string) (RPCOpts, error) {
	privateKey, err := transaction.LoadPrivateKeyFromEnv()
	if err != nil {
		return RPCOpts{}, err
	}
	return RPCOpts{
		Endpoint:   endpoint,
		Timeout:    defaultRPCTimeout,
		PrivateKey: privateKey,
	}, nil
}
