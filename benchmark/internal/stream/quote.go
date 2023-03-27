package stream

// QuoteSource requires some instrumentation
type QuoteSources[T any, R any] interface {
	// Name returns an identifier for the source for printing
	Name() string

	// Run collects stream updates for the context duration. Run should avoid doing any other work besides collecting updates, so as to have accurate timestamps.
	Run(int) ([]RawUpdate[T], error)

	// Process deserializes the messages received by Run into useful formats for comparison.
	Process(updates []RawUpdate[T], removeDuplicates bool) (map[int][]ProcessedUpdate[R], map[int][]ProcessedUpdate[R], error)
}
