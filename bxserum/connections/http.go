package connections

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HTTPError struct {
	Code    int         `json:"code"`
	Details interface{} `json:"details"`
	Message string      `json:"message"`
}

func (h HTTPError) Error() string {
	return h.Message
}

// HTTP response for GET request
func HTTPGet[T any](client *http.Client, url string) (*T, error) { // TODO change func names HTTPGet
	if client == nil {
		return nil, errors.New("client is nil, please create one using a `NewHTTPClient` function")
	}
	httpResp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode != 200 {
		return nil, httpUnmarshalError(httpResp)
	}

	return httpUnmarshalResponse[T](httpResp)
}

// HTTP response for GET request
func HTTPGetWithClient[T any](url string, client *http.Client) (*T, error) {
	httpResp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	return httpUnmarshal[T](httpResp)
}

func httpUnmarshal(httpResp *http.Response) error {
	var httpError HTTPError
	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("error unmarshalling response to HTTPError") // TODO write better errors?
	}

	err = json.Unmarshal(body, &httpError)
	if err != nil {
		return err
	}

	return httpError
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
