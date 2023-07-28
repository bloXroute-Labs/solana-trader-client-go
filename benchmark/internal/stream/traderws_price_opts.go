package stream

import "github.com/bloXroute-Labs/solana-trader-client-go/provider"

type TraderWSPriceOpt func(s *tradeWSPrice)

func WithTraderWSClient(w *provider.WSClient) TraderWSPriceOpt {
	return func(s *tradeWSPrice) {
		s.w = w
	}
}

func WithTraderWSMint(m string) TraderWSPriceOpt {
	return func(s *tradeWSPrice) {
		s.mint = m
	}
}
