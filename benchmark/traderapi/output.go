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
	Slot             int
	TraderTimestamps []time.Time
	SolanaAsk        time.Time
	SolanaBid        time.Time
}

// OrderedTimestamps returns timestamps pairs of [serumTS, solanaTS] within the slot. For example, the first TS from Solana Trader API is matched with the sooner of SolanaAsk and SolanaBid, the second with the later, and the rest with zero values.
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

	for i, ts := range d.TraderTimestamps {
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

		line := fmt.Sprintf("slot %v (%v): solana trader %v, solana %v, diff %v", d.Slot, side, formatTS(timestamps[0]), formatTS(timestamps[1]), formatDiff(timestamps[0], timestamps[1]))
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
	if len(slotRange) == 0 {
		return "-"
	}
	sr := SortRange(slotRange)
	return fmt.Sprintf("%v-%v", sr[0], sr[len(sr)-1])
}

// SlotRange enumerate the superset range of slots used in Trader API and Solana updates
func SlotRange(traderResults map[int][]arrival.ProcessedUpdate[arrival.TraderAPIUpdate], solanaResults map[int][]arrival.ProcessedUpdate[arrival.SolanaUpdate]) []int {
	traderSlots := SortRange(traderResults)
	solanaSlots := SortRange(solanaResults)

	slots := make([]int, 0, len(traderResults))
	solanaIndex := 0
	for i := 0; i < len(traderSlots); i++ {
		traderCandidate := traderSlots[i]

		for j := solanaIndex; j < len(solanaSlots); j++ {
			solanaCandidate := solanaSlots[j]
			if solanaCandidate < traderCandidate {
				slots = append(slots, solanaCandidate)
				solanaIndex++
			} else if solanaCandidate == traderCandidate {
				solanaIndex++
			} else {
				break
			}
		}

		slots = append(slots, traderCandidate)
	}

	for j := solanaIndex; j < len(solanaSlots); j++ {
		slots = append(slots, solanaSlots[j])
	}

	return slots
}

// Merge combines Trader API and Solana updates over the specified slots, indicating the difference in slot times and any updates that were not included in the other.
func Merge(slots []int, traderResults map[int][]arrival.ProcessedUpdate[arrival.TraderAPIUpdate], solanaResults map[int][]arrival.ProcessedUpdate[arrival.SolanaUpdate]) ([]Datapoint, map[int][]arrival.ProcessedUpdate[arrival.TraderAPIUpdate], map[int][]arrival.ProcessedUpdate[arrival.SolanaUpdate], error) {
	datapoints := make([]Datapoint, 0)
	leftoverTrader := make(map[int][]arrival.ProcessedUpdate[arrival.TraderAPIUpdate])
	leftoverSolana := make(map[int][]arrival.ProcessedUpdate[arrival.SolanaUpdate])

	for _, slot := range slots {
		traderData, traderOK := traderResults[slot]
		solanaData, solanaOK := solanaResults[slot]

		if !traderOK && !solanaOK {
			return nil, nil, nil, fmt.Errorf("(slot %v) improper slot set provided", slot)
		}

		if !traderOK {
			leftoverSolana[slot] = solanaData
		}

		if !solanaOK {
			leftoverTrader[slot] = traderData
		}

		dp := Datapoint{
			Slot:             slot,
			TraderTimestamps: make([]time.Time, 0, len(traderData)),
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

		for _, su := range traderData {
			dp.TraderTimestamps = append(dp.TraderTimestamps, su.Timestamp)
		}

		datapoints = append(datapoints, dp)
	}

	return datapoints, leftoverTrader, leftoverSolana, nil
}

func PrintSummary(runtime time.Duration, traderEndpoint string, solanaEndpoint string, datapoints []Datapoint) {
	traderFaster := 0
	solanaFaster := 0
	totalTraderLead := 0
	totalSolanaLead := 0
	traderUnmatched := 0
	solanaUnmatched := 0
	total := 0

	for _, dp := range datapoints {
		timestamps, _ := dp.OrderedTimestamps()

		for _, matchedTs := range timestamps {
			total++

			if len(matchedTs) < 2 {
				traderUnmatched++
				continue
			}

			traderTs := matchedTs[0]
			solanaTs := matchedTs[1]

			// sometimes only 1 of asks and bids are matched
			if traderTs.IsZero() && solanaTs.IsZero() {
				continue
			}

			// skip cases where one or other timestamp is zero
			if traderTs.IsZero() {
				solanaUnmatched++
				continue
			}

			if solanaTs.IsZero() {
				traderUnmatched++
				continue
			}

			if traderTs.Before(solanaTs) {
				traderFaster++
				totalTraderLead += int(solanaTs.Sub(traderTs).Milliseconds())
				continue
			}

			if solanaTs.Before(traderTs) {
				solanaFaster++
				totalSolanaLead += int(traderTs.Sub(solanaTs).Milliseconds())
				continue
			}
		}
	}

	averageTraderLead := 0
	if traderFaster > 0 {
		averageTraderLead = totalTraderLead / traderFaster
	}
	averageSolanaLead := 0
	if solanaFaster > 0 {
		averageSolanaLead = totalSolanaLead / solanaFaster
	}

	fmt.Println("Run time: ", runtime)
	fmt.Println("Endpoints:")
	fmt.Println("    ", traderEndpoint, " [traderapi]")
	fmt.Println("    ", solanaEndpoint, " [solana]")
	fmt.Println()

	fmt.Println("Total updates: ", total)
	fmt.Println()

	fmt.Println("Faster counts: ")
	fmt.Println(fmt.Sprintf("    %-6d  %v", traderFaster, traderEndpoint))
	fmt.Println(fmt.Sprintf("    %-6d  %v", solanaFaster, solanaEndpoint))

	fmt.Println("Average difference( ms): ")
	fmt.Println(fmt.Sprintf("    %-6s  %v", fmt.Sprintf("%vms", averageTraderLead), traderEndpoint))
	fmt.Println(fmt.Sprintf("    %-6s  %v", fmt.Sprintf("%vms", averageSolanaLead), solanaEndpoint))

	fmt.Println("Unmatched updates: ")
	fmt.Println("(updates from each stream without a corresponding result on the other)")
	fmt.Println(fmt.Sprintf("    %-6d  %v", traderUnmatched, traderEndpoint))
	fmt.Println(fmt.Sprintf("    %-6d  %v", solanaUnmatched, solanaEndpoint))
}

func formatTS(ts time.Time) string {
	if ts.IsZero() {
		return "n/a"
	} else {
		return ts.Format(tsFormat)
	}
}

func formatDiff(traderTS time.Time, solanaTS time.Time) string {
	if traderTS.IsZero() || solanaTS.IsZero() {
		return "n/a"
	} else {
		return traderTS.Sub(solanaTS).String()
	}
}
