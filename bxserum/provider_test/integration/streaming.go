package integration

import (
	"context"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	streamExpectEntries = 3
	streamExpectTimeout = 10 * time.Second
)

func TestGetOrderbookStream(
	t *testing.T,
	connectFn func(ctx context.Context, market string, limit uint32, orderbookCh chan *pb.GetOrderbookStreamResponse),
) {
	// no timeout: channel read timeouts are sufficient
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// normal subscription
	orderbookCh := make(chan *pb.GetOrderbookStreamResponse)
	go connectFn(ctx, "SOLUSDC", 0, orderbookCh)

	for i := 0; i < streamExpectEntries; i++ {
		orderbook := ReadChan[*pb.GetOrderbookStreamResponse](t, orderbookCh, streamExpectTimeout)
		require.NotNil(t, orderbook)

		assertSOLUSDCOrderbook(t, "SOL/USDC", orderbook.Orderbook)
	}
	cancel()
}

func ReadChan[T any](t *testing.T, c chan T, timeout time.Duration) T {
	select {
	case res := <-c:
		return res
	case <-time.After(timeout):
		assert.Fail(t, "no messages on channel")
	}
	return *new(T)
}
