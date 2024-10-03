package main

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	"github.com/manifoldco/promptui"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"

	"github.com/bloXroute-Labs/solana-trader-client-go/examples/config"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	"github.com/bloXroute-Labs/solana-trader-proto/common"

	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	log "github.com/sirupsen/logrus"
)

const (
	sideAsk      = "ask"
	typeLimit    = "limit"
	computePrice = 100000
	computeLimit = 5000
)

type EnvironmentVariables struct {
	privateKey        string
	publicKey         string
	openOrdersAddress string
	payer             string
}

var Environment EnvironmentVariables

func main() {
	utils.InitLogger()

	Environment = initializeEnvironmentVariables()

	listAllEndpoints()

	envPrompt := promptui.Select{
		Label: "Select environment",
		Items: []string{"mainnet", "testnet", "local"},
	}

	_, environment, err := envPrompt.Run()
	if err != nil {
		panic(fmt.Errorf("prompt failed: %v", err))
	}

	for {
		client := setupGRPCClient(config.Env(environment))
		if err != nil {
			log.Fatalf("failed to setup GRPC client: %v", err)
		}

		var names []string
		for name := range ExampleEndpoints {
			names = append(names, name)
		}

		// Choose example
		examplePrompt := promptui.Select{
			Label: "Select example to run",
			Items: names,
		}

		_, exampleName, err := examplePrompt.Run()
		if err != nil {
			fmt.Println("signal interrupt detected")
			os.Exit(1)
		}

		exampleStruct := ExampleEndpoints[exampleName]

		if exampleName == "runAllExamples" {
			for _, content := range ExampleEndpoints {
				if !content.requiresAdditionalEnvironmentVars {
					if failed := content.run(client); failed {
						log.Errorf(fmt.Sprintf("example '%s' failed", exampleName))
						time.Sleep(1 * time.Second)
					}
					time.Sleep(1 * time.Second)
					log.Printf("example '%s' completed successfully\n", exampleName)
				}
			}
		}

		log.Printf("running example: %s\n", exampleName)
		if failed := exampleStruct.run(client); failed {
			log.Errorf(fmt.Sprintf("example '%s' failed", exampleName))
			time.Sleep(1 * time.Second)
		} else {

			time.Sleep(1 * time.Second)
			log.Printf("example '%s' completed successfully\n", exampleName)
		}
	}

}

func initializeEnvironmentVariables() EnvironmentVariables {
	if os.Getenv("AUTH_HEADER") == "" {
		log.Fatal("must specify bloXroute authorization header!")
	}

	privateKey, ok := os.LookupEnv("PRIVATE_KEY")
	if !ok {
		log.Errorf(fmt.Sprintf("PRIVATE_KEY environment variable not set, cannot run any examples that require tx submission"))
	}

	// you must specify:
	//	- PRIVATE_KEY (by default loaded during provider.NewClient()) to sign transactions
	// 	- PUBLIC_KEY to indicate which account you wish to trade from
	//	- OPEN_ORDERS to indicate your Serum account to speed up lookups (optional in actual usage)
	ownerAddr, ok := os.LookupEnv("PUBLIC_KEY")
	if !ok {
		log.Warnf(fmt.Sprintf("PUBLIC_KEY environment variable not set: will skip place/cancel/settle examples"))
	}

	ooAddr, ok := os.LookupEnv("OPEN_ORDERS")
	if !ok {
		log.Errorf("OPEN_ORDERS environment variable not set: requests will be slower")
	}

	payerAddr, ok := os.LookupEnv("PAYER")
	if !ok {
		log.Warnf("PAYER environment variable not set: will be set to owner address")
		payerAddr = ownerAddr
	}

	return EnvironmentVariables{
		privateKey:        privateKey,
		publicKey:         ownerAddr,
		openOrdersAddress: ooAddr,
		payer:             payerAddr,
	}
}

func setupGRPCClient(env config.Env) *provider.GRPCClient {
	var g *provider.GRPCClient
	var err error
	switch env {
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

	return g
}

func listAllEndpoints() {
	fmt.Println(fmt.Sprintf("Available Endpoints (see docs for more info: https://docs.bloxroute.com/solana/trader-api-v2) \n"))

	var names []string
	for name := range ExampleEndpoints {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		ex := ExampleEndpoints[name]
		var extraStr string
		if ex.requiresAdditionalEnvironmentVars {
			extraStr = " (requires additional environment variables to be enabled)"
		}

		fmt.Printf("  %-40s %s%s\n", name, ex.description, extraStr)
	}
}

type ExampleFunc func(*provider.GRPCClient) bool

var ExampleEndpoints = map[string]struct {
	run                               ExampleFunc
	description                       string
	requiresAdditionalEnvironmentVars bool
}{
	"getPools": {
		run:         callPoolsGRPC,
		description: "fetch all available markets",
	},

	"getRaydiumPoolReserve": {
		run:         callRaydiumPoolReserveGRPC,
		description: "get raydium pool reserve",
	},
	"getMarkets": {
		run:         callMarketsGRPC,
		description: "fetch all available markets",
	},
	"getOrderbook": {
		run:         callOrderbookGRPC,
		description: "fetch orderbook for specific market",
	},
	"getMarketDepth": {
		run:         callMarketDepthGRPC,
		description: "get market depth",
	},
	"getOpenOrders": {
		run:         callOpenOrdersGRPC,
		description: "get open orders",
	},

	"getTickers": {
		run:         callTickersGRPC,
		description: "get tickers",
	},

	"getTransaction": {
		run:         callGetTransactionGRPC,
		description: "get tickers",
	},

	"getRateLimit": {
		run:         callGetRateLimitGRPC,
		description: "get rate limit",
	},
	"getRaydiumPools": {
		run:         callRaydiumPoolsGRPC,
		description: "get raydium pools",
	},
	"getPrice": {
		run:         callPriceGRPC,
		description: "get raydium pools",
	},
	"getRaydiumPrices": {
		run:         callRaydiumPricesGRPC,
		description: "get raydium prices",
	},
	"getJupiterPrices": {
		run:         callJupiterPricesGRPC,
		description: "get jupiter prices",
	},
	"orderbookStream": {
		run:         callOrderbookGRPCStream,
		description: "stream orderbook updates (slow example)",
	},
	"marketDepthStream": {
		run:         callMarketDepthGRPCStream,
		description: "stream market depth updates (slow example)",
	},
	"getTickersStream": {
		run:         callGetTickersGRPCStream,
		description: "stream get tickers",
	},
	"getPricesStream": {
		run:         callPricesGRPCStream,
		description: "stream prices",
	},
	"getTradesStream": {
		run:         callTradesGRPCStream,
		description: "stream trades",
	},
	"getSwapsStream": {
		run:         callSwapsGRPCStream,
		description: "stream swaps",
	},
	"getNewRaydiumPoolStream": {
		run:         callGetNewRaydiumPoolsStream,
		description: "stream new raydium pools",
	},
	"getNewRaydiumPoolsStreamWithCPMM": {
		run:         callGetNewRaydiumPoolsStreamWithCPMM,
		description: "stream new raydium pools with cpmm enabled",
	},
	"getUnsettled": {
		run:         callUnsettledGRPC,
		description: "get unsettled",
	},
	"getAccountBalance": {
		run:         callGetAccountBalanceGRPC,
		description: "get account balance",
	},

	"getQuotes": {
		run:         callGetQuotes,
		description: "get quotes",
	},

	"getRaydiumQuotes": {
		run:         callGetRaydiumQuotes,
		description: "get raydium quotes",
	},

	"getPumpFunQuotes": {
		run:         callGetPumpFunQuotes,
		description: "get pump fun quotes",
	},

	"getJupiterQuotes": {
		run:         callGetJupiterQuotes,
		description: "get jupiter quotes",
	},

	"recentBlockhashStream": {
		run:         callRecentBlockHashGRPCStream,
		description: "recent blockhash stream",
	},
	"poolReservesStream": {
		run:         callPoolReservesGRPCStream,
		description: "recent blockhash stream",
	},
	"blockStream": {
		run:         callBlockGRPCStream,
		description: "block stream",
	},
	"getPriorityFee": {
		run:         callGetPriorityFeeGRPC,
		description: "get priority fee",
	},
	"getPriorityFeeStream": {
		run:         callGetPriorityFeeGRPCStream,
		description: "get priority fee stream",
	},
	"getPumpFunNewTokenStream": {
		run:         callGetpumpFunNewTokenGRPCStreamWrap,
		description: "get pump fun new token stream",
	},

	"getBundleTipStream": {
		run:         callGetBundleTipGRPCStream,
		description: "get bundle tip stream",
	},

	"getTokenAccounts": {
		run:                               callGetTokenAccountsGRPCWrap,
		description:                       "get token accounts",
		requiresAdditionalEnvironmentVars: true,
	},

	"placeOrderWithBundle": {
		run:                               callPlaceOrderBundleWrap,
		description:                       "place a new order (openbook)",
		requiresAdditionalEnvironmentVars: true,
	},

	"placeOrderWithStakedRPCs": {
		run:                               callPlaceOrderWithStakedRPCsWrap,
		description:                       "place order (openbook) with staked rpcs and tip",
		requiresAdditionalEnvironmentVars: true,
	},

	"placeOrderWithBundleBatch": {
		run:                               callPlaceOrderBundleWithBatchWrap,
		description:                       "place order (openbook) with bundle with batch",
		requiresAdditionalEnvironmentVars: true,
	},

	"placeOrderWithPriorityFee": {
		run:                               callPlaceOrderGRPCWithPriorityFeeWrap,
		description:                       "place order (openbook) with priority fee",
		requiresAdditionalEnvironmentVars: true,
	},

	"replaceByClientOrderID": {
		run:                               callReplaceByClientOrderIDWrap,
		description:                       "replace order by client id (openbook)",
		requiresAdditionalEnvironmentVars: true,
	},

	"replaceOrder": {
		run:                               callReplaceOrderWrap,
		description:                       "replace order (openbook)",
		requiresAdditionalEnvironmentVars: true,
	},

	"tradeSwap": {
		run:                               callTradeSwapWrap,
		description:                       "trade swap",
		requiresAdditionalEnvironmentVars: true,
	},

	"routeTradeSwap": {
		run:                               callRouteTradeSwapWrap,
		description:                       "route trade swap",
		requiresAdditionalEnvironmentVars: true,
	},

	"raydiumTradeSwap": {
		run:                               callRaydiumSwapWrap,
		description:                       "raydium trade swap",
		requiresAdditionalEnvironmentVars: true,
	},

	"jupiterTradeSwap": {
		run:                               callJupiterSwapWrap,
		description:                       "jupiter trade swap",
		requiresAdditionalEnvironmentVars: true,
	},

	"pumpFunSwap": {
		run:                               callPostPumpFunSwapWrap,
		description:                       "pump fun swap",
		requiresAdditionalEnvironmentVars: true,
	},

	"raydiumRouteSwap": {
		run:                               callRaydiumRouteSwapWrap,
		description:                       "raydium route swap",
		requiresAdditionalEnvironmentVars: true,
	},

	"jupiterRouteSwap": {
		run:                               callJupiterRouteSwapWrap,
		description:                       "call jupiter route swap",
		requiresAdditionalEnvironmentVars: true,
	},

	"raydiumSwapWithInstructions": {
		run:                               callRaydiumSwapInstructionsWrap,
		description:                       "call raydium swap with instructions",
		requiresAdditionalEnvironmentVars: true,
	},
	"jupiterSwapWithInstructions": {
		run:                               callJupiterSwapInstructionsWrap,
		description:                       "call jupiter swap with instructions",
		requiresAdditionalEnvironmentVars: true,
	},

	"orderLifeCycleTest": {
		run:                               orderLifecycleTestWrap,
		description:                       "order lifecycle test",
		requiresAdditionalEnvironmentVars: true,
	},
	"cancelAll": {
		run:                               cancelAllWrap,
		description:                       "cancel all test (run order lifecycle before)",
		requiresAdditionalEnvironmentVars: true,
	},
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orders, err := g.GetOpenOrdersV2(ctx, "SOLUSDC",
		"FFqDwRq8B4hhFKRqx7N1M6Dg6vU699hVqeynDeYJdPj5", "", "", 0)
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

func callGetTokenAccountsGRPCWrap(g *provider.GRPCClient) bool {
	return callGetTokenAccountsGRPC(g, Environment.openOrdersAddress)

}

func callGetTokenAccountsGRPC(g *provider.GRPCClient, ownerAddr string) bool {
	response, err := g.GetTokenAccounts(context.Background(), &pb.GetTokenAccountsRequest{
		OwnerAddress: ownerAddr,
	})
	if err != nil {
		log.Errorf("error with GetTokenAccounts request %v", err)
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
		// prints too much info
		log.Traceln(pools)
	}

	fmt.Println()
	return false
}

func callGetRateLimitGRPC(g *provider.GRPCClient) bool {
	tx, err := g.GetRateLimit(context.Background(), &pb.GetRateLimitRequest{})
	if err != nil {
		log.Errorf("error with GetRateLimit request: %v", err)
		return true
	} else {
		log.Info(tx)
	}

	return false
}

func callGetTransactionGRPC(g *provider.GRPCClient) bool {
	tx, err := g.GetTransaction(context.Background(), &pb.GetTransactionRequest{
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

func callRaydiumPoolReserveGRPC(g *provider.GRPCClient) bool {
	pools, err := g.GetRaydiumPoolReserve(context.Background(), &pb.GetRaydiumPoolReserveRequest{
		PairsOrAddresses: []string{
			"HZ1znC9XBasm9AMDhGocd9EHSyH8Pyj1EUdiPb4WnZjo",
			"D8wAxwpH2aKaEGBKfeGdnQbCc2s54NrRvTDXCK98VAeT",
			"DdpuaJgjB2RptGMnfnCZVmC4vkKsMV6ytRa2gggQtCWt",
			"AVs9TA4nWDzfPJE9gGVNJMVhcQy3V9PGazuz33BfG2RA",
			"58oQChx4yWmvKdwLLZzBi4ChoCc2fqCUWBkwMihLYQo2",
			"7XawhbbxtsRcQA8KTkHT9f9nc6d69UwqCDh6U5EEbEmX",
		},
	})
	if err != nil {
		log.Errorf("error with GetRaydiumPoolReserve request for Raydium: %v", err)
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
		// prints too much info
		log.Traceln(pools)
	}

	fmt.Println()
	return false
}

func callPriceGRPC(g *provider.GRPCClient) bool {
	prices, err := g.GetPrice(context.Background(), []string{"So11111111111111111111111111111111111111112", "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263"})
	if err != nil {
		log.Errorf("error with GetPrice request for SOL and BONK: %v", err)
		return true
	} else {
		log.Info(prices)
	}

	fmt.Println()
	return false
}

func callRaydiumPricesGRPC(g *provider.GRPCClient) bool {
	prices, err := g.GetRaydiumPrices(context.Background(), &pb.GetRaydiumPricesRequest{
		Tokens: []string{"So11111111111111111111111111111111111111112", "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263"},
	})
	if err != nil {
		log.Errorf("error with GetRaydiumPrices request for SOL and BONK: %v", err)
		return true
	} else {
		log.Info(prices)
	}

	fmt.Println()
	return false
}

func callJupiterPricesGRPC(g *provider.GRPCClient) bool {
	prices, err := g.GetJupiterPrices(context.Background(), &pb.GetJupiterPricesRequest{
		Tokens: []string{"So11111111111111111111111111111111111111112", "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263"},
	})
	if err != nil {
		log.Errorf("error with GetJupiterPrices request for SOL and BONK: %v", err)
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

	inToken := "So11111111111111111111111111111111111111112"
	outToken := "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
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

	inToken := "So11111111111111111111111111111111111111112"
	outToken := "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
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

func callGetPumpFunQuotes(g *provider.GRPCClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	amount := 0.01
	slippage := float64(5)

	quotes, err := g.GetPumpFunQuotes(ctx, &pb.GetPumpFunQuotesRequest{
		QuoteType:           "buy",
		BondingCurveAddress: "Dga6eouREJ4kLHMqWWtccGGPsGebexuBYrcepBVd494q",
		MintAddress:         "9QG5NHnfqQCyZ9SKhz7BzfjPseTFWaApmAtBTziXLanY",
		Amount:              amount,
		Slippage:            slippage,
	})
	if err != nil {
		return true
	}

	log.Infof("best quote for PumpFun is %v", quotes)

	fmt.Println()
	return false
}

func callGetJupiterQuotes(g *provider.GRPCClient) bool {
	log.Info("starting get Jupiter quotes test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inToken := "So11111111111111111111111111111111111111112"
	outToken := "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
	amount := 0.01
	slippage := float64(5)

	quotes, err := g.GetJupiterQuotes(ctx, &pb.GetJupiterQuotesRequest{
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

func callOrderbookGRPCStream(g *provider.GRPCClient) bool {
	log.Info("starting orderbook stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream error response
	stream, err := g.GetOrderbookStream(ctx, []string{"", "xxx"}, 3, pb.Project_P_OPENBOOK)
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
	stream, err := g.GetPoolReservesStream(ctx, &pb.GetPoolReservesStreamRequest{
		Projects: []pb.Project{pb.Project_P_RAYDIUM},
		Pools: []string{
			"HZ1znC9XBasm9AMDhGocd9EHSyH8Pyj1EUdiPb4WnZjo",
			"D8wAxwpH2aKaEGBKfeGdnQbCc2s54NrRvTDXCK98VAeT",
			"DdpuaJgjB2RptGMnfnCZVmC4vkKsMV6ytRa2gggQtCWt",
			"AVs9TA4nWDzfPJE9gGVNJMVhcQy3V9PGazuz33BfG2RA",
			"58oQChx4yWmvKdwLLZzBi4ChoCc2fqCUWBkwMihLYQo2",
			"7XawhbbxtsRcQA8KTkHT9f9nc6d69UwqCDh6U5EEbEmX",
		},
	})

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
	//marketAddr = "RAY/USDC"
	//marketAddr = "9Lyhks5bQQxb9EyyX55NtgKQzpM4WK7JCmeaWuQ5MoXD"

	orderPrice  = float64(170200)
	orderAmount = float64(0.001)
)

func orderLifecycleTestWrap(g *provider.GRPCClient) bool {
	return orderLifecycleTest(g, Environment.publicKey, Environment.payer, Environment.openOrdersAddress)
}

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

	clientID, failed := callPlaceOrderGRPC(g, ownerAddr, payerAddr, ooAddr, sideAsk, typeLimit)

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

func callPlaceOrderGRPC(g *provider.GRPCClient, ownerAddr, payerAddr, ooAddr string, orderSide string, orderType string) (uint64, bool) {
	log.Info("starting place order")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// generate a random clientOrderID for this order
	rand.Seed(time.Now().UnixNano())
	clientOrderID := rand.Uint64()

	opts := provider.PostOrderOpts{
		ClientOrderID:     clientOrderID,
		OpenOrdersAddress: ooAddr,
		SkipPreFlight:     config.BoolPtr(true),
	}

	// create order without actually submitting
	response, err := g.PostOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderType, orderAmount, orderPrice, nil, opts)
	if err != nil {
		log.Errorf("failed to create order (%v)", err)
		return 0, true
	}
	log.Infof("created unsigned place order transaction: %v", response.Transaction)

	// sign/submit transaction after creation
	sig, err := g.SubmitOrderV2(ctx, ownerAddr, ownerAddr, marketAddr,
		orderSide, orderType, orderAmount, orderPrice, nil, opts)
	if err != nil {
		log.Errorf("failed to submit order (%v)", err)
		return 0, true
	}

	log.Infof("placed order %v with clientOrderID %v", sig, clientOrderID)
	return clientOrderID, false
}

func callPlaceOrderBundleWrap(g *provider.GRPCClient) bool {
	return callPlaceOrderBundle(g, Environment.publicKey, Environment.payer, sideAsk,
		computeLimit, computePrice, typeLimit, 100000)
}

func callPlaceOrderBundle(g *provider.GRPCClient, ownerAddr, payerAddr,
	orderSide string, computeLimit uint32, computePrice uint64, orderType string, tipAmount uint64) bool {
	log.Info("starting place order with bundle")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// generate a random clientOrderID for this order
	rand.Seed(time.Now().UnixNano())
	clientOrderID := rand.Uint64()

	opts := provider.PostOrderOpts{
		ClientOrderID: clientOrderID,
		SkipPreFlight: config.BoolPtr(true),
	}

	// create order without actually submitting
	response, err := g.PostOrderV2WithPriorityFee(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderType,
		orderAmount, orderPrice, computeLimit, computePrice, &tipAmount, opts)
	if err != nil {
		log.Errorf("failed to create order (%v)", err)
		return true
	}
	log.Infof("created unsigned place order transaction: %v", response.Transaction)

	resp, err := g.SignAndSubmit(ctx, &pb.TransactionMessage{
		Content: response.Transaction.Content}, true, true, false)
	if err != nil {
		log.Errorf("failed to sign and submit order (%v)", err)
		return true
	}

	log.Infof("submitted bundle order to trader api %v", resp)

	return false
}

func callPlaceOrderWithStakedRPCsWrap(g *provider.GRPCClient) bool {
	return callPlaceOrderWithStakedRPCs(g, Environment.publicKey, Environment.payer, Environment.openOrdersAddress, sideAsk,
		10000, 10000, typeLimit, uint64(1100000))
}

func callPlaceOrderWithStakedRPCs(g *provider.GRPCClient, ownerAddr, payerAddr, _ string,
	orderSide string, computeLimit uint32, computePrice uint64, orderType string, tipAmount uint64) bool {
	log.Info("starting place order with bundle")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// generate a random clientOrderID for this order
	rand.Seed(time.Now().UnixNano())
	clientOrderID := rand.Uint64()

	opts := provider.PostOrderOpts{
		ClientOrderID: clientOrderID,
		SkipPreFlight: config.BoolPtr(true),
	}

	// create order without actually submitting
	response, err := g.PostOrderV2WithPriorityFee(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderType,
		orderAmount, orderPrice, computeLimit, computePrice, &tipAmount, opts)
	if err != nil {
		log.Errorf("failed to create order (%v)", err)
		return true
	}
	log.Infof("created unsigned place order transaction: %v", response.Transaction)

	resp, err := g.SignAndSubmit(ctx, &pb.TransactionMessage{
		Content: response.Transaction.Content}, true, false, true)
	if err != nil {
		log.Errorf("failed to sign and submit order (%v)", err)
		return true
	}

	log.Infof("submitted order to trader api with staked rpcs %v", resp)

	return false
}

func callPlaceOrderBundleWithBatchWrap(g *provider.GRPCClient) bool {
	return callPlaceOrderBundleWithBatch(g, Environment.publicKey, Environment.payer, Environment.openOrdersAddress, sideAsk,
		computeLimit, computePrice, typeLimit, uint64(1000000))
}

func callPlaceOrderBundleWithBatch(g *provider.GRPCClient, ownerAddr, payerAddr, _ string,
	orderSide string, computeLimit uint32, computePrice uint64, orderType string, tipAmount uint64) bool {
	log.Info("starting to place order with bundle")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// generate a random clientOrderID for this order
	rand.Seed(time.Now().UnixNano())
	clientOrderID := rand.Uint64()

	opts := provider.PostOrderOpts{
		ClientOrderID: clientOrderID,
		SkipPreFlight: config.BoolPtr(true),
	}

	// create order without actually submitting
	response, err := g.PostOrderV2WithPriorityFee(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderType,
		orderAmount, orderPrice, computeLimit, computePrice, &tipAmount, opts)
	if err != nil {
		log.Errorf("failed to create order (%v)", err)
		return true
	}
	log.Infof("created unsigned place order transaction: %v", response.Transaction)

	signedTx, err := transaction.SignTx(response.GetTransaction().Content)
	if err != nil {
		panic(err)
	}

	useBundle := true

	batchEntry := pb.PostSubmitRequestEntry{
		Transaction:   &pb.TransactionMessage{Content: signedTx},
		SkipPreFlight: true,
	}

	batchRequest := pb.PostSubmitBatchRequest{
		Entries:        []*pb.PostSubmitRequestEntry{&batchEntry},
		SubmitStrategy: 1,
		UseBundle:      &useBundle,
	}

	batchResp, err := g.PostSubmitBatchV2(ctx, &batchRequest)
	if err != nil {
		panic(err)
	}

	log.Infof("successfully placed bundle batch order with signature : %s", batchResp.Transactions[0].Signature)

	if err != nil {
		return false
	}

	return false
}

func callPlaceOrderGRPCWithPriorityFeeWrap(g *provider.GRPCClient) bool {
	return callPlaceOrderGRPCWithPriorityFee(g, Environment.publicKey, Environment.payer, Environment.openOrdersAddress,
		sideAsk, computeLimit, computePrice, typeLimit)
}

func callPlaceOrderGRPCWithPriorityFee(g *provider.GRPCClient, ownerAddr, payerAddr, ooAddr string, orderSide string,
	computeLimit uint32, computePrice uint64, orderType string) bool {
	log.Info("starting place order")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// generate a random clientOrderID for this order
	rand.Seed(time.Now().UnixNano())
	clientOrderID := rand.Uint64()

	opts := provider.PostOrderOpts{
		ClientOrderID: clientOrderID}

	// sign/submit transaction after creation
	sig, err := g.SubmitOrderV2WithPriorityFee(ctx, ownerAddr, ownerAddr, marketAddr,
		orderSide, orderType, orderAmount, orderPrice, 0, 0, nil, opts)
	if err != nil {
		log.Errorf("failed to submit order (%v)", err)
		return true
	}

	log.Infof("placed order %v with clientOrderID %v", sig, clientOrderID)
	return false
}

func callCancelByClientOrderIDGRPC(g *provider.GRPCClient, ownerAddr, ooAddr string, clientID uint64) bool {
	log.Info("starting cancel order by client order ID")
	time.Sleep(30 * time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig, err := g.SubmitCancelOrderV2(ctx, "", clientID, sideAsk, ownerAddr,
		marketAddr, ooAddr, provider.SubmitOpts{
			SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
			SkipPreFlight:  config.BoolPtr(true),
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

func cancelAllWrap(g *provider.GRPCClient) bool {
	return cancelAll(g, Environment.publicKey, Environment.payer, Environment.openOrdersAddress, sideAsk, typeLimit)
}

func cancelAll(g *provider.GRPCClient, ownerAddr, payerAddr, ooAddr string, orderSide string, orderType string) bool {
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
	sig, err := g.SubmitOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderType, orderAmount, orderPrice, nil, opts)
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("submitting place order #1, signature %s", sig)

	opts.ClientOrderID = clientOrderID2
	sig, err = g.SubmitOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderType, orderAmount, orderPrice, nil, opts)
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
	sigs, err := g.SubmitCancelOrderV2(ctx, "", 0, sideAsk, ownerAddr, marketAddr, ooAddr, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  config.BoolPtr(true),
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

func callReplaceByClientOrderIDWrap(g *provider.GRPCClient) bool {
	return callReplaceByClientOrderID(g, Environment.publicKey, Environment.payer, Environment.openOrdersAddress, sideAsk, typeLimit)
}

func callReplaceByClientOrderID(g *provider.GRPCClient, ownerAddr, payerAddr, ooAddr string, orderSide string, orderType string) bool {
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
	sig, err := g.SubmitOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderType, orderAmount, orderPrice, nil, opts)
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
	sig, err = g.SubmitReplaceOrderV2(ctx, "", ownerAddr, payerAddr, marketAddr, orderSide, orderType, orderAmount, orderPrice/2, opts)
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
	sigs, err := g.SubmitCancelOrderV2(ctx, "", 0, sideAsk, ownerAddr, marketAddr, ooAddr, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  config.BoolPtr(true),
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

func callReplaceOrderWrap(g *provider.GRPCClient) bool {
	return callReplaceOrder(g, Environment.publicKey, Environment.payer, Environment.openOrdersAddress, sideAsk, typeLimit)
}

func callReplaceOrder(g *provider.GRPCClient, ownerAddr, payerAddr, ooAddr string, orderSide string, orderType string) bool {
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
	sig, err := g.SubmitOrderV2(ctx, ownerAddr, payerAddr, marketAddr, orderSide, orderType, orderAmount, orderPrice, nil, opts)
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
	var found1 *pb.OrderV2

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
	sig, err = g.SubmitReplaceOrderV2(ctx, found1.OrderID, ownerAddr, payerAddr, marketAddr, orderSide, orderType, orderAmount, orderPrice/2, opts)
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
	var found2 *pb.OrderV2

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
	sigs, err := g.SubmitCancelOrderV2(ctx, "", 0, sideAsk, ownerAddr, marketAddr, ooAddr, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  config.BoolPtr(true),
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
	sig, err := g.SubmitTradeSwap(ctx, ownerAddr, "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		"So11111111111111111111111111111111111111112", 0.01, 0.1, pb.Project_P_RAYDIUM, provider.SubmitOpts{
			SubmitStrategy: pb.SubmitStrategy_P_ABORT_ON_FIRST_ERROR,
			SkipPreFlight:  config.BoolPtr(false),
		})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("trade swap transaction signature : %s", sig)
	return false
}

func callRaydiumSwapWrap(g *provider.GRPCClient) bool {
	return callRaydiumSwap(g, Environment.publicKey)
}

func callRaydiumSwap(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting Raydium swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("Raydium swap")
	sig, err := g.SubmitRaydiumSwap(ctx, &pb.PostRaydiumSwapRequest{
		OwnerAddress: ownerAddr,
		InToken:      "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		OutToken:     "So11111111111111111111111111111111111111112",
		Slippage:     0.1,
		InAmount:     0.01,
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_ABORT_ON_FIRST_ERROR,
		SkipPreFlight:  config.BoolPtr(false),
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Raydium swap transaction signature : %s", sig)
	return false
}

func callPostPumpFunSwapWrap(g *provider.GRPCClient) bool {
	return callPostPumpFunSwap(Environment.publicKey)
}

func callPostPumpFunSwap(ownerAddr string) bool {
	log.Info("starting PostPumpFunSwap test")
	g, err := provider.NewGRPCClientPumpNY()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("PumpFun swap")
	sig, err := g.SubmitPostPumpFunSwap(ctx, &pb.PostPumpFunSwapRequest{
		UserAddress:         ownerAddr,
		BondingCurveAddress: "7BcRpqUC7AF5Xsc3QEpCb8xmoi2X1LpwjUBNThbjWvyo",
		TokenAddress:        "BAHY8ocERNc5j6LqkYav1Prr8GBGsHvBV5X3dWPhsgXw",
		TokenAmount:         10,
		SolThreshold:        0.0001,
		IsBuy:               false,
		ComputeLimit:        0,
		ComputePrice:        0,
		Tip:                 nil,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("PumpFun swap transaction signature : %s", sig)
	return false
}

func callJupiterSwapWrap(g *provider.GRPCClient) bool {
	return callJupiterSwap(g, Environment.publicKey)
}

func callJupiterSwap(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting Jupiter swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("Jupiter swap")
	sig, err := g.SubmitJupiterSwap(ctx, &pb.PostJupiterSwapRequest{
		OwnerAddress: ownerAddr,
		InToken:      "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		OutToken:     "So11111111111111111111111111111111111111112",
		Slippage:     0.1,
		InAmount:     0.01,
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_ABORT_ON_FIRST_ERROR,
		SkipPreFlight:  config.BoolPtr(false),
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Jupiter swap transaction signature : %s", sig)
	return false
}

func callJupiterSwapInstructionsWrap(g *provider.GRPCClient) bool {
	return callJupiterSwapInstructions(g, Environment.publicKey, uint64(1100), false)
}

func callJupiterSwapInstructions(g *provider.GRPCClient, ownerAddr string, tipAmount uint64, useBundle bool) bool {
	log.Info("starting Jupiter swap instructions test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("Jupiter swap")
	sig, err := g.SubmitJupiterSwapInstructions(ctx, &pb.PostJupiterSwapInstructionsRequest{
		OwnerAddress: ownerAddr,
		InToken:      "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		OutToken:     "So11111111111111111111111111111111111111112",
		Slippage:     0.4,
		InAmount:     0.001,
		Tip:          &tipAmount,
	}, useBundle, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  config.BoolPtr(false),
	})

	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Jupiter swap transaction with instructions signature : %s", sig)
	return false
}

func callRaydiumSwapInstructionsWrap(g *provider.GRPCClient) bool {
	return callRaydiumSwapInstructions(g, Environment.publicKey, uint64(1100), true)
}

func callRaydiumSwapInstructions(g *provider.GRPCClient, ownerAddr string, tipAmount uint64, useBundle bool) bool {
	log.Info("starting Raydium swap instructions test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("Raydium swap")
	sig, err := g.SubmitRaydiumSwapInstructions(ctx, &pb.PostRaydiumSwapInstructionsRequest{
		OwnerAddress: ownerAddr,
		InToken:      "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		OutToken:     "So11111111111111111111111111111111111111112",
		Slippage:     0.4,
		InAmount:     0.001,
		Tip:          &tipAmount,
	}, useBundle, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  config.BoolPtr(false),
	})

	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Raydium swap transaction with instructions signature : %s", sig)
	return false
}

func callTradeSwapWrap(g *provider.GRPCClient) bool {
	return callTradeSwap(g, Environment.publicKey)
}
func callRouteTradeSwapWrap(g *provider.GRPCClient) bool {
	return callRouteTradeSwap(g, Environment.publicKey)
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
				InToken:      "So11111111111111111111111111111111111111112",
				OutToken:     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
				InAmount:     0.01,
				OutAmountMin: 0.007505,
				OutAmount:    0.0074,
				Project: &pb.StepProject{
					Label: "Raydium",
					Id:    "58oQChx4yWmvKdwLLZzBi4ChoCc2fqCUWBkwMihLYQo2",
				},
			},
		},
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_ABORT_ON_FIRST_ERROR,
		SkipPreFlight:  config.BoolPtr(false),
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("route trade swap transaction signature : %s", sig)
	return false
}

func callRaydiumRouteSwapWrap(g *provider.GRPCClient) bool {
	return callRaydiumRouteSwap(g, Environment.publicKey)
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
				InToken:      "So11111111111111111111111111111111111111112",
				OutToken:     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
				InAmount:     0.01,
				OutAmountMin: 0.007505,
				OutAmount:    0.0074,
				Project: &pb.StepProject{
					Label: "Raydium",
					Id:    "58oQChx4yWmvKdwLLZzBi4ChoCc2fqCUWBkwMihLYQo2",
				},
			},
		},
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_ABORT_ON_FIRST_ERROR,
		SkipPreFlight:  config.BoolPtr(false),
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Raydium route swap transaction signature : %s", sig)
	return false
}

func callJupiterRouteSwapWrap(g *provider.GRPCClient) bool {
	return callJupiterRouteSwap(g, Environment.publicKey)
}

func callJupiterRouteSwap(g *provider.GRPCClient, ownerAddr string) bool {
	log.Info("starting Jupiter route swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("Jupiter route swap")
	sig, err := g.SubmitJupiterRouteSwap(ctx, &pb.PostJupiterRouteSwapRequest{
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
		SubmitStrategy: pb.SubmitStrategy_P_ABORT_ON_FIRST_ERROR,
		SkipPreFlight:  config.BoolPtr(false),
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Jupiter route swap transaction signature : %s", sig)
	return false
}

func callGetpumpFunNewTokenGRPCStreamWrap(g *provider.GRPCClient) bool {
	gg, err := provider.NewGRPCClientPumpNY()
	if err != nil {
		log.Fatal(err)
	}

	mint, res := callGetPumpFunNewTokensGRPCStream(gg)
	res = callGetPumpFunSwapsGRPCStream(gg, mint)
	return res
}

func callGetPumpFunNewTokensGRPCStream(g *provider.GRPCClient) (string, bool) {
	log.Info("starting GetPumpFunNewTokens stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := g.GetPumpFunNewTokensStream(ctx, &pb.GetPumpFunNewTokensStreamRequest{})
	if err != nil {
		log.Errorf("error with GetPumpFunNewTokens stream request: %v", err)
		return "", true
	}

	ch := stream.Channel(0)
	mint := ""
	for i := 1; i <= 1; i++ {
		v, ok := <-ch
		if !ok {
			return "", true
		}
		log.Infof("response %v received", v)
		mint = v.Mint
	}
	return mint, false
}

func callGetPumpFunSwapsGRPCStream(g *provider.GRPCClient, mint string) bool {
	log.Info("starting GetPumpFunSwaps stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := g.GetPumpFunSwapsStream(ctx, &pb.GetPumpFunSwapsStreamRequest{
		Tokens: []string{mint},
	})
	if err != nil {
		log.Errorf("error with GetPumpFunSwaps stream request: %v", err)
		return true
	}

	ch := stream.Channel(0)
	for i := 1; i <= 1; i++ {
		v, ok := <-ch
		if !ok {
			return true
		}
		log.Infof("response %v received", v)

	}
	return false
}

func callGetTickersGRPCStream(g *provider.GRPCClient) bool {
	log.Info("starting get ticker stream")

	ch := make(chan *pb.GetTickersStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := g.GetTickersStream(ctx, &pb.GetTickersStreamRequest{
		Project: pb.Project_P_OPENBOOK,
		Markets: []string{"BONK/SOL", "wSOL/RAY", "BONK/RAY", "RAY/USDC",
			"SOL/USDC", "SOL/USDC",
			"RAY/USDC", "USDT/USDC"},
	})

	if err != nil {
		log.Errorf("error with GetPrices stream request: %v", err)
		return true
	}
	stream.Into(ch)
	for i := 1; i <= 1; i++ {
		v, ok := <-ch
		if !ok {
			// channel closed
			return true
		}
		log.Infof("response %v received ", v)
	}
	return false
}

func callPricesGRPCStream(g *provider.GRPCClient) bool {
	log.Info("starting get prices stream")

	ch := make(chan *pb.GetPricesStreamResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := g.GetPricesStream(ctx, []pb.Project{pb.Project_P_RAYDIUM}, []string{"So11111111111111111111111111111111111111112", "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263"})

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

func callGetNewRaydiumPoolsStream(g *provider.GRPCClient) bool {
	log.Info("starting get new raydium pools stream without cpmm")

	ch := make(chan *pb.GetNewRaydiumPoolsResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := g.GetNewRaydiumPoolsStream(ctx, false)
	if err != nil {
		log.Errorf("error with GetNewRaydiumPools stream request: %v", err)
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

func callGetNewRaydiumPoolsStreamWithCPMM(g *provider.GRPCClient) bool {
	log.Info("starting get new raydium pools stream with cpmm")

	ch := make(chan *pb.GetNewRaydiumPoolsResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := g.GetNewRaydiumPoolsStream(ctx, true)
	if err != nil {
		log.Errorf("error with GetNewRaydiumPools stream request: %v", err)
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

func callGetPriorityFeeGRPCStream(g *provider.GRPCClient) bool {
	log.Info("starting priority fee stream")

	ch := make(chan *pb.GetPriorityFeeResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := g.GetPriorityFeeStream(ctx, pb.Project_P_RAYDIUM, nil)
	if err != nil {
		log.Errorf("error with GetPriorityFee stream request: %v", err)
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

func callGetPriorityFeeGRPC(g *provider.GRPCClient) bool {
	log.Info("starting priority fee test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	response, err := g.GetPriorityFee(ctx, &pb.GetPriorityFeeRequest{})
	if err != nil {
		log.Errorf("error with GetPriorityFee request: %v", err)
		return true
	}
	log.Infof("response received: %v", response)
	return false
}

func callGetBundleTipGRPCStream(g *provider.GRPCClient) bool {
	log.Info("starting get bundle tip stream")

	ch := make(chan *pb.GetBundleTipResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := g.GetBundleTipStream(ctx)
	if err != nil {
		log.Errorf("error with GetBundleTip stream request: %v", err)
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
