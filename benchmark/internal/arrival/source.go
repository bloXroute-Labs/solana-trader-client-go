package arrival

import (
	"context"
	"time"
)

// Source represents any streaming interface that provides timestamped updates for comparison.
type Source[T any, R any] interface {
	// Run collects stream updates for the context duration. Run should avoid doing any other work besides collecting updates, so as to have accurate timestamps.
	Run(context.Context) ([]StreamUpdate[T], error)

	// Process deserializes the messages received by Run into useful formats for comparison.
	Process(updates []StreamUpdate[T], removeDuplicates bool) (map[int][]ProcessedUpdate[R], map[int][]ProcessedUpdate[R], error)
}

type StreamUpdate[T any] struct {
	Timestamp time.Time
	Data      T
}

func NewStreamUpdate[T any](data T) StreamUpdate[T] {
	return StreamUpdate[T]{
		Timestamp: time.Now(),
		Data:      data,
	}
}

type ProcessedUpdate[T any] struct {
	Timestamp time.Time
	Slot      int
	Data      T
}
