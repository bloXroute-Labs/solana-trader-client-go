package provider

import (
	"errors"
	"time"

	"github.com/bloXroute-Labs/serum-client-go/bxserum/transaction"
	"github.com/gagliardetto/solana-go"
)

const (
	MainnetSerumAPIHTTP = "https://virginia.solana.dex.blxrbdn.com"
	MainnetSerumAPIWS   = "wss://virginia.solana.dex.blxrbdn.com/ws"
	MainnetSerumAPIGRPC = "virginia.solana.dex.blxrbdn.com:9000"
	TestnetSerumAPIHTTP = "http://serum-nlb-53baf45ef9775263.elb.us-east-1.amazonaws.com/"
	//TestnetSerumAPIWS   = "ws://serum-nlb-53baf45ef9775263.elb.us-east-1.amazonaws.com/ws"
	TestnetSerumAPIWS   = "ws://44.203.225.191:1809/ws"
	TestnetSerumAPIGRPC = "serum-nlb-53baf45ef9775263.elb.us-east-1.amazonaws.com:9000"
	defaultRPCTimeout   = 7 * time.Second
)

var ErrPrivateKeyNotFound = errors.New("private key not provided for signing transaction")

type PostOrderOpts struct {
	OpenOrdersAddress string
	ClientOrderID     uint64
	SkipPreFlight     bool
}

type RPCOpts struct {
	Endpoint   string
	Timeout    time.Duration
	PrivateKey *solana.PrivateKey
}

func DefaultRPCOpts(endpoint string) RPCOpts {
	var spk *solana.PrivateKey
	privateKey, err := transaction.LoadPrivateKeyFromEnv()
	if err == nil {
		spk = &privateKey
	}
	return RPCOpts{
		Endpoint:   endpoint,
		Timeout:    defaultRPCTimeout,
		PrivateKey: spk,
	}
}
