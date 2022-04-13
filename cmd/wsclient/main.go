package main

import (
	"fmt"
	"github.com/bloXroute-Labs/serum-api/bxserum/provider"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func main() {
	callWS()
	callWSStream()
}

func callWS() {
	w, err := provider.NewWSClient()
	if err != nil {
		log.Fatalf("error dialing WS client - %v", err)
		return
	}
	defer w.Close()

	// Unary response
	orderbook, err := w.GetOrderbook("ETH-USDT")
	if err != nil {
		log.Errorf("error with GetOrderbook request for ETH-USDT - %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()

	orderbook, err = w.GetOrderbook("SOLUSDT")
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL-USDT - %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()

	orderbook, err = w.GetOrderbook("SOL:USDC")
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL:USDC - %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()
}

func callWSStream() {
	w, err := provider.NewWSClient()
	if err != nil {
		log.Fatalf("error dialing WS client - %v", err)
		return
	}
	defer w.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	orderbookChan := make(chan *pb.GetOrderbookStreamResponse)

	// Stream response
	err = w.GetOrderbookStream(ctx, "SOL/USDC", orderbookChan)
	if err != nil {
		log.Errorf("error with GetOrderbookStream request for SOL/USDC - %v", err)
	} else {
		for i := 1; i <= 5; i++ {
			<-orderbookChan
			fmt.Printf("response %v received\n", i)
		}
	}
}
