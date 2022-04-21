package integration

import (
	"context"
	"github.com/bloXroute-Labs/serum-api/bxserum/provider"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHTTPClient_Requests(t *testing.T) {
	h := provider.NewHTTPClient()

	TestGetOrderbook(
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
}
