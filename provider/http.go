package provider

import (
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/connections"
	pb "github.com/bloXroute-Labs/solana-trader-client-go/proto"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
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

func (h *HTTPClient) GetAuthHeader() string {
	return h.authHeader
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
func (h *HTTPClient) GetOrderbook(market string, limit uint32) (*pb.GetOrderbookResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/orderbooks/%s?limit=%v", h.baseURL, market, limit)
	orderbook := new(pb.GetOrderbookResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOrderbookResponse](url, h.httpClient, orderbook, h.GetAuthHeader()); err != nil {
		return nil, err
	}

	return orderbook, nil
}

// GetTrades returns the requested market's currently executing trades. Set limit to 0 for all trades.
func (h *HTTPClient) GetTrades(market string, limit uint32) (*pb.GetTradesResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/trades/%s?limit=%v", h.baseURL, market, limit)
	marketTrades := new(pb.GetTradesResponse)
	if err := connections.HTTPGetWithClient[*pb.GetTradesResponse](url, h.httpClient, marketTrades, h.GetAuthHeader()); err != nil {
		return nil, err
	}

	return marketTrades, nil
}

// GetPools returns pools for given projects.
func (h *HTTPClient) GetPools(projects []string) (*pb.GetPoolsResponse, error) {
	projectsArg := ""
	if projects != nil && len(projects) > 0 {
		for i, project := range projects {
			arg := "projects=" + project
			if i == 0 {
				projectsArg = projectsArg + "?" + arg
			} else {
				projectsArg = projectsArg + "&" + arg
			}
		}
	}
	url := fmt.Sprintf("%s/api/v1/market/pools%s", h.baseURL, projectsArg)
	pools := new(pb.GetPoolsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetPoolsResponse](url, h.httpClient, pools, h.GetAuthHeader()); err != nil {
		return nil, err
	}

	return pools, nil
}

// GetTickers returns the requested market tickets. Set market to "" for all markets.
func (h *HTTPClient) GetTickers(market string) (*pb.GetTickersResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/tickers/%s", h.baseURL, market)
	tickers := new(pb.GetTickersResponse)
	if err := connections.HTTPGetWithClient[*pb.GetTickersResponse](url, h.httpClient, tickers, h.GetAuthHeader()); err != nil {
		return nil, err
	}

	return tickers, nil
}

// GetOpenOrders returns all opened orders by owner address and market
func (h *HTTPClient) GetOpenOrders(market string, owner string, openOrdersAddress string) (*pb.GetOpenOrdersResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/openorders/%s?address=%s&openOrdersAddress=%s", h.baseURL, market, owner, openOrdersAddress)
	orders := new(pb.GetOpenOrdersResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOpenOrdersResponse](url, h.httpClient, orders, h.GetAuthHeader()); err != nil {
		return nil, err
	}

	return orders, nil
}

// GetMarkets returns the list of all available named markets
func (h *HTTPClient) GetMarkets() (*pb.GetMarketsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/markets", h.baseURL)
	markets := new(pb.GetMarketsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetMarketsResponse](url, h.httpClient, markets, h.GetAuthHeader()); err != nil {
		return nil, err
	}

	return markets, nil
}

// GetUnsettled returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (h *HTTPClient) GetUnsettled(market string, owner string) (*pb.GetUnsettledResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/unsettled/%s?owner=%s", h.baseURL, market, owner)
	result := new(pb.GetUnsettledResponse)
	if err := connections.HTTPGetWithClient[*pb.GetUnsettledResponse](url, h.httpClient, result, h.GetAuthHeader()); err != nil {
		return nil, err
	}

	return result, nil
}

// GetAccountBalance returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (h *HTTPClient) GetAccountBalance(owner string) (*pb.GetAccountBalanceResponse, error) {
	url := fmt.Sprintf("%s/api/v1/account/balance?ownerAddress=%s", h.baseURL, owner)
	result := new(pb.GetAccountBalanceResponse)
	if err := connections.HTTPGetWithClient[*pb.GetAccountBalanceResponse](url, h.httpClient, result, h.GetAuthHeader()); err != nil {
		return nil, err
	}

	return result, nil
}

// GetPrice returns the USDC price of requested tokens
func (h *HTTPClient) GetPrice(tokens []string) (*pb.GetPriceResponse, error) {
	tokensArg := ""
	if tokens != nil && len(tokens) > 0 {
		for i, token := range tokens {
			arg := "tokens=" + token
			if i == 0 {
				tokensArg = tokensArg + "?" + arg
			} else {
				tokensArg = tokensArg + "&" + arg
			}
		}
	}
	url := fmt.Sprintf("%s/api/v1/market/price%s", h.baseURL, tokensArg)
	pools := new(pb.GetPriceResponse)
	if err := connections.HTTPGetWithClient[*pb.GetPriceResponse](url, h.httpClient, pools, h.GetAuthHeader()); err != nil {
		return nil, err
	}

	return pools, nil
}

// signAndSubmit signs the given transaction and submits it.
func (h *HTTPClient) signAndSubmit(tx string, skipPreFlight bool) (string, error) {
	if h.privateKey == nil {
		return "", ErrPrivateKeyNotFound
	}
	txBase64, err := transaction.SignTxWithPrivateKey(tx, *h.privateKey)
	if err != nil {
		return "", err
	}

	response, err := h.PostSubmit(txBase64, skipPreFlight)
	if err != nil {
		return "", err
	}

	return response.Signature, nil
}

// PostTradeSwap PostOrder returns a partially signed transaction for submitting a swap request
func (h *HTTPClient) PostTradeSwap(owner, inToken, outToken string, inAmount, slippage float64, project pb.Project) (*pb.TradeSwapResponse, error) {
	url := fmt.Sprintf("%s/api/v1/amm/trade-swap", h.baseURL)
	request := &pb.TradeSwapRequest{
		Owner:    owner,
		InToken:  inToken,
		OutToken: outToken,
		InAmount: inAmount,
		Slippage: slippage,
		Project:  project,
	}

	var response pb.TradeSwapResponse
	err := connections.HTTPPostWithClient[*pb.TradeSwapResponse](url, h.httpClient, request, &response, h.GetAuthHeader())
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostOrder returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (h *HTTPClient) PostOrder(owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/place", h.baseURL)
	request := &pb.PostOrderRequest{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              types,
		Amount:            amount,
		Price:             price,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
	}

	var response pb.PostOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostOrderResponse](url, h.httpClient, request, &response, h.GetAuthHeader())
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostSubmit posts the transaction string to the Solana network.
func (h *HTTPClient) PostSubmit(txBase64 string, skipPreFlight bool) (*pb.PostSubmitResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/submit", h.baseURL)
	request := &pb.PostSubmitRequest{Transaction: txBase64, SkipPreFlight: skipPreFlight}

	var response pb.PostSubmitResponse
	err := connections.HTTPPostWithClient[*pb.PostSubmitResponse](url, h.httpClient, request, &response, h.GetAuthHeader())
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitTradeSwap builds a TradeSwap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitTradeSwap(owner, inToken, outToken string, inAmount, slippage float64, projectStr string, skipPreFlight bool) ([]string, error) {
	project, err := ProjectFromString(projectStr)
	if err != nil {
		return []string{}, err
	}
	resp, err := h.PostTradeSwap(owner, inToken, outToken, inAmount, slippage, project)
	if err != nil {
		return []string{}, err
	}

	var signatures []string
	for _, tx := range resp.Transactions {
		signature, err := h.signAndSubmit(tx, skipPreFlight)
		if err != nil {
			return signatures, err
		}

		signatures = append(signatures, signature)
	}

	return signatures, nil
}

// SubmitOrder builds a Serum market order, signs it, and submits to the network.
func (h *HTTPClient) SubmitOrder(owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := h.PostOrder(owner, payer, market, side, types, amount, price, opts)
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
) (*pb.PostCancelOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/cancel", h.baseURL)
	request := &pb.PostCancelOrderRequest{
		OrderID:           orderID,
		Side:              side,
		OwnerAddress:      owner,
		MarketAddress:     market,
		OpenOrdersAddress: openOrders,
	}

	var response pb.PostCancelOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostCancelOrderResponse](url, h.httpClient, request, &response, h.GetAuthHeader())
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
	skipPreFlight bool,
) (string, error) {
	order, err := h.PostCancelOrder(orderID, side, owner, market, openOrders)
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
) (*pb.PostCancelOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/cancelbyid", h.baseURL)
	request := &pb.PostCancelByClientOrderIDRequest{
		ClientOrderID:     clientOrderID,
		OwnerAddress:      owner,
		MarketAddress:     market,
		OpenOrdersAddress: openOrders,
	}

	var response pb.PostCancelOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostCancelOrderResponse](url, h.httpClient, request, &response, h.GetAuthHeader())
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
	skipPreFlight bool,
) (string, error) {
	order, err := h.PostCancelByClientOrderID(clientOrderID, owner, market, openOrders)
	if err != nil {
		return "", err
	}

	return h.signAndSubmit(order.Transaction, skipPreFlight)
}

func (h *HTTPClient) PostCancelAll(market, owner string, openOrdersAddresses []string) (*pb.PostCancelAllResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/cancelall", h.baseURL)
	request := &pb.PostCancelAllRequest{
		Market:              market,
		OwnerAddress:        owner,
		OpenOrdersAddresses: openOrdersAddresses,
	}

	var response pb.PostCancelAllResponse
	err := connections.HTTPPostWithClient[*pb.PostCancelAllResponse](url, h.httpClient, request, &response, h.GetAuthHeader())
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (h *HTTPClient) SubmitCancelAll(market, owner string, openOrders []string, skipPreFlight bool) ([]string, error) {
	orders, err := h.PostCancelAll(market, owner, openOrders)
	if err != nil {
		return nil, err
	}

	var signatures []string
	for _, tx := range orders.Transactions {
		signature, err := h.signAndSubmit(tx, skipPreFlight)
		if err != nil {
			return signatures, err
		}

		signatures = append(signatures, signature)
	}

	return signatures, nil
}

// PostSettle returns a partially signed transaction for settling market funds. Typically, you want to use SubmitSettle instead of this.
func (h *HTTPClient) PostSettle(owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string) (*pb.PostSettleResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/settle", h.baseURL)
	request := &pb.PostSettleRequest{
		OwnerAddress:      owner,
		Market:            market,
		BaseTokenWallet:   baseTokenWallet,
		QuoteTokenWallet:  quoteTokenWallet,
		OpenOrdersAddress: openOrdersAccount,
	}

	var response pb.PostSettleResponse
	err := connections.HTTPPostWithClient[*pb.PostSettleResponse](url, h.httpClient, request, &response, h.GetAuthHeader())
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitSettle builds a market SubmitSettle transaction, signs it, and submits to the network.
func (h *HTTPClient) SubmitSettle(owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, skipPreflight bool) (string, error) {
	order, err := h.PostSettle(owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount)
	if err != nil {
		return "", err
	}

	return h.signAndSubmit(order.Transaction, skipPreflight)
}

func (h *HTTPClient) PostReplaceByClientOrderID(owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/replacebyclientid", h.baseURL)
	request := &pb.PostOrderRequest{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              types,
		Amount:            amount,
		Price:             price,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
	}

	var response pb.PostOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostOrderResponse](url, h.httpClient, request, &response, h.GetAuthHeader())
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (h *HTTPClient) SubmitReplaceByClientOrderID(owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := h.PostReplaceByClientOrderID(owner, payer, market, side, types, amount, price, opts)
	if err != nil {
		return "", err
	}

	return h.signAndSubmit(order.Transaction, opts.SkipPreFlight)
}

func (h *HTTPClient) PostReplaceOrder(orderID, owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/replace", h.baseURL)
	request := &pb.PostReplaceOrderRequest{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              types,
		Amount:            amount,
		Price:             price,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
		OrderID:           orderID,
	}

	var response pb.PostOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostOrderResponse](url, h.httpClient, request, &response, h.GetAuthHeader())
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (h *HTTPClient) SubmitReplaceOrder(orderID, owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := h.PostReplaceOrder(orderID, owner, payer, market, side, types, amount, price, opts)
	if err != nil {
		return "", err
	}

	return h.signAndSubmit(order.Transaction, opts.SkipPreFlight)
}

// GetRecentBlockHash subscribes to a stream for getting recent block hash.
func (h *HTTPClient) GetRecentBlockHash() (*pb.GetRecentBlockHashResponse, error) {
	url := fmt.Sprintf("%s/api/v1/system/blockhash", h.baseURL)
	response := new(pb.GetRecentBlockHashResponse)
	if err := connections.HTTPGetWithClient[*pb.GetRecentBlockHashResponse](url, h.httpClient, response, h.GetAuthHeader()); err != nil {
		return nil, err
	}

	return response, nil
}
