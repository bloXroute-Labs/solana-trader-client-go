package transaction

import (
	"fmt"
	"github.com/gagliardetto/solana-go"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	"os"
)

// SignTx uses the environment variable for `PRIVATE_KEY` to sign the message content and replace the zero signature
func SignTx(unsignedTxBase64 string) (string, error) {
	privateKeyBase58, ok := os.LookupEnv("PRIVATE_KEY")
	if !ok {
		return "", fmt.Errorf("env variable `PRIVATE_KEY` not set")
	}

	privateKey, err := solana.PrivateKeyFromBase58(privateKeyBase58)
	if err != nil {
		return "", err
	}

	return SignTxWithPrivateKey(unsignedTxBase64, privateKey)
}

// SignTxWithPrivateKey uses the provided private key to sign the message content and replace the zero signature
func SignTxWithPrivateKey(unsignedTxBase64 string, privateKey solana.PrivateKey) (string, error) {
	unsignedTxBytes, err := solanarpc.DataBytesOrJSONFromBase64(unsignedTxBase64)
	if err != nil {
		return "", err
	}
	unsignedTx := solanarpc.TransactionWithMeta{Transaction: unsignedTxBytes}
	solanaTx, err := unsignedTx.GetTransaction()
	if err != nil {
		return "", err
	}

	err = signTx(solanaTx, privateKey)
	if err != nil {
		return "", err
	}

	return solanaTx.ToBase64()
}

func signTx(solanaTx *solana.Transaction, privateKey solana.PrivateKey) error {
	signaturesRequired := int(solanaTx.Message.Header.NumRequiredSignatures)
	signaturesPresent := len(solanaTx.Signatures)
	if signaturesPresent != signaturesRequired {
		return fmt.Errorf("transaction requires %v signatures and has %v signatures", signaturesRequired, signaturesPresent)
	}

	err := replaceZeroSignature(solanaTx, privateKey)
	if err != nil {
		return err
	}
	return nil
}

func replaceZeroSignature(tx *solana.Transaction, privateKey solana.PrivateKey) error {
	messageContent, err := tx.Message.MarshalBinary()
	if err != nil {
		return fmt.Errorf("unable to encode message for signing: %w", err)
	}

	signedMessageContent, err := privateKey.Sign(messageContent)
	if err != nil {
		return fmt.Errorf("unable to sign message: %v", err)
	}

	for i, sig := range tx.Signatures {
		if sig.IsZero() {
			tx.Signatures[i] = signedMessageContent
			break
		}
	}

	return nil
}
