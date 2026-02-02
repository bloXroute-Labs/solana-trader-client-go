package main

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bloXroute-Labs/solana-trader-client-go/connections"
	"github.com/bloXroute-Labs/solana-trader-client-go/examples/config"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/gagliardetto/solana-go"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/programs/system"
	log "github.com/sirupsen/logrus"
)

func main() {
	utils.InitLogger()

	environment := "testnet" // "mainnet", "testnet", "local"
	region := "ny"           // "ny", "uk", "fr"

	client := setupHTTPClient(config.Env(environment), config.HTTPUrls[config.Region(region)])

	submitNcountWithContentType(client, 1, "application/json") // warm up

	iterationCount := 5_000
	//iterationCount = 1

	roundTripStatsJson, bodyUnmarshalStatsJson := submitNcountWithContentType(client, iterationCount, "application/json")
	roundTripStatsText, bodyUnmarshalStatsText := submitNcountWithContentType(client, iterationCount, "text/plain")

	formatStr := "%53s | roundTrip(ns) avg=%11s p50=%11s p90=%11s p99=%11s errors=%3d | bodyUnmarshal(ns) avg=%11s p50=%11s p90=%11s p99=%11s errors=%3d"

	log.Infof(
		formatStr,
		fmt.Sprintf("Post Submit Performance for JSON over %d iterations", iterationCount),
		formatIntWithCommas(roundTripStatsJson.avg),
		formatIntWithCommas(roundTripStatsJson.p50),
		formatIntWithCommas(roundTripStatsJson.p90),
		formatIntWithCommas(roundTripStatsJson.p99),
		roundTripStatsJson.errorsCount,
		formatIntWithCommas(bodyUnmarshalStatsJson.avg),
		formatIntWithCommas(bodyUnmarshalStatsJson.p50),
		formatIntWithCommas(bodyUnmarshalStatsJson.p90),
		formatIntWithCommas(bodyUnmarshalStatsJson.p99),
		bodyUnmarshalStatsJson.errorsCount,
	)

	log.Infof(
		formatStr,
		fmt.Sprintf("Post Submit Performance for TEXT over %d iterations", iterationCount),
		formatIntWithCommas(roundTripStatsText.avg),
		formatIntWithCommas(roundTripStatsText.p50),
		formatIntWithCommas(roundTripStatsText.p90),
		formatIntWithCommas(roundTripStatsText.p99),
		roundTripStatsText.errorsCount,
		formatIntWithCommas(bodyUnmarshalStatsText.avg),
		formatIntWithCommas(bodyUnmarshalStatsText.p50),
		formatIntWithCommas(bodyUnmarshalStatsText.p90),
		formatIntWithCommas(bodyUnmarshalStatsText.p99),
		bodyUnmarshalStatsText.errorsCount,
	)

	log.Infof(
		formatStr,
		"Post Submit Performance diff JSON-TEXT",
		formatIntWithCommas(roundTripStatsJson.avg-roundTripStatsText.avg),
		formatIntWithCommas(roundTripStatsJson.p50-roundTripStatsText.p50),
		formatIntWithCommas(roundTripStatsJson.p90-roundTripStatsText.p90),
		formatIntWithCommas(roundTripStatsJson.p99-roundTripStatsText.p99),
		roundTripStatsJson.errorsCount-roundTripStatsText.errorsCount,
		formatIntWithCommas(bodyUnmarshalStatsJson.avg-bodyUnmarshalStatsText.avg),
		formatIntWithCommas(bodyUnmarshalStatsJson.p50-bodyUnmarshalStatsText.p50),
		formatIntWithCommas(bodyUnmarshalStatsJson.p90-bodyUnmarshalStatsText.p90),
		formatIntWithCommas(bodyUnmarshalStatsJson.p99-bodyUnmarshalStatsText.p99),
		bodyUnmarshalStatsJson.errorsCount-bodyUnmarshalStatsText.errorsCount,
	)

}

type PostSubmitPerf struct {
	iteration               int
	roundTripDurationNs     int64
	bodyUnmarshalDurationNs int64
}

func submitNcountWithContentType(h provider.HTTPClientTraderAPI, iterationCount int, contentType string) (roundTripStats, bodyUnmarshalStats durationSummary) {
	connections.GlobalPostSubmitContentType = contentType

	postSubmitPerfs := make([]PostSubmitPerf, iterationCount)

	for i := 0; i < iterationCount; i++ {
		startTime := time.Now()
		callPostSubmit(h)
		endTime := time.Now()
		postSubmitPerfs[i] = PostSubmitPerf{
			iteration:               i,
			roundTripDurationNs:     endTime.Sub(startTime).Nanoseconds(),
			bodyUnmarshalDurationNs: connections.GlobalBodyUnmarshalDurationNs,
		}
	}

	roundTripStats = summarizeDurations(extractRoundTripDurations(postSubmitPerfs))
	bodyUnmarshalStats = summarizeDurations(extractBodyUnmarshalDurations(postSubmitPerfs))

	return
}

func setupHTTPClient(env config.Env, endpoint string) provider.HTTPClientTraderAPI {
	var h provider.HTTPClientTraderAPI
	var err error

	switch env {
	case config.EnvLocal:
		h = provider.NewHTTPLocal()
	case config.EnvTestnet:
		h = provider.NewHTTPTestnet()
	case config.EnvMainnet:
		h, err = provider.NewHTTPClientFullService(endpoint)
	}
	if err != nil {
		log.Fatalf("error dialing HTTP client: %v", err)
	}

	return h
}

func callPostSubmit(h provider.HTTPClientTraderAPI) bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	response, err := h.GetRecentBlockHash(ctx)
	if err != nil {
		log.Errorf("error with GetRecentBlockHash request: %v", err)
		return true
	}

	bh := solana.MustHashFromBase58(response.BlockHash)

	wlt := solana.NewWallet()
	privateKey, err := transaction.LoadPrivateKeyFromEnv()
	priceLimitIx, err := computebudget.NewSetComputeUnitPriceInstruction(uint64(200000000)).ValidateAndBuild()
	if err != nil {
		return false
	}

	instructions := []solana.Instruction{priceLimitIx}
	for i := 0; i < 40; i++ {
		instructions = append(instructions, system.NewTransferInstruction(100000, privateKey.PublicKey(), solana.MustPublicKeyFromBase58("FZwLKcQupnTy2CbaVMGGsutxDtjv9CqYVDJxiNZSj5Xi")).Build())
	}
	tx1, err := solana.NewTransaction(instructions, bh, solana.TransactionPayer(privateKey.PublicKey()))
	if err != nil {
		return false
	}

	tx1.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(privateKey.PublicKey()) {
			return &privateKey
		}
		return &wlt.PrivateKey
	})

	_, err = h.SignAndSubmit(ctx, &pb.TransactionMessage{
		Content: tx1.MustToBase64()}, false, false, true)
	if err != nil {
		//log.Errorf("failed to sign and submit order (%v)", err)
		return true
	}

	//log.Infof("submitted bundle order to trader api %v", resp)

	return false
}

type durationSummary struct {
	avg         int64
	p50         int64
	p90         int64
	p99         int64
	errorsCount int
}

func extractRoundTripDurations(perfs []PostSubmitPerf) []int64 {
	durations := make([]int64, len(perfs))
	for i, perf := range perfs {
		durations[i] = perf.roundTripDurationNs
	}
	return durations
}

func extractBodyUnmarshalDurations(perfs []PostSubmitPerf) []int64 {
	durations := make([]int64, len(perfs))
	for i, perf := range perfs {
		durations[i] = perf.bodyUnmarshalDurationNs
	}
	return durations
}

func summarizeDurations(durations []int64) durationSummary {
	if len(durations) == 0 {
		return durationSummary{}
	}

	sorted := make([]int64, 0, len(durations))
	for _, d := range durations {
		if d > 0 {
			sorted = append(sorted, d)
		}
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	var sum int64
	for _, value := range sorted {
		sum += value
	}

	return durationSummary{
		avg:         sum / int64(len(sorted)),
		p50:         percentileNearestRank(sorted, 0.50),
		p90:         percentileNearestRank(sorted, 0.90),
		p99:         percentileNearestRank(sorted, 0.99),
		errorsCount: len(durations) - len(sorted),
	}
}

func percentileNearestRank(sorted []int64, p float64) int64 {
	if len(sorted) == 0 {
		return 0
	}
	if p <= 0 {
		return sorted[0]
	}
	if p >= 1 {
		return sorted[len(sorted)-1]
	}

	index := int(math.Ceil(p*float64(len(sorted)))) - 1
	if index < 0 {
		index = 0
	} else if index >= len(sorted) {
		index = len(sorted) - 1
	}
	return sorted[index]
}

func formatIntWithCommas(value int64) string {
	negative := value < 0
	if negative {
		value = -value
	}

	text := strconv.FormatInt(value, 10)
	if len(text) <= 3 {
		if negative {
			return "-" + text
		}
		return text
	}

	firstGroup := len(text) % 3
	if firstGroup == 0 {
		firstGroup = 3
	}

	var b strings.Builder
	b.Grow(len(text) + (len(text)-1)/3)
	if negative {
		b.WriteByte('-')
	}

	b.WriteString(text[:firstGroup])
	for i := firstGroup; i < len(text); i += 3 {
		b.WriteByte(',')
		b.WriteString(text[i : i+3])
	}

	return b.String()
}
