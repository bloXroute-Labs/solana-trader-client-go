package main

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/examples/config"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	"github.com/bloXroute-Labs/solana-trader-proto/common"
	"math/rand"
	"net/http"
	"os"
	"time"

	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	log "github.com/sirupsen/logrus"
)

func httpClient() *provider.HTTPClient {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	var h *provider.HTTPClient

	switch cfg.Env {
	case config.EnvLocal:
		h = provider.NewHTTPLocal()
	case config.EnvTestnet:
		h = provider.NewHTTPTestnet()
	case config.EnvMainnet:
		h = provider.NewHTTPClient()
	}
	return h
}

func httpClientWithTimeout(timeout time.Duration) *provider.HTTPClient {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	var h *provider.HTTPClient
	client := &http.Client{Timeout: timeout}

	switch cfg.Env {
	case config.EnvLocal:
		h = provider.NewHTTPLocal()
	case config.EnvTestnet:
		h = provider.NewHTTPClientWithOpts(client, provider.DefaultRPCOpts(provider.TestnetHTTP))
	case config.EnvMainnet:
		h = provider.NewHTTPClientWithOpts(client, provider.DefaultRPCOpts(provider.MainnetHTTP))
	}
	return h
}

func main() {
	utils.InitLogger()
	failed := run()
	if failed {
		log.Fatal("one or multiple examples failed")
	}
}

func run() bool {
	var failed bool

	// informational methods
	failed = failed || logCall("callMarketsHTTP", func() bool { return callMarketsHTTP() })
	failed = failed || logCall("callOrderbookHTTP", func() bool { return callOrderbookHTTP() })
	failed = failed || logCall("callMarketDepthHTTP", func() bool { return callMarketDepthHTTP() })
	failed = failed || logCall("callOpenOrdersHTTP", func() bool { return callOpenOrdersHTTP() })
	failed = failed || logCall("callTradesHTTP", func() bool { return callTradesHTTP() })
	failed = failed || logCall("callPoolsHTTP", func() bool { return callPoolsHTTP() })
	failed = failed || logCall("callPriceHTTP", func() bool { return callPriceHTTP() })
	failed = failed || logCall("callTickersHTTP", func() bool { return callTickersHTTP() })
	failed = failed || logCall("callUnsettledHTTP", func() bool { return callUnsettledHTTP() })
	failed = failed || logCall("callGetAccountBalanceHTTP", func() bool { return callGetAccountBalanceHTTP() })
	failed = failed || logCall("callGetQuotes", func() bool { return callGetQuotes() })
	failed = failed || logCall("callDriftOrderbookHTTP", func() bool { return callDriftOrderbookHTTP() })

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	if !cfg.RunTrades {
		log.Info("skipping trades due to config")
		return failed
	}

	// calls below this place an order and immediately cancel it
	// you must specify:
	//	- PRIVATE_KEY (by default loaded during provider.NewGRPCClient()) to sign transactions
	// 	- PUBLIC_KEY to indicate which account you wish to trade from
	//	- OPEN_ORDERS to indicate your Serum account to speed up lookups (optional in actual usage)
	ownerAddr, ok := os.LookupEnv("PUBLIC_KEY")
	if !ok {
		log.Infof("PUBLIC_KEY environment variable not set: will skip place/cancel/settle examples")
		return failed
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
	failed = failed || logCall("callPlaceOrderHTTP", func() bool { return fail })
	failed = failed || logCall("callCancelByClientOrderIDHTTP", func() bool { return callCancelByClientOrderIDHTTP(ownerAddr, ooAddr, clientOrderID) })
	failed = failed || logCall("callPostSettleHTTP", func() bool { return callPostSettleHTTP(ownerAddr, ooAddr) })
	failed = failed || logCall("cancelAll", func() bool { return cancelAll(ownerAddr, payerAddr, ooAddr) })
	failed = failed || logCall("callReplaceByClientOrderID", func() bool { return callReplaceByClientOrderID(ownerAddr, payerAddr, ooAddr) })
	failed = failed || logCall("callReplaceOrder", func() bool { return callReplaceOrder(ownerAddr, payerAddr, ooAddr) })
	failed = failed || logCall("callGetRecentBlockHash", func() bool { return callGetRecentBlockHash() })
	failed = failed || logCall("callTradeSwap", func() bool { return callTradeSwap(ownerAddr) })
	failed = failed || logCall("callRouteTradeSwap", func() bool { return callRouteTradeSwap(ownerAddr) })

	return failed
}

func logCall(name string, call func() bool) bool {
	log.Infof("Executing `%s'...", name)

	result := call()
	if result {
		log.Errorf("`%s' failed", name)
	}

	return result
}

func callMarketsHTTP() bool {
	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	markets, err := h.GetMarkets(ctx)
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
	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orderbook, err := h.GetOrderbook(ctx, "SOL-USDT", 0, pb.Project_P_OPENBOOK)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL-USDT: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()

	orderbook, err = h.GetOrderbook(ctx, "SOLUSDT", 2, pb.Project_P_OPENBOOK)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOLUSDT: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()

	orderbook, err = h.GetOrderbook(ctx, "SOL:USDC", 3, pb.Project_P_OPENBOOK)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL:USDC: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	return false
}

func callMarketDepthHTTP() bool {
	h := httpClient()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mdData, err := h.GetMarketDepth(ctx, "SOL-USDC", 0, pb.Project_P_OPENBOOK)
	if err != nil {
		log.Errorf("error with GetMarketDepth request for SOL-USDC: %v", err)
		return true
	} else {
		log.Info(mdData)
	}

	return false
}

func callOpenOrdersHTTP() bool {
	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	orders, err := h.GetOpenOrders(ctx, "SOLUSDT", "HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc", "", pb.Project_P_OPENBOOK)
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
	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.GetUnsettled(ctx, "SOLUSDT", "HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc", pb.Project_P_OPENBOOK)
	if err != nil {
		log.Errorf("error with GetUnsettled request for SOLUSDT: %v", err)
		return true
	} else {
		log.Info(response)
	}

	fmt.Println()
	return false
}

func callGetAccountBalanceHTTP() bool {
	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.GetAccountBalance(ctx, "F75gCEckFAyeeCWA9FQMkmLCmke7ehvBnZeVZ3QgvJR7")
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
	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	trades, err := h.GetTrades(ctx, "SOLUSDT", 5, pb.Project_P_OPENBOOK)
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
	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pools, err := h.GetPools(ctx, []pb.Project{pb.Project_P_RAYDIUM})
	if err != nil {
		log.Errorf("error with GetPools request for Raydium: %v", err)
		return true
	} else {
		log.Info(pools)
	}

	fmt.Println()
	return false
}

func callPriceHTTP() bool {
	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prices, err := h.GetPrice(ctx, []string{"SOL", "ETH"})
	if err != nil {
		log.Errorf("error with GetPrice request for SOL and ETH: %v", err)
		return true
	} else {
		log.Info(prices)
	}

	fmt.Println()
	return false
}

func callTickersHTTP() bool {
	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tickers, err := h.GetTickers(ctx, "SOLUSDT", pb.Project_P_OPENBOOK)
	if err != nil {
		log.Errorf("error with GetTickers request for SOLUSDT: %v", err)
		return true
	} else {
		log.Info(tickers)
	}

	fmt.Println()
	return false
}

func callGetQuotes() bool {
	h := httpClientWithTimeout(time.Second * 60)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	inToken := "SOL"
	outToken := "USDC"
	amount := 0.01
	slippage := float64(5)
	limit := 5

	quotes, err := h.GetQuotes(ctx, inToken, outToken, amount, slippage, int32(limit), []pb.Project{pb.Project_P_ALL})
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

const (
	// SOL/USDC market
	marketAddr = "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT"

	orderSide   = pb.Side_S_ASK
	orderType   = common.OrderType_OT_LIMIT
	orderPrice  = float64(170200)
	orderAmount = float64(0.1)
)

func callPlaceOrderHTTP(ownerAddr, ooAddr string) (uint64, bool) {
	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// generate a random clientOrderId for this order
	rand.Seed(time.Now().UnixNano())
	clientOrderID := rand.Uint64()

	opts := provider.PostOrderOpts{
		ClientOrderID:     clientOrderID,
		OpenOrdersAddress: ooAddr,
	}

	// create order without actually submitting
	response, err := h.PostOrder(ctx, ownerAddr, ownerAddr, marketAddr, orderSide, []common.OrderType{orderType}, orderAmount, orderPrice, pb.Project_P_OPENBOOK, opts)
	if err != nil {
		log.Errorf("failed to create order (%v)", err)
		return 0, true
	}
	log.Infof("created unsigned place order transaction: %v", response.Transaction)

	// sign/submit transaction after creation
	sig, err := h.SubmitOrder(ctx, ownerAddr, ownerAddr, marketAddr,
		orderSide, []common.OrderType{orderType}, orderAmount,
		orderPrice, pb.Project_P_OPENBOOK, opts)
	if err != nil {
		log.Errorf("failed to submit order (%v)", err)
		return 0, true
	}

	log.Infof("placed order %v with clientOrderID %v", sig, clientOrderID)

	return clientOrderID, false
}

func callCancelByClientOrderIDHTTP(ownerAddr, ooAddr string, clientOrderID uint64) bool {
	time.Sleep(60 * time.Second)
	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := h.SubmitCancelByClientOrderID(ctx, clientOrderID, ownerAddr,
		marketAddr, ooAddr, pb.Project_P_OPENBOOK, true)
	if err != nil {
		log.Errorf("failed to cancel order by client ID (%v)", err)
		return true
	}

	log.Infof("canceled order for clientOrderID %v", clientOrderID)
	return false
}

func callPostSettleHTTP(ownerAddr, ooAddr string) bool {
	time.Sleep(60 * time.Second)
	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := h.SubmitSettle(ctx, ownerAddr, "SOL/USDC", "F75gCEckFAyeeCWA9FQMkmLCmke7ehvBnZeVZ3QgvJR7", "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv", ooAddr, pb.Project_P_OPENBOOK, false)
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

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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
	sig, err := h.SubmitOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []common.OrderType{orderType}, orderAmount, orderPrice, pb.Project_P_OPENBOOK, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #1, signature %s", sig)

	opts.ClientOrderID = clientOrderID2
	sig, err = h.SubmitOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []common.OrderType{orderType}, orderAmount, orderPrice, pb.Project_P_OPENBOOK, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check orders are there
	orders, err := h.GetOpenOrders(ctx, marketAddr, ownerAddr, "", pb.Project_P_OPENBOOK)
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
	sigs, err := h.SubmitCancelAll(ctx, marketAddr, ownerAddr, []string{ooAddr}, pb.Project_P_OPENBOOK, provider.SubmitOpts{
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

	time.Sleep(time.Minute)

	orders, err = h.GetOpenOrders(ctx, marketAddr, ownerAddr, "", pb.Project_P_OPENBOOK)
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

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
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
	sig, err := h.SubmitOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []common.OrderType{orderType}, orderAmount, orderPrice, pb.Project_P_OPENBOOK, opts)
	if err != nil {
		log.Error(err)
		return true
	} else {
		log.Infof("submitting place order #1, signature %s", sig)
	}
	time.Sleep(time.Minute)
	// Check order is there
	orders, err := h.GetOpenOrders(ctx, marketAddr, ownerAddr, "", pb.Project_P_OPENBOOK)
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
	sig, err = h.SubmitReplaceByClientOrderID(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []common.OrderType{orderType}, orderAmount, orderPrice/2, pb.Project_P_OPENBOOK, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check order #2 is in orderbook
	orders, err = h.GetOpenOrders(ctx, marketAddr, ownerAddr, "", pb.Project_P_OPENBOOK)
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
	sigs, err := h.SubmitCancelAll(ctx, marketAddr, ownerAddr, []string{ooAddr}, pb.Project_P_OPENBOOK, provider.SubmitOpts{
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

func callReplaceOrder(ownerAddr, payerAddr, ooAddr string) bool {
	log.Info("starting replace order test")
	fmt.Println()

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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
	sig, err := h.SubmitOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []common.OrderType{orderType}, orderAmount, orderPrice, pb.Project_P_OPENBOOK, opts)
	if err != nil {
		log.Error(err)
		return true
	} else {
		log.Infof("submitting place order #1, signature %s", sig)
	}
	time.Sleep(time.Minute)
	// Check orders are there
	orders, err := h.GetOpenOrders(ctx, marketAddr, ownerAddr, "", pb.Project_P_OPENBOOK)
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
	sig, err = h.SubmitReplaceOrder(ctx, found1.OrderID, ownerAddr, payerAddr, marketAddr, orderSide, []common.OrderType{orderType}, orderAmount, orderPrice/2, pb.Project_P_OPENBOOK, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check orders are there
	orders, err = h.GetOpenOrders(ctx, marketAddr, ownerAddr, "", pb.Project_P_OPENBOOK)
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
	sigs, err := h.SubmitCancelAll(ctx, marketAddr, ownerAddr, []string{ooAddr}, pb.Project_P_OPENBOOK, provider.SubmitOpts{
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

func callGetRecentBlockHash() bool {
	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	hash, err := h.GetRecentBlockHash(ctx)
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

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("trade swap")
	sig, err := h.SubmitTradeSwap(ctx, ownerAddr, "USDC", "SOL",
		0.01, 0.1, "raydium", provider.SubmitOpts{
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

func callRouteTradeSwap(ownerAddr string) bool {
	log.Info("starting route trade swap test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("route trade swap")
	sig, err := h.SubmitRouteTradeSwap(ctx, &pb.RouteTradeSwapRequest{
		OwnerAddress: ownerAddr,
		Project:      pb.Project_P_RAYDIUM,
		Steps: []*pb.RouteStep{
			{
				// FIDA-RAY pool address

				InToken:      "FIDA",
				OutToken:     "RAY",
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

func callDriftOrderbookHTTP() bool {
	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orderbook, err := h.GetPerpOrderbook(ctx, "SOL-PERP", 0, pb.Project_P_DRIFT)
	if err != nil {
		log.Errorf("error with GetPerpOrderbook request for SOL-PERP: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()
	return false
}
