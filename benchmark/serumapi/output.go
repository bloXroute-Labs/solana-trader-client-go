package main

import (
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/arrival"
	gserum "github.com/gagliardetto/solana-go/programs/serum"
	"golang.org/x/exp/maps"
	"sort"
	"strconv"
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

	setBidFirst := func() {
		tsList[1][1] = d.SolanaAsk
		tsList[0][1] = d.SolanaBid
		firstSide = gserum.SideBid
	}

	setAskFirst := func() {
		tsList[0][1] = d.SolanaAsk
		tsList[1][1] = d.SolanaBid
	}

	// zeros should come at end
	if d.SolanaAsk.IsZero() && !d.SolanaBid.IsZero() {
		setBidFirst()
	} else if !d.SolanaAsk.IsZero() && d.SolanaBid.IsZero() {
		setAskFirst()
	} else if d.SolanaBid.Before(d.SolanaAsk) {
		setBidFirst()
	} else {
		setAskFirst()
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

func SortRange[T any](slotRange map[int]T) []int {
	slots := maps.Keys(slotRange)
	sort.Ints(slots)
	return slots
}

func FormatSortRange[T any](slotRange map[int]T) string {
	sr := SortRange(slotRange)
	return fmt.Sprintf("%v-%v", sr[0], sr[len(sr)-1])
}

// SlotRange enumerate the superset range of slots used in Serum and Solana updates
func SlotRange(serumResults map[int][]arrival.ProcessedUpdate[arrival.SerumUpdate], solanaResults map[int][]arrival.ProcessedUpdate[arrival.SolanaUpdate]) []int {
	serumSlots := SortRange(serumResults)
	solanaSlots := SortRange(solanaResults)

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
func Merge(slots []int, serumResults map[int][]arrival.ProcessedUpdate[arrival.SerumUpdate], solanaResults map[int][]arrival.ProcessedUpdate[arrival.SolanaUpdate]) ([]Datapoint, map[int][]arrival.ProcessedUpdate[arrival.SerumUpdate], map[int][]arrival.ProcessedUpdate[arrival.SolanaUpdate], error) {
	datapoints := make([]Datapoint, 0)
	leftoverSerum := make(map[int][]arrival.ProcessedUpdate[arrival.SerumUpdate])
	leftoverSolana := make(map[int][]arrival.ProcessedUpdate[arrival.SolanaUpdate])

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
			if su.Data.Side == gserum.SideAsk {
				dp.SolanaAsk = su.Timestamp
			} else if su.Data.Side == gserum.SideBid {
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

func PrintSummary(runtime time.Duration, serumEndpoint string, solanaEndpoint string, datapoints []Datapoint) {
	serumFaster := 0
	solanaFaster := 0
	totalSerumLead := 0
	totalSolanaLead := 0
	serumUnmatched := 0
	solanaUnmatched := 0
	total := 0

	for _, dp := range datapoints {
		timestamps, _ := dp.OrderedTimestamps()

		for _, matchedTs := range timestamps {
			total++

			if len(matchedTs) < 2 {
				serumUnmatched++
				continue
			}

			serumTs := matchedTs[0]
			solanaTs := matchedTs[1]

			// sometimes only 1 of asks and bids are matched
			if serumTs.IsZero() && solanaTs.IsZero() {
				continue
			}

			// skip cases where one or other timestamp is zero
			if serumTs.IsZero() {
				solanaUnmatched++
				continue
			}

			if solanaTs.IsZero() {
				serumUnmatched++
				continue
			}

			if serumTs.Before(solanaTs) {
				serumFaster++
				totalSerumLead += int(solanaTs.Sub(serumTs).Milliseconds())
				continue
			}

			if solanaTs.Before(serumTs) {
				solanaFaster++
				totalSolanaLead += int(serumTs.Sub(solanaTs).Milliseconds())
				continue
			}
		}
	}

	averageSerumLead := 0
	if serumFaster > 0 {
		averageSerumLead = totalSerumLead / serumFaster
	}
	averageSolanaLead := 0
	if solanaFaster > 0 {
		averageSolanaLead = totalSolanaLead / solanaFaster
	}

	fmt.Println("Run time: ", runtime)
	fmt.Println("Endpoints:")
	fmt.Println("    ", serumEndpoint, " [serum]")
	fmt.Println("    ", solanaEndpoint, " [solana]")
	fmt.Println()

	fmt.Println("Total updates: ", total)
	fmt.Println()

	fmt.Println("Faster counts: ")
	fmt.Println(fmt.Sprintf("    %-6d  %v", serumFaster, serumEndpoint))
	fmt.Println(fmt.Sprintf("    %-6d  %v", solanaFaster, solanaEndpoint))

	fmt.Println("Average difference( ms): ")
	fmt.Println(fmt.Sprintf("    %-6s  %v", fmt.Sprintf("%vms", averageSerumLead), serumEndpoint))
	fmt.Println(fmt.Sprintf("    %-6s  %v", fmt.Sprintf("%vms", averageSolanaLead), solanaEndpoint))

	fmt.Println("Unmatched updates: ")
	fmt.Println("(updates from each stream without a corresponding result on the other)")
	fmt.Println(fmt.Sprintf("    %-6d  %v", serumUnmatched, serumEndpoint))
	fmt.Println(fmt.Sprintf("    %-6d  %v", solanaUnmatched, solanaEndpoint))
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
