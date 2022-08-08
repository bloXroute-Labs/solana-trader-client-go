package integration

import (
	"context"
	"os"
	"testing"

	"github.com/bloXroute-Labs/serum-client-go/bxserum/provider"
	pb "github.com/bloXroute-Labs/serum-client-go/proto"
	"github.com/stretchr/testify/require"
)

func providerForEnv(t *testing.T) *provider.WSClient {
	env, ok := os.LookupEnv("SERUM_API_ENV")
	if !ok {
		env = "prod"
	}

	switch env {
	case "test":
		w, err := provider.NewWSClientTestnet()
		require.Nil(t, err)
		return w
	default: // test
		w, err := provider.NewWSClient()
		require.Nil(t, err)
		return w
	}
}

func TestWSClient_Requests(t *testing.T) {
	w := providerForEnv(t)
	defer func(w *provider.WSClient) {
		_ = w.Close()
	}(w)

	t.Run("orderbook", func(t *testing.T) {
		testGetOrderbook(
			t,
			func(ctx context.Context, market string, limit uint32) *pb.GetOrderbookResponse {
				orderbook, err := w.GetOrderbook(ctx, market, limit)
				require.Nil(t, err)

				return orderbook
			},
			func(ctx context.Context, market string, limit uint32) string {
				_, err := w.GetOrderbook(ctx, market, limit)
				require.NotNil(t, err)

				return err.Error()
			},
		)
	})

	t.Run("markets", func(t *testing.T) {
		testGetMarkets(
			t,
			func(ctx context.Context) *pb.GetMarketsResponse {
				markets, err := w.GetMarkets(ctx)
				require.Nil(t, err)

				return markets
			},
		)

	})

	t.Run("openOrders", func(t *testing.T) {
		testGetOpenOrders(
			t,
			func(ctx context.Context, market string, owner string) *pb.GetOpenOrdersResponse {
				orders, err := w.GetOpenOrders(ctx, market, owner)
				require.Nil(t, err)
				return orders
			},
		)

	})

	t.Run("unsettled", func(t *testing.T) {
		testUnsettled(
			t,
			func(ctx context.Context, market string, owner string) *pb.GetUnsettledResponse {
				response, err := w.GetUnsettled(ctx, market, owner)
				require.Nil(t, err)
				return response
			},
		)

	})

	t.Run("tickers", func(t *testing.T) {
		testGetTickers(
			t,
			func(ctx context.Context, market string) *pb.GetTickersResponse {
				_, _ = w.GetOrderbook(ctx, market, 1) // warmup tickers

				tickers, err := w.GetTickers(ctx, market)
				require.Nil(t, err)
				return tickers
			})
	})

	t.Run("submit order", func(t *testing.T) {
		testSubmitOrder(
			t,
			func(ctx context.Context, owner, payer, market string, side pb.Side, amount, price float64, opts provider.PostOrderOpts) string {
				txHash, err := w.SubmitOrder(ctx, owner, payer, market, side, []pb.OrderType{pb.OrderType_OT_LIMIT}, amount, price, opts)
				require.Nil(t, err, "unexpected error %v", err)
				return txHash
			},
			func(ctx context.Context, owner, payer, market string, side pb.Side, amount, price float64, opts provider.PostOrderOpts) string {
				_, err := w.SubmitOrder(ctx, owner, payer, market, side, []pb.OrderType{pb.OrderType_OT_LIMIT}, amount, price, opts)
				require.NotNil(t, err)

				return err.Error()
			})
	})
}

func TestWSClient_Streams(t *testing.T) {
	w := providerForEnv(t)
	defer func(w *provider.WSClient) {
		_ = w.Close()
	}(w)

	t.Run("orderbooks stream", func(t *testing.T) {
		testGetOrderbookStream(
			t,
			func(ctx context.Context, markets []string, limit uint32, orderbookCh chan *pb.GetOrderbooksStreamResponse) {
				_, err := w.GetOrderbooksStream(ctx, markets, limit)
				require.Nil(t, err)
			},
			func(ctx context.Context, markets []string, limit uint32) string {
				_, err := w.GetOrderbooksStream(ctx, markets, limit)
				require.NotNil(t, err)
				return err.Error()
			},
		)
	})

	t.Run("order status stream", func(t *testing.T) {
		testGetOrderStatusStream(
			t,
			func(ctx context.Context, market string, ownerAddress string) string {
				_, err := w.GetOrderStatusStream(ctx, market, ownerAddress)
				require.NotNil(t, err)

				return err.Error()
			},
		)
	})
}
