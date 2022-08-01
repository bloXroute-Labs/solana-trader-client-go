package utils

import (
	"context"
	"golang.org/x/sync/errgroup"
)

type result[T any] struct {
	index int
	value T
}

func AsyncGather[T any, R any](ctx context.Context, inputs []T, apply func(int, context.Context, T) (R, error)) ([]R, error) {
	ch := make(chan result[R], len(inputs))
	group, groupCtx := errgroup.WithContext(ctx)

	for i, input := range inputs {
		i := i
		input := input

		group.Go(func() error {
			v, err := apply(i, groupCtx, input)
			if err != nil {
				return err
			}

			ch <- result[R]{
				index: i,
				value: v,
			}
			return nil
		})
	}

	err := group.Wait()
	if err != nil {
		return nil, err
	}

	outputs := make([]R, len(inputs))
	for range inputs {
		r := <-ch
		outputs[r.index] = r.value
	}
	return outputs, nil
}
