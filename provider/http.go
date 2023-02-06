package provider

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/connections"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/bloXroute-Labs/solana-trader-proto/common"
	"github.com/gagliardetto/solana-go"
	"net/http"
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
	opts := DefaultRPCOpts(MainnetHTTP)
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

// GetOrderbook returns the requested market's orderbook (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (h *HTTPClient) GetOrderbook(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetOrderbookResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/orderbooks/%s?limit=%v&project=%v", h.baseURL, market, limit, project)
	orderbook := new(pb.GetOrderbookResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOrderbookResponse](ctx, url, h.httpClient, orderbook, h.authHeader); err != nil {
		return nil, err
	}

	return orderbook, nil
}

// GetMarketDepth returns the requested market's coalesced price data (e.g. asks and bids). Set limit to 0 for all bids / asks.
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

// GetOpenOrders returns all opened orders by owner address and market
func (h *HTTPClient) GetOpenOrders(ctx context.Context, market string, owner string, openOrdersAddress string, project pb.Project) (*pb.GetOpenOrdersResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/openorders/%s?address=%s&openOrdersAddress=%s&project=%s", h.baseURL, market, owner, openOrdersAddress, project)
	orders := new(pb.GetOpenOrdersResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOpenOrdersResponse](ctx, url, h.httpClient, orders, h.authHeader); err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOpenPerpOrders returns all opened perp orders
func (h *HTTPClient) GetOpenPerpOrders(ctx context.Context, request *pb.GetOpenPerpOrdersRequest) (*pb.GetOpenPerpOrdersResponse, error) {
	contractsString := convertSliceArgument("contracts", false, request.Contracts)
	url := fmt.Sprintf("%s/api/v1/trade/perp/open-orders?ownerAddress=%s&accountAddress=%s&project=%s%s",
		h.baseURL, request.OwnerAddress, request.AccountAddress, request.Project, contractsString)
	orders := new(pb.GetOpenPerpOrdersResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOpenPerpOrdersResponse](ctx, url, h.httpClient, orders, h.authHeader); err != nil {
		return nil, err
	}

	return orders, nil
}

// PostCancelPerpOrder returns a partially signed transaction for canceling perp order
func (h *HTTPClient) PostCancelPerpOrder(ctx context.Context, request *pb.PostCancelPerpOrderRequest) (*pb.PostCancelPerpOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/perp/cancelbyid", h.baseURL)
	response := new(pb.PostCancelPerpOrderResponse)
	if err := connections.HTTPPostWithClient[*pb.PostCancelPerpOrderResponse](ctx, url, h.httpClient, request, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// PostCancelPerpOrders returns a partially signed transaction for canceling all perp orders of a user
func (h *HTTPClient) PostCancelPerpOrders(ctx context.Context, request *pb.PostCancelPerpOrderRequest) (*pb.PostCancelPerpOrdersResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/perp/cancel", h.baseURL)
	response := new(pb.PostCancelPerpOrdersResponse)
	if err := connections.HTTPPostWithClient[*pb.PostCancelPerpOrdersResponse](ctx, url, h.httpClient, request, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// PostCreateUser returns a partially signed transaction for creating a user
func (h *HTTPClient) PostCreateUser(ctx context.Context, request *pb.PostCreateUserRequest) (*pb.PostCreateUserResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/user", h.baseURL)
	response := new(pb.PostCreateUserResponse)
	if err := connections.HTTPPostWithClient[*pb.PostCreateUserResponse](ctx, url, h.httpClient, request, response, h.authHeader); err != nil {
		return nil, err
	}
	return response, nil
}

// GetUser returns a user's info
func (h *HTTPClient) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/user?ownerAddress=%s&project=%s", h.baseURL, request.OwnerAddress, request.Project)
	resp := new(pb.GetUserResponse)
	if err := connections.HTTPGetWithClient[*pb.GetUserResponse](ctx, url, h.httpClient, resp, h.authHeader); err != nil {
		return nil, err
	}

	return resp, nil
}

// PostDepositCollateral returns a partially signed transaction for posting collateral
func (h *HTTPClient) PostDepositCollateral(ctx context.Context, request *pb.PostDepositCollateralRequest) (*pb.PostDepositCollateralResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/perp/collateral/deposit", h.baseURL)
	response := new(pb.PostDepositCollateralResponse)
	if err := connections.HTTPPostWithClient[*pb.PostDepositCollateralResponse](ctx, url, h.httpClient, request, response, h.authHeader); err != nil {
		return nil, err
	}
	return response, nil
}

// PostWithdrawCollateral returns a partially signed transaction for withdrawing collateral
func (h *HTTPClient) PostWithdrawCollateral(ctx context.Context, request *pb.PostWithdrawCollateralRequest) (*pb.PostWithdrawCollateralResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/perp/collateral/withdraw", h.baseURL)
	response := new(pb.PostWithdrawCollateralResponse)
	if err := connections.HTTPPostWithClient[*pb.PostWithdrawCollateralResponse](ctx, url, h.httpClient, request, response, h.authHeader); err != nil {
		return nil, err
	}
	return response, nil

}

// GetPerpPositions returns all perp positions by owner address and market
func (h *HTTPClient) GetPerpPositions(ctx context.Context, request *pb.GetPerpPositionsRequest) (*pb.GetPerpPositionsResponse, error) {
	var strs []string
	for _, c := range request.Contracts {
		strs = append(strs, fmt.Sprint(c))
	}

	contractsArg := convertStrSliceArgument("contracts", false, strs)
	url := fmt.Sprintf("%s/api/v1/trade/perp/positions?ownerAddress=%s&accountAddress=%s&project=%s%s", h.baseURL,
		request.OwnerAddress, request.AccountAddress, request.Project, contractsArg)
	positions := new(pb.GetPerpPositionsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetPerpPositionsResponse](ctx, url, h.httpClient, positions, h.authHeader); err != nil {
		return nil, err
	}

	return positions, nil
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
	url := fmt.Sprintf("%s/api/v1/account/balance?ownerAddress=%s", h.baseURL, owner)
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
func (h *HTTPClient) PostSubmit(ctx context.Context, txBase64 string, skipPreFlight bool) (*pb.PostSubmitResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/submit", h.baseURL)
	request := &pb.PostSubmitRequest{Transaction: &pb.TransactionMessage{Content: txBase64}, SkipPreFlight: skipPreFlight}

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

// signAndSubmit signs the given transaction and submits it.
func (h *HTTPClient) signAndSubmit(ctx context.Context, tx *pb.TransactionMessage, skipPreFlight bool) (string, error) {
	if h.privateKey == nil {
		return "", ErrPrivateKeyNotFound
	}
	txBase64, err := transaction.SignTxWithPrivateKey(tx.Content, *h.privateKey)
	if err != nil {
		return "", err
	}

	response, err := h.PostSubmit(ctx, txBase64, skipPreFlight)
	if err != nil {
		return "", err
	}

	return response.Signature, nil
}

// signAndSubmitBatch signs the given transactions and submits them.
func (h *HTTPClient) signAndSubmitBatch(ctx context.Context, transactions []*pb.TransactionMessage, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	if h.privateKey == nil {
		return nil, ErrPrivateKeyNotFound
	}
	batchRequest, err := buildBatchRequest(transactions, *h.privateKey, opts)
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
func (h *HTTPClient) SubmitTradeSwap(ctx context.Context, owner, inToken, outToken string, inAmount, slippage float64, projectStr string, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	project, err := ProjectFromString(projectStr)
	if err != nil {
		return nil, err
	}
	resp, err := h.PostTradeSwap(ctx, owner, inToken, outToken, inAmount, slippage, project)
	if err != nil {
		return nil, err
	}
	return h.signAndSubmitBatch(ctx, resp.Transactions, opts)
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
	return h.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitDepositCollateral builds a deposit collateral transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitDepositCollateral(ctx context.Context, request *pb.PostDepositCollateralRequest, skipPreFlight bool) (string, error) {
	resp, err := h.PostDepositCollateral(ctx, request)
	if err != nil {
		return "", err
	}
	return h.signAndSubmit(ctx, &pb.TransactionMessage{
		Content: resp.Transaction,
	}, skipPreFlight)
}

// SubmitWithdrawCollateral builds a withdrawal collateral transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitWithdrawCollateral(ctx context.Context, request *pb.PostWithdrawCollateralRequest, skipPreFlight bool) (string, error) {
	resp, err := h.PostWithdrawCollateral(ctx, request)
	if err != nil {
		return "", err
	}
	return h.signAndSubmit(ctx, &pb.TransactionMessage{
		Content: resp.Transaction,
	}, skipPreFlight)
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

// PostPerpOrder returns a partially signed transaction for placing a perp order. Typically, you want to use SubmitPerpOrder instead of this.
func (h *HTTPClient) PostPerpOrder(ctx context.Context, request *pb.PostPerpOrderRequest) (*pb.PostPerpOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/perp/order", h.baseURL)

	var response pb.PostPerpOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostPerpOrderResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitPerpOrder builds a perp order, signs it, and submits to the network.
func (h *HTTPClient) SubmitPerpOrder(ctx context.Context, request *pb.PostPerpOrderRequest, opts PostOrderOpts) (string, error) {
	order, err := h.PostPerpOrder(ctx, request)
	if err != nil {
		return "", err
	}

	sig, err := h.signAndSubmit(ctx, &pb.TransactionMessage{
		Content: order.Transaction,
	}, opts.SkipPreFlight)
	return sig, err
}

// SubmitOrder builds a Serum market order, signs it, and submits to the network.
func (h *HTTPClient) SubmitOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error) {
	order, err := h.PostOrder(ctx, owner, payer, market, side, types, amount, price, project, opts)
	if err != nil {
		return "", err
	}

	sig, err := h.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
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

	return h.signAndSubmit(ctx, order.Transaction, skipPreFlight)
}

// PostClosePerpPositions builds cancel perp positions txn.
func (h *HTTPClient) PostClosePerpPositions(ctx context.Context, request *pb.PostClosePerpPositionsRequest) (*pb.PostClosePerpPositionsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/perp/close", h.baseURL)
	var response pb.PostClosePerpPositionsResponse
	err := connections.HTTPPostWithClient[*pb.PostClosePerpPositionsResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// SubmitClosePerpPositions builds a close perp positions txn, signs and submits it to the network.
func (h *HTTPClient) SubmitClosePerpPositions(ctx context.Context, request *pb.PostClosePerpPositionsRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	order, err := h.PostClosePerpPositions(ctx, request)
	if err != nil {
		return nil, err
	}

	var msgs []*pb.TransactionMessage
	for _, txn := range order.Transactions {
		msgs = append(msgs, &pb.TransactionMessage{Content: txn})
	}

	return h.signAndSubmitBatch(ctx, msgs, opts)
}

// SubmitCancelPerpOrder builds a cancels perp orders txn, signs and submits it to the network.
func (h *HTTPClient) SubmitCancelPerpOrder(ctx context.Context, request *pb.PostCancelPerpOrderRequest, skipPreFlight bool) (string, error) {
	order, err := h.PostCancelPerpOrder(ctx, request)
	if err != nil {
		return "", err
	}

	return h.signAndSubmit(ctx, &pb.TransactionMessage{
		Content: order.Transaction,
	}, skipPreFlight)
}

// SubmitCreateUser builds a create-user txn, signs and submits it to the network.
func (h *HTTPClient) SubmitCreateUser(ctx context.Context, request *pb.PostCreateUserRequest, skipPreFlight bool) (string, error) {
	order, err := h.PostCreateUser(ctx, request)
	if err != nil {
		return "", err
	}

	return h.signAndSubmit(ctx, &pb.TransactionMessage{
		Content: order.Transaction,
	}, skipPreFlight)
}

// SubmitPostPerpOrder builds a create-user txn, signs and submits it to the network.
func (h *HTTPClient) SubmitPostPerpOrder(ctx context.Context, request *pb.PostPerpOrderRequest, skipPreFlight bool) (string, error) {
	order, err := h.PostPerpOrder(ctx, request)
	if err != nil {
		return "", err
	}

	return h.signAndSubmit(ctx, &pb.TransactionMessage{
		Content: order.Transaction,
	}, skipPreFlight)
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

	return h.signAndSubmit(ctx, order.Transaction, skipPreFlight)
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
	return h.signAndSubmitBatch(ctx, orders.Transactions, opts)
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

	return h.signAndSubmit(ctx, order.Transaction, skipPreflight)
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

	return h.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
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

	return h.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
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

// GetPerpOrderbook returns the current state of perpetual contract orderbook.
func (h *HTTPClient) GetPerpOrderbook(ctx context.Context, request *pb.GetPerpOrderbookRequest) (*pb.GetPerpOrderbookResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/perp/%s?limit=%v&project=%v", request.Market, request.Limit, request.Project)
	orderbook := new(pb.GetPerpOrderbookResponse)
	if err := connections.HTTPGetWithClient[*pb.GetPerpOrderbookResponse](ctx, url, h.httpClient, orderbook, h.authHeader); err != nil {
		return nil, err
	}

	return orderbook, nil
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
