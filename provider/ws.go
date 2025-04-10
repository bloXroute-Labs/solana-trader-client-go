package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/bloXroute-Labs/solana-trader-client-go/connections"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/gagliardetto/solana-go"
)

func (w *WSClient) RecentBlockHash(ctx context.Context) (*pb.GetRecentBlockHashResponse, error) {
	return w.recentBlockHashStore.get(ctx)
}

// GetTransaction returns details of a recent transaction
func (w *WSClient) GetTransaction(ctx context.Context, request *pb.GetTransactionRequest) (*pb.GetTransactionResponse, error) {
	var response pb.GetTransactionResponse
	err := w.conn.Request(ctx, "GetTransaction", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetRateLimit returns details of an account rate-limits
func (w *WSClient) GetRateLimit(ctx context.Context, request *pb.GetRateLimitRequest) (*pb.GetRateLimitResponse, error) {
	var response pb.GetRateLimitResponse
	err := w.conn.Request(ctx, "GetRateLimit", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetRaydiumPoolReserve returns pools details for a given set of pairs or addresses on Raydium
func (w *WSClient) GetRaydiumPoolReserve(ctx context.Context, req *pb.GetRaydiumPoolReserveRequest) (*pb.GetRaydiumPoolReserveResponse, error) {
	var response pb.GetRaydiumPoolReserveResponse
	err := w.conn.Request(ctx, "GetRaydiumPoolReserve", req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetRaydiumPools returns pools on Raydium
func (w *WSClient) GetRaydiumPools(ctx context.Context, request *pb.GetRaydiumPoolsRequest) (*pb.GetRaydiumPoolsResponse, error) {
	var response pb.GetRaydiumPoolsResponse
	err := w.conn.Request(ctx, "GetRaydiumPools", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetRaydiumQuotes returns the possible amount(s) of outToken for an inToken and the route to achieve it on Raydium
func (w *WSClient) GetRaydiumQuotes(ctx context.Context, request *pb.GetRaydiumQuotesRequest) (*pb.GetRaydiumQuotesResponse, error) {
	var response pb.GetRaydiumQuotesResponse
	err := w.conn.Request(ctx, "GetRaydiumQuotes", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetRaydiumQuotesCPMM returns the possible amount(s) of outToken for an inToken and the route to achieve it on Raydium CPMM pool
func (w *WSClient) GetRaydiumQuotesCPMM(ctx context.Context, request *pb.GetRaydiumCPMMQuotesRequest) (*pb.GetRaydiumCPMMQuotesResponse, error) {
	var response pb.GetRaydiumCPMMQuotesResponse
	err := w.conn.Request(ctx, "GetRaydiumQuotesCPMM", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetPumpFunQuotes returns the best quotes for swapping a token on PumpFun platform
func (w *WSClient) GetPumpFunQuotes(ctx context.Context, request *pb.GetPumpFunQuotesRequest) (*pb.GetPumpFunQuotesResponse, error) {
	var response pb.GetPumpFunQuotesResponse
	err := w.conn.Request(ctx, "GetPumpFunQuotes", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetRaydiumPrices returns the USDC price of requested tokens on Raydium
func (w *WSClient) GetRaydiumPrices(ctx context.Context, request *pb.GetRaydiumPricesRequest) (*pb.GetRaydiumPricesResponse, error) {
	var response pb.GetRaydiumPricesResponse
	err := w.conn.Request(ctx, "GetRaydiumPrices", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetRaydiumCLMMQuotes returns the CLMM quotes on Raydium
func (w *WSClient) GetRaydiumCLMMQuotes(ctx context.Context, request *pb.GetRaydiumCLMMQuotesRequest) (*pb.GetRaydiumCLMMQuotesResponse, error) {
	var response pb.GetRaydiumCLMMQuotesResponse
	err := w.conn.Request(ctx, "GetRaydiumCLMMQuotes", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetRaydiumCLMMPools returns the CLMM pools on Raydium
func (w *WSClient) GetRaydiumCLMMPools(ctx context.Context, request *pb.GetRaydiumCLMMPoolsRequest) (*pb.GetRaydiumCLMMPoolsResponse, error) {
	var response pb.GetRaydiumCLMMPoolsResponse
	err := w.conn.Request(ctx, "GetRaydiumCLMMPools", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostRaydiumCLMMSwap returns a partially signed transaction(s) for submitting a swap request on Raydium
func (w *WSClient) PostRaydiumCLMMSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest) (*pb.PostRaydiumSwapResponse, error) {
	var response pb.PostRaydiumSwapResponse
	err := w.conn.Request(ctx, "PostRaydiumCLMMSwap", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostRaydiumCLMMRouteSwap returns a partially signed transaction(s) for submitting a route swap request on Raydium
func (w *WSClient) PostRaydiumCLMMRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest) (*pb.PostRaydiumRouteSwapResponse, error) {
	var response pb.PostRaydiumRouteSwapResponse
	err := w.conn.Request(ctx, "PostRaydiumCLMMRouteSwap", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostRaydiumSwap returns a partially signed transaction(s) for submitting a swap request on Raydium
func (w *WSClient) PostRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest) (*pb.PostRaydiumSwapResponse, error) {
	var response pb.PostRaydiumSwapResponse
	err := w.conn.Request(ctx, "PostRaydiumSwap", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostRaydiumSwapCPMM returns a partially signed transaction(s) for submitting a swap request on Raydium
func (w *WSClient) PostRaydiumSwapCPMM(ctx context.Context, request *pb.PostRaydiumCPMMSwapRequest) (*pb.PostRaydiumCPMMSwapResponse, error) {
	var response pb.PostRaydiumCPMMSwapResponse
	err := w.conn.Request(ctx, "PostRaydiumCPMMSwap", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostPumpFunSwap returns a partially signed transaction(s) for submitting a swap request on Pumpdotfun platform
func (w *WSClient) PostPumpFunSwap(ctx context.Context, request *pb.PostPumpFunSwapRequest) (*pb.PostPumpFunSwapResponse, error) {
	var response pb.PostPumpFunSwapResponse
	err := w.conn.Request(ctx, "PostPumpFunSwap", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostRaydiumRouteSwap returns a partially signed transaction(s) for submitting a swap request on Raydium
func (w *WSClient) PostRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest) (*pb.PostRaydiumRouteSwapResponse, error) {
	var response pb.PostRaydiumRouteSwapResponse
	err := w.conn.Request(ctx, "PostRaydiumRouteSwap", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitRaydiumCLMMSwap builds a Raydium Swap transaction then signs it, and submits to the network.
func (w *WSClient) SubmitRaydiumCLMMSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := w.PostRaydiumCLMMSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return w.SignAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// SubmitRaydiumCLMMRouteSwap builds a Raydium RouteSwap transaction then signs it, and submits to the network.
func (w *WSClient) SubmitRaydiumCLMMRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := w.PostRaydiumCLMMRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return w.SignAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// GetJupiterQuotes returns the possible amount(s) of outToken for an inToken and the route to achieve it on Jupiter
func (w *WSClient) GetJupiterQuotes(ctx context.Context, request *pb.GetJupiterQuotesRequest) (*pb.GetJupiterQuotesResponse, error) {
	var response pb.GetJupiterQuotesResponse
	err := w.conn.Request(ctx, "GetJupiterQuotes", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetJupiterPrices returns the USDC price of requested tokens on Jupiter
func (w *WSClient) GetJupiterPrices(ctx context.Context, request *pb.GetJupiterPricesRequest) (*pb.GetJupiterPricesResponse, error) {
	var response pb.GetJupiterPricesResponse
	err := w.conn.Request(ctx, "GetJupiterPrices", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostJupiterSwap returns a partially signed transaction(s) for submitting a swap request on Jupiter
func (w *WSClient) PostJupiterSwap(ctx context.Context, request *pb.PostJupiterSwapRequest) (*pb.PostJupiterSwapResponse, error) {
	var response pb.PostJupiterSwapResponse
	err := w.conn.Request(ctx, "PostJupiterSwap", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostJupiterSwapInstructions returns instructions to build a transaction and submit it on jupiter
func (w *WSClient) PostJupiterSwapInstructions(ctx context.Context, request *pb.PostJupiterSwapInstructionsRequest) (*pb.PostJupiterSwapInstructionsResponse, error) {
	var response pb.PostJupiterSwapInstructionsResponse
	err := w.conn.Request(ctx, "PostJupiterSwapInstructions", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostRaydiumSwapInstructions returns instructions to build a transaction and submit it on raydium
func (w *WSClient) PostRaydiumSwapInstructions(ctx context.Context, request *pb.PostRaydiumSwapInstructionsRequest) (*pb.PostRaydiumSwapInstructionsResponse, error) {
	var response pb.PostRaydiumSwapInstructionsResponse
	err := w.conn.Request(ctx, "PostRaydiumSwapInstructions", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostJupiterRouteSwap returns a partially signed transaction(s) for submitting a swap request on Jupiter
func (w *WSClient) PostJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest) (*pb.PostJupiterRouteSwapResponse, error) {
	var response pb.PostJupiterRouteSwapResponse
	err := w.conn.Request(ctx, "PostJupiterRouteSwap", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetPools returns pools for given projects.
func (w *WSClient) GetPools(ctx context.Context, projects []pb.Project) (*pb.GetPoolsResponse, error) {
	response := pb.GetPoolsResponse{}
	err := w.conn.Request(ctx, "GetPools", &pb.GetPoolsRequest{Projects: projects}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetTokenAccounts returns all tokens associated with the owner address
func (w *WSClient) GetTokenAccounts(ctx context.Context, req *pb.GetTokenAccountsRequest) (*pb.GetTokenAccountsResponse, error) {
	var response pb.GetTokenAccountsResponse
	err := w.conn.Request(ctx, "GetTokenAccounts", req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetPriorityFee returns an suggested priority fee based on a given percentile
func (w *WSClient) GetPriorityFee(ctx context.Context, project pb.Project, percentile *float64) (*pb.GetPriorityFeeResponse, error) {
	request := &pb.GetPriorityFeeRequest{
		Project: project,
	}
	if percentile != nil {
		request.Percentile = percentile
	}
	var response pb.GetPriorityFeeResponse
	err := w.conn.Request(ctx, "GetPriorityFee", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetPriorityFeeByProgram returns priority fees based on a given list of programs
func (w *WSClient) GetPriorityFeeByProgram(ctx context.Context, programs []string) (*pb.GetPriorityFeeByProgramResponse, error) {
	request := &pb.GetPriorityFeeByProgramRequest{
		Programs: programs,
	}
	var response pb.GetPriorityFeeByProgramResponse
	err := w.conn.Request(ctx, "GetPriorityFeeByProgram", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetLeaderSchedule returns leader schedule given max Slots
func (w *WSClient) GetLeaderSchedule(ctx context.Context, maxSlots uint64) (*pb.GetLeaderScheduleResponse, error) {
	request := &pb.GetLeaderScheduleRequest{
		MaxSlots: maxSlots,
	}
	var response pb.GetLeaderScheduleResponse
	err := w.conn.Request(ctx, "GetLeaderSchedule", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostSubmit posts the transaction string to the Solana network.
func (w *WSClient) PostSubmit(ctx context.Context, txBase64 string, opts PostSubmitOpts) (*pb.PostSubmitResponse, error) {
	if w.privateKey == nil {
		return &pb.PostSubmitResponse{}, ErrPrivateKeyNotFound
	}

	request := &pb.PostSubmitRequest{
		Transaction: &pb.TransactionMessage{
			Content: txBase64,
		},
		SkipPreFlight:          opts.SkipPreFlight,
		FrontRunningProtection: &opts.FrontRunningProtection,
		UseStakedRPCs:          &opts.UseStakedRPCs,
		AllowBackRun:           &opts.AllowBackRun,
		RevenueAddress:         &opts.RevenueAddress,
		Sniping:                &opts.Sniping,
	}
	var response pb.PostSubmitResponse
	err := w.conn.Request(ctx, "PostSubmit", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (w *WSClient) PostSubmitSnipeV2(ctx context.Context, request *pb.PostSubmitSnipeRequest) (*pb.PostSubmitSnipeResponse, error) {
	if w.privateKey == nil {
		return &pb.PostSubmitSnipeResponse{}, ErrPrivateKeyNotFound
	}

	var response pb.PostSubmitSnipeResponse
	err := w.conn.Request(ctx, "PostSubmitSnipeV2", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostSubmitBatch posts a bundle of transactions string based on a specific SubmitStrategy to the Solana network.
func (w *WSClient) PostSubmitBatch(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error) {
	var response pb.PostSubmitBatchResponse
	err := w.conn.Request(ctx, "PostSubmitBatch", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostSubmitV2 posts the transaction string to the Solana network.
func (w *WSClient) PostSubmitV2(ctx context.Context, txBase64 string, opts PostSubmitOpts) (*pb.PostSubmitResponse, error) {
	if w.privateKey == nil {
		return &pb.PostSubmitResponse{}, ErrPrivateKeyNotFound
	}

	txBase64, err := transaction.SignTxWithPrivateKey(txBase64, *w.privateKey)
	if err != nil {
		return &pb.PostSubmitResponse{}, err
	}

	request := &pb.PostSubmitRequest{
		Transaction: &pb.TransactionMessage{
			Content: txBase64,
		},
		SkipPreFlight:          opts.SkipPreFlight,
		FrontRunningProtection: &opts.FrontRunningProtection,
		UseStakedRPCs:          &opts.UseStakedRPCs,
		AllowBackRun:           &opts.AllowBackRun,
		RevenueAddress:         &opts.RevenueAddress,
		Sniping:                &opts.Sniping,
	}
	var response pb.PostSubmitResponse
	err = w.conn.Request(ctx, "PostSubmitV2", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostSubmitBatchV2 posts a bundle of transactions string based on a specific SubmitStrategy to the Solana network.
func (w *WSClient) PostSubmitBatchV2(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error) {
	var response pb.PostSubmitBatchResponse
	err := w.conn.Request(ctx, "PostSubmitBatchV2", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SignAndSubmit signs the given transaction and submits it.
func (w *WSClient) SignAndSubmit(ctx context.Context, tx *pb.TransactionMessage,
	skipPreFlight bool, frontRunningProtection bool, useStakedRPCs bool) (string, error) {
	if w.privateKey == nil {
		return "", ErrPrivateKeyNotFound
	}

	txBase64, err := transaction.SignTxWithPrivateKey(tx.Content, *w.privateKey)
	if err != nil {
		return "", err
	}

	response, err := w.PostSubmit(ctx, txBase64, PostSubmitOpts{
		SkipPreFlight:          skipPreFlight,
		FrontRunningProtection: frontRunningProtection,
		UseStakedRPCs:          useStakedRPCs,
		// Other fields default to zero values
	})
	if err != nil {
		return "", err
	}

	return response.Signature, nil
}

func (w *WSClient) SignAndSubmitSnipe(ctx context.Context, transactions []*pb.TransactionMessage, useStakedRPCs bool) ([]string, error) {
	if w.privateKey == nil {
		return nil, ErrPrivateKeyNotFound
	}

	entries := make([]*pb.PostSubmitRequestEntry, len(transactions))

	for i, tx := range transactions {
		txBase64, err := transaction.SignTxWithPrivateKey(tx.Content, *w.privateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to sign transaction: %w", err)
		}

		entries[i] = &pb.PostSubmitRequestEntry{
			Transaction: &pb.TransactionMessage{
				Content:   txBase64,
				IsCleanup: tx.IsCleanup,
			},
			SkipPreFlight: false,
		}
	}

	snipeRequest := &pb.PostSubmitSnipeRequest{
		Entries:       entries,
		UseStakedRPCs: &useStakedRPCs,
	}

	response, err := w.PostSubmitSnipeV2(ctx, snipeRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to submit snipe request: %w", err)
	}

	signatures := make([]string, 0, len(response.Transactions))
	for _, entry := range response.Transactions {
		if entry.Submitted {
			signatures = append(signatures, entry.Signature)
		}
	}

	return signatures, nil
}

func (w *WSClient) SignAndSubmitPaladin(ctx context.Context, tx *pb.TransactionMessage, revertProtection *bool) (string, error) {
	if w.privateKey == nil {
		return "", ErrPrivateKeyNotFound
	}

	txBase64, err := transaction.SignTxWithPrivateKey(tx.Content, *w.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	paladinRequest := &pb.PostSubmitPaladinRequest{
		Transaction: &pb.TransactionMessageV2{
			Content: txBase64,
		},
		RevertProtection: revertProtection,
	}

	var response pb.PostSubmitResponse
	err = w.conn.Request(ctx, "PostSubmitPaladinV2", paladinRequest, &response)
	if err != nil {
		return "", err
	}
	return response.Signature, nil
}

// SignAndSubmitBatch signs the given transactions and submits them.
func (w *WSClient) SignAndSubmitBatch(ctx context.Context, transactions []*pb.TransactionMessage, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	if w.privateKey == nil {
		return nil, ErrPrivateKeyNotFound
	}

	if len(transactions) == 1 {
		signature, err := w.SignAndSubmit(ctx, transactions[0], *opts.SkipPreFlight, false, false)
		if err != nil {
			return nil, err
		}
		return &pb.PostSubmitBatchResponse{
			Transactions: []*pb.PostSubmitBatchResponseEntry{
				{
					Signature: signature,
					Error:     "",
					Submitted: true,
				},
			},
		}, nil
	}

	batchRequest, err := buildBatchRequest(transactions, *w.privateKey, useBundle, opts)
	if err != nil {
		return nil, err
	}
	return w.PostSubmitBatch(ctx, batchRequest)
}

// SubmitRaydiumSwap builds a Raydium Swap transaction then signs it, and submits to the network.
func (w *WSClient) SubmitRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := w.PostRaydiumSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return w.SignAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// SubmitRaydiumSwapCPMM builds a Raydium Swap CPMM transaction then signs it, and submits to the network.
func (w *WSClient) SubmitRaydiumSwapCPMM(ctx context.Context, request *pb.PostRaydiumCPMMSwapRequest) (string, error) {
	resp, err := w.PostRaydiumSwapCPMM(ctx, request)
	if err != nil {
		return "", err
	}

	sig, err := w.SignAndSubmit(ctx, resp.Transaction, true, false, false)
	if err != nil {
		return "", err
	}

	return sig, nil

}

// SubmitPostPumpFunSwap builds a pumpfun Swap transaction then signs it, and submits to the network.
func (w *WSClient) SubmitPostPumpFunSwap(ctx context.Context, request *pb.PostPumpFunSwapRequest) (string, error) {
	resp, err := w.PostPumpFunSwap(ctx, request)
	if err != nil {
		return "", err
	}
	return w.SignAndSubmit(ctx, &pb.TransactionMessage{
		Content: resp.Transaction.Content,
	}, false, false, false)
}

// SubmitRaydiumRouteSwap builds a Raydium RouteSwap transaction then signs it, and submits to the network.
func (w *WSClient) SubmitRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := w.PostRaydiumRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return w.SignAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// SubmitJupiterSwap builds a Jupiter Swap transaction then signs it, and submits to the network.
func (w *WSClient) SubmitJupiterSwap(ctx context.Context, request *pb.PostJupiterSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := w.PostJupiterSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return w.SignAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// SubmitJupiterSwapInstructions builds a Jupiter Swap transaction then signs it, and submits to the network.
func (w *WSClient) SubmitJupiterSwapInstructions(ctx context.Context, request *pb.PostJupiterSwapInstructionsRequest, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	swapInstructions, err := w.PostJupiterSwapInstructions(ctx, request)
	if err != nil {
		return nil, err
	}

	txBuilder := solana.NewTransactionBuilder()

	addressLookupTable, err := utils.ConvertProtoAddressLookupTable(swapInstructions.AddressLookupTableAddresses)
	if err != nil {
		return nil, err
	}

	txBuilder.WithOpt(solana.TransactionAddressTables(addressLookupTable))

	instructions, err := utils.ConvertJupiterInstructions(swapInstructions.Instructions)
	if err != nil {
		return nil, err
	}

	for _, inst := range instructions {
		txBuilder.AddInstruction(inst)
	}

	txBuilder.SetFeePayer(w.privateKey.PublicKey())
	blockHash, err := w.RecentBlockHash(ctx)

	if err != nil {
		panic(fmt.Errorf("server error: could not retrieve block hash: %w", err))
	}

	hash, err := solana.HashFromBase58(blockHash.BlockHash)
	if err != nil {
		return nil, err
	}

	txBuilder.SetRecentBlockHash(hash)
	tx, err := txBuilder.Build()
	if err != nil {
		return nil, err
	}

	err = transaction.PartialSign(tx, w.privateKey.PublicKey(), make(map[solana.PublicKey]solana.PrivateKey))
	if err != nil {
		return nil, err
	}

	var txToBeSigned []*pb.TransactionMessage

	txBase64, err := tx.ToBase64()
	if err != nil {
		return nil, err
	}

	txToBeSigned = append(txToBeSigned, &pb.TransactionMessage{
		Content:   txBase64,
		IsCleanup: false,
	})

	return w.SignAndSubmitBatch(ctx, txToBeSigned, useBundle, opts)
}

// SubmitRaydiumSwapInstructions builds a Raydium Swap transaction then signs it, and submits to the network.
func (w *WSClient) SubmitRaydiumSwapInstructions(ctx context.Context, request *pb.PostRaydiumSwapInstructionsRequest, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	swapInstructions, err := w.PostRaydiumSwapInstructions(ctx, request)
	if err != nil {
		return nil, err
	}

	instructions, err := utils.ConvertRaydiumInstructions(swapInstructions.Instructions)
	if err != nil {
		return nil, err
	}
	txBuilder := solana.NewTransactionBuilder()

	for _, inst := range instructions {
		txBuilder.AddInstruction(inst)
	}

	txBuilder.SetFeePayer(w.privateKey.PublicKey())
	blockHash, err := w.RecentBlockHash(ctx)

	if err != nil {
		panic(fmt.Errorf("server error: could not retrieve block hash: %w", err))
	}

	hash, err := solana.HashFromBase58(blockHash.BlockHash)
	if err != nil {
		return nil, err
	}

	txBuilder.SetRecentBlockHash(hash)
	tx, err := txBuilder.Build()
	if err != nil {
		return nil, err
	}

	err = transaction.PartialSign(tx, w.privateKey.PublicKey(), make(map[solana.PublicKey]solana.PrivateKey))
	if err != nil {
		return nil, err
	}

	var txToBeSigned []*pb.TransactionMessage

	txBase64, err := tx.ToBase64()
	if err != nil {
		return nil, err
	}

	txToBeSigned = append(txToBeSigned, &pb.TransactionMessage{
		Content:   txBase64,
		IsCleanup: false,
	})

	return w.SignAndSubmitBatch(ctx, txToBeSigned, useBundle, opts)
}

// SubmitJupiterRouteSwap builds a Jupiter RouteSwap transaction then signs it, and submits to the network.
func (w *WSClient) SubmitJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := w.PostJupiterRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return w.SignAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

func (w *WSClient) Close() error {
	return w.conn.Close(errors.New("shutdown requested"))
}

// GetPumpFunSwapsStream subscribes to a stream for swap events related to a set of pumpdotfun tokens
func (w *WSClient) GetPumpFunSwapsStream(ctx context.Context, req *pb.GetPumpFunSwapsStreamRequest) (connections.Streamer[*pb.GetPumpFunSwapsStreamResponse], error) {
	return connections.WSStreamProto(w.conn, ctx, "GetPumpFunSwapsStream", req, func() *pb.GetPumpFunSwapsStreamResponse {
		var v pb.GetPumpFunSwapsStreamResponse
		return &v
	})
}

// GetPumpFunNewAmmPoolStream subscribes to a stream for new AMM pool events
func (w *WSClient) GetPumpFunNewAmmPoolStream(ctx context.Context, req *pb.GetPumpFunNewAmmPoolStreamRequest) (connections.Streamer[*pb.GetPumpFunNewAmmPoolStreamResponse], error) {
	return connections.WSStreamProto(w.conn, ctx, "GetPumpFunNewAmmPoolStream", req, func() *pb.GetPumpFunNewAmmPoolStreamResponse {
		var v pb.GetPumpFunNewAmmPoolStreamResponse
		return &v
	})
}

// GetPumpFunNewTokensStream subscribes to a stream for pumpdotfun's new pool events
func (w *WSClient) GetPumpFunNewTokensStream(ctx context.Context, req *pb.GetPumpFunNewTokensStreamRequest) (connections.Streamer[*pb.GetPumpFunNewTokensStreamResponse], error) {
	return connections.WSStreamProto(w.conn, ctx, "GetPumpFunNewTokensStream", req, func() *pb.GetPumpFunNewTokensStreamResponse {
		var v pb.GetPumpFunNewTokensStreamResponse
		return &v
	})
}

// GetNewRaydiumPoolsStream subscribes to a stream for new Raydium Pools when they are created with
// option to include Raydium cpmm amm.
func (w *WSClient) GetNewRaydiumPoolsStream(ctx context.Context, includeCPMM bool) (connections.Streamer[*pb.GetNewRaydiumPoolsResponse], error) {
	return connections.WSStreamProto(w.conn, ctx, "GetNewRaydiumPoolsStream",
		&pb.GetNewRaydiumPoolsRequest{
			IncludeCPMM: &includeCPMM,
		},
		func() *pb.GetNewRaydiumPoolsResponse {
			var v pb.GetNewRaydiumPoolsResponse
			return &v
		})
}

// GetNewRaydiumPoolsByTranasctionStream subscribes to a stream for new Raydium Pools when they are created with some
// more detailed info compared to the standard straem, while sacrificing some speed
func (w *WSClient) GetNewRaydiumPoolsByTransactionStream(ctx context.Context) (connections.Streamer[*pb.GetNewRaydiumPoolsByTransactionResponse], error) {
	return connections.WSStreamProto(w.conn, ctx, "GetNewRaydiumPoolsByTransactionStream",
		&pb.GetNewRaydiumPoolsByTransactionRequest{},
		func() *pb.GetNewRaydiumPoolsByTransactionResponse {
			var v pb.GetNewRaydiumPoolsByTransactionResponse
			return &v
		})
}

// GetRecentBlockHashStream subscribes to a stream for getting recent block hash.
func (w *WSClient) GetRecentBlockHashStream(ctx context.Context) (connections.Streamer[*pb.GetRecentBlockHashResponse], error) {
	return connections.WSStreamProto(w.conn, ctx, "GetRecentBlockHashStream", &pb.GetRecentBlockHashRequest{}, func() *pb.GetRecentBlockHashResponse {
		return &pb.GetRecentBlockHashResponse{}
	})
}

// GetQuotesStream subscribes to a stream for getting recent quotes of tokens of interest.
func (w *WSClient) GetQuotesStream(ctx context.Context, projects []pb.Project, tokenPairs []*pb.TokenPair) (connections.Streamer[*pb.GetQuotesStreamResponse], error) {
	return connections.WSStreamProto(w.conn, ctx, "GetQuotesStream", &pb.GetQuotesStreamRequest{
		Projects:   projects,
		TokenPairs: tokenPairs,
	}, func() *pb.GetQuotesStreamResponse {
		return &pb.GetQuotesStreamResponse{}
	})
}

// GetPoolReservesStream subscribes to a stream for getting recent quotes of tokens of interest.
func (w *WSClient) GetPoolReservesStream(ctx context.Context, request *pb.GetPoolReservesStreamRequest) (connections.Streamer[*pb.GetPoolReservesStreamResponse], error) {
	return connections.WSStreamProto(w.conn, ctx, "GetPoolReservesStream", request, func() *pb.GetPoolReservesStreamResponse {
		return &pb.GetPoolReservesStreamResponse{}
	})
}

// GetPricesStream subscribes to a stream for getting recent quotes of tokens of interest.
func (w *WSClient) GetPricesStream(ctx context.Context, projects []pb.Project, tokens []string) (connections.Streamer[*pb.GetPricesStreamResponse], error) {
	return connections.WSStreamProto(w.conn, ctx, "GetPricesStream", &pb.GetPricesStreamRequest{
		Projects: projects,
		Tokens:   tokens,
	}, func() *pb.GetPricesStreamResponse {
		return &pb.GetPricesStreamResponse{}
	})
}

// GetSwapsStream subscribes to a stream for getting recent swaps on projects & markets of interest.
func (w *WSClient) GetSwapsStream(
	ctx context.Context,
	projects []pb.Project,
	markets []string,
	includeFailed bool,
) (connections.Streamer[*pb.GetSwapsStreamResponse], error) {
	return connections.WSStreamProto(w.conn, ctx, "GetSwapsStream", &pb.GetSwapsStreamRequest{
		Projects:      projects,
		Pools:         markets,
		IncludeFailed: includeFailed,
	}, func() *pb.GetSwapsStreamResponse {
		return &pb.GetSwapsStreamResponse{}
	})
}

// GetBlockStream subscribes to a stream for getting recent blocks.
func (w *WSClient) GetBlockStream(ctx context.Context) (connections.Streamer[*pb.GetBlockStreamResponse], error) {
	newResponse := func() *pb.GetBlockStreamResponse {
		return &pb.GetBlockStreamResponse{}
	}
	return connections.WSStreamProto(w.conn, ctx, "GetBlockStream", &pb.GetBlockStreamRequest{}, newResponse)
}

// GetPriorityFeeStream subscribes to a stream for getting a recent priority fee estimate based on a percentile.
func (w *WSClient) GetPriorityFeeStream(ctx context.Context, project pb.Project, percentile *float64) (connections.Streamer[*pb.GetPriorityFeeResponse], error) {
	request := &pb.GetPriorityFeeRequest{
		Project: project,
	}
	if percentile != nil {
		request.Percentile = percentile
	}
	return connections.WSStreamProto(w.conn, ctx, "GetPriorityFeeStream", request, func() *pb.GetPriorityFeeResponse {
		return &pb.GetPriorityFeeResponse{}
	})
}

func (w *WSClient) GetPriorityFeeByProgramStream(ctx context.Context, programs []string) (connections.Streamer[*pb.GetPriorityFeeByProgramResponse], error) {
	request := &pb.GetPriorityFeeByProgramRequest{
		Programs: programs,
	}

	return connections.WSStreamProto(w.conn, ctx, "GetPriorityFeeByProgramStream", request, func() *pb.GetPriorityFeeByProgramResponse {
		return &pb.GetPriorityFeeByProgramResponse{}
	})
}

// GetBundleTipStream subscribes to a stream of recent bundle tip percentiles
func (w *WSClient) GetBundleTipStream(ctx context.Context) (connections.Streamer[*pb.GetBundleTipResponse], error) {
	newResponse := func() *pb.GetBundleTipResponse {
		return &pb.GetBundleTipResponse{}
	}
	return connections.WSStreamProto(w.conn, ctx, "GetBundleTipStream", &pb.GetBundleTipRequest{}, newResponse)
}

// GetRecentBlockHash returns recent block hash.
func (w *WSClient) GetRecentBlockHash(ctx context.Context, request *pb.GetRecentBlockHashRequest) (*pb.GetRecentBlockHashResponse, error) {
	var response pb.GetRecentBlockHashResponse
	err := w.conn.Request(ctx, "GetRecentBlockHash", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetRecentBlockHashV2 returns recent block hash, supports optional offset.
func (w *WSClient) GetRecentBlockHashV2(ctx context.Context, request *pb.GetRecentBlockHashRequestV2) (*pb.GetRecentBlockHashResponseV2, error) {
	var response pb.GetRecentBlockHashResponseV2
	err := w.conn.Request(ctx, "GetRecentBlockHashV2", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
