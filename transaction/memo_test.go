package transaction

import (
	"github.com/gagliardetto/solana-go"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddMemo(t *testing.T) {
	privateKey, err := solana.NewRandomPrivateKey()
	if err != nil {
		log.Fatal(err)
	}
	privateKeys := make(map[solana.PublicKey]solana.PrivateKey)
	privateKeys[privateKey.PublicKey()] = privateKey
	if err != nil {
		log.Fatal(err)
	}
	encodedTxn, err := AddMemo(
		[]solana.Instruction{},
		"new memo by dev",
		solana.MustHashFromBase58("A1xapHMk7Y9tj2NuVKw1ddKASsCce2M5EyD1xXo3RWr1"),
		solana.MustPublicKeyFromBase58(privateKey.PublicKey().String()),
		privateKeys,
	)
	require.NoError(t, err)
	require.NotEmpty(t, encodedTxn)
}

func TestAddMemoToSerializedTxn(t *testing.T) {
	privateKey, err := solana.NewRandomPrivateKey()
	if err != nil {
		log.Fatal(err)
	}
	privateKeys := make(map[solana.PublicKey]solana.PrivateKey)
	privateKeys[privateKey.PublicKey()] = privateKey
	if err != nil {
		log.Fatal(err)
	}
	encodedTxn, err := AddMemo(
		[]solana.Instruction{},
		"new memo by dev",
		solana.MustHashFromBase58("A1xapHMk7Y9tj2NuVKw1ddKASsCce2M5EyD1xXo3RWr1"),
		solana.MustPublicKeyFromBase58(privateKey.PublicKey().String()),
		privateKeys,
	)
	require.NoError(t, err)
	require.NotEmpty(t, encodedTxn)
	if err != nil {
		log.Fatal(err)
	}
	owner := privateKey.PublicKey()
	encodedTxn2, err := AddMemoToSerializedTxn(encodedTxn, "my memo", owner, privateKeys)
	require.NoError(t, err)
	require.NotEmpty(t, encodedTxn2)
}
