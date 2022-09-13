package provider

import (
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	"os"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
)

func TestHTTP_New(t *testing.T) {
	pk, err := solana.NewRandomPrivateKey()
	assert.NotNil(t, pk)
	assert.Nil(t, err)

	os.Setenv("PRIVATE_KEY", pk.String())

	c := provider.NewHTTPClient()
	assert.NotNil(t, c)
	assert.Nil(t, err)

	os.Unsetenv("PRIVATE_KEY")
}

func TestHTTP_NewWithOpts(t *testing.T) {
	pk, err := solana.NewRandomPrivateKey()
	assert.NotNil(t, pk)
	assert.Nil(t, err)

	o := provider.RPCOpts{
		Endpoint:   provider.MainnetHTTP,
		Timeout:    time.Second,
		PrivateKey: &pk,
	}

	c := provider.NewHTTPClientWithOpts(nil, o)
	assert.NotNil(t, c)
}
