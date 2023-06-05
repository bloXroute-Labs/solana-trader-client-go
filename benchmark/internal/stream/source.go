package stream

import (
	"context"
	"sync"
	"time"
)

// Source represents any streaming interface that provides timestamped updates for comparison.
type Source[T any, R any] interface {
	// Name returns an identifier for the source for printing
	Name() string

	// Run collects stream updates for the context duration. Run should avoid doing any other work besides collecting updates to have accurate timestamps.
	Run(context.Context) ([]RawUpdate[T], error)

	// Process deserializes the messages received by Run into useful formats for comparison.
	Process(updates []RawUpdate[T], removeDuplicates bool) (results map[int][]ProcessedUpdate[R], duplicates map[int][]ProcessedUpdate[R], err error)
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

type DurationUpdate[T any] struct {
	Start time.Time
	Data  T
}

func NewDurationUpdate[T any](start time.Time, data T) RawUpdate[DurationUpdate[T]] {
	return NewRawUpdate(DurationUpdate[T]{
		Start: start,
		Data:  data,
	})
}

type ProcessedUpdate[T any] struct {
	Timestamp time.Time
	Slot      int
	Data      T
}

func collectOrderedUpdates[T comparable](ctx context.Context, ticker *time.Ticker, requestFn func() (T, error), zero T, onError func(err error)) ([]RawUpdate[DurationUpdate[T]], error) {
	messages := make([]*RawUpdate[DurationUpdate[T]], 0)
	m := sync.Mutex{}
	for {
		select {
		case <-ticker.C:
			// to account for variances in response time, enforce an order
			go func() {
				m.Lock()

				start := time.Now()
				du := NewDurationUpdate[T](start, zero)
				duPtr := &du
				messages = append(messages, duPtr)

				m.Unlock()

				res, err := requestFn()
				if err != nil {
					onError(err)
					return
				}

				duPtr.Data.Data = res
				duPtr.Timestamp = time.Now()
			}()
		case <-ctx.Done():
			returnMessages := make([]RawUpdate[DurationUpdate[T]], 0)
			for _, message := range messages {
				if message.Data.Data != zero {
					returnMessages = append(returnMessages, *message)
				}
			}
			return returnMessages, nil
		}
	}
}
