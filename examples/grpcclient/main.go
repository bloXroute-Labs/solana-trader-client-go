package main

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/examples/config"

	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	"github.com/gagliardetto/solana-go"
	"math/rand"
	"os"
	"time"

	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	log "github.com/sirupsen/logrus"
)

func main() {
	utils.InitLogger()
	os.Setenv("PUBLIC_KEY", "2JJQHAYdogfB1fE1ftcvFcsQAXSgQQKkafCwZczWdSWd")
	os.Setenv("PRIVATE_KEY", "3EhZ4Epe6QrcDKQRucdftv6vWXMnpTKDV4mekSPWZEcZnJV4huzesLHwASdVUzoGyQ8evywwomGHQZiYr91fdm6y")
	os.Setenv("AUTH_HEADER", "ZDIxYzE0NmItZWYxNi00ZmFmLTg5YWUtMzYwMTk4YzUyZmM4OjEwOWE5MzEzZDc2Yjg3MzczYjdjZDdhNmZkZGE3ZDg5")
	os.Setenv("API_ENV", "local")
	failed := run()
	if failed {
		log.Fatal("one or multiple examples failed")
	}
}

func run() bool {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	var g *provider.GRPCClient

	switch cfg.Env {
	case config.EnvLocal:
		g, err = provider.NewGRPCLocal()
	case config.EnvTestnet:
		g, err = provider.NewGRPCTestnet()
	case config.EnvMainnet:

		g, err = provider.NewGRPCClient()
	}
	if err != nil {
		log.Fatalf("error dialing GRPC client: %v", err)
	}

	var failed bool

	// informational methods
	//failed = failed || callMarketsGRPC(g)
	//failed = failed || callOrderbookGRPC(g)
	//failed = failed || callOpenOrdersGRPC(g)
	//failed = failed || callTickersGRPC(g)
	//failed = failed || callPoolsGRPC(g)
	//failed = failed || callPriceGRPC(g)
	failed = failed || callOrderbookGRPCStream(g)
	failed = failed || callPricesGRPCStream(g)

	// trade stream can be slow
	if cfg.RunTradeStream {
		failed = failed || callTradesGRPCStream(g)
	}

	failed = failed || callUnsettledGRPC(g)
	failed = failed || callGetAccountBalanceGRPC(g)
	failed = failed || callGetQuotes(g)
	failed = failed || callRecentBlockHashGRPCStream(g)
	failed = failed || callQuotesGRPCStream(g)
	failed = failed || callPoolReservesGRPCStream(g)
	failed = failed || callSwapsGRPCStream(g)

	if !cfg.RunTrades {
		log.Info("skipping trades due to config")
		return failed
	}

	// calls below this place an order and immediately cancel it
	// you must specify:
	//	- PRIVATE_KEY (by default loaded during provider.NewGRPCClient()) to sign transactions
	// 	- PUBLIC_KEY to indicate which account you wish to trade from
	//	- OPEN_ORDERS to indicate your Serum account to speed up lookups (optional)
	//ownerAddr, ok := os.LookupEnv("PUBLIC_KEY")
	//if !ok {
	//	log.Infof("PUBLIC_KEY environment variable not set: will skip place/cancel/settle examples")
	//	return failed
	//}
	//
	//ooAddr, ok := os.LookupEnv("OPEN_ORDERS")
	//if !ok {
	//	log.Infof("OPEN_ORDERS environment variable not set: requests will be slower")
	//}

	//payerAddr, ok := os.LookupEnv("PAYER")
	//if !ok {
	//	log.Infof("PAYER environment variable not set: will be set to owner address")
	//	payerAddr = ownerAddr
	//}
	//if !ok {
	//	log.Infof("OPEN_ORDERS environment variable not set: requests will be slower")
	//}

	//failed = failed || orderLifecycleTest(g, ownerAddr, payerAddr, ooAddr)
	//failed = failed || cancelAll(g, ownerAddr, payerAddr, ooAddr)
	//failed = failed || callReplaceByClientOrderID(g, ownerAddr, payerAddr, ooAddr)
	//failed = failed || callReplaceOrder(g, ownerAddr, payerAddr, ooAddr)
	//failed = failed || callTradeSwap(g, ownerAddr)
	//failed = failed || callRouteTradeSwap(g, ownerAddr)
	//failed = failed || callAddMemoWithInstructions(g, ownerAddr)
	//failed = failed || callAddMemoToSerializedTxn(g, ownerAddr) //failed = failed || orderLifecycleTest(g, ownerAddr, payerAddr, ooAddr)
	//failed = failed || cancelAll(g, ownerAddr, payerAddr, ooAddr)
	//failed = failed || callReplaceByClientOrderID(g, ownerAddr, payerAddr, ooAddr)
	//failed = failed || callReplaceOrder(g, ownerAddr, payerAddr, ooAddr)
	//failed = failed || callTradeSwap(g, ownerAddr)
	//failed = failed || callRouteTradeSwap(g, ownerAddr)
	//failed = failed || callAddMemoWithInstructions(g, ownerAddr)
	//failed = failed || callAddMemoToSerializedTxn(g, ownerAddr)

	return failed
}

func callMarketsGRPC(g *provider.GRPCClient) bool {
	markets, err := g.GetMarkets(context.Background())
	if err != nil {
		log.Errorf("error with GetMarkets request: %v", err)
		return true
	} else {
		log.Info(markets)
	}

	fmt.Println()
	return false
}

func callOrderbookGRPC(g *provider.GRPCClient) bool {
	orderbook, err := g.GetOrderbook(context.Background(), "ETH-USDT", 0)
	if err != nil {
		log.Errorf("error with GetOrderbook request for ETH-USDT: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()

	orderbook, err = g.GetOrderbook(context.Background(), "SOLUSDT", 2)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOLUSDT: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()

	orderbook, err = g.GetOrderbook(context.Background(), "SOL:USDC", 3)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL:USDC: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()
	return false
}

func callOpenOrdersGRPC(g *provider.GRPCClient) bool {
	orders, err := g.GetOpenOrders(context.Background(), "SOLUSDC", "FFqDwRq8B4hhFKRqx7N1M6Dg6vU699hVqeynDeYJdPj5", "")
	if err != nil {
		log.Errorf("error with GetOrders request for SOLUSDC: %v", err)
		return true
	} else {
		log.Info(orders)
	}

	fmt.Println()
	return false
}

func callUnsettledGRPC(g *provider.GRPCClient) bool {
	response, err := g.GetUnsettled(context.Background(), "SOLUSDC", "HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc")
	if err != nil {
		log.Errorf("error with GetUnsettled request for SOLUSDC: %v", err)
		return true
	} else {
		log.Info(response)
	}

	fmt.Println()
	return false
}

func callGetAccountBalanceGRPC(g *provider.GRPCClient) bool {
	response, err := g.GetAccountBalance(context.Background(), "HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc")
	if err != nil {
		log.Errorf("error with GetAccountBalance request for HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc: %v", err)
		return true
	} else {
		log.Info(response)
	}

	fmt.Println()
	return false
}

func callTickersGRPC(g *provider.GRPCClient) bool {
	orders, err := g.GetTickers(context.Background(), "SOLUSDC")
	if err != nil {
		log.Errorf("error with GetTickers request for SOLUSDC: %v", err)
		return true
	} else {
		log.Info(orders)
	}

	fmt.Println()
	return false
}

func callPoolsGRPC(g *provider.GRPCClient) bool {
	pools, err := g.GetPools(context.Background(), []pb.Project{pb.Project_P_RAYDIUM})
	if err != nil {
		log.Errorf("error with GetPools request for Raydium: %v", err)
		return true
	} else {
		log.Info(pools)
	}

	fmt.Println()
	return false
}

func callPriceGRPC(g *provider.GRPCClient) bool {
	prices, err := g.GetPrice(context.Background(), []string{"SOL", "ETH"})
	if err != nil {
		log.Errorf("error with GetPrice request for SOL and ETH: %v", err)
		return true
	} else {
		log.Info(prices)
	}

	fmt.Println()
	return false
}

func callGetQuotes(g *provider.GRPCClient) bool {
	log.Info("starting get quotes test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inToken := "SOL"
	outToken := "USDC"
	amount := 0.01
	slippage := float64(5)
	limit := 5

	quotes, err := g.GetQuotes(ctx, inToken, outToken, amount, slippage, int32(limit), []pb.Project{pb.Project_P_ALL})
	if err != nil {
		log.Errorf("error with GetQuotes request for %s to %s: %v", inToken, outToken, err)
		return true
	}

	if len(quotes.Quotes) != 2 {
		log.Errorf("did not get back 2 quotes, got %v quotes", len(quotes.Quotes))
		return true
	}
	for _, quote := range quotes.Quotes {
		if len(quote.Routes) == 0 {
			log.Errorf("no routes gotten for project %s", quote.Project)
			return true
		} else {
			log.Infof("best route for project %s: %v", quote.Project, quote.Routes[0])
		}
	}

	fmt.Println()
	return false
}

func callOrderbookGRPCStream(g *provider.GRPCClient) bool {
	log.Info("starting orderbook stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream error response
	stream, err := g.GetOrderbookStream(ctx, []string{"SOL/USDC", "xxx"}, 3)
	if err != nil {
		log.Errorf("connection could not be established. error: %v", err)
		return true
	}

	_, err = stream()
	if err != nil {
		//demonstration purposes only. will swallow
		log.Infof("subscription error: %v", err)
	} else {
		log.Error("subscription should have returned an error")
		return true
	}

	// Stream ok response
	stream, err = g.GetOrderbookStream(ctx, []string{"SOL/USDC", "SOL-USDT"}, 3)
	if err != nil {
		log.Errorf("connection could not be established. error: %v", err)
		return true
	}

	_, err = stream()
	if err != nil {
		log.Errorf("subscription error: %v", err)
		return true
	}

	orderbookCh := stream.Channel(0)
	for i := 1; i <= 2; i++ {
		data, ok := <-orderbookCh
		if !ok {
			// channel closed
			return true
		}

		log.Infof("response %v received, data %v ", i, data)
	}

	return false
}

func callTradesGRPCStream(g *provider.GRPCClient) bool {
	log.Info("starting trades stream")

	tradesChan := make(chan *pb.GetTradesStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := g.GetTradesStream(ctx, "SOL/USDC", 3)
	if err != nil {
		log.Errorf("error with GetTrades stream request for SOL/USDC: %v", err)
		return true
	}
	stream.Into(tradesChan)
	for i := 1; i <= 3; i++ {
		_, ok := <-tradesChan
		if !ok {
			// channel closed
			return true
		}
		log.Infof("response %v received", i)
	}
	return false
}

func callRecentBlockHashGRPCStream(g *provider.GRPCClient) bool {
	log.Info("starting recent block hash stream")

	ch := make(chan *pb.GetRecentBlockHashResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := g.GetRecentBlockHashStream(ctx)
	if err != nil {
		log.Errorf("error with GetRecentBlockHash stream request: %v", err)
		return true
	}
	stream.Into(ch)
	for i := 1; i <= 3; i++ {
		_, ok := <-ch
		if !ok {
			// channel closed
			return true
		}
		log.Infof("response %v received", i)
	}
	return false
}

func callQuotesGRPCStream(g *provider.GRPCClient) bool {
	log.Info("starting get quotes stream")

	ch := make(chan *pb.GetQuotesStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := g.GetQuotesStream(ctx, []pb.Project{pb.Project_P_RAYDIUM}, []*pb.TokenPair{{
		InToken:  "SOL",
		OutToken: "USDC",
		InAmount: 1,
	}})

	if err != nil {
		log.Errorf("error with GetQuotes stream request: %v", err)
		return true
	}
	stream.Into(ch)
	for i := 1; i <= 3; i++ {
		_, ok := <-ch
		if !ok {
			// channel closed
			return true
		}
		log.Infof("response %v received", i)
	}
	return false
}

func callPoolReservesGRPCStream(g *provider.GRPCClient) bool {
	log.Info("starting get pool reserves stream")

	ch := make(chan *pb.GetPoolReservesStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := g.GetPoolReservesStream(ctx, []pb.Project{pb.Project_P_RAYDIUM})

	if err != nil {
		log.Errorf("error with GetPoolReserves stream request: %v", err)
		return true
	}
	stream.Into(ch)
	for i := 1; i <= 3; i++ {
		_, ok := <-ch
		if !ok {
			// channel closed
			return true
		}
		log.Infof("response %v received", i)
	}
	return false
}

const (
	// SOL/USDC market
	marketAddr = "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT"

	orderSide   = pb.Side_S_ASK
	orderType   = pb.OrderType_OT_LIMIT
	orderPrice  = float64(170200)
	orderAmount = float64(0.1)
)

func orderLifecycleTest(g *provider.GRPCClient, ownerAddr, payerAddr, ooAddr string) bool {
	log.Info("starting order lifecycle test")
	fmt.Println()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan *pb.GetOrderStatusStreamResponse)
	errCh := make(chan error)
	stream, err := g.GetOrderStatusStream(ctx, marketAddr, ownerAddr)
	if err != nil {
		log.Errorf("error getting order status stream %v", err)
		errCh <- err
	}
	stream.Into(ch)

	time.Sleep(time.Second * 10)

	clientID, failed := callPlaceOrderGRPC(g, ownerAddr, payerAddr, ooAddr)

	if failed {
		return true
	}

	select {
	case update := <-ch:
		if update.OrderInfo.OrderStatus == pb.OrderStatus_OS_OPEN {
			log.Infof("order went to orderbook (`OPEN`) successfully")
		} else {
			log.Errorf("order should be `OPEN` but is %s", update.OrderInfo.OrderStatus.String())
		}
	case <-errCh:
		return true
	case <-time.After(time.Second * 60):
		log.Error("no updates after placing order")
		return true
	}

	fmt.Println()
	time.Sleep(time.Second * 10)

	failed = callCancelByClientOrderIDGRPC(g, ownerAddr, ooAddr, clientID)
	if failed {
		return true
	}

	select {
	case update := <-ch:
		if update.OrderInfo.OrderStatus == pb.OrderStatus_OS_CANCELLED {
			log.Infof("order cancelled (`CANCELLED`) successfully")
		} else {
			log.Errorf("order should be `CANCELLED` but is %s", update.OrderInfo.OrderStatus.String())
		}
	case <-errCh:
		return true
	case <-time.After(time.Second * 30):
		log.Error("no updates after cancelling order")
		return true
	}

	fmt.Println()
	return callPostSettleGRPC(g, ownerAddr, ooAddr)
}

func callPlaceOrderGRPC(g *provider.GRPCClient, ownerAddr, payerAddr, ooAddr string) (uint64, bool) {
	log.Info("starting place order")

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
	response, err := g.PostOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Errorf("failed to create order (%v)", err)
		return 0, true
	}
	log.Infof("created unsigned place order transaction: %v", response.Transaction)

	// sign/submit transaction after creation
	sig, err := g.SubmitOrder(ctx, ownerAddr, ownerAddr, marketAddr,
		orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Errorf("failed to submit order (%v)", err)
		return 0, true
	}

	log.Infof("placed order %v with clientOrderID %v", sig, clientOrderID)
	return clientOrderID, false
}

func callCancelByClientOrderIDGRPC(g *provider.GRPCClient, ownerAddr, ooAddr string, clientID uint64) bool {
	log.Info("starting cancel order by client order ID")
	time.Sleep(30 * time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig, err := g.SubmitCancelByClientOrderID(ctx, clientID, ownerAddr,
		marketAddr, ooAddr, true)
	if err != nil {
		log.Errorf("failed to cancel order by client order ID (%v)", err)
		return true
	}

	log.Infof("canceled order %v with clientOrderID %v", sig, clientID)
	return false
}

func callPostSettleGRPC(g *provider.GRPCClient, ownerAddr, ooAddr string) bool {
	log.Info("starting post settle")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig, err := g.SubmitSettle(ctx, ownerAddr, "SOL/USDC", "F75gCEckFAyeeCWA9FQMkmLCmke7ehvBnZeVZ3QgvJR7", "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv", ooAddr, false)
	if err != nil {
		log.Errorf("error with post transaction stream request for SOL/USDC: %v", err)
		return true
	}

	log.Infof("response signature received: %v", sig)
	return false
}

func cancelAll(g *provider.GRPCClient, ownerAddr, payerAddr, ooAddr string) bool {
	log.Info("starting cancel all test")
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
	log.Info("placing orders")
	sig, err := g.SubmitOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #1, signature %s", sig)

	opts.ClientOrderID = clientOrderID2
	sig, err = g.SubmitOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check orders are there
	orders, err := g.GetOpenOrders(ctx, marketAddr, ownerAddr, "")
	if err != nil {
		log.Error(err)
		return true
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
		log.Error("one/both orders not found in orderbook")
		return true
	}
	log.Info("2 orders placed successfully")

	// Cancel all the orders
	log.Info("cancelling the orders")
	sigs, err := g.SubmitCancelAll(ctx, marketAddr, ownerAddr, []string{ooAddr}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  true,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	for _, tx := range sigs.Transactions {
		log.Infof("placing cancel order(s) %s", tx.Signature)
	}

	time.Sleep(time.Second * 30)

	orders, err = g.GetOpenOrders(ctx, marketAddr, ownerAddr, "")
	if err != nil {
		log.Error(err)
		return true
	}
	if len(orders.Orders) != 0 {
		log.Errorf("%v orders in ob not cancelled", len(orders.Orders))
		return true
	}
	log.Info("orders cancelled")

	fmt.Println()
	return callPostSettleGRPC(g, ownerAddr, ooAddr)
}

func callReplaceByClientOrderID(g *provider.GRPCClient, ownerAddr, payerAddr, ooAddr string) bool {
	log.Info("starting replace by client order ID test")
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
	log.Info("placing order")
	sig, err := g.SubmitOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #1, signature %s", sig)
	time.Sleep(time.Minute * 1)

	// Check order is there
	orders, err := g.GetOpenOrders(ctx, marketAddr, ownerAddr, "")
	if err != nil {
		log.Error(err)
		return true
	}
	found1 := false

	for _, order := range orders.Orders {
		if order.ClientOrderID == fmt.Sprintf("%v", clientOrderID1) {
			found1 = true
			break
		}
	}
	if !found1 {
		log.Error("order not found in orderbook")
		return true
	}
	log.Info("order placed successfully")

	// replacing order
	sig, err = g.SubmitReplaceByClientOrderID(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice/2, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check order #2 is in orderbook
	orders, err = g.GetOpenOrders(ctx, marketAddr, ownerAddr, "")
	if err != nil {
		log.Error(err)
		return true
	}
	found2 := false

	for _, order := range orders.Orders {
		if order.ClientOrderID == fmt.Sprintf("%v", clientOrderID1) && order.Price == orderPrice/2 {
			found2 = true
			break
		}
	}
	if !(found2) {
		log.Error("order #2 not found in orderbook")
		return true
	}
	log.Info("order #2 placed successfully")

	// Cancel all the orders
	log.Info("cancelling the orders")
	sigs, err := g.SubmitCancelAll(ctx, marketAddr, ownerAddr, []string{ooAddr}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  true,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	for _, tx := range sigs.Transactions {
		log.Infof("placing cancel order(s) %s", tx.Signature)
	}
	return false
}

func callReplaceOrder(g *provider.GRPCClient, ownerAddr, payerAddr, ooAddr string) bool {
	log.Info("starting replace order test")
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
	log.Info("placing order")
	sig, err := g.SubmitOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #1, signature %s", sig)
	time.Sleep(time.Minute)
	// Check orders are there
	orders, err := g.GetOpenOrders(ctx, marketAddr, ownerAddr, "")
	if err != nil {
		log.Error(err)
		return true
	}
	var found1 *pb.Order

	for _, order := range orders.Orders {
		if order.ClientOrderID == fmt.Sprintf("%v", clientOrderID1) {
			found1 = order
			break
		}
	}
	if found1 == nil {
		log.Error("order not found in orderbook")
		return true
	} else {
		log.Info("order placed successfully")
	}

	opts.ClientOrderID = clientOrderID2
	sig, err = g.SubmitReplaceOrder(ctx, found1.OrderID, ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice/2, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check orders are there
	orders, err = g.GetOpenOrders(ctx, marketAddr, ownerAddr, "")
	if err != nil {
		log.Error(err)
		return true
	}
	var found2 *pb.Order

	for _, order := range orders.Orders {
		if order.ClientOrderID == fmt.Sprintf("%v", clientOrderID2) {
			found2 = order
			break
		}
	}
	if found2 == nil {
		log.Error("order 2 not found in orderbook")
		return true
	} else {
		log.Info("order 2 placed successfully")
	}

	// Cancel all the orders
	log.Info("cancelling the orders")
	sigs, err := g.SubmitCancelAll(ctx, marketAddr, ownerAddr, []string{ooAddr}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  true,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	for _, tx := range sigs.Transactions {
		log.Infof("placing cancel order(s) %s", tx.Signature)
	}
	return false
}

func callTradeSwap(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting trade swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("trade swap")
	sig, err := g.SubmitTradeSwap(ctx, ownerAddr, "USDC",
		"SOL", 0.01, 0.1, pb.Project_P_RAYDIUM, provider.SubmitOpts{
			SubmitStrategy: pb.SubmitStrategy_P_ABORT_ON_FIRST_ERROR,
			SkipPreFlight:  false,
		})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("trade swap transaction signature : %s", sig)
	return false
}

func callRouteTradeSwap(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting route trade swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("route trade swap")
	sig, err := g.SubmitRouteTradeSwap(ctx, &pb.RouteTradeSwapRequest{
		OwnerAddress: ownerAddr,
		Project:      pb.Project_P_RAYDIUM,
		Steps: []*pb.RouteStep{
			{
				// FIDA-RAY pool address
				InToken:  "FIDA",
				OutToken: "RAY",

				InAmount:     0.01,
				OutAmountMin: 0.007505,
				OutAmount:    0.0074,
			},
			{
				// RAY-USDC pool address
				InToken:      "RAY",
				OutToken:     "USDC",
				InAmount:     0.007505,
				OutAmount:    0.004043,
				OutAmountMin: 0.004000,
			},
		},
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_ABORT_ON_FIRST_ERROR,
		SkipPreFlight:  false,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("route trade swap transaction signature : %s", sig)
	return false
}

func callAddMemoWithInstructions(g *provider.GRPCClient, ownerAddr string) bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	privateKey, err := transaction.LoadPrivateKeyFromEnv()
	if err != nil {
		log.Error(err)
		return false
	}
	privateKeys := make(map[solana.PublicKey]solana.PrivateKey)
	privateKeys[privateKey.PublicKey()] = privateKey
	blockHashResp, err := g.RecentBlockHash(ctx)
	if err != nil {
		log.Error(err)
		return false
	}
	encodedTxn, err := transaction.AddMemo(
		[]solana.Instruction{},
		"new memo by dev",
		solana.MustHashFromBase58(blockHashResp.BlockHash),
		solana.MustPublicKeyFromBase58(ownerAddr),
		privateKeys,
	)
	response, err := g.PostSubmit(ctx, &pb.TransactionMessage{Content: encodedTxn}, false)
	if err != nil {
		log.Error(err)
		return false
	}
	log.Infof("response.signature : %s", response.Signature)
	return true
}

func callAddMemoToSerializedTxn(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("add memo to serialized tx")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	privateKey, err := transaction.LoadPrivateKeyFromEnv()
	if err != nil {
		log.Error(err)
		return false
	}
	privateKeys := make(map[solana.PublicKey]solana.PrivateKey)
	privateKeys[privateKey.PublicKey()] = privateKey
	blockHashResp, err := g.RecentBlockHash(ctx)
	if err != nil {
		log.Error(err)
		return false
	}
	encodedTxn, err := transaction.AddMemo(
		[]solana.Instruction{},
		"new memo by dev",
		solana.MustHashFromBase58(blockHashResp.BlockHash),
		solana.MustPublicKeyFromBase58(ownerAddr),
		privateKeys,
	)

	encodedTxn2, err := transaction.AddMemoToSerializedTxn(encodedTxn, "new memo by dev2", solana.MustPublicKeyFromBase58(ownerAddr), privateKeys)
	if err != nil {
		log.Error(err)
		return false
	}
	log.Infof("encodedTxn2 : %s", encodedTxn2)
	response, err := g.PostSubmit(ctx, &pb.TransactionMessage{Content: encodedTxn2}, false)
	if err != nil {
		log.Error(err)
		return false
	}
	log.Infof("response.signature : %s", response.Signature)

	return true
}

func callPricesGRPCStream(g *provider.GRPCClient) bool {
	log.Info("starting get prices stream")

	ch := make(chan *pb.GetPricesStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := g.GetPricesStream(ctx, []pb.Project{pb.Project_P_RAYDIUM}, []string{"SOL"})

	if err != nil {
		log.Errorf("error with GetPrices stream request: %v", err)
		return true
	}
	stream.Into(ch)
	for i := 1; i <= 3; i++ {
		_, ok := <-ch
		if !ok {
			// channel closed
			return true
		}
		log.Infof("response %v received", i)
	}
	return false
}

func callSwapsGRPCStream(g *provider.GRPCClient) bool {
	log.Info("starting get swaps stream")

	ch := make(chan *pb.GetSwapsStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := g.GetSwapsStream(ctx, []pb.Project{pb.Project_P_RAYDIUM}, []string{"58oQChx4yWmvKdwLLZzBi4ChoCc2fqCUWBkwMihLYQo2"}) // SOL-USDC Raydium pool
	if err != nil {
		log.Errorf("error with GetSwaps stream request: %v", err)
		return true
	}
	stream.Into(ch)
	for i := 1; i <= 3; i++ {
		_, ok := <-ch
		if !ok {
			// channel closed
			return true
		}

		log.Infof("response %v received", i)
	}
	return false
}
