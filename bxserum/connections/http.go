package connections

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// HTTP response for GET request with default client
func HTTPGet[T any](url string) (*T, error) {
	client := &http.Client{Timeout: time.Second * 7}
	return HTTPGetWithClient[T](url, client)
}

// HTTP response for GET request
func HTTPGetWithClient[T any](url string, client *http.Client) (*T, error) {
	httpResp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	return httpUnmarshal[T](httpResp)
}

func httpUnmarshal[T any](httpResp *http.Response) (*T, error) {
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
