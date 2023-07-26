package transaction

import (
	"github.com/gagliardetto/solana-go"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
)

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
	txbuilder := solana.NewTransactionBuilder()
	txbuilder.AddInstruction(&solana.GenericInstruction{
		AccountValues: nil,
		ProgID:        solana.PublicKey{},
		DataBytes:     nil,
	})
	txbuilder.SetRecentBlockHash(solana.MustHashFromBase58("A1xapHMk7Y9tj2NuVKw1ddKASsCce2M5EyD1xXo3RWr1"))
	txbuilder.SetFeePayer(privateKey.PublicKey())
	tx, err := txbuilder.Build()
	if err != nil {
		log.Fatal(err)
	}

	encodedTxn := tx.MustToBase64()
	require.NoError(t, err)
	require.NotEmpty(t, encodedTxn)
	if err != nil {
		log.Fatal(err)
	}
	encodedTxn2, err := AddMemoAndSign(encodedTxn, privateKey)
	require.NoError(t, err)
	require.NotEmpty(t, encodedTxn2)

	// validate
	signedTxBytes, err := solanarpc.DataBytesOrJSONFromBase64(encodedTxn2)
	if err != nil {
		log.Fatal(err)
	}
	unsignedTx := solanarpc.TransactionWithMeta{Transaction: signedTxBytes}
	solanaTx, err := unsignedTx.GetTransaction()

	require.Equal(t, 2, len(solanaTx.Message.Instructions))
	program, err := solanaTx.Message.Program(solanaTx.Message.Instructions[1].ProgramIDIndex)
	if err != nil {
		log.Fatal(err)
	}
	require.Equal(t, TraderAPIMemoProgram, program)

}
