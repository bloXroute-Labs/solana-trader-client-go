package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/bloXroute-Labs/serum-api/bxserum/provider"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	log "github.com/sirupsen/logrus"
)

func main() {
	callOrderbookGRPC()
	callOpenOrdersGRPC()
	callTradesGRPC()
	callTickersGRPC()
	callOrderbookGRPCStream()
	callTradesGRPCStream()
	callUnsettledGRPC()

	ownerAddr, ok := os.LookupEnv("PUBLIC_KEY")
	if !ok {
		log.Infof("PUBLIC_KEY environment variable not set")
		log.Infof("Skipping Place and Cancel Order examples")
		return
	}

	ooAddr, ok := os.LookupEnv("OPEN_ORDERS")
	if !ok {
		log.Infof("OPEN_ORDERS environment variable not set")
		log.Infof("Skipping Place, Cancel and Settle examples")
		return
	}

	clientID := callPlaceOrderGRPC(ownerAddr, ooAddr)
	callCancelByClientOrderIDGRPC(ownerAddr, ooAddr, clientID)
	callPostSettleGRPC()
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

func callOpenOrdersGRPC() {
	g, err := provider.NewGRPCClient()
	if err != nil {
		log.Fatalf("error dialing GRPC client: %v", err)
		return
	}

	orders, err := g.GetOpenOrders(context.Background(), "SOLUSDC", "HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc")
	if err != nil {
		log.Errorf("error with GetOrders request for SOLUSDC: %v", err)
	} else {
		fmt.Println(orders)
	}

	fmt.Println()

}

func callUnsettledGRPC() {
	g, err := provider.NewGRPCClient()
	if err != nil {
		log.Fatalf("error dialing GRPC client: %v", err)
		return
	}

	response, err := g.GetUnsettled(context.Background(), "SOLUSDC", "HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc")
	if err != nil {
		log.Errorf("error with GetOrders request for SOLUSDC: %v", err)
	} else {
		fmt.Println(response)
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

	orderbookChan := make(chan *pb.GetOrderbooksStreamResponse)
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

const (
	// SOL/USDC market
	marketAddr = "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT"

	orderSide   = pb.Side_S_ASK
	orderType   = pb.OrderType_OT_LIMIT
	orderPrice  = float64(170200)
	orderAmount = float64(0.1)
)

func callPlaceOrderGRPC(ownerAddr, ooAddr string) uint64 {
	fmt.Println("starting place order")

	g, err := provider.NewGRPCClient()
	if err != nil {
		log.Errorf("error dialing GRPC client (%v)", err)
		return 0
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// generate a random clientId for this order
	rand.Seed(time.Now().UnixNano())
	clientOrderID := rand.Uint64()

	opts := provider.PostOrderOpts{
		ClientOrderID:     clientOrderID,
		OpenOrdersAddress: ooAddr,
	}

	sig, err := g.SubmitOrder(ctx, ownerAddr, ownerAddr, marketAddr,
		orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Fatalf("failed to submit order (%v)", err)
	}

	fmt.Printf("placed order %v with clientOrderID %v\n", sig, clientOrderID)

	return clientOrderID
}

func callCancelByClientOrderIDGRPC(ownerAddr, ooAddr string, clientID uint64) {
	fmt.Println("starting cancel order by client order ID")
	time.Sleep(30 * time.Second)
	g, err := provider.NewGRPCClient()
	if err != nil {
		log.Errorf("error dialing GRPC client (%v)", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err = g.SubmitCancelByClientOrderID(ctx, clientID, ownerAddr,
		marketAddr, ooAddr, true)
	if err != nil {
		log.Fatalf("failed to cancel order by client order ID (%v)", err)
	}

	fmt.Printf("canceled order for clientID %v\n", clientID)
}

func callPostSettleGRPC() {
	fmt.Println("starting post settle")
	g, err := provider.NewGRPCClient()
	if err != nil {
		log.Fatalf("error dialing GRPC client - %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//publicKey, _ := os.LookupEnv("PUBLIC_KEY")
	// Stream response
	sig, err := g.SettleFunds(ctx, "F75gCEckFAyeeCWA9FQMkmLCmke7ehvBnZeVZ3QgvJR7", "SOL/USDC", "F75gCEckFAyeeCWA9FQMkmLCmke7ehvBnZeVZ3QgvJR7", "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv", "")
	if err != nil {
		log.Errorf("error with post transaction stream request for SOL/USDC: %v", err)
		return
	}

	fmt.Printf("response signature received: %v", sig)
}
