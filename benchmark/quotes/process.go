package main

import (
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/actor"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/stream"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"math"
	"time"
)

type benchmarkResult struct {
	mint             string
	swaps            []actor.SwapEvent
	jupiterUpdates   []stream.RawUpdate[stream.DurationUpdate[stream.JupiterPriceResponse]]
	tradeWSUpdates   []stream.RawUpdate[*pb.GetPricesStreamResponse]
	tradeHTTPUpdates []stream.RawUpdate[stream.DurationUpdate[*pb.GetPriceResponse]]
}

func (br benchmarkResult) PrintSummary() {
	fmt.Println("Trader API vs. Jupiter API Benchmark")
	fmt.Println()

	fmt.Println("Swaps placed: ", len(br.swaps))
	for _, swap := range br.swaps {
		fmt.Printf("%v: %v\n", swap.Timestamp.Format(time.StampMilli), swap.Signature)
	}
	fmt.Println()

	fmt.Println("Jupiter: ", len(br.jupiterUpdates), " samples")
	if len(br.jupiterUpdates) > 0 {
		startTime := br.firstJupiter().Start
		endTime := br.lastJupiter().Start
		fmt.Printf("Start time: %v\n", startTime.Format(time.StampMilli))
		fmt.Printf("End time: %v\n", endTime.Format(time.StampMilli))
		fmt.Printf("Expected slots: %v\n", expectedSlots(startTime, endTime))
		fmt.Printf("Slot range: %v => %v\n", br.firstJupiter().Data.ContextSlot, br.lastJupiter().Data.ContextSlot)
		fmt.Printf("Initial price: %v\n", br.firstJupiter().Data.Price(br.mint))
		fmt.Printf("Final price: %v\n", br.lastJupiter().Data.Price(br.mint))
		fmt.Printf("Distinct prices: %v\n", 1)
	}
	fmt.Println()

	fmt.Println("Trader WS: ", len(br.tradeWSUpdates), " samples")
	if len(br.tradeWSUpdates) > 0 {
		startTime := br.firstWS().Timestamp
		endTime := br.lastWS().Timestamp
		fmt.Printf("Start time: %v\n", startTime.Format(time.StampMilli))
		fmt.Printf("End time: %v\n", endTime.Format(time.StampMilli))
		fmt.Printf("Expected slots: %v\n", expectedSlots(startTime, endTime))
		fmt.Printf("Slot range: %v => %v\n", br.firstWS().Data.Slot, br.lastWS().Data.Slot)
		fmt.Printf("Initial buy price: %v\n", br.firstWS().Data.Price.Buy)
		fmt.Printf("Initial sell price: %v\n", br.firstWS().Data.Price.Sell)
		fmt.Printf("Final buy price: %v\n", br.lastWS().Data.Price.Buy)
		fmt.Printf("Final sell price: %v\n", br.lastWS().Data.Price.Sell)
		fmt.Printf("Distinct buy prices: %v\n", br.distinctWSBuy())
		fmt.Printf("Distinct sell prices: %v\n", br.distinctWSSell())
	}
	fmt.Println()

	fmt.Println("Trader HTTP: ", len(br.tradeHTTPUpdates), " samples")
	if len(br.tradeHTTPUpdates) > 0 {
		startTime := br.firstHTTP().Start
		endTime := br.lastHTTP().Start
		fmt.Printf("Start time: %v\n", startTime.Format(time.StampMilli))
		fmt.Printf("End time: %v\n", endTime.Format(time.StampMilli))
		fmt.Printf("Expected slots: %v\n", expectedSlots(startTime, endTime))
		fmt.Printf("Initial buy price: %v\n", br.firstHTTP().Data.TokenPrices[0].Buy)
		fmt.Printf("Initial sell price: %v\n", br.firstHTTP().Data.TokenPrices[0].Sell)
		fmt.Printf("Final buy price: %v\n", br.lastHTTP().Data.TokenPrices[0].Buy)
		fmt.Printf("Final sell price: %v\n", br.lastHTTP().Data.TokenPrices[0].Sell)
		fmt.Printf("Distinct buy prices: %v\n", br.distinctHTTPBuy())
		fmt.Printf("Distinct sell prices: %v\n", br.distinctHTTPSell())
	}
	fmt.Println()
}

func (br benchmarkResult) PrintRaw() {
	fmt.Println("jupiter API")
	for _, update := range br.jupiterUpdates {
		fmt.Printf("[%v] %v => %v: %v\n", update.Data.Data.ContextSlot, update.Data.Start, update.Timestamp, update.Data.Data.PriceInfo[br.mint].Price)
	}

	fmt.Println("traderWS")
	for _, update := range br.tradeWSUpdates {
		fmt.Printf("[%v] %v: B %v | S %v\n", update.Data.Slot, update.Timestamp, update.Data.Price.Buy, update.Data.Price.Sell)
	}

	fmt.Println("traderHTTP")
	for _, update := range br.tradeHTTPUpdates {
		fmt.Printf("%v => %v: B %v | S %v\n", update.Data.Start, update.Timestamp, update.Data.Data.TokenPrices[0].Buy, update.Data.Data.TokenPrices[0].Sell)
	}
}

func (br benchmarkResult) firstJupiter() stream.DurationUpdate[stream.JupiterPriceResponse] {
	return br.jupiterUpdates[0].Data
}

func (br benchmarkResult) lastJupiter() stream.DurationUpdate[stream.JupiterPriceResponse] {
	return br.jupiterUpdates[len(br.jupiterUpdates)-1].Data
}

func (br benchmarkResult) distinctJupiter() int {
	return distinctPrices(br.jupiterUpdates, func(v stream.RawUpdate[stream.DurationUpdate[stream.JupiterPriceResponse]]) float64 {
		return v.Data.Data.Price(br.mint)
	})
}

func (br benchmarkResult) firstWS() stream.RawUpdate[*pb.GetPricesStreamResponse] {
	return br.tradeWSUpdates[0]
}

func (br benchmarkResult) lastWS() stream.RawUpdate[*pb.GetPricesStreamResponse] {
	return br.tradeWSUpdates[len(br.tradeWSUpdates)-1]
}

func (br benchmarkResult) distinctWSBuy() int {
	return distinctPrices(br.tradeWSUpdates, func(v stream.RawUpdate[*pb.GetPricesStreamResponse]) float64 {
		return v.Data.Price.Buy
	})
}

func (br benchmarkResult) distinctWSSell() int {
	return distinctPrices(br.tradeWSUpdates, func(v stream.RawUpdate[*pb.GetPricesStreamResponse]) float64 {
		return v.Data.Price.Sell
	})
}

func (br benchmarkResult) firstHTTP() stream.DurationUpdate[*pb.GetPriceResponse] {
	return br.tradeHTTPUpdates[0].Data
}

func (br benchmarkResult) lastHTTP() stream.DurationUpdate[*pb.GetPriceResponse] {
	return br.tradeHTTPUpdates[len(br.tradeHTTPUpdates)-1].Data
}

func (br benchmarkResult) distinctHTTPBuy() int {
	return distinctPrices(br.tradeHTTPUpdates, func(v stream.RawUpdate[stream.DurationUpdate[*pb.GetPriceResponse]]) float64 {
		return v.Data.Data.TokenPrices[0].Buy
	})
}

func (br benchmarkResult) distinctHTTPSell() int {
	return distinctPrices(br.tradeHTTPUpdates, func(v stream.RawUpdate[stream.DurationUpdate[*pb.GetPriceResponse]]) float64 {
		return v.Data.Data.TokenPrices[0].Sell
	})
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

func expectedSlots(start, end time.Time) float64 {
	return math.Round(end.Sub(start).Seconds() * 5)
}
