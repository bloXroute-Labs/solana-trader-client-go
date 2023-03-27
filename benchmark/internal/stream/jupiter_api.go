package stream

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const (
	defaultAPIEndpoint = "https://quote-api.jup.ag/v4/quote"
	defaultAmount      = 100
	defaultSlippage    = 5
)

type jupiterAPIStream struct {
	ctx         context.Context
	client      *http.Client
	endpoint    string
	inputMint   string
	outputMint  string
	amount      float64
	slippageBps int64
}

type JupiterOpt func(s *jupiterAPIStream)

func NewJupiterAPI(ctx context.Context, opts ...JupiterOpt) (Source[JupiterResponse, JupiterResponse], error) {
	j := &jupiterAPIStream{
		ctx:         ctx,
		client:      &http.Client{},
		endpoint:    defaultAPIEndpoint,
		amount:      defaultAmount,
		slippageBps: defaultSlippage,
	}

	for _, o := range opts {
		o(j)
	}

	if j.inputMint == "" || j.outputMint == "" {
		return nil, errors.New("base and quote token are mandatory")
	}

	return j, nil
}

func WithJupiterEndpoint(endpoint string) JupiterOpt {
	return func(s *jupiterAPIStream) {
		s.endpoint = endpoint
	}
}

func WithJupiterTokenPair(inputMint, outputToken string) JupiterOpt {
	return func(s *jupiterAPIStream) {
		s.inputMint = inputMint
		s.outputMint = outputToken
	}
}

func WithJupiterAmount(amount float64) JupiterOpt {
	return func(s *jupiterAPIStream) {
		s.amount = amount
	}
}

func WithJupiterSlippage(slippage int64) JupiterOpt {
	return func(s *jupiterAPIStream) {
		s.slippageBps = slippage
	}
}

func (j *jupiterAPIStream) Name() string {
	return "jupiter"
}

func (j *jupiterAPIStream) Run(ctx context.Context) ([]RawUpdate[JupiterResponse], error) {
	// TODO: improve
	j.ctx = ctx

	res, err := j.fetch()
	if err != nil {
		return nil, err
	}

	return []RawUpdate[JupiterResponse]{NewRawUpdate(*res)}, nil
}

func (j *jupiterAPIStream) Process(updates []RawUpdate[JupiterResponse], removeDuplicates bool) (map[int][]ProcessedUpdate[JupiterResponse], map[int][]ProcessedUpdate[JupiterResponse], error) {
	// TODO implement me
	panic("implement me")
}

func (j *jupiterAPIStream) collectOne() {

}

func (j *jupiterAPIStream) fetch() (*JupiterResponse, error) {
	url := fmt.Sprintf("%v?inputMint=%v&outputMint=%v&amount=%v&slippageBps=%v", j.endpoint, j.inputMint, j.outputMint, j.amount, j.slippageBps)
	req, err := http.NewRequestWithContext(j.ctx, http.MethodGet, url, nil)
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

	var jr JupiterResponse
	err = json.Unmarshal(b, &jr)
	if err != nil {
		return nil, err
	}

	return &jr, nil
}
