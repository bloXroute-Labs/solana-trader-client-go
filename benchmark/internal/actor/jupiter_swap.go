package actor

import (
	"context"
	"errors"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/logger"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"go.uber.org/zap"
	"time"
)

const (
	defaultInterval       = time.Second
	defaultInitialTimeout = 2 * time.Second
	defaultSlippage       = 1
	defaultAmount         = 0.1
)

// places orders to affect liquidity and trigger price changes
type jupiterSwap struct {
	interval       time.Duration
	initialTimeout time.Duration
	inputMint      string
	outputMint     string
	amount         float64
	slippage       float64
	publicKey      string

	client *provider.HTTPClient
}

func NewJupiterSwap(opts ...JupiterOpt) (Liquidity, error) {
	j := &jupiterSwap{
		amount:         defaultAmount,
		interval:       defaultInterval,
		initialTimeout: defaultInitialTimeout,
		slippage:       defaultSlippage,
		client:         provider.NewHTTPClient(),
	}

	for _, o := range opts {
		o(j)
	}

	if j.inputMint == "" || j.outputMint == "" {
		return nil, errors.New("input and output mints are mandatory")
	}

	if j.publicKey == "" {
		return nil, errors.New("public key is mandatory")
	}

	return j, nil
}

func (j *jupiterSwap) log() *zap.SugaredLogger {
	return logger.Log().With("source", "jupiterActor")
}

func (j *jupiterSwap) Swap(ctx context.Context, iterations int) error {
	submitOpts := provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  true,
	}

	for i := 0; i < iterations; i++ {
		res, err := j.client.SubmitTradeSwap(ctx, j.publicKey, j.inputMint, j.outputMint, j.amount, j.slippage, pb.Project_P_JUPITER, submitOpts)
		if err != nil {
			return err
		}

		j.log().Infow("submitted transactions", "results", res.Transactions)

		time.Sleep(j.interval)
	}

	return nil
}
