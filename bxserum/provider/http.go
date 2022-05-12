package provider

import (
	"fmt"
	"github.com/bloXroute-Labs/serum-api/bxserum/connections"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/bloXroute-Labs/serum-api/utils"
	"net/http"
	"time"
)

type HTTPClient struct {
	pb.UnimplementedApiServer

	baseURL    string
	httpClient *http.Client
	requestID  utils.RequestID
}

// Connects to Mainnet Serum API
func NewHTTPClient() *HTTPClient {
	return NewHTTPClientWithTimeout(time.Second * 7)
}

// Connects to Mainnet Serum API
func NewHTTPClientWithTimeout(timeout time.Duration) *HTTPClient {
	return NewHTTPClientWithEndpoint("http://174.129.154.164:1809", nil, timeout)
}

// Connects to Testnet Serum API
func NewHTTPTestnet() *HTTPClient {
	panic("implement me")
}

// Connects to custom Serum API (set client to nil to use default client)
func NewHTTPClientWithEndpoint(endpoint string, client *http.Client, timeout time.Duration) *HTTPClient {
	if client == nil {
		client = &http.Client{Timeout: time.Second * 7}
	}
	return &HTTPClient{baseURL: endpoint, httpClient: client}
}

// Set limit to 0 to get all bids/asks
func (h *HTTPClient) GetOrderbook(market string, limit uint32) (*pb.GetOrderbookResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/orderbooks/%s?limit=%v", h.baseURL, market, limit)
	orderbook := new(pb.GetOrderbookResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOrderbookResponse](url, h.httpClient, orderbook); err != nil {
		return nil, err
	}

	return orderbook, nil
}

// Set limit to 0 to get all trades
func (h *HTTPClient) GetTrades(market string, limit uint32) (*pb.GetTradesResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/trades/%s?limit=%v", h.baseURL, market, limit)
	marketTrades := new(pb.GetTradesResponse)
	if err := connections.HTTPGetWithClient[*pb.GetTradesResponse](url, h.httpClient, marketTrades); err != nil {
		return nil, err
	}

	return marketTrades, nil
}

// Set market to empty string to get all tickers
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

func (h *HTTPClient) GetMarkets() (*pb.GetMarketsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/markets", h.baseURL)
	markets := new(pb.GetMarketsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetMarketsResponse](url, h.httpClient, markets); err != nil {
		return nil, err
	}

	return markets, nil
}

func (h *HTTPClient) PostOrder() *pb.PostOrderRequest
