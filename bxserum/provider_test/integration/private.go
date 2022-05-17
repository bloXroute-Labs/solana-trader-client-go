package integration

import (
	"context"
	"github.com/bloXroute-Labs/serum-api/bxserum/provider"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"testing"
)

func testSubmitOrder(
	t *testing.T,
	submitOrderFn func(ctx context.Context, owner, payer, market string, side pb.Side, amount, price float64, opts provider.PostOrderOpts) string,
	submitOrderErrFn func(ctx context.Context, owner, payer, market string, side pb.Side, amount, price float64, opts provider.PostOrderOpts) string,
) {

}
