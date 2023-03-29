package stream

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	quoteAPIEndpoint = "https://quote-api.jup.ag/v4/quote"
	priceAPIEndpoint = "https://price.jup.ag/v4/price"
	defaultAmount    = 100
	defaultSlippage  = 5
	defaultInterval  = time.Second
)

type jupiterAPIStream struct {
	client      *http.Client
	mint        string
	amount      float64
	slippageBps int64
	interval    time.Duration
}

func NewJupiterAPI(opts ...JupiterOpt) (Source[JupiterPriceResponse, JupiterPriceResponse], error) {
	j := &jupiterAPIStream{
		client:      &http.Client{},
		amount:      defaultAmount,
		slippageBps: defaultSlippage,
		interval:    defaultInterval,
	}

	for _, o := range opts {
		o(j)
	}

	if j.mint == "" {
		return nil, errors.New("mint token is mandatory")
	}

	return j, nil
}

func (j *jupiterAPIStream) Name() string {
	return "jupiter"
}

func (j *jupiterAPIStream) Run(ctx context.Context) ([]RawUpdate[JupiterPriceResponse], error) {
	res, err := j.fetchPrice(ctx)
	if err != nil {
		return nil, err
	}

	return []RawUpdate[JupiterPriceResponse]{NewRawUpdate(*res)}, nil
}

func (j *jupiterAPIStream) Process(updates []RawUpdate[JupiterPriceResponse], removeDuplicates bool) (map[int][]ProcessedUpdate[JupiterPriceResponse], map[int][]ProcessedUpdate[JupiterPriceResponse], error) {
	// TODO implement me
	panic("implement me")
}

// func (j *jupiterAPIStream) fetchQuote(ctx context.Context) (*JupiterQuoteResponse, error) {
// 	url := fmt.Sprintf("%v?inputMint=%v&outputMint=%v&amount=%v&slippageBps=%v", quoteAPIEndpoint, j.inputMint, j.outputMint, j.amount, j.slippageBps)
// 	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	res, err := j.client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	defer func(Body io.ReadCloser) {
// 		_ = Body.Close()
// 	}(res.Body)
//
// 	b, err := io.ReadAll(res.Body)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// TODO: move this to processing
// 	var jr JupiterQuoteResponse
// 	err = json.Unmarshal(b, &jr)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &jr, nil
// }

func (j *jupiterAPIStream) fetchPrice(ctx context.Context) (*JupiterPriceResponse, error) {
	url := fmt.Sprintf("%v?ids=%v", priceAPIEndpoint, j.mint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

	// TODO: move this to processing
	var jr JupiterPriceResponse
	err = json.Unmarshal(b, &jr)
	if err != nil {
		return nil, err
	}

	return &jr, nil
}
