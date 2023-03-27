package stream

// QuoteSource requires some instrumentation
type QuoteSource[T any, R any] interface {
	Source[T, R]
}
