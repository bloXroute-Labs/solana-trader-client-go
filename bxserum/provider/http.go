package provider

import (
	"fmt"
	"net/http"

	"github.com/bloXroute-Labs/serum-api/bxserum/connections"
	"github.com/bloXroute-Labs/serum-api/bxserum/transaction"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/bloXroute-Labs/serum-api/utils"
	"github.com/gagliardetto/solana-go"
)

type HTTPClient struct {
	pb.UnimplementedApiServer

	baseURL    string
	httpClient *http.Client
	requestID  utils.RequestID
	privateKey solana.PrivateKey
}

// NewHTTPClient connects to Mainnet Serum API
func NewHTTPClient() (*HTTPClient, error) {
	opts, err := DefaultRPCOpts(MainnetSerumAPIHTTP)
	if err != nil {
		return nil, err
	}
	return NewHTTPClientWithOpts(nil, opts), nil
}

// NewHTTPTestnet connects to Testnet Serum API
func NewHTTPTestnet() (*HTTPClient, error) {
	opts, err := DefaultRPCOpts(TestnetSerumAPIHTTP)
	if err != nil {
		return nil, err
	}
	return NewHTTPClientWithOpts(nil, opts), nil
}

// NewHTTPClientWithOpts connects to custom Serum API (set client to nil to use default client)
func NewHTTPClientWithOpts(client *http.Client, opts RPCOpts) *HTTPClient {
	if client == nil {
		client = &http.Client{Timeout: opts.Timeout}
	}

	return &HTTPClient{
		baseURL:    opts.Endpoint,
		httpClient: client,
		privateKey: opts.PrivateKey,
	}
}

// GetOrderbook returns the requested market's orderbook (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (h *HTTPClient) GetOrderbook(market string, limit uint32) (*pb.GetOrderbookResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/orderbooks/%s?limit=%v", h.baseURL, market, limit)
	orderbook := new(pb.GetOrderbookResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOrderbookResponse](url, h.httpClient, orderbook); err != nil {
		return nil, err
	}

	return orderbook, nil
}

// GetTrades returns the requested market's currently executing trades. Set limit to 0 for all trades.
func (h *HTTPClient) GetTrades(market string, limit uint32) (*pb.GetTradesResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/trades/%s?limit=%v", h.baseURL, market, limit)
	marketTrades := new(pb.GetTradesResponse)
	if err := connections.HTTPGetWithClient[*pb.GetTradesResponse](url, h.httpClient, marketTrades); err != nil {
		return nil, err
	}

	return marketTrades, nil
}

// GetTickers returns the requested market tickets. Set market to "" for all markets.
func (h *HTTPClient) GetTickers(market string) (*pb.GetTickersResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/tickers/%s", h.baseURL, market)
	tickers := new(pb.GetTickersResponse)
	if err := connections.HTTPGetWithClient[*pb.GetTickersResponse](url, h.httpClient, tickers); err != nil {
		return nil, err
	}

	return tickers, nil
}

// GetOpenOrders returns all opened orders by owner address and market
func (h *HTTPClient) GetOpenOrders(market string, owner string) (*pb.GetOpenOrdersResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/openorders/%s?address=%s", h.baseURL, market, owner)
	orders := new(pb.GetOpenOrdersResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOpenOrdersResponse](url, h.httpClient, orders); err != nil {
		return nil, err
	}

	return orders, nil
}

// GetMarkets returns the list of all available named markets
func (h *HTTPClient) GetMarkets() (*pb.GetMarketsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/markets", h.baseURL)
	markets := new(pb.GetMarketsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetMarketsResponse](url, h.httpClient, markets); err != nil {
		return nil, err
	}

	return markets, nil
}

// GetUnsettled returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (h *HTTPClient) GetUnsettled(market string, owner string) (*pb.GetUnsettledResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/unsettled/%s?owner=%s", h.baseURL, market, owner)
	result := new(pb.GetUnsettledResponse)
	if err := connections.HTTPGetWithClient[*pb.GetUnsettledResponse](url, h.httpClient, result); err != nil {
		return nil, err
	}

	return result, nil
}

// signAndSubmit signs the given transaction and submits it.
func (h *HTTPClient) signAndSubmit(tx string) (string, error) {
	txBase64, err := transaction.SignTxWithPrivateKey(tx, h.privateKey)
	if err != nil {
		return "", err
	}

	response, err := h.PostSubmit(txBase64)
	if err != nil {
		return "", err
	}

	return response.Signature, nil
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
	err := connections.HTTPPostWithClient[*pb.PostOrderResponse](url, h.httpClient, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostSubmit posts the transaction string to the Solana network.
func (h *HTTPClient) PostSubmit(txBase64 string) (*pb.PostSubmitResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/submit", h.baseURL)
	request := &pb.PostSubmitRequest{Transaction: txBase64}

	var response pb.PostSubmitResponse
	err := connections.HTTPPostWithClient[*pb.PostSubmitResponse](url, h.httpClient, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitOrder builds a Serum market order, signs it, and submits to the network.
func (h *HTTPClient) SubmitOrder(owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (string, string, error) {
	order, err := h.PostOrder(owner, payer, market, side, types, amount, price, opts)
	if err != nil {
		return "", "", err
	}

	sig, err := h.signAndSubmit(order.Transaction)
	return sig, order.OpenOrdersAddress, err
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
	err := connections.HTTPPostWithClient[*pb.PostCancelOrderResponse](url, h.httpClient, request, &response)
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
) (string, error) {
	order, err := h.PostCancelOrder(orderID, side, owner, market, openOrders)
	if err != nil {
		return "", err
	}

	return h.signAndSubmit(order.Transaction)
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
	err := connections.HTTPPostWithClient[*pb.PostCancelOrderResponse](url, h.httpClient, request, &response)
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
) (string, error) {
	order, err := h.PostCancelByClientOrderID(clientOrderID, owner, market, openOrders)
	if err != nil {
		return "", err
	}

	return h.signAndSubmit(order.Transaction)
}

// PostSettle returns a partially signed transaction for settling market funds. Typically, you want to use SettleFunds instead of this.
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
	err := connections.HTTPPostWithClient[*pb.PostSettleResponse](url, h.httpClient, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SettleFunds builds a market SettleFunds transaction, signs it, and submits to the network.
func (h *HTTPClient) SettleFunds(owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string) (string, error) {
	order, err := h.PostSettle(owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount)
	if err != nil {
		return "", err
	}

	txBase64, err := transaction.SignTxWithPrivateKey(order.Transaction, h.privateKey)
	if err != nil {
		return "", err
	}

	response, err := h.PostSubmit(txBase64)
	if err != nil {
		return "", err
	}
	return response.Signature, nil
}
