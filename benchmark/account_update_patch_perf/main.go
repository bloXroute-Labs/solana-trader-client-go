package main

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var runDuration int

func createHash(update *pb.GetNewRaydiumPoolsResponse) string {
	data := fmt.Sprintf("%d%s%s%d%s%s%d%s%s%d%s%s%s%s%s%d%s%s%s%s%s%s%s%s%d%s%s%s%s%s%s%s%s%d",
		update.Slot,
		update.Pool.Pool,
		update.Pool.PoolAddress,
		update.Pool.Token1Reserves,
		update.Pool.Token1MintAddress,
		update.Pool.Token1MintSymbol,
		update.Pool.Token2Reserves,
		update.Pool.Token2MintAddress,
		update.Pool.Token2MintSymbol,
		update.Pool.OpenTime,
		update.Pool.PoolType,
		update.Pool.LiquidityPoolKeys.Id,
		update.Pool.LiquidityPoolKeys.BaseMint,
		update.Pool.LiquidityPoolKeys.QuoteMint,
		update.Pool.LiquidityPoolKeys.LpMint,
		update.Pool.LiquidityPoolKeys.Version,
		update.Pool.LiquidityPoolKeys.ProgramID,
		update.Pool.LiquidityPoolKeys.Authority,
		update.Pool.LiquidityPoolKeys.BaseVault,
		update.Pool.LiquidityPoolKeys.QuoteVault,
		update.Pool.LiquidityPoolKeys.LpVault,
		update.Pool.LiquidityPoolKeys.OpenOrders,
		update.Pool.LiquidityPoolKeys.TargetOrders,
		update.Pool.LiquidityPoolKeys.WithdrawQueue,
		update.Pool.LiquidityPoolKeys.MarketVersion,
		update.Pool.LiquidityPoolKeys.MarketProgramID,
		update.Pool.LiquidityPoolKeys.MarketID,
		update.Pool.LiquidityPoolKeys.MarketAuthority,
		update.Pool.LiquidityPoolKeys.MarketBaseVault,
		update.Pool.LiquidityPoolKeys.MarketQuoteVault,
		update.Pool.LiquidityPoolKeys.MarketBids,
		update.Pool.LiquidityPoolKeys.MarketAsks,
		update.Pool.LiquidityPoolKeys.MarketEventQueue,
		update.Pool.LiquidityPoolKeys.TradeFeeRate,
	)
	// return fmt.Sprintf("%x%s", sha256.Sum256([]byte(data)), data)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(data)))
}

func main() {

	flag.IntVar(&runDuration, "duration", 1, "run duration in minutes")
	flag.Parse()

	var ch = make(chan accountUpdate, 1e7)
	var wg = new(sync.WaitGroup)
	var ctx, _ = context.WithTimeout(context.Background(), time.Duration(runDuration)*time.Minute)

	wg.Add(2)
	go accountUpdatesPatchLib(ctx, wg, ch, "127.0.0.1:9001", sourceAccountUpdateStream, false)
	go accountUpdatesPatchLib(ctx, wg, ch, "127.0.0.1:9000", sourceAccountUpdatePatchStream, false)

	wg.Wait()
	close(ch)

	var mp = make(map[string]accountUpdate)
	var timeDiffsUs = make([]int64, 0)

	for upd := range ch {

		hash := createHash(upd.update)

		if v, ok := mp[hash]; !ok {
			mp[hash] = upd
		} else {
			if v.source == upd.source {
				fmt.Println("update from the same source", v.source)
				continue
			}

			diff := upd.rcvTime.Sub(v.rcvTime).Microseconds()
			if upd.source == sourceAccountUpdateStream { // VANILLA WIN
				// sourceAccountUpdatePatchStream is WINNER (-)
				timeDiffsUs = append(timeDiffsUs, -diff)
			} else { // not VANILLA WIN
				// sourceAccountUpdateStream is WINNER (+)
				timeDiffsUs = append(timeDiffsUs, diff)
			}

			delete(mp, hash)
		}
	}

	// Did we have pool events with no match?
	streamCount := 0
	patchStreamCount := 0
	for hash, entry := range mp {
		fmt.Printf("Hash: %s, Entry: %+v\n", hash, entry)
		if entry.source == sourceAccountUpdateStream {
			streamCount++
		} else if entry.source == sourceAccountUpdatePatchStream {
			patchStreamCount++
		}
	}

	fmt.Println("----------------")
	fmt.Println("Benchmark comparing contents of AccountUpdateStream and AccountUpdatePatchStream")
	fmt.Println("- sdk run for GetNewRaydiumPoolsStream")
	fmt.Println("Duration: ", runDuration, "minutes")
	fmt.Println("----------------")
	fmt.Printf("Entries from AccountUpdateStream with no match: %d\n", streamCount)
	fmt.Printf("Entries from AccountUpdatePatchStream with no match: %d\n", patchStreamCount)

	timeDiffsMs := []float64{}
	for _, timestamp := range timeDiffsUs {
		timeDiffsMs = append(timeDiffsMs, float64(timestamp)/1000)
	}

	fmt.Println("----------------")
	fmt.Println("Number of pairs that matched:", len(timeDiffsUs))

	fmt.Println("----------------")
	fmt.Println("values with (+) is AccountUpdatesStream as the winner")
	fmt.Println("values with (-) is AccountUpdatesPatchStream as the winner")
	calculatePercentiles(timeDiffsUs)

	fmt.Println("----------------")
	fmt.Println("Grouping by millisecond")
	fmt.Println(makeHistogram(timeDiffsMs))
	fmt.Println("----------------")
	fmt.Println("")
}

const (
	sourceAccountUpdateStream = iota
	sourceAccountUpdatePatchStream
)

type source int

type accountUpdate struct {
	source    source
	rcvTime   time.Time
	update    *pb.GetNewRaydiumPoolsResponse
}

func accountUpdatesPatchLib(ctx context.Context, wg *sync.WaitGroup, ch chan accountUpdate, addr string, src source, testnet bool) {
	defer wg.Done()

	var client *provider.GRPCClient
	var err error
	if testnet {
		client, err = provider.NewGRPCTestnet()
		if err != nil {
			panic("failed to get grpc client")
		}

	} else {
		opts := provider.DefaultRPCOpts(addr)
		client, err = provider.NewGRPCClientWithOpts(opts)
		if err != nil {
			panic("failed to get grpc client")
		}
	}

	chRaydium := make(chan *pb.GetNewRaydiumPoolsResponse)
	stream, err := client.GetNewRaydiumPoolsStream(ctx, true)
	if err != nil {
		panic(fmt.Sprintf("failed to get new Raydium pools stream: %s", err))
	}
	stream.Into(chRaydium)

	for {
		select {
		case <-ctx.Done():
			return
		case pool := <-chRaydium:
			fmt.Println("New Raydium pool:", pool.Pool.PoolAddress)
			select {
			case ch <- accountUpdate{
				source:    src,
				rcvTime:   time.Now(),
				update:    pool,
			}:
			default:
				panic("accountUpdatesPatchLib: send into recv chan failed: full")
			}
		}
	}
}

const (
	numBins  = 22
	binSize  = 1
	minValue = -10
	maxValue = 10
)

func makeHistogram(data []float64) string {
	if len(data) == 0 {
		return ""
	}
	histogram := make([]int, numBins)

	for _, value := range data {
		binIndex := getBinIndex(value, minValue, maxValue)
		if binIndex >= 0 && binIndex < numBins {
			histogram[binIndex]++
		}
	}

	result := ""
	for i, count := range histogram {
		percentage := float64(count) / float64(len(data)) * 100
		binStart := getBinStart(i)
		binEnd := getBinEnd(i)
		bar := getHistogramBar(int(percentage))

		result += fmt.Sprintf("%3s <-> %3s  %6.2f%%  %4s %s\n", binStart, binEnd, percentage, strconv.Itoa(count), bar)
	}

	return result
}

func getHistogramBar(count int) string {
	barLength := count / 2
	bar := strings.Repeat("█", barLength)
	remainder := count % 2
	if remainder > 0 {
		bar += "▏"
	}
	return bar
}

func getBinIndex(value, minValue, maxValue float64) int {
	if value < minValue {
		return 0
	}
	if value > maxValue {
		return numBins - 1
	}

	return int(math.Ceil(value+maxValue) / binSize)
}

func getBinStart(binIndex int) string {
	if binIndex == 0 {
		return "-∞"
	}
	return fmt.Sprintf("%d", -10+(binIndex-1)*binSize)
}

func getBinEnd(binIndex int) string {
	if binIndex == 0 {
		return fmt.Sprintf("%d", -10)
	}
	if binIndex == numBins-1 {
		return "+∞"
	}

	return fmt.Sprintf("%d", -10+binIndex*binSize)
}

var percentiles = []int{10, 25, 50, 75, 90}

func calculatePercentiles(data []int64) {
	if len(data) == 0 {
		return
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i] < data[j]
	})

	for _, percentile := range percentiles {
		fmt.Printf("P%d: %dμs\n", percentile, calculatePercentile(data, float64(percentile)))
	}
}

func calculatePercentile(data []int64, percentile float64) int64 {
	index := int((percentile / 100) * float64(len(data)-1))
	return data[index]
}
