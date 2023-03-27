package stream

type JupiterResponse struct {
	Routes      []jupiterRoute `json:"data"`
	TimeTaken   float64        `json:"timeTaken"`
	ContextSlot uint64         `json:"contextSlot"`
}

type jupiterRoute struct {
	InAmount             string              `json:"inAmount"`
	OutAmount            string              `json:"outAmount"`
	PriceImpactPct       float64             `json:"priceImpactPct"`
	MarketInfos          []jupiterMarketInfo `json:"marketInfos"`
	Amount               string              `json:"amount"`
	SlippageBps          int                 `json:"slippageBps"`
	OtherAmountThreshold string              `json:"otherAmountThreshold"`
	SwapMode             string              `json:"swapMode"`
}

type jupiterMarketInfo struct {
	ID                 string     `json:"id"`
	Label              string     `json:"label"`
	InputMint          string     `json:"inputMint"`
	OutputMint         string     `json:"outputMint"`
	NotEnoughLiquidity bool       `json:"notEnoughLiquidity"`
	InAmount           string     `json:"inAmount"`
	OutAmount          string     `json:"outAmount"`
	PriceImpactPct     float64    `json:"priceImpactPct"`
	LpFee              jupiterFee `json:"lpFee"`
	PlatformFee        jupiterFee `json:"platformFee"`
}

type jupiterFee struct {
	Amount string `json:"amount"`
	Mint   string `json:"mint"`
	Pct    int    `json:"pct"`
}
