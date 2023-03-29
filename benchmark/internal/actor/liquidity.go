package actor

import (
	"context"
	"time"
)

type Liquidity interface {
	Swap(ctx context.Context, iterations int) ([]SwapEvent, error)
}

type SwapEvent struct {
	Timestamp time.Time
	Signature string
}
