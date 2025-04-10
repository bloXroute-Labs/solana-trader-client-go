package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"

	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	"github.com/manifoldco/promptui"

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
	computePrice = 200000
	computeLimit = 100000
)

var Environment config.EnvironmentVariables

func main() {
	utils.InitLogger()

	Environment = config.InitializeEnvironmentVariables()

	listAllEndpoints()

	envPrompt := promptui.Select{
		Label: "Select environment",
		Items: []string{"mainnet", "testnet", "local"},
	}

	_, environment, err := envPrompt.Run()
	if err != nil {
		panic(fmt.Errorf("prompt failed: %v", err))
	}

	regionPrompt := promptui.Select{
		Label: "Select region",
		Items: []string{"ny", "uk"},
	}

	_, region, err := regionPrompt.Run()
	if err != nil {
		panic(fmt.Errorf("prompt failed: %v", err))
	}

	for {
		client := setupGRPCClient(config.Env(environment), config.GRPCUrls[config.Region(region)])
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
						log.Errorf("%s", fmt.Sprintf("example '%s' failed", exampleName))
						time.Sleep(1 * time.Second)
					}
					time.Sleep(1 * time.Second)
					log.Printf("example '%s' completed successfully\n", exampleName)
				}
			}
		}

	rerun:
		log.Printf("running example: %s\n", exampleName)
		if failed := exampleStruct.run(client); failed {
			log.Errorf("%s", fmt.Sprintf("example '%s' failed", exampleName))
			time.Sleep(1 * time.Second)
		} else {

			time.Sleep(1 * time.Second)
			log.Printf("example '%s' completed successfully\n", exampleName)
		}

		rerunPrompt := promptui.Select{
			Label: "Choose an option",
			Items: []string{"Rerun Function", "Return to Main Menu"},
		}

		_, result, err := rerunPrompt.Run()
		if err != nil {
			fmt.Println("Prompt failed, returning to main menu")
			break
		}

		if result == "Rerun Function" {
			goto rerun
		}

	}

}

func setupGRPCClient(env config.Env, endpoint string) provider.GRPCClientTraderAPI {
	var g provider.GRPCClientTraderAPI
	var err error

	switch env {
	case config.EnvLocal:
		g, err = provider.NewGRPCLocal()
	case config.EnvTestnet:
		g, err = provider.NewGRPCTestnet()
	case config.EnvMainnet:
		g, err = provider.NewGRPCClientFullService(endpoint)
	}
	if err != nil {
		log.Fatalf("error dialing GRPC client: %v", err)
	}

	return g
}

func listAllEndpoints() {
	fmt.Println(fmt.Printf("Available Endpoints (see docs for more info: https://docs.bloxroute.com/solana/trader-api-v2) \n"))

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

type ExampleFunc func(provider.GRPCClientTraderAPI) bool

var ExampleEndpoints = map[string]struct {
	run                               ExampleFunc
	description                       string
	requiresAdditionalEnvironmentVars bool
}{
	"getPools": {
		run:         callPoolsGRPC,
		description: "fetch all available markets",
	},
	"getRaydiumCLMMPools": {
		run:         callRaydiumCLMMPoolsGRPC,
		description: "fetch all available markets",
	},
	"getRaydiumPoolReserve": {
		run:         callRaydiumPoolReserveGRPC,
		description: "get raydium pool reserve",
	},
	"getTransaction": {
		run:         callGetTransactionGRPC,
		description: "get tickers",
	},
	"getRecentBlockHash": {
		run:         callRecentBlockHash,
		description: "get recent blockhash",
	},

	"getRecentBlockHashV2": {
		run:         callGetRecentBlockHashV2GRPCWrap,
		description: "get recent blockhash v2",
	},

	"getRateLimit": {
		run:         callGetRateLimitGRPC,
		description: "get rate limit",
	},
	"getRaydiumPools": {
		run:         callRaydiumPoolsGRPC,
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
	"getPricesStream": {
		run:         callPricesGRPCStream,
		description: "stream prices",
	},
	"getSwapsStream": {
		run:         callSwapsGRPCStream,
		description: "stream swaps",
	},
	"getNewRaydiumPoolStream": {
		run:         callGetNewRaydiumPoolsStream,
		description: "stream new raydium pools",
	},

	"getNewRaydiumPoolByTransactionStream": {
		run:         callGetNewRaydiumPoolsByTransactionStream,
		description: "stream new raydium pools (by transaction updates)",
	},
	"getNewRaydiumPoolsStreamWithCPMM": {
		run:         callGetNewRaydiumPoolsStreamWithCPMM,
		description: "stream new raydium pools with cpmm enabled",
	},

	"getRaydiumQuotes": {
		run:         callGetRaydiumQuotes,
		description: "get raydium quotes",
	},

	"getRaydiumQuotesCPMM": {
		run:         callGetRaydiumCPMMQuotes,
		description: "get raydium cpmm quotes",
	},
	"getRaydiumQuotesCLMM": {
		run:         callGetRaydiumCLMMQuotes,
		description: "get raydium cpmm quotes",
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
	"getPriorityFeeByProgram": {
		run:         callGetPriorityFeeByProgramGRPC,
		description: "get priority fee by program",
	},
	"getGetLeaderSchedule": {
		run:         callGetLeaderScheduleGRPC,
		description: "get leader schedule",
	},
	"getPriorityFeeByProgramStream": {
		run:         callGetPriorityFeeByProgramGRPCStream,
		description: "get priority fee by program stream",
	},
	"getPumpFunNewTokenStream": {
		run:         callGetpumpFunNewTokenGRPCStreamWrap,
		description: "get pump fun new token stream",
	},
	"getPumpFunNewAmmPoolStream": {
		run:         callGetPumpFunNewAmmPoolGRPCStream,
		description: "get new amm pools on pump swap",
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

	"callTestSubmitSnipe": {
		run:                               callTestSubmitSnipeWrap,
		description:                       "submit sniping transactions example",
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

	"raydiumCLMMSwap": {
		run:                               callRaydiumCLMMSwapGRPCWrap,
		description:                       "raydium clmm swap",
		requiresAdditionalEnvironmentVars: true,
	},

	"raydiumCLMMRouteSwap": {
		run:                               callRaydiumCLMMRouteSwapGRPCWrap,
		description:                       "raydium clmm route swap",
		requiresAdditionalEnvironmentVars: true,
	},

	"raydiumCPMMSwap": {
		run:                               callRaydiumCPMMSwapGRPCWrap,
		description:                       "raydium cpmm swap",
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
}

func callGetTokenAccountsGRPCWrap(g provider.GRPCClientTraderAPI) bool {
	return callGetTokenAccountsGRPC(g, Environment.OpenOrdersAddress)

}

func callGetTokenAccountsGRPC(g provider.GRPCClientTraderAPI, ownerAddr string) bool {
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

func callPoolsGRPC(g provider.GRPCClientTraderAPI) bool {
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

func callGetRateLimitGRPC(g provider.GRPCClientTraderAPI) bool {
	tx, err := g.GetRateLimit(context.Background(), &pb.GetRateLimitRequest{})
	if err != nil {
		log.Errorf("error with GetRateLimit request: %v", err)
		return true
	} else {
		log.Info(tx)
	}

	return false
}

func callGetTransactionGRPC(g provider.GRPCClientTraderAPI) bool {
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

func callRaydiumPoolReserveGRPC(g provider.GRPCClientTraderAPI) bool {
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

func callRaydiumPoolsGRPC(g provider.GRPCClientTraderAPI) bool {
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

func callRaydiumCLMMPoolsGRPC(g provider.GRPCClientTraderAPI) bool {
	pools, err := g.GetRaydiumCLMMPools(context.Background(), &pb.GetRaydiumCLMMPoolsRequest{})
	if err != nil {
		log.Errorf("error with callRaydiumCLMMPoolsGRPC request for Raydium: %v", err)
		return true
	} else {
		log.Info(pools)
	}

	fmt.Println()
	return false
}

func callRaydiumPricesGRPC(g provider.GRPCClientTraderAPI) bool {
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

func callJupiterPricesGRPC(g provider.GRPCClientTraderAPI) bool {
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

func callGetRaydiumQuotes(g provider.GRPCClientTraderAPI) bool {
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

func callGetPumpFunQuotes(g provider.GRPCClientTraderAPI) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	amount := 0.01

	quotes, err := g.GetPumpFunQuotes(ctx, &pb.GetPumpFunQuotesRequest{
		QuoteType:           "buy",
		BondingCurveAddress: "Fh8fnZUVEpPStJ2hKFNNjMAyuyvoJLMouENawg4DYCBc",
		MintAddress:         "2DEsbYgW94AtZxgUfYXoL8DqJAorsLrEWZdSfriipump",
		Amount:              amount,
	})
	if err != nil {
		return true
	}

	log.Infof("best quote for PumpFun is %v", quotes)

	fmt.Println()
	return false
}

func callGetRaydiumCLMMQuotes(g provider.GRPCClientTraderAPI) bool {
	log.Info("starting get Raydium CLMMQ quotes test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inToken := "SOL"
	outToken := "USDT"
	amount := 0.01
	slippage := float64(5)

	quotes, err := g.GetRaydiumCLMMQuotes(ctx, &pb.GetRaydiumCLMMQuotesRequest{
		InToken:  inToken,
		OutToken: outToken,
		InAmount: amount,
		Slippage: slippage,
	})
	if err != nil {
		log.Errorf("error with GetRaydiumCLMMQuotes request for %s to %s: %v", inToken, outToken, err)
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

func callGetRaydiumCPMMQuotes(g provider.GRPCClientTraderAPI) bool {
	log.Info("starting get Raydium quotes test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inToken := "So11111111111111111111111111111111111111112"
	outToken := "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
	amount := 0.01
	slippage := float64(5)

	quotes, err := g.GetRaydiumQuotesCPMM(ctx, &pb.GetRaydiumCPMMQuotesRequest{
		InToken:  inToken,
		OutToken: outToken,
		InAmount: amount,
		Slippage: slippage,
	})
	if err != nil {
		log.Errorf("error with GetQuotes request for %s to %s: %v", inToken, outToken, err)
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

func callGetJupiterQuotes(g provider.GRPCClientTraderAPI) bool {
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

func callRecentBlockHashGRPCStream(g provider.GRPCClientTraderAPI) bool {
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

func callRecentBlockHash(g provider.GRPCClientTraderAPI) bool {
	log.Info("starting recent block hash")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// response
	response, err := g.GetRecentBlockHash(ctx)
	if err != nil {
		log.Errorf("error with GetRecentBlockHash request: %v", err)
		return true
	}

	log.Infof("response %v received", response)
	return false
}

func callGetRecentBlockHashV2GRPCWrap(g provider.GRPCClientTraderAPI) bool {
	var failed bool
	for i := 0; i < 2; i++ {
		failed = callRecentBlockHashV2GRPC(g, uint64(i))
	}

	return failed
}

func callRecentBlockHashV2GRPC(g provider.GRPCClientTraderAPI, offset uint64) bool {
	log.Info("starting recent block hash V2")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// response
	response, err := g.GetRecentBlockHashV2(ctx, offset)
	if err != nil {
		log.Errorf("error with GetRecentBlockHashV2 request: %v", err)
		return true
	}

	log.Infof("V2 response %v received", response)
	return false
}

func callPoolReservesGRPCStream(g provider.GRPCClientTraderAPI) bool {
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

func callPlaceOrderBundleWrap(g provider.GRPCClientTraderAPI) bool {
	return callPlaceOrderBundle(g, Environment.PublicKey, Environment.Payer, sideAsk,
		computeLimit, computePrice, typeLimit, 100000)
}

func callPlaceOrderBundle(g provider.GRPCClientTraderAPI, ownerAddr, payerAddr,
	orderSide string, computeLimit uint32, computePrice uint64, orderType string, tipAmount uint64) bool {
	log.Info("starting place order with bundle")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := g.GetRecentBlockHash(ctx)
	if err != nil {
		log.Errorf("error with GetRecentBlockHash request: %v", err)
		return true
	}

	bh := solana.MustHashFromBase58(response.BlockHash)

	wlt := solana.NewWallet()
	privateKey, err := transaction.LoadPrivateKeyFromEnv()
	priceLimitIx, err := computebudget.NewSetComputeUnitPriceInstruction(uint64(200000000)).ValidateAndBuild()
	if err != nil {
		return false
	}

	tx1, err := solana.NewTransaction([]solana.Instruction{
		priceLimitIx,
		system.NewTransferInstruction(10000000, privateKey.PublicKey(), solana.MustPublicKeyFromBase58("FZwLKcQupnTy2CbaVMGGsutxDtjv9CqYVDJxiNZSj5Xi")).Build(),
	}, bh, solana.TransactionPayer(privateKey.PublicKey()))
	if err != nil {
		return false
	}

	tx1.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(privateKey.PublicKey()) {
			return &privateKey
		}
		return &wlt.PrivateKey
	})
	tru := true
	resp, err := g.SignAndSubmitPaladin(ctx, &pb.TransactionMessage{
		Content: tx1.MustToBase64()}, &tru)
	if err != nil {
		log.Errorf("failed to sign and submit order (%v)", err)
		return true
	}

	log.Infof("submitted bundle order to trader api %v", resp)

	return false
}

func callRaydiumSwapWrap(g provider.GRPCClientTraderAPI) bool {
	return callRaydiumSwap(g, Environment.PublicKey)
}

func callRaydiumSwap(g provider.GRPCClientTraderAPI, ownerAddr string) bool {
	log.Info("starting Raydium swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tip := uint64(2000000)

	log.Info("Raydium swap")
	sig, err := g.SubmitRaydiumSwap(ctx, &pb.PostRaydiumSwapRequest{
		OwnerAddress: ownerAddr,
		InToken:      "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		OutToken:     "So11111111111111111111111111111111111111112",
		Slippage:     0.3,
		InAmount:     0.01,
		Tip:          &tip,
		ComputePrice: computePrice,
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_ABORT_ON_FIRST_ERROR,
		SkipPreFlight:  config.BoolPtr(false),
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Raydium swap transaction signature : %s, time: %s", sig, time.Now().String())
	return false
}

func callPostPumpFunSwapWrap(g provider.GRPCClientTraderAPI) bool {
	return callPostPumpFunSwap(Environment.PublicKey)
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
		BondingCurveAddress: "Fh8fnZUVEpPStJ2hKFNNjMAyuyvoJLMouENawg4DYCBc",
		TokenAddress:        "2DEsbYgW94AtZxgUfYXoL8DqJAorsLrEWZdSfriipump",
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

func callRaydiumCLMMSwapGRPCWrap(g provider.GRPCClientTraderAPI) bool {
	return callRaydiumCLMMSwapGRPC(g, Environment.PublicKey)
}

func callRaydiumCLMMSwapGRPC(g provider.GRPCClientTraderAPI, ownerAddr string) bool {
	log.Info("starting Raydium CLMM swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tip := uint64(10000)
	log.Info("Raydium swap")
	sig, err := g.SubmitRaydiumCLMMSwap(ctx, &pb.PostRaydiumSwapRequest{
		OwnerAddress: ownerAddr,
		InToken:      "So11111111111111111111111111111111111111112",
		OutToken:     "HDa3zJc12ahykSsBRvgiWzr6WLEByf36yzKKbVvy4gnF",
		Slippage:     0.1,
		InAmount:     0.0089,
		Tip:          &tip,
		ComputePrice: 10000,
		ComputeLimit: 300000,
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_ABORT_ON_FIRST_ERROR,
		SkipPreFlight:  config.BoolPtr(true),
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Raydium swap transaction signature : %s", sig)
	return false
}

func callRaydiumCPMMSwapGRPCWrap(g provider.GRPCClientTraderAPI) bool {
	return callRaydiumCPMMSwapGRPC(g, Environment.PublicKey)
}

func callRaydiumCPMMSwapGRPC(g provider.GRPCClientTraderAPI, ownerAddr string) bool {
	log.Info("starting Raydium swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tip := uint64(2000000)

	log.Info("Raydium swap")
	sig, err := g.SubmitRaydiumSwapCPMM(ctx, &pb.PostRaydiumCPMMSwapRequest{
		OwnerAddress: ownerAddr,
		InToken:      "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		OutToken:     "So11111111111111111111111111111111111111112",
		Slippage:     0.5,
		InAmount:     0.01,
		ComputePrice: computePrice,
		ComputeLimit: computeLimit,
		Tip:          &tip,
	})

	if err != nil {
		log.Error(err)
		return true
	}

	log.Infof("Raydium Swap CPMM transaction signature : %s", sig)
	return false
}

func callJupiterSwapWrap(g provider.GRPCClientTraderAPI) bool {
	return callJupiterSwap(g, Environment.PublicKey)
}

func callJupiterSwap(g provider.GRPCClientTraderAPI, ownerAddr string) bool {
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

func callJupiterSwapInstructionsWrap(g provider.GRPCClientTraderAPI) bool {
	return callJupiterSwapInstructions(g, Environment.PublicKey, uint64(1100), false)
}

func callJupiterSwapInstructions(g provider.GRPCClientTraderAPI, ownerAddr string, tipAmount uint64, useBundle bool) bool {
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

func callRaydiumSwapInstructionsWrap(g provider.GRPCClientTraderAPI) bool {
	return callRaydiumSwapInstructions(g, Environment.PublicKey, uint64(1100), true)
}

func callRaydiumSwapInstructions(g provider.GRPCClientTraderAPI, ownerAddr string, tipAmount uint64, useBundle bool) bool {
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

func callRaydiumRouteSwapWrap(g provider.GRPCClientTraderAPI) bool {
	return callRaydiumRouteSwap(g, Environment.PublicKey)
}

func callRaydiumRouteSwap(g provider.GRPCClientTraderAPI, ownerAddr string) bool {
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

func callRaydiumCLMMRouteSwapGRPCWrap(g provider.GRPCClientTraderAPI) bool {
	return callRaydiumCLMMRouteSwapGRPC(g, Environment.PublicKey)
}

func callRaydiumCLMMRouteSwapGRPC(g provider.GRPCClientTraderAPI, ownerAddr string) bool {
	log.Info("starting Raydium route CLMM swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("Raydium route swap")
	sig, err := g.SubmitRaydiumCLMMRouteSwap(ctx, &pb.PostRaydiumRouteSwapRequest{
		OwnerAddress: ownerAddr,
		Slippage:     0.1,
		Steps: []*pb.RaydiumRouteStep{
			{
				InToken:      "HDa3zJc12ahykSsBRvgiWzr6WLEByf36yzKKbVvy4gnF",
				OutToken:     "So11111111111111111111111111111111111111112",
				InAmount:     0.000303,
				OutAmountMin: 0.0006005,
				OutAmount:    0.00064,
			},
		},
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_ABORT_ON_FIRST_ERROR,
		SkipPreFlight:  config.BoolPtr(true),
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Raydium route swap transaction signature : %s", sig)
	return false
}

func callJupiterRouteSwapWrap(g provider.GRPCClientTraderAPI) bool {
	return callJupiterRouteSwap(g, Environment.PublicKey)
}

func callJupiterRouteSwap(g provider.GRPCClientTraderAPI, ownerAddr string) bool {
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

func callGetpumpFunNewTokenGRPCStreamWrap(g provider.GRPCClientTraderAPI) bool {
	gg, err := provider.NewGRPCClientPumpNY()
	if err != nil {
		log.Fatal(err)
	}

	mint, res := callGetPumpFunNewTokensGRPCStream(gg)
	res = callGetPumpFunSwapsGRPCStream(gg, mint)
	return res
}

func callGetPumpFunNewTokensGRPCStream(g provider.GRPCClientTraderAPI) (string, bool) {
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

func callGetPumpFunNewAmmPoolGRPCStream(g provider.GRPCClientTraderAPI) bool {
	gg, err := provider.NewGRPCClientPumpNY()
	if err != nil {
		log.Errorf("failed to create pump provider: %v", err)
		return true
	}

	log.Info("starting GetPumpFunNewAmmPool stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := gg.GetPumpFunNewAmmPoolStream(ctx, &pb.GetPumpFunNewAmmPoolStreamRequest{})
	if err != nil {
		log.Errorf("error with GetPumpFunNewAmmPool stream request: %v", err)
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

func callGetPumpFunSwapsGRPCStream(g provider.GRPCClientTraderAPI, mint string) bool {
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

func callPricesGRPCStream(g provider.GRPCClientTraderAPI) bool {
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

func callSwapsGRPCStream(g provider.GRPCClientTraderAPI) bool {
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

func callGetNewRaydiumPoolsStream(g provider.GRPCClientTraderAPI) bool {
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

func callGetNewRaydiumPoolsByTransactionStream(g provider.GRPCClientTraderAPI) bool {
	log.Info("starting get new raydium pools stream without cpmm")

	ch := make(chan *pb.GetNewRaydiumPoolsByTransactionResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Stream response
	stream, err := g.GetNewRaydiumPoolsByTransactionStream(ctx, false)
	if err != nil {
		log.Errorf("error with GetNewRaydiumPools stream request: %v", err)
		return true
	}
	stream.Into(ch)
	for i := 1; i <= 3; i++ {
		resp, ok := <-ch
		if !ok {
			// channel closed
			return true
		}

		log.Infof("response %v received", i)

		fmt.Println(resp)
	}
	return false
}

func callGetNewRaydiumPoolsStreamWithCPMM(g provider.GRPCClientTraderAPI) bool {
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

func callBlockGRPCStream(g provider.GRPCClientTraderAPI) bool {
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

func callGetPriorityFeeGRPCStream(g provider.GRPCClientTraderAPI) bool {
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

func callGetPriorityFeeByProgramGRPCStream(g provider.GRPCClientTraderAPI) bool {
	programs := []string{
		"JUP6LkbZbjS1jKKwapdHNy74zcZ3tLUZoi5QNyVTaV4",
		"CAMMCzo5YL8w4VFF8KVHrK22GGUsp5VTaW7grrKgrWqK",
		"CPMMoo8L3F4NbTegBCKVNunggL7H1ZpdTHKxQB5qKP1C",
	}

	log.Info("starting priority fee by program stream")

	ch := make(chan *pb.GetPriorityFeeByProgramResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := g.GetPriorityFeeByProgramStream(ctx, programs)
	if err != nil {
		log.Errorf("error with GetPriorityFeeByProgram stream request: %v", err)
		return true
	}

	stream.Into(ch)
	for i := 1; i <= 1; i++ {
		response, ok := <-ch
		if !ok {
			return true
		}
		log.Infof("response %v received: %+v", i, response)
	}

	return false
}

func callGetPriorityFeeGRPC(g provider.GRPCClientTraderAPI) bool {
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

func callGetPriorityFeeByProgramGRPC(g provider.GRPCClientTraderAPI) bool {
	log.Info("starting priority fee by program test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	RaydiumCLMM := "CAMMCzo5YL8w4VFF8KVHrK22GGUsp5VTaW7grrKgrWqK"
	RaydiumCPMM := "CPMMoo8L3F4NbTegBCKVNunggL7H1ZpdTHKxQB5qKP1C"

	response, err := g.GetPriorityFeeByProgram(ctx, &pb.GetPriorityFeeByProgramRequest{Programs: []string{RaydiumCLMM, RaydiumCPMM}})
	if err != nil {
		log.Errorf("error with GetPriorityFeeByProgram request: %v", err)
		return true
	}
	log.Infof("response received: %v", response)
	return false
}

func callGetLeaderScheduleGRPC(g provider.GRPCClientTraderAPI) bool {
	log.Info("starting get leader schedule test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := g.GetLeaderSchedule(ctx, &pb.GetLeaderScheduleRequest{
		MaxSlots: 0,
	})
	if err != nil {
		log.Errorf("error with GetLeaderSchedule request: %v", err)
		return true
	}
	log.Infof("response received: %v", response)
	return false
}

func callGetBundleTipGRPCStream(g provider.GRPCClientTraderAPI) bool {
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

func callTestSubmitSnipeWrap(g provider.GRPCClientTraderAPI) bool {
	return callTestSubmitSnipe(g, Environment.PublicKey)
}

func callTestSubmitSnipe(g provider.GRPCClientTraderAPI, ownerAddr string) bool {
	ownerKey, err := solana.PublicKeyFromBase58(ownerAddr)
	if err != nil {
		log.Errorf("Please set Public key environment variable: %v", err)
		return true
	}

	log.Info("starting submit snipe test")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := g.GetRecentBlockHash(ctx)
	if err != nil {
		log.Errorf("error with GetRecentBlockHash request: %v", err)
		return true
	}
	blockHash := solana.MustHashFromBase58(response.BlockHash)

	smallTip := uint64(100_000)
	stakedTipThreshold := uint64(1_000_000)
	tipWallet := solana.MustPublicKeyFromBase58("HWEoBxYs7ssKuudEjzjmpfJVX7Dvi7wescFsVx2L5yoY")
	jitoTipWallet := solana.MustPublicKeyFromBase58("96gYZGLnJYVFmbjzopPSU6QiEV5fGqZNyN9nmNhvrZU5")

	transactions := make([]*pb.TransactionMessage, 2)

	// First transfer to jito tip wallet, then to bloxroute.
	tx1, err := solana.NewTransaction([]solana.Instruction{
		system.NewTransferInstruction(
			smallTip,
			ownerKey,
			jitoTipWallet,
		).Build(),
		system.NewTransferInstruction(
			smallTip,
			ownerKey,
			tipWallet,
		).Build(),
	}, blockHash, solana.TransactionPayer(ownerKey))
	if err != nil {
		log.Errorf("failed to create first transaction: %v", err)
		return true
	}

	// Second transfer to bloxroute directly, with a tip big enough to propegate directly as a staked transaction. > 1_000_000
	tx2, err := solana.NewTransaction([]solana.Instruction{
		system.NewTransferInstruction(
			stakedTipThreshold,
			ownerKey,
			tipWallet,
		).Build(),
	}, blockHash, solana.TransactionPayer(ownerKey))
	if err != nil {
		log.Errorf("failed to create second transaction: %v", err)
		return true
	}

	transactions[0] = &pb.TransactionMessage{
		Content:   tx1.MustToBase64(),
		IsCleanup: false,
	}
	transactions[1] = &pb.TransactionMessage{
		Content:   tx2.MustToBase64(),
		IsCleanup: false,
	}

	signatures, err := g.SignAndSubmitSnipe(ctx, transactions, true)
	if err != nil {
		log.Errorf("failed to submit snipe request: %v", err)
		return true
	}

	log.Infof("snipe signatures: %v", signatures)
	return false
}
