package transaction

import "time"

type BatchSummary struct {
	Best            int
	LostTransaction []int
}

type BlockStatus struct {
	ExecutionTime time.Time
	Slot          uint64
	Position      int
	Found         bool
}
