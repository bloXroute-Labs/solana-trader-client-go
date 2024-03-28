package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bloXroute-Labs/solana-trader-client-go/connections"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/bloXroute-Labs/solana-trader-proto/common"
	"github.com/gagliardetto/solana-go"
)

type HTTPClient struct {
	pb.UnimplementedApiServer

	baseURL    string
	httpClient *http.Client
	requestID  utils.RequestID
	privateKey *solana.PrivateKey
	authHeader string
}

// NewHTTPClient connects to Mainnet Trader API
func NewHTTPClient() *HTTPClient {
	opts := DefaultRPCOpts(MainnetNYHTTP)
	return NewHTTPClientWithOpts(nil, opts)
}

// NewHTTPTestnet connects to Testnet Trader API
func NewHTTPTestnet() *HTTPClient {
	opts := DefaultRPCOpts(TestnetHTTP)
	return NewHTTPClientWithOpts(nil, opts)
}

// NewHTTPDevnet connects to Devnet Trader API
func NewHTTPDevnet() *HTTPClient {
	opts := DefaultRPCOpts(DevnetHTTP)
	return NewHTTPClientWithOpts(nil, opts)
}

// NewHTTPLocal connects to local Trader API
func NewHTTPLocal() *HTTPClient {
	opts := DefaultRPCOpts(LocalHTTP)
	return NewHTTPClientWithOpts(nil, opts)
}

// NewHTTPClientWithOpts connects to custom Trader API (set client to nil to use default client)
func NewHTTPClientWithOpts(client *http.Client, opts RPCOpts) *HTTPClient {
	if client == nil {
		client = &http.Client{}
	}

	return &HTTPClient{
		baseURL:    opts.Endpoint,
		httpClient: client,
		privateKey: opts.PrivateKey,
		authHeader: opts.AuthHeader,
	}
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
func (h *HTTPClient) GetRateLimit(ctx context.Context, request *pb.GetRateLimitRequest) (*pb.GetRateLimitResponse, error) {
	url := fmt.Sprintf("%s/api/v2/account/rate-limit", h.baseURL)
	response := new(pb.GetRateLimitResponse)
	if err := connections.HTTPGetWithClient[*pb.GetRateLimitResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// GetRaydiumPools returns pools on Raydium
func (h *HTTPClient) GetRaydiumPools(ctx context.Context, request *pb.GetRaydiumPoolsRequest) (*pb.GetRaydiumPoolsResponse, error) {
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
	url := fmt.Sprintf("%s/api/v2/jupiter/quotes?inToken=%s&outToken=%s&inAmount=%v&slippage=%v",
		h.baseURL, request.InToken, request.OutToken, request.InAmount, request.Slippage)
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

// GetOrderbook returns the requested market's orderbook (e.h. asks and bids). Set limit to 0 for all bids / asks.
func (h *HTTPClient) GetOrderbook(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetOrderbookResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/orderbooks/%s?limit=%v&project=%v", h.baseURL, market, limit, project)
	orderbook := new(pb.GetOrderbookResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOrderbookResponse](ctx, url, h.httpClient, orderbook, h.authHeader); err != nil {
		return nil, err
	}

	return orderbook, nil
}

// GetMarketDepth returns the requested market's coalesced price data (e.h. asks and bids). Set limit to 0 for all bids / asks.
func (h *HTTPClient) GetMarketDepth(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetMarketDepthResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/depth/%s?limit=%v&project=%v", h.baseURL, market, limit, project)
	mktDepth := new(pb.GetMarketDepthResponse)
	if err := connections.HTTPGetWithClient[*pb.GetMarketDepthResponse](ctx, url, h.httpClient, mktDepth, h.authHeader); err != nil {
		return nil, err
	}

	return mktDepth, nil
}

// GetTrades returns the requested market's currently executing trades. Set limit to 0 for all trades.
func (h *HTTPClient) GetTrades(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetTradesResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/trades/%s?limit=%v&project=%v", h.baseURL, market, limit, project)
	marketTrades := new(pb.GetTradesResponse)
	if err := connections.HTTPGetWithClient[*pb.GetTradesResponse](ctx, url, h.httpClient, marketTrades, h.authHeader); err != nil {
		return nil, err
	}

	return marketTrades, nil
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

// GetTickers returns the requested market tickets. Set market to "" for all markets.
func (h *HTTPClient) GetTickers(ctx context.Context, market string, project pb.Project) (*pb.GetTickersResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/tickers/%s?project=%v", h.baseURL, market, project)
	tickers := new(pb.GetTickersResponse)
	if err := connections.HTTPGetWithClient[*pb.GetTickersResponse](ctx, url, h.httpClient, tickers, h.authHeader); err != nil {
		return nil, err
	}

	return tickers, nil
}

// GetOpenOrders returns all open orders by owner address and market
func (h *HTTPClient) GetOpenOrders(ctx context.Context, market string, owner string, openOrdersAddress string, project pb.Project) (*pb.GetOpenOrdersResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/openorders/%s?address=%s&openOrdersAddress=%s&project=%s", h.baseURL, market, owner, openOrdersAddress, project)
	orders := new(pb.GetOpenOrdersResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOpenOrdersResponse](ctx, url, h.httpClient, orders, h.authHeader); err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOrderByID returns an order by id
func (h *HTTPClient) GetOrderByID(ctx context.Context, in *pb.GetOrderByIDRequest) (*pb.GetOrderByIDResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/orderbyid/%s?market=%s&project=%s", h.baseURL, in.OrderID, in.Market, in.Project)
	orders := new(pb.GetOrderByIDResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOrderByIDResponse](ctx, url, h.httpClient, orders, h.authHeader); err != nil {
		return nil, err
	}

	return orders, nil
}

// GetMarkets returns the list of all available named markets
func (h *HTTPClient) GetMarkets(ctx context.Context) (*pb.GetMarketsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/markets", h.baseURL)
	markets := new(pb.GetMarketsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetMarketsResponse](ctx, url, h.httpClient, markets, h.authHeader); err != nil {
		return nil, err
	}

	return markets, nil
}

// GetUnsettled returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (h *HTTPClient) GetUnsettled(ctx context.Context, market string, owner string, project pb.Project) (*pb.GetUnsettledResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/unsettled/%s?ownerAddress=%s&project=%s", h.baseURL, market, owner, project)
	result := new(pb.GetUnsettledResponse)
	if err := connections.HTTPGetWithClient[*pb.GetUnsettledResponse](ctx, url, h.httpClient, result, h.authHeader); err != nil {
		return nil, err
	}

	return result, nil
}

// GetAccountBalance returns all OpenOrders accounts for a given market with the amountsctx context.Context,  of unsettled funds
func (h *HTTPClient) GetAccountBalance(ctx context.Context, owner string) (*pb.GetAccountBalanceResponse, error) {
	url := fmt.Sprintf("%s/api/v2/balance?ownerAddress=%s", h.baseURL, owner)
	result := new(pb.GetAccountBalanceResponse)
	if err := connections.HTTPGetWithClient[*pb.GetAccountBalanceResponse](ctx, url, h.httpClient, result, h.authHeader); err != nil {
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
func (h *HTTPClient) PostSubmit(ctx context.Context, txBase64 string, skipPreFlight bool, frontRunningProtection bool) (*pb.PostSubmitResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/submit", h.baseURL)
	request := &pb.PostSubmitRequest{Transaction: &pb.TransactionMessage{Content: txBase64},
		SkipPreFlight:          skipPreFlight,
		FrontRunningProtection: &frontRunningProtection}

	var response pb.PostSubmitResponse
	err := connections.HTTPPostWithClient[*pb.PostSubmitResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
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
func (h *HTTPClient) PostSubmitV2(ctx context.Context, txBase64 string, skipPreFlight bool, frontRunningProtection bool) (*pb.PostSubmitResponse, error) {
	url := fmt.Sprintf("%s/api/v2/submit", h.baseURL)
	request := &pb.PostSubmitRequest{Transaction: &pb.TransactionMessage{Content: txBase64},
		SkipPreFlight:          skipPreFlight,
		FrontRunningProtection: &frontRunningProtection}

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
func (h *HTTPClient) SignAndSubmit(ctx context.Context, tx *pb.TransactionMessage, skipPreFlight bool, frontRunningProtection bool) (string, error) {
	if h.privateKey == nil {
		return "", ErrPrivateKeyNotFound
	}
	txBase64, err := transaction.SignTxWithPrivateKey(tx.Content, *h.privateKey)
	if err != nil {
		return "", err
	}

	response, err := h.PostSubmit(ctx, txBase64, skipPreFlight, frontRunningProtection)
	if err != nil {
		return "", err
	}

	return response.Signature, nil
}

// SignAndSubmitBatch signs the given transactions and submits them.
func (h *HTTPClient) SignAndSubmitBatch(ctx context.Context, transactions []*pb.TransactionMessage, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	if h.privateKey == nil {
		return nil, ErrPrivateKeyNotFound
	}
	batchRequest, err := buildBatchRequest(transactions, *h.privateKey, useBundle, opts)
	if err != nil {
		return nil, err
	}
	return h.PostSubmitBatch(ctx, batchRequest)
}

// PostTradeSwap returns a partially signed transaction for submitting a swap request
func (h *HTTPClient) PostTradeSwap(ctx context.Context, ownerAddress, inToken, outToken string, inAmount, slippage float64, project pb.Project) (*pb.TradeSwapResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/swap", h.baseURL)
	request := &pb.TradeSwapRequest{
		OwnerAddress: ownerAddress,
		InToken:      inToken,
		OutToken:     outToken,
		InAmount:     inAmount,
		Slippage:     slippage,
		Project:      project,
	}

	var response pb.TradeSwapResponse
	err := connections.HTTPPostWithClient[*pb.TradeSwapResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitTradeSwap builds a TradeSwap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitTradeSwap(ctx context.Context, owner, inToken, outToken string, inAmount, slippage float64, project pb.Project, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := h.PostTradeSwap(ctx, owner, inToken, outToken, inAmount, slippage, project)
	if err != nil {
		return nil, err
	}
	return h.SignAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// PostRouteTradeSwap returns a partially signed transaction(s) for submitting a route swap request
func (h *HTTPClient) PostRouteTradeSwap(ctx context.Context, request *pb.RouteTradeSwapRequest) (*pb.TradeSwapResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/route-swap", h.baseURL)

	var response pb.TradeSwapResponse
	err := connections.HTTPPostWithClient[*pb.TradeSwapResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitRouteTradeSwap builds a RouteTradeSwap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitRouteTradeSwap(ctx context.Context, request *pb.RouteTradeSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := h.PostRouteTradeSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return h.SignAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// SubmitRaydiumSwap builds a Raydium Swap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := h.PostRaydiumSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return h.SignAndSubmitBatch(ctx, resp.Transactions, false, opts)
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

	instructions, err := utils.ConvertProtoInstructionsToSolanaInstructions(swapInstructions.Instructions)
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

// SubmitJupiterRouteSwap builds a Jupiter RouteSwap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := h.PostJupiterRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return h.SignAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// PostOrder returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (h *HTTPClient) PostOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/place", h.baseURL)
	request := &pb.PostOrderRequest{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              types,
		Amount:            amount,
		Price:             price,
		Project:           project,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
	}

	var response pb.PostOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostOrderResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitOrder builds a Serum market order, signs it, and submits to the network.
func (h *HTTPClient) SubmitOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error) {
	order, err := h.PostOrder(ctx, owner, payer, market, side, types, amount, price, project, opts)
	if err != nil {
		return "", err
	}

	skipPreFlight := true
	if opts.SkipPreFlight != nil {
		skipPreFlight = *opts.SkipPreFlight
	}
	sig, err := h.SignAndSubmit(ctx, order.Transaction, skipPreFlight, false)
	if err != nil {
		return "", err
	}
	return sig, err
}

// PostCancelOrder builds a Serum cancel order.
func (h *HTTPClient) PostCancelOrder(
	ctx context.Context,
	orderID string,
	side pb.Side,
	owner,
	market,
	openOrders string,
	project pb.Project,
) (*pb.PostCancelOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/cancel", h.baseURL)
	request := &pb.PostCancelOrderRequest{
		OrderID:           orderID,
		Side:              side,
		OwnerAddress:      owner,
		MarketAddress:     market,
		Project:           project,
		OpenOrdersAddress: openOrders,
	}

	var response pb.PostCancelOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostCancelOrderResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// SubmitCancelOrder builds a Serum cancel order, signs and submits it to the network.
func (h *HTTPClient) SubmitCancelOrder(
	ctx context.Context,
	orderID string,
	side pb.Side,
	owner,
	market,
	openOrders string,
	project pb.Project,
	skipPreFlight bool,
) (string, error) {
	order, err := h.PostCancelOrder(ctx, orderID, side, owner, market, openOrders, project)
	if err != nil {
		return "", err
	}

	return h.SignAndSubmit(ctx, order.Transaction, skipPreFlight, false)
}

// PostCancelByClientOrderID builds a Serum cancel order by client ID.
func (h *HTTPClient) PostCancelByClientOrderID(
	ctx context.Context,
	clientOrderID uint64,
	owner,
	market,
	openOrders string,
	project pb.Project,
) (*pb.PostCancelOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/cancelbyid", h.baseURL)
	request := &pb.PostCancelByClientOrderIDRequest{
		ClientOrderID:     clientOrderID,
		OwnerAddress:      owner,
		MarketAddress:     market,
		Project:           project,
		OpenOrdersAddress: openOrders,
	}

	var response pb.PostCancelOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostCancelOrderResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// SubmitCancelByClientOrderID builds a Serum cancel order by client ID, signs and submits it to the network.
func (h *HTTPClient) SubmitCancelByClientOrderID(
	ctx context.Context,
	clientOrderID uint64,
	owner,
	market,
	openOrders string,
	project pb.Project,
	skipPreFlight bool,
) (string, error) {
	order, err := h.PostCancelByClientOrderID(ctx, clientOrderID, owner, market, openOrders, project)
	if err != nil {
		return "", err
	}

	return h.SignAndSubmit(ctx, order.Transaction, skipPreFlight, false)
}

func (h *HTTPClient) PostCancelAll(ctx context.Context, market, owner string, openOrdersAddresses []string, project pb.Project) (*pb.PostCancelAllResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/cancelall", h.baseURL)
	request := &pb.PostCancelAllRequest{
		Market:              market,
		OwnerAddress:        owner,
		OpenOrdersAddresses: openOrdersAddresses,
		Project:             project,
	}

	var response pb.PostCancelAllResponse
	err := connections.HTTPPostWithClient[*pb.PostCancelAllResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (h *HTTPClient) SubmitCancelAll(ctx context.Context, market, owner string, openOrders []string, project pb.Project, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	orders, err := h.PostCancelAll(ctx, market, owner, openOrders, project)
	if err != nil {
		return nil, err
	}
	return h.SignAndSubmitBatch(ctx, orders.Transactions, false, opts)
}

// PostSettle returns a partially signed transaction for settling market funds. Typically, you want to use SubmitSettle instead of this.
func (h *HTTPClient) PostSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, project pb.Project) (*pb.PostSettleResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/settle", h.baseURL)
	request := &pb.PostSettleRequest{
		OwnerAddress:      owner,
		Market:            market,
		BaseTokenWallet:   baseTokenWallet,
		QuoteTokenWallet:  quoteTokenWallet,
		OpenOrdersAddress: openOrdersAccount,
		Project:           project,
	}

	var response pb.PostSettleResponse
	err := connections.HTTPPostWithClient[*pb.PostSettleResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitSettle builds a market SubmitSettle transaction, signs it, and submits to the network.
func (h *HTTPClient) SubmitSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, project pb.Project, skipPreflight bool) (string, error) {
	order, err := h.PostSettle(ctx, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount, project)
	if err != nil {
		return "", err
	}

	return h.SignAndSubmit(ctx, order.Transaction, skipPreflight, false)
}

func (h *HTTPClient) PostReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/replacebyclientid", h.baseURL)
	request := &pb.PostOrderRequest{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              types,
		Amount:            amount,
		Price:             price,
		Project:           project,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
	}

	var response pb.PostOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostOrderResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (h *HTTPClient) SubmitReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error) {
	order, err := h.PostReplaceByClientOrderID(ctx, owner, payer, market, side, types, amount, price, project, opts)
	if err != nil {
		return "", err
	}
	skipPreFlight := true
	if opts.SkipPreFlight != nil {
		skipPreFlight = *opts.SkipPreFlight
	}
	return h.SignAndSubmit(ctx, order.Transaction, skipPreFlight, false)
}

func (h *HTTPClient) PostReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/replace", h.baseURL)
	request := &pb.PostReplaceOrderRequest{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              types,
		Amount:            amount,
		Price:             price,
		Project:           project,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
		OrderID:           orderID,
	}

	var response pb.PostOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostOrderResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (h *HTTPClient) SubmitReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error) {
	order, err := h.PostReplaceOrder(ctx, orderID, owner, payer, market, side, types, amount, price, project, opts)
	if err != nil {
		return "", err
	}
	skipPreFlight := true
	if opts.SkipPreFlight != nil {
		skipPreFlight = *opts.SkipPreFlight
	}
	return h.SignAndSubmit(ctx, order.Transaction, skipPreFlight, false)
}

// GetRecentBlockHash subscribes to a stream for getting recent block hash.
func (h *HTTPClient) GetRecentBlockHash(ctx context.Context) (*pb.GetRecentBlockHashResponse, error) {
	url := fmt.Sprintf("%s/api/v1/system/blockhash", h.baseURL)
	response := new(pb.GetRecentBlockHashResponse)
	if err := connections.HTTPGetWithClient[*pb.GetRecentBlockHashResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

func (h *HTTPClient) GetPriorityFee(ctx context.Context, percentile *float64) (*pb.GetPriorityFeeResponse, error) {
	url := fmt.Sprintf("%s/api/v2/system/priority-fee", h.baseURL)
	if percentile != nil {
		url = fmt.Sprintf("%s/api/v2/system/priority-fee?percentile=%v", h.baseURL, *percentile)
	}
	response := new(pb.GetPriorityFeeResponse)
	if err := connections.HTTPGetWithClient[*pb.GetPriorityFeeResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

//V2 Openbook

// GetMarketsV2 returns the list of all available named markets
func (h *HTTPClient) GetMarketsV2(ctx context.Context) (*pb.GetMarketsResponseV2, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/markets", h.baseURL)
	markets := new(pb.GetMarketsResponseV2)
	if err := connections.HTTPGetWithClient[*pb.GetMarketsResponseV2](ctx, url, h.httpClient, markets, h.authHeader); err != nil {
		return nil, err
	}

	return markets, nil
}

// GetOrderbookV2 returns the requested market's orderbook (e.h. asks and bids). Set limit to 0 for all bids / asks.
func (h *HTTPClient) GetOrderbookV2(ctx context.Context, market string, limit uint32) (*pb.GetOrderbookResponseV2, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/orderbooks/%s?limit=%v", h.baseURL, market, limit)
	orderbook := new(pb.GetOrderbookResponseV2)
	if err := connections.HTTPGetWithClient[*pb.GetOrderbookResponseV2](ctx, url, h.httpClient, orderbook, h.authHeader); err != nil {
		return nil, err
	}

	return orderbook, nil
}

// GetMarketDepthV2 returns the requested market's coalesced price data (e.h. asks and bids). Set limit to 0 for all bids / asks.
func (h *HTTPClient) GetMarketDepthV2(ctx context.Context, market string, limit uint32) (*pb.GetMarketDepthResponseV2, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/depth/%s?limit=%v", h.baseURL, market, limit)
	mktDepth := new(pb.GetMarketDepthResponseV2)
	if err := connections.HTTPGetWithClient[*pb.GetMarketDepthResponseV2](ctx, url, h.httpClient, mktDepth, h.authHeader); err != nil {
		return nil, err
	}

	return mktDepth, nil
}

// GetTickersV2 returns the requested market tickets. Set market to "" for all markets.
func (h *HTTPClient) GetTickersV2(ctx context.Context, market string) (*pb.GetTickersResponseV2, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/tickers/%s", h.baseURL, market)
	tickers := new(pb.GetTickersResponseV2)
	if err := connections.HTTPGetWithClient[*pb.GetTickersResponseV2](ctx, url, h.httpClient, tickers, h.authHeader); err != nil {
		return nil, err
	}

	return tickers, nil
}

// GetOpenOrdersV2 returns all open orders by owner address and market
func (h *HTTPClient) GetOpenOrdersV2(ctx context.Context, market string, owner string, openOrdersAddress string, orderID string, clientOrderID uint64) (*pb.GetOpenOrdersResponse, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/open-orders/%s?address=%s&openOrdersAddress=%s&orderID=%s&clientOrderID=%v",
		h.baseURL, market, owner, openOrdersAddress, orderID, clientOrderID)

	orders := new(pb.GetOpenOrdersResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOpenOrdersResponse](ctx, url, h.httpClient, orders, h.authHeader); err != nil {
		return nil, err
	}

	return orders, nil
}

// GetUnsettledV2 returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (h *HTTPClient) GetUnsettledV2(ctx context.Context, market string, owner string) (*pb.GetUnsettledResponse, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/unsettled/%s?ownerAddress=%s", h.baseURL, market, owner)
	result := new(pb.GetUnsettledResponse)
	if err := connections.HTTPGetWithClient[*pb.GetUnsettledResponse](ctx, url, h.httpClient, result, h.authHeader); err != nil {
		return nil, err
	}

	return result, nil
}

// GetBundleResult subscribes to a stream for getting recent block hash.
func (h *HTTPClient) GetBundleResult(ctx context.Context, uuid string) (*pb.GetBundleResultResponse, error) {
	url := fmt.Sprintf("%s/api/v2/trade/bundle-result/%s", h.baseURL, uuid)

	response := new(pb.GetBundleResultResponse)
	if err := connections.HTTPGetWithClient[*pb.GetBundleResultResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// PostOrderV2 returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (h *HTTPClient) PostOrderV2(ctx context.Context, owner, payer, market string, side string, orderType string, amount,
	price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/place", h.baseURL)
	request := &pb.PostOrderRequestV2{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              orderType,
		Amount:            amount,
		Price:             price,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
	}

	var response pb.PostOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostOrderResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostOrderV2WithPriorityFee returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (h *HTTPClient) PostOrderV2WithPriorityFee(ctx context.Context, owner, payer, market string, side string,
	orderType string, amount, price float64, computeLimit uint32, computePrice uint64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/place", h.baseURL)
	request := &pb.PostOrderRequestV2{
		OwnerAddress: owner,
		PayerAddress: payer,

		Market:            market,
		Side:              side,
		Type:              orderType,
		Amount:            amount,
		Price:             price,
		ComputeLimit:      computeLimit,
		ComputePrice:      computePrice,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
	}

	var response pb.PostOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostOrderResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitOrderV2 builds a Serum market order, signs it, and submits to the network.
func (h *HTTPClient) SubmitOrderV2(ctx context.Context, owner, payer, market string, side string, orderType string,
	amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := h.PostOrderV2(ctx, owner, payer, market, side, orderType, amount, price, opts)
	if err != nil {
		return "", err
	}
	skipPreFlight := true
	if opts.SkipPreFlight != nil {
		skipPreFlight = *opts.SkipPreFlight
	}
	sig, err := h.SignAndSubmit(ctx, order.Transaction, skipPreFlight, false)
	return sig, err
}

// SubmitOrderV2WithPriorityFee builds a Serum market order, signs it, and submits to the network.
func (h *HTTPClient) SubmitOrderV2WithPriorityFee(ctx context.Context, owner, payer, market string, side string,
	orderType string, amount, price float64, computeLimit uint32, computePrice uint64, opts PostOrderOpts) (string, error) {
	order, err := h.PostOrderV2WithPriorityFee(ctx, owner, payer, market, side, orderType, amount, price,
		computeLimit, computePrice, opts)
	if err != nil {
		return "", err
	}
	skipPreFlight := true
	if opts.SkipPreFlight != nil {
		skipPreFlight = *opts.SkipPreFlight
	}
	sig, err := h.SignAndSubmit(ctx, order.Transaction, skipPreFlight, false)
	return sig, err
}

// PostCancelOrderV2 builds a Serum cancel order.
func (h *HTTPClient) PostCancelOrderV2(
	ctx context.Context,
	orderID string,
	clientOrderID uint64,
	side string,
	owner,
	market,
	openOrders string,
) (*pb.PostCancelOrderResponseV2, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/cancel", h.baseURL)
	request := &pb.PostCancelOrderRequestV2{
		OrderID:           orderID,
		ClientOrderID:     clientOrderID,
		Side:              side,
		OwnerAddress:      owner,
		MarketAddress:     market,
		OpenOrdersAddress: openOrders,
	}

	var response pb.PostCancelOrderResponseV2
	err := connections.HTTPPostWithClient[*pb.PostCancelOrderResponseV2](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// SubmitCancelOrderV2 builds a Serum cancel order, signs and submits it to the network.
func (h *HTTPClient) SubmitCancelOrderV2(
	ctx context.Context,
	orderID string,
	clientOrderID uint64,
	side string,
	owner,
	market,
	openOrders string,
	opts SubmitOpts,
) (*pb.PostSubmitBatchResponse, error) {
	order, err := h.PostCancelOrderV2(ctx, orderID, clientOrderID, side, owner, market, openOrders)
	if err != nil {
		return nil, err
	}

	return h.SignAndSubmitBatch(ctx, order.Transactions, false, opts)
}

// PostSettleV2 returns a partially signed transaction for settling market funds. Typically, you want to use SubmitSettle instead of this.
func (h *HTTPClient) PostSettleV2(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string) (*pb.PostSettleResponse, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/settle", h.baseURL)
	request := &pb.PostSettleRequestV2{
		OwnerAddress:      owner,
		Market:            market,
		BaseTokenWallet:   baseTokenWallet,
		QuoteTokenWallet:  quoteTokenWallet,
		OpenOrdersAddress: openOrdersAccount,
	}

	var response pb.PostSettleResponse
	err := connections.HTTPPostWithClient[*pb.PostSettleResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitSettleV2 builds a market SubmitSettle transaction, signs it, and submits to the network.
func (h *HTTPClient) SubmitSettleV2(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, skipPreflight bool) (string, error) {
	order, err := h.PostSettleV2(ctx, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount)
	if err != nil {
		return "", err
	}

	return h.SignAndSubmit(ctx, order.Transaction, skipPreflight, false)
}

func (h *HTTPClient) PostReplaceOrderV2(ctx context.Context, orderID, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/replace", h.baseURL)
	request := &pb.PostReplaceOrderRequestV2{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              orderType,
		Amount:            amount,
		Price:             price,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
		OrderID:           orderID,
	}

	var response pb.PostOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostOrderResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (h *HTTPClient) SubmitReplaceOrderV2(ctx context.Context, orderID, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := h.PostReplaceOrderV2(ctx, orderID, owner, payer, market, side, orderType, amount, price, opts)
	if err != nil {
		return "", err
	}
	skipPreFlight := true
	if opts.SkipPreFlight != nil {
		skipPreFlight = *opts.SkipPreFlight
	}
	return h.SignAndSubmit(ctx, order.Transaction, skipPreFlight, false)
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
