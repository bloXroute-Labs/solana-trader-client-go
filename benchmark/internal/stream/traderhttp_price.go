package stream

import (
	"context"
	"errors"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/logger"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"go.uber.org/zap"
	"time"
)

type traderHTTPPriceStream struct {
	h        *provider.HTTPClient
	mint     string
	ticker   *time.Ticker
	interval time.Duration
}

func NewTraderHTTPPriceStream(opts ...TraderHTTPPriceOpt) (Source[DurationUpdate[*pb.GetPriceResponse], TraderAPIUpdate], error) {
	s := &traderHTTPPriceStream{
		h:        provider.NewHTTPClient(),
		interval: defaultInterval,
	}

	for _, o := range opts {
		o(s)
	}

	if s.mint == "" {
		return nil, errors.New("mint is mandatory")
	}

	return s, nil
}

func (s traderHTTPPriceStream) log() *zap.SugaredLogger {
	return logger.Log().With("source", "traderapi")
}

func (s traderHTTPPriceStream) Name() string {
	return fmt.Sprintf("traderapi")
}

// Run stops when parent ctx is canceled
func (s traderHTTPPriceStream) Run(parent context.Context) ([]RawUpdate[DurationUpdate[*pb.GetPriceResponse]], error) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	ticker := s.ticker
	if ticker == nil {
		ticker = time.NewTicker(s.interval)
	}
	messages := make([]RawUpdate[DurationUpdate[*pb.GetPriceResponse]], 0)
	for {
		select {
		case <-ticker.C:
			go func() {
				start := time.Now()
				price, err := s.h.GetPrice(ctx, []string{s.mint})
				if err != nil {
					return
				}
				messages = append(messages, NewDurationUpdate(start, price))
			}()
		case <-ctx.Done():
			return messages, nil
		}
	}
}

func (s traderHTTPPriceStream) Process(updates []RawUpdate[DurationUpdate[*pb.GetPriceResponse]], removeDuplicates bool) (map[int][]ProcessedUpdate[TraderAPIUpdate], map[int][]ProcessedUpdate[TraderAPIUpdate], error) {
	return nil, nil, nil
}
