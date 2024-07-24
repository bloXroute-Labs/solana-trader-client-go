package stream

import (
	"context"
	pb "github.com/bloXroute-Labs/bxproto-private/go/solana-geyser-broadcast/v1"
	"github.com/bloXroute-Labs/bxprovider"
	"github.com/bloXroute-Labs/bxprovider/geyser"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/logger"
	"time"
)

type GeyserRawUpdateGRPC struct {
	Data *pb.Account
}

type GeyserUpdateGRPC struct {
	Account   *pb.Account
	Timestamp time.Time
}

type GeyserGRPCStream struct {
	source geyser.AccountsPatchSource
	pools  []string
}

func (s *GeyserGRPCStream) Name() string {
	return "GeyserGRPCStream"
}

func NewGeyserGRPCStream(ctx context.Context, grpcAddress string, pools []string) (*GeyserGRPCStream, error) {
	logger := new(bxprovider.StdLAdapter)
	source := geyser.NewAccountsPatchSource(ctx, grpcAddress, logger,
		geyser.WithAccountsFilter[*pb.AccountSubscribeResponse, *geyser.AccountsPatchRecv](pools))

	return &GeyserGRPCStream{
		source: source,
		pools:  pools,
	}, nil
}

func (s *GeyserGRPCStream) Run(ctx context.Context, updatesChan chan<- RawUpdate[GeyserRawUpdateGRPC]) error {
	updateChan, healthChan := s.source.Listen()

	for {
		select {
		case update, ok := <-updateChan:
			if !ok {
				return nil
			}
			updatesChan <- NewRawUpdate(GeyserRawUpdateGRPC{
				Data: update,
			})
		case health := <-healthChan:
			logger.Log().Infow("Geyser health update", "health", health)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *GeyserGRPCStream) Process(updates []RawUpdate[GeyserRawUpdateGRPC], removeDuplicates bool) (map[int][]ProcessedUpdate[GeyserUpdateGRPC], map[int][]ProcessedUpdate[GeyserUpdateGRPC], error) {
	results := make(map[int][]ProcessedUpdate[GeyserUpdateGRPC])
	duplicates := make(map[int][]ProcessedUpdate[GeyserUpdateGRPC])

	var lastUpdate *GeyserUpdateGRPC

	for _, update := range updates {
		slot := int(update.Data.Data.GetSlot())
		timestamp := update.Timestamp // You might want to use a more accurate timestamp if available

		currentUpdate := GeyserUpdateGRPC{
			Account:   update.Data.Data,
			Timestamp: timestamp,
		}

		pu := ProcessedUpdate[GeyserUpdateGRPC]{
			Timestamp: timestamp,
			Slot:      slot,
			Data:      currentUpdate,
		}

		if removeDuplicates && lastUpdate != nil {
			if isEqualGeyser(currentUpdate.Account, lastUpdate.Account) {
				duplicates[slot] = append(duplicates[slot], pu)
				continue
			}
		}

		results[slot] = append(results[slot], pu)
		lastUpdate = &currentUpdate
	}

	return results, duplicates, nil
}

func isEqualGeyser(a, b *pb.Account) bool {
	// Implement comparison logic for Geyser accounts
	// This is a simplified example, adjust according to your needs
	return a.GetPubkey() == b.GetPubkey() &&
		a.GetOwner() == b.GetOwner() &&
		a.GetLamports() == b.GetLamports()
}
