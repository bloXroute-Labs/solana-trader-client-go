package main

import (
	"context"
	"fmt"
	"github.com/manifoldco/promptui"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"

	"github.com/bloXroute-Labs/solana-trader-client-go/examples/config"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"

	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/bloXroute-Labs/solana-trader-proto/common"
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

		client := setupHTTPClient(config.Env(environment))
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

func setupHTTPClient(env config.Env) *provider.HTTPClient {

	var h *provider.HTTPClient
	var err error
	switch env {
	case config.EnvLocal:
		h = provider.NewHTTPLocal()
	case config.EnvTestnet:
		h = provider.NewHTTPTestnet()
	case config.EnvMainnet:
		h = provider.NewHTTPClient()
	}
	if err != nil {
		log.Fatalf("error connecting to HTTP client: %v", err)
	}

	return h
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

type ExampleFunc func(client *provider.HTTPClient) bool

var ExampleEndpoints = map[string]struct {
	run                               ExampleFunc
	description                       string
	requiresAdditionalEnvironmentVars bool
}{
	"getPools": {
		run:         callPoolsHTTP,
		description: "fetch all available markets",
	},

	"getTrades": {
		run:         callTradesHTTP,
		description: "get trades",
	},

	"getRaydiumPoolReserve": {
		run:         callRaydiumPoolReserveHTTP,
		description: "get raydium pool reserve",
	},
	"getMarkets": {
		run:         callMarketsHTTP,
		description: "fetch all available markets",
	},
	"getOrderbook": {
		run:         callOrderbookHTTP,
		description: "fetch orderbook for specific market",
	},
	"getMarketDepth": {
		run:         callMarketDepthHTTP,
		description: "get market depth",
	},
	"getTickers": {
		run:         callTickersHTTP,
		description: "get tickers",
	},

	"getTransaction": {
		run:         callGetTransactionHTTP,
		description: "get tickers",
	},

	"getRateLimit": {
		run:         callGetRateLimitHTTP,
		description: "get rate limit",
	},
	"getRaydiumPools": {
		run:         callRaydiumPoolsHTTP,
		description: "get raydium pools",
	},
	"getPrice": {
		run:         callPriceHTTP,
		description: "get raydium pools",
	},
	"getRecentBlockhash": {
		run:         callGetRecentBlockHashHTTP,
		description: "get recent blockhash",
	},
	"getRaydiumPrices": {
		run:         callRaydiumPricesHTTP,
		description: "get raydium prices",
	},
	"getJupiterPrices": {
		run:         callJupiterPricesHTTP,
		description: "get jupiter prices",
	},

	"getUnsettled": {
		run:         callUnsettledHTTP,
		description: "get unsettled",
	},
	"getAccountBalance": {
		run:         callGetAccountBalanceHTTP,
		description: "get account balance",
	},
	"getQuotes": {
		run:         callGetQuotesHTTP,
		description: "get quotes",
	},

	"getRaydiumQuotes": {
		run:         callGetRaydiumQuotesHTTP,
		description: "get raydium quotes",
	},

	"getPumpFunQuotes": {
		run:         callGetPumpFunQuotesHTTP,
		description: "get pump fun quotes",
	},

	"getJupiterQuotes": {
		run:         callGetJupiterQuotesHTTP,
		description: "get jupiter quotes",
	},

	"getPriorityFee": {
		run:         callGetPriorityFeeHTTP,
		description: "get priority fee",
	},

	"getTokenAccounts": {
		run:                               callGetTokenAccountsHTTPWrap,
		description:                       "get token accounts",
		requiresAdditionalEnvironmentVars: true,
	},

	"placeOrderWithBundle": {
		run:                               callPlaceOrderHTTPWrap,
		description:                       "place a new order (openbook)",
		requiresAdditionalEnvironmentVars: true,
	},

	"placeOrderWithStakedRPCs": {
		run:                               callPlaceOrderWithStakedRPCsHTTPWrap,
		description:                       "place order (openbook) with staked rpcs and tip",
		requiresAdditionalEnvironmentVars: true,
	},

	"placeOrderWithBundleBatch": {
		run:                               callPlaceOrderBundleUsingBatchHTTPWithWrap,
		description:                       "place order (openbook) with bundle with batch",
		requiresAdditionalEnvironmentVars: true,
	},

	"placeOrderWithPriorityFee": {
		run:                               callPlaceOrderHTTPWithPriorityFeeWrap,
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

	"cancelAll": {
		run:                               cancelAllWrap,
		description:                       "cancel all test (run order lifecycle before)",
		requiresAdditionalEnvironmentVars: true,
	},
}

func callMarketsHTTP(h *provider.HTTPClient) bool {
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

func callOrderbookHTTP(h *provider.HTTPClient) bool {
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

func callMarketDepthHTTP(h *provider.HTTPClient) bool {
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

func callUnsettledHTTP(h *provider.HTTPClient) bool {
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

func callGetAccountBalanceHTTP(h *provider.HTTPClient) bool {
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

func callGetTokenAccountsHTTPWrap(h *provider.HTTPClient) bool {
	return callGetTokenAccountsHTTP(h, Environment.publicKey)
}

func callGetTokenAccountsHTTP(h *provider.HTTPClient, ownerAddr string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.GetTokenAccounts(ctx, &pb.GetTokenAccountsRequest{
		OwnerAddress: ownerAddr,
	})
	if err != nil {
		log.Errorf("error with GetTokenAccounts request for : %v", err)
		return true
	} else {
		log.Info(response)
	}

	fmt.Println()
	return false
}

func callTradesHTTP(h *provider.HTTPClient) bool {
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

func callPoolsHTTP(h *provider.HTTPClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pools, err := h.GetPools(ctx, []pb.Project{pb.Project_P_RAYDIUM})
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

func callRaydiumPoolReserveHTTP(h *provider.HTTPClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pools, err := h.GetRaydiumPoolReserve(ctx, &pb.GetRaydiumPoolReserveRequest{
		PairsOrAddresses: []string{"HZ1znC9XBasm9AMDhGocd9EHSyH8Pyj1EUdiPb4WnZjo",
			"D8wAxwpH2aKaEGBKfeGdnQbCc2s54NrRvTDXCK98VAeT", "DdpuaJgjB2RptGMnfnCZVmC4vkKsMV6ytRa2gggQtCWt"},
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

func callRaydiumPoolsHTTP(h *provider.HTTPClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pools, err := h.GetRaydiumPools(ctx, &pb.GetRaydiumPoolsRequest{})
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

func callGetRateLimitHTTP(h *provider.HTTPClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := h.GetRateLimit(ctx, &pb.GetRateLimitRequest{})
	if err != nil {
		log.Errorf("error with GetTransaction request: %v", err)
		return true
	} else {
		log.Info(tx)
	}

	fmt.Println()
	return false
}

func callGetTransactionHTTP(h *provider.HTTPClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := h.GetTransaction(ctx, &pb.GetTransactionRequest{
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

func callPriceHTTP(h *provider.HTTPClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prices, err := h.GetPrice(ctx, []string{"So11111111111111111111111111111111111111112", "DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263"})
	if err != nil {
		log.Errorf("error with GetPrice request for SOL and BONK: %v", err)
		return true
	} else {
		log.Info(prices)
	}

	fmt.Println()
	return false
}

func callRaydiumPricesHTTP(h *provider.HTTPClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prices, err := h.GetRaydiumPrices(ctx, &pb.GetRaydiumPricesRequest{
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

func callJupiterPricesHTTP(h *provider.HTTPClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prices, err := h.GetJupiterPrices(ctx, &pb.GetJupiterPricesRequest{
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

func callTickersHTTP(h *provider.HTTPClient) bool {
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

func callGetQuotesHTTP(h *provider.HTTPClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	inToken := "So11111111111111111111111111111111111111112"
	outToken := "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
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

func callGetRaydiumQuotesHTTP(h *provider.HTTPClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	inToken := "So11111111111111111111111111111111111111112"
	outToken := "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
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

func callGetPumpFunQuotesHTTP(h *provider.HTTPClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	amount := 0.01
	slippage := float64(5)

	quotes, err := h.GetPumpFunQuotes(ctx, &pb.GetPumpFunQuotesRequest{
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

func callGetJupiterQuotesHTTP(h *provider.HTTPClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	inToken := "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
	outToken := "So11111111111111111111111111111111111111112"
	amount := 0.01
	slippage := float64(5)

	quotes, err := h.GetJupiterQuotes(ctx, &pb.GetJupiterQuotesRequest{
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

	for _, route := range quotes.Routes {
		log.Infof("best route for Jupiter is %v", route)
	}

	fmt.Println()
	return false
}

const (
	// SOL/USDC market
	marketAddr = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh6"

	orderPrice  = float64(170200)
	orderAmount = float64(0.1)
)

func callPlaceOrderHTTPWrap(h *provider.HTTPClient) bool {
	_, ok := callPlaceOrderHTTP(h, Environment.publicKey, Environment.openOrdersAddress, sideAsk, typeLimit)
	return ok
}

func callPlaceOrderHTTP(h *provider.HTTPClient, ownerAddr, ooAddr string, orderSide string, orderType string) (uint64, bool) {
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

func callPlaceOrderHTTPWithPriorityFeeWrap(h *provider.HTTPClient) bool {
	return callPlaceOrderHTTPWithPriorityFee(h, Environment.publicKey, Environment.openOrdersAddress, sideAsk, typeLimit,
		computeLimit, computePrice)
}

func callPlaceOrderHTTPWithPriorityFee(h *provider.HTTPClient, ownerAddr, ooAddr string, orderSide string, orderType string,
	computeLimit uint32, computePrice uint64) bool {
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
	response, err := h.PostOrderV2WithPriorityFee(ctx, ownerAddr, ownerAddr, marketAddr, orderSide, orderType,
		orderAmount, orderPrice, computeLimit, computePrice, opts)
	if err != nil {
		log.Errorf("failed to create order (%v)", err)
		return true
	}
	log.Infof("created unsigned place order transaction: %v", response.Transaction)

	// sign/submit transaction after creation
	sig, err := h.SubmitOrderV2WithPriorityFee(ctx, ownerAddr, ownerAddr, marketAddr,
		orderSide, orderType, orderAmount, orderPrice, computeLimit, computePrice, opts)
	if err != nil {
		log.Errorf("failed to submit order (%v)", err)
	}

	log.Infof("placed order %v with clientOrderID %v", sig, clientOrderID)

	return false
}

func callPlaceOrderBundleUsingBatchHTTPWithWrap(h *provider.HTTPClient) bool {
	return callPlaceOrderBundleUsingBatchHTTP(h, Environment.publicKey, 100000)
}

func callPlaceOrderBundleUsingBatchHTTP(h *provider.HTTPClient, ownerAddr string, bundleTip uint64) bool {
	log.Info("starting placing order with bundle, using a raydium swap")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := &pb.PostRaydiumSwapRequest{
		OwnerAddress: ownerAddr,
		InToken:      "SOL",
		OutToken:     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		Slippage:     0.2,
		InAmount:     0.01,
		Tip:          &bundleTip,
	}

	resp, err := h.PostRaydiumSwap(ctx, request)
	if err != nil {
		log.Error(fmt.Errorf("failed to post raydium swap: %w", err))
		return true
	}

	signedTx, err := transaction.SignTx(resp.Transactions[0].Content)
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

	batchResp, err := h.PostSubmitBatchV2(ctx, &batchRequest)
	if err != nil {
		panic(err)
	}

	log.Infof("successfully placed bundle batch order with signature : %s", batchResp.Transactions[0].Signature)

	return false
}

func callPlaceOrderBundleWrap(h *provider.HTTPClient) bool {
	return callPlaceOrderBundle(h, Environment.publicKey, 100000)
}

func callPlaceOrderBundle(h *provider.HTTPClient, ownerAddr string, bundleTip uint64) bool {
	log.Info("starting placing order with bundle, using a raydium swap")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := &pb.PostRaydiumSwapRequest{
		OwnerAddress: ownerAddr,
		InToken:      "SOL",
		OutToken:     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		Slippage:     0.2,
		InAmount:     0.01,
		Tip:          &bundleTip,
	}

	resp, err := h.PostRaydiumSwap(ctx, request)
	if err != nil {
		log.Error(fmt.Errorf("failed to post raydium swap: %w", err))
		return true
	}

	tx, err := h.SignAndSubmit(ctx, &pb.TransactionMessage{Content: resp.Transactions[0].Content, IsCleanup: false},
		true,
		true, false)
	if err != nil {
		panic(err)
	}

	log.Infof("successfully placed bundle batch order with signature : %s", tx)

	return false
}

func callPlaceOrderWithStakedRPCsHTTPWrap(h *provider.HTTPClient) bool {
	return callPlaceOrderWithStakedRPCsHTTP(h, Environment.publicKey, 100000)
}

func callPlaceOrderWithStakedRPCsHTTP(h *provider.HTTPClient, ownerAddr string, bundleTip uint64) bool {
	log.Info("starting placing raydium swap with staked rpcs")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := &pb.PostRaydiumSwapRequest{
		OwnerAddress: ownerAddr,
		InToken:      "SOL",
		OutToken:     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		Slippage:     0.5,
		InAmount:     0.01,
		Tip:          &bundleTip,
	}

	resp, err := h.PostRaydiumSwap(ctx, request)
	if err != nil {
		log.Error(fmt.Errorf("failed to post raydium swap: %w", err))
		return true
	}

	tx, err := h.SignAndSubmit(ctx, &pb.TransactionMessage{Content: resp.Transactions[0].Content, IsCleanup: false},
		true,
		false, true)
	if err != nil {
		panic(err)
	}

	log.Infof("successfully placed raydium swap using staked rpcs with signature : %s", tx)

	return false
}

func callCancelByClientOrderIDHTTP(h *provider.HTTPClient, ownerAddr, ooAddr string, clientOrderID uint64) bool {
	time.Sleep(60 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := h.SubmitCancelOrderV2(ctx, "", clientOrderID, sideAsk, ownerAddr,
		marketAddr, ooAddr, provider.SubmitOpts{
			SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
			SkipPreFlight:  config.BoolPtr(false),
		})
	if err != nil {
		log.Errorf("failed to cancel order by client ID (%v)", err)
		return true
	}

	log.Infof("canceled order for clientOrderID %v", clientOrderID)
	return false
}

func callPostSettleHTTPWrap(h *provider.HTTPClient) bool {
	return callPostSettleHTTP(h, Environment.publicKey, Environment.openOrdersAddress)
}

func callPostSettleHTTP(h *provider.HTTPClient, ownerAddr, ooAddr string) bool {
	time.Sleep(60 * time.Second)
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

func cancelAllWrap(h *provider.HTTPClient) bool {
	return cancelAll(h, Environment.publicKey, Environment.payer, Environment.openOrdersAddress, sideAsk, typeLimit)
}

func cancelAll(h *provider.HTTPClient, ownerAddr, payerAddr, ooAddr string, orderSide string, orderType string) bool {
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
		SkipPreFlight:  config.BoolPtr(true),
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
	callPostSettleHTTP(h, ownerAddr, ooAddr)
	return false
}

func callReplaceByClientOrderIDWrap(h *provider.HTTPClient) bool {
	return callReplaceByClientOrderID(h, Environment.publicKey, Environment.payer, Environment.openOrdersAddress, sideAsk, typeLimit)
}

func callReplaceByClientOrderID(h *provider.HTTPClient, ownerAddr, payerAddr, ooAddr string, orderSide string, orderType string) bool {
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

func callReplaceOrderWrap(h *provider.HTTPClient) bool {
	return callReplaceOrder(h, Environment.publicKey, Environment.payer, Environment.openOrdersAddress, sideAsk, typeLimit)
}

func callReplaceOrder(h *provider.HTTPClient, ownerAddr, payerAddr, ooAddr string, orderSide string, orderType string) bool {
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

func callGetRecentBlockHashHTTP(h *provider.HTTPClient) bool {
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

func callTradeSwapWrap(h *provider.HTTPClient) bool {
	return callTradeSwap(h, Environment.publicKey)
}

func callTradeSwap(h *provider.HTTPClient, ownerAddr string) bool {
	log.Info("starting trade swap test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("trade swap")
	sig, err := h.SubmitTradeSwap(ctx, ownerAddr, "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", "So11111111111111111111111111111111111111112",
		0.01, 0.1, pb.Project_P_RAYDIUM, provider.SubmitOpts{
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

func callRaydiumSwapWrap(h *provider.HTTPClient) bool {
	return callRaydiumSwap(h, Environment.publicKey)
}

func callRaydiumSwap(h *provider.HTTPClient, ownerAddr string) bool {
	log.Info("starting Raydium swap test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("Raydium swap")
	sig, err := h.SubmitRaydiumSwap(ctx, &pb.PostRaydiumSwapRequest{
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

func callPostPumpFunSwapWrap(h *provider.HTTPClient) bool {
	return callPostPumpFunSwap(h, Environment.publicKey)
}

func callPostPumpFunSwap(h *provider.HTTPClient, ownerAddr string) bool {
	log.Info("starting PostPumpFunSwap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	h = provider.NewHTTPClientPumpNY()

	log.Info("PumpFun swap")
	sig, err := h.SubmitPostPumpFunSwap(ctx, &pb.PostPumpFunSwapRequest{
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

func callRaydiumRouteSwapWrap(h *provider.HTTPClient) bool {
	return callRaydiumRouteSwap(h, Environment.publicKey)
}

func callRaydiumRouteSwap(h *provider.HTTPClient, ownerAddr string) bool {
	log.Info("starting Raydium route swap test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("Raydium route swap")
	sig, err := h.SubmitRaydiumRouteSwap(ctx, &pb.PostRaydiumRouteSwapRequest{
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

func callJupiterRouteSwapWrap(h *provider.HTTPClient) bool {
	return callJupiterRouteSwap(h, Environment.publicKey)
}

func callJupiterRouteSwap(h *provider.HTTPClient, ownerAddr string) bool {
	log.Info("starting Jupiter route swap test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("Jupiter route swap")
	sig, err := h.SubmitJupiterRouteSwap(ctx, &pb.PostJupiterRouteSwapRequest{
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

func callJupiterSwapWrap(h *provider.HTTPClient) bool {
	return callJupiterSwap(h, Environment.publicKey)
}

func callJupiterSwap(h *provider.HTTPClient, ownerAddr string) bool {
	log.Info("starting Jupiter swap test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("Jupiter swap")
	sig, err := h.SubmitJupiterSwap(ctx, &pb.PostJupiterSwapRequest{
		OwnerAddress: ownerAddr,
		InToken:      "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		OutToken:     "So11111111111111111111111111111111111111112",
		Slippage:     0.4,
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

func callJupiterSwapInstructionsWrap(h *provider.HTTPClient) bool {
	tip := uint64(100000)
	return callJupiterSwapInstructions(h, Environment.publicKey, &tip, true)
}

func callJupiterSwapInstructions(h *provider.HTTPClient, ownerAddr string, tipAmount *uint64, useBundle bool) bool {
	log.Info("starting Jupiter swap test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("Jupiter swap")
	sig, err := h.SubmitJupiterSwapInstructions(ctx, &pb.PostJupiterSwapInstructionsRequest{
		OwnerAddress: ownerAddr,
		InToken:      "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		OutToken:     "So11111111111111111111111111111111111111112",
		Slippage:     0.4,
		InAmount:     0.01,
		Tip:          tipAmount,
	}, useBundle, provider.SubmitOpts{
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

func callRaydiumSwapInstructionsWrap(h *provider.HTTPClient) bool {
	tip := uint64(100000)
	return callRaydiumSwapInstructions(h, Environment.publicKey, &tip, true)
}

func callRaydiumSwapInstructions(h *provider.HTTPClient, ownerAddr string, tipAmount *uint64, useBundle bool) bool {
	log.Info("starting Raydium swap test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("Raydium swap")
	sig, err := h.SubmitRaydiumSwapInstructions(ctx, &pb.PostRaydiumSwapInstructionsRequest{
		OwnerAddress: ownerAddr,
		InToken:      "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		OutToken:     "So11111111111111111111111111111111111111112",
		Slippage:     0.4,
		InAmount:     0.01,
		Tip:          tipAmount,
	}, useBundle, provider.SubmitOpts{
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

func callRouteTradeSwapWrap(h *provider.HTTPClient) bool {
	return callRouteTradeSwap(h, Environment.publicKey)
}

func callRouteTradeSwap(h *provider.HTTPClient, ownerAddr string) bool {
	log.Info("starting route trade swap test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("route trade swap")
	sig, err := h.SubmitRouteTradeSwap(ctx, &pb.RouteTradeSwapRequest{
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

func callGetPriorityFeeHTTP(h *provider.HTTPClient) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pf, err := h.GetPriorityFee(ctx, pb.Project_P_RAYDIUM, nil)
	if err != nil {
		log.Errorf("error with GetPriorityFee request: %v", err)
		return true
	}

	log.Infof("priority fee: %v", pf)
	return false
}
