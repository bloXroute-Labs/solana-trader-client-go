package transaction

import (
	"fmt"
	"github.com/gagliardetto/solana-go"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	"os"
)

// SignTxMessage uses the environment variable for `PRIVATE_KEY` to sign the transaction message
func SignTxMessage(txMessage string) (string, error) {
	pKeyStr := os.Getenv("PRIVATE_KEY")
	if pKeyStr == "" {
		return "", fmt.Errorf("env variable `PRIVATE_KEY` not set")
	}
	pKey, err := solana.PrivateKeyFromBase58(pKeyStr)
	if err != nil {
		return "", err
	}

	txBytes, err := solanarpc.DataBytesOrJSONFromBase64(txMessage)
	if err != nil {
		return "", err
	}
	tx := solanarpc.TransactionWithMeta{Transaction: txBytes}
	solanaTx, err := tx.GetTransaction()

	err = signTx(solanaTx, pKey)
	if err != nil {
		return "", err
	}

	return solanaTx.ToBase64()
}

func signTx(tx *solana.Transaction, privateKey solana.PrivateKey) error {
	signaturesRequired := int(tx.Message.Header.NumRequiredSignatures)
	signaturesPresent := len(tx.Signatures)
	if signaturesPresent != signaturesRequired-1 {
		return fmt.Errorf("transaction requires %v instruction and has %v signatures, should need exactly one more signature", signaturesRequired, signaturesPresent)
	}

	_, err := tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(privateKey.PublicKey()) {
			return &privateKey
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
