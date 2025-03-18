package provider

import (
	"context"
	"fmt"

	"github.com/bloXroute-Labs/solana-trader-client-go/connections"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/gagliardetto/solana-go"
)

func (g *GRPCClient) RecentBlockHash(ctx context.Context) (*pb.GetRecentBlockHashResponse, error) {
	return g.recentBlockHashStore.get(ctx)
}

// GetRecentBlockHash returns recent block hash.
func (g *GRPCClient) GetRecentBlockHash(ctx context.Context) (*pb.GetRecentBlockHashResponse, error) {
	return g.apiClient.GetRecentBlockHash(ctx, &pb.GetRecentBlockHashRequest{})
}

// GetRecentBlockHash returns recent block hash, supports optional offset.
func (g *GRPCClient) GetRecentBlockHashV2(ctx context.Context, offset uint64) (*pb.GetRecentBlockHashResponseV2, error) {
	return g.apiClient.GetRecentBlockHashV2(ctx, &pb.GetRecentBlockHashRequestV2{Offset: offset})
}

// GetPriorityFee returns a priority fee estimate for a given percentile
func (g *GRPCClient) GetPriorityFee(ctx context.Context, request *pb.GetPriorityFeeRequest) (*pb.GetPriorityFeeResponse, error) {
	return g.apiClient.GetPriorityFee(ctx, request)
}

// GetPriorityFeeByProgram returns priority fees for given programs
func (g *GRPCClient) GetPriorityFeeByProgram(ctx context.Context, request *pb.GetPriorityFeeByProgramRequest) (*pb.GetPriorityFeeByProgramResponse, error) {
	return g.apiClient.GetPriorityFeeByProgram(ctx, request)
}

// GetLeaderSchedule returns leader schedule for given max slots
func (g *GRPCClient) GetLeaderSchedule(ctx context.Context, request *pb.GetLeaderScheduleRequest) (*pb.GetLeaderScheduleResponse, error) {
	return g.apiClient.GetLeaderSchedule(ctx, request)
}

// GetRateLimit returns details of an account rate-limits
func (g *GRPCClient) GetRateLimit(ctx context.Context, request *pb.GetRateLimitRequest) (*pb.GetRateLimitResponse, error) {
	return g.apiClient.GetRateLimit(ctx, request)
}

// GetTransaction returns details of a recent transaction
func (g *GRPCClient) GetTransaction(ctx context.Context, request *pb.GetTransactionRequest) (*pb.GetTransactionResponse, error) {
	return g.apiClient.GetTransaction(ctx, request)
}

// GetRaydiumPoolReserve returns pools details for a given set of pairs or addresses on Raydium
func (g *GRPCClient) GetRaydiumPoolReserve(ctx context.Context, req *pb.GetRaydiumPoolReserveRequest) (*pb.GetRaydiumPoolReserveResponse, error) {
	return g.apiClient.GetRaydiumPoolReserve(ctx, req)
}

// GetRaydiumPools returns pools on Raydium
func (g *GRPCClient) GetRaydiumPools(ctx context.Context, request *pb.GetRaydiumPoolsRequest) (*pb.GetRaydiumPoolsResponse, error) {
	return g.apiClient.GetRaydiumPools(ctx, request)
}

// GetRaydiumQuotes returns the possible amount(s) of outToken for an inToken and the route to achieve it on Raydium
func (g *GRPCClient) GetRaydiumQuotes(ctx context.Context, request *pb.GetRaydiumQuotesRequest) (*pb.GetRaydiumQuotesResponse, error) {
	return g.apiClient.GetRaydiumQuotes(ctx, request)
}

// GetRaydiumQuotesCPMM returns the possible amount(s) of outToken for an inToken and the route to achieve it on Raydium CPMM pool
func (g *GRPCClient) GetRaydiumQuotesCPMM(ctx context.Context, request *pb.GetRaydiumCPMMQuotesRequest) (*pb.GetRaydiumCPMMQuotesResponse, error) {
	return g.apiClient.GetRaydiumCPMMQuotes(ctx, request)
}

// GetPumpFunQuotes returns the best quotes for swapping a token on PumpFun platform
func (g *GRPCClient) GetPumpFunQuotes(ctx context.Context, request *pb.GetPumpFunQuotesRequest) (*pb.GetPumpFunQuotesResponse, error) {
	return g.apiClient.GetPumpFunQuotes(ctx, request)
}

// GetRaydiumPrices returns the USDC price of requested tokens on Raydium
func (g *GRPCClient) GetRaydiumPrices(ctx context.Context, request *pb.GetRaydiumPricesRequest) (*pb.GetRaydiumPricesResponse, error) {
	return g.apiClient.GetRaydiumPrices(ctx, request)
}

// SubmitRaydiumCLMMSwap builds a Raydium Swap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitRaydiumCLMMSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.apiClient.PostRaydiumCLMMSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// SubmitRaydiumCLMMRouteSwap builds a Raydium RouteSwap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitRaydiumCLMMRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.apiClient.PostRaydiumCLMMRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// PostRaydiumSwap returns a partially signed transaction(s) for submitting a swap request on Raydium
func (g *GRPCClient) PostRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest) (*pb.PostRaydiumSwapResponse, error) {
	return g.apiClient.PostRaydiumSwap(ctx, request)
}

// PostPumpFunSwap returns a partially signed transaction(s) for submitting a swap request on Pumpdotfun platform
func (g *GRPCClient) PostPumpFunSwap(ctx context.Context, request *pb.PostPumpFunSwapRequest) (*pb.PostPumpFunSwapResponse, error) {
	return g.apiClient.PostPumpFunSwap(ctx, request)
}

// PostRaydiumSwapCPMM returns a partially signed transaction(s) for submitting a swap request on Raydium CPMM Pool
func (g *GRPCClient) PostRaydiumSwapCPMM(ctx context.Context, request *pb.PostRaydiumCPMMSwapRequest) (*pb.PostRaydiumCPMMSwapResponse, error) {
	return g.apiClient.PostRaydiumCPMMSwap(ctx, request)
}

// PostRaydiumRouteSwap returns a partially signed transaction(s) for submitting a swap request on Raydium
func (g *GRPCClient) PostRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest) (*pb.PostRaydiumRouteSwapResponse, error) {
	return g.apiClient.PostRaydiumRouteSwap(ctx, request)
}

// GetJupiterQuotes returns the possible amount(s) of outToken for an inToken and the route to achieve it on Jupiter
func (g *GRPCClient) GetJupiterQuotes(ctx context.Context, request *pb.GetJupiterQuotesRequest) (*pb.GetJupiterQuotesResponse, error) {
	return g.apiClient.GetJupiterQuotes(ctx, request)
}

// GetJupiterPrices returns the USDC price of requested tokens on Jupiter
func (g *GRPCClient) GetJupiterPrices(ctx context.Context, request *pb.GetJupiterPricesRequest) (*pb.GetJupiterPricesResponse, error) {
	return g.apiClient.GetJupiterPrices(ctx, request)
}

// PostJupiterSwap returns a partially signed transaction(s) for submitting a swap request on Jupiter
func (g *GRPCClient) PostJupiterSwap(ctx context.Context, request *pb.PostJupiterSwapRequest) (*pb.PostJupiterSwapResponse, error) {
	return g.apiClient.PostJupiterSwap(ctx, request)
}

// PostJupiterSwapInstructions returns instructions to build a transaction and submit it on jupiter
func (g *GRPCClient) PostJupiterSwapInstructions(ctx context.Context, request *pb.PostJupiterSwapInstructionsRequest) (*pb.PostJupiterSwapInstructionsResponse, error) {
	return g.apiClient.PostJupiterSwapInstructions(ctx, request)
}

// PostRaydiumSwapInstructions returns instructions to build a transaction and submit it on raydium
func (g *GRPCClient) PostRaydiumSwapInstructions(ctx context.Context, request *pb.PostRaydiumSwapInstructionsRequest) (*pb.PostRaydiumSwapInstructionsResponse, error) {
	return g.apiClient.PostRaydiumSwapInstructions(ctx, request)
}

// PostJupiterRouteSwap returns a partially signed transaction(s) for submitting a swap request on Jupiter
func (g *GRPCClient) PostJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest) (*pb.PostJupiterRouteSwapResponse, error) {
	return g.apiClient.PostJupiterRouteSwap(ctx, request)
}

// GetPools returns pools for given projects.
func (g *GRPCClient) GetPools(ctx context.Context, projects []pb.Project) (*pb.GetPoolsResponse, error) {
	return g.apiClient.GetPools(ctx, &pb.GetPoolsRequest{Projects: projects})
}

// GetTokenAccounts returns all tokens associated with the owner address
func (g *GRPCClient) GetTokenAccounts(ctx context.Context, req *pb.GetTokenAccountsRequest) (*pb.GetTokenAccountsResponse, error) {
	return g.apiClient.GetTokenAccounts(ctx, req)
}

// GetPrice returns the USDC price of requested tokens
func (g *GRPCClient) GetPrice(ctx context.Context, tokens []string) (*pb.GetPriceResponse, error) {
	return g.apiClient.GetPrice(ctx, &pb.GetPriceRequest{Tokens: tokens})
}

// SignAndSubmit signs the given transaction and submits it.
func (g *GRPCClient) SignAndSubmit(ctx context.Context, tx *pb.TransactionMessage,
	skipPreFlight bool, frontRunningProtection bool, useStakedRPCs bool) (string, error) {
	if g.privateKey == nil {
		return "", ErrPrivateKeyNotFound
	}
	txBase64, err := transaction.SignTxWithPrivateKey(tx.Content, *g.privateKey)
	if err != nil {
		return "", err
	}

	response, err := g.PostSubmit(ctx, &pb.TransactionMessage{
		Content:   txBase64,
		IsCleanup: tx.IsCleanup,
	}, PostSubmitOpts{
		SkipPreFlight:          skipPreFlight,
		FrontRunningProtection: frontRunningProtection,
		UseStakedRPCs:          useStakedRPCs,
		// Using zero values for other fields:
		// AllowBackRun:    false,
		// RevenueAddress:  "",
		// Sniping:         false,
		// AllowRevert:     false,
	})
	if err != nil {
		return "", err
	}

	return response.Signature, nil
}

func (g *GRPCClient) SignAndSubmitSnipe(ctx context.Context, transactions []*pb.TransactionMessage, useStakedRPCs bool) ([]string, error) {
	if g.privateKey == nil {
		return nil, ErrPrivateKeyNotFound
	}

	entries := make([]*pb.PostSubmitRequestEntry, len(transactions))

	for i, tx := range transactions {
		txBase64, err := transaction.SignTxWithPrivateKey(tx.Content, *g.privateKey)
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

	response, err := g.apiClient.PostSubmitSnipeV2(ctx, snipeRequest)
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

func (g *GRPCClient) SignAndSubmitPaladin(ctx context.Context, tx *pb.TransactionMessage) (string, error) {
	if g.privateKey == nil {
		return "", ErrPrivateKeyNotFound
	}
	txBase64, err := transaction.SignTxWithPrivateKey(tx.Content, *g.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}
	paladinRequest := &pb.PostSubmitPaladinRequest{
		Transaction: &pb.TransactionMessageV2{
			Content: txBase64,
		},
	}
	response, err := g.apiClient.PostSubmitPaladinV2(ctx, paladinRequest)
	if err != nil {
		return "", fmt.Errorf("failed to submit paladin request: %w", err)
	}
	return response.Signature, nil
}

// signAndSubmitBatch signs the given transactions and submits them.
func (g *GRPCClient) signAndSubmitBatch(ctx context.Context, transactions []*pb.TransactionMessage, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	if g.privateKey == nil {
		return nil, ErrPrivateKeyNotFound
	}

	if len(transactions) == 1 {
		println("here")
		signature, err := g.SignAndSubmit(ctx, transactions[0], *opts.SkipPreFlight, false, false)
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

	batchRequest, err := buildBatchRequest(transactions, *g.privateKey, useBundle, opts)
	if err != nil {
		return nil, err
	}

	return g.PostSubmitBatch(ctx, batchRequest)
}

// PostSubmit posts the transaction string to the Solana network.
func (g *GRPCClient) PostSubmit(ctx context.Context, tx *pb.TransactionMessage, opts PostSubmitOpts) (*pb.PostSubmitResponse, error) {
	return g.apiClient.PostSubmit(ctx, &pb.PostSubmitRequest{
		Transaction:            tx,
		SkipPreFlight:          opts.SkipPreFlight,
		FrontRunningProtection: &opts.FrontRunningProtection,
		UseStakedRPCs:          &opts.UseStakedRPCs,
		AllowBackRun:           &opts.AllowBackRun,
		RevenueAddress:         &opts.RevenueAddress,
		Sniping:                &opts.Sniping,
	})
}

// PostSubmitBatch posts a bundle of transactions string based on a specific SubmitStrategy to the Solana network.
func (g *GRPCClient) PostSubmitBatch(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error) {
	return g.apiClient.PostSubmitBatch(ctx, request)
}

// PostSubmitV2 posts the transaction string to the Solana network.
func (g *GRPCClient) PostSubmitV2(ctx context.Context, tx *pb.TransactionMessage, opts PostSubmitOpts) (*pb.PostSubmitResponse, error) {
	return g.apiClient.PostSubmitV2(ctx, &pb.PostSubmitRequest{
		Transaction:            tx,
		SkipPreFlight:          opts.SkipPreFlight,
		FrontRunningProtection: &opts.FrontRunningProtection,
	})
}

// PostSubmitSnipeV2 posts the transaction string to the Solana network.
func (g *GRPCClient) PostSubmitSnipeV2(ctx context.Context, request *pb.PostSubmitSnipeRequest) (*pb.PostSubmitSnipeResponse, error) {
	return g.apiClient.PostSubmitSnipeV2(ctx, request)
}

// PostSubmitBatchV2 posts a bundle of transactions string based on a specific SubmitStrategy to the Solana network.
func (g *GRPCClient) PostSubmitBatchV2(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error) {
	return g.apiClient.PostSubmitBatchV2(ctx, request)
}

// SubmitRaydiumSwap builds a Raydium Swap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostRaydiumSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// SubmitPostPumpFunSwap builds a pumpfun Swap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitPostPumpFunSwap(ctx context.Context, request *pb.PostPumpFunSwapRequest) (string, error) {
	resp, err := g.PostPumpFunSwap(ctx, request)
	if err != nil {
		return "", err
	}
	return g.SignAndSubmit(ctx, &pb.TransactionMessage{
		Content: resp.Transaction.Content,
	}, false, false, false)
}

// SubmitRaydiumSwapCPMM builds a Raydium Swap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitRaydiumSwapCPMM(ctx context.Context, request *pb.PostRaydiumCPMMSwapRequest) (string, error) {
	resp, err := g.PostRaydiumSwapCPMM(ctx, request)
	if err != nil {
		return "", err
	}

	sig, err := g.SignAndSubmit(ctx, resp.Transaction, true, false, false)
	if err != nil {
		return "", err
	}

	return sig, nil
}

// SubmitRaydiumRouteSwap builds a Raydium RouteSwap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostRaydiumRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// SubmitJupiterSwap builds a Jupiter Swap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitJupiterSwap(ctx context.Context, request *pb.PostJupiterSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostJupiterSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// SubmitJupiterSwapInstructions builds a Jupiter Swap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitJupiterSwapInstructions(ctx context.Context, request *pb.PostJupiterSwapInstructionsRequest, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	swapInstructions, err := g.PostJupiterSwapInstructions(ctx, request)
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

	txBuilder.SetFeePayer(g.privateKey.PublicKey())
	blockHash, err := g.RecentBlockHash(ctx)

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

	err = transaction.PartialSign(tx, g.privateKey.PublicKey(), make(map[solana.PublicKey]solana.PrivateKey))
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

	return g.signAndSubmitBatch(ctx, txToBeSigned, useBundle, opts)
}

// SubmitRaydiumSwapInstructions builds a Raydium Swap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitRaydiumSwapInstructions(ctx context.Context, request *pb.PostRaydiumSwapInstructionsRequest, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	swapInstructions, err := g.PostRaydiumSwapInstructions(ctx, request)
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

	txBuilder.SetFeePayer(g.privateKey.PublicKey())
	blockHash, err := g.RecentBlockHash(ctx)

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

	err = transaction.PartialSign(tx, g.privateKey.PublicKey(), make(map[solana.PublicKey]solana.PrivateKey))
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

	return g.signAndSubmitBatch(ctx, txToBeSigned, useBundle, opts)
}

// SubmitJupiterRouteSwap builds a Jupiter RouteSwap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostJupiterRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// GetPumpFunSwapsStream subscribes to a stream for swap events related to a set of pumpdotfun tokens
func (g *GRPCClient) GetPumpFunSwapsStream(ctx context.Context, req *pb.GetPumpFunSwapsStreamRequest) (connections.Streamer[*pb.GetPumpFunSwapsStreamResponse], error) {
	stream, err := g.apiClient.GetPumpFunSwapsStream(ctx, req)
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetPumpFunSwapsStreamResponse](stream, ""), nil
}

// GetPumpFunNewTokensStream subscribes to a stream for pumpdotfun's new pool events
func (g *GRPCClient) GetPumpFunNewTokensStream(ctx context.Context, req *pb.GetPumpFunNewTokensStreamRequest) (connections.Streamer[*pb.GetPumpFunNewTokensStreamResponse], error) {
	stream, err := g.apiClient.GetPumpFunNewTokensStream(ctx, req)
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetPumpFunNewTokensStreamResponse](stream, ""), nil
}

// GetRecentBlockHashStream subscribes to a stream for getting recent block hash.
func (g *GRPCClient) GetRecentBlockHashStream(ctx context.Context) (connections.Streamer[*pb.GetRecentBlockHashResponse], error) {
	stream, err := g.apiClient.GetRecentBlockHashStream(ctx, &pb.GetRecentBlockHashRequest{})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetRecentBlockHashResponse](stream, ""), nil
}

// GetQuotesStream subscribes to a stream for getting recent quotes of tokens of interest.
func (g *GRPCClient) GetQuotesStream(ctx context.Context, projects []pb.Project, tokenPairs []*pb.TokenPair) (connections.Streamer[*pb.GetQuotesStreamResponse], error) {
	stream, err := g.apiClient.GetQuotesStream(ctx, &pb.GetQuotesStreamRequest{
		Projects:   projects,
		TokenPairs: tokenPairs,
	})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetQuotesStreamResponse](stream, ""), nil
}

// GetPoolReservesStream subscribes to a stream for getting recent quotes of tokens of interest.
func (g *GRPCClient) GetPoolReservesStream(ctx context.Context, request *pb.GetPoolReservesStreamRequest) (connections.Streamer[*pb.GetPoolReservesStreamResponse], error) {
	stream, err := g.apiClient.GetPoolReservesStream(ctx, request)
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetPoolReservesStreamResponse](stream, ""), nil
}

// GetPricesStream subscribes to a stream for getting recent prices of tokens of interest.
func (g *GRPCClient) GetPricesStream(ctx context.Context, projects []pb.Project, tokens []string) (connections.Streamer[*pb.GetPricesStreamResponse], error) {
	stream, err := g.apiClient.GetPricesStream(ctx, &pb.GetPricesStreamRequest{
		Projects: projects,
		Tokens:   tokens,
	})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetPricesStreamResponse](stream, ""), nil
}

// GetSwapsStream subscribes to a stream for getting recent swaps on projects & markets of interest.
func (g *GRPCClient) GetSwapsStream(
	ctx context.Context,
	projects []pb.Project,
	markets []string,
	includeFailed bool,
) (connections.Streamer[*pb.GetSwapsStreamResponse], error) {
	stream, err := g.apiClient.GetSwapsStream(ctx, &pb.GetSwapsStreamRequest{
		Projects:      projects,
		Pools:         markets,
		IncludeFailed: includeFailed,
	})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetSwapsStreamResponse](stream, ""), nil
}

// GetNewRaydiumPoolsStream subscribes to a stream for getting recent swaps on projects & markets of interest with
// option to include Raydium cpmm amm.
func (g *GRPCClient) GetNewRaydiumPoolsStream(
	ctx context.Context, includeCPMM bool,
) (connections.Streamer[*pb.GetNewRaydiumPoolsResponse], error) {
	stream, err := g.apiClient.GetNewRaydiumPoolsStream(ctx, &pb.GetNewRaydiumPoolsRequest{
		IncludeCPMM: &includeCPMM,
	})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetNewRaydiumPoolsResponse](stream, ""), nil
}

// GetNewRaydiumPoolsStreamByTransaction subscribes to a stream for getting recent swaps on projects & markets of interest.
// The ByTransaction option gives a bit more pool info while sacrificing speed

func (g *GRPCClient) GetNewRaydiumPoolsByTransactionStream(
	ctx context.Context, includeCPMM bool,
) (connections.Streamer[*pb.GetNewRaydiumPoolsByTransactionResponse], error) {
	stream, err := g.apiClient.GetNewRaydiumPoolsByTransactionStream(ctx, &pb.GetNewRaydiumPoolsByTransactionRequest{})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetNewRaydiumPoolsByTransactionResponse](stream, ""), nil
}

// GetBlockStream subscribes to a stream for getting recent blocks.
func (g *GRPCClient) GetBlockStream(ctx context.Context) (connections.Streamer[*pb.GetBlockStreamResponse], error) {
	stream, err := g.apiClient.GetBlockStream(ctx, &pb.GetBlockStreamRequest{})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetBlockStreamResponse](stream, ""), nil
}

// GetPriorityFeeStream subscribes to a stream of priority fees for a given percentile
func (g *GRPCClient) GetPriorityFeeStream(ctx context.Context, project pb.Project, percentile *float64) (connections.Streamer[*pb.GetPriorityFeeResponse], error) {
	request := &pb.GetPriorityFeeRequest{
		Project: project,
	}
	if percentile != nil {
		request.Percentile = percentile
	}
	stream, err := g.apiClient.GetPriorityFeeStream(ctx, request)
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetPriorityFeeResponse](stream, fmt.Sprint(percentile)), nil
}

func (g *GRPCClient) GetPriorityFeeByProgramStream(ctx context.Context, programs []string) (connections.Streamer[*pb.GetPriorityFeeByProgramResponse], error) {
	request := &pb.GetPriorityFeeByProgramRequest{
		Programs: programs,
	}

	stream, err := g.apiClient.GetPriorityFeeByProgramStream(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("GetPriorityFeeByProgramStream error: %v", err)
	}

	return connections.GRPCStream[pb.GetPriorityFeeByProgramResponse](stream, "programs"), nil
}

// GetBundleTipStream subscribes to a stream of bundle tip percentiles
func (g *GRPCClient) GetBundleTipStream(ctx context.Context) (connections.Streamer[*pb.GetBundleTipResponse], error) {
	stream, err := g.apiClient.GetBundleTipStream(ctx, &pb.GetBundleTipRequest{})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetBundleTipResponse](stream, ""), nil
}

// GetRaydiumCLMMQuotes returns the CLMM quotes on Raydium
func (g *GRPCClient) GetRaydiumCLMMQuotes(ctx context.Context, request *pb.GetRaydiumCLMMQuotesRequest) (*pb.GetRaydiumCLMMQuotesResponse, error) {
	quotes, err := g.apiClient.GetRaydiumCLMMQuotes(ctx, request)
	if err != nil {
		return nil, err
	}

	return quotes, nil
}

// GetRaydiumCLMMPools returns the CLMM pools on Raydium
func (g *GRPCClient) GetRaydiumCLMMPools(ctx context.Context, request *pb.GetRaydiumCLMMPoolsRequest) (*pb.GetRaydiumCLMMPoolsResponse, error) {
	pools, err := g.apiClient.GetRaydiumCLMMPools(ctx, request)
	if err != nil {
		return nil, err
	}

	return pools, nil
}
