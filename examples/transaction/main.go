package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bloXroute-Labs/serum-api/bxserum/transaction"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
	sendandconfirmtransaction "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	solanaws "github.com/gagliardetto/solana-go/rpc/ws"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	rpcEndpoint      = solanarpc.MainNetBeta_RPC
	recipientAddress = "FmZ9kC8bRVsFTgAWrXUyGHp3dN3HtMxJmoi2ijdaYGwi"
)

type txConfirmation struct {
	TxHash string `json:"txHash"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rpcClient := solanarpc.New(rpcEndpoint)
	wsClient, err := solanaws.Connect(ctx, solanarpc.MainNetBeta_WS)
	if err != nil {
		log.Fatal(err)
	}

	defer cancel()
	recentBlockhash, err := rpcClient.GetRecentBlockhash(ctx, solanarpc.CommitmentFinalized)
	if err != nil {
		log.Fatal(err)
	}

	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		log.Fatalf("env variable `PRIVATE_KEY` not set")
	}

	unsignedTx, err := unsignedTransaction(privateKey, recentBlockhash)
	if err != nil {
		log.Fatal(err)
	}
	unsignedTxMessageBase64 := unsignedTx.Message.ToBase64()

	signedTx, err := transaction.SignTxMessage(unsignedTxMessageBase64)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tx message signed, ready to send tx")

	signature, err := sendAndConfirmTx(context.Background(), signedTx, rpcClient, wsClient)
	if err != nil {
		log.Fatalf("transaction not sent successfully: %v", err)
	}

	fmt.Printf("tx %s sent and confirmed successfully\n", signature.String())
}

func unsignedTransaction(privateKey string, recentBlockHash *solanarpc.GetRecentBlockhashResult) (*solana.Transaction, error) {
	pKey, err := solana.PrivateKeyFromBase58(privateKey)
	if err != nil {
		return nil, err
	}
	recipient := solana.MustPublicKeyFromBase58(recipientAddress)

	return solana.NewTransaction([]solana.Instruction{
		system.NewTransferInstruction(1, pKey.PublicKey(), recipient).Build(),
	}, recentBlockHash.Value.Blockhash)
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

func confirmTx(signature solana.Signature) (string, error) {
	url := fmt.Sprintf("https://public-api.solscan.io/transaction/%s", signature.String())
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var confirmation txConfirmation
	err = json.Unmarshal(b, &confirmation)
	if err != nil {
		return "", err
	}

	return confirmation.TxHash, nil
}
