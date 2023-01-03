package provider

import (
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/connections"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
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
		client = &http.Client{Timeout: opts.Timeout}
	}

	return &HTTPClient{
		baseURL:    opts.Endpoint,
		httpClient: client,
		privateKey: opts.PrivateKey,
		authHeader: opts.AuthHeader,
	}
}

// GetOrderbook returns the requested market's orderbook (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (h *HTTPClient) GetOrderbook(market string, limit uint32, project pb.Project) (*pb.GetOrderbookResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/orderbooks/%s?limit=%v&project=%v", h.baseURL, market, limit, project)
	orderbook := new(pb.GetOrderbookResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOrderbookResponse](url, h.httpClient, orderbook, h.authHeader); err != nil {
		return nil, err
	}

	return orderbook, nil
}

// GetMarketDepth returns the requested market's coalesced price data (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (h *HTTPClient) GetMarketDepth(market string, limit uint32, project pb.Project) (*pb.GetMarketDepthResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/depth/%s?limit=%v&project=%v", h.baseURL, market, limit, project)
	mktDepth := new(pb.GetMarketDepthResponse)
	if err := connections.HTTPGetWithClient[*pb.GetMarketDepthResponse](url, h.httpClient, mktDepth, h.authHeader); err != nil {
		return nil, err
	}

	return mktDepth, nil
}

// GetTrades returns the requested market's currently executing trades. Set limit to 0 for all trades.
func (h *HTTPClient) GetTrades(market string, limit uint32, project pb.Project) (*pb.GetTradesResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/trades/%s?limit=%v&project=%v", h.baseURL, market, limit, project)
	marketTrades := new(pb.GetTradesResponse)
	if err := connections.HTTPGetWithClient[*pb.GetTradesResponse](url, h.httpClient, marketTrades, h.authHeader); err != nil {
		return nil, err
	}

	return marketTrades, nil
}

// GetPools returns pools for given projects.
func (h *HTTPClient) GetPools(projects []pb.Project) (*pb.GetPoolsResponse, error) {
	projectsArg := convertSliceArgument("projects", true, projects)
	url := fmt.Sprintf("%s/api/v1/market/pools%s", h.baseURL, projectsArg)
	pools := new(pb.GetPoolsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetPoolsResponse](url, h.httpClient, pools, h.authHeader); err != nil {
		return nil, err
	}

	return pools, nil
}

// GetTickers returns the requested market tickets. Set market to "" for all markets.
func (h *HTTPClient) GetTickers(market string, project pb.Project) (*pb.GetTickersResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/tickers/%s?project=%v", h.baseURL, market, project)
	tickers := new(pb.GetTickersResponse)
	if err := connections.HTTPGetWithClient[*pb.GetTickersResponse](url, h.httpClient, tickers, h.authHeader); err != nil {
		return nil, err
	}

	return tickers, nil
}

// GetOpenOrders returns all opened orders by owner address and market
func (h *HTTPClient) GetOpenOrders(market string, owner string, openOrdersAddress string, project pb.Project) (*pb.GetOpenOrdersResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/openorders/%s?address=%s&openOrdersAddress=%s&project=%s", h.baseURL, market, owner, openOrdersAddress, project)
	orders := new(pb.GetOpenOrdersResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOpenOrdersResponse](url, h.httpClient, orders, h.authHeader); err != nil {
		return nil, err
	}

	return orders, nil
}

// GetMarkets returns the list of all available named markets
func (h *HTTPClient) GetMarkets() (*pb.GetMarketsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/markets", h.baseURL)
	markets := new(pb.GetMarketsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetMarketsResponse](url, h.httpClient, markets, h.authHeader); err != nil {
		return nil, err
	}

	return markets, nil
}

// GetUnsettled returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (h *HTTPClient) GetUnsettled(market string, owner string, project pb.Project) (*pb.GetUnsettledResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/unsettled/%s?ownerAddress=%s&project=%s", h.baseURL, market, owner, project)
	result := new(pb.GetUnsettledResponse)
	if err := connections.HTTPGetWithClient[*pb.GetUnsettledResponse](url, h.httpClient, result, h.authHeader); err != nil {
		return nil, err
	}

	return result, nil
}

// GetAccountBalance returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (h *HTTPClient) GetAccountBalance(owner string) (*pb.GetAccountBalanceResponse, error) {
	url := fmt.Sprintf("%s/api/v1/account/balance?ownerAddress=%s", h.baseURL, owner)
	result := new(pb.GetAccountBalanceResponse)
	if err := connections.HTTPGetWithClient[*pb.GetAccountBalanceResponse](url, h.httpClient, result, h.authHeader); err != nil {
		return nil, err
	}

	return result, nil
}

// GetPrice returns the USDC price of requested tokens
func (h *HTTPClient) GetPrice(tokens []string) (*pb.GetPriceResponse, error) {
	tokensArg := convertStrSliceArgument("tokens", true, tokens)
	url := fmt.Sprintf("%s/api/v1/market/price%s", h.baseURL, tokensArg)
	pools := new(pb.GetPriceResponse)
	if err := connections.HTTPGetWithClient[*pb.GetPriceResponse](url, h.httpClient, pools, h.authHeader); err != nil {
		return nil, err
	}

	return pools, nil
}

// GetQuotes returns the possible amount(s) of outToken for an inToken and the route to achieve it
func (h *HTTPClient) GetQuotes(inToken, outToken string, inAmount, slippage float64, limit int32, projects []pb.Project) (*pb.GetQuotesResponse, error) {
	projectString := convertSliceArgument("projects", false, projects)

	url := fmt.Sprintf("%s/api/v1/market/quote?inToken=%s&outToken=%s&inAmount=%v&slippage=%v&limit=%v%s",
		h.baseURL, inToken, outToken, inAmount, slippage, limit, projectString)
	result := new(pb.GetQuotesResponse)
	if err := connections.HTTPGetWithClient[*pb.GetQuotesResponse](url, h.httpClient, result, h.authHeader); err != nil {
		return nil, err
	}

	return result, nil
}

// PostSubmit posts the transaction string to the Solana network.
func (h *HTTPClient) PostSubmit(txBase64 string, skipPreFlight bool) (*pb.PostSubmitResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/submit", h.baseURL)
	request := &pb.PostSubmitRequest{Transaction: &pb.TransactionMessage{Content: txBase64}, SkipPreFlight: skipPreFlight}

	var response pb.PostSubmitResponse
	err := connections.HTTPPostWithClient[*pb.PostSubmitResponse](url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostSubmitBatch posts a bundle of transactions string based on a specific SubmitStrategy to the Solana network.
func (h *HTTPClient) PostSubmitBatch(request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/submit-batch", h.baseURL)

	var response pb.PostSubmitBatchResponse
	err := connections.HTTPPostWithClient[*pb.PostSubmitBatchResponse](url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// signAndSubmit signs the given transaction and submits it.
func (h *HTTPClient) signAndSubmit(tx *pb.TransactionMessage, skipPreFlight bool) (string, error) {
	if h.privateKey == nil {
		return "", ErrPrivateKeyNotFound
	}
	txBase64, err := transaction.SignTxWithPrivateKey(tx.Content, *h.privateKey)
	if err != nil {
		return "", err
	}

	response, err := h.PostSubmit(txBase64, skipPreFlight)
	if err != nil {
		return "", err
	}

	return response.Signature, nil
}

// signAndSubmitBatch signs the given transactions and submits them.
func (h *HTTPClient) signAndSubmitBatch(transactions []*pb.TransactionMessage, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	if h.privateKey == nil {
		return nil, ErrPrivateKeyNotFound
	}
	batchRequest, err := buildBatchRequest(transactions, *h.privateKey, opts)
	if err != nil {
		return nil, err
	}
	return h.PostSubmitBatch(batchRequest)
}

// PostTradeSwap PostOrder returns a partially signed transaction for submitting a swap request
func (h *HTTPClient) PostTradeSwap(ownerAddress, inToken, outToken string, inAmount, slippage float64, project pb.Project) (*pb.TradeSwapResponse, error) {
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
	err := connections.HTTPPostWithClient[*pb.TradeSwapResponse](url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitTradeSwap builds a TradeSwap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitTradeSwap(owner, inToken, outToken string, inAmount, slippage float64, projectStr string, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	project, err := ProjectFromString(projectStr)
	if err != nil {
		return nil, err
	}
	resp, err := h.PostTradeSwap(owner, inToken, outToken, inAmount, slippage, project)
	if err != nil {
		return nil, err
	}
	return h.signAndSubmitBatch(resp.Transactions, opts)
}

// PostRouteTradeSwap returns a partially signed transaction(s) for submitting a route swap request
func (h *HTTPClient) PostRouteTradeSwap(request *pb.RouteTradeSwapRequest) (*pb.TradeSwapResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/route-swap", h.baseURL)

	var response pb.TradeSwapResponse
	err := connections.HTTPPostWithClient[*pb.TradeSwapResponse](url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitRouteTradeSwap builds a RouteTradeSwap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitRouteTradeSwap(request *pb.RouteTradeSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := h.PostRouteTradeSwap(request)
	if err != nil {
		return nil, err
	}
	return h.signAndSubmitBatch(resp.Transactions, opts)
}

// PostOrder returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (h *HTTPClient) PostOrder(owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
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
	err := connections.HTTPPostWithClient[*pb.PostOrderResponse](url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitOrder builds a Serum market order, signs it, and submits to the network.
func (h *HTTPClient) SubmitOrder(owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error) {
	order, err := h.PostOrder(owner, payer, market, side, types, amount, price, project, opts)
	if err != nil {
		return "", err
	}

	sig, err := h.signAndSubmit(order.Transaction, opts.SkipPreFlight)
	return sig, err
}

// PostCancelOrder builds a Serum cancel order.
func (h *HTTPClient) PostCancelOrder(
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
	err := connections.HTTPPostWithClient[*pb.PostCancelOrderResponse](url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// SubmitCancelOrder builds a Serum cancel order, signs and submits it to the network.
func (h *HTTPClient) SubmitCancelOrder(
	orderID string,
	side pb.Side,
	owner,
	market,
	openOrders string,
	project pb.Project,
	skipPreFlight bool,
) (string, error) {
	order, err := h.PostCancelOrder(orderID, side, owner, market, openOrders, project)
	if err != nil {
		return "", err
	}

	return h.signAndSubmit(order.Transaction, skipPreFlight)
}

// PostCancelByClientOrderID builds a Serum cancel order by client ID.
func (h *HTTPClient) PostCancelByClientOrderID(
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
	err := connections.HTTPPostWithClient[*pb.PostCancelOrderResponse](url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// SubmitCancelByClientOrderID builds a Serum cancel order by client ID, signs and submits it to the network.
func (h *HTTPClient) SubmitCancelByClientOrderID(
	clientOrderID uint64,
	owner,
	market,
	openOrders string,
	project pb.Project,
	skipPreFlight bool,
) (string, error) {
	order, err := h.PostCancelByClientOrderID(clientOrderID, owner, market, openOrders, project)
	if err != nil {
		return "", err
	}

	return h.signAndSubmit(order.Transaction, skipPreFlight)
}

func (h *HTTPClient) PostCancelAll(market, owner string, openOrdersAddresses []string, project pb.Project) (*pb.PostCancelAllResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/cancelall", h.baseURL)
	request := &pb.PostCancelAllRequest{
		Market:              market,
		OwnerAddress:        owner,
		OpenOrdersAddresses: openOrdersAddresses,
		Project:             project,
	}

	var response pb.PostCancelAllResponse
	err := connections.HTTPPostWithClient[*pb.PostCancelAllResponse](url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (h *HTTPClient) SubmitCancelAll(market, owner string, openOrders []string, project pb.Project, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	orders, err := h.PostCancelAll(market, owner, openOrders, project)
	if err != nil {
		return nil, err
	}
	return h.signAndSubmitBatch(orders.Transactions, opts)
}

// PostSettle returns a partially signed transaction for settling market funds. Typically, you want to use SubmitSettle instead of this.
func (h *HTTPClient) PostSettle(owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, project pb.Project) (*pb.PostSettleResponse, error) {
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
	err := connections.HTTPPostWithClient[*pb.PostSettleResponse](url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitSettle builds a market SubmitSettle transaction, signs it, and submits to the network.
func (h *HTTPClient) SubmitSettle(owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, project pb.Project, skipPreflight bool) (string, error) {
	order, err := h.PostSettle(owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount, project)
	if err != nil {
		return "", err
	}

	return h.signAndSubmit(order.Transaction, skipPreflight)
}

func (h *HTTPClient) PostReplaceByClientOrderID(owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
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
	err := connections.HTTPPostWithClient[*pb.PostOrderResponse](url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (h *HTTPClient) SubmitReplaceByClientOrderID(owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error) {
	order, err := h.PostReplaceByClientOrderID(owner, payer, market, side, types, amount, price, project, opts)
	if err != nil {
		return "", err
	}

	return h.signAndSubmit(order.Transaction, opts.SkipPreFlight)
}

func (h *HTTPClient) PostReplaceOrder(orderID, owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
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
	err := connections.HTTPPostWithClient[*pb.PostOrderResponse](url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (h *HTTPClient) SubmitReplaceOrder(orderID, owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error) {
	order, err := h.PostReplaceOrder(orderID, owner, payer, market, side, types, amount, price, project, opts)
	if err != nil {
		return "", err
	}

	return h.signAndSubmit(order.Transaction, opts.SkipPreFlight)
}

// GetRecentBlockHash subscribes to a stream for getting recent block hash.
func (h *HTTPClient) GetRecentBlockHash() (*pb.GetRecentBlockHashResponse, error) {
	url := fmt.Sprintf("%s/api/v1/system/blockhash", h.baseURL)
	response := new(pb.GetRecentBlockHashResponse)
	if err := connections.HTTPGetWithClient[*pb.GetRecentBlockHashResponse](url, h.httpClient, response, h.authHeader); err != nil {
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
