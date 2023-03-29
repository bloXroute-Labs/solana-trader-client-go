package stream

import "github.com/bloXroute-Labs/solana-trader-client-go/provider"

type TraderWSPriceOpt func(s *traderPriceStream)

func WithTraderWSClient(w *provider.WSClient) TraderWSPriceOpt {
	return func(s *traderPriceStream) {
		s.w = w
	}
}

func WithTraderWSMint(m string) TraderWSPriceOpt {
	return func(s *traderPriceStream) {
		s.mint = m
	}
}
