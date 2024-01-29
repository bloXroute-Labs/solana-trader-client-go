package utils

import (
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
)

const (
	BloxrouteTipAddress = "AFT8VayE7qr8MoQsW3wHsDS83HhEvhGWdbNSHRKeUDfQ"
	jitoTipAddress      = "96gYZGLnJYVFmbjzopPSU6QiEV5fGqZNyN9nmNhvrZU5"
	memoProgramIDV2     = "MemoSq4gqABAXKb96qnH8TysNcWxMyWCqXgDLGmfcHr"
)

func CreateJitoTipTx(privateKey solana.PrivateKey, recentBlockHash solana.Hash) (*solana.Transaction, error) {
	// https://jito-labs.gitbook.io/mev/searcher-resources/bundles
	// as of 1/17/2024, all bundle requests must include a transaction at the end of the bundle tipping 1000 lamports
	// to jitoTipAddress

	recipient := solana.MustPublicKeyFromBase58(jitoTipAddress)

	tx, err := solana.NewTransaction([]solana.Instruction{
		system.NewTransferInstruction(1000, privateKey.PublicKey(), recipient).Build(),
		transaction.CreateTraderAPIMemoInstruction(transaction.BxMemoMarkerMsg)}, recentBlockHash)
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

// CreateBloxrouteTipTransactionToUseBundles creates a transaction you can use to when using PostSubmitJitoBundle endpoints.
// This transaction should be the LAST transaction in your submission bundle
func CreateBloxrouteTipTransactionToUseBundles(privateKey solana.PrivateKey, tipAmount uint64, recentBlockHash solana.Hash) (*solana.Transaction, error) {
	// https://jito-labs.gitbook.io/mev/searcher-resources/bundles
	// as of 1/17/2024, all bundle requests must include a transaction at the end of the bundle tipping 1000 lamports
	// to jitoTipAddress

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

// CreateBloxrouteTipInstructionToUseJitoBundles creates a transaction you can use to when using PostSubmitJitoBundle endpoints.
// This instruction should be the LAST transaction in your submission bundle
func CreateBloxrouteTipInstructionToUseJitoBundles(senderAddress solana.PublicKey, tipAmount uint64) (solana.Instruction, error) {
	recipient := solana.MustPublicKeyFromBase58(BloxrouteTipAddress)

	return system.NewTransferInstruction(tipAmount, senderAddress, recipient).Build(), nil
}
