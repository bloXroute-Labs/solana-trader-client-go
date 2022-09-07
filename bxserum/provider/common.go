package provider

import (
	"errors"
	"os"
	"time"

	"github.com/bloXroute-Labs/serum-client-go/bxserum/transaction"
	"github.com/gagliardetto/solana-go"
)

const (
	MainnetSerumAPIHTTP = "https://virginia.solana.dex.blxrbdn.com"
	MainnetSerumAPIWS   = "wss://virginia.solana.dex.blxrbdn.com/ws"
	MainnetSerumAPIGRPC = "virginia.solana.dex.blxrbdn.com:443"
	TestnetSerumAPIHTTP = "http://serum-nlb-5a2c3912804344a3.elb.us-east-1.amazonaws.com/"
	TestnetSerumAPIWS   = "ws://serum-nlb-5a2c3912804344a3.elb.us-east-1.amazonaws.com/ws"
	TestnetSerumAPIGRPC = "serum-nlb-5a2c3912804344a3.elb.us-east-1.amazonaws.com:80"
	DevnetSerumAPIHTTP  = "http://serum-nlb-53baf45ef9775263.elb.us-east-1.amazonaws.com/"
	DevnetSerumAPIWS    = "ws://serum-nlb-53baf45ef9775263.elb.us-east-1.amazonaws.com/ws"
	DevnetSerumAPIGRPC  = "serum-nlb-53baf45ef9775263.elb.us-east-1.amazonaws.com:80"
	LocalSerumAPIWS     = "ws://localhost:9000/ws"
	LocalSerumAPIGRPC   = "localhost:9000"
	LocalSerumAPIHTTP   = "http://127.0.0.1:9000"
	defaultRPCTimeout   = 7 * time.Second
)

var ErrPrivateKeyNotFound = errors.New("private key not provided for signing transaction")

type PostOrderOpts struct {
	OpenOrdersAddress string
	ClientOrderID     uint64
	SkipPreFlight     bool
}

type RPCOpts struct {
	Endpoint       string
	Timeout        time.Duration
	PrivateKey     *solana.PrivateKey
	AuthHeader     string
	CacheBlockHash time.Duration
}

func DefaultRPCOpts(endpoint string) RPCOpts {
	var spk *solana.PrivateKey
	privateKey, err := transaction.LoadPrivateKeyFromEnv()
	if err == nil {
		spk = &privateKey
	}
	return RPCOpts{
		Endpoint:       endpoint,
		Timeout:        defaultRPCTimeout,
		PrivateKey:     spk,
		AuthHeader:     os.Getenv("AUTH_HEADER"),
		CacheBlockHash: time.Second * 45, // to be on safe side
	}
}
