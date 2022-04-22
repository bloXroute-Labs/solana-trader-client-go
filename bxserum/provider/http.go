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
	return NewHTTPClientWithEndpoint("http://174.129.154.164:1809", nil)
}

// Connects to Testnet Serum API
func NewHTTPTestnet() *HTTPClient {
	panic("implement me")
}

// Connects to custom Serum API (set client to nil to use default client)
func NewHTTPClientWithEndpoint(endpoint string, client *http.Client) *HTTPClient {
	if client == nil {
		client = &http.Client{Timeout: time.Second * 7}
	}
	return &HTTPClient{baseURL: endpoint, httpClient: client}
}

// Set limit to 0 to get all bids/asks
func (h *HTTPClient) GetOrderbook(market string, limit uint32) (*pb.GetOrderbookResponse, error) {
	url := h.baseURL + fmt.Sprintf("/api/v1/market/orderbooks/%s", market)
	if limit != 0 {
		url += fmt.Sprintf("?limit=%v", limit)
	}
	orderbook := new(pb.GetOrderbookResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOrderbookResponse](url, h.httpClient, orderbook); err != nil {
		return nil, err
	}

	return orderbook, nil
}

func (h *HTTPClient) GetMarkets() (*pb.GetMarketsResponse, error) {
	url := h.baseURL + fmt.Sprintf("/api/v1/market/markets")
	markets := new(pb.GetMarketsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetMarketsResponse](url, h.httpClient, markets); err != nil {
		return nil, err
	}

	return markets, nil
}
