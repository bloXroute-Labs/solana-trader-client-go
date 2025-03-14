package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/manifoldco/promptui"

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

		client := setupHTTPClient(config.Env(environment), config.HTTPUrls[config.Region(region)])
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

	rerun:
		log.Printf("running example: %s\n", exampleName)
		if failed := exampleStruct.run(client); failed {
			log.Errorf(fmt.Sprintf("example '%s' failed", exampleName))
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

func setupHTTPClient(env config.Env, endpoint string) provider.HTTPClientTraderAPI {

	var h provider.HTTPClientTraderAPI
	var err error

	switch env {
	case config.EnvLocal:
		h = provider.NewHTTPLocal()
	case config.EnvTestnet:
		h = provider.NewHTTPTestnet()
	case config.EnvMainnet:
		h, err = provider.NewHTTPClientFullService(endpoint)
	}
	if err != nil {
		log.Fatalf("error dialing HTTP client: %v", err)
	}

	return h
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

type ExampleFunc func(api provider.HTTPClientTraderAPI) bool

var ExampleEndpoints = map[string]struct {
	run                               ExampleFunc
	description                       string
	requiresAdditionalEnvironmentVars bool
}{
	"getPools": {
		run:         callPoolsHTTP,
		description: "fetch all available markets",
	},
	"getRaydiumCLMMPools": {
		run:         callRaydiumCLMMPools,
		description: "fetch all available markets",
	},

	"getRaydiumPoolReserve": {
		run:         callRaydiumPoolReserveHTTP,
		description: "get raydium pool reserve",
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

	"getQuotes": {
		run:         callGetQuotesHTTP,
		description: "get quotes",
	},

	"getRaydiumQuotes": {
		run:         callGetRaydiumQuotesHTTP,
		description: "get raydium quotes",
	},

	"getRaydiumCLMMQuotes": {
		run:         callGetRaydiumCLMMQuotes,
		description: "get raydium quotes",
	},

	"getRaydiumCPMMQuotes": {
		run:         callGetRaydiumCPMMQuotes,
		description: "get raydium quotes",
	},

	"getPumpFunQuotes": {
		run:         callGetPumpFunQuotesHTTP,
		description: "get pump fun quotes",
	},

	"getJupiterQuotes": {
		run:         callGetJupiterQuotes,
		description: "get jupiter quotes",
	},

	"getPriorityFee": {
		run:         callGetPriorityFeeHTTP,
		description: "get priority fee",
	},

	"getPriorityFeeByProgram": {
		run:         callGetPriorityFeeByProgramHTTP,
		description: "get priority fee by program",
	},

	"getLeaderSchedule": {
		run:         callGetLeaderScheduleHTTP,
		description: "get leader schedule",
	},

	"callTestSubmitSnipe": {
		run:                               callTestSubmitSnipeHTTPWrap,
		description:                       "test submit snipe",
		requiresAdditionalEnvironmentVars: true,
	},

	"getTokenAccounts": {
		run:                               callGetTokenAccountsHTTPWrap,
		description:                       "get token accounts",
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
	"raydiumCLMMSwap": {
		run:                               callRaydiumCLMMSwapHTTPWrap,
		description:                       "raydium clmm swap",
		requiresAdditionalEnvironmentVars: true,
	},

	"raydiumCLMMRouteSwap": {
		run:                               callRaydiumCLMMRouteSwapWrap,
		description:                       "raydium clmm route swap",
		requiresAdditionalEnvironmentVars: true,
	},

	"raydiumCPMMSwap": {
		run:                               callRaydiumSwapCPMMWrap,
		description:                       "raydium cpmm swap",
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

func callGetTokenAccountsHTTPWrap(h provider.HTTPClientTraderAPI) bool {
	return callGetTokenAccountsHTTP(h, Environment.PublicKey)
}

func callGetTokenAccountsHTTP(h provider.HTTPClientTraderAPI, ownerAddr string) bool {
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

func callPoolsHTTP(h provider.HTTPClientTraderAPI) bool {
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

func callRaydiumPoolReserveHTTP(h provider.HTTPClientTraderAPI) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pools, err := h.GetRaydiumPoolReserve(ctx, &pb.GetRaydiumPoolReserveRequest{
		PairsOrAddresses: []string{"66cxXqzCpFttLCdMBXYykjfCEVQKag8Cv1oB5KEacd5b"},
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

func callRaydiumPoolsHTTP(h provider.HTTPClientTraderAPI) bool {
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

func callGetRateLimitHTTP(h provider.HTTPClientTraderAPI) bool {
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

func callRaydiumCLMMPools(h provider.HTTPClientTraderAPI) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pools, err := h.GetRaydiumCLMMPools(ctx, &pb.GetRaydiumCLMMPoolsRequest{})
	if err != nil {
		log.Errorf("error with GetRaydiumCLMMPools request for Raydium: %v", err)
		return true
	} else {
		log.Info(pools)
	}

	fmt.Println()
	return false
}

func callGetTransactionHTTP(h provider.HTTPClientTraderAPI) bool {
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

func callPriceHTTP(h provider.HTTPClientTraderAPI) bool {
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

func callRaydiumPricesHTTP(h provider.HTTPClientTraderAPI) bool {
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

func callJupiterPricesHTTP(h provider.HTTPClientTraderAPI) bool {
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

func callGetQuotesHTTP(h provider.HTTPClientTraderAPI) bool {
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

func callGetRaydiumQuotesHTTP(h provider.HTTPClientTraderAPI) bool {
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

func callGetPumpFunQuotesHTTP(h provider.HTTPClientTraderAPI) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	amount := 0.01

	quotes, err := h.GetPumpFunQuotes(ctx, &pb.GetPumpFunQuotesRequest{
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

func callGetRaydiumCLMMQuotes(h provider.HTTPClientTraderAPI) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	inToken := "SOL"
	outToken := "USDT"
	amount := 0.01
	slippage := float64(5)

	quotes, err := h.GetRaydiumCLMMQuotes(ctx, &pb.GetRaydiumCLMMQuotesRequest{
		InToken:  inToken,
		OutToken: outToken,
		InAmount: amount,
		Slippage: slippage,
	})
	if err != nil {
		log.Errorf("error with GetRaydiumCLMMQuotes request for %s to %s: %v", inToken, outToken, err)
		return true
	}

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

func callGetRaydiumCPMMQuotes(h provider.HTTPClientTraderAPI) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	inToken := "So11111111111111111111111111111111111111112"
	outToken := "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
	amount := 0.01
	slippage := float64(5)

	quotes, err := h.GetRaydiumQuotesCPMM(ctx, &pb.GetRaydiumCPMMQuotesRequest{
		InToken:  inToken,
		OutToken: outToken,
		InAmount: amount,
		Slippage: slippage,
	})
	if err != nil {
		log.Errorf("error with GetQuotesCPMM request for %s to %s: %v", inToken, outToken, err)
		return true
	}

	if err != nil {
		log.Errorf("error with GetRaydiumQuotesCPMM request for %s to %s: %v", inToken, outToken, err)
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

func callGetJupiterQuotes(h provider.HTTPClientTraderAPI) bool {
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

func callGetRecentBlockHashHTTP(h provider.HTTPClientTraderAPI) bool {
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

func callRaydiumSwapWrap(h provider.HTTPClientTraderAPI) bool {
	return callRaydiumSwap(h, Environment.PublicKey)
}

func callRaydiumSwap(h provider.HTTPClientTraderAPI, ownerAddr string) bool {
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

func callRaydiumCLMMSwapHTTPWrap(h provider.HTTPClientTraderAPI) bool {
	return callRaydiumCLMMSwapHTTP(h, Environment.PublicKey)
}

func callRaydiumCLMMSwapHTTP(h provider.HTTPClientTraderAPI, ownerAddr string) bool {
	log.Info("starting Raydium CLMM swap test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("Raydium CLMM swap")
	sig, err := h.SubmitRaydiumCLMMSwap(ctx, &pb.PostRaydiumSwapRequest{
		OwnerAddress: ownerAddr,
		InToken:      "USDT",
		OutToken:     "SOL",
		Slippage:     0.1,
		InAmount:     0.01,
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

func callPostPumpFunSwapWrap(h provider.HTTPClientTraderAPI) bool {
	return callPostPumpFunSwap(h, Environment.PublicKey)
}

func callPostPumpFunSwap(h provider.HTTPClientTraderAPI, ownerAddr string) bool {
	log.Info("starting PostPumpFunSwap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	h = provider.NewHTTPClientPumpNY()

	log.Info("PumpFun swap")
	sig, err := h.SubmitPostPumpFunSwap(ctx, &pb.PostPumpFunSwapRequest{
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

func callRaydiumSwapCPMMWrap(h provider.HTTPClientTraderAPI) bool {
	return callRaydiumSwapCPMM(h, Environment.PublicKey)
}

func callRaydiumSwapCPMM(h provider.HTTPClientTraderAPI, ownerAddr string) bool {
	log.Info("starting Raydium swap test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tip := uint64(2000000)

	log.Info("Raydium swap")
	sig, err := h.SubmitRaydiumSwapCPMM(ctx, &pb.PostRaydiumCPMMSwapRequest{
		OwnerAddress: ownerAddr,
		InToken:      "So11111111111111111111111111111111111111112",
		OutToken:     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		Slippage:     0.1,
		InAmount:     0.01,
		ComputePrice: computePrice,
		ComputeLimit: computeLimit,
		Tip:          &tip,
	})
	if err != nil {
		log.Error(err)
		return true
	}

	// 		OutToken:     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
	log.Infof("Raydium swap transaction signature : %s", sig)
	return false
}

func callRaydiumRouteSwapWrap(h provider.HTTPClientTraderAPI) bool {
	return callRaydiumRouteSwap(h, Environment.PublicKey)
}

func callRaydiumRouteSwap(h provider.HTTPClientTraderAPI, ownerAddr string) bool {
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

func callRaydiumCLMMRouteSwapWrap(h provider.HTTPClientTraderAPI) bool {
	return callRaydiumCLMMRouteSwap(h, Environment.PublicKey)
}

func callRaydiumCLMMRouteSwap(h provider.HTTPClientTraderAPI, ownerAddr string) bool {
	log.Info("starting Raydium CLMM route swap test")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info("Raydium route swap")
	sig, err := h.SubmitRaydiumCLMMRouteSwap(ctx, &pb.PostRaydiumRouteSwapRequest{
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
		SkipPreFlight:  config.BoolPtr(true),
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Raydium route swap transaction signature : %s", sig)
	return false
}

func callJupiterRouteSwapWrap(h provider.HTTPClientTraderAPI) bool {
	return callJupiterRouteSwap(h, Environment.PublicKey)
}

func callJupiterRouteSwap(h provider.HTTPClientTraderAPI, ownerAddr string) bool {
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

func callJupiterSwapWrap(h provider.HTTPClientTraderAPI) bool {
	return callJupiterSwap(h, Environment.PublicKey)
}

func callJupiterSwap(h provider.HTTPClientTraderAPI, ownerAddr string) bool {
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

func callJupiterSwapInstructionsWrap(h provider.HTTPClientTraderAPI) bool {
	tip := uint64(100000)
	return callJupiterSwapInstructions(h, Environment.PublicKey, &tip, true)
}

func callJupiterSwapInstructions(h provider.HTTPClientTraderAPI, ownerAddr string, tipAmount *uint64, useBundle bool) bool {
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

func callRaydiumSwapInstructionsWrap(h provider.HTTPClientTraderAPI) bool {
	tip := uint64(100000)
	return callRaydiumSwapInstructions(h, Environment.PublicKey, &tip, true)
}

func callRaydiumSwapInstructions(h provider.HTTPClientTraderAPI, ownerAddr string, tipAmount *uint64, useBundle bool) bool {
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

func callGetPriorityFeeHTTP(h provider.HTTPClientTraderAPI) bool {
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

func callGetPriorityFeeByProgramHTTP(h provider.HTTPClientTraderAPI) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	RaydiumCLMM := "CAMMCzo5YL8w4VFF8KVHrK22GGUsp5VTaW7grrKgrWqK"
	RaydiumCPMM := "CPMMoo8L3F4NbTegBCKVNunggL7H1ZpdTHKxQB5qKP1C"

	pf, err := h.GetPriorityFeeByProgram(ctx, []string{RaydiumCLMM, RaydiumCPMM})
	if err != nil {
		log.Errorf("error with GetPriorityFeeByProgram request: %v", err)
		return true
	}

	log.Infof("priority fee by program: %v", pf)
	return false
}

func callGetLeaderScheduleHTTP(h provider.HTTPClientTraderAPI) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	leaderSchedule, err := h.GetLeaderSchedule(ctx, 0)
	if err != nil {
		log.Errorf("error with GetLeaderSchedule request: %v", err)
		return true
	}

	log.Infof("leader schedule: %v", leaderSchedule)
	return false
}

func callTestSubmitSnipeHTTPWrap(h provider.HTTPClientTraderAPI) bool {
	return callTestSubmitSnipeHTTP(h, Environment.PublicKey)
}

func callTestSubmitSnipeHTTP(h provider.HTTPClientTraderAPI, ownerAddr string) bool {
	ownerKey, err := solana.PublicKeyFromBase58(ownerAddr)
	if err != nil {
		log.Errorf("Please set Public key environment variable: %v", err)
		return true
	}

	log.Info("starting submit snipe test")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := h.GetRecentBlockHash(ctx)
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

	// Prepare unsigned transaction messages
	transactions[0] = &pb.TransactionMessage{
		Content:   tx1.MustToBase64(),
		IsCleanup: false,
	}
	transactions[1] = &pb.TransactionMessage{
		Content:   tx2.MustToBase64(),
		IsCleanup: false,
	}

	// Submit snipe request
	signatures, err := h.SignAndSubmitSnipe(ctx, transactions, true)
	if err != nil {
		log.Errorf("failed to submit snipe request: %v", err)
		return true
	}

	log.Infof("snipe signatures: %v", signatures)
	return false
}
