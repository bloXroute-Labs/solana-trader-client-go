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
	ticker      *time.Ticker
	interval    time.Duration
}

func NewJupiterAPI(opts ...JupiterOpt) (Source[DurationUpdate[JupiterPriceResponse], JupiterPriceResponse], error) {
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

func (j *jupiterAPIStream) Run(parent context.Context) ([]RawUpdate[DurationUpdate[JupiterPriceResponse]], error) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	ticker := j.ticker
	if ticker == nil {
		ticker = time.NewTicker(j.interval)
	}
	messages := make([]RawUpdate[DurationUpdate[JupiterPriceResponse]], 0)
	for {
		select {
		case <-ticker.C:
			go func() {
				start := time.Now()
				res, err := j.fetchPrice(ctx)
				if err != nil {
					return
				}

				messages = append(messages, NewDurationUpdate(start, res))
			}()
		case <-ctx.Done():
			return messages, nil
		}
	}
}

func (j *jupiterAPIStream) Process(updates []RawUpdate[DurationUpdate[JupiterPriceResponse]], removeDuplicates bool) (map[int][]ProcessedUpdate[JupiterPriceResponse], map[int][]ProcessedUpdate[JupiterPriceResponse], error) {
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

func (j *jupiterAPIStream) fetchPrice(ctx context.Context) (jr JupiterPriceResponse, err error) {
	url := fmt.Sprintf("%v?ids=%v", priceAPIEndpoint, j.mint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return
	}

	res, err := j.client.Do(req)
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	// TODO: move this to processing
	err = json.Unmarshal(b, &jr)
	if err != nil {
		return
	}

	return jr, nil
}
