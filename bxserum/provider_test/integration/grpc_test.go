package integration

import (
	"context"
	"github.com/bloXroute-Labs/serum-api/bxserum/provider"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"
	"testing"
)

func TestGRPCClient_Requests(t *testing.T) {
	g, err := provider.NewGRPCClient()
	require.Nil(t, err)

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

	testGetMarkets(
		t,
		func(ctx context.Context) *pb.GetMarketsResponse {
			markets, err := g.GetMarkets(ctx)
			require.Nil(t, err)

			return markets
		},
	)
}

func TestGRPCClient_Streams(t *testing.T) {
	g, err := provider.NewGRPCClient()
	require.Nil(t, err)

	testGetOrderbookStream(
		t,
		func(ctx context.Context, market string, limit uint32, orderbookCh chan *pb.GetOrderbookStreamResponse) {
			err := g.GetOrderbookStream(ctx, market, limit, orderbookCh)
			require.Nil(t, err)
		},
		func(ctx context.Context, market string, limit uint32) string {
			orderbookCh := make(chan *pb.GetOrderbookStreamResponse)
			err := g.GetOrderbookStream(ctx, market, limit, orderbookCh)
			require.NotNil(t, err)

			return err.Error()
		},
	)
}