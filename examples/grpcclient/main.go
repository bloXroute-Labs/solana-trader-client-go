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
	g, err := provider.NewGRPCTestnet()
	if err != nil {
		log.Fatalf("error dialing GRPC client: %v", err)
		return
	}

	callMarketsGRPC(g)
	callOrderbookGRPC(g)
	callOpenOrdersGRPC(g)
	callTickersGRPC(g)
	callOrderbookGRPCStream(g)
	callTradesGRPCStream(g)
	callUnsettledGRPC(g)

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

	clientID := callPlaceOrderGRPC(g, ownerAddr, ooAddr)
	callCancelByClientOrderIDGRPC(g, ownerAddr, ooAddr, clientID)
	callPostSettleGRPC(g, ownerAddr, ooAddr)
}

func callMarketsGRPC(g *provider.GRPCClient) {
	markets, err := g.GetMarkets(context.Background())
	if err != nil {
		log.Errorf("error with GetMarkets request: %v", err)
	} else {
		fmt.Println(markets)
	}

	fmt.Println()
}

func callOrderbookGRPC(g *provider.GRPCClient) {
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

func callOpenOrdersGRPC(g *provider.GRPCClient) {
	orders, err := g.GetOpenOrders(context.Background(), "SOLUSDC", "FFqDwRq8B4hhFKRqx7N1M6Dg6vU699hVqeynDeYJdPj5")
	if err != nil {
		log.Errorf("error with GetOrders request for SOLUSDC: %v", err)
	} else {
		fmt.Println(orders)
	}

	fmt.Println()
}

func callUnsettledGRPC(g *provider.GRPCClient) {
	response, err := g.GetUnsettled(context.Background(), "SOLUSDC", "HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc")
	if err != nil {
		log.Errorf("error with GetOrders request for SOLUSDC: %v", err)
	} else {
		fmt.Println(response)
	}

	fmt.Println()

}

func callTickersGRPC(g *provider.GRPCClient) {
	orders, err := g.GetTickers(context.Background(), "SOLUSDC")
	if err != nil {
		log.Errorf("error with GetTickers request for SOLUSDC: %v", err)
	} else {
		fmt.Println(orders)
	}

	fmt.Println()
}

func callOrderbookGRPCStream(g *provider.GRPCClient) {
	fmt.Println("starting orderbook stream")

	orderbookChan := make(chan *pb.GetOrderbooksStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	err := g.GetOrderbookStream(ctx, "SOL/USDC", 3, orderbookChan)
	if err != nil {
		log.Errorf("error with GetOrderbook stream request for SOL/USDC: %v", err)
	} else {
		for i := 1; i <= 5; i++ {
			<-orderbookChan
			fmt.Printf("response %v received\n", i)
		}
	}
}

func callTradesGRPCStream(g *provider.GRPCClient) {
	fmt.Println("starting trades stream")

	tradesChan := make(chan *pb.GetTradesStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	err := g.GetTradesStream(ctx, "SOL/USDC", 3, tradesChan)
	if err != nil {
		log.Errorf("error with GetTrades stream request for SOL/USDC: %v", err)
	} else {
		for i := 1; i <= 3; i++ {
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

func callPlaceOrderGRPC(g *provider.GRPCClient, ownerAddr, ooAddr string) uint64 {
	fmt.Println("starting place order")

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

func callCancelByClientOrderIDGRPC(g *provider.GRPCClient, ownerAddr, ooAddr string, clientID uint64) {
	fmt.Println("starting cancel order by client order ID")
	time.Sleep(30 * time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := g.SubmitCancelByClientOrderID(ctx, clientID, ownerAddr,
		marketAddr, ooAddr, true)
	if err != nil {
		log.Fatalf("failed to cancel order by client order ID (%v)", err)
	}

	fmt.Printf("canceled order for clientID %v\n", clientID)
}

func callPostSettleGRPC(g *provider.GRPCClient, ownerAddr, ooAddr string) {
	fmt.Println("starting post settle")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig, err := g.SubmitSettle(ctx, ownerAddr, "SOL/USDC", "F75gCEckFAyeeCWA9FQMkmLCmke7ehvBnZeVZ3QgvJR7", "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv", ooAddr, false)
	if err != nil {
		log.Errorf("error with post transaction stream request for SOL/USDC: %v", err)
		return
	}

	fmt.Printf("response signature received: %v", sig)
}
