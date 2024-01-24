package main

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/examples/config"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	utils.InitLogger()
	run()
}

func run() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	var g *provider.GRPCClient

	switch cfg.Env {
	case config.EnvLocal:
		g, err = provider.NewGRPCLocal()
	case config.EnvTestnet:
		g, err = provider.NewGRPCTestnet()
	case config.EnvMainnet:
		g, err = provider.NewGRPCClient()
	}
	if err != nil {
		log.Fatalf("error dialing GRPC client: %v", err)
	}

	privateKey := os.Getenv("PRIVATE_KEY")
	if privateKey == "" {
		log.Fatalf("env variable `PRIVATE_KEY` not set")
	}
	recipient := os.Getenv("RECIPIENT_ADDRESS")
	if recipient == "" {
		log.Fatalf("env variable `RECIPIENT_ADDRESS` not set")
	}
	ctx := context.Background()
	recentBlockhash, err := g.RecentBlockHash(ctx)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := createTx(privateKey, recipient, recentBlockhash.BlockHash)
	txBase64, err := tx.ToBase64()
	if err != nil {
		log.Fatalf("transaction not converted to bytes successfully: %v", err)
	}
	response, err := g.PostSubmit(ctx, &pb.TransactionMessage{
		Content:   txBase64,
		IsCleanup: false,
	}, true)
	if err != nil {
		panic(err)
	}
	fmt.Println("response.Signature : ", response.Signature)
}

func createTx(privateKeyStr string, recipientAddress, recentBlockHash string) (*solana.Transaction, error) {
	privateKey, err := solana.PrivateKeyFromBase58(privateKeyStr)
	if err != nil {
		return nil, err
	}
	recipient := solana.MustPublicKeyFromBase58(recipientAddress)

	tx, err := solana.NewTransaction([]solana.Instruction{
		system.NewTransferInstruction(1, privateKey.PublicKey(), recipient).Build(),
	}, solana.MustHashFromBase58(recentBlockHash))
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
