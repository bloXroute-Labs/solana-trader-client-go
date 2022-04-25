package integration

import (
	"context"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const publicRequestTimeout = 5 * time.Second

// if limit == -1: assume no limit
func testGetOrderbook(
	t *testing.T,
	getOrderbookFn func(ctx context.Context, market string, limit uint32) *pb.GetOrderbookResponse,
	getOrderbookErrFn func(ctx context.Context, market string, limit uint32) string,
) {

	// simple query
	ctx, cancel := context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	orderbook := getOrderbookFn(ctx, "SOLUSDC", 0)

	assertSOLUSDCOrderbook(t, "SOLUSDC", orderbook)

	// try other appearances of orderbook format
	ctx, cancel = context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	orderbook = getOrderbookFn(ctx, "SOL-USDC", 0)

	assertSOLUSDCOrderbook(t, "SOL-USDC", orderbook)

	// use limit
	ctx, cancel = context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	orderbook = getOrderbookFn(ctx, "SOL-USDC", 2)

	assertSOLUSDCOrderbook(t, "SOL-USDC", orderbook)
	assert.Equal(t, 2, len(orderbook.Asks))
	assert.Equal(t, 2, len(orderbook.Bids))

	// another orderbook
	ctx, cancel = context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	orderbook = getOrderbookFn(ctx, "SOLUSDC", 2)

	assertOrderbook(t, "SOLUSDC", "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT", orderbook)
	assert.Equal(t, 2, len(orderbook.Asks))
	assert.Equal(t, 2, len(orderbook.Bids))

	// market doesn't exist
	ctx, cancel = context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	errMessage := getOrderbookErrFn(ctx, "market-doesnt-exist", 0)

	assert.Equal(t, "provided market name/address is not found", errMessage)

	// another orderbook by address
	ctx, cancel = context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	orderbook = getOrderbookFn(ctx, "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT", 2)
	assertOrderbook(t, "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT", "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT", orderbook)
	assert.Equal(t, 2, len(orderbook.Asks))
	assert.Equal(t, 2, len(orderbook.Bids))
}

func testGetMarkets(t *testing.T, getMarketsFn func(ctx context.Context) *pb.GetMarketsResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	markets := getMarketsFn(ctx)
	assertMarketPresent(t, "SOL/USDT", pb.MarketStatus_MS_ONLINE, "HWHvQhFmJB3NUcu1aihKmrKegfVxBEHzwVX6yZCKEsi1", markets)

	ctx, cancel = context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	assertMarketPresent(t, "MATH/USDT", pb.MarketStatus_MS_ONLINE, "2WghiBkDL2yRhHdvm8CpprrkmfguuQGJTCDfPSudKBAZ", markets)

	ctx, cancel = context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	assertMarketPresent(t, "SWAG/USDT", pb.MarketStatus_MS_ONLINE, "J2XSt77XWim5HwtUM8RUwQvmRXNZsbMKpp5GTKpHafvf", markets)

	ctx, cancel = context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	assertMarketNotPresent(t, "market-doesnt-exist", markets)
}

func assertSOLUSDCOrderbook(t *testing.T, name string, orderbook *pb.GetOrderbookResponse) {
	assertOrderbook(t, name, "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT", orderbook)
}

func assertOrderbook(t *testing.T, expectedName string, expectedAddress string, orderbook *pb.GetOrderbookResponse) {
	require.NotNil(t, orderbook)

	assert.Equal(t, expectedName, orderbook.Market)
	assert.Equal(t, expectedAddress, orderbook.MarketAddress)
	assert.Greater(t, len(orderbook.Bids), 1)
	assert.Greater(t, len(orderbook.Asks), 1)
}

func assertMarketPresent(t *testing.T, marketName string, status pb.MarketStatus, address string, markets *pb.GetMarketsResponse) {
	require.NotNil(t, markets)

	market, ok := markets.Markets[marketName]
	require.True(t, ok)

	assert.Equal(t, marketName, market.Market)
	assert.Equal(t, status, market.Status)
	assert.Equal(t, address, market.Address)
}

func assertMarketNotPresent(t *testing.T, marketName string, markets *pb.GetMarketsResponse) {
	require.NotNil(t, markets)

	_, ok := markets.Markets[marketName]
	assert.False(t, ok)
}
