package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bloXroute-Labs/serum-client-go/bxserum/provider"
	pb "github.com/bloXroute-Labs/serum-client-go/proto"
	log "github.com/sirupsen/logrus"
)

func main() {
	g, err := provider.NewGRPCTestnet()
	if err != nil {
		log.Fatalf("error dialing GRPC client: %v", err)
		return
	}

	// informational methods
	callMarketsGRPC(g)
	callOrderbookGRPC(g)
	callOpenOrdersGRPC(g)
	callTickersGRPC(g)
	callOrderbookGRPCStream(g)
	callTradesGRPCStream(g)
	callUnsettledGRPC(g)
	callGetAccountBalanceGRPC(g)

	// calls below this place an order and immediately cancel it
	// you must specify:
	//	- PRIVATE_KEY (by default loaded during provider.NewGRPCClient()) to sign transactions
	// 	- PUBLIC_KEY to indicate which account you wish to trade from
	//	- OPEN_ORDERS to indicate your Serum account to speed up lookups (optional)
	ownerAddr, ok := os.LookupEnv("PUBLIC_KEY")
	if !ok {
		log.Infof("PUBLIC_KEY environment variable not set: will skip place/cancel/settle examples")
		return
	}

	ooAddr, ok := os.LookupEnv("OPEN_ORDERS")
	if !ok {
		log.Infof("OPEN_ORDERS environment variable not set: requests will be slower")
	}

	payerAddr, ok := os.LookupEnv("PAYER")
	if !ok {
		log.Infof("PAYER environment variable not set: will be set to owner address")
		payerAddr = ownerAddr
	}

	orderLifecycleTest(g, ownerAddr, ooAddr)
	cancelAll(g, ownerAddr, payerAddr, ooAddr)
	callReplaceByClientOrderID(g, ownerAddr, payerAddr, ooAddr)
	callReplaceOrder(g, ownerAddr, payerAddr, ooAddr)
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

func callGetAccountBalanceGRPC(g *provider.GRPCClient) {
	response, err := g.GetAccountBalance(context.Background(), "HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc")
	if err != nil {
		log.Errorf("error with GetAccountBalance request for HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc: %v", err)
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := g.GetOrderbookStream(ctx, []string{"SOL/USDC", "xxx", "SOL-USDT"}, 3)
	if err != nil {
		log.Errorf("error with GetOrderbook stream request for SOL/USDC: %v", err)
		return
	}

	orderbookCh := stream.Channel(0)
	for i := 1; i <= 5; i++ {
		data, ok := <-orderbookCh
		if !ok {
			// channel closed
			return
		}

		fmt.Printf("response %v received, data %v \n", i, data)
	}
}

func callTradesGRPCStream(g *provider.GRPCClient) {
	fmt.Println("starting trades stream")

	tradesChan := make(chan *pb.GetTradesStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := g.GetTradesStream(ctx, "SOL/USDC", 3)
	if err != nil {
		log.Errorf("error with GetTrades stream request for SOL/USDC: %v", err)
	}
	stream.Into(tradesChan)
	for i := 1; i <= 3; i++ {
		_, ok := <-tradesChan
		if !ok {
			// channel closed
			return
		}
		fmt.Printf("response %v received\n", i)
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

func orderLifecycleTest(g *provider.GRPCClient, ownerAddr string, ooAddr string) {
	fmt.Println("\nstarting order lifecycle test")
	fmt.Println()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan *pb.GetOrderStatusStreamResponse)
	go func() {
		stream, err := g.GetOrderStatusStream(ctx, marketAddr, ownerAddr)
		if err != nil {
			log.Fatalf("error getting order status stream %v", err)
		}
		stream.Into(ch)
	}()

	time.Sleep(time.Second * 10)

	clientID := callPlaceOrderGRPC(g, ownerAddr, ooAddr)

	select {
	case update := <-ch:
		if update.OrderInfo.OrderStatus == pb.OrderStatus_OS_OPEN {
			log.Infof("order went to orderbook (`OPEN`) successfully")
		} else {
			log.Errorf("order should be `OPEN` but is %s", update.OrderInfo.OrderStatus.String())
		}
	case <-time.After(time.Second * 30):
		log.Error("no updates after placing order")
		return
	}

	fmt.Println()
	time.Sleep(time.Second * 10)

	callCancelByClientOrderIDGRPC(g, ownerAddr, ooAddr, clientID)

	select {
	case update := <-ch:
		if update.OrderInfo.OrderStatus == pb.OrderStatus_OS_CANCELLED {
			log.Infof("order cancelled (`CANCELLED`) successfully")
		} else {
			log.Errorf("order should be `CANCELLED` but is %s", update.OrderInfo.OrderStatus.String())
		}
	case <-time.After(time.Second * 30):
		log.Error("no updates after cancelling order")
		return
	}

	fmt.Println()
	callPostSettleGRPC(g, ownerAddr, ooAddr)
}

func callPlaceOrderGRPC(g *provider.GRPCClient, ownerAddr, ooAddr string) uint64 {
	fmt.Println("starting place order")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// generate a random clientOrderID for this order
	rand.Seed(time.Now().UnixNano())
	clientOrderID := rand.Uint64()

	opts := provider.PostOrderOpts{
		ClientOrderID:     clientOrderID,
		OpenOrdersAddress: ooAddr,
	}

	// create order without actually submitting
	response, err := g.PostOrder(ctx, ownerAddr, ownerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Fatalf("failed to create order (%v)", err)
	}
	fmt.Printf("created unsigned place order transaction: %v\n", response.Transaction)

	// sign/submit transaction after creation
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

	fmt.Printf("response signature received: %v\n", sig)
}

func cancelAll(g *provider.GRPCClient, ownerAddr, payerAddr, ooAddr string) {
	fmt.Println("\nstarting cancel all test")
	fmt.Println()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rand.Seed(time.Now().UnixNano())
	clientOrderID1 := rand.Uint64()
	clientOrderID2 := rand.Uint64()
	opts := provider.PostOrderOpts{
		ClientOrderID:     clientOrderID1,
		OpenOrdersAddress: ooAddr,
		SkipPreFlight:     true,
	}

	// Place 2 orders in orderbook
	fmt.Println("placing orders")
	sig, err := g.SubmitOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("submitting place order #1, signature %s", sig)

	opts.ClientOrderID = clientOrderID2
	sig, err = g.SubmitOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check orders are there
	orders, err := g.GetOpenOrders(ctx, marketAddr, ownerAddr)
	if err != nil {
		log.Fatal(err)
	}
	found1 := false
	found2 := false

	for _, order := range orders.Orders {
		if order.ClientOrderID == fmt.Sprintf("%v", clientOrderID1) {
			found1 = true
			continue
		}
		if order.ClientOrderID == fmt.Sprintf("%v", clientOrderID2) {
			found2 = true
		}
	}
	if !(found1 && found2) {
		log.Fatal("one/both orders not found in orderbook")
	}
	fmt.Println("2 orders placed successfully")

	// Cancel all the orders
	fmt.Println("\ncancelling the orders")
	sigs, err := g.SubmitCancelAll(ctx, marketAddr, ownerAddr, []string{ooAddr}, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("placing cancel order(s) %s", strings.Join(sigs, ", "))

	time.Sleep(time.Second * 30)

	orders, err = g.GetOpenOrders(ctx, marketAddr, ownerAddr)
	if err != nil {
		log.Fatal(err)
	}
	if len(orders.Orders) != 0 {
		log.Errorf("%v orders in ob not cancelled", len(orders.Orders))
		return
	}
	fmt.Println("orders cancelled")

	fmt.Println()
	callPostSettleGRPC(g, ownerAddr, ooAddr)
}

func callReplaceByClientOrderID(g *provider.GRPCClient, ownerAddr, payerAddr, ooAddr string) {
	fmt.Println("\nstarting replace by client order ID test")
	fmt.Println()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rand.Seed(time.Now().UnixNano())
	clientOrderID1 := rand.Uint64()
	opts := provider.PostOrderOpts{
		ClientOrderID:     clientOrderID1,
		OpenOrdersAddress: ooAddr,
		SkipPreFlight:     true,
	}

	// Place order in orderbook
	fmt.Println("placing order")
	sig, err := g.SubmitOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("submitting place order #1, signature %s", sig)
	time.Sleep(time.Minute * 1)
	// Check order is there
	orders, err := g.GetOpenOrders(ctx, marketAddr, ownerAddr)
	if err != nil {
		log.Fatal(err)
	}
	found1 := false

	for _, order := range orders.Orders {
		if order.ClientOrderID == fmt.Sprintf("%v", clientOrderID1) {
			found1 = true
			break
		}
	}
	if !found1 {
		log.Fatal("order not found in orderbook")
	}
	fmt.Println("order placed successfully")

	// replacing order
	sig, err = g.SubmitReplaceByClientOrderID(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice/2, opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check order #2 is in orderbook
	orders, err = g.GetOpenOrders(ctx, marketAddr, ownerAddr)
	if err != nil {
		log.Fatal(err)
	}
	found2 := false

	for _, order := range orders.Orders {
		if order.ClientOrderID == fmt.Sprintf("%v", clientOrderID1) && order.Price == orderPrice/2 {
			found2 = true
			break
		}
	}
	if !(found2) {
		log.Fatal("order #2 not found in orderbook")
	}
	fmt.Println("order #2 placed successfully")

	// Cancel all the orders
	fmt.Println("\ncancelling the orders")
	sigs, err := g.SubmitCancelAll(ctx, marketAddr, ownerAddr, []string{ooAddr}, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("placing cancel order(s) %s", strings.Join(sigs, ", "))
}

func callReplaceOrder(g *provider.GRPCClient, ownerAddr, payerAddr, ooAddr string) {
	fmt.Println("\nstarting replace order test")
	fmt.Println()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rand.Seed(time.Now().UnixNano())
	clientOrderID1 := rand.Uint64()
	clientOrderID2 := rand.Uint64()
	opts := provider.PostOrderOpts{
		ClientOrderID:     clientOrderID1,
		OpenOrdersAddress: ooAddr,
		SkipPreFlight:     true,
	}

	// Place order in orderbook
	fmt.Println("placing order")
	sig, err := g.SubmitOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("submitting place order #1, signature %s", sig)
	time.Sleep(time.Minute)
	// Check orders are there
	orders, err := g.GetOpenOrders(ctx, marketAddr, ownerAddr)
	if err != nil {
		log.Fatal(err)
	}
	var found1 *pb.Order

	for _, order := range orders.Orders {
		if order.ClientOrderID == fmt.Sprintf("%v", clientOrderID1) {
			found1 = order
			break
		}
	}
	if found1 == nil {
		log.Fatal("order not found in orderbook")
	} else {
		fmt.Println("order placed successfully")
	}

	opts.ClientOrderID = clientOrderID2
	sig, err = g.SubmitReplaceOrder(ctx, found1.OrderID, ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice/2, opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check orders are there
	orders, err = g.GetOpenOrders(ctx, marketAddr, ownerAddr)
	if err != nil {
		log.Fatal(err)
	}
	var found2 *pb.Order

	for _, order := range orders.Orders {
		if order.ClientOrderID == fmt.Sprintf("%v", clientOrderID2) {
			found2 = order
			break
		}
	}
	if found2 == nil {
		log.Fatal("order 2 not found in orderbook")
	} else {
		fmt.Println("order 2 placed successfully")
	}

	// Cancel all the orders
	fmt.Println("\ncancelling the orders")
	sigs, err := g.SubmitCancelAll(ctx, marketAddr, ownerAddr, []string{ooAddr}, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("placing cancel order(s) %s", strings.Join(sigs, ", "))
}
