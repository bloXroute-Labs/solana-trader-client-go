package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/bloXroute-Labs/serum-api/bxserum/transaction"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	sendandconfirmtransaction "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	solanaws "github.com/gagliardetto/solana-go/rpc/ws"
	"log"
	"os"
)

const (
	rpcEndpoint      = solanarpc.MainNetBeta_RPC
	wsEndpoint       = solanarpc.MainNetBeta_WS
	recipientAddress = "FmZ9kC8bRVsFTgAWrXUyGHp3dN3HtMxJmoi2ijdaYGwi"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// get clients
	rpcClient := solanarpc.New(rpcEndpoint)
	wsClient, err := solanaws.Connect(ctx, wsEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	// get recent block hash
	recentBlockhash, err := rpcClient.GetRecentBlockhash(ctx, solanarpc.CommitmentFinalized)
	if err != nil {
		log.Fatal(err)
	}

	// get private key from env variable
	privateKeyBase58 := os.Getenv("PRIVATE_KEY")
	if privateKeyBase58 == "" {
		log.Fatal("env variable `PRIVATE_KEY` not set")
	}
	privateKey, err := solana.PrivateKeyFromBase58(privateKeyBase58)
	if err != nil {
		log.Fatal(err)
	}

	// create unsigned tx using block hash and private key
	unsignedTx, err := unsignedTransaction(privateKey.PublicKey(), recentBlockhash)
	if err != nil {
		log.Fatal(err)
	}
	unsignedTxBytes, err := partialMarshal(unsignedTx)
	unsignedTxBase64 := base64.StdEncoding.EncodeToString(unsignedTxBytes)

	// sign tx
	signedTx, err := transaction.SignTx(unsignedTxBase64) // gets private key from environment variable
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tx signed, ready to send")

	// send and confirm tx
	signature, err := sendAndConfirmTx(context.Background(), signedTx, rpcClient, wsClient)
	if err != nil {
		log.Fatalf("tx not sent successfully: %v", err)
	}
	fmt.Printf("tx sent and confirmed successfully, signature: %s\n", signature.String())
}

// creates a transaction with a zero signature (private key only used to get public key)
func unsignedTransaction(publicKey solana.PublicKey, recentBlockHash *solanarpc.GetRecentBlockhashResult) (*solana.Transaction, error) {
	recipient := solana.MustPublicKeyFromBase58(recipientAddress)

	tx, err := solana.NewTransaction([]solana.Instruction{
		system.NewTransferInstruction(1, publicKey, recipient).Build(),
	}, recentBlockHash.Value.Blockhash)
	if err != nil {
		return nil, err
	}

	tx.Signatures = append(tx.Signatures, solana.Signature{}) // adding a zero signature
	return tx, nil
}

// marshals transaction without checking number of signatures
func partialMarshal(tx *solana.Transaction) ([]byte, error) {
	messageBytes, err := tx.Message.MarshalBinary()
	if err != nil {
		return nil, err
	}

	var signatureCount []byte
	bin.EncodeCompactU16Length(&signatureCount, len(tx.Signatures))

	output := make([]byte, 0, len(signatureCount)+len(signatureCount)*64+len(messageBytes))
	output = append(output, signatureCount...) // signatureCount | signatures | message
	for _, sig := range tx.Signatures {
		output = append(output, sig[:]...)
	}
	output = append(output, messageBytes...)

	return output, nil
}

func sendAndConfirmTx(ctx context.Context, txBase64 string, rpcClient *solanarpc.Client, wsClient *solanaws.Client) (solana.Signature, error) {
	txBytes, err := solanarpc.DataBytesOrJSONFromBase64(txBase64)
	if err != nil {
		return solana.Signature{}, err
	}

	tx, err := solanarpc.TransactionWithMeta{Transaction: txBytes}.GetTransaction()
	if err != nil {
		return solana.Signature{}, err
	}

	return sendandconfirmtransaction.SendAndConfirmTransaction(ctx, rpcClient, wsClient, tx)
}
