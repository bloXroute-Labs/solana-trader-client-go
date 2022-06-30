package integration

import (
	"context"
	"github.com/bloXroute-Labs/serum-client-go/bxserum/provider_test/bxassert" // different import than serum-api
	pb "github.com/bloXroute-Labs/serum-client-go/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// Very similar to stream.go in serum-api except for two differences
const (
	streamExpectEntries = 3
	streamExpectTimeout = 60 * time.Second // #1. longer timeout than serum-api
)

func testGetOrderbookStream(
	t *testing.T,
	connectFn func(ctx context.Context, market string, limit uint32, orderbookCh chan *pb.GetOrderbooksStreamResponse),
	connectFnErr func(ctx context.Context, market string, limit uint32) string,
) {
	// no timeout: channel read timeouts are sufficient
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// normal subscription
	orderbookCh := make(chan *pb.GetOrderbooksStreamResponse)
	go connectFn(ctx, "SOLUSDC", 0, orderbookCh)

	for i := 0; i < streamExpectEntries; i++ {
		orderbook := bxassert.ReadChan[*pb.GetOrderbooksStreamResponse](t, orderbookCh, streamExpectTimeout)
		require.NotNil(t, orderbook)

		assertSOLUSDCOrderbook(t, "SOL/USDC", orderbook.Orderbook)
	}
	cancel()

	// unknown market
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	errMessage := connectFnErr(ctx, "market-doesnt-exist", 0) // #2. cleaned up code here
	assert.Equal(t, "provided market name/address was not found", errMessage)
}

func testGetOrderStatusStream(t *testing.T, connectFnErr func(ctx context.Context, market string, ownerAddress string) string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// unknown market
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	errMessage := connectFnErr(ctx, "market-doesnt-exist", "FFqDwRq8B4hhFKRqx7N1M6Dg6vU699hVqeynDeYJdPj5") // #2. cleaned up code here
	assert.Equal(t, "provided market name/address was not found", errMessage)

	errMessage = connectFnErr(ctx, "SOLUSDC", "abcd")
	assert.Equal(t, "invalid len base58 public key string", errMessage)
}
