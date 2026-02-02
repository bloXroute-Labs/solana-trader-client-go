package main

import (
	"context"
	"time"

	"github.com/bloXroute-Labs/solana-trader-client-go/connections"
	"github.com/bloXroute-Labs/solana-trader-client-go/examples/config"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/gagliardetto/solana-go"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/programs/system"
	log "github.com/sirupsen/logrus"
)

func main() {
	utils.InitLogger()

	environment := "testnet" // "mainnet", "testnet", "local"
	region := "ny"           // "ny", "uk", "fr"

	client := setupHTTPClient(config.Env(environment), config.HTTPUrls[config.Region(region)])

	submitNcountWithContentType(client, 1, "application/json", "postSubmit-cold-start-json")

	iterationCount := 5_000
	iterationCount = 1

	//time.Sleep(11 * time.Second)
	submitNcountWithContentType(client, iterationCount, "application/json", "postSubmit-1-json")
	//time.Sleep(11 * time.Second)
	submitNcountWithContentType(client, iterationCount, "text/plain", "postSubmit-1-text")
}

func submitNcountWithContentType(h provider.HTTPClientTraderAPI, iterationCount int, contentType, correlationID string) {
	connections.GlobalPostSubmitContentType = contentType
	connections.GlobalCorrelationID = correlationID

	durationsSum := int64(0)

	for i := 0; i < iterationCount; i++ {
		startTime := time.Now()
		callPostSubmit(h)
		endTime := time.Now()
		durationsSum += endTime.Sub(startTime).Nanoseconds()
	}
	log.Infof("Post Submit Performance for %s over %d iterations, AVG ns: %d", correlationID, iterationCount, durationsSum/int64(iterationCount))
}

func setupHTTPClient(env config.Env, endpoint string) provider.HTTPClientTraderAPI {
	var h provider.HTTPClientTraderAPI
	var err error

	switch env {
	case config.EnvLocal:
		h = provider.NewHTTPLocal()
	case config.EnvTestnet:
		h = provider.NewHTTPTestnet()
	case config.EnvMainnet:
		h, err = provider.NewHTTPClientFullService(endpoint)
	}
	if err != nil {
		log.Fatalf("error dialing HTTP client: %v", err)
	}

	return h
}

func callPostSubmit(h provider.HTTPClientTraderAPI) bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := h.GetRecentBlockHash(ctx)
	if err != nil {
		log.Errorf("error with GetRecentBlockHash request: %v", err)
		return true
	}

	bh := solana.MustHashFromBase58(response.BlockHash)

	wlt := solana.NewWallet()
	privateKey, err := transaction.LoadPrivateKeyFromEnv()
	priceLimitIx, err := computebudget.NewSetComputeUnitPriceInstruction(uint64(200000000)).ValidateAndBuild()
	if err != nil {
		return false
	}

	instructions := []solana.Instruction{priceLimitIx}
	for i := 0; i < 40; i++ {
		instructions = append(instructions, system.NewTransferInstruction(100000, privateKey.PublicKey(), solana.MustPublicKeyFromBase58("FZwLKcQupnTy2CbaVMGGsutxDtjv9CqYVDJxiNZSj5Xi")).Build())
	}
	tx1, err := solana.NewTransaction(instructions, bh, solana.TransactionPayer(privateKey.PublicKey()))
	if err != nil {
		return false
	}

	tx1.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(privateKey.PublicKey()) {
			return &privateKey
		}
		return &wlt.PrivateKey
	})

	_, err = h.SignAndSubmit(ctx, &pb.TransactionMessage{
		Content: tx1.MustToBase64()}, false, false, true)
	if err != nil {
		//log.Errorf("failed to sign and submit order (%v)", err)
		return true
	}

	//log.Infof("submitted bundle order to trader api %v", resp)

	return false
}
