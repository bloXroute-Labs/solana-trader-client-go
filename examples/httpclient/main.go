package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bloXroute-Labs/serum-client-go/bxserum/provider"
	pb "github.com/bloXroute-Labs/serum-client-go/proto"
	log "github.com/sirupsen/logrus"
)

func main() {
	// informational methods
	callMarketsHTTP()
	callOrderbookHTTP()
	callOpenOrdersHTTP()
	callTradesHTTP()
	callTickersHTTP()
	callUnsettledHTTP()
	callGetAccountBalanceHTTP()

	// calls below this place an order and immediately cancel it
	// you must specify:
	//	- PRIVATE_KEY (by default loaded during provider.NewGRPCClient()) to sign transactions
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

	clientOrderID := callPlaceOrderHTTP(ownerAddr, ooAddr)
	callCancelByClientOrderIDHTTP(ownerAddr, ooAddr, clientOrderID)
	callPostSettleHTTP(ownerAddr, ooAddr)

	cancelAll(ownerAddr, payerAddr, ooAddr)
}

func callMarketsHTTP() {
	h := provider.NewHTTPClient()

	markets, err := h.GetMarkets()
	if err != nil {
		log.Errorf("error with GetMarkets request: %v", err)
	} else {
		fmt.Println(markets)
	}

	fmt.Println()
}

func callOrderbookHTTP() {
	h := provider.NewHTTPClient()

	orderbook, err := h.GetOrderbook("ETH-USDT", 0)
	if err != nil {
		log.Errorf("error with GetOrderbook request for ETH-USDT: %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()

	orderbook, err = h.GetOrderbook("SOLUSDT", 2)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOLUSDT: %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()

	orderbook, err = h.GetOrderbook("SOL:USDC", 3)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL:USDC: %v", err)
	} else {
		fmt.Println(orderbook)
	}
}

func callOpenOrdersHTTP() {
	client := &http.Client{Timeout: time.Second * 60}
	opts := provider.DefaultRPCOpts(provider.MainnetSerumAPIHTTP)
	h := provider.NewHTTPClientWithOpts(client, opts)

	orders, err := h.GetOpenOrders("SOLUSDT", "HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc")
	if err != nil {
		log.Errorf("error with GetOrders request for SOLUSDT: %v", err)
	} else {
		fmt.Println(orders)
	}

	fmt.Println()
}

func callUnsettledHTTP() {
	client := &http.Client{Timeout: time.Second * 60}
	opts := provider.DefaultRPCOpts(provider.MainnetSerumAPIHTTP)
	h := provider.NewHTTPClientWithOpts(client, opts)

	response, err := h.GetUnsettled("SOLUSDT", "HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc")
	if err != nil {
		log.Errorf("error with GetOrders request for SOLUSDT: %v", err)
	} else {
		fmt.Println(response)
	}

	fmt.Println()
}

func callGetAccountBalanceHTTP() {
	client := &http.Client{Timeout: time.Second * 60}
	opts := provider.DefaultRPCOpts(provider.MainnetSerumAPIHTTP)
	h := provider.NewHTTPClientWithOpts(client, opts)

	response, err := h.GetAccountBalance("F75gCEckFAyeeCWA9FQMkmLCmke7ehvBnZeVZ3QgvJR7")
	if err != nil {
		log.Errorf("error with GetAccountBalance request for HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc: %v", err)
	} else {
		fmt.Println(response)
	}

	fmt.Println()
}

func callTradesHTTP() {
	h := provider.NewHTTPClient()

	trades, err := h.GetTrades("SOLUSDT", 5)
	if err != nil {
		log.Errorf("error with GetTrades request for SOLUSDT: %v", err)
	} else {
		fmt.Println(trades)
	}

	fmt.Println()
}

func callTickersHTTP() {
	h := provider.NewHTTPClient()

	tickers, err := h.GetTickers("SOLUSDT")
	if err != nil {
		log.Errorf("error with GetTickers request for SOLUSDT: %v", err)
	} else {
		fmt.Println(tickers)
	}

	fmt.Println()
}

const (
	// SOL/USDC market
	marketAddr = "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT"

	orderSide   = pb.Side_S_ASK
	orderType   = pb.OrderType_OT_LIMIT
	orderPrice  = float64(170200)
	orderAmount = float64(0.1)
)

func callPlaceOrderHTTP(ownerAddr, ooAddr string) uint64 {
	client := &http.Client{Timeout: time.Second * 30}
	rpcOpts := provider.DefaultRPCOpts(provider.MainnetSerumAPIHTTP)
	h := provider.NewHTTPClientWithOpts(client, rpcOpts)

	// generate a random clientOrderId for this order
	rand.Seed(time.Now().UnixNano())
	clientOrderID := rand.Uint64()

	opts := provider.PostOrderOpts{
		ClientOrderID:     clientOrderID,
		OpenOrdersAddress: ooAddr,
	}

	// create order without actually submitting
	response, err := h.PostOrder(ownerAddr, ownerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Fatalf("failed to create order (%v)", err)
	}
	fmt.Printf("created unsigned place order transaction: %v", response.Transaction)

	// sign/submit transaction after creation
	sig, err := h.SubmitOrder(ownerAddr, ownerAddr, marketAddr,
		orderSide, []pb.OrderType{orderType}, orderAmount,
		orderPrice, opts)
	if err != nil {
		log.Fatalf("failed to submit order (%v)", err)
	}

	fmt.Printf("placed order %v with clientOrderID %v\n", sig, clientOrderID)

	return clientOrderID
}

func callCancelByClientOrderIDHTTP(ownerAddr, ooAddr string, clientOrderID uint64) {
	time.Sleep(60 * time.Second)
	client := &http.Client{Timeout: time.Second * 30}
	opts := provider.DefaultRPCOpts(provider.MainnetSerumAPIHTTP)
	h := provider.NewHTTPClientWithOpts(client, opts)

	_, err := h.SubmitCancelByClientOrderID(clientOrderID, ownerAddr,
		marketAddr, ooAddr, true)
	if err != nil {
		log.Fatalf("failed to cancel order by client ID (%v)", err)
	}

	fmt.Printf("canceled order for clientOrderID %v\n", clientOrderID)
}

func callPostSettleHTTP(ownerAddr, ooAddr string) {
	time.Sleep(60 * time.Second)
	client := &http.Client{Timeout: time.Second * 30}
	opts := provider.DefaultRPCOpts(provider.MainnetSerumAPIHTTP)
	h := provider.NewHTTPClientWithOpts(client, opts)

	sig, err := h.SubmitSettle(ownerAddr, "SOL/USDC", "F75gCEckFAyeeCWA9FQMkmLCmke7ehvBnZeVZ3QgvJR7", "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv", ooAddr, false)
	if err != nil {
		log.Errorf("error with post transaction stream request for SOL/USDC: %v", err)
		return
	}

	fmt.Printf("response signature received: %v", sig)
}

func cancelAll(owner, payer, ooAddr string) {
	client := &http.Client{Timeout: time.Second * 30}
	rpcOpts := provider.DefaultRPCOpts(provider.MainnetSerumAPIHTTP)
	h := provider.NewHTTPClientWithOpts(client, rpcOpts)

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
	sig, err := h.SubmitOrder(owner, payer, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("submitting place order #1, signature %s", sig)

	opts.ClientOrderID = clientOrderID2
	sig, err = h.SubmitOrder(owner, payer, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check orders are there
	orders, err := h.GetOpenOrders(marketAddr, owner)
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
	fmt.Println("\ncancelling all orders")
	sigs, err := h.SubmitCancelAll(marketAddr, owner, ooAddr, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("placing cancel order(s) %s", strings.Join(sigs, ", "))

	time.Sleep(time.Minute)

	orders, err = h.GetOpenOrders(marketAddr, owner)
	if err != nil {
		log.Fatal(err)
	}
	if len(orders.Orders) != 0 {
		log.Errorf("%v orders in ob not cancelled", len(orders.Orders))
		return
	}
	fmt.Println("all orders in ob cancelled")
}
