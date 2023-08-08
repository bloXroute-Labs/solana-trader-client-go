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

const (
	sideBid      = "bid"
	sideAsk      = "ask"
	typeLimit    = "limit"
	typeIOC      = "ioc"
	typePostOnly = "postonly"
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
		h = provider.NewHTTPClientWithOpts(client, provider.DefaultRPCOpts(provider.MainnetVirginiaHTTP))
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
	failed = failed || logCall("callRaydiumPools ", func() bool { return callRaydiumPools() })
	failed = failed || logCall("callRaydiumPrices", func() bool { return callRaydiumPrices() })
	failed = failed || logCall("callJupiterPrices", func() bool { return callJupiterPrices() })
	failed = failed || logCall("callPriceHTTP", func() bool { return callPriceHTTP() })
	failed = failed || logCall("callTickersHTTP", func() bool { return callTickersHTTP() })
	failed = failed || logCall("callUnsettledHTTP", func() bool { return callUnsettledHTTP() })
	failed = failed || logCall("callGetAccountBalanceHTTP", func() bool { return callGetAccountBalanceHTTP() })
	failed = failed || logCall("callGetQuotesHTTP", func() bool { return callGetQuotesHTTP() })
	failed = failed || logCall("callGetRaydiumQuotes", func() bool { return callGetRaydiumQuotes() })
	failed = failed || logCall("callGetJupiterQuotes", func() bool { return callGetJupiterQuotes() })
	failed = failed || logCall("callDriftPerpOrderbookHTTP", func() bool { return callDriftPerpOrderbookHTTP() })
	failed = failed || logCall("callDriftGetMarginOrderbookHTTP", func() bool { return callDriftGetMarginOrderbookHTTP() })
	failed = failed || logCall("callDriftMarketDepthHTTP", func() bool { return callDriftMarketDepthHTTP() })

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// calls below this place an order and immediately cancel it
	// you must specify:
	//	- PRIVATE_KEY (by default loaded during provider.NewClient()) to sign transactions
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

	if cfg.RunTrades {
		// Order lifecycle
		clientOrderID, fail := callPlaceOrderHTTP(ownerAddr, ooAddr, sideAsk, typeLimit)
		failed = failed || logCall("callPlaceOrderHTTP", func() bool { return fail })
		failed = failed || logCall("callCancelByClientOrderIDHTTP", func() bool { return callCancelByClientOrderIDHTTP(ownerAddr, ooAddr, clientOrderID) })
		failed = failed || logCall("callPostSettleHTTP", func() bool { return callPostSettleHTTP(ownerAddr, ooAddr) })
		failed = failed || logCall("cancelAll", func() bool { return cancelAll(ownerAddr, payerAddr, ooAddr, sideAsk, typeLimit) })
		failed = failed || logCall("callReplaceByClientOrderID", func() bool { return callReplaceByClientOrderID(ownerAddr, payerAddr, ooAddr, sideAsk, typeLimit) })
		failed = failed || logCall("callReplaceOrder", func() bool { return callReplaceOrder(ownerAddr, payerAddr, ooAddr, sideAsk, typeLimit) })
		failed = failed || logCall("callGetRecentBlockHash", func() bool { return callGetRecentBlockHash() })
		failed = failed || logCall("callTradeSwap", func() bool { return callTradeSwap(ownerAddr) })
		failed = failed || logCall("callRouteTradeSwap", func() bool { return callRouteTradeSwap(ownerAddr) })
		failed = failed || logCall("callRaydiumTradeSwap", func() bool { return callRaydiumSwap(ownerAddr) })
		failed = failed || logCall("callJupiterTradeSwap", func() bool { return callJupiterSwap(ownerAddr) })
		failed = failed || logCall("callRaydiumRouteTradeSwap", func() bool { return callRaydiumRouteSwap(ownerAddr) })
		failed = failed || logCall("callJupiterRouteTradeSwap", func() bool { return callJupiterRouteSwap(ownerAddr) })
	}

	failed = failed || logCall("callGetOpenPerpOrders", func() bool { return callGetOpenPerpOrders(ownerAddr) })
	failed = failed || logCall("callGetDriftOpenMarginOrders", func() bool { return callGetDriftOpenMarginOrders(ownerAddr) })
	failed = failed || logCall("callGetPerpPositions", func() bool { return callGetPerpPositions(ownerAddr) })
	failed = failed || logCall("callGetDriftPerpPositions", func() bool { return callGetDriftPerpPositions(ownerAddr) })
	failed = failed || logCall("callGetUser", func() bool { return callGetUser(ownerAddr) })

	failed = failed || logCall("callGetOpenPerpOrder", func() bool { return callGetOpenPerpOrder(ownerAddr) })
	failed = failed || logCall("callGetDriftOpenPerpOrders", func() bool { return callGetDriftOpenPerpOrders(ownerAddr) })
	failed = failed || logCall("callGetAssets", func() bool { return callGetAssets(ownerAddr) })
	failed = failed || logCall("callGetPerpContracts", func() bool { return callGetPerpContracts() })
	failed = failed || logCall("callGetDriftMarkets", func() bool { return callGetDriftMarkets() })

	failed = failed || logCall("callGetDriftAssets", func() bool { return callGetDriftAssets(ownerAddr) })
	failed = failed || logCall("callGetDriftPerpContracts", func() bool { return callGetDriftPerpContracts() })
	failed = failed || logCall("callGetDriftPerpOrderbook", func() bool { return callGetDriftPerpOrderbook() })
	failed = failed || logCall("callGetDriftUser", func() bool { return callGetDriftUser(ownerAddr) })
	failed = failed || logCall("callGetDriftOpenPerpOrder", func() bool { return callGetDriftOpenPerpOrder(ownerAddr) })
	failed = failed || logCall("callGetDriftOpenMarginOrder", func() bool { return callGetDriftOpenMarginOrder(ownerAddr) })

	if cfg.RunPerpTrades {
		failed = failed || logCall("callCancelPerpOrder", func() bool { return callCancelPerpOrder(ownerAddr) })
		failed = failed || logCall("callDriftCancelPerpOrder", func() bool { return callDriftCancelPerpOrder(ownerAddr) })
		failed = failed || logCall("callCancelDriftMarginOrder", func() bool { return callCancelDriftMarginOrder(ownerAddr) })
		failed = failed || logCall("callClosePerpPositions", func() bool { return callClosePerpPositions(ownerAddr) })
		failed = failed || logCall("callCreateUser", func() bool { return callCreateUser(ownerAddr) })
		failed = failed || logCall("callManageCollateralDeposit", func() bool { return callManageCollateralDeposit() })
		failed = failed || logCall("callPostPerpOrder", func() bool { return callPostPerpOrder(ownerAddr) })
		failed = failed || logCall("callPostDriftPerpOrder", func() bool { return callPostDriftPerpOrder(ownerAddr) })
		failed = failed || logCall("callPostModifyOrder", func() bool { return callPostModifyOrder(ownerAddr) })
		failed = failed || logCall("callPostMarginOrder", func() bool { return callPostMarginOrder(ownerAddr) })
		failed = failed || logCall("callManageCollateralWithdraw", func() bool { return callManageCollateralWithdraw() })
		failed = failed || logCall("callManageCollateralTransfer", func() bool { return callManageCollateralTransfer() })
		failed = failed || logCall("callDriftEnableMarginTrading", func() bool { return callDriftEnableMarginTrading(ownerAddr) })
		failed = failed || logCall("callPostSettlePNL", func() bool { return callPostSettlePNL(ownerAddr) })
		failed = failed || logCall("callPostSettlePNLs", func() bool { return callPostSettlePNLs(ownerAddr) })
		failed = failed || logCall("callPostLiquidatePerp", func() bool { return callPostLiquidatePerp(ownerAddr) })

		failed = failed || logCall("callPostCloseDriftPerpPositions", func() bool { return callPostCloseDriftPerpPositions(ownerAddr) })
		failed = failed || logCall("callPostCreateDriftUser", func() bool { return callPostCreateDriftUser(ownerAddr) })
		failed = failed || logCall("callPostDriftManageCollateralDeposit", func() bool { return callPostDriftManageCollateralDeposit() })
		failed = failed || logCall("callPostDriftManageCollateralWithdraw", func() bool { return callPostDriftManageCollateralWithdraw() })
		failed = failed || logCall("callPostDriftManageCollateralTransfer", func() bool { return callPostDriftManageCollateralTransfer() })
		failed = failed || logCall("callPostDriftSettlePNL", func() bool { return callPostDriftSettlePNL(ownerAddr) })
		failed = failed || logCall("callPostDriftSettlePNLs", func() bool { return callPostDriftSettlePNLs(ownerAddr) })
		failed = failed || logCall("callPostLiquidateDriftPerp", func() bool { return callPostLiquidateDriftPerp(ownerAddr) })
	}
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

	markets, err := h.GetMarketsV2(ctx)
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
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	orderbook, err := h.GetOrderbookV2(ctx, "SOL-USDT", 0)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL-USDT: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()

	orderbook, err = h.GetOrderbookV2(ctx, "SOLUSDT", 2)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOLUSDT: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()

	orderbook, err = h.GetOrderbookV2(ctx, "SOL:USDC", 3)
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

	mdData, err := h.GetMarketDepthV2(ctx, "SOL-USDC", 0)
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

	orders, err := h.GetOpenOrdersV2(ctx, "SOLUSDT", "HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc", "", "", 0)
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

	response, err := h.GetUnsettledV2(ctx, "SOLUSDT", "HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc")
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

func callRaydiumPools() bool {
	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pools, err := h.GetRaydiumPools(ctx, &pb.GetRaydiumPoolsRequest{})
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

func callRaydiumPrices() bool {
	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prices, err := h.GetRaydiumPrices(ctx, &pb.GetRaydiumPricesRequest{
		Tokens: []string{"SOL", "ETH"},
	})
	if err != nil {
		log.Errorf("error with GetRaydiumPrices request for SOL and ETH: %v", err)
		return true
	} else {
		log.Info(prices)
	}

	fmt.Println()
	return false
}

func callJupiterPrices() bool {
	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prices, err := h.GetJupiterPrices(ctx, &pb.GetJupiterPricesRequest{
		Tokens: []string{"SOL", "ETH"},
	})
	if err != nil {
		log.Errorf("error with GetJupiterPrices request for SOL and ETH: %v", err)
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

	tickers, err := h.GetTickersV2(ctx, "SOLUSDT")
	if err != nil {
		log.Errorf("error with GetTickers request for SOLUSDT: %v", err)
		return true
	} else {
		log.Info(tickers)
	}

	fmt.Println()
	return false
}

func callGetQuotesHTTP() bool {
	h := httpClientWithTimeout(time.Second * 60)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	inToken := "SOL"
	outToken := "USDT"
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

func callGetRaydiumQuotes() bool {
	h := httpClientWithTimeout(time.Second * 60)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	inToken := "SOL"
	outToken := "USDT"
	amount := 0.01
	slippage := float64(5)

	quotes, err := h.GetRaydiumQuotes(ctx, &pb.GetRaydiumQuotesRequest{
		InToken:  inToken,
		OutToken: outToken,
		InAmount: amount,
		Slippage: slippage,
	})
	if err != nil {
		log.Errorf("error with GetQuotes request for %s to %s: %v", inToken, outToken, err)
		return true
	}

	if err != nil {
		log.Errorf("error with GetRaydiumQuotes request for %s to %s: %v", inToken, outToken, err)
		return true
	}

	if len(quotes.Routes) != 1 {
		log.Errorf("did not get back 1 quotes, got %v quotes", len(quotes.Routes))
		return true
	}
	for _, route := range quotes.Routes {
		log.Infof("best route for Raydium is %v", route)
	}

	fmt.Println()
	return false
}

func callGetJupiterQuotes() bool {
	h := httpClientWithTimeout(time.Second * 60)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	inToken := "SOL"
	outToken := "USDT"
	amount := 0.01
	slippage := float64(5)
	limit := int32(3)

	quotes, err := h.GetJupiterQuotes(ctx, &pb.GetJupiterQuotesRequest{
		InToken:  inToken,
		OutToken: outToken,
		InAmount: amount,
		Slippage: slippage,
		Limit:    limit,
	})
	if err != nil {
		log.Errorf("error with GetQuotes request for %s to %s: %v", inToken, outToken, err)
		return true
	}

	if err != nil {
		log.Errorf("error with GetJupiterQuotes request for %s to %s: %v", inToken, outToken, err)
		return true
	}

	if len(quotes.Routes) != 3 {
		log.Errorf("did not get back 3 quotes, got %v quotes", len(quotes.Routes))
		return true
	}
	for _, route := range quotes.Routes {
		log.Infof("best route for Jupiter is %v", route)
	}

	fmt.Println()
	return false
}

const (
	// SOL/USDC market
	marketAddr = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh6"

	orderSide   = pb.Side_S_ASK
	orderPrice  = float64(170200)
	orderAmount = float64(0.1)
)

func callPlaceOrderHTTP(ownerAddr, ooAddr string, orderSide string, orderType string) (uint64, bool) {
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
	response, err := h.PostOrderV2(ctx, ownerAddr, ownerAddr, marketAddr, orderSide, orderType, orderAmount, orderPrice, opts)
	if err != nil {
		log.Errorf("failed to create order (%v)", err)
		return 0, true
	}
	log.Infof("created unsigned place order transaction: %v", response.Transaction)

	// sign/submit transaction after creation
	sig, err := h.SubmitOrderV2(ctx, ownerAddr, ownerAddr, marketAddr,
		orderSide, orderType, orderAmount,
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
	h := httpClient()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := h.SubmitCancelOrderV2(ctx, "", clientOrderID, sideAsk, ownerAddr,
		marketAddr, ooAddr, provider.SubmitOpts{
			SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
			SkipPreFlight:  false,
		})
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

	sig, err := h.SubmitSettleV2(ctx, ownerAddr, "SOL/USDC", "F75gCEckFAyeeCWA9FQMkmLCmke7ehvBnZeVZ3QgvJR7", "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv", ooAddr, false)
	if err != nil {
		log.Errorf("error with post transaction stream request for SOL/USDC: %v", err)
		return true
	}

	log.Infof("response signature received: %v", sig)
	return false
}

func cancelAll(ownerAddr, payerAddr, ooAddr string, orderSide string, orderType string) bool {
	log.Info("starting cancel all test")
	fmt.Println()

	h := httpClient()
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
	sig, err := h.SubmitOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderType, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #1, signature %s", sig)

	opts.ClientOrderID = clientOrderID2
	sig, err = h.SubmitOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderType, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check orders are there
	orders, err := h.GetOpenOrdersV2(ctx, marketAddr, ownerAddr, "", "", 0)
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
	sigs, err := h.SubmitCancelOrderV2(ctx, "", 0, sideAsk, ownerAddr, marketAddr, ooAddr, provider.SubmitOpts{
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

	orders, err = h.GetOpenOrdersV2(ctx, marketAddr, ownerAddr, "", "", 0)
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

func callReplaceByClientOrderID(ownerAddr, payerAddr, ooAddr string, orderSide string, orderType string) bool {
	log.Info("starting replace by client order ID test")
	fmt.Println()

	h := httpClient()
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
	sig, err := h.SubmitOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderType, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	} else {
		log.Infof("submitting place order #1, signature %s", sig)
	}
	time.Sleep(time.Minute)
	// Check order is there
	orders, err := h.GetOpenOrdersV2(ctx, marketAddr, ownerAddr, "", "", 0)
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
	sig, err = h.SubmitReplaceOrderV2(ctx, "", ownerAddr, payerAddr, marketAddr, orderSide, orderType, orderAmount, orderPrice/2, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check order #2 is in orderbook
	orders, err = h.GetOpenOrdersV2(ctx, marketAddr, ownerAddr, "", "", 0)
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
	sigs, err := h.SubmitCancelOrderV2(ctx, "", 0, sideAsk, ownerAddr, marketAddr, ooAddr, provider.SubmitOpts{
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

func callReplaceOrder(ownerAddr, payerAddr, ooAddr string, orderSide string, orderType string) bool {
	log.Info("starting replace order test")
	fmt.Println()

	h := httpClient()
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
	sig, err := h.SubmitOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderType, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	} else {
		log.Infof("submitting place order #1, signature %s", sig)
	}
	time.Sleep(time.Minute)
	// Check orders are there
	orders, err := h.GetOpenOrdersV2(ctx, marketAddr, ownerAddr, "", "", 0)
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
	sig, err = h.SubmitReplaceOrderV2(ctx, found1.OrderID, ownerAddr, payerAddr, marketAddr, orderSide, typeLimit, orderAmount, orderPrice/2, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check orders are there
	orders, err = h.GetOpenOrdersV2(ctx, marketAddr, ownerAddr, "", "", 0)
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
	sigs, err := h.SubmitCancelOrderV2(ctx, "", 0, sideAsk, ownerAddr, marketAddr, ooAddr, provider.SubmitOpts{
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
	sig, err := h.SubmitTradeSwap(ctx, ownerAddr, "USDT", "SOL",
		0.01, 0.1, pb.Project_P_RAYDIUM, provider.SubmitOpts{
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

func callRaydiumSwap(ownerAddr string) bool {
	log.Info("starting Raydium swap test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("Raydium swap")
	sig, err := h.SubmitRaydiumSwap(ctx, &pb.PostRaydiumSwapRequest{
		OwnerAddress: ownerAddr,
		InToken:      "USDT",
		OutToken:     "SOL",
		Slippage:     0.1,
		InAmount:     0.01,
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_ABORT_ON_FIRST_ERROR,
		SkipPreFlight:  false,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Raydium swap transaction signature : %s", sig)
	return false
}

func callRaydiumRouteSwap(ownerAddr string) bool {
	log.Info("starting Raydium route swap test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("Raydium route swap")
	sig, err := h.SubmitRaydiumRouteSwap(ctx, &pb.PostRaydiumRouteSwapRequest{
		OwnerAddress: ownerAddr,
		Slippage:     0.1,
		Steps: []*pb.RaydiumRouteStep{
			{
				InToken:      "FIDA",
				OutToken:     "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R",
				InAmount:     0.01,
				OutAmountMin: 0.007505,
				OutAmount:    0.0074,
			},
			{
				InToken:      "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R",
				OutToken:     "USDT",
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
	log.Infof("Raydium route swap transaction signature : %s", sig)
	return false
}

func callJupiterRouteSwap(ownerAddr string) bool {
	log.Info("starting Jupiter route swap test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("Jupiter route swap")
	sig, err := h.SubmitJupiterRouteSwap(ctx, &pb.PostJupiterRouteSwapRequest{
		OwnerAddress: ownerAddr,
		Slippage:     0.1,
		Steps: []*pb.JupiterRouteStep{
			{
				InToken:      "FIDA",
				OutToken:     "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R",
				InAmount:     0.01,
				OutAmountMin: 0.007505,
				OutAmount:    0.0074,
			},
			{
				InToken:      "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R",
				OutToken:     "USDT",
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
	log.Infof("Jupiter route swap transaction signature : %s", sig)
	return false
}

func callJupiterSwap(ownerAddr string) bool {
	log.Info("starting Jupiter swap test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("Jupiter swap")
	sig, err := h.SubmitJupiterSwap(ctx, &pb.PostJupiterSwapRequest{
		OwnerAddress: ownerAddr,
		InToken:      "USDT",
		OutToken:     "SOL",
		Slippage:     0.1,
		InAmount:     0.01,
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_ABORT_ON_FIRST_ERROR,
		SkipPreFlight:  false,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Jupiter swap transaction signature : %s", sig)
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
		Slippage:     0.1,
		Steps: []*pb.RouteStep{
			{
				Project: &pb.StepProject{
					Label: "Raydium",
					Id:    "",
				},
				InToken:      "FIDA",
				OutToken:     "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R",
				InAmount:     0.01,
				OutAmountMin: 0.007505,
				OutAmount:    0.0074,
			},
			{
				Project: &pb.StepProject{
					Label: "Raydium",
					Id:    "",
				},
				InToken:      "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R",
				OutToken:     "USDT",
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

func callDriftPerpOrderbookHTTP() bool {
	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orderbook, err := h.GetPerpOrderbook(ctx, &pb.GetPerpOrderbookRequest{
		Contract: common.PerpContract_SOL_PERP,
		Limit:    0,
		Project:  pb.Project_P_DRIFT,
	})
	if err != nil {
		log.Errorf("error with GetPerpOrderbook request for SOL-PERP: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()
	return false
}

func callDriftMarketDepthHTTP() bool {
	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	marketDepth, err := h.GetDriftMarketDepth(ctx, &pb.GetDriftMarketDepthRequest{
		Contract: "SOL_PERP",
		Limit:    0,
	})
	if err != nil {
		log.Errorf("error with GetDriftMarketDepth request for SOL_PERP: %v", err)
		return true
	} else {
		log.Info(marketDepth)
	}

	fmt.Println()
	return false
}

func callDriftGetMarginOrderbookHTTP() bool {
	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	orderbook, err := h.GetDriftMarginOrderbook(ctx, &pb.GetDriftMarginOrderbookRequest{
		Market:   "SOL",
		Limit:    0,
		Metadata: true,
	})
	if err != nil {
		log.Errorf("error with GetPerpOrderbook request for SOL-PERP: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()
	return false
}

func callGetOpenPerpOrders(ownerAddr string) bool {
	log.Info("starting callGetOpenPerpOrders test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := h.GetOpenPerpOrders(ctx, &pb.GetOpenPerpOrdersRequest{
		OwnerAddress:   ownerAddr,
		AccountAddress: "",
		Contracts:      []common.PerpContract{common.PerpContract_SOL_PERP},
		Project:        pb.Project_P_DRIFT,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("GetOpenPerpOrders resp : %s", user)
	return false
}

func callGetDriftOpenMarginOrders(ownerAddr string) bool {
	log.Info("starting callGetDriftOpenMarginOrders test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := h.GetDriftOpenMarginOrders(ctx, &pb.GetDriftOpenMarginOrdersRequest{
		OwnerAddress:   ownerAddr,
		AccountAddress: "",
		Markets:        []string{"SOL"},
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetDriftOpenMarginOrders resp : %s", user)
	return false
}

func callGetDriftOpenMarginOrder(ownerAddr string) bool {
	log.Info("starting callGetDriftOpenMarginOrder test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := h.GetDriftOpenMarginOrder(ctx, &pb.GetDriftOpenMarginOrderRequest{
		OwnerAddress:   ownerAddr,
		AccountAddress: "",
		ClientOrderID:  13,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetDriftOpenMarginOrder resp : %s", user)
	return false
}

func callGetPerpPositions(ownerAddr string) bool {
	log.Info("starting callGetPerpPositions test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := h.GetPerpPositions(ctx, &pb.GetPerpPositionsRequest{
		OwnerAddress:   ownerAddr,
		AccountAddress: "",
		Contracts:      []common.PerpContract{common.PerpContract_SOL_PERP},
		Project:        pb.Project_P_DRIFT,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("GetPerpPositions resp : %s", user)
	return false
}

func callGetDriftPerpPositions(ownerAddr string) bool {
	log.Info("starting callGetDriftPerpPositions test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := h.GetDriftPerpPositions(ctx, &pb.GetDriftPerpPositionsRequest{
		OwnerAddress:   ownerAddr,
		AccountAddress: "",
		Contracts:      []string{"SOL_PERP"},
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("GetDriftPerpPositions resp : %s", user)
	return false
}

func callGetUser(ownerAddr string) bool {
	log.Info("starting callGetUser test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := h.GetUser(ctx, &pb.GetUserRequest{
		OwnerAddress: ownerAddr,
		Project:      pb.Project_P_DRIFT,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("GetUser resp : %s", user)
	return false
}

func callGetDriftUser(ownerAddr string) bool {
	log.Info("starting callGetDriftUser test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := h.GetDriftUser(ctx, &pb.GetDriftUserRequest{
		OwnerAddress: ownerAddr,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetDriftUser resp : %s", user)
	return false
}

func callCancelPerpOrder(ownerAddr string) bool {
	log.Info("starting callCancelPerpOrder test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := h.PostCancelPerpOrder(ctx, &pb.PostCancelPerpOrderRequest{
		Project:       pb.Project_P_DRIFT,
		OwnerAddress:  ownerAddr,
		OrderID:       1,
		ClientOrderID: 0,
		Contract:      common.PerpContract_SOL_PERP,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callCancelPerpOrder signature : %s", sig)
	return false
}

func callDriftCancelPerpOrder(ownerAddr string) bool {
	log.Info("starting callDriftCancelPerpOrder test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := h.PostDriftCancelPerpOrder(ctx, &pb.PostDriftCancelPerpOrderRequest{
		OwnerAddress:  ownerAddr,
		OrderID:       1,
		ClientOrderID: 0,
		Contract:      "SOL_PERP",
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callDriftCancelPerpOrder signature : %s", sig)
	return false
}

func callCancelDriftMarginOrder(ownerAddr string) bool {
	log.Info("starting callCancelDriftMarginOrder test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := h.PostCancelDriftMarginOrder(ctx, &pb.PostCancelDriftMarginOrderRequest{
		OwnerAddress:  ownerAddr,
		OrderID:       1,
		ClientOrderID: 0,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callCancelDriftMarginOrder signature : %s", sig)
	return false
}

func callClosePerpPositions(ownerAddr string) bool {
	log.Info("starting callClosePerpPositions test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	sig, err := h.PostClosePerpPositions(ctx, &pb.PostClosePerpPositionsRequest{
		Project:      pb.Project_P_DRIFT,
		OwnerAddress: ownerAddr,
		Contracts:    []common.PerpContract{common.PerpContract_SOL_PERP},
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callClosePerpPositions signature : %s", sig)
	return false
}

func callPostCloseDriftPerpPositions(ownerAddr string) bool {
	log.Info("starting callPostCloseDriftPerpPositions test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	sig, err := h.PostCloseDriftPerpPositions(ctx, &pb.PostCloseDriftPerpPositionsRequest{
		OwnerAddress: ownerAddr,
		Contracts:    []string{"SOL_PERP"},
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostCloseDriftPerpPositions signature : %s", sig)
	return false
}

func callCreateUser(ownerAddr string) bool {
	log.Info("starting callCreateUser test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := h.PostCreateUser(ctx, &pb.PostCreateUserRequest{
		Project:      pb.Project_P_DRIFT,
		OwnerAddress: ownerAddr,
		Action:       "create",
		SubAccountID: 10,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callCreateUser signature : %s", sig)
	return false
}

func callPostCreateDriftUser(ownerAddr string) bool {
	log.Info("starting callPostCreateDriftUser test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := h.PostCreateDriftUser(ctx, &pb.PostCreateDriftUserRequest{
		OwnerAddress: ownerAddr,
		Action:       "create",
		SubAccountID: 10,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostCreateDriftUser signature : %s", sig)
	return false
}

func callPostPerpOrder(ownerAddr string) bool {
	log.Info("starting callPostPerpOrder test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	request := &pb.PostPerpOrderRequest{
		Project:        pb.Project_P_DRIFT,
		OwnerAddress:   ownerAddr,
		Contract:       common.PerpContract_SOL_PERP,
		AccountAddress: "",
		PositionSide:   common.PerpPositionSide_PS_SHORT,
		Slippage:       10,
		Type:           common.PerpOrderType_POT_LIMIT,
		Amount:         1,
		Price:          1000,
		ClientOrderID:  2,
	}
	sig, err := h.PostPerpOrder(ctx, request)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostPerpOrder signature : %s", sig)
	return false
}

func callPostDriftPerpOrder(ownerAddr string) bool {
	log.Info("starting callPostDriftPerpOrder test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	request := &pb.PostDriftPerpOrderRequest{
		OwnerAddress:   ownerAddr,
		Contract:       "SOL_PERP",
		AccountAddress: "",
		PositionSide:   "SHORT",
		Slippage:       10,
		Type:           "LIMIT",
		Amount:         1,
		Price:          1000,
		ClientOrderID:  2,
	}
	sig, err := h.PostDriftPerpOrder(ctx, request)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostDriftPerpOrder signature : %s", sig)
	return false
}

func callPostModifyOrder(ownerAddr string) bool {
	log.Info("starting callPostModifyOrder test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	request := &pb.PostModifyDriftOrderRequest{
		OwnerAddress:    ownerAddr,
		AccountAddress:  "",
		NewLimitPrice:   1000,
		NewPositionSide: "long",
		OrderID:         1,
	}
	sig, err := h.PostModifyDriftOrder(ctx, request)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostModifyOrder signature : %s", sig)
	return false
}

func callPostMarginOrder(ownerAddr string) bool {
	log.Info("starting callPostMarginOrder test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	request := &pb.PostDriftMarginOrderRequest{
		OwnerAddress:   ownerAddr,
		Market:         "SOL",
		AccountAddress: "",
		PositionSide:   "short",
		Slippage:       10,
		Type:           "limit",
		Amount:         1,
		Price:          1000,
		ClientOrderID:  2,
	}
	sig, err := h.PostDriftMarginOrder(ctx, request)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostMarginOrder signature : %s", sig)
	return false
}

func callManageCollateralWithdraw() bool {
	log.Info("starting callManageCollateralWithdraw test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := h.PostManageCollateral(ctx, &pb.PostManageCollateralRequest{
		Project:        pb.Project_P_DRIFT,
		Amount:         1,
		AccountAddress: "61bvX2qCwzPKNztgVQF3ktDHM2hZGdivCE28RrC99EAS",
		Type:           common.PerpCollateralType_PCT_WITHDRAWAL,
		Token:          common.PerpCollateralToken_PCTK_SOL,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callManageCollateralWithdraw signature : %s", sig)
	return false
}

func callPostDriftManageCollateralWithdraw() bool {
	log.Info("starting callPostDriftManageCollateralWithdraw test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := h.PostDriftManageCollateral(ctx, &pb.PostDriftManageCollateralRequest{
		Amount:         1,
		AccountAddress: "61bvX2qCwzPKNztgVQF3ktDHM2hZGdivCE28RrC99EAS",
		Type:           "WITHDRAWAL",
		Token:          "SOL",
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostDriftManageCollateralWithdraw signature : %s", sig)
	return false
}

func callManageCollateralTransfer() bool {
	log.Info("starting callManageCollateralTransfer test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := h.PostManageCollateral(ctx, &pb.PostManageCollateralRequest{
		Project:          pb.Project_P_DRIFT,
		Amount:           1,
		AccountAddress:   "61bvX2qCwzPKNztgVQF3ktDHM2hZGdivCE28RrC99EAS",
		Type:             common.PerpCollateralType_PCT_TRANSFER,
		Token:            common.PerpCollateralToken_PCTK_SOL,
		ToAccountAddress: "BTHDMaruPPTyUAZDv6w11qSMtyNAaNX6zFTPPepY863V",
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callManageCollateralTransfer signature : %s", sig)
	return false
}

func callPostDriftManageCollateralTransfer() bool {
	log.Info("starting callPostDriftManageCollateralTransfer test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := h.PostDriftManageCollateral(ctx, &pb.PostDriftManageCollateralRequest{
		Amount:           1,
		AccountAddress:   "61bvX2qCwzPKNztgVQF3ktDHM2hZGdivCE28RrC99EAS",
		Type:             "TRANSFER",
		Token:            "SOL",
		ToAccountAddress: "BTHDMaruPPTyUAZDv6w11qSMtyNAaNX6zFTPPepY863V",
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostDriftManageCollateralTransfer signature : %s", sig)
	return false
}

func callDriftEnableMarginTrading(ownerAddress string) bool {
	log.Info("starting callDriftEnableMarginTrading transfer test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := httpClient().PostDriftEnableMarginTrading(ctx, &pb.PostDriftEnableMarginTradingRequest{
		OwnerAddress: ownerAddress,
		EnableMargin: true,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callDriftEnableMarginTrading signature : %s", sig)
	return false
}

func callManageCollateralDeposit() bool {
	log.Info("starting callManageCollateral Deposit test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := h.PostManageCollateral(ctx, &pb.PostManageCollateralRequest{
		Project:        pb.Project_P_DRIFT,
		Amount:         1,
		AccountAddress: "61bvX2qCwzPKNztgVQF3ktDHM2hZGdivCE28RrC99EAS",
		Type:           common.PerpCollateralType_PCT_DEPOSIT,
		Token:          common.PerpCollateralToken_PCTK_SOL,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callManageCollateral signature : %s", sig)
	return false
}

func callPostDriftManageCollateralDeposit() bool {
	log.Info("starting callPostDriftManageCollateralDeposit Deposit test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := h.PostDriftManageCollateral(ctx, &pb.PostDriftManageCollateralRequest{
		Amount:         1,
		AccountAddress: "61bvX2qCwzPKNztgVQF3ktDHM2hZGdivCE28RrC99EAS",
		Type:           "DEPOSIT",
		Token:          "SOL",
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostDriftManageCollateralDeposit signature : %s", sig)
	return false
}

func callGetOpenPerpOrder(ownerAddr string) bool {
	log.Info("starting callGetOpenPerpOrder test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := h.GetOpenPerpOrder(ctx, &pb.GetOpenPerpOrderRequest{
		OwnerAddress:   ownerAddr,
		AccountAddress: "",
		Project:        pb.Project_P_DRIFT,
		OrderID:        1,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("GetOpenPerpOrder resp : %s", user)
	return false
}

func callGetDriftOpenPerpOrder(ownerAddr string) bool {
	log.Info("starting callGetDriftOpenPerpOrder test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := h.GetDriftOpenPerpOrder(ctx, &pb.GetDriftOpenPerpOrderRequest{
		OwnerAddress:   ownerAddr,
		AccountAddress: "",
		OrderID:        1,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetDriftOpenPerpOrder resp : %s", user)
	return false
}

func callGetDriftOpenPerpOrders(ownerAddr string) bool {
	log.Info("starting callGetDriftOpenPerpOrders test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := h.GetDriftOpenPerpOrders(ctx, &pb.GetDriftOpenPerpOrdersRequest{
		OwnerAddress:   ownerAddr,
		AccountAddress: "",
		Contracts:      []string{"SOL_PERP"},
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("GetDriftOpenPerpOrders resp : %s", user)
	return false
}

func callPostSettlePNL(ownerAddr string) bool {
	log.Info("starting callPostSettlePNL deposit test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := h.PostSettlePNL(ctx, &pb.PostSettlePNLRequest{
		Project:               pb.Project_P_DRIFT,
		OwnerAddress:          ownerAddr,
		SettleeAccountAddress: "9UnwdvTf5EfGeLyLrF4GZDUs7LKRUeJQzW7qsDVGQ8sS",
		Contract:              common.PerpContract_SOL_PERP,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostSettlePNL signature : %s", sig)
	return false
}

func callPostDriftSettlePNL(ownerAddr string) bool {
	log.Info("starting callPostDriftSettlePNL deposit test")

	h := httpClient()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := h.PostDriftSettlePNL(ctx, &pb.PostDriftSettlePNLRequest{
		OwnerAddress:          ownerAddr,
		SettleeAccountAddress: "9UnwdvTf5EfGeLyLrF4GZDUs7LKRUeJQzW7qsDVGQ8sS",
		Contract:              "SOL_PERP",
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostDriftSettlePNL signature : %s", sig)
	return false
}

func callPostSettlePNLs(ownerAddr string) bool {
	log.Info("starting callPostSettlePNLs test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := httpClient().PostSettlePNLs(ctx, &pb.PostSettlePNLsRequest{
		Project:                 pb.Project_P_DRIFT,
		OwnerAddress:            ownerAddr,
		SettleeAccountAddresses: []string{"9UnwdvTf5EfGeLyLrF4GZDUs7LKRUeJQzW7qsDVGQ8sS"},
		Contract:                common.PerpContract_SOL_PERP,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostSettlePNL signature : %s", sig)
	return false
}

func callPostDriftSettlePNLs(ownerAddr string) bool {
	log.Info("starting callPostDriftSettlePNLs test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := httpClient().PostDriftSettlePNLs(ctx, &pb.PostDriftSettlePNLsRequest{
		OwnerAddress:            ownerAddr,
		SettleeAccountAddresses: []string{"9UnwdvTf5EfGeLyLrF4GZDUs7LKRUeJQzW7qsDVGQ8sS"},
		Contract:                "SOL_PERP",
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostDriftSettlePNLs signature : %s", sig)
	return false
}

func callGetAssets(ownerAddr string) bool {
	log.Info("starting callGetAssets test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := httpClient().GetAssets(ctx, &pb.GetAssetsRequest{
		OwnerAddress:   ownerAddr,
		AccountAddress: "",
		Project:        pb.Project_P_DRIFT,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetAssets resp : %s", user)
	return false
}

func callGetDriftAssets(ownerAddr string) bool {
	log.Info("starting callGetDriftAssets test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := httpClient().GetDriftAssets(ctx, &pb.GetDriftAssetsRequest{
		OwnerAddress:   ownerAddr,
		AccountAddress: "",
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetDriftAssets resp : %s", user)
	return false
}

func callGetPerpContracts() bool {
	log.Info("starting callGetPerpContracts test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := httpClient().GetPerpContracts(ctx, &pb.GetPerpContractsRequest{
		Project: pb.Project_P_DRIFT,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetPerpContracts resp : %s", user)
	return false
}

func callGetDriftPerpContracts() bool {
	log.Info("starting callGetDriftPerpContracts test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := httpClient().GetDriftPerpContracts(ctx, &pb.GetDriftPerpContractsRequest{})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetDriftPerpContracts resp : %s", user)
	return false
}

func callGetDriftPerpOrderbook() bool {
	log.Info("starting callGetDriftPerpOrderbook test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := httpClient().GetDriftPerpOrderbook(ctx, &pb.GetDriftPerpOrderbookRequest{
		Contract: "SOL_PERP",
		Limit:    12,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetDriftPerpOrderbook resp : %s", user)
	return false
}

func callGetDriftMarkets() bool {
	log.Info("starting callGetDriftMarkets test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := httpClient().GetDriftMarkets(ctx, &pb.GetDriftMarketsRequest{})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("GetDriftMarkets resp : %s", user)
	return false
}

func callPostLiquidatePerp(ownerAddr string) bool {
	log.Info("starting callPostLiquidatePerp test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := httpClient().PostLiquidatePerp(ctx, &pb.PostLiquidatePerpRequest{
		Project:               pb.Project_P_DRIFT,
		OwnerAddress:          ownerAddr,
		Amount:                1,
		Contract:              common.PerpContract_SOL_PERP,
		SettleeAccountAddress: "9UnwdvTf5EfGeLyLrF4GZDUs7LKRUeJQzW7qsDVGQ8sS",
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostLiquidatePerp signature : %s", sig)
	return false
}

func callPostLiquidateDriftPerp(ownerAddr string) bool {
	log.Info("starting callPostLiquidateDriftPerp test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := httpClient().PostLiquidateDriftPerp(ctx, &pb.PostLiquidateDriftPerpRequest{
		OwnerAddress:          ownerAddr,
		Amount:                1,
		Contract:              "SOL_PERP",
		SettleeAccountAddress: "9UnwdvTf5EfGeLyLrF4GZDUs7LKRUeJQzW7qsDVGQ8sS",
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostLiquidateDriftPerp signature : %s", sig)
	return false
}
