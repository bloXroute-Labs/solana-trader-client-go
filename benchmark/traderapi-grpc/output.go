package main

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/stream"
)

const tsFormat = "15:04:05.999999"

type Datapoint struct {
	Slot            int
	TraderTimestamp time.Time
	GeyserTimestamp time.Time
}

func (d Datapoint) OrderedTimestamps() [][]time.Time {
	return [][]time.Time{
		{d.TraderTimestamp, d.GeyserTimestamp},
	}
}

func (d Datapoint) FormatCSV() [][]string {
	diff := formatDiff(d.TraderTimestamp, d.GeyserTimestamp)
	line := []string{
		strconv.Itoa(d.Slot),
		diff,
		formatTS(d.TraderTimestamp),
		formatTS(d.GeyserTimestamp),
	}
	return [][]string{line}
}

func SlotRange(traderResults map[int][]stream.ProcessedUpdate[stream.TraderAPIUpdateGRPC],
	geyserResults map[int][]stream.ProcessedUpdate[stream.GeyserUpdateGRPC]) []int {
	traderSlots := SortRange(traderResults)
	geyserSlots := SortRange(geyserResults)

	slots := make([]int, 0, len(traderResults)+len(geyserResults))
	geyserIndex := 0

	for i := 0; i < len(traderSlots); i++ {
		traderCandidate := traderSlots[i]

		for j := geyserIndex; j < len(geyserSlots); j++ {
			geyserCandidate := geyserSlots[j]
			if geyserCandidate < traderCandidate {
				slots = append(slots, geyserCandidate)
				geyserIndex++
			} else if geyserCandidate == traderCandidate {
				geyserIndex++
			} else {
				break
			}
		}

		slots = append(slots, traderCandidate)
	}

	for j := geyserIndex; j < len(geyserSlots); j++ {
		slots = append(slots, geyserSlots[j])
	}

	return slots
}

func Merge(slots []int,
	traderResults map[int][]stream.ProcessedUpdate[stream.TraderAPIUpdateGRPC],
	geyserResults map[int][]stream.ProcessedUpdate[stream.GeyserUpdateGRPC]) ([]Datapoint, map[int][]stream.ProcessedUpdate[stream.TraderAPIUpdateGRPC], map[int][]stream.ProcessedUpdate[stream.GeyserUpdateGRPC], error) {

	datapoints := make([]Datapoint, 0)
	leftoverTrader := make(map[int][]stream.ProcessedUpdate[stream.TraderAPIUpdateGRPC])
	leftoverGeyser := make(map[int][]stream.ProcessedUpdate[stream.GeyserUpdateGRPC])

	for _, slot := range slots {
		traderUpdates, traderOK := traderResults[slot]
		geyserUpdates, geyserOK := geyserResults[slot]

		if !traderOK && !geyserOK {
			return nil, nil, nil, fmt.Errorf("(slot %v) improper slot set provided", slot)
		}

		var latestTraderUpdate stream.ProcessedUpdate[stream.TraderAPIUpdateGRPC]
		var latestGeyserUpdate stream.ProcessedUpdate[stream.GeyserUpdateGRPC]

		if traderOK && len(traderUpdates) > 0 {
			latestTraderUpdate = traderUpdates[0]
			for _, update := range traderUpdates {
				if update.Timestamp.After(latestTraderUpdate.Timestamp) {
					latestTraderUpdate = update
				}
			}
		}

		if geyserOK && len(geyserUpdates) > 0 {
			latestGeyserUpdate = geyserUpdates[0]
			for _, update := range geyserUpdates {
				if update.Timestamp.After(latestGeyserUpdate.Timestamp) {
					latestGeyserUpdate = update
				}
			}
		}

		// Create a datapoint if we have updates from both sources
		if traderOK && geyserOK {
			dp := Datapoint{
				Slot:            slot,
				TraderTimestamp: latestTraderUpdate.Timestamp,
				GeyserTimestamp: latestGeyserUpdate.Timestamp,
			}
			datapoints = append(datapoints, dp)
		} else if traderOK {
			leftoverTrader[slot] = traderUpdates
		} else if geyserOK {
			leftoverGeyser[slot] = geyserUpdates
		}
	}

	return datapoints, leftoverTrader, leftoverGeyser, nil
}

func PrintSummary(runtime time.Duration, traderEndpoint string, geyserEndpoint string, datapoints []Datapoint) {
	traderFaster := 0
	geyserFaster := 0
	equalTimes := 0
	totalTraderLead := 0
	totalGeyserLead := 0
	totalDiff := 0
	minDiff := int(^uint(0) >> 1) // Max int value
	maxDiff := 0
	traderUnmatched := 0
	geyserUnmatched := 0
	totalTraderUpdates := 0
	totalGeyserUpdates := 0
	total := 0

	for _, dp := range datapoints {
		total++

		if dp.TraderTimestamp.IsZero() && dp.GeyserTimestamp.IsZero() {
			continue
		}

		if dp.TraderTimestamp.IsZero() {
			geyserUnmatched++
			totalGeyserUpdates++
			continue
		}

		if dp.GeyserTimestamp.IsZero() {
			traderUnmatched++
			totalTraderUpdates++
			continue
		}

		totalTraderUpdates++
		totalGeyserUpdates++

		diff := int(dp.TraderTimestamp.Sub(dp.GeyserTimestamp).Milliseconds())
		absDiff := abs(diff)
		totalDiff += absDiff

		if absDiff < minDiff {
			minDiff = absDiff
		}
		if absDiff > maxDiff {
			maxDiff = absDiff
		}

		if diff < 0 {
			traderFaster++
			totalTraderLead += -diff
		} else if diff > 0 {
			geyserFaster++
			totalGeyserLead += diff
		} else {
			equalTimes++
		}
	}

	averageTraderLead := 0
	if traderFaster > 0 {
		averageTraderLead = totalTraderLead / traderFaster
	}
	averageGeyserLead := 0
	if geyserFaster > 0 {
		averageGeyserLead = totalGeyserLead / geyserFaster
	}
	averageDiff := 0
	if total > 0 {
		averageDiff = totalDiff / total
	}

	fmt.Println("Run time: ", runtime)
	fmt.Println("Endpoints:")
	fmt.Println("    ", traderEndpoint, " [traderapi]")
	fmt.Println("    ", geyserEndpoint, " [geyser]")
	fmt.Println()

	fmt.Println("Total updates: ", total)
	fmt.Printf("Trader updates: %d\n", totalTraderUpdates)
	fmt.Printf("Geyser updates: %d\n", totalGeyserUpdates)
	fmt.Println()

	fmt.Println("Faster counts: ")
	fmt.Printf("    %-6d  %v\n", traderFaster, traderEndpoint)
	fmt.Printf("    %-6d  %v\n", geyserFaster, geyserEndpoint)
	fmt.Printf("    %-6d  Equal\n", equalTimes)

	fmt.Println("Average difference when faster (ms): ")
	fmt.Printf("    %-6d  %v\n", averageTraderLead, traderEndpoint)
	fmt.Printf("    %-6d  %v\n", averageGeyserLead, geyserEndpoint)

	fmt.Printf("Overall average difference: %d ms\n", averageDiff)
	fmt.Printf("Minimum difference: %d ms\n", minDiff)
	fmt.Printf("Maximum difference: %d ms\n", maxDiff)
	fmt.Println()

	fmt.Println("Unmatched updates: ")
	fmt.Println("(updates from each stream without a corresponding result on the other)")
	fmt.Printf("    %-6d  %v\n", traderUnmatched, traderEndpoint)
	fmt.Printf("    %-6d  %v\n", geyserUnmatched, geyserEndpoint)

	traderAvgSpeed := float64(runtime.Milliseconds()) / float64(totalTraderUpdates)
	geyserAvgSpeed := float64(runtime.Milliseconds()) / float64(totalGeyserUpdates)

	fmt.Println("\nAverage speed (ms per update):")
	fmt.Printf("    %-6.2f  %v\n", traderAvgSpeed, traderEndpoint)
	fmt.Printf("    %-6.2f  %v\n", geyserAvgSpeed, geyserEndpoint)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func SortRange[T any](updates map[int][]stream.ProcessedUpdate[T]) []int {
	slots := make([]int, 0, len(updates))
	for slot := range updates {
		slots = append(slots, slot)
	}
	sort.Ints(slots)
	return slots
}

func formatTS(ts time.Time) string {
	if ts.IsZero() {
		return "n/a"
	} else {
		return ts.Format(tsFormat)
	}
}

func formatDiff(traderTS time.Time, geyserTS time.Time) string {
	if traderTS.IsZero() || geyserTS.IsZero() {
		return "n/a"
	} else {
		return traderTS.Sub(geyserTS).String()
	}
}
