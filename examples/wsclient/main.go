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
	failed = failed || logCall("callPriceWS", func() bool { return callPriceWS(w) })
	failed = failed || logCall("callOpenOrdersWS", func() bool { return callOpenOrdersWS(w) })
	failed = failed || logCall("callTickersWS", func() bool { return callTickersWS(w) })
	failed = failed || logCall("callUnsettledWS", func() bool { return callUnsettledWS(w) })
	failed = failed || logCall("callAccountBalanceWS", func() bool { return callAccountBalanceWS(w) })
	failed = failed || logCall("callGetQuotes", func() bool { return callGetQuotes(w) })
	failed = failed || logCall("callDriftOrderbookWS", func() bool { return callDriftOrderbookWS(w) })

	// streaming methods
	failed = failed || logCall("callOrderbookWSStream", func() bool { return callOrderbookWSStream(w) })
	failed = failed || logCall("callMarketDepthWSStream", func() bool { return callMarketDepthWSStream(w) })
	failed = failed || logCall("callRecentBlockHashWSStream", func() bool { return callRecentBlockHashWSStream(w) })
	failed = failed || logCall("callPoolReservesWSStream", func() bool { return callPoolReservesWSStream(w) })
	failed = failed || logCall("callBlockWSStream", func() bool { return callBlockWSStream(w) })
	failed = failed || logCall("callDriftOrderbookWSStream", func() bool { return callDriftOrderbookWSStream(w) })
	failed = failed || logCall("callDriftGetPerpTradesStream", func() bool { return callDriftGetPerpTradesStream(w) })

	if cfg.RunSlowStream {
		failed = failed || logCall("callPricesWSStream", func() bool { return callPricesWSStream(w) })
		failed = failed || logCall("callSwapsWSStream", func() bool { return callSwapsWSStream(w) })
		failed = failed || logCall("callTradesWSStream", func() bool { return callTradesWSStream(w) })
	}

	if !cfg.RunTrades {
		log.Info("skipping trades due to config")
		return failed
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

	failed = failed || logCall("orderLifecycleTest", func() bool { return orderLifecycleTest(w, ownerAddr, payerAddr, ooAddr) })
	failed = failed || logCall("cancelAll", func() bool { return cancelAll(w, ownerAddr, payerAddr, ooAddr) })
	failed = failed || logCall("callReplaceByClientOrderID", func() bool { return callReplaceByClientOrderID(w, ownerAddr, payerAddr, ooAddr) })
	failed = failed || logCall("callReplaceOrder", func() bool { return callReplaceOrder(w, ownerAddr, payerAddr, ooAddr) })
	failed = failed || logCall("callRecentBlockHashWSStream", func() bool { return callRecentBlockHashWSStream(w) })
	failed = failed || logCall("callTradeSwap", func() bool { return callTradeSwap(w, ownerAddr) })
	failed = failed || logCall("callRouteTradeSwap", func() bool { return callRouteTradeSwap(w, ownerAddr) })

	failed = failed || logCall("callGetOpenPerpOrders", func() bool { return callGetOpenPerpOrders(w, ownerAddr) })
	failed = failed || logCall("callGetPerpPositions", func() bool { return callGetPerpPositions(w, ownerAddr) })
	failed = failed || logCall("callGetUser", func() bool { return callGetUser(w, ownerAddr) })

	failed = failed || logCall("callGetOpenPerpOrder", func() bool { return callGetOpenPerpOrder(w, ownerAddr) })
	failed = failed || logCall("callGetAssets", func() bool { return callGetAssets(w, ownerAddr) })
	failed = failed || logCall("callGetContracts", func() bool { return callGetContracts(w) })

	if cfg.RunPerpTrades {
		failed = failed || logCall("callCancelPerpOrder", func() bool { return callCancelPerpOrder(w, ownerAddr) })
		failed = failed || logCall("callClosePerpPositions", func() bool { return callClosePerpPositions(w, ownerAddr) })
		failed = failed || logCall("callCreateUser", func() bool { return callCreateUser(w, ownerAddr) })
		failed = failed || logCall("callManageCollateralDeposit", func() bool { return callManageCollateralDeposit(w, ownerAddr) })
		failed = failed || logCall("callPostPerpOrder", func() bool { return callPostPerpOrder(w, ownerAddr) })
		failed = failed || logCall("callManageCollateralWithdraw", func() bool { return callManageCollateralWithdraw(w) })
		failed = failed || logCall("callManageCollateralTransfer", func() bool { return callManageCollateralTransfer(w) })

		failed = failed || logCall("callPostSettlePNL", func() bool { return callPostSettlePNL(w, ownerAddr) })
		failed = failed || logCall("callPostSettlePNLs", func() bool { return callPostSettlePNLs(w, ownerAddr) })
		failed = failed || logCall("callPostLiquidatePerp", func() bool { return callPostLiquidatePerp(w, ownerAddr) })
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

	markets, err := w.GetMarkets(context.Background())
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

	orderbook, err := w.GetOrderbook(context.Background(), "SOL-USDT", 0, pb.Project_P_OPENBOOK)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL-USDT: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()

	orderbook, err = w.GetOrderbook(context.Background(), "SOLUSDT", 2, pb.Project_P_OPENBOOK)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL-USDT: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()

	orderbook, err = w.GetOrderbook(context.Background(), "SOL:USDC", 3, pb.Project_P_OPENBOOK)
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

	mktDepth, err := w.GetMarketDepth(context.Background(), "SOL:USDC", 3, pb.Project_P_OPENBOOK)
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

func callPriceWS(w *provider.WSClient) bool {
	log.Info("fetching prices...")

	pools, err := w.GetPrice(context.Background(), []string{"SOL", "ETH"})
	if err != nil {
		log.Errorf("error with GetPrice request for SOL and ETH: %v", err)
		return true
	} else {
		log.Info(pools)
	}

	fmt.Println()
	return false
}

func callOpenOrdersWS(w *provider.WSClient) bool {
	log.Info("fetching open orders...")

	orders, err := w.GetOpenOrders(context.Background(), "SOLUSDC", "FFqDwRq8B4hhFKRqx7N1M6Dg6vU699hVqeynDeYJdPj5", "", pb.Project_P_OPENBOOK)
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

	response, err := w.GetUnsettled(context.Background(), "SOLUSDC", "AFT8VayE7qr8MoQsW3wHsDS83HhEvhGWdbNSHRKeUDfQ", pb.Project_P_OPENBOOK)
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

	tickers, err := w.GetTickers(context.Background(), "SOLUSDC", pb.Project_P_OPENBOOK)
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
	outToken := "USDC"
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
	marketAddr = "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT"

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

	// create order without actually submitting
	response, err := w.PostOrder(context.Background(), ownerAddr, payerAddr, marketAddr, orderSide, []common.OrderType{orderType}, orderAmount, orderPrice, pb.Project_P_OPENBOOK, opts)
	if err != nil {
		log.Errorf("failed to create order (%v)", err)
		return 0, true
	}
	log.Infof("created unsigned place order transaction: %v", response.Transaction)

	// sign/submit transaction after creation
	sig, err := w.SubmitOrder(context.Background(), ownerAddr, payerAddr, marketAddr,
		orderSide, []common.OrderType{orderType}, orderAmount,
		orderPrice, pb.Project_P_OPENBOOK, opts)
	if err != nil {
		log.Errorf("failed to submit order (%v)", err)
		return 0, true
	}

	log.Infof("placed order %v with clientOrderID %v", sig, clientOrderID)

	return clientOrderID, false
}

func callCancelByClientOrderIDWS(w *provider.WSClient, ownerAddr, ooAddr string, clientOrderID uint64) bool {
	log.Info("trying to cancel order")

	_, err := w.SubmitCancelByClientOrderID(context.Background(), clientOrderID, ownerAddr,
		marketAddr, ooAddr, pb.Project_P_OPENBOOK, true)
	if err != nil {
		log.Errorf("failed to cancel order by client ID (%v)", err)
		return true
	}

	log.Infof("canceled order for clientOrderID %v", clientOrderID)
	return false
}

func callPostSettleWS(w *provider.WSClient, ownerAddr, ooAddr string) bool {
	log.Info("starting post settle")

	sig, err := w.SubmitSettle(context.Background(), ownerAddr, "SOL/USDC", "F75gCEckFAyeeCWA9FQMkmLCmke7ehvBnZeVZ3QgvJR7", "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv", ooAddr, pb.Project_P_OPENBOOK, false)
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
	sig, err := w.SubmitOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []common.OrderType{orderType}, orderAmount, orderPrice, pb.Project_P_OPENBOOK, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #1, signature %s", sig)

	opts.ClientOrderID = clientOrderID2
	sig, err = w.SubmitOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []common.OrderType{orderType}, orderAmount, orderPrice, pb.Project_P_OPENBOOK, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check orders are there
	orders, err := w.GetOpenOrders(ctx, marketAddr, ownerAddr, "", pb.Project_P_OPENBOOK)
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
	sigs, err := w.SubmitCancelAll(ctx, marketAddr, ownerAddr, []string{ooAddr}, pb.Project_P_OPENBOOK, provider.SubmitOpts{
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

	orders, err = w.GetOpenOrders(ctx, marketAddr, ownerAddr, "", pb.Project_P_OPENBOOK)
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
	sig, err := w.SubmitOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []common.OrderType{orderType}, orderAmount, orderPrice, pb.Project_P_OPENBOOK, opts)
	if err != nil {
		log.Error(err)
		return true
	} else {
		log.Infof("submitting place order #1, signature %s", sig)
	}
	time.Sleep(time.Minute)
	// Check order is there
	orders, err := w.GetOpenOrders(ctx, marketAddr, ownerAddr, "", pb.Project_P_OPENBOOK)
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
	sig, err = w.SubmitReplaceByClientOrderID(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []common.OrderType{orderType}, orderAmount, orderPrice/2, pb.Project_P_OPENBOOK, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check order #2 is in orderbook
	orders, err = w.GetOpenOrders(ctx, marketAddr, ownerAddr, "", pb.Project_P_OPENBOOK)
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
	sigs, err := w.SubmitCancelAll(ctx, marketAddr, ownerAddr, []string{ooAddr}, pb.Project_P_OPENBOOK, provider.SubmitOpts{
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
	sig, err := w.SubmitOrder(ctx, ownerAddr, payerAddr, marketAddr, orderSide, []common.OrderType{orderType}, orderAmount, orderPrice, pb.Project_P_OPENBOOK, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #1, signature %s", sig)

	// Check orders are there
	orders, err := w.GetOpenOrders(ctx, marketAddr, ownerAddr, "", pb.Project_P_OPENBOOK)
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
	sig, err = w.SubmitReplaceOrder(ctx, found1.OrderID, ownerAddr, payerAddr, marketAddr, orderSide, []common.OrderType{orderType}, orderAmount, orderPrice/2, pb.Project_P_OPENBOOK, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check orders are there
	orders, err = w.GetOpenOrders(ctx, marketAddr, ownerAddr, "", pb.Project_P_OPENBOOK)
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
	sigs, err := w.SubmitCancelAll(ctx, marketAddr, ownerAddr, []string{ooAddr}, pb.Project_P_OPENBOOK, provider.SubmitOpts{
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

func callTradeSwap(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting trade swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("trade swap")
	sig, err := w.SubmitTradeSwap(ctx, ownerAddr, "USDC",
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

func callRouteTradeSwap(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting route trade swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("route trade swap")
	sig, err := w.SubmitRouteTradeSwap(ctx, &pb.RouteTradeSwapRequest{
		OwnerAddress: ownerAddr,
		Project:      pb.Project_P_RAYDIUM,
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
				OutToken:     "USDC",
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

func callDriftOrderbookWS(w *provider.WSClient) bool {
	log.Info("fetching drift orderbooks...")

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

func callDriftOrderbookWSStream(w *provider.WSClient) bool {
	log.Info("starting drift orderbook stream")

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

func callCancelPerpOrder(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callCancelPerpOrder test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.SubmitCancelPerpOrder(ctx, &pb.PostCancelPerpOrderRequest{
		Project:       pb.Project_P_DRIFT,
		OwnerAddress:  ownerAddr,
		OrderID:       1,
		ClientOrderID: 0,
		Contract:      common.PerpContract_SOL_PERP,
	}, false)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callCancelPerpOrder signature : %s", sig)
	return false
}

func callClosePerpPositions(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callClosePerpPositions test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	sig, err := w.SubmitClosePerpPositions(ctx, &pb.PostClosePerpPositionsRequest{
		Project:      pb.Project_P_DRIFT,
		OwnerAddress: ownerAddr,
		Contracts:    []common.PerpContract{common.PerpContract_SOL_PERP},
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  true,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callClosePerpPositions signature : %s", sig)
	return false
}

func callCreateUser(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callCreateUser test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.SubmitCreateUser(ctx, &pb.PostCreateUserRequest{
		Project:      pb.Project_P_DRIFT,
		OwnerAddress: ownerAddr,
	}, false)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callCreateUser signature : %s", sig)
	return false
}

func callPostPerpOrder(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callPostPerpOrder test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	request := &pb.PostPerpOrderRequest{
		Project:        pb.Project_P_DRIFT,
		OwnerAddress:   ownerAddr,
		PayerAddress:   ownerAddr,
		Contract:       common.PerpContract_SOL_PERP,
		AccountAddress: "",
		PositionSide:   common.PositionSide_PS_SHORT,
		Slippage:       10,
		Type:           common.DriftOrderType_POT_LIMIT,
		Amount:         1,
		Price:          1000,
		ClientOrderID:  2,
	}
	sig, err := w.SubmitPostPerpOrder(ctx, request, false)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostPerpOrder signature : %s", sig)
	return false
}

func callManageCollateralWithdraw(w *provider.WSClient) bool {
	log.Info("starting callManageCollateralWithdraw withdraw test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.SubmitManageCollateral(ctx, &pb.PostManageCollateralRequest{
		Project:        pb.Project_P_DRIFT,
		Amount:         1,
		AccountAddress: "61bvX2qCwzPKNztgVQF3ktDHM2hZGdivCE28RrC99EAS",
		Type:           common.PerpCollateralType_PCT_WITHDRAWAL,
		Token:          common.PerpCollateralToken_PCTK_SOL,
	}, false)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callManageCollateral signature : %s", sig)
	return false
}

func callManageCollateralTransfer(w *provider.WSClient) bool {
	log.Info("starting callManageCollateralTransfer withdraw test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.SubmitManageCollateral(ctx, &pb.PostManageCollateralRequest{
		Project:          pb.Project_P_DRIFT,
		Amount:           1,
		AccountAddress:   "61bvX2qCwzPKNztgVQF3ktDHM2hZGdivCE28RrC99EAS",
		Type:             common.PerpCollateralType_PCT_TRANSFER,
		Token:            common.PerpCollateralToken_PCTK_SOL,
		ToAccountAddress: "BTHDMaruPPTyUAZDv6w11qSMtyNAaNX6zFTPPepY863V",
	}, false)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callManageCollateral signature : %s", sig)
	return false
}

func callManageCollateralDeposit(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callManageCollateralDeposit deposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.SubmitManageCollateral(ctx, &pb.PostManageCollateralRequest{
		Project:        pb.Project_P_DRIFT,
		Amount:         1,
		AccountAddress: "61bvX2qCwzPKNztgVQF3ktDHM2hZGdivCE28RrC99EAS",
		Type:           common.PerpCollateralType_PCT_DEPOSIT,
		Token:          common.PerpCollateralToken_PCTK_SOL,
	}, false)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callManageCollateral signature : %s", sig)
	return false
}

func callGetOpenPerpOrder(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callGetOpenPerpOrder test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := w.GetOpenPerpOrder(ctx, &pb.GetOpenPerpOrderRequest{
		OwnerAddress:   ownerAddr,
		AccountAddress: "",
		Contract:       common.PerpContract_SOL_PERP,
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

func callPostSettlePNL(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting SubmitPostSettlePNL deposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.SubmitPostSettlePNL(ctx, &pb.PostSettlePNLRequest{
		Project:               pb.Project_P_DRIFT,
		OwnerAddress:          ownerAddr,
		SettleeAccountAddress: ownerAddr,
		Contract:              common.PerpContract_SOL_PERP,
	}, false)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("SubmitPostSettlePNL signature : %s", sig)
	return false
}

func callPostSettlePNLs(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callPostSettlePNLs deposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.SubmitPostSettlePNLs(ctx, &pb.PostSettlePNLsRequest{
		Project:                 pb.Project_P_DRIFT,
		OwnerAddress:            ownerAddr,
		SettleeAccountAddresses: []string{ownerAddr},
		Contract:                common.PerpContract_SOL_PERP,
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  true,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostSettlePNLs signature : %s", sig)
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

func callGetContracts(w *provider.WSClient) bool {
	log.Info("starting callGetContracts test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := w.GetPerpContracts(ctx, &pb.GetPerpContractsRequest{
		Project: pb.Project_P_DRIFT,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetAssets resp : %s", user)
	return false
}

func callPostLiquidatePerp(w *provider.WSClient, ownerAddr string) bool {
	log.Info("starting callPostLiquidatePerp deposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := w.SubmitPostLiquidatePerp(ctx, &pb.PostLiquidatePerpRequest{
		Project:               pb.Project_P_DRIFT,
		OwnerAddress:          ownerAddr,
		Amount:                1,
		Contract:              common.PerpContract_SOL_PERP,
		SettleeAccountAddress: ownerAddr,
	}, false)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callPostLiquidatePerp signature : %s", sig)
	return false
}
