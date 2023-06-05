package stream

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/logger"
	"go.uber.org/zap"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"
)

const (
	quoteAPIEndpoint = "https://quote-api.jup.ag/v4/quote"
	priceAPIEndpoint = "https://price.jup.ag/v4/price"
	defaultAmount    = 1
	defaultSlippage  = 5
	defaultDecimals  = 8
	defaultInterval  = time.Second
	usdcMint         = "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
	usdcDecimals     = 6
)

type JupiterAPIStream struct {
	client          *http.Client
	mint            string
	amount          float64 // amount of USDC to provide for the swap (in $)
	adjustedAmount  int     // amount of USDC to provide for the swap (in units)
	decimals        int     // e.g. Jupiter requires 1000000 to indicate 1 USDC. This one's for the other mint.
	priceAdjustment float64 // number to multiply final price by to account for decimals.
	slippageBps     int64
	ticker          *time.Ticker
	interval        time.Duration
}

func NewJupiterAPI(opts ...JupiterOpt) (Source[DurationUpdate[*JupiterPriceResponse], QuoteResult], error) {
	j := &JupiterAPIStream{
		client:      &http.Client{},
		amount:      defaultAmount,
		decimals:    defaultDecimals,
		slippageBps: defaultSlippage,
		interval:    defaultInterval,
	}

	for _, o := range opts {
		o(j)
	}

	if j.mint == "" {
		return nil, errors.New("mint token is mandatory")
	}

	j.adjustedAmount = int(j.amount * math.Pow(10, float64(usdcDecimals)))
	j.priceAdjustment = math.Pow(10, float64(j.decimals-usdcDecimals))
	return j, nil
}

func (j *JupiterAPIStream) log() *zap.SugaredLogger {
	return logger.Log().With("source", "jupiterApi")
}

func (j *JupiterAPIStream) Name() string {
	return "jupiter"
}

func (j *JupiterAPIStream) Run(parent context.Context) ([]RawUpdate[DurationUpdate[*JupiterPriceResponse]], error) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	ticker := j.ticker
	if ticker == nil {
		ticker = time.NewTicker(j.interval)
	}

	return collectOrderedUpdates(ctx, ticker, func() (*JupiterPriceResponse, error) {
		res, err := j.FetchQuote(ctx)
		if err != nil {
			return nil, err
		}

		return &res, err
	}, nil, func(err error) {
		j.log().Errorw("could not fetch price", "err", err)
	})
}

func (j *JupiterAPIStream) Process(updates []RawUpdate[DurationUpdate[*JupiterPriceResponse]], removeDuplicates bool) (results map[int][]ProcessedUpdate[QuoteResult], duplicates map[int][]ProcessedUpdate[QuoteResult], err error) {
	results = make(map[int][]ProcessedUpdate[QuoteResult])
	duplicates = make(map[int][]ProcessedUpdate[QuoteResult])

	lastPrice := -1.

	for _, update := range updates {
		slot := update.Data.Data.ContextSlot
		price := update.Data.Data.Price(j.mint)

		qr := QuoteResult{
			Elapsed:   update.Timestamp.Sub(update.Data.Start),
			BuyPrice:  price,
			SellPrice: price,
			Source:    "jupiter",
		}
		pu := ProcessedUpdate[QuoteResult]{
			Timestamp: update.Data.Start,
			Slot:      slot,
			Data:      qr,
		}

		if price == lastPrice {
			duplicates[slot] = append(duplicates[slot], pu)
			if removeDuplicates {
				continue
			}
		}

		lastPrice = price
		results[slot] = append(results[slot], pu)
	}

	return
}

// FetchQuote is used to specify 1 USDC instead of 0.000001
func (j *JupiterAPIStream) FetchQuote(ctx context.Context) (jr JupiterPriceResponse, err error) {
	url := fmt.Sprintf("%v?inputMint=%v&outputMint=%v&amount=%v&slippageBps=%v", quoteAPIEndpoint, usdcMint, j.mint, j.adjustedAmount, j.slippageBps)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return
	}

	res, err := j.client.Do(req)
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	var quoteResponse JupiterQuoteResponse
	err = json.Unmarshal(b, &quoteResponse)
	if err != nil {
		return
	}

	if len(quoteResponse.Routes) == 0 {
		err = errors.New("no quotes found")
		return
	}

	bestRoute := quoteResponse.Routes[0]
	inAmount, err := strconv.ParseFloat(bestRoute.InAmount, 64)
	if err != nil {
		return
	}
	outAmount, err := strconv.ParseFloat(bestRoute.OutAmount, 64)
	if err != nil {
		return
	}
	price := inAmount / outAmount * j.priceAdjustment

	jr = JupiterPriceResponse{
		PriceInfo: map[string]JupiterPriceInfo{
			j.mint: {
				ID:            j.mint,
				MintSymbol:    "",
				VsToken:       usdcMint,
				VsTokenSymbol: "USDC",
				Price:         price,
			},
		},
		TimeTaken:   quoteResponse.TimeTaken,
		ContextSlot: int(quoteResponse.ContextSlot),
	}
	return
}

// FetchPrice returns a price based off of swapping 0.0000001 USDC (the minimum possible unit). Trader API does 1 USDC.
func (j *JupiterAPIStream) FetchPrice(ctx context.Context) (jr JupiterPriceResponse, err error) {
	url := fmt.Sprintf("%v?ids=%v", priceAPIEndpoint, j.mint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return
	}

	res, err := j.client.Do(req)
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &jr)
	if err != nil {
		return
	}

	return jr, nil
}
