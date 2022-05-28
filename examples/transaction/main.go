package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/bloXroute-Labs/serum-api/borsh/serumborsh"
	"github.com/bloXroute-Labs/serum-api/service"
	"github.com/bloXroute-Labs/serum-api/utils"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	solanarpc "github.com/gagliardetto/solana-go/rpc"
)

const (
	defaultRpcEndpoint = "https://solana-api.projectserum.com"
	wsEndpoint         = solanarpc.MainNetBeta_WS

	// SOL/USDC market
	marketAddress = "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT"
)

type txConfirmation struct {
	TxHash string `json:"txHash"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rpcEndpoint, ok := os.LookupEnv("RPC_ENDPOINT")
	if !ok {
		rpcEndpoint = defaultRpcEndpoint
	}

	publicKeyStr, ok := os.LookupEnv("PUBLIC_KEY")
	if !ok {
		log.Fatal("PUBLIC_KEY environment variable not set")
	}

	privateKeyStr, ok := os.LookupEnv("PRIVATE_KEY")
	if !ok {
		log.Fatal("PRIVATE_KEY environment variable not set")
	}
	privateKey, err := solana.PrivateKeyFromBase58(privateKeyStr)
	if err != nil {
		panic(err)
	}

	ooAddress, _ := os.LookupEnv("OPEN_ORDERS")

	marketService, err := service.NewMarket(utils.MarketsFileFlag.Value)
	if err != nil {
		panic(err)
	}

	solanaService, err := service.NewSolana(rpcEndpoint, wsEndpoint)
	if err != nil {
		panic(err)
	}

	serumService := service.NewSerum(ctx, marketService, solanaService, false)

	// generate a random clientId for this order
	rand.Seed(time.Now().UnixNano())
	clientId := rand.Uint64()

	// get partially signed transaction
	txBytes, _, err := serumService.PlaceOrder(ctx, "SOL/USDC", publicKeyStr,
		publicKeyStr, serumborsh.OSSell, 0.1, 170200, ooAddress, clientId)
	if err != nil {
		panic(err)
	}

	fmt.Println("partial transaction received")
	signAndConfirmTx(solanaService, txBytes, privateKey, false)

	fmt.Printf("attempting to cancel order with clientId %x\n", clientId)

	// try to cancel the order
	txCancel, err := serumService.CancelOrderByClientID(ctx, clientId,
		marketAddress, publicKeyStr, ooAddress, nil)
	if err != nil {
		log.Fatalf("failed to create a CancelOrder tx (%v)", err)
	}

	signAndConfirmTx(solanaService, txCancel, privateKey, true)
}

func signAndConfirmTx(solanaService service.Solana, tx []byte, privateKey solana.PrivateKey, skipPreFlight bool) {
	txBytes, txSignSig := signTx(tx, privateKey)

	txSig, err := sendTx(solanaService, txBytes, skipPreFlight)
	if err != nil {
		log.Fatalf("failed to send tx (%v)", err)
	}
	fmt.Printf("tx %s sent successfully\n", txSig)

	time.Sleep(time.Second * 25)

	txConf, err := confirmTx(txSignSig)
	if err != nil {
		log.Printf("failed to confirm tx (%v)", err)
		return
	}

	if txConf == txSignSig.String() {
		fmt.Printf("tx %s confirmed successfully\n", txConf)
	} else {
		log.Fatalf("tx from confirmation %s not equal to tx sent %s",
			txConf, txSignSig.String())
	}
}

func signTx(txBytes []byte, privateKey solana.PrivateKey) ([]byte, solana.Signature) {
	var transaction solana.Transaction

	// decode transaction for signing
	err := bin.NewBinDecoder(txBytes).Decode(&transaction)
	if err != nil {
		panic(err)
	}

	// sign with final private key
	messageBytes, err := transaction.Message.MarshalBinary()
	if err != nil {
		panic(err)
	}

	signature, err := privateKey.Sign(messageBytes)
	if err != nil {
		panic(err)
	}

	for i, signed := range transaction.Signatures {
		if signed.IsZero() {
			transaction.Signatures[i] = signature
			break
		}
	}

	signedTxBytes, err := transaction.MarshalBinary()
	if err != nil {
		panic(err)
	}

	err = transaction.VerifySignatures()
	if err != nil {
		panic(err)
	}

	return signedTxBytes, signature
}

func sendTx(solanaService service.Solana, transactionBytes []byte, skipPreFlight bool) (solana.Signature, error) {
	txBase64 := base64.StdEncoding.EncodeToString(transactionBytes)
	fmt.Println("transaction has been resigned and verified, submitting...")
	fmt.Println(txBase64)

	return solanaService.SendTransaction(context.Background(), txBase64, skipPreFlight)
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
