package main

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/serum-api/bxserum/provider"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	log "github.com/sirupsen/logrus"
)

func main() {
	callOrderbookGRPC()
	callOrdersGRPC()
	callTradesGRPC()
	callTickersGRPC()
	callOrderbookGRPCStream()
	callTradesGRPCStream()
}

func callOrderbookGRPC() {
	g, err := provider.NewGRPCClient()
	if err != nil {
		log.Fatalf("error dialing GRPC client: %v", err)
		return
	}

	// Unary response
	orderbook, err := g.GetOrderbook(context.Background(), "ETH-USDT", 0)
	if err != nil {
		log.Errorf("error with GetOrderbook request for ETH-USDT: %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()

	orderbook, err = g.GetOrderbook(context.Background(), "SOLUSDT", 2)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOLUSDT: %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()

	orderbook, err = g.GetOrderbook(context.Background(), "SOL:USDC", 3)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL:USDC: %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()
}

func callOrdersGRPC() {
	g, err := provider.NewGRPCClient()
	if err != nil {
		log.Fatalf("error dialing GRPC client: %v", err)
		return
	}

	orders, err := g.GetOrders(context.Background(), "SOLUSDC", "HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc")
	if err != nil {
		log.Errorf("error with GetOrders request for SOLUSDC: %v", err)
	} else {
		fmt.Println(orders)
	}

	fmt.Println()

}

func callTradesGRPC() {
	g, err := provider.NewGRPCClient()
	if err != nil {
		log.Fatalf("error dialing GRPC client: %v", err)
		return
	}

	trades, err := g.GetTrades(context.Background(), "SOLUSDC", 5)
	if err != nil {
		log.Errorf("error with GetTrades request for SOLUSDC: %v", err)
	} else {
		fmt.Println(trades)
	}

	fmt.Println()

}

func callTickersGRPC() {
	g, err := provider.NewGRPCClient()
	if err != nil {
		log.Fatalf("error dialing GRPC client: %v", err)
		return
	}

	orders, err := g.GetTickers(context.Background(), "SOLUSDC")
	if err != nil {
		log.Errorf("error with GetTickers request for SOLUSDC: %v", err)
	} else {
		fmt.Println(orders)
	}

	fmt.Println()

}

func callOrderbookGRPCStream() {
	fmt.Println("starting orderbook stream")
	g, err := provider.NewGRPCClient()
	if err != nil {
		log.Fatalf("error dialing GRPC client - %v", err)
		return
	}

	orderbookChan := make(chan *pb.GetOrderbookStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	err = g.GetOrderbookStream(ctx, "SOL/USDC", 3, orderbookChan)
	if err != nil {
		log.Errorf("error with GetOrderbook stream request for SOL/USDC: %v", err)
	} else {
		for i := 1; i <= 5; i++ {
			<-orderbookChan
			fmt.Printf("response %v received\n", i)
		}
	}
}

func callTradesGRPCStream() {
	fmt.Println("starting trades stream")
	g, err := provider.NewGRPCClient()
	if err != nil {
		log.Fatalf("error dialing GRPC client - %v", err)
		return
	}

	tradesChan := make(chan *pb.GetTradesStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	err = g.GetTradesStream(ctx, "SOL/USDC", 3, tradesChan)
	if err != nil {
		log.Errorf("error with GetTrades stream request for SOL/USDC: %v", err)
	} else {
		for i := 1; i <= 5; i++ {
			<-tradesChan
			fmt.Printf("response %v received\n", i)
		}
	}
}
