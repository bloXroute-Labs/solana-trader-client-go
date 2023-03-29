package stream

import "github.com/bloXroute-Labs/solana-trader-client-go/provider"

type TraderHTTPPriceOpt func(s *traderHTTPPriceStream)

func WithTraderHTTPClient(h *provider.HTTPClient) TraderHTTPPriceOpt {
	return func(s *traderHTTPPriceStream) {
		s.h = h
	}
}

func WithTraderHTTPMint(m string) TraderHTTPPriceOpt {
	return func(s *traderHTTPPriceStream) {
		s.mint = m
	}
}
