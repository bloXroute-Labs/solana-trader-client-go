package transaction

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
	rpc_endpoint     = solanarpc.MainNetBeta_RPC
	recipientAddress = "FmZ9kC8bRVsFTgAWrXUyGHp3dN3HtMxJmoi2ijdaYGwi"
)

type txConfirmation struct {
	TxHash string `json:"txHash"`
}

func main() {
	rpcClient := solanarpc.New(rpc_endpoint)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	privateKeyPath := os.Getenv("PRIVATE_KEY_PATH")
	if privateKeyPath == "" {
		log.Fatalf("env variable `PRIVATE_KEY_PATH` not set")
	}

	recentBlockhash, err := rpcClient.GetRecentBlockhash(ctx, solanarpc.CommitmentFinalized)
	if err != nil {
		log.Fatal(err)
	}

	unsignedTx, err := unsignedTx(privateKeyPath, recentBlockhash)
	signedTx, err := transaction.SignTx(unsignedTx.MustToBase64())
	if err != nil {
		return
	}

	signature, err := sendTx(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("transaction not sent successfully: %v", err)
	}

	fmt.Printf("tx %s sent successfully\n", signature.String())
}

func unsignedTx(privateKeyPath string, recentBlockHash *solanarpc.GetRecentBlockhashResult) (*solana.Transaction, error) {
	privateKey, err := solana.PrivateKeyFromSolanaKeygenFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	recipient := solana.MustPublicKeyFromBase58(recipientAddress)

	tx, err := solana.NewTransaction([]solana.Instruction{
		system.NewTransferInstruction(1, privateKey.PublicKey(), recipient).Build(),
	}, recentBlockHash.Value.Blockhash)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func sendTx(ctx context.Context, txBase64 string) (solana.Signature, error) {
	rpcClient := solanarpc.New(solanarpc.MainNetBeta_RPC)
	wsClient, err := solanaws.Connect(ctx, solanarpc.MainNetBeta_WS)
	if err != nil {
		return solana.Signature{}, err
	}

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
