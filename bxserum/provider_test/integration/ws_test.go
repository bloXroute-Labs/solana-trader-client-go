package integration

import (
	"context"
	"github.com/bloXroute-Labs/serum-api/bxserum/provider"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestWSClient_Requests(t *testing.T) {
	w, err := provider.NewWSClient()
	require.Nil(t, err)

	TestGetOrderbook(
		t,
		func(ctx context.Context, market string, limit uint32) *pb.GetOrderbookResponse {
			orderbook, err := w.GetOrderbook(market, limit)
			require.Nil(t, err)

			return orderbook
		},
		func(ctx context.Context, market string, limit uint32) string {
			_, err := w.GetOrderbook(market, limit)
			require.NotNil(t, err)
			assert.Equal(t, err.Error(), "\"provided market name/address is not found\"")

			return "provided market name/address is not found"
		},
	)
}

func TestWSClient_Streams(t *testing.T) {
	w, err := provider.NewWSClient()
	require.Nil(t, err)

	TestGetOrderbookStream(
		t,
		func(ctx context.Context, market string, limit uint32, orderbookCh chan *pb.GetOrderbookStreamResponse) {
			err := w.GetOrderbookStream(ctx, market, limit, orderbookCh)
			require.Nil(t, err)
		},
	)
}
