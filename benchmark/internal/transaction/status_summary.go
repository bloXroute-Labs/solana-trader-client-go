package transaction

type BatchSummary struct {
	Best            int
	LostTransaction []int
}

type BlockStatus struct {
	Slot     uint64
	Position int
	Found    bool
}
