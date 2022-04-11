package main

import (
	"fmt"
	"github.com/bloXroute-Labs/serum-api/bxserum/provider"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func main() {
	callWebsocket()
}

func callWebsocket() {
	w, err := provider.NewWSClient("ws://174.129.154.164:1810/ws")
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer w.Close()

	// One time request
	orderBook, err := w.GetOrderbook("ETH/USDT")
	if err != nil {
		log.Errorf("error with GetOrderbook request for ETH/USDT - %v", err)
	} else {
		fmt.Println(orderBook)
	}

	// Stream request
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	orderBookChan := make(chan *pb.GetOrderbookStreamResponse)

	err = w.GetOrderbookStream(ctx, "SOL/USDT", orderBookChan)
	if err != nil {
		log.Errorf("error with GetOrderbookStream request for SOL/USDT - %v", err)
	} else {
		for i := 1; i <= 5; i++ {
			<-orderBookChan
			fmt.Printf("response %v received\n", i)
		}
	}
}
