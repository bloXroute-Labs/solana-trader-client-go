package transaction

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
)

const ReceipientAddress = "EpC5oyb2kx848Skd6guerBfHEyTTSSPZFFdWEhxtLxyt"

func CreateSampleTx(privateKey solana.PrivateKey, recentBlockHash solana.Hash) (*solana.Transaction, error) {

	recipient := solana.MustPublicKeyFromBase58(ReceipientAddress)

	tx, err := solana.NewTransaction([]solana.Instruction{
		system.NewTransferInstruction(1, privateKey.PublicKey(), recipient).Build(),
	}, recentBlockHash)
	if err != nil {
		return nil, err
	}

	signatures, err := tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(privateKey.PublicKey()) {
			return &privateKey
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	tx.Signatures = signatures

	return tx, nil
}
