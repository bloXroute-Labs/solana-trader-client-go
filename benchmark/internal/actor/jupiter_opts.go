package actor

import (
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	"time"
)

type JupiterOpt func(s *jupiterSwap)

func WithJupiterTokenPair(inputMint, outputMint string) JupiterOpt {
	return func(s *jupiterSwap) {
		s.inputMint = inputMint
		s.outputMint = outputMint
	}
}

func WithJupiterAmount(amount float64) JupiterOpt {
	return func(s *jupiterSwap) {
		s.amount = amount
	}
}

func WithJupiterInterval(interval time.Duration) JupiterOpt {
	return func(s *jupiterSwap) {
		s.interval = interval
	}
}

func WithJupiterIterations(iterations int) JupiterOpt {
	return func(s *jupiterSwap) {
		s.iterations = iterations
	}
}

func WithJupiterClient(client *provider.HTTPClient) JupiterOpt {
	return func(s *jupiterSwap) {
		s.client = client
	}
}
