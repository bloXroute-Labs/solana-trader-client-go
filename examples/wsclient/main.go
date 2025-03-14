package main

import (
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
	"github.com/bloXroute-Labs/solana-trader-proto/common"

	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
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
		client := setupWSClient(config.Env(environment), config.WSUrls[config.Region(region)])
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

func setupWSClient(env config.Env, endpoint string) provider.WSClientTraderAPI {
	var w provider.WSClientTraderAPI
	var err error

	switch env {
	case config.EnvLocal:
		w, err = provider.NewWSClientLocal()
	case config.EnvTestnet:
		w, err = provider.NewWSClientTestnet()
	case config.EnvMainnet:
		w, err = provider.NewWSClientFullService(endpoint)
	}
	if err != nil {
		log.Fatalf("error dialing HTTP client: %v", err)
	}

	return w
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

type ExampleFunc func(api provider.WSClientTraderAPI) bool

var ExampleEndpoints = map[string]struct {
	run                               ExampleFunc
	description                       string
	requiresAdditionalEnvironmentVars bool
}{
	"getPools": {
		run:         callPoolsWS,
		description: "fetch all available markets",
	},

	"getPoolsCLMM": {
		run:         callRaydiumCLMMPoolsWS,
		description: "fetch all available markets",
	},

	"getRaydiumPoolReserve": {
		run:         callRaydiumPoolReserveWS,
		description: "get raydium pool reserve",
	},

	"getTransaction": {
		run:         callGetTransactionWS,
		description: "get tickers",
	},

	"getRateLimit": {
		run:         callGetRateLimitWS,
		description: "get rate limit",
	},
	"getRaydiumPools": {
		run:         callRaydiumPoolsWS,
		description: "get raydium pools",
	},
	"getPrice": {
		run:         callPriceWS,
		description: "get raydium pools",
	},
	"getRaydiumPrices": {
		run:         callRaydiumPricesWS,
		description: "get raydium prices",
	},
	"getJupiterPrices": {
		run:         callJupiterPricesWS,
		description: "get jupiter prices",
	},
	"getPricesStream": {
		run:         callPricesWSStream,
		description: "stream prices",
	},
	"getSwapsStream": {
		run:         callSwapsWSStream,
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

	"getQuotes": {
		run:         callGetQuotesWS,
		description: "get quotes",
	},
	"getRecentBlockhash": {
		run:         callGetRecentBlockHashWS,
		description: "get quotes",
	},
	"getRecentBlockHashV2": {
		run:         callGetRecentBlockHashV2WSWrap,
		description: "get quotes",
	},

	"getRaydiumQuotes": {
		run:         callGetRaydiumQuotes,
		description: "get raydium quotes",
	},

	"getRaydiumCPMMQuotes": {
		run:         callGetRaydiumCPMMQuotes,
		description: "get raydium quotes",
	},

	"getRaydiumCLMMQuotes": {
		run:         callGetRaydiumCLMMQuotes,
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
		run:         callRecentBlockHashWSStream,
		description: "recent blockhash stream",
	},
	"poolReservesStream": {
		run:         callPoolReservesWSStream,
		description: "recent blockhash stream",
	},
	"blockStream": {
		run:         callBlockWSStream,
		description: "block stream",
	},
	"getPriorityFee": {
		run:         callGetPriorityFeeWS,
		description: "get priority fee",
	},
	"getLeaderSchedule": {
		run:         callGetLeaderScheduleWS,
		description: "get leader schedule",
	},
	"getPriorityFeeStream": {
		run:         callGetPriorityFeeWSStream,
		description: "get priority fee stream",
	},
	"getPriorityFeeByProgram": {
		run:         callGetPriorityFeeByProgramWS,
		description: "get priority fee by program",
	},
	"getPriorityFeeByProgramStream": {
		run:         callGetPriorityFeeByProgramWSStream,
		description: "get priority fee by program stream",
	},
	"getPumpFunNewTokenStream": {
		run:         callGetPumpFunNewTokensWSStreamWrap,
		description: "get pump fun new token stream",
	},

	"getBundleTipStream": {
		run:         callGetBundleTipWSStream,
		description: "get bundle tip stream",
	},

	"getTokenAccounts": {
		run:                               callGetTokenAccountsWSWrap,
		description:                       "get token accounts",
		requiresAdditionalEnvironmentVars: true,
	},

	"callTestSubmitSnipe": {
		run:                               callTestSubmitSnipeWSWrap,
		description:                       "call test submit snipe",
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

	"raydiumCLMMSwap": {
		run:                               callRaydiumCLMMSwapWSWrap,
		description:                       "raydium clmm swap",
		requiresAdditionalEnvironmentVars: true,
	},

	"raydiumCLMMRouteSwap": {
		run:                               callRaydiumCLMMRouteSwapWSWrap,
		description:                       "raydium clmm route swap",
		requiresAdditionalEnvironmentVars: true,
	},

	"raydiumCPMMSwap": {
		run:                               callRaydiumCPMMSwapWSWrap,
		description:                       "raydium cpmm swap",
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

	"runAllExamples": {},
}

func callPoolsWS(w provider.WSClientTraderAPI) bool {
	log.Info("fetching pools...")

	pools, err := w.GetPools(context.Background(), []pb.Project{pb.Project_P_RAYDIUM})
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

func callGetRateLimitWS(w provider.WSClientTraderAPI) bool {
	log.Info("calling callGetRateLimit...")

	tx, err := w.GetRateLimit(context.Background(), &pb.GetRateLimitRequest{})
	if err != nil {
		log.Errorf("error with GetTransaction request: %v", err)
		return true
	} else {
		log.Info(tx)
	}

	fmt.Println()
	return false
}

func callGetTransactionWS(w provider.WSClientTraderAPI) bool {
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

func callRaydiumPoolReserveWS(w provider.WSClientTraderAPI) bool {
	log.Info("fetching raydium pool reserve...")

	pools, err := w.GetRaydiumPoolReserve(context.Background(), &pb.GetRaydiumPoolReserveRequest{
		PairsOrAddresses: []string{
			"HZ1znC9XBasm9AMDhGocd9EHSyH8Pyj1EUdiPb4WnZjo",
			"D8wAxwpH2aKaEGBKfeGdnQbCc2s54NrRvTDXCK98VAeT",
			"DdpuaJgjB2RptGMnfnCZVmC4vkKsMV6ytRa2gggQtCWt",
			"AVs9TA4nWDzfPJE9gGVNJMVhcQy3V9PGazuz33BfG2RA",
			"58oQChx4yWmvKdwLLZzBi4ChoCc2fqCUWBkwMihLYQo2",
		},
	})
	if err != nil {
		log.Errorf("error with GetRaydiumPools request for Raydium: %v", err)
		return true
	} else {
		log.Info(pools)
	}

	fmt.Println()
	return false
}

func callRaydiumPoolsWS(w provider.WSClientTraderAPI) bool {
	log.Info("fetching Raydium pools...")

	pools, err := w.GetRaydiumPools(context.Background(), &pb.GetRaydiumPoolsRequest{})
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

func callRaydiumCLMMPoolsWS(w provider.WSClientTraderAPI) bool {
	log.Info("fetching Raydium CLMM pools...")

	pools, err := w.GetRaydiumCLMMPools(context.Background(), &pb.GetRaydiumCLMMPoolsRequest{})
	if err != nil {
		log.Errorf("error with GetRaydiumCLMMPools request for Raydium: %v", err)
		return true
	} else {
		log.Info(pools)
	}

	fmt.Println()
	return false
}

func callPriceWS(w provider.WSClientTraderAPI) bool {
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

func callRaydiumPricesWS(w provider.WSClientTraderAPI) bool {
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

func callJupiterPricesWS(w provider.WSClientTraderAPI) bool {
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

func callGetTokenAccountsWSWrap(w provider.WSClientTraderAPI) bool {
	return callGetTokenAccountsWS(w, Environment.PublicKey)
}

func callGetTokenAccountsWS(w provider.WSClientTraderAPI, ownerAddr string) bool {
	log.Info("fetching token accounts...")

	response, err := w.GetTokenAccounts(context.Background(), &pb.GetTokenAccountsRequest{OwnerAddress: ownerAddr})
	if err != nil {
		log.Errorf("error with GetTokenAccounts request %v", err)
		return true
	} else {
		log.Info(response)
	}

	fmt.Println()
	return false
}

func callGetQuotesWS(w provider.WSClientTraderAPI) bool {
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

func callGetRaydiumQuotes(w provider.WSClientTraderAPI) bool {
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

func callGetPumpFunQuotes(w provider.WSClientTraderAPI) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	amount := 0.01

	quotes, err := w.GetPumpFunQuotes(ctx, &pb.GetPumpFunQuotesRequest{
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

func callGetRaydiumCLMMQuotes(w provider.WSClientTraderAPI) bool {
	log.Info("fetching Raydium CLMM quotes...")

	inToken := "SOL"
	outToken := "USDT"
	amount := 0.01
	slippage := float64(5)

	quotes, err := w.GetRaydiumCLMMQuotes(context.Background(), &pb.GetRaydiumCLMMQuotesRequest{
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
		log.Errorf("did not get back 1 quote, got %v quotes", len(quotes.Routes))
		return true
	}
	for _, route := range quotes.Routes {
		log.Infof("best route for Raydium is %v", route)
	}

	fmt.Println()
	return false
}

func callGetJupiterQuotes(w provider.WSClientTraderAPI) bool {
	log.Info("fetching Jupiter quotes...")

	inToken := "So11111111111111111111111111111111111111112"
	outToken := "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
	amount := 0.01
	slippage := float64(5)

	quotes, err := w.GetJupiterQuotes(context.Background(), &pb.GetJupiterQuotesRequest{
		InToken:  inToken,
		OutToken: outToken,
		InAmount: amount,
		Slippage: slippage,
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

func callGetRaydiumCPMMQuotes(w provider.WSClientTraderAPI) bool {
	log.Info("fetching Raydium quotes...")

	inToken := "So11111111111111111111111111111111111111112"
	outToken := "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
	amount := 0.01
	slippage := float64(5)

	quotes, err := w.GetRaydiumQuotesCPMM(context.Background(), &pb.GetRaydiumCPMMQuotesRequest{
		InToken:  inToken,
		OutToken: outToken,
		InAmount: amount,
		Slippage: slippage,
	})
	if err != nil {
		log.Errorf("error with GetRaydiumQuotesCPMM request for %s to %s: %v", inToken, outToken, err)
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

func callGetNewRaydiumPoolsStream(w provider.WSClientTraderAPI) bool {
	log.Info("starting get new raydium pools stream without cpmm")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	poolsChan := make(chan *pb.GetNewRaydiumPoolsResponse)

	stream, err := w.GetNewRaydiumPoolsStream(ctx, false)
	if err != nil {
		log.Errorf("error with GetNewRaydiumPoolsStream: %v", err)
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

func callGetNewRaydiumPoolsByTransactionStream(w provider.WSClientTraderAPI) bool {
	log.Info("starting get new raydium pools stream without cpmm")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	poolsChan := make(chan *pb.GetNewRaydiumPoolsByTransactionResponse)

	stream, err := w.GetNewRaydiumPoolsByTransactionStream(ctx)
	if err != nil {
		log.Errorf("error with GetNewRaydiumPoolsByTransactionStream: %v", err)
		return true
	}

	stream.Into(poolsChan)
	for i := 1; i <= 3; i++ {
		_, ok := <-poolsChan
		if !ok {
			return true
		}
		log.Infof("response %v received", i)
	}
	return false
}

func callGetNewRaydiumPoolsStreamWithCPMM(w provider.WSClientTraderAPI) bool {
	log.Info("starting get new raydium pools stream with cpmm")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	poolsChan := make(chan *pb.GetNewRaydiumPoolsResponse)

	stream, err := w.GetNewRaydiumPoolsStream(ctx, true)
	if err != nil {
		log.Errorf("error with GetNewRaydiumPoolsStream with cpmm: %v", err)
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
func callRecentBlockHashWSStream(w provider.WSClientTraderAPI) bool {
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

func callPoolReservesWSStream(w provider.WSClientTraderAPI) bool {
	log.Info("starting pool reserves stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stream, err := w.GetPoolReservesStream(ctx, &pb.GetPoolReservesStreamRequest{
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

const (
	// SOL/USDC market
	marketAddr = "8BnEgHoWFysVcuFFX7QztDmzuH8r5ZFvyP3sYwn1XTh6"

	orderType   = common.OrderType_OT_LIMIT
	orderPrice  = float64(170200)
	orderAmount = float64(0.1)
)

func callRaydiumSwapWrap(w provider.WSClientTraderAPI) bool {
	return callRaydiumSwap(w, Environment.PublicKey)
}

func callRaydiumSwap(w provider.WSClientTraderAPI, ownerAddr string) bool {
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

func callPostPumpFunSwapWrap(w provider.WSClientTraderAPI) bool {
	return callPostPumpFunSwap(w, Environment.PublicKey)
}

func callPostPumpFunSwap(w provider.WSClientTraderAPI, ownerAddr string) bool {
	log.Info("starting PostPumpFunSwap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wp, err := provider.NewWSClientPumpNY()
	if err != nil {
		panic(err)
	}

	log.Info("PumpFun swap")
	sig, err := wp.SubmitPostPumpFunSwap(ctx, &pb.PostPumpFunSwapRequest{
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

func callRaydiumCLMMSwapWSWrap(w provider.WSClientTraderAPI) bool {
	return callRaydiumCLMMSwapWS(w, Environment.PublicKey)
}

func callRaydiumCLMMSwapWS(w provider.WSClientTraderAPI, ownerAddr string) bool {
	log.Info("starting Raydium CLMM swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig, err := w.SubmitRaydiumCLMMSwap(ctx, &pb.PostRaydiumSwapRequest{
		OwnerAddress: ownerAddr,
		InToken:      "USDT",
		OutToken:     "SOL",
		Slippage:     0.1,
		InAmount:     0.01,
	}, provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  config.BoolPtr(true),
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Raydium CLMM swap transaction signature : %s", sig)
	return false
}

func callRaydiumCPMMSwapWSWrap(w provider.WSClientTraderAPI) bool {
	return callRaydiumSwapCPMMWS(w, Environment.PublicKey)
}

func callRaydiumSwapCPMMWS(w provider.WSClientTraderAPI, ownerAddr string) bool {
	log.Info("starting Raydium swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tip := uint64(2000000)

	sig, err := w.SubmitRaydiumSwapCPMM(ctx, &pb.PostRaydiumCPMMSwapRequest{
		OwnerAddress: ownerAddr,
		InToken:      "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		OutToken:     "So11111111111111111111111111111111111111112",
		Slippage:     0.5,
		InAmount:     0.01,
		Tip:          &tip})

	if err != nil {
		log.Error(err)
		return true
	}

	log.Infof("Raydium CPMM swap transaction signature : %s", sig)
	return false
}

func callRaydiumRouteSwapWrap(w provider.WSClientTraderAPI) bool {
	return callRaydiumRouteSwap(w, Environment.PublicKey)
}

func callRaydiumRouteSwap(w provider.WSClientTraderAPI, ownerAddr string) bool {
	log.Info("starting Raydium swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig, err := w.SubmitRaydiumRouteSwap(ctx, &pb.PostRaydiumRouteSwapRequest{
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

func callJupiterSwapWrap(w provider.WSClientTraderAPI) bool {
	return callJupiterSwap(w, Environment.PublicKey)
}

func callRaydiumCLMMRouteSwapWSWrap(w provider.WSClientTraderAPI) bool {
	return callRaydiumCLMMRouteSwapWS(w, Environment.PublicKey)
}

func callRaydiumCLMMRouteSwapWS(w provider.WSClientTraderAPI, ownerAddr string) bool {
	log.Info("starting Raydium CLMM swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig, err := w.SubmitRaydiumCLMMRouteSwap(ctx, &pb.PostRaydiumRouteSwapRequest{
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
		SkipPreFlight:  config.BoolPtr(true),
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Raydium route swap transaction signature : %s", sig)
	return false
}

func callJupiterSwap(w provider.WSClientTraderAPI, ownerAddr string) bool {
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

func callJupiterSwapInstructionsWrap(w provider.WSClientTraderAPI) bool {
	tip := uint64(100000)
	return callJupiterSwapInstructions(w, Environment.PublicKey, &tip, false)
}

func callJupiterSwapInstructions(w provider.WSClientTraderAPI, ownerAddr string, tipAmount *uint64, useBundle bool) bool {
	log.Info("starting Jupiter swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig, err := w.SubmitJupiterSwapInstructions(ctx, &pb.PostJupiterSwapInstructionsRequest{
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

func callRaydiumSwapInstructionsWrap(w provider.WSClientTraderAPI) bool {
	tip := uint64(100000)
	return callRaydiumSwapInstructions(w, Environment.PublicKey, &tip, false)
}

func callRaydiumSwapInstructions(w provider.WSClientTraderAPI, ownerAddr string, tipAmount *uint64, useBundle bool) bool {
	log.Info("starting Raydium swap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig, err := w.SubmitRaydiumSwapInstructions(ctx, &pb.PostRaydiumSwapInstructionsRequest{
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

func callJupiterRouteSwapWrap(w provider.WSClientTraderAPI) bool {
	return callJupiterRouteSwap(w, Environment.PublicKey)
}

func callJupiterRouteSwap(w provider.WSClientTraderAPI, ownerAddr string) bool {
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

func callPricesWSStream(w provider.WSClientTraderAPI) bool {
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

func callGetPumpFunNewTokensWSStreamWrap(_ provider.WSClientTraderAPI) bool {
	ww, err := provider.NewWSClientPumpNY()
	if err != nil {
		panic(err)
	}
	mint, res := callGetPumpFunNewTokensWSStream(ww)

	if !res {
		return false
	}

	return callGetPumpFunSwapsWSStream(ww, mint)
}

func callGetPumpFunNewTokensWSStream(w provider.WSClientTraderAPI) (string, bool) {
	log.Info("starting GetPumpFunNewTokens stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := w.GetPumpFunNewTokensStream(ctx, &pb.GetPumpFunNewTokensStreamRequest{})
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

func callGetPumpFunSwapsWSStream(w provider.WSClientTraderAPI, mint string) bool {
	log.Info("starting GetPumpFunSwaps stream")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := w.GetPumpFunSwapsStream(ctx, &pb.GetPumpFunSwapsStreamRequest{
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

func callSwapsWSStream(w provider.WSClientTraderAPI) bool {
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

func callBlockWSStream(w provider.WSClientTraderAPI) bool {
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

func callGetPriorityFeeWSStream(w provider.WSClientTraderAPI) bool {
	log.Info("starting get priority fee stream")

	ch := make(chan *pb.GetPriorityFeeResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := w.GetPriorityFeeStream(ctx, pb.Project_P_RAYDIUM, nil)
	if err != nil {
		log.Errorf("error with GetPriorityFee stream request: %v", err)
		return true
	}
	stream.Into(ch)
	for i := 1; i <= 1; i++ {
		_, ok := <-ch
		if !ok {
			return true
		}

		log.Infof("response %v received", i)
	}
	return false
}

func callGetPriorityFeeWS(w provider.WSClientTraderAPI) bool {
	log.Info("fetching priority fee...")

	priorityFee, err := w.GetPriorityFee(context.Background(), pb.Project_P_RAYDIUM, nil)
	if err != nil {
		log.Errorf("error with GetPriorityFee request: %v", err)
		return true
	}

	log.Infof("priority fee: %v", priorityFee)
	return false
}

func callGetLeaderScheduleWS(w provider.WSClientTraderAPI) bool {
	log.Info("fetching leader schedule...")

	leaderSchedule, err := w.GetLeaderSchedule(context.Background(), 0)
	if err != nil {
		log.Errorf("error with GetLeaderSchedule request: %v", err)
		return true
	}

	log.Infof("leader schedule: %v", leaderSchedule)
	return false
}

func callGetBundleTipWSStream(w provider.WSClientTraderAPI) bool {
	log.Info("starting get bundle tip stream")

	ch := make(chan *pb.GetBundleTipResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := w.GetBundleTipStream(ctx)
	if err != nil {
		log.Errorf("error with GetBundleTip stream request: %v", err)
		return true
	}

	stream.Into(ch)
	for i := 1; i <= 1; i++ {
		_, ok := <-ch
		if !ok {
			return true
		}

		log.Infof("response %v received", i)
	}
	return false
}

func callGetPriorityFeeByProgramWS(w provider.WSClientTraderAPI) bool {
	log.Info("fetching priority fee by program...")

	RaydiumCLMM := "CAMMCzo5YL8w4VFF8KVHrK22GGUsp5VTaW7grrKgrWqK"
	RaydiumCPMM := "CPMMoo8L3F4NbTegBCKVNunggL7H1ZpdTHKxQB5qKP1C"

	priorityFee, err := w.GetPriorityFeeByProgram(context.Background(), []string{RaydiumCLMM, RaydiumCPMM})
	if err != nil {
		log.Errorf("error with GetPriorityFeeByProgram request: %v", err)
		return true
	}

	log.Infof("priority fee by program: %v", priorityFee)
	return false
}

func callGetPriorityFeeByProgramWSStream(w provider.WSClientTraderAPI) bool {
	log.Info("starting get priority fee by program stream")

	ch := make(chan *pb.GetPriorityFeeByProgramResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	programs := []string{
		"JUP6LkbZbjS1jKKwapdHNy74zcZ3tLUZoi5QNyVTaV4",
		"CAMMCzo5YL8w4VFF8KVHrK22GGUsp5VTaW7grrKgrWqK",
		"CPMMoo8L3F4NbTegBCKVNunggL7H1ZpdTHKxQB5qKP1C",
	}

	stream, err := w.GetPriorityFeeByProgramStream(ctx, programs)
	if err != nil {
		log.Errorf("error with GetPriorityFeeByProgram stream request: %v", err)
		return true
	}
	stream.Into(ch)
	for i := 1; i <= 1; i++ {
		v, ok := <-ch
		if !ok {
			return true
		}

		log.Infof("response %v received", v)
	}
	return false
}

func callGetRecentBlockHashWS(w provider.WSClientTraderAPI) bool {
	log.Info("starting recent block hash")

	result, err := w.GetRecentBlockHash(context.Background(), &pb.GetRecentBlockHashRequest{})
	if err != nil {
		log.Errorf("error with GetRecentBlockHash request: %v", err)
		return true
	}

	log.Infof("response %v received", result)
	return false
}

func callGetRecentBlockHashV2WSWrap(w provider.WSClientTraderAPI) bool {
	var failed bool
	for i := 0; i < 2; i++ {
		failed = callGetRecentBlockHashV2WS(w, uint64(i))
	}

	return failed
}

func callGetRecentBlockHashV2WS(w provider.WSClientTraderAPI, offset uint64) bool {
	log.Info("starting recent block hash V2")

	result, err := w.GetRecentBlockHashV2(context.Background(), &pb.GetRecentBlockHashRequestV2{Offset: offset})
	if err != nil {
		log.Errorf("error with GetRecentBlockHashV2 request: %v", err)
		return true
	}

	log.Infof("response %v received V2", result)
	return false
}

func callTestSubmitSnipeWSWrap(w provider.WSClientTraderAPI) bool {
	return callTestSubmitSnipeWS(w, Environment.PublicKey)
}

func callTestSubmitSnipeWS(w provider.WSClientTraderAPI, ownerAddr string) bool {
	ownerKey, err := solana.PublicKeyFromBase58(ownerAddr)
	if err != nil {
		log.Errorf("Please set Public key environment variable: %v", err)
		return true
	}

	log.Info("starting submit snipe test")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result, err := w.GetRecentBlockHashV2(context.Background(), &pb.GetRecentBlockHashRequestV2{Offset: 0})
	if err != nil {
		log.Errorf("error with GetRecentBlockHashV2 request: %v", err)
		return true
	}

	blockHash := solana.MustHashFromBase58(result.BlockHash)

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
	signatures, err := w.SignAndSubmitSnipe(ctx, transactions, true)
	if err != nil {
		log.Errorf("failed to submit snipe request: %v", err)
		return true
	}

	log.Infof("snipe signatures: %v", signatures)
	return false
}
