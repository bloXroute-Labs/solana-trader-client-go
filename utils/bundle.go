package utils

import (
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
)

const (
	// BloxrouteTipAddress is from here and may fall out of date from time to time. Check our docs:
	// https://docs.bloxroute.com/solana/trader-api-v2/front-running-protection-and-transaction-bundle
	BloxrouteTipAddress = "HWEoBxYs7ssKuudEjzjmpfJVX7Dvi7wescFsVx2L5yoY"
)

// CreateBloxrouteTipTransactionToUseBundles creates a transaction you can use to when using PostSubmitBundle endpoints.
// This transaction should be the LAST transaction in your submission bundle
func CreateBloxrouteTipTransactionToUseBundles(privateKey solana.PrivateKey, tipAmount uint64, recentBlockHash solana.Hash) (*solana.Transaction, error) {
	recipient := solana.MustPublicKeyFromBase58(BloxrouteTipAddress)

	tx, err := solana.NewTransaction([]solana.Instruction{
		system.NewTransferInstruction(tipAmount, privateKey.PublicKey(), recipient).Build()}, recentBlockHash)
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
