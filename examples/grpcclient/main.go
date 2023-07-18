package main

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/examples/config"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	"github.com/bloXroute-Labs/solana-trader-proto/common"
	"google.golang.org/protobuf/encoding/protojson"
	"math/rand"
	"os"
	"time"

	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	log "github.com/sirupsen/logrus"
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

	failed = failed || logCall("callMarketsGRPC", func() bool { return callMarketsGRPC(g) })
	failed = failed || logCall("callOrderbookGRPC", func() bool { return callOrderbookGRPC(g) })
	failed = failed || logCall("callMarketDepthGRPC", func() bool { return callMarketDepthGRPC(g) })
	failed = failed || logCall("callOpenOrdersGRPC", func() bool { return callOpenOrdersGRPC(g) })
	failed = failed || logCall("callTickersGRPC", func() bool { return callTickersGRPC(g) })

	failed = failed || logCall("callPoolsGRPC", func() bool { return callPoolsGRPC(g) })
	failed = failed || logCall("callRaydiumPoolsGRPC", func() bool { return callRaydiumPoolsGRPC(g) })
	failed = failed || logCall("callPriceGRPC", func() bool { return callPriceGRPC(g) })
	failed = failed || logCall("callRaydiumPricesGRPC", func() bool { return callRaydiumPricesGRPC(g) })
	failed = failed || logCall("callJupiterPricesGRPC", func() bool { return callJupiterPricesGRPC(g) })
	failed = failed || logCall("callDriftPerpOrderbookGRPC", func() bool { return callDriftPerpOrderbookGRPC(g) })
	failed = failed || logCall("callDriftGetMarginOrderbookGRPC", func() bool { return callDriftGetMarginOrderbookGRPC(g) })
	failed = failed || logCall("callDriftMarketDepthGRPC", func() bool { return callDriftMarketDepthGRPC(g) })

	if cfg.RunSlowStream {
		failed = failed || logCall("callOrderbookGRPCStream", func() bool { return callOrderbookGRPCStream(g) })
		failed = failed || logCall("callMarketDepthGRPCStream", func() bool { return callMarketDepthGRPCStream(g) })
	}

	if cfg.RunSlowStream {
		failed = failed || logCall("callPricesGRPCStream", func() bool { return callPricesGRPCStream(g) })
		failed = failed || logCall("callTradesGRPCStream", func() bool { return callTradesGRPCStream(g) })
		failed = failed || logCall("callSwapsGRPCStream", func() bool { return callSwapsGRPCStream(g) })
	}

	failed = failed || logCall("callUnsettledGRPC", func() bool { return callUnsettledGRPC(g) })
	failed = failed || logCall("callGetAccountBalanceGRPC", func() bool { return callGetAccountBalanceGRPC(g) })

	failed = failed || logCall("callGetQuotes", func() bool { return callGetQuotes(g) })
	failed = failed || logCall("callGetRaydiumQuotes", func() bool { return callGetRaydiumQuotes(g) })
	failed = failed || logCall("callGetJupiterQuotes", func() bool { return callGetJupiterQuotes(g) })
	failed = failed || logCall("callRecentBlockHashGRPCStream", func() bool { return callRecentBlockHashGRPCStream(g) })
	failed = failed || logCall("callPoolReservesGRPCStream", func() bool { return callPoolReservesGRPCStream(g) })
	failed = failed || logCall("callBlockGRPCStream", func() bool { return callBlockGRPCStream(g) })
	failed = failed || logCall("callDriftPerpOrderbookGRPCStream", func() bool { return callDriftPerpOrderbookGRPCStream(g) })
	failed = failed || logCall("callDriftMarginOrderbooksGRPCStream", func() bool { return callDriftMarginOrderbooksGRPCStream(g) })
	failed = failed || logCall("callDriftGetPerpTradesStream", func() bool { return callDriftGetPerpTradesStream(g) })
	failed = failed || logCall("callDriftGetMarketDepthsStream", func() bool { return callDriftGetMarketDepthsStream(g) })

	// calls below this place an order and immediately cancel it
	// you must specify:
	//	- PRIVATE_KEY (by default loaded during provider.NewGRPCClient()) to sign transactions
	// 	- PUBLIC_KEY to indicate which account you wish to trade from
	//	- OPEN_ORDERS to indicate your Serum account to speed up lookups (optional)
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
	if !ok {
		log.Infof("OPEN_ORDERS environment variable not set: requests will be slower")
	}

	if cfg.RunTrades {
		failed = failed || logCall("orderLifecycleTest", func() bool { return orderLifecycleTest(g, ownerAddr, payerAddr, ooAddr) })
		failed = failed || logCall("cancelAll", func() bool { return cancelAll(g, ownerAddr, payerAddr, ooAddr) })
		failed = failed || logCall("callReplaceByClientOrderID", func() bool { return callReplaceByClientOrderID(g, ownerAddr, payerAddr, ooAddr) })
		failed = failed || logCall("callReplaceOrder", func() bool { return callReplaceOrder(g, ownerAddr, payerAddr, ooAddr) })
		failed = failed || logCall("callTradeSwap", func() bool { return callTradeSwap(g, ownerAddr) })
		failed = failed || logCall("callRouteTradeSwap", func() bool { return callRouteTradeSwap(g, ownerAddr) })
		failed = failed || logCall("callRaydiumTradeSwap", func() bool { return callRaydiumSwap(g, ownerAddr) })
		failed = failed || logCall("callJupiterTradeSwap", func() bool { return callJupiterSwap(g, ownerAddr) })
		failed = failed || logCall("callRaydiumRouteTradeSwap", func() bool { return callRaydiumRouteSwap(g, ownerAddr) })
		failed = failed || logCall("callJupiterRouteTradeSwap", func() bool { return callJupiterRouteSwap(g, ownerAddr) })

	}

	failed = failed || logCall("callGetOpenPerpOrders", func() bool { return callGetOpenPerpOrders(g, ownerAddr) })
	failed = failed || logCall("callGetDriftOpenMarginOrders", func() bool { return callGetDriftOpenMarginOrders(g, ownerAddr) })
	failed = failed || logCall("callGetPerpPositions", func() bool { return callGetPerpPositions(g, ownerAddr) })
	failed = failed || logCall("callGetDriftPerpPositions", func() bool { return callGetDriftPerpPositions(g, ownerAddr) })
	failed = failed || logCall("callGetUser", func() bool { return callGetUser(g, ownerAddr) })

	failed = failed || logCall("callGetOpenPerpOrder", func() bool { return callGetOpenPerpOrder(g, ownerAddr) })
	failed = failed || logCall("callGetDriftOpenPerpOrders", func() bool { return callGetDriftOpenPerpOrders(g, ownerAddr) })
	failed = failed || logCall("callGetAssets", func() bool { return callGetAssets(g, ownerAddr) })
	failed = failed || logCall("callGetPerpContracts", func() bool { return callGetPerpContracts(g) })
	failed = failed || logCall("callGetDriftMarkets", func() bool { return callGetDriftMarkets(g) })

	failed = failed || logCall("callGetDriftAssets", func() bool { return callGetDriftAssets(g, ownerAddr) })
	failed = failed || logCall("callGetDriftPerpContracts", func() bool { return callGetDriftPerpContracts(g) })
	failed = failed || logCall("callGetDriftPerpOrderbook", func() bool { return callGetDriftPerpOrderbook(g, ownerAddr) })
	failed = failed || logCall("callGetDriftUser", func() bool { return callGetDriftUser(g, ownerAddr) })
	failed = failed || logCall("callGetDriftOpenPerpOrder", func() bool { return callGetDriftOpenPerpOrder(g, ownerAddr) })
	failed = failed || logCall("callGetDriftOpenMarginOrder", func() bool { return callGetDriftOpenMarginOrder(g, ownerAddr) })

	if cfg.RunPerpTrades {
		failed = failed || logCall("callCancelPerpOrder", func() bool { return callCancelPerpOrder(g, ownerAddr) })
		failed = failed || logCall("callDriftCancelPerpOrder", func() bool { return callDriftCancelPerpOrder(g, ownerAddr) })

		failed = failed || logCall("callCancelDriftMarginOrder", func() bool { return callCancelDriftMarginOrder(g, ownerAddr) })
		failed = failed || logCall("callClosePerpPositions", func() bool { return callClosePerpPositions(g, ownerAddr) })
		failed = failed || logCall("callCreateUser", func() bool { return callCreateUser(g, ownerAddr) })
		failed = failed || logCall("callManageCollateralDeposit", func() bool { return callManageCollateralDeposit(g) })
		failed = failed || logCall("callPostPerpOrder", func() bool { return callPostPerpOrder(g, ownerAddr) })
		failed = failed || logCall("callPostModifyOrder", func() bool { return callPostModifyOrder(g, ownerAddr) })
		failed = failed || logCall("callPostMarginOrder", func() bool { return callPostMarginOrder(g, ownerAddr) })
		failed = failed || logCall("callManageCollateralWithdraw", func() bool { return callManageCollateralWithdraw(g) })
		failed = failed || logCall("callManageCollateralTransfer", func() bool { return callManageCollateralTransfer(g) })
		failed = failed || logCall("callDriftEnableMarginTrading", func() bool { return callDriftEnableMarginTrading(g, ownerAddr) })
		failed = failed || logCall("callPostSettlePNL", func() bool { return callPostSettlePNL(g, ownerAddr) })
		failed = failed || logCall("callPostSettlePNLs", func() bool { return callPostSettlePNLs(g, ownerAddr) })
		failed = failed || logCall("callPostLiquidatePerp", func() bool { return callPostLiquidatePerp(g, ownerAddr) })

		failed = failed || logCall("callPostCloseDriftPerpPositions", func() bool { return callPostCloseDriftPerpPositions(g, ownerAddr) })
		failed = failed || logCall("callPostCreateDriftUser", func() bool { return callPostCreateDriftUser(g, ownerAddr) })
		failed = failed || logCall("callPostDriftManageCollateralDeposit", func() bool { return callPostDriftManageCollateralDeposit(g) })
		failed = failed || logCall("callPostDriftManageCollateralWithdraw", func() bool { return callPostDriftManageCollateralWithdraw(g) })
		failed = failed || logCall("callPostDriftManageCollateralTransfer", func() bool { return callPostDriftManageCollateralTransfer(g) })
		failed = failed || logCall("callPostDriftSettlePNL", func() bool { return callPostDriftSettlePNL(g, ownerAddr) })
		failed = failed || logCall("callPostDriftSettlePNLs", func() bool { return callPostDriftSettlePNLs(g, ownerAddr) })
		failed = failed || logCall("callPostLiquidateDriftPerp", func() bool { return callPostLiquidateDriftPerp(g, ownerAddr) })
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

func callMarketsGRPC(g *provider.GRPCClient) bool {
	markets, err := g.GetMarketsV2(context.Background())
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
	orderbook, err := g.GetOrderbookV2(context.Background(), "SOL-USDC", 0)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL-USDC: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()

	orderbook, err = g.GetOrderbookV2(context.Background(), "SOLUSDT", 2)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOLUSDT: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()

	orderbook, err = g.GetOrderbookV2(context.Background(), "SOL:USDC", 3)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL:USDC: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()
	return false
}

func callMarketDepthGRPC(g *provider.GRPCClient) bool {
	mktDepth, err := g.GetMarketDepthV2(context.Background(), "SOL-USDC", 0)
	if err != nil {
		log.Errorf("error with GetMarketDepth request for SOL-USDC: %v", err)
		return true
	} else {
		log.Info(mktDepth)
	}

	fmt.Println()
	return false
}

func callOpenOrdersGRPC(g *provider.GRPCClient) bool {
	orders, err := g.GetOpenOrdersV2(context.Background(), "SOLUSDC", "FFqDwRq8B4hhFKRqx7N1M6Dg6vU699hVqeynDeYJdPj5", "", "", 0)
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
	response, err := g.GetUnsettledV2(context.Background(), "SOLUSDC", "HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc")
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
	orders, err := g.GetTickersV2(context.Background(), "SOLUSDC")
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

func callRaydiumPoolsGRPC(g *provider.GRPCClient) bool {
	pools, err := g.GetRaydiumPools(context.Background(), &pb.GetRaydiumPoolsRequest{})
	if err != nil {
		log.Errorf("error with GetRaydiumPools request for Raydium: %v", err)
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

func callRaydiumPricesGRPC(g *provider.GRPCClient) bool {
	prices, err := g.GetRaydiumPrices(context.Background(), &pb.GetRaydiumPricesRequest{
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

func callJupiterPricesGRPC(g *provider.GRPCClient) bool {
	prices, err := g.GetJupiterPrices(context.Background(), &pb.GetJupiterPricesRequest{
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

func callGetQuotes(g *provider.GRPCClient) bool {
	log.Info("starting get quotes test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inToken := "SOL"
	outToken := "USDT"
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

func callGetRaydiumQuotes(g *provider.GRPCClient) bool {
	log.Info("starting get Raydium quotes test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inToken := "SOL"
	outToken := "USDT"
	amount := 0.01
	slippage := float64(5)

	quotes, err := g.GetRaydiumQuotes(ctx, &pb.GetRaydiumQuotesRequest{
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

func callGetJupiterQuotes(g *provider.GRPCClient) bool {
	log.Info("starting get Jupiter quotes test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inToken := "SOL"
	outToken := "USDT"
	amount := 0.01
	slippage := float64(5)
	limit := int32(3)

	quotes, err := g.GetJupiterQuotes(ctx, &pb.GetJupiterQuotesRequest{
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

func callOrderbookGRPCStream(g *provider.GRPCClient) bool {
	log.Info("starting orderbook stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream error response
	stream, err := g.GetOrderbookStream(ctx, []string{"SOL/USDC", "xxx"}, 3, pb.Project_P_OPENBOOK)
	if err != nil {
		log.Errorf("connection could not be established. error: %v", err)
		return true
	}

	_, err = stream()
	if err != nil {
		// demonstration purposes only. will swallow
		log.Infof("subscription error: %v", err)
	} else {
		log.Error("subscription should have returned an error")
		return true
	}

	// Stream ok response
	stream, err = g.GetOrderbookStream(ctx, []string{"SOL/USDC", "SOL-USDT"}, 3, pb.Project_P_OPENBOOK)
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
	for i := 1; i <= 1; i++ {
		data, ok := <-orderbookCh
		if !ok {
			// channel closed
			return true
		}

		log.Infof("response %v received, data %v ", i, data)
	}

	return false
}

func callMarketDepthGRPCStream(g *provider.GRPCClient) bool {
	log.Info("starting market depth stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream error response
	stream, err := g.GetMarketDepthsStream(ctx, []string{"SOL/USDC", "xxx"}, 3, pb.Project_P_OPENBOOK)
	if err != nil {
		log.Errorf("connection could not be established. error: %v", err)
		return true
	}

	_, err = stream()
	if err != nil {
		// demonstration purposes only. will swallow
		log.Infof("subscription error: %v", err)
	} else {
		log.Error("subscription should have returned an error")
		return true
	}

	// Stream ok response
	stream, err = g.GetMarketDepthsStream(ctx, []string{"SOL/USDC", "SOL-USDT"}, 3, pb.Project_P_OPENBOOK)
	if err != nil {
		log.Errorf("connection could not be established. error: %v", err)
		return true
	}

	_, err = stream()
	if err != nil {
		log.Errorf("subscription error: %v", err)
		return true
	}

	ordermktDepthUpdateCh := stream.Channel(0)
	for i := 1; i <= 1; i++ {
		data, ok := <-ordermktDepthUpdateCh
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
	stream, err := g.GetTradesStream(ctx, "SOL/USDC", 3, pb.Project_P_OPENBOOK)
	if err != nil {
		log.Errorf("error with GetTrades stream request for SOL/USDC: %v", err)
		return true
	}
	stream.Into(tradesChan)
	for i := 1; i <= 1; i++ {
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

const (
	// SOL/USDC market
	marketAddr = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh6"

	orderSide   = pb.Side_S_ASK
	orderType   = common.OrderType_OT_LIMIT
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
	stream, err := g.GetOrderStatusStream(ctx, marketAddr, ownerAddr, pb.Project_P_OPENBOOK)
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
		SkipPreFlight:     true,
	}

	// create order without actually submitting
	response, err := g.PostOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderAmount, orderPrice, opts)
	if err != nil {
		log.Errorf("failed to create order (%v)", err)
		return 0, true
	}
	log.Infof("created unsigned place order transaction: %v", response.Transaction)

	// sign/submit transaction after creation
	sig, err := g.SubmitOrderV2(ctx, ownerAddr, ownerAddr, marketAddr,
		orderSide, orderAmount, orderPrice, opts)
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

	sig, err := g.SubmitCancelOrderV2(ctx, "", clientID, pb.Side_S_ASK, ownerAddr,
		marketAddr, ooAddr, provider.SubmitOpts{
			SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
			SkipPreFlight:  true,
		})
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

	sig, err := g.SubmitSettleV2(ctx, ownerAddr, "SOL/USDC", "F75gCEckFAyeeCWA9FQMkmLCmke7ehvBnZeVZ3QgvJR7", "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv", ooAddr, false)
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
	sig, err := g.SubmitOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #1, signature %s", sig)

	opts.ClientOrderID = clientOrderID2
	sig, err = g.SubmitOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check orders are there
	orders, err := g.GetOpenOrdersV2(ctx, marketAddr, ownerAddr, "", "", 0)
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
	sigs, err := g.SubmitCancelOrderV2(ctx, "", 0, pb.Side_S_ASK, ownerAddr, marketAddr, ooAddr, provider.SubmitOpts{
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

	orders, err = g.GetOpenOrdersV2(ctx, marketAddr, ownerAddr, "", "", 0)
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
	sig, err := g.SubmitOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #1, signature %s", sig)
	time.Sleep(time.Minute * 1)

	// Check order is there
	orders, err := g.GetOpenOrdersV2(ctx, marketAddr, ownerAddr, "", "", 0)
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
	sig, err = g.SubmitReplaceOrderV2(ctx, "", ownerAddr, payerAddr, marketAddr, orderSide, orderAmount, orderPrice/2, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check order #2 is in orderbook
	orders, err = g.GetOpenOrdersV2(ctx, marketAddr, ownerAddr, "", "", 0)
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
	sigs, err := g.SubmitCancelOrderV2(ctx, "", 0, pb.Side_S_ASK, ownerAddr, marketAddr, ooAddr, provider.SubmitOpts{
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
	sig, err := g.SubmitOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderAmount, orderPrice, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #1, signature %s", sig)
	time.Sleep(time.Minute)
	// Check orders are there
	orders, err := g.GetOpenOrdersV2(ctx, marketAddr, ownerAddr, "", "", 0)
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
	sig, err = g.SubmitReplaceOrderV2(ctx, found1.OrderID, ownerAddr, payerAddr, marketAddr, orderSide, orderAmount, orderPrice/2, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #2, signature %s", sig)

	time.Sleep(time.Minute)

	// Check orders are there
	orders, err = g.GetOpenOrdersV2(ctx, marketAddr, ownerAddr, "", "", 0)
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
	sigs, err := g.SubmitCancelOrderV2(ctx, "", 0, pb.Side_S_ASK, ownerAddr, marketAddr, ooAddr, provider.SubmitOpts{
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
	sig, err := g.SubmitTradeSwap(ctx, ownerAddr, "USDT",
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

func callRaydiumSwap(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting Raydium swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("Raydium swap")
	sig, err := g.SubmitRaydiumSwap(ctx, &pb.PostRaydiumSwapRequest{
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

func callJupiterSwap(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting Jupiter swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("Jupiter swap")
	sig, err := g.SubmitJupiterSwap(ctx, &pb.PostJupiterSwapRequest{
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

func callRouteTradeSwap(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting route trade swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("route trade swap")
	sig, err := g.SubmitRouteTradeSwap(ctx, &pb.RouteTradeSwapRequest{
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

func callRaydiumRouteSwap(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting Raydium route swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("Raydium route  swap")
	sig, err := g.SubmitRaydiumRouteSwap(ctx, &pb.PostRaydiumRouteSwapRequest{
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

func callJupiterRouteSwap(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting Jupiter route swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("Jupiter route  swap")
	sig, err := g.SubmitJupiterRouteSwap(ctx, &pb.PostJupiterRouteSwapRequest{
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

func callSwapsGRPCStream(g *provider.GRPCClient) bool {
	log.Info("starting get swaps stream")

	ch := make(chan *pb.GetSwapsStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := g.GetSwapsStream(ctx, []pb.Project{pb.Project_P_RAYDIUM}, []string{"58oQChx4yWmvKdwLLZzBi4ChoCc2fqCUWBkwMihLYQo2"}, true) // SOL-USDC Raydium pool
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

func callBlockGRPCStream(g *provider.GRPCClient) bool {
	log.Info("starting get block stream")

	ch := make(chan *pb.GetBlockStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := g.GetBlockStream(ctx)
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

func callDriftPerpOrderbookGRPCStream(g *provider.GRPCClient) bool {
	log.Info("starting get Drift perp orderbook stream")

	ch := make(chan *pb.GetPerpOrderbooksStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := g.GetPerpOrderbooksStream(ctx, &pb.GetPerpOrderbooksRequest{
		Contracts: []common.PerpContract{common.PerpContract_SOL_PERP},
		Limit:     0,
		Project:   pb.Project_P_DRIFT,
	})
	if err != nil {
		log.Errorf("error with GetPerpOrderbooksStream stream request: %v", err)
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

func callDriftMarginOrderbooksGRPCStream(g *provider.GRPCClient) bool {
	log.Info("starting get Drift spot orderbook stream")

	ch := make(chan *pb.GetDriftMarginOrderbooksStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := g.GetDriftMarginOrderbooksStream(ctx, &pb.GetDriftMarginOrderbooksRequest{
		Markets: []string{"SOL"},
		Limit:   0,
	})
	if err != nil {
		log.Errorf("error with GetDriftMarginOrderbooksStream request: %v", err)
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

func callDriftGetPerpTradesStream(g *provider.GRPCClient) bool {
	log.Info("starting get Drift PerpTrades stream")

	ch := make(chan *pb.GetPerpTradesStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := g.GetPerpTradesStream(ctx, &pb.GetPerpTradesStreamRequest{
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

func callDriftGetMarketDepthsStream(g *provider.GRPCClient) bool {
	log.Info("starting get Drift MarketDepth stream")

	ch := make(chan *pb.GetDriftMarketDepthStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := g.GetDriftMarketDepthsStream(ctx, &pb.GetDriftMarketDepthsStreamRequest{
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

func callDriftPerpOrderbookGRPC(g *provider.GRPCClient) bool {
	orderbook, err := g.GetPerpOrderbook(context.Background(), &pb.GetPerpOrderbookRequest{
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

func callDriftMarketDepthGRPC(g *provider.GRPCClient) bool {
	marketDepth, err := g.GetDriftMarketDepth(context.Background(), &pb.GetDriftMarketDepthRequest{
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

func callDriftGetMarginOrderbookGRPC(g *provider.GRPCClient) bool {
	orderbook, err := g.GetDriftMarginOrderbook(context.Background(), &pb.GetDriftMarginOrderbookRequest{
		Market:   "SOL",
		Limit:    0,
		Metadata: true,
	})
	if err != nil {
		log.Errorf("error with callDriftMarginOrderbookGRPC request for SOL-MARGIN: %v", err)
		return true
	} else {
		log.Info(orderbook)
	}

	fmt.Println()
	return false
}

func callGetOpenPerpOrders(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callGetOpenPerpOrders test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	perpOrders, err := g.GetOpenPerpOrders(ctx, &pb.GetOpenPerpOrdersRequest{
		OwnerAddress:   ownerAddr,
		AccountAddress: "",
		Contracts:      []common.PerpContract{common.PerpContract_SOL_PERP},
		Project:        pb.Project_P_DRIFT,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetOpenPerpOrders resp : %s", perpOrders)
	return false
}

func callGetDriftOpenMarginOrders(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callGetDriftOpenMarginOrders test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	openMarginOrders, err := g.GetDriftOpenMarginOrders(ctx, &pb.GetDriftOpenMarginOrdersRequest{
		OwnerAddress:   ownerAddr,
		AccountAddress: "",
		Markets:        []string{"SOL"},
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetDriftOpenMarginOrders resp : %s", openMarginOrders)
	return false
}

func callGetDriftOpenMarginOrder(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callGetDriftOpenMarginOrder test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	openMarginOrders, err := g.GetDriftOpenMarginOrder(ctx, &pb.GetDriftOpenMarginOrderRequest{
		OwnerAddress:   ownerAddr,
		AccountAddress: "",
		ClientOrderID:  13,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetDriftOpenMarginOrder resp : %s", openMarginOrders)
	return false
}

func callGetPerpPositions(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callGetPerpPositions test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	perpPositions, err := g.GetPerpPositions(ctx, &pb.GetPerpPositionsRequest{
		OwnerAddress:   ownerAddr,
		AccountAddress: "",
		Contracts:      []common.PerpContract{common.PerpContract_SOL_PERP},
		Project:        pb.Project_P_DRIFT,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("GetPerpPositions resp : %s", perpPositions)
	return false
}

func callGetDriftPerpPositions(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callGetDriftPerpPositions test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	perpPositions, err := g.GetDriftPerpPositions(ctx, &pb.GetDriftPerpPositionsRequest{
		OwnerAddress:   ownerAddr,
		AccountAddress: "",
		Contracts:      []string{"SOL_PERP"},
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("GetDriftPerpPositions resp : %s", perpPositions)
	return false
}

func callGetUser(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callGetUser test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := g.GetUser(ctx, &pb.GetUserRequest{
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

func callGetDriftUser(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callGetDriftUser test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := g.GetDriftUser(ctx, &pb.GetDriftUserRequest{
		OwnerAddress: ownerAddr,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetDriftUser resp : %s", user)
	return false
}

func callCancelPerpOrder(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callCancelPerpOrder test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := g.PostCancelPerpOrder(ctx, &pb.PostCancelPerpOrderRequest{
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

func callDriftCancelPerpOrder(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callDriftCancelPerpOrder test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := g.PostDriftCancelPerpOrder(ctx, &pb.PostDriftCancelPerpOrderRequest{
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

func callCancelDriftMarginOrder(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callCancelDriftMarginOrder test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := g.PostCancelDriftMarginOrder(ctx, &pb.PostCancelDriftMarginOrderRequest{
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

func callClosePerpPositions(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callClosePerpPositions test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	sig, err := g.PostClosePerpPositions(ctx, &pb.PostClosePerpPositionsRequest{
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

func callPostCloseDriftPerpPositions(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callPostCloseDriftPerpPositions test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	sig, err := g.PostCloseDriftPerpPositions(ctx, &pb.PostCloseDriftPerpPositionsRequest{
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

func callCreateUser(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callCreateUser test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := g.PostCreateUser(ctx, &pb.PostCreateUserRequest{
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

func callPostCreateDriftUser(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callPostCreateDriftUser test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := g.PostCreateDriftUser(ctx, &pb.PostCreateDriftUserRequest{
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

func callPostPerpOrder(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callPostPerpOrder test")

	request, _ := g.PostPerpOrder(context.Background(), &pb.PostPerpOrderRequest{
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
	})
	log.Infof("callPostPerpOrder request : %s", request)
	return false
}

func callPostModifyOrder(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callPostModifyOrder test")

	request, _ := g.PostModifyDriftOrder(context.Background(), &pb.PostModifyDriftOrderRequest{
		OwnerAddress:    ownerAddr,
		AccountAddress:  "",
		NewLimitPrice:   1000,
		NewPositionSide: "long",
		OrderID:         1,
	})
	log.Infof("callPostModifyOrder request : %s", request)
	return false
}

func callPostMarginOrder(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callPostMarginOrder test")

	req, _ := g.PostDriftMarginOrder(context.Background(), &pb.PostDriftMarginOrderRequest{
		OwnerAddress:   ownerAddr,
		Market:         "SOL",
		AccountAddress: "",
		PositionSide:   "short",
		Slippage:       10,
		Type:           "limit",
		Amount:         1,
		Price:          1000,
		ClientOrderID:  2,
	})
	b, _ := protojson.Marshal(req)
	log.Infof("callPostMarginOrder request : %s", string(b))
	return false
}

func callManageCollateralWithdraw(g *provider.GRPCClient) bool {
	log.Info("starting callManageCollateral withdraw test")

	request, err := g.PostManageCollateral(context.Background(), &pb.PostManageCollateralRequest{
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
	b, _ := protojson.Marshal(request)
	log.Infof("callManageCollateral request : %s", string(b))
	return false
}

func callPostDriftManageCollateralWithdraw(g *provider.GRPCClient) bool {
	log.Info("starting callPostDriftManageCollateralWithdraw test")

	request, err := g.PostManageCollateral(context.Background(), &pb.PostManageCollateralRequest{
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
	b, _ := protojson.Marshal(request)
	log.Infof("callManageCollateral request : %s", string(b))
	return false
}

func callManageCollateralTransfer(g *provider.GRPCClient) bool {
	log.Info("starting callManageCollateral transfer test")

	sig, err := g.PostManageCollateral(context.Background(), &pb.PostManageCollateralRequest{
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

func callPostDriftManageCollateralTransfer(g *provider.GRPCClient) bool {
	log.Info("starting callPostDriftManageCollateralTransfer test")

	sig, err := g.PostDriftManageCollateral(context.Background(), &pb.PostDriftManageCollateralRequest{
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

func callDriftEnableMarginTrading(g *provider.GRPCClient, ownerAddress string) bool {
	log.Info("starting callDriftEnableMarginTrading transfer test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := g.PostDriftEnableMarginTrading(ctx, &pb.PostDriftEnableMarginTradingRequest{
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

func callManageCollateralDeposit(g *provider.GRPCClient) bool {
	log.Info("starting callManageCollateral deposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := g.PostManageCollateral(ctx, &pb.PostManageCollateralRequest{
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

func callPostDriftManageCollateralDeposit(g *provider.GRPCClient) bool {
	log.Info("starting callPostDriftManageCollateralDeposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := g.PostDriftManageCollateral(ctx, &pb.PostDriftManageCollateralRequest{
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

func callGetOpenPerpOrder(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callGetOpenPerpOrder test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := g.GetOpenPerpOrder(ctx, &pb.GetOpenPerpOrderRequest{
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

func callGetDriftOpenPerpOrder(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callGetDriftOpenPerpOrder test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := g.GetDriftOpenPerpOrder(ctx, &pb.GetDriftOpenPerpOrderRequest{
		OwnerAddress:   ownerAddr,
		AccountAddress: "",
		OrderID:        1,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("GetDriftOpenPerpOrder resp : %s", user)
	return false
}

func callGetDriftOpenPerpOrders(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callGetDriftOpenPerpOrders test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := g.GetDriftOpenPerpOrders(ctx, &pb.GetDriftOpenPerpOrdersRequest{
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

func callPostSettlePNL(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callPostSettlePNL deposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := g.PostSettlePNL(ctx, &pb.PostSettlePNLRequest{
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

func callPostDriftSettlePNL(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callPostDriftSettlePNL deposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := g.PostDriftSettlePNL(ctx, &pb.PostDriftSettlePNLRequest{
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

func callPostSettlePNLs(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callPostSettlePNLs deposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := g.PostSettlePNLs(ctx, &pb.PostSettlePNLsRequest{
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

func callPostDriftSettlePNLs(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callPostDriftSettlePNLs deposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := g.PostDriftSettlePNLs(ctx, &pb.PostDriftSettlePNLsRequest{
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

func callGetAssets(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callGetAssets test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := g.GetAssets(ctx, &pb.GetAssetsRequest{
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

func callGetDriftAssets(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callGetDriftAssets test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := g.GetDriftAssets(ctx, &pb.GetDriftAssetsRequest{
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

func callGetPerpContracts(g *provider.GRPCClient) bool {
	log.Info("starting callGetPerpContracts test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := g.GetPerpContracts(ctx, &pb.GetPerpContractsRequest{
		Project: pb.Project_P_DRIFT,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetPerpContracts resp : %s", user)
	return false
}

func callGetDriftPerpContracts(g *provider.GRPCClient) bool {
	log.Info("starting callGetDriftPerpContracts test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := g.GetDriftPerpContracts(ctx, &pb.GetDriftPerpContractsRequest{})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetDriftPerpContracts resp : %s", user)
	return false
}

func callGetDriftMarkets(g *provider.GRPCClient) bool {
	log.Info("starting callGetDriftMarkets test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	user, err := g.GetDriftMarkets(ctx, &pb.GetDriftMarketsRequest{})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("callGetDriftMarkets resp : %s", user)
	return false
}

func callPostLiquidatePerp(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callPostLiquidatePerp deposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := g.PostLiquidatePerp(ctx, &pb.PostLiquidatePerpRequest{
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

func callGetDriftPerpOrderbook(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callGetDriftPerpOrderbook deposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := g.GetDriftPerpOrderbook(ctx, &pb.GetDriftPerpOrderbookRequest{
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

func callPostLiquidateDriftPerp(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting callPostLiquidateDriftPerp deposit test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sig, err := g.PostLiquidateDriftPerp(ctx, &pb.PostLiquidateDriftPerpRequest{
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
