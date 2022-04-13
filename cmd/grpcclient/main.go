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
		log.Errorf("error dialing GRPC client - %v", err)
	}

	// One time request
	orderbook, err := g.GetOrderbook(context.Background(), "ETH-USDT")
	if err != nil {
		log.Errorf("error with GetOrderbook request for ETH/USDT - %v", err)
	} else {
		fmt.Println(orderbook)
	}
	fmt.Println()

	orderbook, err = g.GetOrderbook(context.Background(), "SOL-USDT")
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL/USDT - %v", err)
	} else {
		fmt.Println(orderbook)
	}
	fmt.Println()

	g, err = provider.NewGRPCClientWithEndpoint("http://1.1.1.1:1809")
	if err != nil {
		fmt.Println("successfully got error when dialing invalid IP")
	} else {
		fmt.Println("did not get error when dialing invalid IP")
	}
}

func callGRPCStream() {
	g, err := provider.NewGRPCClient()
	if err != nil {
		log.Errorf("error dialing GRPC client - %v", err)
	}

	// Stream response
	respChan := make(chan *pb.GetOrderbookStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = g.GetOrderbookStream(ctx, "ETH/USDC", respChan)
	if err != nil {
		log.Errorf("error with GetOrderbook stream request for ETH/USDT - %v", err)
	} else {
		for i := 0; i < 5; i++ {
			fmt.Println(<-respChan)
			fmt.Println()
		}
	}
}
