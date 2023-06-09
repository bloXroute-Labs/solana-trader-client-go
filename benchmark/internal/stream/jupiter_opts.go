package stream

import "time"

type JupiterOpt func(s *JupiterAPIStream)

func WithJupiterToken(mint string, decimals int) JupiterOpt {
	return func(s *JupiterAPIStream) {
		s.mint = mint
		s.decimals = decimals
	}
}

func WithJupiterTicker(t *time.Ticker) JupiterOpt {
	return func(s *JupiterAPIStream) {
		s.ticker = t
	}
}

func WithJupiterAmount(amount float64) JupiterOpt {
	return func(s *JupiterAPIStream) {
		s.amount = amount
	}
}

func WithJupiterSlippage(slippage int64) JupiterOpt {
	return func(s *JupiterAPIStream) {
		s.slippageBps = slippage
	}
}

func WithJupiterInterval(interval time.Duration) JupiterOpt {
	return func(s *JupiterAPIStream) {
		s.interval = interval
	}
}
