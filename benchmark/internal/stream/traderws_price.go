package stream

import (
	"context"
	"errors"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/logger"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"go.uber.org/zap"
)

type tradeWSPrice struct {
	w    *provider.WSClient
	mint string
}

func NewTraderWSPrice(opts ...TraderWSPriceOpt) (Source[*pb.GetPricesStreamResponse, QuoteResult], error) {
	s := &tradeWSPrice{}

	for _, o := range opts {
		o(s)
	}

	if s.mint == "" {
		return nil, errors.New("mint is mandatory")
	}

	if s.w == nil {
		w, err := provider.NewWSClient()
		if err != nil {
			return nil, err
		}
		s.w = w
	}

	return s, nil
}

func (s tradeWSPrice) log() *zap.SugaredLogger {
	return logger.Log().With("source", "traderapi")
}

func (s tradeWSPrice) Name() string {
	return fmt.Sprintf("traderapi")
}

// Run stops when parent ctx is canceled
func (s tradeWSPrice) Run(parent context.Context) ([]RawUpdate[*pb.GetPricesStreamResponse], error) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	stream, err := s.w.GetPricesStream(ctx, []pb.Project{pb.Project_P_JUPITER}, []string{s.mint})
	if err != nil {
		return nil, err
	}

	ch := stream.Channel(10)

	messages := make([]RawUpdate[*pb.GetPricesStreamResponse], 0)
	for {
		select {
		case msg := <-ch:
			messages = append(messages, NewRawUpdate(msg))
		case <-ctx.Done():
			err = s.w.Close()
			if err != nil {
				s.log().Errorw("could not close connection", "err", err)
			}
			return messages, nil
		}
	}
}

func (s tradeWSPrice) Process(updates []RawUpdate[*pb.GetPricesStreamResponse], removeDuplicates bool) (results map[int][]ProcessedUpdate[QuoteResult], duplicates map[int][]ProcessedUpdate[QuoteResult], err error) {
	results = make(map[int][]ProcessedUpdate[QuoteResult])
	duplicates = make(map[int][]ProcessedUpdate[QuoteResult])

	lastBuyPrice := -1.
	lastSellPrice := -1.
	for _, update := range updates {
		slot := int(update.Data.Slot)
		buyPrice := update.Data.Price.Buy
		sellPrice := update.Data.Price.Sell

		qr := QuoteResult{
			Elapsed:   0,
			BuyPrice:  buyPrice,
			SellPrice: sellPrice,
			Source:    "traderWS",
		}
		pu := ProcessedUpdate[QuoteResult]{
			Timestamp: update.Timestamp,
			Slot:      slot,
			Data:      qr,
		}

		if buyPrice == lastBuyPrice && sellPrice == lastSellPrice {
			duplicates[slot] = append(duplicates[slot], pu)
			if removeDuplicates {
				continue
			}
		}

		lastBuyPrice = buyPrice
		lastSellPrice = sellPrice
		results[slot] = append(results[slot], pu)
	}

	return
}
