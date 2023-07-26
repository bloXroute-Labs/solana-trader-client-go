package stream

import "time"

type QuoteResult struct {
	Elapsed   time.Duration
	BuyPrice  float64
	SellPrice float64
	Source    string
}
