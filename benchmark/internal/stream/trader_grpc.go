package stream

import (
	"context"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/logger"
	"github.com/bloXroute-Labs/solana-trader-client-go/connections"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"time"
)

type TraderAPIRawUpdateGRPC struct {
	Data *pb.GetPoolReservesStreamResponse
}

type TraderAPIUpdateGRPC struct {
	Reserves  *pb.PoolReserves
	Timestamp time.Time
}

type TraderAPIGRPCStream struct {
	client *provider.GRPCClient
	pools  []string
}

func (s *TraderAPIGRPCStream) Name() string {
	return "TraderAPIGRPCStream"
}

func NewTraderAPIGRPCStream(grpcAddress, authHeader string, pools []string, useTLS bool) (*TraderAPIGRPCStream, error) {
	logger.Log().Infow("UseTLS: ", "TLS: ", useTLS)
	opts := provider.RPCOpts{
		Endpoint:   grpcAddress,
		AuthHeader: authHeader,
		UseTLS:     useTLS,
	}

	client, err := provider.NewGRPCClientWithOpts(opts)
	if err != nil {
		return nil, err
	}

	return &TraderAPIGRPCStream{
		client: client,
		pools:  pools,
	}, nil
}

func (s *TraderAPIGRPCStream) Run(ctx context.Context, updatesChan chan<- RawUpdate[TraderAPIRawUpdateGRPC], pools []string) error {
	request := &pb.GetPoolReservesStreamRequest{
		Projects: []pb.Project{pb.Project_P_RAYDIUM},
		Pools:    pools,
	}

	const maxRetries = 3
	var stream connections.Streamer[*pb.GetPoolReservesStreamResponse]
	var err error

	for attempt := 0; attempt < maxRetries; attempt++ {
		stream, err = s.client.GetPoolReservesStream(ctx, request)
		if err == nil {
			logger.Log().Infow("Established connection with trader api gRPC server", "stream", stream)
			break
		}

		logger.Log().Infow("Error connecting to trader api gRPC, retrying", "error", err, "attempt", attempt+1)

		if attempt < maxRetries-1 {
			time.Sleep(3 * time.Second * time.Duration(attempt+1)) // Simple linear backoff
		}
	}

	if err != nil {
		logger.Log().Infow("Failed to connect to trader api gRPC after retries", "error", err)
		return err
	}

	stream, err = s.client.GetPoolReservesStream(ctx, request)

	ch := make(chan *pb.GetPoolReservesStreamResponse)
	stream.Into(ch)

	for {
		select {
		case update, ok := <-ch:
			if !ok {
				return nil
			}
			updatesChan <- NewRawUpdate(TraderAPIRawUpdateGRPC{
				Data: update,
			})
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *TraderAPIGRPCStream) Process(updates []RawUpdate[TraderAPIRawUpdateGRPC], removeDuplicates bool) (map[int][]ProcessedUpdate[TraderAPIUpdateGRPC], map[int][]ProcessedUpdate[TraderAPIUpdateGRPC], error) {
	results := make(map[int][]ProcessedUpdate[TraderAPIUpdateGRPC])
	duplicates := make(map[int][]ProcessedUpdate[TraderAPIUpdateGRPC])

	var lastUpdate *TraderAPIUpdateGRPC

	for _, update := range updates {
		slot := int(update.Data.Data.Slot)
		timestamp := update.Timestamp

		currentUpdate := TraderAPIUpdateGRPC{
			Reserves:  update.Data.Data.Reserves,
			Timestamp: timestamp,
		}

		pu := ProcessedUpdate[TraderAPIUpdateGRPC]{
			Timestamp: timestamp,
			Slot:      slot,
			Data:      currentUpdate,
		}

		if removeDuplicates && lastUpdate != nil {
			if isEqual(currentUpdate.Reserves, lastUpdate.Reserves) {
				duplicates[slot] = append(duplicates[slot], pu)
				continue
			}
		}

		results[slot] = append(results[slot], pu)
		lastUpdate = &currentUpdate
	}

	return results, duplicates, nil
}

func isEqual(a, b *pb.PoolReserves) bool {
	return a.Token1Reserves == b.Token1Reserves &&
		a.Token1Address == b.Token1Address &&
		a.Token2Reserves == b.Token2Reserves &&
		a.Token2Address == b.Token2Address &&
		a.PoolAddress == b.PoolAddress &&
		a.Project == b.Project
}
