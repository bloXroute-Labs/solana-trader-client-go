package output

import (
	"fmt"
	"golang.org/x/exp/maps"
	"sort"
)

func SortRange[T any](slotRange map[int]T) []int {
	slots := maps.Keys(slotRange)
	sort.Ints(slots)
	return slots
}

func FormatSortRange[T any](slotRange map[int]T) string {
	if len(slotRange) == 0 {
		return "-"
	}
	sr := SortRange(slotRange)
	return fmt.Sprintf("%v-%v", sr[0], sr[len(sr)-1])
}
