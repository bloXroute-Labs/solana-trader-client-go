package integration

import (
	"context"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// Same as public.go in serum-api
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

	assert.Equal(t, "provided market name/address was not found", errMessage)

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
	assertMarketPresent(t, "MATH/USDT", pb.MarketStatus_MS_ONLINE, "CkvNfATB7nky8zPLuwS9bgcFbVRkQdkd5zuKEovyo9rs", markets)

	ctx, cancel = context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	assertMarketPresent(t, "SWAG/USDT", pb.MarketStatus_MS_ONLINE, "6URQ4zFWvPm1fhJCKKWorrh8X3mmTFiDDyXEUmSf8Rb2", markets)

	ctx, cancel = context.WithTimeout(context.Background(), publicRequestTimeout)
	defer cancel()
	assertMarketNotPresent(t, "market-doesnt-exist", markets)
}

func testGetOrders(t *testing.T, getOrdersFn func(ctx context.Context, market string, owner string) *pb.GetOrdersResponse) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	orders := getOrdersFn(ctx, "SOLUSDC", "AFT8VayE7qr8MoQsW3wHsDS83HhEvhGWdbNSHRKeUDfQ")
	assertOrder(t, "SOLUSDC", orders)

}

func testGetTickers(t *testing.T, getTickersFn func(ctx context.Context, market string) *pb.GetTickersResponse) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tickers := getTickersFn(ctx, "SOLUSDC")
	assertTickers(t, "SOLUSDC", tickers)

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

func assertOrder(t *testing.T, expectedName string, orders *pb.GetOrdersResponse) {
	require.NotEmpty(t, orders.GetOrders())
	order := *(orders.Orders[0])
	assert.Equal(t, expectedName, order.Market)
}

func assertTickers(t *testing.T, expectedName string, tickers *pb.GetTickersResponse) {
	require.NotEmpty(t, tickers.Tickers)
	ticker := *(tickers.Tickers[0])
	assert.Equal(t, expectedName, ticker.Market)
	assert.Greater(t, ticker.Bid, float64(0))
	assert.Greater(t, ticker.BidSize, float64(0))
	assert.Greater(t, ticker.Ask, float64(0))
	assert.Greater(t, ticker.AskSize, float64(0))
}

func assertTrades(t *testing.T, orders *pb.GetTradesResponse) {
	trades := orders.Trades
	require.NotEmpty(t, trades)
	trade := trades[0]
	assert.Greater(t, 0, trade.Size)
	assert.Greater(t, 0, trade.Price)
}
