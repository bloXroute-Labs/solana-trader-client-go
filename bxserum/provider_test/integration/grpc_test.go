package integration

import (
	"context"
	"testing"

	"github.com/bloXroute-Labs/serum-client-go/bxserum/provider"
	pb "github.com/bloXroute-Labs/serum-client-go/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"
)

// Unary response
func TestGRPCClient_Requests(t *testing.T) {
	g, err := provider.NewGRPCClient()
	require.Nil(t, err)

	t.Run("orderbook", func(t *testing.T) {
		testGetOrderbook(
			t,
			func(ctx context.Context, market string, limit uint32) *pb.GetOrderbookResponse {
				orderbook, err := g.GetOrderbook(ctx, market, limit)
				require.Nil(t, err)

				return orderbook
			},
			func(ctx context.Context, market string, limit uint32) string {
				_, err := g.GetOrderbook(ctx, market, limit)
				require.NotNil(t, err)

				grpcStatus, ok := status.FromError(err)
				require.True(t, ok)

				return grpcStatus.Message()
			},
		)
	})

	t.Run("markets", func(t *testing.T) {
		testGetMarkets(
			t,
			func(ctx context.Context) *pb.GetMarketsResponse {
				markets, err := g.GetMarkets(ctx)
				require.Nil(t, err)

				return markets
			},
		)
	})

	t.Run("openOrders", func(t *testing.T) {
		testGetOpenOrders(
			t,
			func(ctx context.Context, market string, owner string) *pb.GetOpenOrdersResponse {
				orders, err := g.GetOpenOrders(ctx, market, owner)
				require.Nil(t, err)
				return orders
			},
		)
	})

	t.Run("unsettled", func(t *testing.T) {
		testUnsettled(
			t,
			func(ctx context.Context, market string, owner string) *pb.GetUnsettledResponse {
				response, err := g.GetUnsettled(ctx, market, owner)
				require.Nil(t, err)
				return response
			},
		)
	})

	t.Run("tickers", func(t *testing.T) {
		testGetTickers(
			t,
			func(ctx context.Context, market string) *pb.GetTickersResponse {
				_, _ = g.GetOrderbook(ctx, market, 2) // warm up tickers
				response, err := g.GetTickers(ctx, market)

				require.Nil(t, err, "unexpected error=%v", err)
				return response
			})
	})

	t.Run("submit order", func(t *testing.T) {
		testSubmitOrder(
			t,
			func(ctx context.Context, owner, payer, market string, side pb.Side, amount, price float64, opts provider.PostOrderOpts) string {
				txHash, err := g.SubmitOrder(ctx, owner, payer, market, side, []pb.OrderType{pb.OrderType_OT_LIMIT}, amount, price, opts)
				require.Nil(t, err, "unexpected error %v", err)
				return txHash
			},
			func(ctx context.Context, owner, payer, market string, side pb.Side, amount, price float64, opts provider.PostOrderOpts) string {
				_, err := g.SubmitOrder(ctx, owner, payer, market, side, []pb.OrderType{pb.OrderType_OT_LIMIT}, amount, price, opts)
				require.NotNil(t, err)

				s, ok := status.FromError(err)
				require.True(t, ok)
				return s.Message()
			})
	})

}

// Stream response
func TestGRPCClient_Streams(t *testing.T) {
	g, err := provider.NewGRPCClient()
	require.Nil(t, err)

	testGetOrderbookStream(
		t,
		func(ctx context.Context, market string, limit uint32, orderbookCh chan *pb.GetOrderbooksStreamResponse) {
			err := g.GetOrderbookStream(ctx, market, limit, orderbookCh)
			require.Nil(t, err)
		},
		func(ctx context.Context, market string, limit uint32) string {
			orderbookCh := make(chan *pb.GetOrderbooksStreamResponse)
			err := g.GetOrderbookStream(ctx, market, limit, orderbookCh)
			require.NotNil(t, err)

			grpcStatus, ok := status.FromError(err)
			require.True(t, ok)

			return grpcStatus.Message()
		},
	)
}
