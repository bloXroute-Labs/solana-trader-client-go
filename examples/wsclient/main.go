package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bloXroute-Labs/serum-client-go/bxserum/provider"
	api "github.com/bloXroute-Labs/serum-client-go/proto"
	pb "github.com/bloXroute-Labs/serum-client-go/proto"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func main() {
	w, err := provider.NewWSClientTestnet()
	if err != nil {
		log.Fatalf("error dialing WS client: %v", err)
		return
	}
	defer func(w *provider.WSClient) {
		err := w.Close()
		if err != nil {
			panic(err)
		}
	}(w)

	// informational requests
	callMarketsWS(w)
	callOrderbookWS(w)
	callTradesWS(w)
	callOpenOrdersWS(w)
	callTickersWS(w)
	callUnsettledWS(w)
	callAccountBalanceWS(w)

	// streaming methods
	callOrderbookWSStream(w)
	callTradesWSStream(w)

	// calls below this place an order and immediately cancel it
	// you must specify:
	//	- PRIVATE_KEY (by default loaded during provider.NewWSClient()) to sign transactions
	// 	- PUBLIC_KEY to indicate which account you wish to trade from
	//	- OPEN_ORDERS to indicate your Serum account to speed up lookups (optional in actual usage)
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

	orderLifecycleTest(w, ownerAddr, ooAddr)
	cancelAll(w, ownerAddr, payerAddr, ooAddr)
	callReplaceByClientOrderID(w, ownerAddr, payerAddr, ooAddr)
	callReplaceOrder(w, ownerAddr, payerAddr, ooAddr)
}

func callMarketsWS(w *provider.WSClient) {
	fmt.Println("fetching markets...")

	markets, err := w.GetMarkets(context.Background())
	if err != nil {
		log.Errorf("error with GetMarkets request: %v", err)
	} else {
		fmt.Println(markets)
	}

	fmt.Println()
}

func callOrderbookWS(w *provider.WSClient) {
	fmt.Println("fetching orderbooks...")

	orderbook, err := w.GetOrderbook(context.Background(), "ETH-USDT", 0)
	if err != nil {
		log.Errorf("error with GetOrderbook request for ETH-USDT: %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()

	orderbook, err = w.GetOrderbook(context.Background(), "SOLUSDT", 2)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL-USDT: %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()

	orderbook, err = w.GetOrderbook(context.Background(), "SOL:USDC", 3)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL:USDC: %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()
}

func callTradesWS(w *provider.WSClient) {
	fmt.Println("fetching trades...")

	trades, err := w.GetTrades(context.Background(), "SOLUSDC", 3)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL:USDC: %v", err)
	} else {
		fmt.Println(trades)
	}

	fmt.Println()
}

func callOpenOrdersWS(w *provider.WSClient) {
	fmt.Println("fetching open orders...")

	orders, err := w.GetOpenOrders(context.Background(), "SOLUSDC", "FFqDwRq8B4hhFKRqx7N1M6Dg6vU699hVqeynDeYJdPj5")
	if err != nil {
		log.Errorf("error with GetOrders request for SOL-USDT: %v", err)
	} else {
		fmt.Println(orders)
	}

	fmt.Println()
}

func callUnsettledWS(w *provider.WSClient) {
	fmt.Println("fetching unsettled...")

	response, err := w.GetUnsettled(context.Background(), "SOLUSDC", "AFT8VayE7qr8MoQsW3wHsDS83HhEvhGWdbNSHRKeUDfQ")
	if err != nil {
		log.Errorf("error with GetOrders request for SOL-USDT: %v", err)
	} else {
		fmt.Println(response)
	}

	fmt.Println()
}

func callAccountBalanceWS(w *provider.WSClient) {
	fmt.Println("fetching balances...")

	response, err := w.GetAccountBalance(context.Background(), "AFT8VayE7qr8MoQsW3wHsDS83HhEvhGWdbNSHRKeUDfQ")
	if err != nil {
		log.Errorf("error with GetAccountBalance request for AFT8VayE7qr8MoQsW3wHsDS83HhEvhGWdbNSHRKeUDfQ: %v", err)
	} else {
		fmt.Println(response)
	}

	fmt.Println()
}

func callTickersWS(w *provider.WSClient) {
	fmt.Println("fetching tickers...")

	tickers, err := w.GetTickers(context.Background(), "SOLUSDC")
	if err != nil {
		log.Errorf("error with GetTickers request for SOL-USDT: %v", err)
	} else {
		fmt.Println(tickers)
	}

	fmt.Println()
}

// Stream response
func callOrderbookWSStream(w *provider.WSClient) {
	fmt.Println("starting orderbook stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := w.GetOrderbooksStream(ctx, []string{"SOL/USDC"}, 3)
	if err != nil {
		log.Errorf("error with GetOrderbooksStream request for SOL/USDC: %v", err)
	}

	orderbookCh := stream.Channel(0)
	for i := 1; i <= 5; i++ {
		_, ok := <-orderbookCh
		if !ok {
			return
		}
		fmt.Printf("response %v received\n", i)
	}
}

func callTradesWSStream(w *provider.WSClient) {
	fmt.Println("starting trades stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tradesChan := make(chan *pb.GetTradesStreamResponse)
	stream, err := w.GetTradesStream(ctx, "SOL/USDC", 3)
	if err != nil {
		log.Errorf("error with GetTradesStream request for SOL/USDC: %v", err)
	}

	stream.Into(tradesChan)
	for i := 1; i <= 5; i++ {
		_, ok := <-tradesChan
		if !ok {
			return
		}
		fmt.Printf("response %v received\n", i)
	}
}

const (
	// SOL/USDC market
	marketAddr = "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT"

	orderSide   = api.Side_S_ASK
	orderType   = api.OrderType_OT_LIMIT
	orderPrice  = float64(170200)
	orderAmount = float64(0.1)
)

func orderLifecycleTest(w *provider.WSClient, ownerAddr, ooAddr string) {
	fmt.Println("\nstarting order lifecycle test")
	fmt.Println()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan *pb.GetOrderStatusStreamResponse)
	go func() {
		stream, err := w.GetOrderStatusStream(ctx, marketAddr, ownerAddr)
		if err != nil {
			log.Fatalf("error getting order status stream %v", err)
		}
		stream.Into(ch)
	}()

	time.Sleep(time.Second * 10)

	clientOrderID := callPlaceOrderWS(w, ownerAddr, ooAddr)

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

	callCancelByClientOrderIDWS(w, ownerAddr, ooAddr, clientOrderID)

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
	callPostSettleWS(w, ownerAddr, ooAddr)
}

func callPlaceOrderWS(w *provider.WSClient, ownerAddr, ooAddr string) uint64 {
	fmt.Println("trying to place an order")

	// generate a random clientOrderId for this order
	rand.Seed(time.Now().UnixNano())
	clientOrderID := rand.Uint64()

	opts := provider.PostOrderOpts{
		ClientOrderID:     clientOrderID,
		OpenOrdersAddress: ooAddr,
	}

	// create order without actually submitting
	response, err := w.PostOrder(context.Background(), ownerAddr, ownerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Fatalf("failed to create order (%v)", err)
	}
	fmt.Printf("created unsigned place order transaction: %v\n", response.Transaction)

	// sign/submit transaction after creation
	sig, err := w.SubmitOrder(context.Background(), ownerAddr, ownerAddr, marketAddr,
		orderSide, []api.OrderType{orderType}, orderAmount,
		orderPrice, opts)
	if err != nil {
		log.Fatalf("failed to submit order (%v)", err)
	}

	fmt.Printf("placed order %v with clientOrderID %v\n", sig, clientOrderID)

	return clientOrderID
}

func callCancelByClientOrderIDWS(w *provider.WSClient, ownerAddr, ooAddr string, clientOrderID uint64) {
	fmt.Println("trying to cancel order")

	_, err := w.SubmitCancelByClientOrderID(context.Background(), clientOrderID, ownerAddr,
		marketAddr, ooAddr, true)
	if err != nil {
		log.Fatalf("failed to cancel order by client ID (%v)", err)
	}

	fmt.Printf("canceled order for clientOrderID %v\n", clientOrderID)
}

func callPostSettleWS(w *provider.WSClient, ownerAddr, ooAddr string) {
	fmt.Println("starting post settle")

	sig, err := w.SubmitSettle(context.Background(), ownerAddr, "SOL/USDC", "F75gCEckFAyeeCWA9FQMkmLCmke7ehvBnZeVZ3QgvJR7", "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv", ooAddr, false)
	if err != nil {
		log.Errorf("error with post transaction stream request for SOL/USDC: %v", err)
		return
	}

	fmt.Printf("response signature received: %v\n", sig)
}

func cancelAll(w *provider.WSClient, ownerAddr, payerAddr, ooAddr string) {
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
	sig, err := w.SubmitOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("submitting place order #1, signature %s", sig)

	opts.ClientOrderID = clientOrderID2
	sig, err = w.SubmitOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check orders are there
	orders, err := w.GetOpenOrders(ctx, marketAddr, ownerAddr)
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
	sigs, err := w.SubmitCancelAll(ctx, marketAddr, ownerAddr, []string{ooAddr}, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("placing cancel order(s) %s", strings.Join(sigs, ", "))

	time.Sleep(time.Minute)

	orders, err = w.GetOpenOrders(ctx, marketAddr, ownerAddr)
	if err != nil {
		log.Fatal(err)
	}
	if len(orders.Orders) != 0 {
		log.Errorf("%v orders in ob not cancelled", len(orders.Orders))
		return
	}
	fmt.Println("orders cancelled")

	fmt.Println()
	callPostSettleWS(w, ownerAddr, ooAddr)
}

func callReplaceByClientOrderID(w *provider.WSClient, ownerAddr, payerAddr, ooAddr string) {
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
	sig, err := w.SubmitOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Infof("submitting place order #1, signature %s", sig)
	}
	time.Sleep(time.Minute)
	// Check order is there
	orders, err := w.GetOpenOrders(ctx, marketAddr, ownerAddr)
	if err != nil {
		log.Fatal(err)
	}
	found1 := false

	for _, order := range orders.Orders {
		if order.ClientOrderID == fmt.Sprintf("%v", clientOrderID1) {
			found1 = true
			continue
		}
	}
	if !(found1) {
		log.Fatal("order not found in orderbook")
	}
	fmt.Println("order placed successfully")

	// replacing order
	sig, err = w.SubmitReplaceByClientOrderID(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice/2, opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check order #2 is in orderbook
	orders, err = w.GetOpenOrders(ctx, marketAddr, ownerAddr)
	if err != nil {
		log.Fatal(err)
	}
	found2 := false

	for _, order := range orders.Orders {
		if order.ClientOrderID == fmt.Sprintf("%v", clientOrderID1) && order.Price == orderPrice/2 {
			found2 = true
		}
	}
	if !(found2) {
		log.Fatal("order #2 not found in orderbook")
	} else {
		fmt.Println("order #2 placed successfully")
	}
	time.Sleep(time.Minute)
	// Cancel all the orders
	fmt.Println("\ncancelling the orders")
	sigs, err := w.SubmitCancelAll(ctx, marketAddr, ownerAddr, []string{ooAddr}, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("placing cancel order(s) %s", strings.Join(sigs, ", "))
}

func callReplaceOrder(w *provider.WSClient, ownerAddr, payerAddr, ooAddr string) {
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
	sig, err := w.SubmitOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("submitting place order #1, signature %s", sig)

	// Check orders are there
	orders, err := w.GetOpenOrders(ctx, marketAddr, ownerAddr)
	if err != nil {
		log.Fatal(err)
	}
	var found1 *pb.Order

	for _, order := range orders.Orders {
		if order.ClientOrderID == fmt.Sprintf("%v", clientOrderID1) {
			found1 = order
			continue
		}
	}
	if found1 == nil {
		log.Fatal("order not found in orderbook")
	} else {
		fmt.Println("order placed successfully")
	}

	opts.ClientOrderID = clientOrderID2
	sig, err = w.SubmitReplaceOrder(ctx, found1.OrderID, ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice/2, opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check orders are there
	orders, err = w.GetOpenOrders(ctx, marketAddr, ownerAddr)
	if err != nil {
		log.Fatal(err)
	}
	var found2 *pb.Order

	for _, order := range orders.Orders {
		if order.ClientOrderID == fmt.Sprintf("%v", clientOrderID2) {
			found2 = order
		}
	}
	if found2 == nil {
		log.Fatal("order 2 not found in orderbook")
	} else {
		fmt.Println("order 2 placed successfully")
	}

	// Cancel all the orders
	fmt.Println("\ncancelling the orders")
	sigs, err := w.SubmitCancelAll(ctx, marketAddr, ownerAddr, []string{ooAddr}, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("placing cancel order(s) %s", strings.Join(sigs, ", "))
}
