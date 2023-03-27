package stream

import (
	"context"
	"time"
)

// Source represents any streaming interface that provides timestamped updates for comparison.
type Source[T any, R any] interface {
	// Name returns an identifier for the source for printing
	Name() string

	// Run collects stream updates for the context duration. Run should avoid doing any other work besides collecting updates, so as to have accurate timestamps.
	Run(context.Context) ([]RawUpdate[T], error)

	// Process deserializes the messages received by Run into useful formats for comparison.
	Process(updates []RawUpdate[T], removeDuplicates bool) (map[int][]ProcessedUpdate[R], map[int][]ProcessedUpdate[R], error)
}

type RawUpdate[T any] struct {
	Timestamp time.Time
	Data      T
}

func NewRawUpdate[T any](data T) RawUpdate[T] {
	return RawUpdate[T]{
		Timestamp: time.Now(),
		Data:      data,
	}
}

type ProcessedUpdate[T any] struct {
	Timestamp time.Time
	Slot      int
	Data      T
}
