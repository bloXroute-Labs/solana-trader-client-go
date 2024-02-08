package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/bloXroute-Labs/solana-trader-client-go/examples/config"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	"github.com/bloXroute-Labs/solana-trader-proto/common"

	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	sideAsk      = "ask"
	typeLimit    = "limit"
	computeLimit = 10000
	computePrice = 2000
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

	// streaming methods
	failed = failed || logCall("callSwapsWSStream", func() bool { return callSwapsWSStream(w) })
	failed = failed || logCall("callOrderbookWSStream", func() bool { return callOrderbookWSStream(w) })
	failed = failed || logCall("callMarketDepthWSStream", func() bool { return callMarketDepthWSStream(w) })
	failed = failed || logCall("callRecentBlockHashWSStream", func() bool { return callRecentBlockHashWSStream(w) })
	failed = failed || logCall("callPoolReservesWSStream", func() bool { return callPoolReservesWSStream(w) })
	failed = failed || logCall("callBlockWSStream", func() bool { return callBlockWSStream(w) })
	failed = failed || logCall("callBlockWSStream", func() bool { return callBlockWSStream(w) })

	if cfg.RunSlowStream {
		failed = failed || logCall("callPricesWSStream", func() bool { return callPricesWSStream(w) })
		failed = failed || logCall("callSwapsWSStream", func() bool { return callSwapsWSStream(w) })
		failed = failed || logCall("callTradesWSStream", func() bool { return callTradesWSStream(w) })
		failed = failed || logCall("callGetNewRaydiumPoolsStream", func() bool { return callGetNewRaydiumPoolsStream(w) })
	}

	// informational requests
	failed = failed || logCall("callMarketsWS", func() bool { return callMarketsWS(w) })
	failed = failed || logCall("callOrderbookWS", func() bool { return callOrderbookWS(w) })
	failed = failed || logCall("callMarketDepthWS", func() bool { return callMarketDepthWS(w) })
	failed = failed || logCall("callTradesWS", func() bool { return callTradesWS(w) })
	failed = failed || logCall("callPoolsWS", func() bool { return callPoolsWS(w) })
	failed = failed || logCall("callRaydiumPools", func() bool { return callRaydiumPoolsWS(w) })
	failed = failed || logCall("callGetTransactionWS", func() bool { return callGetTransactionWS(w) })
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

	// streaming methods
	failed = failed || logCall("callOrderbookWSStream", func() bool { return callOrderbookWSStream(w) })
	failed = failed || logCall("callMarketDepthWSStream", func() bool { return callMarketDepthWSStream(w) })
	failed = failed || logCall("callRecentBlockHashWSStream", func() bool { return callRecentBlockHashWSStream(w) })
	failed = failed || logCall("callPoolReservesWSStream", func() bool { return callPoolReservesWSStream(w) })
	failed = failed || logCall("callBlockWSStream", func() bool { return callBlockWSStream(w) })

	if cfg.RunSlowStream {
		failed = failed || logCall("callPricesWSStream", func() bool { return callPricesWSStream(w) })
		failed = failed || logCall("callSwapsWSStream", func() bool { return callSwapsWSStream(w) })
		failed = failed || logCall("callTradesWSStream", func() bool { return callTradesWSStream(w) })
		failed = failed || logCall("callGetNewRaydiumPoolsStream", func() bool { return callGetNewRaydiumPoolsStream(w) })
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
		failed = failed || logCall("callReplaceOrder", func() bool { return callReplaceOrder(w, ownerAddr, payerAddr, ooAddr, sideAsk, typeLimit) })
		failed = failed || logCall("callRecentBlockHashWSStream", func() bool { return callRecentBlockHashWSStream(w) })
		failed = failed || logCall("callTradeSwap", func() bool { return callTradeSwap(w, ownerAddr) })
		failed = failed || logCall("callTradeSwapWithPriorityFee", func() bool { return callTradeSwapWithPriorityFee(w, ownerAddr, computeLimit, computePrice) })
		failed = failed || logCall("callRouteTradeSwap", func() bool { return callRouteTradeSwap(w, ownerAddr) })
		failed = failed || logCall("callRaydiumTradeSwap", func() bool { return callRaydiumSwap(w, ownerAddr) })
		failed = failed || logCall("callJupiterTradeSwap", func() bool { return callJupiterSwap(w, ownerAddr) })
		failed = failed || logCall("callRaydiumRouteTradeSwap", func() bool { return callRaydiumRouteSwap(w, ownerAddr) })
		failed = failed || logCall("callJupiterRouteTradeSwap", func() bool { return callJupiterRouteSwap(w, ownerAddr) })
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

func callGetTransactionWS(w *provider.WSClient) bool {
	log.Info("calling GetTransaction...")

	tx, err := w.GetTransaction(context.Background(), &pb.GetTransactionRequest{
		Signature: "2s48MnhH54GfJbRwwiEK7iWKoEh3uNbS2zDEVBPNu7DaCjPXe3bfqo6RuCg9NgHRFDn3L28sMVfEh65xevf4o5W3",
	})
	if err != nil {
		log.Errorf("error with GetTransaction request: %v", err)
		return true
	} else {
		log.Info(tx)
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

	pools, err := w.GetPrice(context.Background(), []string{"So11111111111111111111111111111111111111112", "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263"})
	if err != nil {
		log.Errorf("error with GetPrice request for SOL and BONK: %v", err)
		return true
	} else {
		log.Info(pools)
	}

	return false
}

func callRaydiumPricesWS(w *provider.WSClient) bool {
	log.Info("fetching Raydium prices...")

	pools, err := w.GetRaydiumPrices(context.Background(), &pb.GetRaydiumPricesRequest{
		Tokens: []string{"So11111111111111111111111111111111111111112", "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263"},
	})
	if err != nil {
		log.Errorf("error with GetRaydiumPrices request for SOL and BONK: %v", err)
		return true
	} else {
		log.Info(pools)
	}

	return false
}

func callJupiterPricesWS(w *provider.WSClient) bool {
	log.Info("fetching Jupiter prices...")

	pools, err := w.GetJupiterPrices(context.Background(), &pb.GetJupiterPricesRequest{
		Tokens: []string{"So11111111111111111111111111111111111111112", "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263"},
	})
	if err != nil {
		log.Errorf("error with GetJupiterPrices request for SOL and BONK: %v", err)
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

	inToken := "So11111111111111111111111111111111111111112"
	outToken := "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
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

	inToken := "So11111111111111111111111111111111111111112"
	outToken := "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
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

	inToken := "So11111111111111111111111111111111111111112"
	outToken := "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
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

	if len(quotes.Routes) == 0 {
		log.Errorf("did not get any quotes, got %v quotes", len(quotes.Routes))
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

func callGetNewRaydiumPoolsStream(w *provider.WSClient) bool {
	log.Info("starting get new raydium pools stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	poolsChan := make(chan *pb.GetNewRaydiumPoolsResponse)

	stream, err := w.GetNewRaydiumPoolsStream(ctx)
	if err != nil {
		log.Errorf("error with GetTradesStream request for SOL/USDC: %v", err)
		return true
	}

	stream.Into(poolsChan)
	for i := 1; i <= 1; i++ {
		_, ok := <-poolsChan
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

	clientOrderID, fail := callPlaceOrderWS(w, ownerAddr, payerAddr, ooAddr, sideAsk, typeLimit)
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

	fail = callCancelByClientOrderIDWS(w, ownerAddr, ooAddr, clientOrderID, sideAsk)
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

func callPlaceOrderWS(w *provider.WSClient, ownerAddr, payerAddr, ooAddr string, orderSide string, orderType string) (uint64, bool) {
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
		orderSide, orderType, orderAmount, orderPrice, opts)
	if err != nil {
		log.Errorf("failed to submit order (%v)", err)
		return 0, true
	}

	log.Infof("placed order %v with clientOrderID %v", sig, clientOrderID)

	return clientOrderID, false
}

func callCancelByClientOrderIDWS(w *provider.WSClient, ownerAddr, ooAddr string, clientOrderID uint64, orderSide string) bool {
	log.Info("trying to cancel order")

	_, err := w.SubmitCancelOrderV2(context.Background(), &pb.PostCancelOrderRequestV2{
		OrderID:           "",
		Side:              orderSide,
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

func cancelAll(w *provider.WSClient, ownerAddr, payerAddr, ooAddr string, orderSide string, orderType string) bool {
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
		SkipPreFlight:     config.BoolPtr(true),
	}

	// Place 2 orders in orderbook
	log.Info("placing orders")
	sig, err := w.SubmitOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderType, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #1, signature %s", sig)

	opts.ClientOrderID = clientOrderID2
	sig, err = w.SubmitOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderType, orderAmount, orderPrice, opts)
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

func callReplaceByClientOrderID(w *provider.WSClient, ownerAddr, payerAddr, ooAddr string, orderSide string, orderType string) bool {
	log.Info("starting replace by client order ID test")
	fmt.Println()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rand.Seed(time.Now().UnixNano())
	clientOrderID1 := rand.Uint64()
	opts := provider.PostOrderOpts{
		ClientOrderID:     clientOrderID1,
		OpenOrdersAddress: ooAddr,
		SkipPreFlight:     config.BoolPtr(true),
	}

	// Place order in orderbook
	log.Info("placing order")
	sig, err := w.SubmitOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderType, orderAmount, orderPrice, opts)
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
	sig, err = w.SubmitReplaceOrderV2(ctx, "", ownerAddr, payerAddr, marketAddr, orderSide, orderType, orderAmount, orderPrice/2, opts)
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
		Side:              orderSide,
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

func callReplaceOrder(w *provider.WSClient, ownerAddr, payerAddr, ooAddr string, orderSide string, orderType string) bool {
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
		SkipPreFlight:     config.BoolPtr(true),
	}

	// Place order in orderbook
	log.Info("placing order")
	sig, err := w.SubmitOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderType, orderAmount, orderPrice, opts)
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
	sig, err = w.SubmitReplaceOrderV2(ctx, found1.OrderID, ownerAddr, payerAddr, marketAddr, orderSide, orderType, orderAmount, orderPrice/2, opts)
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
		Side:              orderSide,
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
	sig, err := w.SubmitTradeSwap(ctx, ownerAddr, "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		"So11111111111111111111111111111111111111112", 0.01, 0.1, "raydium", provider.SubmitOpts{
			SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
			SkipPreFlight:  config.BoolPtr(false),
		})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("trade swap transaction signature : %s", sig)
	return false
}

func callTradeSwapWithPriorityFee(w *provider.WSClient, ownerAddr string, computeLimit uint32, computePrice uint64) bool {
	log.Info("starting trade swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("trade swap")
	sig, err := w.SubmitTradeSwapWithPriorityFee(ctx, ownerAddr, "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		"So11111111111111111111111111111111111111112", 0.01, 0.1, "raydium", computeLimit, computePrice,
		provider.SubmitOpts{
			SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
			SkipPreFlight:  config.BoolPtr(false),
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
		InToken:      "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		OutToken:     "So11111111111111111111111111111111111111112",
		Slippage:     0.1,
		InAmount:     0.01,
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  config.BoolPtr(false),
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
				OutToken:     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
				InAmount:     0.007505,
				OutAmount:    0.004043,
				OutAmountMin: 0.004000,
			},
		},
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  config.BoolPtr(false),
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
				OutToken:     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
				InAmount:     0.007505,
				OutAmount:    0.004043,
				OutAmountMin: 0.004000,
			},
		},
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  config.BoolPtr(false),
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
		InToken:      "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		OutToken:     "So11111111111111111111111111111111111111112",
		Slippage:     0.1,
		InAmount:     0.01,
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  config.BoolPtr(false),
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
		Slippage:     0.25,
		Steps: []*pb.JupiterRouteStep{
			{
				Project: &pb.StepProject{
					Label: "Raydium",
					Id:    "61acRgpURKTU8LKPJKs6WQa18KzD9ogavXzjxfD84KLu",
				},
				InToken:      "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
				OutToken:     "So11111111111111111111111111111111111111112",
				InAmount:     0.01,
				OutAmountMin: 0.000123117,
				OutAmount:    0.000123425,
				Fee: &common.Fee{
					Amount:  0.000025,
					Mint:    "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
					Percent: 0.0025062656,
				},
			},
		},
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  config.BoolPtr(false),
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

	stream, err := w.GetPricesStream(ctx, []pb.Project{pb.Project_P_RAYDIUM}, []string{"So11111111111111111111111111111111111111112"})
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
