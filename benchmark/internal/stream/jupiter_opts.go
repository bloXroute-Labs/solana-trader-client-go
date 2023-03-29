package stream

import "time"

type JupiterOpt func(s *jupiterAPIStream)

func WithJupiterToken(mint string) JupiterOpt {
	return func(s *jupiterAPIStream) {
		s.mint = mint
	}
}

func WithJupiterTicker(t *time.Ticker) JupiterOpt {
	return func(s *jupiterAPIStream) {
		s.ticker = t
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

func WithJupiterInterval(interval time.Duration) JupiterOpt {
	return func(s *jupiterAPIStream) {
		s.interval = interval
	}
}
