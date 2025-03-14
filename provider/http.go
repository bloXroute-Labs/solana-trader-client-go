package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/bloXroute-Labs/solana-trader-client-go/connections"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/gagliardetto/solana-go"
)

// GetRaydiumCLMMQuotes returns the CLMM quotes on Raydium
func (h *HTTPClient) GetRaydiumCLMMQuotes(ctx context.Context, request *pb.GetRaydiumCLMMQuotesRequest) (*pb.GetRaydiumCLMMQuotesResponse, error) {
	url := fmt.Sprintf("%s/api/v2/raydium/clmm-quotes?inToken=%s&outToken=%s&inAmount=%v&slippage=%v",
		h.baseURL, request.InToken, request.OutToken, request.InAmount, request.Slippage)
	response := new(pb.GetRaydiumCLMMQuotesResponse)
	if err := connections.HTTPGetWithClient[*pb.GetRaydiumCLMMQuotesResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// GetRaydiumCLMMPools returns the CLMM pools on Raydium
func (h *HTTPClient) GetRaydiumCLMMPools(ctx context.Context, request *pb.GetRaydiumCLMMPoolsRequest) (*pb.GetRaydiumCLMMPoolsResponse, error) {
	url := fmt.Sprintf("%s/api/v2/raydium/clmm-pools?pairOrAddress=%s", h.baseURL, request.PairOrAddress)
	pools := new(pb.GetRaydiumCLMMPoolsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetRaydiumCLMMPoolsResponse](ctx, url, h.httpClient, pools, h.authHeader); err != nil {
		return nil, err
	}

	return pools, nil
}

// PostRaydiumCLMMSwap returns a partially signed transaction(s) for submitting a swap request on Raydium
func (h *HTTPClient) PostRaydiumCLMMSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest) (*pb.PostRaydiumSwapResponse, error) {
	url := fmt.Sprintf("%s/api/v2/raydium/clmm-swap", h.baseURL)
	var response pb.PostRaydiumSwapResponse
	err := connections.HTTPPostWithClient[*pb.PostRaydiumSwapResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// PostRaydiumCLMMRouteSwap returns a partially signed transaction(s) for submitting a route swap request on Raydium
func (h *HTTPClient) PostRaydiumCLMMRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest) (*pb.PostRaydiumRouteSwapResponse, error) {
	url := fmt.Sprintf("%s/api/v2/raydium/clmm-route-swap", h.baseURL)
	var response pb.PostRaydiumRouteSwapResponse
	err := connections.HTTPPostWithClient[*pb.PostRaydiumRouteSwapResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetTransaction returns details of a recent transaction
func (h *HTTPClient) GetTransaction(ctx context.Context, request *pb.GetTransactionRequest) (*pb.GetTransactionResponse, error) {
	url := fmt.Sprintf("%s/api/v2/transaction?signature=%s", h.baseURL, request.Signature)
	response := new(pb.GetTransactionResponse)
	if err := connections.HTTPGetWithClient[*pb.GetTransactionResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// GetRateLimit returns details of an account rate-limits
func (h *HTTPClient) GetRateLimit(ctx context.Context, _ *pb.GetRateLimitRequest) (*pb.GetRateLimitResponse, error) {
	url := fmt.Sprintf("%s/api/v2/rate-limit", h.baseURL)
	response := new(pb.GetRateLimitResponse)
	if err := connections.HTTPGetWithClient[*pb.GetRateLimitResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// GetRaydiumPoolReserve returns pools details for a given set of pairs or addresses on Raydium
func (h *HTTPClient) GetRaydiumPoolReserve(ctx context.Context, req *pb.GetRaydiumPoolReserveRequest) (*pb.GetRaydiumPoolReserveResponse, error) {
	pairsOrAddressesArg := convertStrSliceArgument("pairsOrAddresses", true, req.GetPairsOrAddresses())
	url := fmt.Sprintf("%s/api/v2/raydium/pool-reserves%s", h.baseURL, pairsOrAddressesArg)
	pools := new(pb.GetRaydiumPoolReserveResponse)
	if err := connections.HTTPGetWithClient[*pb.GetRaydiumPoolReserveResponse](ctx, url, h.httpClient, pools, h.authHeader); err != nil {
		return nil, err
	}

	return pools, nil
}

// GetRaydiumPools returns pools on Raydium
func (h *HTTPClient) GetRaydiumPools(ctx context.Context, _ *pb.GetRaydiumPoolsRequest) (*pb.GetRaydiumPoolsResponse, error) {
	url := fmt.Sprintf("%s/api/v2/raydium/pools", h.baseURL)
	pools := new(pb.GetRaydiumPoolsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetRaydiumPoolsResponse](ctx, url, h.httpClient, pools, h.authHeader); err != nil {
		return nil, err
	}

	return pools, nil
}

// GetRaydiumQuotes returns the possible amount(s) of outToken for an inToken and the route to achieve it on Raydium
func (h *HTTPClient) GetRaydiumQuotes(ctx context.Context, request *pb.GetRaydiumQuotesRequest) (*pb.GetRaydiumQuotesResponse, error) {
	url := fmt.Sprintf("%s/api/v2/raydium/quotes?inToken=%s&outToken=%s&inAmount=%v&slippage=%v",
		h.baseURL, request.InToken, request.OutToken, request.InAmount, request.Slippage)
	response := new(pb.GetRaydiumQuotesResponse)
	if err := connections.HTTPGetWithClient[*pb.GetRaydiumQuotesResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// GetRaydiumQuotesCPMM returns the possible amount(s) of outToken for an inToken and the route to achieve it on Raydium CPMM Pools
func (h *HTTPClient) GetRaydiumQuotesCPMM(ctx context.Context, request *pb.GetRaydiumCPMMQuotesRequest) (*pb.GetRaydiumCPMMQuotesResponse, error) {
	url := fmt.Sprintf("%s/api/v2/raydium/cpmm-quotes?inToken=%s&outToken=%s&inAmount=%v&slippage=%v",
		h.baseURL, request.InToken, request.OutToken, request.InAmount, request.Slippage)
	response := new(pb.GetRaydiumCPMMQuotesResponse)
	if err := connections.HTTPGetWithClient[*pb.GetRaydiumCPMMQuotesResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// GetPumpFunQuotes returns the best quotes for swapping a token on PumpFun platform
func (h *HTTPClient) GetPumpFunQuotes(ctx context.Context, request *pb.GetPumpFunQuotesRequest) (*pb.GetPumpFunQuotesResponse, error) {
	url := fmt.Sprintf("%s/api/v2/pumpfun/quotes?mintAddress=%s&quoteType=%s&amount=%f&bondingCurveAddress=%s",
		h.baseURL, request.MintAddress, request.QuoteType, request.Amount, request.BondingCurveAddress)
	response := new(pb.GetPumpFunQuotesResponse)
	if err := connections.HTTPGetWithClient[*pb.GetPumpFunQuotesResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// GetRaydiumPrices returns the USDC price of requested tokens on Raydium
func (h *HTTPClient) GetRaydiumPrices(ctx context.Context, request *pb.GetRaydiumPricesRequest) (*pb.GetRaydiumPricesResponse, error) {
	tokensArg := convertStrSliceArgument("tokens", true, request.Tokens)
	url := fmt.Sprintf("%s/api/v2/raydium/prices%s", h.baseURL, tokensArg)
	respons := new(pb.GetRaydiumPricesResponse)
	if err := connections.HTTPGetWithClient[*pb.GetRaydiumPricesResponse](ctx, url, h.httpClient, respons, h.authHeader); err != nil {
		return nil, err
	}

	return respons, nil
}

// SubmitRaydiumCLMMSwap builds a Raydium CLMM Swap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitRaydiumCLMMSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := h.PostRaydiumCLMMSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return h.SignAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// SubmitRaydiumCLMMRouteSwap builds a Raydium CLMM RouteSwap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitRaydiumCLMMRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := h.PostRaydiumCLMMRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return h.SignAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// PostRaydiumSwap returns a partially signed transaction(s) for submitting a swap request on Raydium
func (h *HTTPClient) PostRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest) (*pb.PostRaydiumSwapResponse, error) {
	url := fmt.Sprintf("%s/api/v2/raydium/swap", h.baseURL)
	var response pb.PostRaydiumSwapResponse
	err := connections.HTTPPostWithClient[*pb.PostRaydiumSwapResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// PostRaydiumSwapCPMM returns a partially signed transaction(s) for submitting a swap request on Raydium
func (h *HTTPClient) PostRaydiumCPMMSwap(ctx context.Context, request *pb.PostRaydiumCPMMSwapRequest) (*pb.PostRaydiumCPMMSwapResponse, error) {
	url := fmt.Sprintf("%s/api/v2/raydium/cpmm-swap", h.baseURL)
	var response pb.PostRaydiumCPMMSwapResponse
	err := connections.HTTPPostWithClient[*pb.PostRaydiumCPMMSwapResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// PostPumpFunSwap returns a partially signed transaction(s) for submitting a swap request on Pumpdotfun platform
func (h *HTTPClient) PostPumpFunSwap(ctx context.Context, request *pb.PostPumpFunSwapRequest) (*pb.PostPumpFunSwapResponse, error) {
	url := fmt.Sprintf("%s/api/v2/pumpfun/swap", h.baseURL)
	var response pb.PostPumpFunSwapResponse
	err := connections.HTTPPostWithClient[*pb.PostPumpFunSwapResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// PostRaydiumRouteSwap returns a partially signed transaction(s) for submitting a swap request on Raydium
func (h *HTTPClient) PostRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest) (*pb.PostRaydiumRouteSwapResponse, error) {
	url := fmt.Sprintf("%s/api/v2/raydium/route-swap", h.baseURL)
	var response pb.PostRaydiumRouteSwapResponse
	err := connections.HTTPPostWithClient[*pb.PostRaydiumRouteSwapResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetJupiterQuotes returns the possible amount(s) of outToken for an inToken and the route to achieve it on Jupiter
func (h *HTTPClient) GetJupiterQuotes(ctx context.Context, request *pb.GetJupiterQuotesRequest) (*pb.GetJupiterQuotesResponse, error) {
	url := fmt.Sprintf("%s/api/v2/jupiter/quotes?inToken=%s&outToken=%s&inAmount=%v&slippage=%v", h.baseURL, request.InToken, request.OutToken, request.InAmount, request.Slippage)
	response := new(pb.GetJupiterQuotesResponse)
	if err := connections.HTTPGetWithClient[*pb.GetJupiterQuotesResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// GetJupiterPrices returns the USDC price of requested tokens on Jupiter
func (h *HTTPClient) GetJupiterPrices(ctx context.Context, request *pb.GetJupiterPricesRequest) (*pb.GetJupiterPricesResponse, error) {
	tokensArg := convertStrSliceArgument("tokens", true, request.Tokens)
	url := fmt.Sprintf("%s/api/v2/jupiter/prices%s", h.baseURL, tokensArg)
	response := new(pb.GetJupiterPricesResponse)
	if err := connections.HTTPGetWithClient[*pb.GetJupiterPricesResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// PostJupiterSwap returns a partially signed transaction(s) for submitting a swap request on Jupiter
func (h *HTTPClient) PostJupiterSwap(ctx context.Context, request *pb.PostJupiterSwapRequest) (*pb.PostJupiterSwapResponse, error) {
	url := fmt.Sprintf("%s/api/v2/jupiter/swap", h.baseURL)
	var response pb.PostJupiterSwapResponse
	err := connections.HTTPPostWithClient[*pb.PostJupiterSwapResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// PostJupiterSwapInstructions returns a list of instructions that can be used to construct a custom transaction for a jupiter swap
func (h *HTTPClient) PostJupiterSwapInstructions(ctx context.Context, request *pb.PostJupiterSwapInstructionsRequest) (*pb.PostJupiterSwapInstructionsResponse, error) {
	url := fmt.Sprintf("%s/api/v2/jupiter/swap-instructions", h.baseURL)
	var response pb.PostJupiterSwapInstructionsResponse
	err := connections.HTTPPostWithClient[*pb.PostJupiterSwapInstructionsResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// PostRaydiumSwapInstructions returns a list of instructions that can be used to construct a custom transaction for a raydium swap
func (h *HTTPClient) PostRaydiumSwapInstructions(ctx context.Context, request *pb.PostRaydiumSwapInstructionsRequest) (*pb.PostRaydiumSwapInstructionsResponse, error) {
	url := fmt.Sprintf("%s/api/v2/raydium/swap-instructions", h.baseURL)
	var response pb.PostRaydiumSwapInstructionsResponse
	err := connections.HTTPPostWithClient[*pb.PostRaydiumSwapInstructionsResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// PostJupiterRouteSwap returns a partially signed transaction(s) for submitting a swap request on Jupiter
func (h *HTTPClient) PostJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest) (*pb.PostJupiterRouteSwapResponse, error) {
	url := fmt.Sprintf("%s/api/v2/jupiter/route-swap", h.baseURL)
	var response pb.PostJupiterRouteSwapResponse
	err := connections.HTTPPostWithClient[*pb.PostJupiterRouteSwapResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetPools returns pools for given projects.
func (h *HTTPClient) GetPools(ctx context.Context, projects []pb.Project) (*pb.GetPoolsResponse, error) {
	projectsArg := convertSliceArgument("projects", true, projects)
	url := fmt.Sprintf("%s/api/v1/market/pools%s", h.baseURL, projectsArg)
	pools := new(pb.GetPoolsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetPoolsResponse](ctx, url, h.httpClient, pools, h.authHeader); err != nil {
		return nil, err
	}

	return pools, nil
}

func (h *HTTPClient) GetTokenAccounts(ctx context.Context, req *pb.GetTokenAccountsRequest) (*pb.GetTokenAccountsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/account/token-accounts?ownerAddress=%s", h.baseURL, req.OwnerAddress)
	result := new(pb.GetTokenAccountsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetTokenAccountsResponse](ctx, url, h.httpClient, result, h.authHeader); err != nil {
		return nil, err
	}

	return result, nil
}

// GetPrice returns the USDC price of requested tokens
func (h *HTTPClient) GetPrice(ctx context.Context, tokens []string) (*pb.GetPriceResponse, error) {
	tokensArg := convertStrSliceArgument("tokens", true, tokens)
	url := fmt.Sprintf("%s/api/v1/market/price%s", h.baseURL, tokensArg)
	pools := new(pb.GetPriceResponse)
	if err := connections.HTTPGetWithClient[*pb.GetPriceResponse](ctx, url, h.httpClient, pools, h.authHeader); err != nil {
		return nil, err
	}

	return pools, nil
}

// GetQuotes returns the possible amount(s) of outToken for an inToken and the route to achieve it
func (h *HTTPClient) GetQuotes(ctx context.Context, inToken, outToken string, inAmount, slippage float64, limit int32, projects []pb.Project) (*pb.GetQuotesResponse, error) {
	projectString := convertSliceArgument("projects", false, projects)

	url := fmt.Sprintf("%s/api/v1/market/quote?inToken=%s&outToken=%s&inAmount=%v&slippage=%v&limit=%v%s",
		h.baseURL, inToken, outToken, inAmount, slippage, limit, projectString)
	result := new(pb.GetQuotesResponse)
	if err := connections.HTTPGetWithClient[*pb.GetQuotesResponse](ctx, url, h.httpClient, result, h.authHeader); err != nil {
		return nil, err
	}

	return result, nil
}

// PostSubmit posts the transaction string to the Solana network.
func (h *HTTPClient) PostSubmit(ctx context.Context, txBase64 string, opts PostSubmitOpts) (*pb.PostSubmitResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/submit", h.baseURL)
	request := &pb.PostSubmitRequest{
		Transaction:            &pb.TransactionMessage{Content: txBase64},
		SkipPreFlight:          opts.SkipPreFlight,
		FrontRunningProtection: &opts.FrontRunningProtection,
		UseStakedRPCs:          &opts.UseStakedRPCs,
		AllowBackRun:           &opts.AllowBackRun,
		RevenueAddress:         &opts.RevenueAddress,
		Sniping:                &opts.Sniping,
	}

	var response pb.PostSubmitResponse
	err := connections.HTTPPostWithClient[*pb.PostSubmitResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (h *HTTPClient) PostSubmitSnipeV2(ctx context.Context, request *pb.PostSubmitSnipeRequest) (*pb.PostSubmitSnipeResponse, error) {
	url := fmt.Sprintf("%s/api/v2/submit-snipe", h.baseURL)

	var response pb.PostSubmitSnipeResponse
	err := connections.HTTPPostWithClient[*pb.PostSubmitSnipeResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostSubmitBatch posts a bundle of transactions string based on a specific SubmitStrategy to the Solana network.
func (h *HTTPClient) PostSubmitBatch(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/submit-batch", h.baseURL)

	var response pb.PostSubmitBatchResponse
	err := connections.HTTPPostWithClient[*pb.PostSubmitBatchResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostSubmitV2 posts the transaction string to the Solana network.
func (h *HTTPClient) PostSubmitV2(ctx context.Context, txBase64 string, opts PostSubmitOpts) (*pb.PostSubmitResponse, error) {
	url := fmt.Sprintf("%s/api/v2/submit", h.baseURL)
	request := &pb.PostSubmitRequest{
		Transaction:            &pb.TransactionMessage{Content: txBase64},
		SkipPreFlight:          opts.SkipPreFlight,
		FrontRunningProtection: &opts.FrontRunningProtection,
		UseStakedRPCs:          &opts.UseStakedRPCs,
		AllowBackRun:           &opts.AllowBackRun,
		RevenueAddress:         &opts.RevenueAddress,
		Sniping:                &opts.Sniping,
	}

	var response pb.PostSubmitResponse
	err := connections.HTTPPostWithClient[*pb.PostSubmitResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostSubmitBatchV2 posts a bundle of transactions string based on a specific SubmitStrategy to the Solana network.
func (h *HTTPClient) PostSubmitBatchV2(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error) {
	url := fmt.Sprintf("%s/api/v2/submit-batch", h.baseURL)

	var response pb.PostSubmitBatchResponse
	err := connections.HTTPPostWithClient[*pb.PostSubmitBatchResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SignAndSubmit signs the given transaction and submits it.
func (h *HTTPClient) SignAndSubmit(ctx context.Context, tx *pb.TransactionMessage,
	skipPreFlight bool, frontRunningProtection bool, useStakedRPCs bool) (string, error) {
	if h.privateKey == nil {
		return "", ErrPrivateKeyNotFound
	}
	txBase64, err := transaction.SignTxWithPrivateKey(tx.Content, *h.privateKey)
	if err != nil {
		return "", err
	}

	response, err := h.PostSubmit(ctx, txBase64, PostSubmitOpts{
		SkipPreFlight:          skipPreFlight,
		FrontRunningProtection: frontRunningProtection,
		UseStakedRPCs:          useStakedRPCs,
		// Other fields default to zero values
	})
	if err != nil {
		return "", err
	}

	return response.Signature, nil
}

func (h *HTTPClient) SignAndSubmitSnipe(ctx context.Context, transactions []*pb.TransactionMessage, useStakedRPCs bool) ([]string, error) {
	if h.privateKey == nil {
		return nil, ErrPrivateKeyNotFound
	}

	entries := make([]*pb.PostSubmitRequestEntry, len(transactions))

	for i, tx := range transactions {
		txBase64, err := transaction.SignTxWithPrivateKey(tx.Content, *h.privateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to sign transaction: %w", err)
		}

		entries[i] = &pb.PostSubmitRequestEntry{
			Transaction: &pb.TransactionMessage{
				Content:   txBase64,
				IsCleanup: tx.IsCleanup,
			},
			SkipPreFlight: false,
		}
	}

	snipeRequest := &pb.PostSubmitSnipeRequest{
		Entries:       entries,
		UseStakedRPCs: &useStakedRPCs,
	}

	response, err := h.PostSubmitSnipeV2(ctx, snipeRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to submit snipe request: %w", err)
	}

	signatures := make([]string, 0, len(response.Transactions))
	for _, entry := range response.Transactions {
		if entry.Submitted {
			signatures = append(signatures, entry.Signature)
		}
	}

	return signatures, nil
}

func (h *HTTPClient) SignAndSubmitPaladin(ctx context.Context, tx *pb.TransactionMessage) (string, error) {
	if h.privateKey == nil {
		return "", ErrPrivateKeyNotFound
	}

	txBase64, err := transaction.SignTxWithPrivateKey(tx.Content, *h.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	paladinRequest := &pb.PostSubmitPaladinRequest{
		Transaction: &pb.TransactionMessageV2{
			Content: txBase64,
		},
	}

	response, err := h.PostSubmitPaladinV2(ctx, paladinRequest)
	if err != nil {
		return "", fmt.Errorf("failed to submit paladin request: %w", err)
	}

	return response.Signature, nil
}

// SignAndSubmitBatch signs the given transactions and submits them.
func (h *HTTPClient) SignAndSubmitBatch(ctx context.Context, transactions []*pb.TransactionMessage, useBundle bool,
	opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	if h.privateKey == nil {
		return nil, ErrPrivateKeyNotFound
	}

	if len(transactions) == 1 {
		signature, err := h.SignAndSubmit(ctx, transactions[0], *opts.SkipPreFlight, false, false)
		if err != nil {
			return nil, err
		}
		return &pb.PostSubmitBatchResponse{
			Transactions: []*pb.PostSubmitBatchResponseEntry{
				{
					Signature: signature,
					Error:     "",
					Submitted: true,
				},
			},
		}, nil
	}

	batchRequest, err := buildBatchRequest(transactions, *h.privateKey, useBundle, opts)
	if err != nil {
		return nil, err
	}
	return h.PostSubmitBatch(ctx, batchRequest)
}

// SubmitRaydiumSwap builds a Raydium Swap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := h.PostRaydiumSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return h.SignAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// SubmitRaydiumSwapCPMM builds a Raydium Swap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitRaydiumSwapCPMM(ctx context.Context, request *pb.PostRaydiumCPMMSwapRequest) (string, error) {
	resp, err := h.PostRaydiumCPMMSwap(ctx, request)
	if err != nil {
		return "", err
	}

	sig, err := h.SignAndSubmit(ctx, resp.Transaction, true, false, false)
	if err != nil {
		return "", err
	}

	return sig, nil
}

// SubmitPostPumpFunSwap builds a pumpfun Swap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitPostPumpFunSwap(ctx context.Context, request *pb.PostPumpFunSwapRequest) (string, error) {
	resp, err := h.PostPumpFunSwap(ctx, request)
	if err != nil {
		return "", err
	}
	return h.SignAndSubmit(ctx, &pb.TransactionMessage{
		Content: resp.Transaction.Content,
	}, false, false, false)
}

// SubmitRaydiumRouteSwap builds a Raydium RouteSwap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := h.PostRaydiumRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return h.SignAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// SubmitJupiterSwap builds a Jupiter Swap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitJupiterSwap(ctx context.Context, request *pb.PostJupiterSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := h.PostJupiterSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return h.SignAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// SubmitJupiterSwapInstructions builds a Jupiter Swap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitJupiterSwapInstructions(ctx context.Context, request *pb.PostJupiterSwapInstructionsRequest, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	swapInstructions, err := h.PostJupiterSwapInstructions(ctx, request)
	if err != nil {
		return nil, err
	}

	txBuilder := solana.NewTransactionBuilder()

	addressLookupTable, err := utils.ConvertProtoAddressLookupTable(swapInstructions.AddressLookupTableAddresses)
	if err != nil {
		return nil, err
	}

	txBuilder.WithOpt(solana.TransactionAddressTables(addressLookupTable))

	instructions, err := utils.ConvertJupiterInstructions(swapInstructions.Instructions)
	if err != nil {
		return nil, err
	}

	for _, inst := range instructions {
		txBuilder.AddInstruction(inst)
	}

	txBuilder.SetFeePayer(h.privateKey.PublicKey())
	blockHash, err := h.GetRecentBlockHash(ctx)

	if err != nil {
		panic(fmt.Errorf("server error: could not retrieve block hash: %w", err))
	}

	hash, err := solana.HashFromBase58(blockHash.BlockHash)
	if err != nil {
		return nil, err
	}

	txBuilder.SetRecentBlockHash(hash)
	tx, err := txBuilder.Build()
	if err != nil {
		return nil, err
	}

	err = transaction.PartialSign(tx, h.privateKey.PublicKey(), make(map[solana.PublicKey]solana.PrivateKey))
	if err != nil {
		return nil, err
	}

	var txToBeSigned []*pb.TransactionMessage

	txBase64, err := tx.ToBase64()
	if err != nil {
		return nil, err
	}

	txToBeSigned = append(txToBeSigned, &pb.TransactionMessage{
		Content:   txBase64,
		IsCleanup: false,
	})

	return h.SignAndSubmitBatch(ctx, txToBeSigned, useBundle, opts)
}

// SubmitRaydiumSwapInstructions builds a Raydium Swap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitRaydiumSwapInstructions(ctx context.Context, request *pb.PostRaydiumSwapInstructionsRequest, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	swapInstructions, err := h.PostRaydiumSwapInstructions(ctx, request)
	if err != nil {
		return nil, err
	}

	instructions, err := utils.ConvertRaydiumInstructions(swapInstructions.Instructions)
	if err != nil {
		return nil, err
	}

	txBuilder := solana.NewTransactionBuilder()
	for _, inst := range instructions {
		txBuilder.AddInstruction(inst)
	}

	txBuilder.SetFeePayer(h.privateKey.PublicKey())
	blockHash, err := h.GetRecentBlockHash(ctx)

	if err != nil {
		panic(fmt.Errorf("server error: could not retrieve block hash: %w", err))
	}

	hash, err := solana.HashFromBase58(blockHash.BlockHash)
	if err != nil {
		return nil, err
	}

	txBuilder.SetRecentBlockHash(hash)
	tx, err := txBuilder.Build()
	if err != nil {
		return nil, err
	}

	err = transaction.PartialSign(tx, h.privateKey.PublicKey(), make(map[solana.PublicKey]solana.PrivateKey))
	if err != nil {
		return nil, err
	}

	var txToBeSigned []*pb.TransactionMessage

	txBase64, err := tx.ToBase64()
	if err != nil {
		return nil, err
	}

	txToBeSigned = append(txToBeSigned, &pb.TransactionMessage{
		Content:   txBase64,
		IsCleanup: false,
	})

	return h.SignAndSubmitBatch(ctx, txToBeSigned, useBundle, opts)
}

// SubmitJupiterRouteSwap builds a Jupiter RouteSwap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := h.PostJupiterRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return h.SignAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// GetRecentBlockHash returns recent block hash.
func (h *HTTPClient) GetRecentBlockHash(ctx context.Context) (*pb.GetRecentBlockHashResponse, error) {
	url := fmt.Sprintf("%s/api/v1/system/blockhash", h.baseURL)
	response := new(pb.GetRecentBlockHashResponse)
	if err := connections.HTTPGetWithClient[*pb.GetRecentBlockHashResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// GetRecentBlockHash returns recent block hash, supports optional offset.
func (h *HTTPClient) GetRecentBlockHashV2(ctx context.Context, offset uint64) (*pb.GetRecentBlockHashResponseV2, error) {
	url := fmt.Sprintf("%s/api/v2/system/blockhash?offset=%d", h.baseURL, offset)
	response := new(pb.GetRecentBlockHashResponseV2)
	if err := connections.HTTPGetWithClient[*pb.GetRecentBlockHashResponseV2](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

func (h *HTTPClient) GetPriorityFee(ctx context.Context, project pb.Project, percentile *float64) (*pb.GetPriorityFeeResponse, error) {
	url := fmt.Sprintf("%s/api/v2/system/priority-fee?project=%v", h.baseURL, project)
	if percentile != nil {
		url = fmt.Sprintf("%s/api/v2/system/priority-fee?project=%v&percentile=%v", h.baseURL, project, *percentile)
	}
	response := new(pb.GetPriorityFeeResponse)
	if err := connections.HTTPGetWithClient[*pb.GetPriorityFeeResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

func (h *HTTPClient) GetPriorityFeeByProgram(ctx context.Context, programs []string) (*pb.GetPriorityFeeByProgramResponse, error) {
	url := fmt.Sprintf("%s/api/v2/system/priority-fee-by-program?programs=%v", h.baseURL, strings.Join(programs, "&programs="))
	response := new(pb.GetPriorityFeeByProgramResponse)
	if err := connections.HTTPGetWithClient[*pb.GetPriorityFeeByProgramResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

func (h *HTTPClient) GetLeaderSchedule(ctx context.Context, maxSlots uint) (*pb.GetLeaderScheduleResponse, error) {
	url := fmt.Sprintf("%s/api/v2/system/leader-schedule?maxSlots=%v", h.baseURL, maxSlots)
	response := new(pb.GetLeaderScheduleResponse)
	if err := connections.HTTPGetWithClient[*pb.GetLeaderScheduleResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil

}

type stringable interface {
	String() string
}

func convertSliceArgument[T stringable](argName string, isFirst bool, s []T) string {
	r := ""
	for i, v := range s {
		arg := fmt.Sprintf("%v=%v", argName, v.String())
		if i == 0 && isFirst {
			r += "?" + arg
		} else {
			r += "&" + arg
		}
	}
	return r
}

func convertStrSliceArgument(argName string, isFirst bool, s []string) string {
	r := ""
	for i, v := range s {
		arg := fmt.Sprintf("%v=%v", argName, v)
		if i == 0 && isFirst {
			r += "?" + arg
		} else {
			r += "&" + arg
		}
	}
	return r
}
