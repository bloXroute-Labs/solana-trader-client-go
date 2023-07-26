package main

import (
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/examples/config"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	"github.com/bloXroute-Labs/solana-trader-proto/common"
	"math/rand"
	"os"
	"time"

	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func main() {
	utils.InitLogger()
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

	var w *provider.WSClient

	switch cfg.Env {
	case config.EnvLocal:
		w, err = provider.NewWSClientLocal()
	case config.EnvTestnet:
		w, err = provider.NewWSClientTestnet()
	case config.EnvMainnet:
		w, err = provider.NewWSClient()
	}
	if err != nil {
		log.Fatalf("error dialing WS client: %v", err)
	}
	defer func(w *provider.WSClient) {
		err := w.Close()
		if err != nil {
			panic(err)
		}
	}(w)

	var failed bool

	// informational requests
	failed = failed || logCall("callMarketsWS", func() bool { return callMarketsWS(w) })
	failed = failed || logCall("callOrderbookWS", func() bool { return callOrderbookWS(w) })
	failed = failed || logCall("callMarketDepthWS", func() bool { return callMarketDepthWS(w) })
	failed = failed || logCall("callTradesWS", func() bool { return callTradesWS(w) })
	failed = failed || logCall("callPoolsWS", func() bool { return callPoolsWS(w) })
	failed = failed || logCall("callRaydiumPools", func() bool { return callRaydiumPoolsWS(w) })
	failed = failed || logCall("callRaydiumPrices", func() bool { return callRaydiumPricesWS(w) })
	failed = failed || logCall("callJupiterPrices", func() bool { return callJupiterPricesWS(w) })
	failed = failed || logCall("callPriceWS", func() bool { return callPriceWS(w) })
	failed = failed || logCall("callOpenOrdersWS", func() bool { return callOpenOrdersWS(w) })
	failed = failed || logCall("callTickersWS", func() bool { return callTickersWS(w) })
	failed = failed || logCall("callUnsettledWS", func() bool { return callUnsettledWS(w) })
	failed = failed || logCall("callAccountBalanceWS", func() bool { return callAccountBalanceWS(w) })
	failed = failed || logCall("callGetQuotes", func() bool { return callGetQuotes(w) })
	failed = failed || logCall("callGetRaydiumQuotes", func() bool { return callGetRaydiumQuotes(w) })
	failed = failed || logCall("callGetJupiterQuotes", func() bool { return callGetJupiterQuotes(w) })
	failed = failed || logCall("callDriftPerpOrderbookWS", func() bool { return callDriftPerpOrderbookWS(w) })
	failed = failed || logCall("callDriftGetMarginOrderbookWS", func() bool { return callDriftGetMarginOrderbookWS(w) })
	failed = failed || logCall("callGetDriftMarketDepthWS", func() bool { return callGetDriftMarketDepthWS(w) })

	// streaming methods
	failed = failed || logCall("callOrderbookWSStream", func() bool { return callOrderbookWSStream(w) })
	failed = failed || logCall("callMarketDepthWSStream", func() bool { return callMarketDepthWSStream(w) })
	failed = failed || logCall("callRecentBlockHashWSStream", func() bool { return callRecentBlockHashWSStream(w) })
	failed = failed || logCall("callPoolReservesWSStream", func() bool { return callPoolReservesWSStream(w) })
	failed = failed || logCall("callBlockWSStream", func() bool { return callBlockWSStream(w) })
	failed = failed || logCall("callDriftPerpOrderbookWSStream", func() bool { return callDriftPerpOrderbookWSStream(w) })
	failed = failed || logCall("callDriftMarginOrderbooksWSStream", func() bool { return callDriftMarginOrderbooksWSStream(w) })
	failed = failed || logCall("callDriftGetPerpTradesStream", func() bool { return callDriftGetPerpTradesStream(w) })
	failed = failed || logCall("callDriftMarketDepthsStream", func() bool { return callDriftMarketDepthsStream(w) })

	if cfg.RunSlowStream {
		failed = failed || logCall("callPricesWSStream", func() bool { return callPricesWSStream(w) })
		failed = failed || logCall("callSwapsWSStream", func() bool { return callSwapsWSStream(w) })
		failed = failed || logCall("callTradesWSStream", func() bool { return callTradesWSStream(w) })
	}

	// calls below this place an order and immediately cancel it
	// you must specify:
	//	- PRIVATE_KEY (by default loaded during provider.NewWSClient()) to sign transactions
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
		/*failed = failed || logCall("orderLifecycleTest", func() bool { return orderLifecycleTest(w, ownerAddr, payerAddr, ooAddr) })
		failed = failed || logCall("cancelAll", func() bool { return cancelAll(w, ownerAddr, payerAddr, ooAddr) })
		failed = failed || logCall("callReplaceByClientOrderID", func() bool { return callReplaceByClientOrderID(w, ownerAddr, payerAddr, ooAddr) })*/
		failed = failed || logCall("callReplaceOrder", func() bool { return callReplaceOrder(w, ownerAddr, payerAddr, ooAddr) })
		failed = failed || logCall("callRecentBlockHashWSStream", func() bool { return callRecentBlockHashWSStream(w) })
		failed = failed || logCall("callTradeSwap", func() bool { return callTradeSwap(w, ownerAddr) })
		failed = failed || logCall("callRouteTradeSwap", func() bool { return callRouteTradeSwap(w, ownerAddr) })
		failed = failed || logCall("callRaydiumTradeSwap", func() bool { return callRaydiumSwap(w, ownerAddr) })
		failed = failed || logCall("callJupiterTradeSwap", func() bool { return callJupiterSwap(w, ownerAddr) })
		failed = failed || logCall("callRaydiumRouteTradeSwap", func() bool { return callRaydiumRouteSwap(w, ownerAddr) })
		failed = failed || logCall("callJupiterRouteTradeSwap", func() bool { return callJupiterRouteSwap(w, ownerAddr) })
	}

	failed = failed || logCall("callGetOpenPerpOrders", func() bool { return callGetOpenPerpOrders(w, ownerAddr) })
	failed = failed || logCall("callGetDriftOpenPerpOrders", func() bool { return callGetDriftOpenPerpOrders(w, ownerAddr) })
	failed = failed || logCall("callGetDriftOpenMarginOrders", func() bool { return callGetDriftOpenMarginOrders(w, ownerAddr) })
	failed = failed || logCall("callGetPerpPositions", func() bool { return callGetPerpPositions(w, ownerAddr) })
	failed = failed || logCall("callGetDriftPerpPositions", func() bool { return callGetDriftPerpPositions(w, ownerAddr) })
	failed = failed || logCall("callGetUser", func() bool { return callGetUser(w, ownerAddr) })

	failed = failed || logCall("callGetOpenPerpOrder", func() bool { return callGetOpenPerpOrder(w, ownerAddr) })
	failed = failed || logCall("callGetAssets", func() bool { return callGetAssets(w, ownerAddr) })
	failed = failed || logCall("callGetPerpContracts", func() bool { return callGetPerpContracts(w) })
	failed = failed || logCall("callGetDriftMarkets", func() bool { return callGetDriftMarkets(w) })

	failed = failed || logCall("callGetDriftAssets", func() bool { return callGetDriftAssets(w, ownerAddr) })
	failed = failed || logCall("callGetDriftPerpContracts", func() bool { return callGetDriftPerpContracts(w) })
	failed = failed || logCall("callGetDriftPerpOrderbook", func() bool { return callGetDriftPerpOrderbook(w, ownerAddr) })
	failed = failed || logCall("callGetDriftUser", func() bool { return callGetDriftUser(w, ownerAddr) })
	failed = failed || logCall("callGetDriftOpenPerpOrder", func() bool { return callGetDriftOpenPerpOrder(w, ownerAddr) })
	failed = failed || logCall("callGetDriftOpenMarginOrder", func() bool { return callGetDriftOpenMarginOrder(w, ownerAddr) })

	if cfg.RunPerpTrades {
		failed = failed || logCall("callCancelPerpOrder", func() bool { return callCancelPerpOrder(w, ownerAddr) })
		failed = failed || logCall("callDriftCancelPerpOrder", func() bool { return callDriftCancelPerpOrder(w, ownerAddr) })
		failed = failed || logCall("callCancelDriftMarginOrder", func() bool { return callCancelDriftMarginOrder(w, ownerAddr) })
		failed = failed || logCall("callClosePerpPositions", func() bool { return callClosePerpPositions(w, ownerAddr) })
		failed = failed || logCall("callCreateUser", func() bool { return callCreateUser(w, ownerAddr) })
		failed = failed || logCall("callManageCollateralDeposit", func() bool { return callManageCollateralDeposit(w) })
		failed = failed || logCall("callPostPerpOrder", func() bool { return callPostPerpOrder(w, ownerAddr) })
		failed = failed || logCall("callPostDriftPerpOrder", func() bool { return callPostDriftPerpOrder(w, ownerAddr) })
		failed = failed || logCall("callPostModifyOrder", func() bool { return callPostModifyOrder(w, ownerAddr) })
		failed = failed || logCall("callPostMarginOrder", func() bool { return callPostMarginOrder(w, ownerAddr) })
		failed = failed || logCall("callManageCollateralWithdraw", func() bool { return callManageCollateralWithdraw(w) })
		failed = failed || logCall("callManageCollateralTransfer", func() bool { return callManageCollateralTransfer(w) })
		failed = failed || logCall("callDriftEnableMarginTrading", func() bool { return callDriftEnableMarginTrading(w, ownerAddr) })

		failed = failed || logCall("callPostSettlePNL", func() bool { return callPostSettlePNL(w, ownerAddr) })
		failed = failed || logCall("callPostSettlePNLs", func() bool { return callPostSettlePNLs(w, ownerAddr) })
		failed = failed || logCall("callPostLiquidatePerp", func() bool { return callPostLiquidatePerp(w, ownerAddr) })

		failed = failed || logCall("callPostCloseDriftPerpPositions", func() bool { return callPostCloseDriftPerpPositions(w, ownerAddr) })
		failed = failed || logCall("callPostCreateDriftUser", func() bool { return callPostCreateDriftUser(w, ownerAddr) })
		failed = failed || logCall("callPostDriftManageCollateralDeposit", func() bool { return callPostDriftManageCollateralDeposit(w) })
		failed = failed || logCall("callPostDriftManageCollateralWithdraw", func() bool { return callPostDriftManageCollateralWithdraw(w) })
		failed = failed || logCall("callPostDriftManageCollateralTransfer", func() bool { return callPostDriftManageCollateralTransfer(w) })
		failed = failed || logCall("callPostDriftSettlePNL", func() bool { return callPostDriftSettlePNL(w, ownerAddr) })
		failed = failed || logCall("callPostDriftSettlePNLs", func() bool { return callPostDriftSettlePNLs(w, ownerAddr) })
		failed = failed || logCall("callPostLiquidateDriftPerp", func() bool { return callPostLiquidateDriftPerp(w, ownerAddr) })
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

func callMarketsWS(w *provider.WSClient) bool {
	log.Info("fetching markets...")

	markets, err := w.GetMarketsV2(context.Background())
	if err != nil {
		log.Errorf("error with GetMarkets request: %v", err)
		return true
	} else {
		log.Info(markets)
	}

	fmt.Println()
	return false
}

func callOrderbookWS(w *provider.WSClient) bool {
	log.Info("fetching orderbooks...")

	orderbook, err := w.GetOrderbookV2(context.Background(), "SOL-USDT", 0)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL-USDT: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()

	orderbook, err = w.GetOrderbookV2(context.Background(), "SOLUSDT", 2)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL-USDT: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()

	orderbook, err = w.GetOrderbookV2(context.Background(), "SOL:USDT", 3)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL:USDC: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()
	return false
}

func callMarketDepthWS(w *provider.WSClient) bool {
	log.Info("fetching market depth data...")

	mktDepth, err := w.GetMarketDepthV2(context.Background(), "SOL:USDT", 3)
	if err != nil {
		log.Errorf("error with GetMarketDepth request for SOL:USDC: %v", err)
		return true
	} else {
		log.Info(mktDepth)
	}

	fmt.Println()
	return false
}

func callTradesWS(w *provider.WSClient) bool {
	log.Info("fetching trades...")

	trades, err := w.GetTrades(context.Background(), "SOLUSDC", 3, pb.Project_P_OPENBOOK)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL:USDC: %v", err)
		return true
	} else {
		log.Info(trades)
	}

	fmt.Println()
	return false
}

func callPoolsWS(w *provider.WSClient) bool {
	log.Info("fetching pools...")

	pools, err := w.GetPools(context.Background(), []pb.Project{pb.Project_P_RAYDIUM})
	if err != nil {
		log.Errorf("error with GetPools request for Raydium: %v", err)
		return true
	} else {
		log.Info(pools)
	}

	fmt.Println()
	return false
}

func callRaydiumPoolsWS(w *provider.WSClient) bool {
	log.Info("fetching Raydium pools...")

	pools, err := w.GetRaydiumPools(context.Background(), &pb.GetRaydiumPoolsRequest{})
	if err != nil {
		log.Errorf("error with GetRaydiumPools request for Raydium: %v", err)
		return true
	} else {
		log.Info(pools)
	}

	fmt.Println()
	return false
}

func callPriceWS(w *provider.WSClient) bool {
	log.Info("fetching prices...")

	pools, err := w.GetPrice(context.Background(), []string{"SOL", "ETH"})
	if err != nil {
		log.Errorf("error with GetPrice request for SOL and ETH: %v", err)
		return true
	} else {
		log.Info(pools)
	}

	return false
}

func callRaydiumPricesWS(w *provider.WSClient) bool {
	log.Info("fetching Raydium prices...")

	pools, err := w.GetRaydiumPrices(context.Background(), &pb.GetRaydiumPricesRequest{
		Tokens: []string{"SOL", "ETH"},
	})
	if err != nil {
		log.Errorf("error with GetRaydiumPrices request for SOL and ETH: %v", err)
		return true
	} else {
		log.Info(pools)
	}

	return false
}

func callJupiterPricesWS(w *provider.WSClient) bool {
	log.Info("fetching Jupiter prices...")

	pools, err := w.GetJupiterPrices(context.Background(), &pb.GetJupiterPricesRequest{
		Tokens: []string{"SOL", "ETH"},
	})
	if err != nil {
		log.Errorf("error with GetJupiterPrices request for SOL and ETH: %v", err)
		return true
	} else {
		log.Info(pools)
	}

	return false
}

func callOpenOrdersWS(w *provider.WSClient) bool {
	log.Info("fetching open orders...")

	orders, err := w.GetOpenOrdersV2(context.Background(), "SOLUSDC", "FFqDwRq8B4hhFKRqx7N1M6Dg6vU699hVqeynDeYJdPj5", "", "", 0)
	if err != nil {
		log.Errorf("error with GetOrders request for SOL-USDT: %v", err)
		return true
	} else {
		log.Info(orders)
	}

	fmt.Println()
	return false
}

func callUnsettledWS(w *provider.WSClient) bool {
	log.Info("fetching unsettled...")

	response, err := w.GetUnsettledV2(context.Background(), "SOLUSDC", "AFT8VayE7qr8MoQsW3wHsDS83HhEvhGWdbNSHRKeUDfQ")
	if err != nil {
		log.Errorf("error with GetUnsettled request for SOL-USDT: %v", err)
		return true
	} else {
		log.Info(response)
	}

	fmt.Println()
	return false
}

func callAccountBalanceWS(w *provider.WSClient) bool {
	log.Info("fetching balances...")

	response, err := w.GetAccountBalance(context.Background(), "AFT8VayE7qr8MoQsW3wHsDS83HhEvhGWdbNSHRKeUDfQ")
	if err != nil {
		log.Errorf("error with GetAccountBalance request for AFT8VayE7qr8MoQsW3wHsDS83HhEvhGWdbNSHRKeUDfQ: %v", err)
		return true
	} else {
		log.Info(response)
	}

	fmt.Println()
	return false
}

func callTickersWS(w *provider.WSClient) bool {
	log.Info("fetching tickers...")

	tickers, err := w.GetTickersV2(context.Background(), "SOLUSDC")
	if err != nil {
		log.Errorf("error with GetTickers request for SOL-USDT: %v", err)
		return true
	} else {
		log.Info(tickers)
	}

	fmt.Println()
	return false
}

func callGetQuotes(w *provider.WSClient) bool {
	log.Info("fetching quotes...")

	inToken := "SOL"
	outToken := "USDT"
	amount := 0.01
	slippage := float64(5)
	limit := 5

	quotes, err := w.GetQuotes(context.Background(), inToken, outToken, amount, slippage, int32(limit), []pb.Project{pb.Project_P_ALL})
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

func callGetRaydiumQuotes(w *provider.WSClient) bool {
	log.Info("fetching Raydium quotes...")

	inToken := "SOL"
	outToken := "USDT"
	amount := 0.01
	slippage := float64(5)

	quotes, err := w.GetRaydiumQuotes(context.Background(), &pb.GetRaydiumQuotesRequest{
		InToken:  inToken,
		OutToken: outToken,
		InAmount: amount,
		Slippage: slippage,
	})
	if err != nil {
		log.Errorf("error with GetRaydiumQuotes request for %s to %s: %v", inToken, outToken, err)
		return true
	}

	if len(quotes.Routes) != 1 {
		log.Errorf("did not get back 1 quote, got %v quotes", len(quotes.Routes))
		return true
	}
	for _, route := range quotes.Routes {
		log.Infof("best route for Raydium is %v", route)
	}

	fmt.Println()
	return false
}

func callGetJupiterQuotes(w *provider.WSClient) bool {
	log.Info("fetching Jupiter quotes...")

	inToken := "SOL"
	outToken := "USDT"
	amount := 0.01
	slippage := float64(5)
	limit := int32(3)

	quotes, err := w.GetJupiterQuotes(context.Background(), &pb.GetJupiterQuotesRequest{
		InToken:  inToken,
		OutToken: outToken,
		InAmount: amount,
		Slippage: slippage,
		Limit:    limit,
	})
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

// Stream response
func callOrderbookWSStream(w *provider.WSClient) bool {
	log.Info("starting orderbook stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := w.GetOrderbooksStream(ctx, []string{"SOL/USDC"}, 3, pb.Project_P_OPENBOOK)
	if err != nil {
		log.Errorf("error with GetOrderbooksStream request for SOL/USDC: %v", err)
		return true
	}

	orderbookCh := stream.Channel(0)
	for i := 1; i <= 1; i++ {
		_, ok := <-orderbookCh
		if !ok {
			return true
		}
		log.Infof("response %v received", i)
	}
	return false
}

// Stream response
func callMarketDepthWSStream(w *provider.WSClient) bool {
	log.Info("starting market depth stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := w.GetMarketDepthsStream(ctx, []string{"SOL/USDC"}, 3, pb.Project_P_OPENBOOK)
	if err != nil {
		log.Errorf("error with GetMarketDepthsStream request for SOL/USDC: %v", err)
		return true
	}

	mktDepthDataCh := stream.Channel(0)
	for i := 1; i <= 1; i++ {
		_, ok := <-mktDepthDataCh
		if !ok {
			return true
		}
		log.Infof("response %v received", i)
	}
	return false
}

func callTradesWSStream(w *provider.WSClient) bool {
	log.Info("starting trades stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tradesChan := make(chan *pb.GetTradesStreamResponse)
	stream, err := w.GetTradesStream(ctx, "SOL/USDC", 3, pb.Project_P_OPENBOOK)
	if err != nil {
		log.Errorf("error with GetTradesStream request for SOL/USDC: %v", err)
		return true
	}

	stream.Into(tradesChan)
	for i := 1; i <= 1; i++ {
		_, ok := <-tradesChan
		if !ok {
			return true
		}
		log.Infof("response %v received", i)
	}
	return false
}

// Stream response
func callRecentBlockHashWSStream(w *provider.WSClient) bool {
	log.Info("starting recent block hash stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := w.GetRecentBlockHashStream(ctx)
	if err != nil {
		log.Errorf("error with GetRecentBlockHashStream request: %v", err)
		return true
	}

	ch := stream.Channel(0)
	for i := 1; i <= 1; i++ {
		_, ok := <-ch
		if !ok {
			return true
		}
		log.Infof("response %v received", i)
	}
	return false
}

func callPoolReservesWSStream(w *provider.WSClient) bool {
	log.Info("starting pool reserves stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := w.GetPoolReservesStream(ctx, []pb.Project{pb.Project_P_RAYDIUM})
	if err != nil {
		log.Errorf("error with GetPoolReserves stream request: %v", err)
		return true
	}

	ch := stream.Channel(0)
	for i := 1; i <= 1; i++ {
		_, ok := <-ch
		if !ok {
			return true
		}
		log.Infof("response %v received", i)
	}
	return false
}

const (
	// SOL/USDC market
	marketAddr = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh6"

	orderSide   = pb.Side_S_ASK
	orderType   = common.OrderType_OT_LIMIT
	orderPrice  = float64(170200)
	orderAmount = float64(0.1)
)

func orderLifecycleTest(w *provider.WSClient, ownerAddr, payerAddr, ooAddr string) bool {
	log.Info("starting order lifecycle test")
	fmt.Println()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan *pb.GetOrderStatusStreamResponse)
	errCh := make(chan error)
	stream, err := w.GetOrderStatusStream(ctx, marketAddr, ownerAddr, pb.Project_P_OPENBOOK)
	if err != nil {
		log.Errorf("error getting order status stream %v", err)
		errCh <- err
	}
	stream.Into(ch)

	time.Sleep(time.Second * 10)

	clientOrderID, fail := callPlaceOrderWS(w, ownerAddr, payerAddr, ooAddr)
	if fail {
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

	fail = callCancelByClientOrderIDWS(w, ownerAddr, ooAddr, clientOrderID)
	if fail {
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
	case <-time.After(time.Second * 60):
		log.Error("no updates after cancelling order")
		return true
	}

	fmt.Println()
	callPostSettleWS(w, ownerAddr, ooAddr)
	return false
}

func callPlaceOrderWS(w *provider.WSClient, ownerAddr, payerAddr, ooAddr string) (uint64, bool) {
	log.Info("trying to place an order")

	// generate a random clientOrderId for this order
	rand.Seed(time.Now().UnixNano())
	clientOrderID := rand.Uint64()

	opts := provider.PostOrderOpts{
		ClientOrderID:     clientOrderID,
		OpenOrdersAddress: ooAddr,
	}

	// sign/submit transaction after creation
	sig, err := w.SubmitOrderV2(context.Background(), ownerAddr, payerAddr, marketAddr,
		orderSide, common.OrderType_OT_LIMIT, orderAmount, orderPrice, opts)
	if err != nil {
		log.Errorf("failed to submit order (%v)", err)
		return 0, true
	}

	log.Infof("placed order %v with clientOrderID %v", sig, clientOrderID)

	return clientOrderID, false
}

func callCancelByClientOrderIDWS(w *provider.WSClient, ownerAddr, ooAddr string, clientOrderID uint64) bool {
	log.Info("trying to cancel order")

	_, err := w.SubmitCancelOrderV2(context.Background(), &pb.PostCancelOrderRequestV2{
		OrderID:           "",
		Side:              pb.Side_S_ASK.String(),
		MarketAddress:     marketAddr,
		OwnerAddress:      ownerAddr,
		OpenOrdersAddress: ooAddr,
		ClientOrderID:     clientOrderID,
	}, true)
	if err != nil {
		log.Errorf("failed to cancel order by client ID (%v)", err)
		return true
	}

	log.Infof("canceled order for clientOrderID %v", clientOrderID)
	return false
}

func callPostSettleWS(w *provider.WSClient, ownerAddr, ooAddr string) bool {
	log.Info("starting post settle")

	sig, err := w.SubmitSettleV2(context.Background(), ownerAddr, "SOL/USDC", "F75gCEckFAyeeCWA9FQMkmLCmke7ehvBnZeVZ3QgvJR7",
		"4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv", ooAddr, false)
	if err != nil {
		log.Errorf("error with post transaction stream request for SOL/USDC: %v", err)
		return true
	}

	log.Infof("response signature received: %v", sig)
	return false
}

func cancelAll(w *provider.WSClient, ownerAddr, payerAddr, ooAddr string) bool {
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
	sig, err := w.SubmitOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, common.OrderType_OT_LIMIT, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #1, signature %s", sig)

	opts.ClientOrderID = clientOrderID2
	sig, err = w.SubmitOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, common.OrderType_OT_LIMIT, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check orders are there
	orders, err := w.GetOpenOrdersV2(ctx, marketAddr, ownerAddr, "", "", 0)
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
	sigs, err := w.SubmitCancelOrderV2(ctx, &pb.PostCancelOrderRequestV2{
		OrderID:           "",
		Side:              pb.Side_S_ASK.String(),
		MarketAddress:     marketAddr,
		OwnerAddress:      ownerAddr,
		OpenOrdersAddress: ooAddr,
		ClientOrderID:     0,
	}, true)
	if err != nil {
		log.Error(err)
		return true
	}
	for _, tx := range sigs.Transactions {
		log.Infof("placing cancel order(s) %s", tx.Signature)
	}

	time.Sleep(time.Minute)

	orders, err = w.GetOpenOrdersV2(ctx, marketAddr, ownerAddr, "", "", 0)
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
	callPostSettleWS(w, ownerAddr, ooAddr)
	return false
}

func callReplaceByClientOrderID(w *provider.WSClient, ownerAddr, payerAddr, ooAddr string) bool {
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
	sig, err := w.SubmitOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, common.OrderType_OT_LIMIT, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	} else {
		log.Infof("submitting place order #1, signature %s", sig)
	}
	time.Sleep(time.Minute)
	// Check order is there
	orders, err := w.GetOpenOrdersV2(ctx, marketAddr, ownerAddr, "", "", 0)
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
	sig, err = w.SubmitReplaceOrderV2(ctx, "", ownerAddr, payerAddr, marketAddr, orderSide, common.OrderType_OT_LIMIT, orderAmount, orderPrice/2, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check order #2 is in orderbook
	orders, err = w.GetOpenOrdersV2(ctx, marketAddr, ownerAddr, "", "", 0)
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
	} else {
		log.Info("order #2 placed successfully")
	}
	time.Sleep(time.Minute)
	// Cancel all the orders
	log.Info("cancelling the orders")
	sigs, err := w.SubmitCancelOrderV2(ctx, &pb.PostCancelOrderRequestV2{
		OrderID:           "",
		Side:              pb.Side_S_ASK.String(),
		MarketAddress:     marketAddr,
		OwnerAddress:      ownerAddr,
		OpenOrdersAddress: ooAddr,
		ClientOrderID:     0,
	}, true)
	if err != nil {
		log.Error(err)
		return true
	}
	for _, tx := range sigs.Transactions {
		log.Infof("placing cancel order(s) %s", tx.Signature)
	}
	return false
}

func callReplaceOrder(w *provider.WSClient, ownerAddr, payerAddr, ooAddr string) bool {
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
	sig, err := w.SubmitOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, common.OrderType_OT_LIMIT, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #1, signature %s", sig)

	time.Sleep(time.Minute)
	// Check orders are there
	orders, err := w.GetOpenOrdersV2(ctx, marketAddr, ownerAddr, "", "", 0)
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
	sig, err = w.SubmitReplaceOrderV2(ctx, found1.OrderID, ownerAddr, payerAddr, marketAddr, pb.Side_S_ASK, common.OrderType_OT_LIMIT, orderAmount, orderPrice/2, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check orders are there
	orders, err = w.GetOpenOrdersV2(ctx, marketAddr, ownerAddr, "", "", 0)
	if err != nil {
		log.Error(err)
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
	sigs, err := w.SubmitCancelOrderV2(ctx, &pb.PostCancelOrderRequestV2{
		OrderID:           "",
		Side:              pb.Side_S_ASK.String(),
		MarketAddress:     marketAddr,
		OwnerAddress:      ownerAddr,
		OpenOrdersAddress: ooAddr,
		ClientOrderID:     0,
	}, true)
	if err != nil {
		log.Error(err)
		return true
	}
	for _, tx := range sigs.Transactions {
		log.Infof("placing cancel order(s) %s", tx.Signature)
	}
	return false
}

func callTradeSwap(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting trade swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("trade swap")
	sig, err := w.SubmitTradeSwap(ctx, ownerAddr, "USDT",
		"SOL", 0.01, 0.1, "raydium", provider.SubmitOpts{
			SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
			SkipPreFlight:  false,
		})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("trade swap transaction signature : %s", sig)
	return false
}

func callRaydiumSwap(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting Raydium swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig, err := w.SubmitRaydiumSwap(ctx, &pb.PostRaydiumSwapRequest{
		OwnerAddress: ownerAddr,
		InToken:      "USDT",
		OutToken:     "SOL",
		Slippage:     0.1,
		InAmount:     0.01,
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  false,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Raydium swap transaction signature : %s", sig)
	return false
}

func callRouteTradeSwap(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting route trade swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("route trade swap")
	sig, err := w.SubmitRouteTradeSwap(ctx, &pb.RouteTradeSwapRequest{
		OwnerAddress: ownerAddr,
		Project:      pb.Project_P_RAYDIUM,
		Slippage:     0.1,
		Steps: []*pb.RouteStep{
			{
				Project: &pb.StepProject{
					Label: "Raydium",
					Id:    "",
				},
				InToken:  "FIDA",
				OutToken: "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R",

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
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  false,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("route trade swap transaction signature : %s", sig)
	return false
}

func callRaydiumRouteSwap(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting Raydium swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig, err := w.SubmitRaydiumRouteSwap(ctx, &pb.PostRaydiumRouteSwapRequest{
		OwnerAddress: ownerAddr,
		Slippage:     0.1,
		Steps: []*pb.RaydiumRouteStep{
			{
				InToken:  "FIDA",
				OutToken: "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R",

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
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  false,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Raydium route swap transaction signature : %s", sig)
	return false
}

func callJupiterSwap(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting Jupiter swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig, err := w.SubmitJupiterSwap(ctx, &pb.PostJupiterSwapRequest{
		OwnerAddress: ownerAddr,
		InToken:      "USDT",
		OutToken:     "SOL",
		Slippage:     0.1,
		InAmount:     0.01,
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  false,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Jupiter swap transaction signature : %s", sig)
	return false
}

func callJupiterRouteSwap(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting Jupiter swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig, err := w.SubmitJupiterRouteSwap(ctx, &pb.PostJupiterRouteSwapRequest{
		OwnerAddress: ownerAddr,
		Slippage:     0.1,
		Steps: []*pb.JupiterRouteStep{
			{
				InToken:  "FIDA",
				OutToken: "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R",

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
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  false,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Jupiter route swap transaction signature : %s", sig)
	return false
}

func callPricesWSStream(w *provider.WSClient) bool {
	log.Info("starting prices stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := w.GetPricesStream(ctx, []pb.Project{pb.Project_P_RAYDIUM}, []string{"SOL"})
	if err != nil {
		log.Errorf("error with GetPrices stream request: %v", err)
		return true
	}

	ch := stream.Channel(0)
	for i := 1; i <= 1; i++ {
		_, ok := <-ch
		if !ok {
			return true
		}
		log.Infof("response %v received", i)
	}
	return false
}

func callSwapsWSStream(w *provider.WSClient) bool {
	log.Info("starting get swaps stream")

	ch := make(chan *pb.GetSwapsStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := w.GetSwapsStream(ctx, []pb.Project{pb.Project_P_RAYDIUM}, []string{"58oQChx4yWmvKdwLLZzBi4ChoCc2fqCUWBkwMihLYQo2"}, true) // SOL-USDC Raydium pool
	if err != nil {
		log.Errorf("error with GetSwaps stream request: %v", err)
		return true
	}
	stream.Into(ch)
	for i := 1; i <= 1; i++ {
		_, ok := <-ch
		if !ok {
			// channel closed
			return true
		}

		log.Infof("response %v received", i)
	}
	return false
}

func callBlockWSStream(w *provider.WSClient) bool {
	log.Info("starting get block stream")

	ch := make(chan *pb.GetBlockStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := w.GetBlockStream(ctx)
	if err != nil {
		log.Errorf("error with GetBlock stream request: %v", err)
		return true
	}
	stream.Into(ch)
	for i := 1; i <= 1; i++ {
		_, ok := <-ch
		if !ok {
			// channel closed
			return true
		}

		log.Infof("response %v received", i)
	}
	return false
}

func callDriftPerpOrderbookWS(w *provider.WSClient) bool {
	log.Info("fetching drift perp orderbooks...")

	orderbook, err := w.GetPerpOrderbook(context.Background(), &pb.GetPerpOrderbookRequest{
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

func callDriftPerpOrderbookWSStream(w *provider.WSClient) bool {
	log.Info("starting drift perp orderbook stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := w.GetPerpOrderbooksStream(ctx, &pb.GetPerpOrderbooksRequest{
		Contracts: []common.PerpContract{common.PerpContract_SOL_PERP},
		Limit:     0,
		Project:   pb.Project_P_DRIFT,
	})
	if err != nil {
		log.Errorf("error with GetPerpOrderbooksStream request for SOL-PERP: %v", err)
		return true
	}

	orderbookCh := stream.Channel(0)
	for i := 1; i <= 1; i++ {
		_, ok := <-orderbookCh
		if !ok {
			return true
		}
		log.Infof("response %v received", i)
	}
	return false
}

func callDriftGetMarginOrderbookWS(w *provider.WSClient) bool {
	log.Info("fetching drift spot orderbooks...")

	orderbook, err := w.GetDriftMarginOrderbook(context.Background(), &pb.GetDriftMarginOrderbookRequest{
		Market:   "SOL",
		Limit:    0,
		Metadata: true,
	})
	if err != nil {
		log.Errorf("error with GetMarginOrderbook request for SOL-MARGIN: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()
	return false
}

func callDriftMarginOrderbooksWSStream(w *provider.WSClient) bool {
	log.Info("starting drift spot orderbook stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := w.GetDriftMarginOrderbooksStream(ctx, &pb.GetDriftMarginOrderbooksRequest{
		Markets: []string{"SOL"},
		Limit:   0,
	})
	if err != nil {
		log.Errorf("error with GetMarginOrderbooksStream request for SOL-MARGIN: %v", err)
		return true
	}

	orderbookCh := stream.Channel(0)
	for i := 1; i <= 1; i++ {
		_, ok := <-orderbookCh
		if !ok {
			return true
		}
		log.Infof("response %v received", i)
	}
	return false
}

func callDriftGetPerpTradesStream(w *provider.WSClient) bool {
	log.Info("starting get Drift PerpTrades stream")

	ch := make(chan *pb.GetPerpTradesStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := w.GetPerpTradesStream(ctx, &pb.GetPerpTradesStreamRequest{
		Contracts: []common.PerpContract{
			common.PerpContract_SOL_PERP, common.PerpContract_ETH_PERP,
			common.PerpContract_BTC_PERP, common.PerpContract_APT_PERP,
		},
		Project: pb.Project_P_DRIFT,
	})
	if err != nil {
		log.Errorf("error with GetPerpTradesStream stream request: %v", err)
		return true
	}
	stream.Into(ch)
	for i := 1; i <= 1; i++ {
		_, ok := <-ch
		if !ok {
			// channel closed
			return true
		}

		log.Infof("response %v received", i)
	}
	return false
}

func callDriftMarketDepthsStream(w *provider.WSClient) bool {
	log.Info("starting get Drift MarketDepth stream")

	ch := make(chan *pb.GetDriftMarketDepthStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := w.GetDriftMarketDepthsStream(ctx, &pb.GetDriftMarketDepthsStreamRequest{
		Contracts: []string{"SOL_PERP", "ETH_PERP"},
		Limit:     2,
	})
	if err != nil {
		log.Errorf("error with GetDriftMarketDepthsStream stream request: %v", err)
		return true
	}
	stream.Into(ch)
	for i := 1; i <= 1; i++ {
		_, ok := <-ch
		if !ok {
			// channel closed
			return true
		}

		log.Infof("response %v received", i)
	}
	return false
}

func callGetDriftMarketDepthWS(w *provider.WSClient) bool {
	log.Info("starting callGetDriftMarketDepthWS test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := w.GetDriftMarketDepth(ctx, &pb.GetDriftMarketDepthRequest{
		Contract: "SOL_PERP",
		Limit:    2,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetDriftMarketDepthWS resp : %s", user)
	return false
}

func callGetOpenPerpOrders(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callGetOpenPerpOrders test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := w.GetOpenPerpOrders(ctx, &pb.GetOpenPerpOrdersRequest{
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

func callGetDriftOpenPerpOrders(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callGetDriftOpenPerpOrders test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := w.GetDriftOpenPerpOrders(ctx, &pb.GetDriftOpenPerpOrdersRequest{
		OwnerAddress:   ownerAddr,
		AccountAddress: "",
		Contracts:      []string{"SOL_PERP"},
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetDriftOpenPerpOrders resp : %s", user)
	return false
}

func callGetDriftOpenMarginOrders(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callGetDriftOpenMarginOrders test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := w.GetDriftOpenMarginOrders(ctx, &pb.GetDriftOpenMarginOrdersRequest{
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

func callGetPerpPositions(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callGetPerpPositions test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := w.GetPerpPositions(ctx, &pb.GetPerpPositionsRequest{
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

func callGetDriftPerpPositions(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callGetDriftPerpPositions test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := w.GetDriftPerpPositions(ctx, &pb.GetDriftPerpPositionsRequest{
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

func callGetUser(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callGetUser test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := w.GetUser(ctx, &pb.GetUserRequest{
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

func callGetDriftUser(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callGetDriftUser test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := w.GetDriftUser(ctx, &pb.GetDriftUserRequest{
		OwnerAddress: ownerAddr,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetDriftUser resp : %s", user)
	return false
}

func callCancelPerpOrder(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callCancelPerpOrder test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.PostCancelPerpOrder(ctx, &pb.PostCancelPerpOrderRequest{
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

func callDriftCancelPerpOrder(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callDriftCancelPerpOrder test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.PostDriftCancelPerpOrder(ctx, &pb.PostDriftCancelPerpOrderRequest{
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

func callCancelDriftMarginOrder(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callCancelDriftMarginOrder test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.PostCancelDriftMarginOrder(ctx, &pb.PostCancelDriftMarginOrderRequest{
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

func callClosePerpPositions(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callClosePerpPositions test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	sig, err := w.PostClosePerpPositions(ctx, &pb.PostClosePerpPositionsRequest{
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

func callPostCloseDriftPerpPositions(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callPostCloseDriftPerpPositions test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	sig, err := w.PostCloseDriftPerpPositions(ctx, &pb.PostCloseDriftPerpPositionsRequest{
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

func callCreateUser(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callCreateUser test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.PostCreateUser(ctx, &pb.PostCreateUserRequest{
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

func callPostCreateDriftUser(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callPostCreateDriftUser test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.PostCreateDriftUser(ctx, &pb.PostCreateDriftUserRequest{
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

func callPostPerpOrder(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callPostPerpOrder test")

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
	sig, err := w.PostPerpOrder(ctx, request)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostPerpOrder signature : %s", sig)
	return false
}

func callPostDriftPerpOrder(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callPostDriftPerpOrder test")

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
	sig, err := w.PostDriftPerpOrder(ctx, request)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostPerpOrder signature : %s", sig)
	return false
}

func callPostModifyOrder(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callPostModifyOrder test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	request := &pb.PostModifyDriftOrderRequest{
		OwnerAddress:    ownerAddr,
		AccountAddress:  "",
		NewLimitPrice:   1000,
		NewPositionSide: "long",
		OrderID:         1,
	}
	sig, err := w.PostModifyDriftOrder(ctx, request)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostModifyOrder signature : %s", sig)
	return false
}

func callPostMarginOrder(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callPostMarginOrder test")

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
	sig, err := w.PostDriftMarginOrder(ctx, request)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostMarginOrder signature : %s", sig)
	return false
}

func callManageCollateralWithdraw(w *provider.WSClient) bool {
	log.Info("starting callManageCollateralWithdraw test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.PostManageCollateral(ctx, &pb.PostManageCollateralRequest{
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

func callPostDriftManageCollateralWithdraw(w *provider.WSClient) bool {
	log.Info("starting callPostDriftManageCollateralWithdraw test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.PostDriftManageCollateral(ctx, &pb.PostDriftManageCollateralRequest{
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

func callManageCollateralTransfer(w *provider.WSClient) bool {
	log.Info("starting callManageCollateralTransfer test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.PostManageCollateral(ctx, &pb.PostManageCollateralRequest{
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
	log.Infof("callManageCollateral signature : %s", sig)
	return false
}

func callPostDriftManageCollateralTransfer(w *provider.WSClient) bool {
	log.Info("starting callPostDriftManageCollateralTransfer test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.PostDriftManageCollateral(ctx, &pb.PostDriftManageCollateralRequest{
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

func callDriftEnableMarginTrading(w *provider.WSClient, ownerAddress string) bool {
	log.Info("starting callDriftEnableMarginTrading transfer test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.PostDriftEnableMarginTrading(ctx, &pb.PostDriftEnableMarginTradingRequest{
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

func callManageCollateralDeposit(w *provider.WSClient) bool {
	log.Info("starting callManageCollateralDeposit deposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.PostManageCollateral(ctx, &pb.PostManageCollateralRequest{
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

func callPostDriftManageCollateralDeposit(w *provider.WSClient) bool {
	log.Info("starting callPostDriftManageCollateralDeposit deposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.PostDriftManageCollateral(ctx, &pb.PostDriftManageCollateralRequest{
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

func callGetOpenPerpOrder(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callGetOpenPerpOrder test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := w.GetOpenPerpOrder(ctx, &pb.GetOpenPerpOrderRequest{
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

func callGetDriftOpenPerpOrder(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callGetDriftOpenPerpOrder test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := w.GetDriftOpenPerpOrder(ctx, &pb.GetDriftOpenPerpOrderRequest{
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

func callPostSettlePNL(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callPostSettlePNL deposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.PostSettlePNL(ctx, &pb.PostSettlePNLRequest{
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

func callGetDriftOpenMarginOrder(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callGetDriftOpenMarginOrder deposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.GetDriftOpenMarginOrder(ctx, &pb.GetDriftOpenMarginOrderRequest{
		OwnerAddress:   ownerAddr,
		AccountAddress: "",
		ClientOrderID:  13,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetDriftOpenMarginOrder signature : %s", sig)
	return false
}

func callPostDriftSettlePNL(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callPostDriftSettlePNL deposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.PostDriftSettlePNL(ctx, &pb.PostDriftSettlePNLRequest{
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

func callPostSettlePNLs(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callPostSettlePNLs test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.PostSettlePNLs(ctx, &pb.PostSettlePNLsRequest{
		Project:                 pb.Project_P_DRIFT,
		OwnerAddress:            ownerAddr,
		SettleeAccountAddresses: []string{"9UnwdvTf5EfGeLyLrF4GZDUs7LKRUeJQzW7qsDVGQ8sS"},
		Contract:                common.PerpContract_SOL_PERP,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostSettlePNLs signature : %s", sig)
	return false
}

func callPostDriftSettlePNLs(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callPostDriftSettlePNLs test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.PostDriftSettlePNLs(ctx, &pb.PostDriftSettlePNLsRequest{
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

func callGetAssets(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callGetAssets test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := w.GetAssets(ctx, &pb.GetAssetsRequest{
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

func callGetDriftAssets(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callGetDriftAssets test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := w.GetDriftAssets(ctx, &pb.GetDriftAssetsRequest{
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

func callGetPerpContracts(w *provider.WSClient) bool {
	log.Info("starting callGetPerpContracts test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := w.GetPerpContracts(ctx, &pb.GetPerpContractsRequest{
		Project: pb.Project_P_DRIFT,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetPerpContracts resp : %s", user)
	return false
}

func callGetDriftPerpContracts(w *provider.WSClient) bool {
	log.Info("starting callGetDriftPerpContracts test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := w.GetDriftPerpContracts(ctx, &pb.GetDriftPerpContractsRequest{})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetDriftPerpContracts resp : %s", user)
	return false
}

func callGetDriftMarkets(w *provider.WSClient) bool {
	log.Info("starting callGetDriftMarkets test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := w.GetDriftMarkets(ctx, &pb.GetDriftMarketsRequest{})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("GetDriftMarkets resp : %s", user)
	return false
}

func callPostLiquidatePerp(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callPostLiquidatePerp deposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.PostLiquidatePerp(ctx, &pb.PostLiquidatePerpRequest{
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

func callPostLiquidateDriftPerp(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callPostLiquidateDriftPerp deposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.PostLiquidateDriftPerp(ctx, &pb.PostLiquidateDriftPerpRequest{
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

func callGetDriftPerpOrderbook(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callGetDriftPerpOrderbook deposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.GetDriftPerpOrderbook(ctx, &pb.GetDriftPerpOrderbookRequest{
		Contract: "SOL_PERP",
		Limit:    12,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetDriftPerpOrderbook signature : %s", sig)
	return false
}
