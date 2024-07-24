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

func (d Datapoint) FormatPrint() []string {
	diff := formatDiff(d.TraderTimestamp, d.GeyserTimestamp)
	line := fmt.Sprintf("slot %v: trader API %v, geyser %v, diff %v",
		d.Slot, formatTS(d.TraderTimestamp), formatTS(d.GeyserTimestamp), diff)
	return []string{line}
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

		var traderUpdate stream.ProcessedUpdate[stream.TraderAPIUpdateGRPC]
		var geyserUpdate stream.ProcessedUpdate[stream.GeyserUpdateGRPC]

		if traderOK && len(traderUpdates) > 0 {
			traderUpdate = traderUpdates[len(traderUpdates)-1]
		}

		if geyserOK && len(geyserUpdates) > 0 {
			geyserUpdate = geyserUpdates[len(geyserUpdates)-1]
		}

		// Create a datapoint if we have updates from both sources
		if traderOK && geyserOK {
			dp := Datapoint{
				Slot:            slot,
				TraderTimestamp: traderUpdate.Timestamp,
				GeyserTimestamp: geyserUpdate.Timestamp,
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
	totalTraderLead := 0
	totalGeyserLead := 0
	traderUnmatched := 0
	geyserUnmatched := 0
	total := 0

	for _, dp := range datapoints {
		total++

		if dp.TraderTimestamp.IsZero() && dp.GeyserTimestamp.IsZero() {
			continue
		}

		if dp.TraderTimestamp.IsZero() {
			geyserUnmatched++
			continue
		}

		if dp.GeyserTimestamp.IsZero() {
			traderUnmatched++
			continue
		}

		if dp.TraderTimestamp.Before(dp.GeyserTimestamp) {
			traderFaster++
			totalTraderLead += int(dp.GeyserTimestamp.Sub(dp.TraderTimestamp).Milliseconds())
		} else if dp.GeyserTimestamp.Before(dp.TraderTimestamp) {
			geyserFaster++
			totalGeyserLead += int(dp.TraderTimestamp.Sub(dp.GeyserTimestamp).Milliseconds())
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

	fmt.Println("Run time: ", runtime)
	fmt.Println("Endpoints:")
	fmt.Println("    ", traderEndpoint, " [traderapi]")
	fmt.Println("    ", geyserEndpoint, " [geyser]")
	fmt.Println()

	fmt.Println("Total updates: ", total)
	fmt.Println()

	fmt.Println("Faster counts: ")
	fmt.Println(fmt.Sprintf("    %-6d  %v", traderFaster, traderEndpoint))
	fmt.Println(fmt.Sprintf("    %-6d  %v", geyserFaster, geyserEndpoint))

	fmt.Println("Average difference (ms): ")
	fmt.Println(fmt.Sprintf("    %-6s  %v", fmt.Sprintf("%vms", averageTraderLead), traderEndpoint))
	fmt.Println(fmt.Sprintf("    %-6s  %v", fmt.Sprintf("%vms", averageGeyserLead), geyserEndpoint))

	fmt.Println("Unmatched updates: ")
	fmt.Println("(updates from each stream without a corresponding result on the other)")
	fmt.Println(fmt.Sprintf("    %-6d  %v", traderUnmatched, traderEndpoint))
	fmt.Println(fmt.Sprintf("    %-6d  %v", geyserUnmatched, geyserEndpoint))
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
