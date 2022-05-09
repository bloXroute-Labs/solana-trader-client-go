package transaction

import (
	"github.com/gagliardetto/solana-go"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
)

func SignTx(transaction string, privateKeyKeygenPath string) (string, error) {
	txBytes, err := solanarpc.DataBytesOrJSONFromBase64(transaction)
	tx := solanarpc.TransactionWithMeta{
		Transaction: txBytes,
	}
	solanaTx, err := tx.GetTransaction()
	if err != nil {
		return "", err
	}
	privKey, err := solana.PrivateKeyFromSolanaKeygenFile(privateKeyKeygenPath)
	if err != nil {
		return "", err
	}

	err = signTx(solanaTx, privKey)
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
