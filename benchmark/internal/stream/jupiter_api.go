package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	defaultAPIEndpoint = "https://quote-api.jup.ag/v4/quote"
)

type jupiterAPIStream struct {
	ctx        context.Context
	client     *http.Client
	endpoint   string
	baseToken  string
	quoteToken string
}

func NewJupiterAPI(ctx context.Context) {

}

func (j *jupiterAPIStream) collectOne() {

}

func (j *jupiterAPIStream) query() (*jupiterResponse, error) {
	s := fmt.Sprintf("%v?inputMint=%v&outputMint=%v&amount=%v&slippageBps=%v")
	req, err := http.NewRequestWithContext(j.ctx, http.MethodGet, s, nil)
	if err != nil {
		return nil, err
	}

	res, err := j.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var jr jupiterResponse
	err = json.Unmarshal(b, &jr)
	if err != nil {
		return nil, err
	}

	return &jr, nil
}
