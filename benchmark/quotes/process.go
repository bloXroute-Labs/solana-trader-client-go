package main

import (
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/actor"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/stream"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
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
		fmt.Printf("Slot range: %v => %v\n", br.firstJupiter().ContextSlot, br.lastJupiter().ContextSlot)
		fmt.Printf("Initial price: %v\n", br.firstJupiter().Price(br.mint))
		fmt.Printf("Final price: %v\n", br.lastJupiter().Price(br.mint))
		fmt.Printf("Distinct prices: %v\n", 1)
	}
	fmt.Println()

	fmt.Println("Trader WS: ", len(br.tradeWSUpdates), " samples")
	if len(br.tradeWSUpdates) > 0 {
		fmt.Printf("Slot range: %v => %v\n", br.firstWS().Slot, br.lastWS().Slot)
		fmt.Printf("Initial buy price: %v\n", br.firstWS().Price.Buy)
		fmt.Printf("Initial sell price: %v\n", br.firstWS().Price.Sell)
		fmt.Printf("Final buy price: %v\n", br.lastWS().Price.Buy)
		fmt.Printf("Final sell price: %v\n", br.lastWS().Price.Sell)
		fmt.Printf("Distinct prices: %v\n", 1)
	}
	fmt.Println()

	fmt.Println("Trader HTTP: ", len(br.tradeWSUpdates), " samples")
	if len(br.tradeHTTPUpdates) > 0 {
		fmt.Printf("Initial buy price: %v\n", br.firstWS().Price.Buy)
		fmt.Printf("Initial sell price: %v\n", br.firstWS().Price.Sell)
		fmt.Printf("Final buy price: %v\n", br.lastWS().Price.Buy)
		fmt.Printf("Final sell price: %v\n", br.lastWS().Price.Sell)
		fmt.Printf("Distinct prices: %v\n", 1)
	}
	fmt.Println()
}

func (br benchmarkResult) PrintSimple() {
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

func (br benchmarkResult) firstJupiter() stream.JupiterPriceResponse {
	return br.jupiterUpdates[0].Data.Data
}

func (br benchmarkResult) lastJupiter() stream.JupiterPriceResponse {
	return br.jupiterUpdates[len(br.jupiterUpdates)-1].Data.Data
}

func (br benchmarkResult) firstWS() *pb.GetPricesStreamResponse {
	return br.tradeWSUpdates[0].Data
}

func (br benchmarkResult) lastWS() *pb.GetPricesStreamResponse {
	return br.tradeWSUpdates[len(br.tradeWSUpdates)-1].Data
}

func (br benchmarkResult) firstHTTP() *pb.GetPriceResponse {
	return br.tradeHTTPUpdates[0].Data.Data
}

func (br benchmarkResult) lastHTTP() *pb.GetPriceResponse {
	return br.tradeHTTPUpdates[len(br.tradeHTTPUpdates)-1].Data.Data
}
