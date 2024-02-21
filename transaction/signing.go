package transaction

import (
	"errors"
	"fmt"
	"github.com/gagliardetto/solana-go"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	"os"
)

// LoadPrivateKeyFromEnv looks up private key from the `PRIVATE_KEY` environment variable
func LoadPrivateKeyFromEnv() (solana.PrivateKey, error) {
	privateKeyBase58, ok := os.LookupEnv("PRIVATE_KEY")
	if !ok {
		return solana.PrivateKey{}, fmt.Errorf("env variable `PRIVATE_KEY` not set")
	}

	return solana.PrivateKeyFromBase58(privateKeyBase58)
}

// SignTx uses the environment variable for `PRIVATE_KEY` to sign the message content and replace the zero signature
func SignTx(unsignedTxBase64 string) (string, error) {
	privateKey, err := LoadPrivateKeyFromEnv()
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
		if signaturesRequired-signaturesPresent == 1 {
			return appendSignature(solanaTx, privateKey)
		}
		return fmt.Errorf("transaction requires %v signatures and has %v signatures", signaturesRequired, signaturesPresent)
	}

	return replaceZeroSignature(solanaTx, privateKey)
}

func appendSignature(solanaTx *solana.Transaction, privateKey solana.PrivateKey) error {
	messageContent, err := solanaTx.Message.MarshalBinary()
	if err != nil {
		return fmt.Errorf("unable to encode message for signing: %w", err)
	}

	signedMessageContent, err := privateKey.Sign(messageContent)
	if err != nil {
		return fmt.Errorf("unable to sign message: %v", err)
	}

	solanaTx.Signatures = append(solanaTx.Signatures, signedMessageContent)
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

	zeroSigIndex := -1
	for i, sig := range tx.Signatures {
		if sig.IsZero() {
			if zeroSigIndex != -1 {
				return errors.New("more than one zero signature provided in transaction")
			}
			zeroSigIndex = i
		}
	}
	if zeroSigIndex == -1 {
		return errors.New("no zero signatures to replace in transaction")
	}

	tx.Signatures[zeroSigIndex] = signedMessageContent
	return nil
}

// PartialSign heavily derived from `solana-go/transaction.go`. Signs the transaction with all available private keys, except
// the main Solana address's
func PartialSign(tx *solana.Transaction, ownerPk solana.PublicKey, privateKeys map[solana.PublicKey]solana.PrivateKey) error {
	requiredSignatures := tx.Message.Header.NumRequiredSignatures
	if uint8(len(privateKeys)) != requiredSignatures-1 {
		// one signature is reserved for the end user to sign the transaction
		return fmt.Errorf("unexpected error: could not generate enough signatures : # of privateKeys : %d vs requiredSignatures: %d", len(privateKeys), requiredSignatures)
	}

	messageBytes, err := tx.Message.MarshalBinary()
	if err != nil {
		return err
	}

	signatures := make([]solana.Signature, 0, requiredSignatures)
	for _, key := range tx.Message.AccountKeys[0:requiredSignatures] {
		if key == ownerPk {
			// if belongs to owner: add empty signature
			signatures = append(signatures, solana.Signature{})
		} else {
			// otherwise, sign
			privateKey, ok := privateKeys[key]
			if !ok {
				return errors.New("private key not found")
			}
			s, err := privateKey.Sign(messageBytes)
			if err != nil {
				return err // TODO: wrap error
			}
			signatures = append(signatures, s)
		}
	}

	tx.Signatures = signatures
	return nil
}
