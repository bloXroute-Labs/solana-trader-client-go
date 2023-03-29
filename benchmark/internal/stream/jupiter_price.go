package stream

type JupiterPriceResponse struct {
	PriceInfo   map[string]JupiterPriceInfo `json:"data"`
	TimeTaken   float64                     `json:"timeTaken"`
	ContextSlot int                         `json:"contextSlot"`
}

func (jr JupiterPriceResponse) Price(mint string) float64 {
	return jr.PriceInfo[mint].Price
}

type JupiterPriceInfo struct {
	ID            string  `json:"id"`
	MintSymbol    string  `json:"mintSymbol"`
	VsToken       string  `json:"vsToken"`
	VsTokenSymbol string  `json:"vsTokenSymbol"`
	Price         float64 `json:"price"`
}
