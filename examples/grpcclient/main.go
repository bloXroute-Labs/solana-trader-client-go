package main

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/serum-api/bxserum/provider"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	log "github.com/sirupsen/logrus"
)

func main() {
	callGRPC()
	callGRPCStream()
}

func callGRPC() {
	g, err := provider.NewGRPCClient()
	if err != nil {
		log.Fatalf("error dialing GRPC client - %v", err)
		return
	}

	// Unary response
	orderbook, err := g.GetOrderbook(context.Background(), "ETH-USDT")
	if err != nil {
		log.Errorf("error with GetOrderbook request for ETH-USDT - %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()

	orderbook, err = g.GetOrderbook(context.Background(), "SOLUSDT")
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOLUSDT - %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()

	orderbook, err = g.GetOrderbook(context.Background(), "SOL:USDC")
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL:USDC - %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()
}

func callGRPCStream() {
	g, err := provider.NewGRPCClient()
	if err != nil {
		log.Fatalf("error dialing GRPC client - %v", err)
		return
	}

	orderbookChan := make(chan *pb.GetOrderbookStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	err = g.GetOrderbookStream(ctx, "SOL/USDC", orderbookChan)
	if err != nil {
		log.Errorf("error with GetOrderbook stream request for SOL/USDC - %v", err)
	} else {
		for i := 1; i <= 5; i++ {
			<-orderbookChan
			fmt.Printf("response %v received\n", i)
		}
	}
}
