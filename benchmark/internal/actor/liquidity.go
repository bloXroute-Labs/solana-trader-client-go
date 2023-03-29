package actor

import "context"

type Liquidity interface {
	Swap(ctx context.Context, iterations int) error
}
