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
	return NewHTTPClientWithEndpoint("http://174.129.154.164:1809")
}

// Connects to Testnet Serum API
func NewHTTPTestnet() *HTTPClient {
	panic("implement me")
}

// Connects to custom Serum API
func NewHTTPClientWithEndpoint(endpoint string) *HTTPClient {
	client := http.Client{Timeout: time.Second * 7}            // TODO should we allow users to set the timeout?
	return &HTTPClient{baseURL: endpoint, httpClient: &client} // TODO handle possible forward slash at end of base url?
}

func (h *HTTPClient) GetOrderbook(market string) (*pb.GetOrderbookResponse, error) {
	url := h.baseURL + fmt.Sprintf("/api/v1/market/orderbooks/%s", market)
	return connections.HTTPGetResponse[pb.GetOrderbookResponse](h.httpClient, url)
}
