package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	"github.com/gagliardetto/solana-go"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/manifoldco/promptui"

	"github.com/bloXroute-Labs/solana-trader-client-go/examples"
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
		Items: []string{"ny", "uk", "fr"},
	}

	_, region, err := regionPrompt.Run()
	if err != nil {
		panic(fmt.Errorf("prompt failed: %v", err))
	}

	for {
		client := setupHTTPClient(config.Env(environment), config.HTTPUrls[config.Region(region)])

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
	"getTransaction": {
		run:         callGetTransactionHTTP,
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
		run:         callGetRateLimitHTTP,
		description: "get rate limit",
	},
	"getRecentBlockhash": {
		run:         callGetRecentBlockHashHTTP,
		description: "get recent blockhash",
	},
	"getAccountBalance": {
		run:         callGetAccountBalanceHTTP,
		description: "get account balance",
	},
	"getPumpFunQuotes": {
		run:         callGetPumpFunQuotesHTTP,
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

	"getPriorityFee": {
		run:         callGetPriorityFeeHTTP,
		description: "get priority fee",
	},

	"getPriorityFeeByProgram": {
		run:         callGetPriorityFeeByProgramHTTP,
		description: "get priority fee by program",
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
}

func callGetAccountBalanceHTTP(h provider.HTTPClientTraderAPI) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.GetAccountBalance(ctx, Environment.PublicKey)
	if err != nil {
		log.Errorf("error with GetAccountBalance request for HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc: %v", err)
		return true
	} else {
		log.Info(response)
	}

	fmt.Println()
	return false
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

func callPostSubmit(h provider.HTTPClientTraderAPI) bool {
	log.Info("starting place order with bundle")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := h.GetRecentBlockHash(ctx)
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

	resp, err := h.SignAndSubmit(ctx, &pb.TransactionMessage{
		Content: tx1.MustToBase64()}, false, false, true)
	if err != nil {
		log.Errorf("failed to sign and submit order (%v)", err)
		return true
	}

	log.Infof("submitted bundle order to trader api %v", resp)

	return false
}

func callPosSubmitWithBatch(h provider.HTTPClientTraderAPI) bool {
	log.Info("starting to place order with bundle")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := h.GetRecentBlockHash(ctx)
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

	batchResp, err := h.PostSubmitBatchV2(ctx, &batchRequest)
	if err != nil {
		panic(err)
	}

	log.Infof("successfully placed bundle batch order with signature : %s", batchResp.Transactions[0].Signature)

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

func callGetPumpFunAmmQuotes(_ provider.HTTPClientTraderAPI) bool {
	g := provider.NewHTTPClientPumpNY()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	quotes, err := g.GetPumpFunAmmQuotes(ctx, &pb.GetPumpFunAmmQuotesRequest{
		InToken:  "So11111111111111111111111111111111111111112",
		InAmount: 0.01,
		OutToken: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		Pool:     "Gf7sXMoP8iRw4iiXmJ1nq4vxcRycbGXy5RL8a8LnTd3v",
		Slippage: 1,
	})
	if err != nil {
		return true
	}

	log.Infof("quote for PumpFun AMM is %v", quotes)

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

func callPostPumpFunSwapWrap(h provider.HTTPClientTraderAPI) bool {
	return callPostPumpFunSwap(h, Environment.PublicKey)
}

func callPostPumpFunSwap(h provider.HTTPClientTraderAPI, ownerAddr string) bool {
	log.Info("starting PostPumpFunSwap test")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	h = provider.NewHTTPClientPumpNY()

	newToken, err := examples.GetPumpFunNewTokenHelper()
	if err != nil {
		panic(err)
	}

	log.Info("PumpFun swap")
	sig, err := h.SubmitPostPumpFunSwap(ctx, &pb.PostPumpFunSwapRequest{
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

func callPostPumpFunAmmSwap(_ provider.HTTPClientTraderAPI) bool {
	log.Info("starting PostPumpFunAmmSwap test")
	h := provider.NewHTTPClientPumpNY()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("PumpFun AMM swap")
	tip := uint64(1000000)
	sig, err := h.SubmitPostPumpFunAmmSwap(ctx, &pb.PostPumpFunAmmSwapRequest{
		OwnerAddress: Environment.PublicKey,
		InToken:      "So11111111111111111111111111111111111111112",
		InAmount:     0.0001,
		OutToken:     "FKW2EGRyTiks14gHQH7pTFopPeSMmPXQ9yo4BgyCpump", // USDC
		Pool:         "8jXaxVEiNmZDwqUKkSLAA82dbLBrvgzzu6WH13nA29si",
		Slippage:     8000,
		Tip:          &tip,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("PumpFun AMM swap transaction signature : %s", sig)
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
		SkipPreFlight: false,
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
		SkipPreFlight: false,
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
		SkipPreFlight: false,
	})
	if err != nil {
		log.Error(err)
		return true
	}
	log.Infof("Jupiter swap transaction signature : %s", sig)
	return false
}
func callGetPriorityFeeHTTP(h provider.HTTPClientTraderAPI) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	nf := 66.0
	pf, err := h.GetPriorityFee(ctx, pb.Project_P_RAYDIUM, &nf)
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

	//RaydiumCLMM := "CAMMCzo5YL8w4VFF8KVHrK22GGUsp5VTaW7grrKgrWqK"
	RaydiumCPMM := "CPMMoo8L3F4NbTegBCKVNunggL7H1ZpdTHKxQB5qKP1C"

	pf, err := h.GetPriorityFeeByProgram(ctx, []string{RaydiumCPMM})
	if err != nil {
		log.Errorf("error with GetPriorityFeeByProgram request: %v", err)
		return true
	}

	log.Infof("priority fee by program: %v", pf)
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
