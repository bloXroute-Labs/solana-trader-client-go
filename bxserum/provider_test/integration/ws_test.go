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

	t.Run("orderbook", func(t *testing.T) {
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
				assert.Equal(t, err.Error(), "\"provided market name/address was not found\"")

				return "provided market name/address was not found"
			},
		)
	})

	t.Run("markets", func(t *testing.T) {
		testGetMarkets(
			t,
			func(ctx context.Context) *pb.GetMarketsResponse {
				markets, err := w.GetMarkets()
				require.Nil(t, err)

				return markets
			},
		)

	})

	t.Run("orders", func(t *testing.T) {
		testGetOrders(
			t,
			func(ctx context.Context, market string, owner string) *pb.GetOrdersResponse {
				orders, err := w.GetOrders(market, owner)
				require.Nil(t, err)
				return orders
			},
		)

	})

	t.Run("tickers", func(t *testing.T) {
		testGetTickers(
			t,
			func(ctx context.Context, market string) *pb.GetTickersResponse {
				_, _ = w.GetOrderbook(market, 1) // warmup tickers

				tickers, err := w.GetTickers(market)
				require.Nil(t, err)
				return tickers
			})
	})

	t.Run("submit order", func(t *testing.T) {
		testSubmitOrder(
			t,
			func(ctx context.Context, owner, payer, market string, side pb.Side, amount, price float64, opts provider.PostOrderOpts) string {
				txHash, err := w.SubmitOrder(owner, payer, market, side, []pb.OrderType{pb.OrderType_OT_LIMIT}, amount, price, opts)
				require.Nil(t, err, "unexpected error %v", err)
				return txHash
			},
			func(ctx context.Context, owner, payer, market string, side pb.Side, amount, price float64, opts provider.PostOrderOpts) string {
				_, err := w.SubmitOrder(owner, payer, market, side, []pb.OrderType{pb.OrderType_OT_LIMIT}, amount, price, opts)
				require.NotNil(t, err)

				// TODO: fix websockets including extra quotes
				errString := err.Error()
				return errString[1 : len(errString)-1]
			})
	})
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
			require.Equal(t, "\"provided market name/address was not found\"", err.Error())

			return "provided market name/address was not found"
		},
	)
}*/
