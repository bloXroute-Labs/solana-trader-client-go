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
	opts, err := provider.DefaultRPCOpts("http://174.129.154.164:1809")
	require.Nil(t, err)

	opts.Timeout = 60 * time.Second
	h := provider.NewHTTPClientWithOpts(nil, opts)

	t.Run("orderbook", func(t *testing.T) {
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
	})

	t.Run("markets", func(t *testing.T) {
		testGetMarkets(
			t,
			func(ctx context.Context) *pb.GetMarketsResponse {
				markets, err := h.GetMarkets()
				require.Nil(t, err)

				return markets
			},
		)
	})

	t.Run("openOrders", func(t *testing.T) {
		testGetOpenOrders(
			t,
			func(ctx context.Context, market string, owner string) *pb.GetOpenOrdersResponse {
				orders, err := h.GetOpenOrders(market, owner)
				require.Nil(t, err)
				return orders
			},
		)
	})

	t.Run("unsettled", func(t *testing.T) {
		testUnsettled(
			t,
			func(ctx context.Context, market string, owner string) *pb.GetUnsettledResponse {
				response, err := h.GetUnsettled(market, owner)
				require.Nil(t, err)
				return response
			},
		)
	})

	t.Run("tickers", func(t *testing.T) {
		testGetTickers(
			t,
			func(ctx context.Context, market string) *pb.GetTickersResponse {
				_, _ = h.GetOrderbook(market, 1) // warm up ticker

				tickers, err := h.GetTickers(market)
				require.Nil(t, err)
				return tickers
			})
	})

	t.Run("submit order", func(t *testing.T) {
		testSubmitOrder(
			t,
			func(ctx context.Context, owner, payer, market string, side pb.Side, amount, price float64, opts provider.PostOrderOpts) string {
				txHash, err := h.SubmitOrder(owner, payer, market, side, []pb.OrderType{pb.OrderType_OT_LIMIT}, amount, price, opts)
				require.Nil(t, err, "unexpected error %v", err)
				return txHash
			},
			func(ctx context.Context, owner, payer, market string, side pb.Side, amount, price float64, opts provider.PostOrderOpts) string {
				_, err := h.SubmitOrder(owner, payer, market, side, []pb.OrderType{pb.OrderType_OT_LIMIT}, amount, price, opts)
				require.NotNil(t, err)

				return err.Error()
			})
	})
}
