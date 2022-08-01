package main

import (
	"fmt"
	"github.com/bloXroute-Labs/serum-client-go/benchmark/internal/arrival"
	gserum "github.com/gagliardetto/solana-go/programs/serum"
	"golang.org/x/exp/maps"
	"sort"
	"strconv"
	"strings"
	"time"
)

const tsFormat = "15:04:05.999999"

type Datapoint struct {
	Slot            int
	SerumTimestamps []time.Time
	SolanaAsk       time.Time
	SolanaBid       time.Time
}

// OrderedTimestamps returns timestamps pairs of [serumTS, solanaTS] within the slot. For example, the first TS from Serum is matched with the sooner of SolanaAsk and SolanaBid, the second with the later, and the rest with zero values.
func (d Datapoint) OrderedTimestamps() ([][]time.Time, gserum.Side) {
	var firstSide gserum.Side = gserum.SideAsk

	tsList := [][]time.Time{
		{time.Time{}, time.Time{}},
		{time.Time{}, time.Time{}},
	}
	if d.SolanaAsk.Before(d.SolanaBid) {
		tsList[0][1] = d.SolanaAsk
		tsList[1][1] = d.SolanaBid
	} else {
		tsList[1][1] = d.SolanaAsk
		tsList[0][1] = d.SolanaBid
		firstSide = gserum.SideBid
	}

	for i, ts := range d.SerumTimestamps {
		if i < 2 {
			tsList[i][0] = ts
		} else {
			tsList = append(tsList, []time.Time{ts, {}})
		}
	}

	return tsList, firstSide
}

func (d Datapoint) FormatPrint() []string {
	ots, firstSide := d.OrderedTimestamps()
	sides := []string{"ask", "bid"}
	if firstSide == gserum.SideBid {
		sides = []string{"bid", "ask"}
	}

	lines := make([]string, 0)
	for i, timestamps := range ots {
		side := "n/a"
		if i < len(sides) {
			side = sides[i]
		}

		line := fmt.Sprintf("slot %v (%v): serum %v, solana %v, diff %v", d.Slot, side, formatTS(timestamps[0]), formatTS(timestamps[1]), formatDiff(timestamps[0], timestamps[1]))
		lines = append(lines, line)
	}
	return lines
}

func (d Datapoint) FormatCSV() [][]string {
	ots, firstSide := d.OrderedTimestamps()
	sides := []string{"ask", "bid"}
	if firstSide == gserum.SideBid {
		sides = []string{"bid", "ask"}
	}

	lines := make([][]string, 0)
	for i, timestamps := range ots {
		side := "n/a"
		if i < len(sides) {
			side = sides[i]
		}

		line := []string{strconv.Itoa(d.Slot), formatDiff(timestamps[0], timestamps[1]), strconv.Itoa(i + 1), formatTS(timestamps[0]), side, formatTS(timestamps[1])}
		lines = append(lines, line)
	}

	return lines
}

// SlotRange enumerate the superset range of slots used in Serum and Solana updates
func SlotRange(serumResults map[int][]arrival.ProcessedUpdate[serumUpdate], solanaResults map[int][]arrival.ProcessedUpdate[solanaUpdate]) []int {
	serumSlots := maps.Keys(serumResults)
	sort.Ints(serumSlots)

	solanaSlots := maps.Keys(solanaResults)
	sort.Ints(solanaSlots)

	slots := make([]int, 0, len(serumResults))
	solanaIndex := 0
	for i := 0; i < len(serumSlots); i++ {
		serumCandidate := serumSlots[i]

		for j := solanaIndex; j < len(solanaSlots); j++ {
			solanaCandidate := solanaSlots[j]
			if solanaCandidate < serumCandidate {
				slots = append(slots, solanaCandidate)
				solanaIndex++
			} else if solanaCandidate == serumCandidate {
				solanaIndex++
			} else {
				break
			}
		}

		slots = append(slots, serumCandidate)
	}

	for j := solanaIndex; j < len(solanaSlots); j++ {
		slots = append(slots, solanaSlots[j])
	}

	return slots
}

// Merge combines Serum and Solana updates over the specified slots, indicating the difference in slot times and any updates that were not included in the other.
func Merge(slots []int, serumResults map[int][]arrival.ProcessedUpdate[serumUpdate], solanaResults map[int][]arrival.ProcessedUpdate[solanaUpdate]) ([]Datapoint, map[int][]arrival.ProcessedUpdate[serumUpdate], map[int][]arrival.ProcessedUpdate[solanaUpdate], error) {
	datapoints := make([]Datapoint, 0)
	leftoverSerum := make(map[int][]arrival.ProcessedUpdate[serumUpdate])
	leftoverSolana := make(map[int][]arrival.ProcessedUpdate[solanaUpdate])

	for _, slot := range slots {
		serumData, serumOK := serumResults[slot]
		solanaData, solanaOK := solanaResults[slot]

		if !serumOK && !solanaOK {
			return nil, nil, nil, fmt.Errorf("(slot %v) improper slot set provided", slot)
		}

		if !serumOK {
			leftoverSolana[slot] = solanaData
		}

		if !solanaOK {
			leftoverSerum[slot] = serumData
		}

		dp := Datapoint{
			Slot:            slot,
			SerumTimestamps: make([]time.Time, 0, len(serumData)),
		}
		if len(solanaData) > 2 {
			return nil, nil, nil, fmt.Errorf("(slot %v) solana data unexpectedly had more than 2 entries: %v", slot, solanaData)
		}
		for _, su := range solanaData {
			if su.Data.side == gserum.SideAsk {
				dp.SolanaAsk = su.Timestamp
			} else if su.Data.side == gserum.SideBid {
				dp.SolanaBid = su.Timestamp
			} else {
				return nil, nil, nil, fmt.Errorf("(slot %v) solana data unknown side: %v", slot, solanaData)
			}
		}

		for _, su := range serumData {
			dp.SerumTimestamps = append(dp.SerumTimestamps, su.Timestamp)
		}

		datapoints = append(datapoints, dp)
	}

	return datapoints, leftoverSerum, leftoverSolana, nil
}

func Print(datapoints []Datapoint, removeUnmatched bool) {
	fmt.Println("comparison points:")

	unmatchedCount := 0
	for _, dp := range datapoints {
		lines := dp.FormatPrint()
		for _, line := range lines {
			if removeUnmatched && strings.Contains(line, "n/a") {
				unmatchedCount++
				continue
			}
			fmt.Println(line)
		}
	}

	fmt.Println("skipped", unmatchedCount, "events without matches")
	fmt.Println()
}

func formatTS(ts time.Time) string {
	if ts.IsZero() {
		return "n/a"
	} else {
		return ts.Format(tsFormat)
	}
}

func formatDiff(serumTS time.Time, solanaTS time.Time) string {
	if serumTS.IsZero() || solanaTS.IsZero() {
		return "n/a"
	} else {
		return serumTS.Sub(solanaTS).String()
	}
}
