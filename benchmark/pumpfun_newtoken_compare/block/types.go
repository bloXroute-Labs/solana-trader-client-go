package block

type FullBlock struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Result struct {
			Context struct {
				Slot int `json:"slot"`
			} `json:"context"`
			Value struct {
				Slot  int `json:"slot"`
				Block struct {
					PreviousBlockhash string `json:"previousBlockhash"`
					Blockhash         string `json:"blockhash"`
					ParentSlot        int    `json:"parentSlot"`
					Txs               []Tx   `json:"transactions"`
					Rewards           []struct {
						Pubkey      string      `json:"pubkey"`
						Lamports    int         `json:"lamports"`
						PostBalance int64       `json:"postBalance"`
						RewardType  string      `json:"rewardType"`
						Commission  interface{} `json:"commission"`
					} `json:"rewards"`
					BlockTime   int `json:"blockTime"`
					BlockHeight int `json:"blockHeight"`
				} `json:"block"`
				Err interface{} `json:"err"`
			} `json:"value"`
		} `json:"result"`
		Subscription int `json:"subscription"`
	} `json:"params"`
}

type Tx struct {
	Transaction []string `json:"transaction"`
	Meta        struct {
		Err    interface{} `json:"err"`
		Status struct {
			Ok interface{} `json:"Ok"`
		} `json:"status"`
		Fee               int           `json:"fee"`
		PreBalances       []int         `json:"preBalances"`
		PostBalances      []int         `json:"postBalances"`
		InnerInstructions []interface{} `json:"innerInstructions"`
		LogMessages       []string      `json:"logMessages"`
		PreTokenBalances  []interface{} `json:"preTokenBalances"`
		PostTokenBalances []interface{} `json:"postTokenBalances"`
		Rewards           []interface{} `json:"rewards"`
		LoadedAddresses   struct {
			Writable []interface{} `json:"writable"`
			Readonly []interface{} `json:"readonly"`
		} `json:"loadedAddresses"`
		ComputeUnitsConsumed int `json:"computeUnitsConsumed"`
	} `json:"meta"`
	Version interface{} `json:"version"`
}

type HeliusTx struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  struct {
		Subscription int64 `json:"subscription"`
		Result       struct {
			Transaction struct {
				Transaction []string `json:"transaction"`
				Meta        struct {
					Err    interface{} `json:"err"`
					Status struct {
						Ok interface{} `json:"Ok"`
					} `json:"status"`
					Fee               int     `json:"fee"`
					PreBalances       []int64 `json:"preBalances"`
					PostBalances      []int64 `json:"postBalances"`
					InnerInstructions []struct {
						Index        int `json:"index"`
						Instructions []struct {
							ProgramIdIndex int    `json:"programIdIndex"`
							Accounts       []int  `json:"accounts"`
							Data           string `json:"data"`
							StackHeight    int    `json:"stackHeight"`
						} `json:"instructions"`
					} `json:"innerInstructions"`
					LogMessages       []string      `json:"logMessages"`
					PreTokenBalances  []interface{} `json:"preTokenBalances"`
					PostTokenBalances []struct {
						AccountIndex  int    `json:"accountIndex"`
						Mint          string `json:"mint"`
						UiTokenAmount struct {
							UiAmount       float64 `json:"uiAmount"`
							Decimals       int     `json:"decimals"`
							Amount         string  `json:"amount"`
							UiAmountString string  `json:"uiAmountString"`
						} `json:"uiTokenAmount"`
						Owner     string `json:"owner"`
						ProgramId string `json:"programId"`
					} `json:"postTokenBalances"`
					Rewards         []interface{} `json:"rewards"`
					LoadedAddresses struct {
						Writable []interface{} `json:"writable"`
						Readonly []interface{} `json:"readonly"`
					} `json:"loadedAddresses"`
					ComputeUnitsConsumed int `json:"computeUnitsConsumed"`
				} `json:"meta"`
				Version int `json:"version"`
			} `json:"transaction"`
			Signature string `json:"signature"`
			Slot      int    `json:"slot"`
		} `json:"result"`
	} `json:"params"`
}
