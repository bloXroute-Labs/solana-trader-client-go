package main

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/examples/config"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	utils.InitLogger()
	failed := run()
	if failed {
		log.Fatal("one or multiple examples failed")
	}
}

func run() bool {
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
	ownerAddr, ok := os.LookupEnv("PUBLIC_KEY")
	if !ok {
		log.Fatal("PUBLIC_KEY environment variable not set: will skip place/cancel/settle examples")
	}

	callJupiterSwap(g, ownerAddr)
	return false
}

func callJupiterSwap(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting Jupiter swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("Jupiter swap")
	sig, err := g.SubmitJupiterSwap(ctx, &pb.PostJupiterSwapRequest{
		OwnerAddress: ownerAddr,
		InToken:      "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		OutToken:     "So11111111111111111111111111111111111111112",
		Slippage:     0.01,
		InAmount:     0.25,
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_ABORT_ON_FIRST_ERROR,
		SkipPreFlight:  false,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Jupiter swap transaction signature : %s", sig)
	return false
}

func callGetJupiterQuotes(g *provider.GRPCClient) bool {
	log.Info("starting get Jupiter quotes test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inToken := "SOL"
	outToken := "USDT"
	amount := 0.01
	slippage := float64(5)
	limit := int32(3)

	quotes, err := g.GetJupiterQuotes(ctx, &pb.GetJupiterQuotesRequest{
		InToken:  inToken,
		OutToken: outToken,
		InAmount: amount,
		Slippage: slippage,
		Limit:    limit,
	})
	if err != nil {
		log.Errorf("error with GetQuotes request for %s to %s: %v", inToken, outToken, err)
		return true
	}

	if err != nil {
		log.Errorf("error with GetJupiterQuotes request for %s to %s: %v", inToken, outToken, err)
		return true
	}

	if len(quotes.Routes) != 3 {
		log.Errorf("did not get back 3 quotes, got %v quotes", len(quotes.Routes))
		return true
	}
	for _, route := range quotes.Routes {
		log.Infof("best route for Jupiter is %v", route)
	}

	fmt.Println()
	return false
}
