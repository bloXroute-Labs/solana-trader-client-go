package integration

import (
	"context"
	"github.com/bloXroute-Labs/serum-api/bxserum/provider"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestHTTPClient_Requests(t *testing.T) {
	h := provider.NewHTTPClientWithTimeout(time.Second * 30)

	testGetOrderbook(
		t,
		func(ctx context.Context, market string, limit uint32) *pb.GetOrderbookResponse {
			orderbook, err := h.GetOrderbook(market, limit)
			require.Nil(t, err)

			return orderbook
		},
		func(ctx context.Context, market string, limit uint32) string {
			_, err := h.GetOrderbook(market, limit)
			require.NotNil(t, err)

			return err.Error()
		},
	)

	testGetMarkets(
		t,
		func(ctx context.Context) *pb.GetMarketsResponse {
			markets, err := h.GetMarkets()
			require.Nil(t, err)

			return markets
		},
	)

	testGetOrders(
		t,
		func(ctx context.Context, market string, owner string) *pb.GetOrdersResponse {
			orders, err := h.GetOrders(market, owner)
			require.Nil(t, err)
			return orders
		},
	)

	testGetTickers(
		t,
		func(ctx context.Context, market string) *pb.GetTickersResponse {
			h.GetOrderbook(market, 1)

			tickers, err := h.GetTickers(market)
			require.Nil(t, err)
			return tickers
		})
}
