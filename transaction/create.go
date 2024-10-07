package transaction

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
)

const ReceipientAddress = "5wiGAqf4BX23XU6jc3MDDZAFoNV5pz61thsUuSgpsAxS"

func CreateSampleTx(privateKey solana.PrivateKey, recentBlockHash solana.Hash, lamports uint64) (*solana.Transaction, error) {

	recipient := solana.MustPublicKeyFromBase58(ReceipientAddress)

	tx, err := solana.NewTransaction([]solana.Instruction{
		system.NewTransferInstruction(lamports, privateKey.PublicKey(), recipient).Build(),
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
