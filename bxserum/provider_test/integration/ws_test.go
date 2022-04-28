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

	testGetOrderbook(
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

	testGetMarkets(
		t,
		func(ctx context.Context) *pb.GetMarketsResponse {
			markets, err := w.GetMarkets()
			require.Nil(t, err)

			return markets
		},
	)
}

// TODO separate WS streams
/*
func TestWSClient_Streams(t *testing.T) {
	w, err := provider.NewWSClient()
	require.Nil(t, err)

	testGetOrderbookStream(
		t,
		func(ctx context.Context, market string, limit uint32, orderbookCh chan *pb.GetOrderbookStreamResponse) {
			err := w.GetOrderbookStream(ctx, market, limit, orderbookCh)
			require.Nil(t, err)
		},
		func(ctx context.Context, market string, limit uint32) string {
			orderbookCh := make(chan *pb.GetOrderbookStreamResponse)
			err := w.GetOrderbookStream(ctx, market, limit, orderbookCh)
			require.NotNil(t, err)
			require.Equal(t, "\"provided market name/address is not found\"", err.Error())

			return "provided market name/address is not found"
		},
	)
}*/
