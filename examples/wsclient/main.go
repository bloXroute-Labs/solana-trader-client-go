package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/manifoldco/promptui"

	"github.com/bloXroute-Labs/solana-trader-client-go/examples"
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

		var names []string
		for name := range ExampleEndpoints {
			names = append(names, name)
		}

		sort.Strings(names)

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
	"getTransaction": {
		run:         callGetTransactionWS,
		description: "get tickers",
	},
	"postSubmit": {
		run:         callPostSubmit,
		description: "get tickers",
	},
	"postSubmitWithBatch": {
		run:         callPosSubmitWithBatch,
		description: "get tickers",
	},
	"getRateLimit": {
		run:         callGetRateLimitWS,
		description: "get rate limit",
	},
	"getJupiterPrices": {
		run:         callJupiterPricesWS,
		description: "get jupiter prices",
	},
	"getAccountBalance": {
		run:         callAccountBalanceWS,
		description: "get account balance",
	},

	"getRecentBlockhash": {
		run:         callGetRecentBlockHashWS,
		description: "get quotes",
	},
	"getRecentBlockHashV2": {
		run:         callGetRecentBlockHashV2WSWrap,
		description: "get quotes",
	},

	"getPumpFunQuotes": {
		run:         callGetPumpFunQuotes,
		description: "get pump fun quotes",
	},

	"getPumpFunAmmQuotes": {
		run:         callGetPumpFunAmmQuotes,
		description: "get pump fun AMM quotes",
	},

	"getJupiterQuotes": {
		run:         callGetJupiterQuotes,
		description: "get jupiter quotes",
	},

	"blockStream": {
		run:         callBlockWSStream,
		description: "block stream",
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
	"getPumpFunAmmSwapStream": {
		run:         callGetPumpFunAmmSwapWSStream,
		description: "get pump swap swaps stream",
	},
	"getPumpFunNewAmmPoolStream": {
		run:         callGetPumpFunNewAmmPoolWSStream,
		description: "get new amm pools on pump swap",
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

	"pumpFunAmmSwap": {
		run:                               callPostPumpFunAmmSwap,
		description:                       "pump fun AMM swap",
		requiresAdditionalEnvironmentVars: true,
	},

	"jupiterRouteSwap": {
		run:                               callJupiterRouteSwapWrap,
		description:                       "call jupiter route swap",
		requiresAdditionalEnvironmentVars: true,
	},

	"jupiterSwapWithInstructions": {
		run:                               callJupiterSwapInstructionsWrap,
		description:                       "call jupiter swap with instructions",
		requiresAdditionalEnvironmentVars: true,
	},

	"runAllExamples": {},
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

func callPostSubmit(w provider.WSClientTraderAPI) bool {
	log.Info("starting place order with bundle")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := w.GetRecentBlockHash(ctx)
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
		system.NewTransferInstruction(10000000, privateKey.PublicKey(), solana.MustPublicKeyFromBase58("HWEoBxYs7ssKuudEjzjmpfJVX7Dvi7wescFsVx2L5yoY")).Build(),
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

	resp, err := w.SignAndSubmit(ctx, &pb.TransactionMessage{
		Content: tx1.MustToBase64()}, false, false, true)
	if err != nil {
		log.Errorf("failed to sign and submit order (%v)", err)
		return true
	}

	log.Infof("submitted bundle order to trader api %v", resp)

	return false
}

func callPosSubmitWithBatch(w provider.WSClientTraderAPI) bool {
	log.Info("starting to place order with bundle")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := w.GetRecentBlockHash(ctx)
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
		system.NewTransferInstruction(1000000, privateKey.PublicKey(), solana.MustPublicKeyFromBase58("FZwLKcQupnTy2CbaVMGGsutxDtjv9CqYVDJxiNZSj5Xi")).Build(),
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
	tx2, err := solana.NewTransaction([]solana.Instruction{
		priceLimitIx,
		system.NewTransferInstruction(100000, privateKey.PublicKey(), privateKey.PublicKey()).Build(),
	}, bh, solana.TransactionPayer(privateKey.PublicKey()))
	if err != nil {
		return false
	}

	tx2.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(privateKey.PublicKey()) {
			return &privateKey
		}
		return &wlt.PrivateKey
	})

	frp := true

	batchEntry := pb.PostSubmitRequestEntry{
		Transaction:   &pb.TransactionMessage{Content: tx2.MustToBase64()},
		SkipPreFlight: true,
	}
	batchEntry2 := pb.PostSubmitRequestEntry{
		Transaction:   &pb.TransactionMessage{Content: tx1.MustToBase64()},
		SkipPreFlight: true,
	}

	batchRequest := pb.PostSubmitBatchRequest{
		Entries:                []*pb.PostSubmitRequestEntry{&batchEntry, &batchEntry2},
		FrontRunningProtection: &frp,
	}

	batchResp, err := w.PostSubmitBatchV2(ctx, &batchRequest)
	if err != nil {
		panic(err)
	}

	log.Infof("successfully placed bundle batch order with signature : %s", batchResp.Transactions[0].Signature)

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

func callAccountBalanceWS(w provider.WSClientTraderAPI) bool {
	log.Info("fetching balances...")

	response, err := w.GetAccountBalance(context.Background(), Environment.PublicKey)
	if err != nil {
		log.Errorf("error with GetAccountBalance request for AFT8VayE7qr8MoQsW3wHsDS83HhEvhGWdbNSHRKeUDfQ: %v", err)
		return true
	} else {
		log.Info(response)
	}

	fmt.Println()
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

func callGetPumpFunAmmQuotes(_ provider.WSClientTraderAPI) bool {
	w, err := provider.NewWSClientPumpNY(provider.MainnetPumpNYWS)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	quotes, err := w.GetPumpFunAmmQuotes(ctx, &pb.GetPumpFunAmmQuotesRequest{
		InToken:  "So11111111111111111111111111111111111111112",
		InAmount: 0.01,
		OutToken: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		Pool:     "Gf7sXMoP8iRw4iiXmJ1nq4vxcRycbGXy5RL8a8LnTd3v",
		Slippage: 1,
	})
	if err != nil {
		return true
	}

	log.Infof("best quote for PumpFun AMM is %v", quotes)

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

func callPostPumpFunSwapWrap(w provider.WSClientTraderAPI) bool {
	return callPostPumpFunSwap(w, Environment.PublicKey)
}

func callPostPumpFunSwap(w provider.WSClientTraderAPI, ownerAddr string) bool {
	log.Info("starting PostPumpFunSwap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wp, err := provider.NewWSClientPumpNY(provider.MainnetPumpNYWS)
	if err != nil {
		panic(err)
	}

	newToken, err := examples.GetPumpFunNewTokenHelper()
	if err != nil {
		panic(err)
	}

	log.Info("PumpFun swap")
	sig, err := wp.SubmitPostPumpFunSwap(ctx, &pb.PostPumpFunSwapRequest{
		UserAddress:         ownerAddr,
		BondingCurveAddress: newToken.BondingCurve,
		TokenAddress:        newToken.Mint,
		Creator:             newToken.Creator,
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

func callPostPumpFunAmmSwap(_ provider.WSClientTraderAPI) bool {
	log.Info("starting PostPumpFunAmmSwap test")
	w, err := provider.NewWSClientPumpNY(provider.MainnetPumpNYWS)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("PumpFun AMM swap")
	tip := uint64(10000)
	sig, err := w.SubmitPostPumpFunAmmSwap(ctx, &pb.PostPumpFunAmmSwapRequest{
		OwnerAddress: Environment.PublicKey,
		InToken:      "So11111111111111111111111111111111111111112",
		InAmount:     0.01,
		OutToken:     "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", // USDC
		Pool:         "Gf7sXMoP8iRw4iiXmJ1nq4vxcRycbGXy5RL8a8LnTd3v",
		Slippage:     1,
		ComputeLimit: 130000,
		ComputePrice: uint64(0.1 * 1000000),
		Tip:          &tip,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("PumpFun AMM swap transaction signature : %s", sig)
	return false
}

func callJupiterSwapWrap(w provider.WSClientTraderAPI) bool {
	return callJupiterSwap(w, Environment.PublicKey)
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
		SkipPreFlight: false,
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
		SkipPreFlight: false,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Jupiter swap transaction signature : %s", sig)
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
		SkipPreFlight: false,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Jupiter route swap transaction signature : %s", sig)
	return false
}

func callGetPumpFunNewTokensWSStreamWrap(_ provider.WSClientTraderAPI) bool {
	ww, err := provider.NewWSClientPumpNY(provider.MainnetPumpNYWS)
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

func callGetPumpFunAmmSwapWSStream(_ provider.WSClientTraderAPI) bool {
	log.Info("starting GetPumpFunAMMSwap stream")
	wp, err := provider.NewWSClientPumpNY(provider.MainnetPumpNYWS)
	if err != nil {
		log.Errorf("failed to create pump fun provider: %v", err)
		return true
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := wp.GetPumpFunAmmSwapStream(ctx, &pb.GetPumpFunAMMSwapStreamRequest{
		Pools: []string{"Gj5t6KjTw3gWW7SrMHEi1ojCkaYHyvLwb17gktf96HNH"},
	})
	if err != nil {
		log.Errorf("error with GetPumpFunAMMSwap stream request: %v", err)
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

func callGetPumpFunNewAmmPoolWSStream(w provider.WSClientTraderAPI) bool {
	log.Info("starting GetPumpFunNewAmmPool stream")

	wp, err := provider.NewWSClientPumpNY(provider.MainnetPumpNYWS)
	if err != nil {
		log.Errorf("failed to create pump fun provider: %v", err)
		return true
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := wp.GetPumpFunNewAmmPoolStream(ctx, &pb.GetPumpFunNewAmmPoolStreamRequest{})
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

	result, err := w.GetRecentBlockHash(context.Background())
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

	result, err := w.GetRecentBlockHashV2(context.Background(), offset)
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

	result, err := w.GetRecentBlockHashV2(context.Background(), 0)
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
