package transaction

import (
	"fmt"
	"github.com/gagliardetto/solana-go"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	"os"
)

// SignTx uses the environment variable for `PRIVATE_KEY` to sign the transaction
func SignTx(transaction string) (string, error) {
	pKeyStr := os.Getenv("PRIVATE_KEY")
	if pKeyStr == "" {
		return "", fmt.Errorf("env variable `PRIVATE_KEY` not set")
	}
	pKey, err := solana.PrivateKeyFromBase58(pKeyStr)
	if err != nil {
		return "", err
	}

	txBytes, err := solanarpc.DataBytesOrJSONFromBase64(transaction)
	tx := solanarpc.TransactionWithMeta{
		Transaction: txBytes,
	}
	solanaTx, err := tx.GetTransaction()
	if err != nil {
		return "", err
	}

	err = signTx(solanaTx, pKey)
	if err != nil {
		return "", err
	}

	return solanaTx.ToBase64()
}

func signTx(tx *solana.Transaction, privateKey solana.PrivateKey) error {
	sigs, err := tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(privateKey.PublicKey()) {
			return &privateKey
		}
		return nil
	})
	if err != nil {
		return err
	}

	tx.Signatures = sigs
	return nil
}
