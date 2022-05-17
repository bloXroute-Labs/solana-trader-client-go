package provider

import (
	"github.com/bloXroute-Labs/serum-api/bxserum/transaction"
	"github.com/gagliardetto/solana-go"
	"time"
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
		Timeout:    time.Second * 7,
		PrivateKey: privateKey,
	}, nil
}
