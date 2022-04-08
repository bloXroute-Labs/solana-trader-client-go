package main

import (
	"fmt"
	"github.com/bloXroute-Labs/serum-api/bxserum/provider"
	"github.com/bloXroute-Labs/serum-api/logger"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"golang.org/x/net/context"
	"log"
)

func main() {
	callWebsocket()
}

func callWebsocket() {
	w, err := provider.NewWSClient("ws://174.129.154.164:1810/ws")
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer w.CloseConn()

	// One time request
	orderBook, err := w.GetOrderbook("ETH/USDT")
	if err != nil {
		logger.Log().Errorw("error with GetOrderbook request for ETH/USDT", "err", err)
	} else {
		fmt.Println(orderBook)
	}

	// Stream request
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	orderBookChan := make(chan *pb.GetOrderbookStreamResponse)

	err = w.GetOrderbookStream(ctx, "SOL/USDT", orderBookChan)
	if err != nil {
		logger.Log().Errorw("error with GetOrderbookStream request for SOL/USDT", "err", err)
	}

	for i := 1; i <= 5; i++ {
		orderBook := <-orderBookChan
		fmt.Println(orderBook)

		fmt.Printf("response %v received\n", i)
	}
}
