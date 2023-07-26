package actor

import (
	"context"
	"errors"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/logger"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	defaultInterval       = time.Second
	defaultInitialTimeout = 2 * time.Second
	defaultAfterTimeout   = 2 * time.Second
	defaultSlippage       = 1
	defaultAmount         = 0.1
)

// places orders to affect liquidity and trigger price changes
type jupiterSwap struct {
	interval       time.Duration
	initialTimeout time.Duration
	afterTimeout   time.Duration
	inputMint      string
	outputMint     string
	amount         float64
	slippage       float64
	publicKey      string
	alternate      bool

	client *provider.HTTPClient
}

func NewJupiterSwap(opts ...JupiterOpt) (Liquidity, error) {
	j := &jupiterSwap{
		amount:         defaultAmount,
		interval:       defaultInterval,
		initialTimeout: defaultInitialTimeout,
		afterTimeout:   defaultAfterTimeout,
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

func (j *jupiterSwap) Swap(ctx context.Context, iterations int) ([]SwapEvent, error) {
	submitOpts := provider.SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  true,
	}

	time.Sleep(j.initialTimeout)

	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	errCh := make(chan error, 1)
	resultCh := make(chan error, iterations)
	signatures := make([]SwapEvent, 0, iterations)
	signatureLock := &sync.Mutex{}
	lastOutAmount := 0.

	j.log().Infow("starting swap submission", "total", iterations)

	for i := 0; i < iterations; i++ {
		select {
		case <-ticker.C:
			go func(i int) {
				j.log().Infow("submitting swap", "count", i)

				var (
					inputMint  = j.inputMint
					outputMint = j.outputMint
					amount     = j.amount
				)

				if j.alternate && i%2 == 1 {
					inputMint, outputMint = outputMint, inputMint
					amount = lastOutAmount
				}

				info := fmt.Sprintf("%v => %v: %v", inputMint, outputMint, amount)

				postResponse, err := j.client.PostTradeSwap(ctx, j.publicKey, inputMint, outputMint, amount, j.slippage, pb.Project_P_JUPITER)
				if err != nil {
					errCh <- fmt.Errorf("error posting swap %v: %w", i, err)
					resultCh <- err
					return
				}

				// technically this can be a race condition, but shouldn't be a concern with the ticker times
				lastOutAmount = postResponse.OutAmount

				submitResponse, err := j.client.SignAndSubmitBatch(ctx, postResponse.Transactions, submitOpts)
				if err != nil {
					errCh <- fmt.Errorf("error submitting swap %v: %w", i, err)
					resultCh <- err
					return
				}

				j.log().Infow("completed swap", "transactions", submitResponse.Transactions)
				signatureLock.Lock()
				for _, transaction := range submitResponse.Transactions {
					signatures = append(signatures, SwapEvent{
						Timestamp: time.Now(),
						Signature: transaction.Signature,
						Info:      info,
					})
				}
				signatureLock.Unlock()
				resultCh <- nil

				time.Sleep(j.interval)
			}(i)
		case err := <-errCh:
			return signatures, err
		case <-ctx.Done():
			return signatures, errors.New("did not complete swaps before timeout")
		}
	}

	for i := 0; i < iterations; i++ {
		select {
		case err := <-resultCh:
			if err != nil {
				return signatures, err
			}
		case <-ctx.Done():
			return signatures, errors.New("did not complete swaps before timeout")
		}
	}

	time.Sleep(j.afterTimeout)
	return signatures, nil
}
