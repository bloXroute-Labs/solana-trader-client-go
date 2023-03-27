package stream

type JupiterQuoteResponse struct {
	Routes      []JupiterRoute `json:"data"`
	TimeTaken   float64        `json:"timeTaken"`
	ContextSlot uint64         `json:"contextSlot"`
}

type JupiterRoute struct {
	InAmount             string              `json:"inAmount"`
	OutAmount            string              `json:"outAmount"`
	PriceImpactPct       float64             `json:"priceImpactPct"`
	MarketInfos          []JupiterMarketInfo `json:"marketInfos"`
	Amount               string              `json:"amount"`
	SlippageBps          int                 `json:"slippageBps"`
	OtherAmountThreshold string              `json:"otherAmountThreshold"`
	SwapMode             string              `json:"swapMode"`
}

type JupiterMarketInfo struct {
	ID                 string     `json:"id"`
	Label              string     `json:"label"`
	InputMint          string     `json:"inputMint"`
	OutputMint         string     `json:"outputMint"`
	NotEnoughLiquidity bool       `json:"notEnoughLiquidity"`
	InAmount           string     `json:"inAmount"`
	OutAmount          string     `json:"outAmount"`
	PriceImpactPct     float64    `json:"priceImpactPct"`
	LpFee              JupiterFee `json:"lpFee"`
	PlatformFee        JupiterFee `json:"platformFee"`
}

type JupiterFee struct {
	Amount string  `json:"amount"`
	Mint   string  `json:"mint"`
	Pct    float64 `json:"pct"`
}
