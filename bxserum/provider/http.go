package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/bloXroute-Labs/serum-api/utils"
	"io/ioutil"
	"net/http"
	"time"
)

type HTTPClient struct {
	pb.UnsafeApiServer // TODO Regular API Server?
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
func NewHTTPClientWithEndpoint(baseURL string) *HTTPClient {
	client := http.Client{Timeout: time.Second * 7}
	return &HTTPClient{baseURL: baseURL, httpClient: &client} // TODO handle possible forward slash at end of base url?
}

func (w *HTTPClient) GetOrderbook(market string) (*pb.GetOrderbookResponse, error) {
	url := w.baseURL + fmt.Sprintf("/api/v1/market/orderbooks/%s", market)
	return getResponse[pb.GetOrderbookResponse](w.httpClient, url)
}

// response for GET request
func getResponse[T any](client *http.Client, url string) (*T, error) {
	if client == nil {
		return nil, errors.New("client is nil, please create one using a `NewHTTPClient` function")
	}
	httpResp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	return unmarshalResponse[T](httpResp)
}

func unmarshalResponse[T any](httpResp *http.Response) (*T, error) {
	if httpResp == nil {
		return nil, fmt.Errorf("HTTP response is nil")
	}

	b, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var resp T
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	return &resp, err
}
