package connections

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// HTTP response for GET request
func HTTPGetResponse[T any](client *http.Client, url string) (*T, error) {
	if client == nil {
		return nil, errors.New("client is nil, please create one using a `NewHTTPClient` function")
	}
	httpResp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	return httpUnmarshalResponse[T](httpResp)
}

func httpUnmarshalResponse[T any](httpResp *http.Response) (*T, error) {
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
