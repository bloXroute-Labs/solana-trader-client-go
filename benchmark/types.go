package benchmark

import "time"

type PumpTxInfo struct {
	TimeSeen time.Time
}

type NewTokenResult struct {
	TraderAPIEventTime  time.Time
	ThirdPartyEventTime time.Time
	BlockTime           time.Time
	Diff                time.Duration
}
