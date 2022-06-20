package provider

import (
	"os"
	"testing"
	"time"

	"github.com/bloXroute-Labs/serum-api/bxserum/provider"
	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
)

func TestWS_New(t *testing.T) {
	pk, err := solana.NewRandomPrivateKey()
	assert.NotNil(t, pk)
	assert.Nil(t, err)

	os.Setenv("PRIVATE_KEY", pk.String())

	c, err := provider.NewWSClient()
	assert.NotNil(t, c)
	assert.Nil(t, err)

	os.Unsetenv("PRIVATE_KEY")
}

func TestWS_NewWithOpts(t *testing.T) {
	pk, err := solana.NewRandomPrivateKey()
	assert.NotNil(t, pk)
	assert.Nil(t, err)

	o := provider.RPCOpts{
		Endpoint:   provider.MainnetSerumAPIWS,
		Timeout:    time.Second,
		PrivateKey: &pk,
	}

	c, err := provider.NewWSClientWithOpts(o)
	assert.NotNil(t, c)
	assert.Nil(t, err)
}
