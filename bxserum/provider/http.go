package provider

import (
	"fmt"
	"github.com/bloXroute-Labs/serum-api/bxserum/connections"
	"github.com/bloXroute-Labs/serum-api/bxserum/transaction"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/bloXroute-Labs/serum-api/utils"
	"github.com/gagliardetto/solana-go"
	"net/http"
	"time"
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
	return NewHTTPClientWithTimeout(time.Second * 7)
}

// NewHTTPClientWithTimeout connects to Mainnet Serum API
func NewHTTPClientWithTimeout(timeout time.Duration) (*HTTPClient, error) {
	return NewHTTPClientWithEndpoint("http://174.129.154.164:1809", nil, timeout)
}

// NewHTTPTestnet connects to Testnet Serum API
func NewHTTPTestnet() (*HTTPClient, error) {
	panic("implement me")
}

// NewHTTPClientWithEndpoint connects to custom Serum API (set client to nil to use default client)
func NewHTTPClientWithEndpoint(endpoint string, client *http.Client, timeout time.Duration) (*HTTPClient, error) {
	if client == nil {
		client = &http.Client{Timeout: timeout}
	}

	privateKey, err := transaction.LoadPrivateKeyFromEnv()
	if err != nil {
		return nil, err
	}

	return &HTTPClient{
		baseURL:    endpoint,
		httpClient: client,
		privateKey: privateKey,
	}, nil
}

// SetPrivateKey sets the clients private key for signing orders
func (h *HTTPClient) SetPrivateKey(privateKey solana.PrivateKey) {
	h.privateKey = privateKey
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

// GetOrders returns all opened orders by owner address and market
func (h *HTTPClient) GetOrders(market string, owner string) (*pb.GetOrdersResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/orders/%s?address=%s", h.baseURL, market, owner)
	orders := new(pb.GetOrdersResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOrdersResponse](url, h.httpClient, orders); err != nil {
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

// PostOrder returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (h *HTTPClient) PostOrder(owner, payer, market string, side pb.Side, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/markets", h.baseURL)
	request := &pb.PostOrderRequest{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              []pb.OrderType{pb.OrderType_OT_LIMIT},
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
func (h *HTTPClient) SubmitOrder(owner, payer, market string, side pb.Side, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := h.PostOrder(owner, payer, market, side, amount, price, opts)
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
