package actor

import (
	"context"
	"errors"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"time"
)

const (
	defaultIterations = 1
	defaultInterval   = 5 * time.Second
	defaultAmount     = 0.1
)

// places orders to affect liquidity and trigger price changes
type jupiterSwap struct {
	iterations int
	interval   time.Duration
	inputMint  string
	outputMint string
	amount     float64
	owner      string

	client *provider.HTTPClient
}

func NewJupiterSwap(opts ...JupiterOpt) (Liquidity, error) {
	j := &jupiterSwap{
		amount:     defaultAmount,
		iterations: defaultIterations,
		interval:   defaultInterval,
		client:     provider.NewHTTPClient(),
	}

	for _, o := range opts {
		o(j)
	}

	if j.inputMint == "" || j.outputMint == "" {
		return nil, errors.New("input and output mints are mandatory")
	}

	return j, nil
}

func (j *jupiterSwap) Swap(ctx context.Context, iterations int) error {
	submitOpts := provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  false,
	}

	for i := 0; i < iterations; i++ {
		res, err := j.client.SubmitTradeSwap(ctx, j.owner, j.inputMint, j.outputMint, j.amount, 0, "jupiter", submitOpts)
		if err != nil {
			return err
		}

		fmt.Println(res)
	}

	return nil
}
