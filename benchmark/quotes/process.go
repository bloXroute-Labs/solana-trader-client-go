package main

import (
	"encoding/csv"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/actor"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/output"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/stream"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"os"
	"sort"
	"strconv"
	"time"
)

type benchmarkResult struct {
	mint                string
	swaps               []actor.SwapEvent
	jupiterRawUpdates   []stream.RawUpdate[stream.DurationUpdate[*stream.JupiterPriceResponse]]
	tradeWSRawUpdates   []stream.RawUpdate[*pb.GetPricesStreamResponse]
	tradeHTTPRawUpdates []stream.RawUpdate[stream.DurationUpdate[*pb.GetPriceResponse]]

	jupiterProcessedUpdates   map[int][]stream.ProcessedUpdate[stream.QuoteResult]
	tradeWSProcessedUpdates   map[int][]stream.ProcessedUpdate[stream.QuoteResult]
	tradeHTTPProcessedUpdates map[int][]stream.ProcessedUpdate[stream.QuoteResult]
}

func (br benchmarkResult) PrintSummary() {
	fmt.Println("Trader API vs. Jupiter API Benchmark")
	fmt.Println()

	fmt.Println("Swaps placed: ", len(br.swaps))
	for _, swap := range br.swaps {
		fmt.Printf("%v: %v\n", swap.Timestamp.Format(time.StampMilli), swap.Signature)
	}
	fmt.Println()

	fmt.Println("Jupiter: ", len(br.jupiterRawUpdates), " samples")
	if len(br.jupiterRawUpdates) > 0 {
		startTime := br.firstJupiter().Start
		endTime := br.lastJupiter().Start
		fmt.Printf("Start time: %v\n", startTime.Format(time.StampMilli))
		fmt.Printf("End time: %v\n", endTime.Format(time.StampMilli))
		fmt.Printf("Slot range: %v => %v\n", br.firstJupiter().Data.ContextSlot, br.lastJupiter().Data.ContextSlot)
		fmt.Printf("Price change: %v => %v \n", br.firstJupiter().Data.Price(br.mint), br.lastJupiter().Data.Price(br.mint))
		fmt.Printf("Distinct prices: %v\n", br.distinctJupiter())
	}
	fmt.Println()

	fmt.Println("Trader WS: ", len(br.tradeWSRawUpdates), " samples")
	if len(br.tradeWSRawUpdates) > 0 {
		startTime := br.firstWS().Timestamp
		endTime := br.lastWS().Timestamp
		fmt.Printf("Start time: %v\n", startTime.Format(time.StampMilli))
		fmt.Printf("End time: %v\n", endTime.Format(time.StampMilli))
		fmt.Printf("Slot range: %v => %v\n", br.firstWS().Data.Slot, br.lastWS().Data.Slot)
		fmt.Printf("Buy change: %v => %v\n", br.firstWS().Data.Price.Buy, br.lastWS().Data.Price.Buy)
		fmt.Printf("Sell change: %v => %v\n", br.firstWS().Data.Price.Sell, br.lastWS().Data.Price.Sell)
		fmt.Printf("Distinct buy prices: %v\n", br.distinctWSBuy())
		fmt.Printf("Distinct sell prices: %v\n", br.distinctWSSell())
	}
	fmt.Println()

	fmt.Println("Trader HTTP: ", len(br.tradeHTTPRawUpdates), " samples")
	if len(br.tradeHTTPRawUpdates) > 0 {
		startTime := br.firstHTTP().Start
		endTime := br.lastHTTP().Start
		fmt.Printf("Start time: %v\n", startTime.Format(time.StampMilli))
		fmt.Printf("End time: %v\n", endTime.Format(time.StampMilli))
		fmt.Printf("Buy change: %v => %v\n", br.firstHTTP().Data.TokenPrices[0].Buy, br.lastHTTP().Data.TokenPrices[0].Buy)
		fmt.Printf("Sell change: %v => %v\n", br.firstHTTP().Data.TokenPrices[0].Sell, br.lastHTTP().Data.TokenPrices[0].Sell)
		fmt.Printf("Distinct buy prices: %v\n", br.distinctHTTPBuy())
		fmt.Printf("Distinct sell prices: %v\n", br.distinctHTTPSell())
	}
	fmt.Println()
}

func (br benchmarkResult) PrintRaw() {
	fmt.Println("jupiter API")
	printProcessedUpdate(br.jupiterProcessedUpdates)
	fmt.Println()

	fmt.Println("traderWS")
	printProcessedUpdate(br.tradeWSProcessedUpdates)
	fmt.Println()

	fmt.Println("traderHTTP")
	printProcessedUpdate(br.tradeHTTPProcessedUpdates)
}

func (br benchmarkResult) WriteCSV(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	w := csv.NewWriter(f)
	defer w.Flush()

	err = w.Write([]string{
		"timestamp",
		"source",
		"slot",
		"processingTime",
		"buy",
		"sell",
	})
	if err != nil {
		return err
	}

	allUpdates := br.allUpdates()

	for _, update := range allUpdates {
		err = w.Write([]string{
			update.Timestamp.Format(time.RFC3339Nano),
			update.Data.Source,
			strconv.Itoa(update.Slot),
			update.Data.Elapsed.String(),
			strconv.FormatFloat(update.Data.BuyPrice, 'f', -1, 64),
			strconv.FormatFloat(update.Data.SellPrice, 'f', -1, 64),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (br benchmarkResult) firstJupiter() stream.DurationUpdate[*stream.JupiterPriceResponse] {
	return br.jupiterRawUpdates[0].Data
}

func (br benchmarkResult) lastJupiter() stream.DurationUpdate[*stream.JupiterPriceResponse] {
	return br.jupiterRawUpdates[len(br.jupiterRawUpdates)-1].Data
}

func (br benchmarkResult) distinctJupiter() int {
	return distinctPrices(br.jupiterRawUpdates, func(v stream.RawUpdate[stream.DurationUpdate[*stream.JupiterPriceResponse]]) float64 {
		return v.Data.Data.Price(br.mint)
	})
}

func (br benchmarkResult) firstWS() stream.RawUpdate[*pb.GetPricesStreamResponse] {
	return br.tradeWSRawUpdates[0]
}

func (br benchmarkResult) lastWS() stream.RawUpdate[*pb.GetPricesStreamResponse] {
	return br.tradeWSRawUpdates[len(br.tradeWSRawUpdates)-1]
}

func (br benchmarkResult) distinctWSBuy() int {
	return distinctPrices(br.tradeWSRawUpdates, func(v stream.RawUpdate[*pb.GetPricesStreamResponse]) float64 {
		return v.Data.Price.Buy
	})
}

func (br benchmarkResult) distinctWSSell() int {
	return distinctPrices(br.tradeWSRawUpdates, func(v stream.RawUpdate[*pb.GetPricesStreamResponse]) float64 {
		return v.Data.Price.Sell
	})
}

func (br benchmarkResult) firstHTTP() stream.DurationUpdate[*pb.GetPriceResponse] {
	return br.tradeHTTPRawUpdates[0].Data
}

func (br benchmarkResult) lastHTTP() stream.DurationUpdate[*pb.GetPriceResponse] {
	return br.tradeHTTPRawUpdates[len(br.tradeHTTPRawUpdates)-1].Data
}

func (br benchmarkResult) distinctHTTPBuy() int {
	return distinctPrices(br.tradeHTTPRawUpdates, func(v stream.RawUpdate[stream.DurationUpdate[*pb.GetPriceResponse]]) float64 {
		return v.Data.Data.TokenPrices[0].Buy
	})
}

func (br benchmarkResult) distinctHTTPSell() int {
	return distinctPrices(br.tradeHTTPRawUpdates, func(v stream.RawUpdate[stream.DurationUpdate[*pb.GetPriceResponse]]) float64 {
		return v.Data.Data.TokenPrices[0].Sell
	})
}

func (br benchmarkResult) allUpdates() []stream.ProcessedUpdate[stream.QuoteResult] {
	allEvents := make([]stream.ProcessedUpdate[stream.QuoteResult], 0)

	for _, updates := range br.jupiterProcessedUpdates {
		for _, update := range updates {
			allEvents = append(allEvents, update)
		}
	}
	for _, updates := range br.tradeWSProcessedUpdates {
		for _, update := range updates {
			allEvents = append(allEvents, update)
		}
	}
	for _, updates := range br.tradeHTTPProcessedUpdates {
		for _, update := range updates {
			allEvents = append(allEvents, update)
		}
	}

	sort.Slice(allEvents, func(i, j int) bool {
		return allEvents[i].Timestamp.Before(allEvents[j].Timestamp)
	})

	return allEvents
}

func distinctPrices[T any](s []T, getter func(T) float64) int {
	prices := make(map[float64]bool)

	count := 0
	for _, v := range s {
		price := getter(v)
		_, ok := prices[price]
		if ok {
			continue
		}

		prices[price] = true
		count++
	}

	return count
}

func printProcessedUpdate(pu map[int][]stream.ProcessedUpdate[stream.QuoteResult]) {
	slots := output.SortRange(pu)
	for _, slot := range slots {
		for _, update := range pu[slot] {
			printLine(update)
		}
	}
}

func printLine(update stream.ProcessedUpdate[stream.QuoteResult]) {
	fmt.Printf("[%v] %v [%vms]: B: %v | S: %v\n", update.Slot, update.Timestamp, update.Data.Elapsed.Milliseconds(), update.Data.BuyPrice, update.Data.SellPrice)
}
