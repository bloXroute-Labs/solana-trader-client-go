package main

import (
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	pb "github.com/bloXroute-Labs/solana-trader-client-go/proto"
	log "github.com/sirupsen/logrus"
)

func main() {
	utils.InitLogger()
	var failed bool
	// informational methods
	failed = failed || callMarketsHTTP()
	failed = failed || callOrderbookHTTP()
	failed = failed || callOpenOrdersHTTP()
	failed = failed || callTradesHTTP()
	failed = failed || callPoolsHTTP()
	failed = failed || callTickersHTTP()
	failed = failed || callUnsettledHTTP()
	failed = failed || callGetAccountBalanceHTTP()

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

	// Order lifecycle
	clientOrderID, fail := callPlaceOrderHTTP(ownerAddr, ooAddr)
	failed = failed || fail
	failed = failed || callCancelByClientOrderIDHTTP(ownerAddr, ooAddr, clientOrderID)
	failed = failed || callPostSettleHTTP(ownerAddr, ooAddr)
	failed = failed || cancelAll(ownerAddr, payerAddr, ooAddr)
	failed = failed || callReplaceByClientOrderID(ownerAddr, payerAddr, ooAddr)
	failed = failed || callReplaceOrder(ownerAddr, payerAddr, ooAddr)
	failed = failed || callGetRecentBlockHash()
	failed = failed || callTradeSwap(ownerAddr)

	if failed {
		log.Fatal("one or multiple examples failed")
	}
}

func callMarketsHTTP() bool {
	h := provider.NewHTTPTestnet()

	markets, err := h.GetMarkets()
	if err != nil {
		log.Errorf("error with GetMarkets request: %v", err)
		return true
	} else {
		log.Info(markets)
	}

	fmt.Println()
	return false
}

func callOrderbookHTTP() bool {
	h := provider.NewHTTPTestnet()

	orderbook, err := h.GetOrderbook("ETH-USDT", 0)
	if err != nil {
		log.Errorf("error with GetOrderbook request for ETH-USDT: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()

	orderbook, err = h.GetOrderbook("SOLUSDT", 2)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOLUSDT: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()

	orderbook, err = h.GetOrderbook("SOL:USDC", 3)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL:USDC: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	return false
}

func callOpenOrdersHTTP() bool {
	client := &http.Client{Timeout: time.Second * 60}
	opts := provider.DefaultRPCOpts(provider.TestnetHTTP)
	h := provider.NewHTTPClientWithOpts(client, opts)

	orders, err := h.GetOpenOrders("SOLUSDT", "HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc", "")
	if err != nil {
		log.Errorf("error with GetOrders request for SOLUSDT: %v", err)
		return true
	} else {
		log.Info(orders)
	}

	fmt.Println()
	return false
}

func callUnsettledHTTP() bool {
	client := &http.Client{Timeout: time.Second * 60}
	opts := provider.DefaultRPCOpts(provider.TestnetHTTP)
	h := provider.NewHTTPClientWithOpts(client, opts)

	response, err := h.GetUnsettled("SOLUSDT", "HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc")
	if err != nil {
		log.Errorf("error with GetOrders request for SOLUSDT: %v", err)
		return true
	} else {
		log.Info(response)
	}

	fmt.Println()
	return false
}

func callGetAccountBalanceHTTP() bool {
	client := &http.Client{Timeout: time.Second * 60}
	opts := provider.DefaultRPCOpts(provider.TestnetHTTP)
	h := provider.NewHTTPClientWithOpts(client, opts)

	response, err := h.GetAccountBalance("F75gCEckFAyeeCWA9FQMkmLCmke7ehvBnZeVZ3QgvJR7")
	if err != nil {
		log.Errorf("error with GetAccountBalance request for HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc: %v", err)
		return true
	} else {
		log.Info(response)
	}

	fmt.Println()
	return false
}

func callTradesHTTP() bool {
	h := provider.NewHTTPTestnet()

	trades, err := h.GetTrades("SOLUSDT", 5)
	if err != nil {
		log.Errorf("error with GetTrades request for SOLUSDT: %v", err)
		return true
	} else {
		log.Info(trades)
	}

	fmt.Println()
	return false
}

func callPoolsHTTP() bool {
	h := provider.NewHTTPTestnet()

	pools, err := h.GetPools([]string{"Radium"})
	if err != nil {
		log.Errorf("error with GetPools request for Radium: %v", err)
		return true
	} else {
		log.Info(pools)
	}

	fmt.Println()
	return false
}

func callTickersHTTP() bool {
	h := provider.NewHTTPTestnet()

	tickers, err := h.GetTickers("SOLUSDT")
	if err != nil {
		log.Errorf("error with GetTickers request for SOLUSDT: %v", err)
		return true
	} else {
		log.Info(tickers)
	}

	fmt.Println()
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

func callPlaceOrderHTTP(ownerAddr, ooAddr string) (uint64, bool) {
	client := &http.Client{Timeout: time.Second * 30}
	rpcOpts := provider.DefaultRPCOpts(provider.TestnetHTTP)
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
		log.Errorf("failed to create order (%v)", err)
		return 0, true
	}
	log.Infof("created unsigned place order transaction: %v", response.Transaction)

	// sign/submit transaction after creation
	sig, err := h.SubmitOrder(ownerAddr, ownerAddr, marketAddr,
		orderSide, []pb.OrderType{orderType}, orderAmount,
		orderPrice, opts)
	if err != nil {
		log.Errorf("failed to submit order (%v)", err)
		return 0, true
	}

	log.Infof("placed order %v with clientOrderID %v", sig, clientOrderID)

	return clientOrderID, false
}

func callCancelByClientOrderIDHTTP(ownerAddr, ooAddr string, clientOrderID uint64) bool {
	time.Sleep(60 * time.Second)
	client := &http.Client{Timeout: time.Second * 30}
	opts := provider.DefaultRPCOpts(provider.TestnetHTTP)
	h := provider.NewHTTPClientWithOpts(client, opts)

	_, err := h.SubmitCancelByClientOrderID(clientOrderID, ownerAddr,
		marketAddr, ooAddr, true)
	if err != nil {
		log.Errorf("failed to cancel order by client ID (%v)", err)
		return true
	}

	log.Infof("canceled order for clientOrderID %v", clientOrderID)
	return false
}

func callPostSettleHTTP(ownerAddr, ooAddr string) bool {
	time.Sleep(60 * time.Second)
	client := &http.Client{Timeout: time.Second * 30}
	opts := provider.DefaultRPCOpts(provider.TestnetHTTP)
	h := provider.NewHTTPClientWithOpts(client, opts)

	sig, err := h.SubmitSettle(ownerAddr, "SOL/USDC", "F75gCEckFAyeeCWA9FQMkmLCmke7ehvBnZeVZ3QgvJR7", "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv", ooAddr, false)
	if err != nil {
		log.Errorf("error with post transaction stream request for SOL/USDC: %v", err)
		return true
	}

	log.Infof("response signature received: %v", sig)
	return false
}

func cancelAll(ownerAddr, payerAddr, ooAddr string) bool {
	log.Info("starting cancel all test")
	fmt.Println()

	client := &http.Client{Timeout: time.Second * 30}
	rpcOpts := provider.DefaultRPCOpts(provider.TestnetHTTP)
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
	log.Info("placing orders")
	sig, err := h.SubmitOrder(ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #1, signature %s", sig)

	opts.ClientOrderID = clientOrderID2
	sig, err = h.SubmitOrder(ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check orders are there
	orders, err := h.GetOpenOrders(marketAddr, ownerAddr, "")
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
	sigs, err := h.SubmitCancelAll(marketAddr, ownerAddr, []string{ooAddr}, true)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("placing cancel order(s) %s", strings.Join(sigs, ", "))

	time.Sleep(time.Minute)

	orders, err = h.GetOpenOrders(marketAddr, ownerAddr, "")
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
	callPostSettleHTTP(ownerAddr, ooAddr)
	return false
}

func callReplaceByClientOrderID(ownerAddr, payerAddr, ooAddr string) bool {
	log.Info("starting replace by client order ID test")
	fmt.Println()

	client := &http.Client{Timeout: time.Second * 60}
	rpcOpts := provider.DefaultRPCOpts(provider.TestnetHTTP)
	h := provider.NewHTTPClientWithOpts(client, rpcOpts)

	rand.Seed(time.Now().UnixNano())
	clientOrderID1 := rand.Uint64()
	opts := provider.PostOrderOpts{
		ClientOrderID:     clientOrderID1,
		OpenOrdersAddress: ooAddr,
		SkipPreFlight:     true,
	}

	// Place order in orderbook
	log.Info("placing order")
	sig, err := h.SubmitOrder(ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	} else {
		log.Infof("submitting place order #1, signature %s", sig)
	}
	time.Sleep(time.Minute)
	// Check order is there
	orders, err := h.GetOpenOrders(marketAddr, ownerAddr, "")
	if err != nil {
		log.Error(err)
		return true
	}
	found1 := false

	for _, order := range orders.Orders {
		if order.ClientOrderID == fmt.Sprintf("%v", clientOrderID1) {
			found1 = true
			continue
		}
	}
	if !(found1) {
		log.Error("order not found in orderbook")
		return true
	}
	log.Info("order placed successfully")

	// replacing order
	sig, err = h.SubmitReplaceByClientOrderID(ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice/2, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check order #2 is in orderbook
	orders, err = h.GetOpenOrders(marketAddr, ownerAddr, "")
	if err != nil {
		log.Error(err)
		return true
	}
	found2 := false

	for _, order := range orders.Orders {
		if order.ClientOrderID == fmt.Sprintf("%v", clientOrderID1) && order.Price == orderPrice/2 {
			found2 = true
		}
	}
	if !(found2) {
		log.Error("order #2 not found in orderbook")
		return true
	}
	log.Info("order #2 placed successfully")

	// Cancel all the orders
	log.Info("cancelling the orders")
	sigs, err := h.SubmitCancelAll(marketAddr, ownerAddr, []string{ooAddr}, true)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("placing cancel order(s) %s", strings.Join(sigs, ", "))
	return false
}

func callReplaceOrder(ownerAddr, payerAddr, ooAddr string) bool {
	log.Info("starting replace order test")
	fmt.Println()

	client := &http.Client{Timeout: time.Second * 30}
	rpcOpts := provider.DefaultRPCOpts(provider.TestnetHTTP)
	h := provider.NewHTTPClientWithOpts(client, rpcOpts)

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
	sig, err := h.SubmitOrder(ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	} else {
		log.Infof("submitting place order #1, signature %s", sig)
	}
	time.Sleep(time.Minute)
	// Check orders are there
	orders, err := h.GetOpenOrders(marketAddr, ownerAddr, "")
	if err != nil {
		log.Error(err)
		return true
	}
	var found1 *pb.Order

	for _, order := range orders.Orders {
		if order.ClientOrderID == fmt.Sprintf("%v", clientOrderID1) {
			found1 = order
			continue
		}
	}
	if found1 == nil {
		log.Error("order not found in orderbook")
		return true
	} else {
		log.Info("order placed successfully")
	}

	opts.ClientOrderID = clientOrderID2
	sig, err = h.SubmitReplaceOrder(found1.OrderID, ownerAddr, payerAddr, marketAddr, orderSide, []pb.OrderType{orderType}, orderAmount, orderPrice/2, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check orders are there
	orders, err = h.GetOpenOrders(marketAddr, ownerAddr, "")
	if err != nil {
		log.Error(err)
		return true
	}
	var found2 *pb.Order

	for _, order := range orders.Orders {
		if order.ClientOrderID == fmt.Sprintf("%v", clientOrderID2) {
			found2 = order
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
	sigs, err := h.SubmitCancelAll(marketAddr, ownerAddr, []string{ooAddr}, true)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("placing cancel order(s) %s", strings.Join(sigs, ", "))
	return false
}

func callGetRecentBlockHash() bool {
	h := provider.NewHTTPTestnet()

	hash, err := h.GetRecentBlockHash()
	if err != nil {
		log.Errorf("error with GetRecentBlockHash request: %v", err)
		return true
	} else {
		log.Info(hash)
	}

	fmt.Println()
	return false
}

func callTradeSwap(ownerAddr string) bool {
	log.Info("starting trade swap test")

	client := &http.Client{Timeout: time.Second * 30}
	rpcOpts := provider.DefaultRPCOpts(provider.TestnetHTTP)
	h := provider.NewHTTPClientWithOpts(client, rpcOpts)

	log.Info("trade swap")
	sig, err := h.SubmitTradeSwap(ownerAddr, "USDC", "SOL",
		0.01, 0.1, "raydium", false)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("trade swap transaction signature : %s", sig)
	return false
}
