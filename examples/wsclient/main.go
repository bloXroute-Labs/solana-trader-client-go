package main

import (
	"fmt"
	"github.com/bloXroute-Labs/serum-api/bxserum/provider"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func main() {
	callOrderbookWS()
	callOpenOrdersWS()
	callTradesWS()
	callTickersWS()
	callOrderbookWSStream()
	callTradesWSStream()
	callUnsettledWS()
}

// Unary response
func callOrderbookWS() {
	w, err := provider.NewWSClient()
	if err != nil {
		log.Fatalf("error dialing WS client: %v", err)
		return
	}
	defer w.Close()

	orderbook, err := w.GetOrderbook("ETH-USDT", 0)
	if err != nil {
		log.Errorf("error with GetOrderbook request for ETH-USDT: %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()

	orderbook, err = w.GetOrderbook("SOLUSDT", 2)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL-USDT: %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()

	orderbook, err = w.GetOrderbook("SOL:USDC", 3)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL:USDC: %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()
}

func callOpenOrdersWS() {
	w, err := provider.NewWSClient()
	if err != nil {
		log.Fatalf("error dialing WS client: %v", err)
		return
	}
	defer w.Close()

	orders, err := w.GetOpenOrders("SOLUSDC", "AFT8VayE7qr8MoQsW3wHsDS83HhEvhGWdbNSHRKeUDfQ")
	if err != nil {
		log.Errorf("error with GetOrders request for SOL-USDT: %v", err)
	} else {
		fmt.Println(orders)
	}

	fmt.Println()
}

func callUnsettledWS() {
	w, err := provider.NewWSClient()
	if err != nil {
		log.Fatalf("error dialing WS client: %v", err)
		return
	}
	defer w.Close()

	response, err := w.GetUnsettled("SOLUSDC", "AFT8VayE7qr8MoQsW3wHsDS83HhEvhGWdbNSHRKeUDfQ")
	if err != nil {
		log.Errorf("error with GetOrders request for SOL-USDT: %v", err)
	} else {
		fmt.Println(response)
	}

	fmt.Println()
}

func callTradesWS() {
	w, err := provider.NewWSClient()
	if err != nil {
		log.Fatalf("error dialing WS client: %v", err)
		return
	}
	defer w.Close()

	trades, err := w.GetTrades("SOLUSDC", 2)
	if err != nil {
		log.Errorf("error with GetTrades request for SOL-USDT: %v", err)
	} else {
		fmt.Println(trades)
	}

	fmt.Println()
}

func callTickersWS() {
	w, err := provider.NewWSClient()
	if err != nil {
		log.Fatalf("error dialing WS client: %v", err)
		return
	}
	defer w.Close()

	tickers, err := w.GetTickers("SOLUSDC")
	if err != nil {
		log.Errorf("error with GetTickers request for SOL-USDT: %v", err)
	} else {
		fmt.Println(tickers)
	}

	fmt.Println()
}

// Stream response
func callOrderbookWSStream() {
	fmt.Println("starting orderbook stream")
	w, err := provider.NewWSClient()
	if err != nil {
		log.Fatalf("error dialing WS client: %v", err)
		return
	}
	defer w.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	orderbookChan := make(chan *pb.GetOrderbookStreamResponse)

	err = w.GetOrderbookStream(ctx, "SOL/USDC", 3, orderbookChan)
	if err != nil {
		log.Errorf("error with GetOrderbookStream request for SOL/USDC: %v", err)
	} else {
		for i := 1; i <= 5; i++ {
			<-orderbookChan
			fmt.Printf("response %v received\n", i)
		}
	}
}

func callTradesWSStream() {
	fmt.Println("starting trades stream")
	w, err := provider.NewWSClient()
	if err != nil {
		log.Fatalf("error dialing WS client: %v", err)
		return
	}
	defer w.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tradesChan := make(chan *pb.GetTradesStreamResponse)

	err = w.GetTradesStream(ctx, "SOL/USDC", 3, tradesChan)
	if err != nil {
		log.Errorf("error with GetTradesStream request for SOL/USDC: %v", err)
	} else {
		for i := 1; i <= 5; i++ {
			<-tradesChan
			fmt.Printf("response %v received\n", i)
		}
	}
}
