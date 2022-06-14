package provider

import (
	"testing"
	"time"

	"github.com/bloXroute-Labs/serum-api/bxserum/provider"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
)

func TestGRPC_New(t *testing.T) {
	c, err := provider.NewGRPCClient()
	assert.NotNil(t, c)
	assert.Nil(t, err)
}

func TestGRPC_NewWithOpts(t *testing.T) {
	pk, err := solana.NewRandomPrivateKey()
	assert.NotNil(t, pk)
	assert.Nil(t, err)

	o := provider.RPCOpts{
		Endpoint:   provider.MainnetSerumAPIGRPC,
		Timeout:    time.Second,
		PrivateKey: pk,
	}
	c, err := provider.NewGRPCClientWithOpts(o)
	assert.NotNil(t, c)
	assert.Nil(t, err)
}
