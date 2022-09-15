package provider

import (
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
)

func TestWS_New(t *testing.T) {
	pk, err := solana.NewRandomPrivateKey()
	assert.NotNil(t, pk)
	assert.Nil(t, err)

	err = os.Setenv("PRIVATE_KEY", pk.String())
	require.Nil(t, err)

	c, err := provider.NewWSClient()
	assert.NotNil(t, c)
	assert.Nil(t, err)

	err = os.Unsetenv("PRIVATE_KEY")
	require.Nil(t, err)
}

func TestWS_NewWithOpts(t *testing.T) {
	pk, err := solana.NewRandomPrivateKey()
	assert.NotNil(t, pk)
	assert.Nil(t, err)

	o := provider.RPCOpts{
		Endpoint:   provider.MainnetWS,
		Timeout:    time.Second,
		PrivateKey: &pk,
	}

	c, err := provider.NewWSClientWithOpts(o)
	assert.NotNil(t, c)
	assert.Nil(t, err)
}
