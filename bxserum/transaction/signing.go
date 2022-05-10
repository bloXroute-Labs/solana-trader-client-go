package transaction

import (
	"encoding/base64"
	"fmt"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"os"
)

// SignTxMessage uses the environment variable for `PRIVATE_KEY` to sign the transaction
func SignTxMessage(txMessage string) (string, error) {
	pKeyStr := os.Getenv("PRIVATE_KEY")
	if pKeyStr == "" {
		return "", fmt.Errorf("env variable `PRIVATE_KEY` not set")
	}
	pKey, err := solana.PrivateKeyFromBase58(pKeyStr)
	if err != nil {
		return "", err
	}

	var m solana.Message
	txBytes, err := base64.StdEncoding.DecodeString(txMessage)
	if err := m.UnmarshalWithDecoder(bin.NewCompactU16Decoder(txBytes)); err != nil {
		return "", err
	}

	solanaTx := &solana.Transaction{Message: m}
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
