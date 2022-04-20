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

func TestGetOrderbook(
	t *testing.T,
	getOrderbookFn func(ctx context.Context, market string, limit uint32) (*pb.GetOrderbookResponse, error),
) {

	// simple query
	ctx, cancel := context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	orderbook, err := getOrderbookFn(ctx, "ETHUSDC", 0)
	require.Nil(t, err)

	assertETHUSDCOrderbook(t, "ETHUSDC", orderbook)

	// try other appearances of orderbook format
	ctx, cancel = context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	orderbook, err = getOrderbookFn(ctx, "ETH-USDC", 0)
	require.Nil(t, err)

	assertETHUSDCOrderbook(t, "ETH-USDC", orderbook)

	// use limit
	ctx, cancel = context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	orderbook, err = getOrderbookFn(ctx, "ETH-USDC", 2)
	require.Nil(t, err)

	assertETHUSDCOrderbook(t, "ETH-USDC", orderbook)
	assert.Equal(t, 2, len(orderbook.Asks))
	assert.Equal(t, 2, len(orderbook.Bids))

	// another orderbook
	ctx, cancel = context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	orderbook, err = getOrderbookFn(ctx, "SOLUSDC", 2)
	require.Nil(t, err)

	assertOrderbook(t, "SOLUSDC", "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT", orderbook)
	assert.Equal(t, 2, len(orderbook.Asks))
	assert.Equal(t, 2, len(orderbook.Bids))

	// market doesn't exist
	ctx, cancel = context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	orderbook, err = getOrderbookFn(ctx, "market-doesnt-exist", 0)
	require.NotNil(t, err)
	assert.Equal(t, "\"provided market name/address is not found\"", err.Error())

	// another orderbook by address
	ctx, cancel = context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	orderbook, err = getOrderbookFn(ctx, "oTVAoRCiHnfEds5MTPerZk6VunEp24bCae8oSVrQmSU", 1)
	require.Nil(t, err)

	assert.Equal(t, 1, len(orderbook.Asks))
	assert.Equal(t, 1, len(orderbook.Bids))
}

func assertETHUSDCOrderbook(t *testing.T, name string, orderbook *pb.GetOrderbookResponse) {
	assertOrderbook(t, name, "4tSvZvnbyzHXLMTiFonMyxZoHmFqau1XArcRCVHLZ5gX", orderbook)
}

func assertOrderbook(t *testing.T, expectedName string, expectedAddress string, orderbook *pb.GetOrderbookResponse) {
	assert.Equal(t, expectedName, orderbook.Market)
	assert.Equal(t, expectedAddress, orderbook.MarketAddress)
	assert.Greater(t, len(orderbook.Bids), 1)
	assert.Greater(t, len(orderbook.Asks), 1)
}

func testGetMarkets(t *testing.T, assertMarketPresent func(t *testing.T, ctx context.Context, market string, addr string), assertMarketNotPresent func(t *testing.T, ctx context.Context, market string)) {
	ctx, cancel := context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	assertMarketPresent(t, ctx, "SOL/USDT", "HWHvQhFmJB3NUcu1aihKmrKegfVxBEHzwVX6yZCKEsi1")

	ctx, cancel = context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	assertMarketPresent(t, ctx, "MATH/USDT", "2WghiBkDL2yRhHdvm8CpprrkmfguuQGJTCDfPSudKBAZ")

	ctx, cancel = context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	assertMarketPresent(t, ctx, "SWAG/USDT", "J2XSt77XWim5HwtUM8RUwQvmRXNZsbMKpp5GTKpHafvf")

	ctx, cancel = context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	assertMarketNotPresent(t, ctx, "market-doesnt-exist")
}
