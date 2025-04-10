package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/gagliardetto/solana-go"
)

type HTTPClientTraderAPI interface {
	GetRaydiumCLMMQuotes(ctx context.Context, request *pb.GetRaydiumCLMMQuotesRequest) (*pb.GetRaydiumCLMMQuotesResponse, error)
	GetRaydiumCLMMPools(ctx context.Context, request *pb.GetRaydiumCLMMPoolsRequest) (*pb.GetRaydiumCLMMPoolsResponse, error)
	PostRaydiumCLMMSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest) (*pb.PostRaydiumSwapResponse, error)
	PostRaydiumCLMMRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest) (*pb.PostRaydiumRouteSwapResponse, error)
	GetTransaction(ctx context.Context, request *pb.GetTransactionRequest) (*pb.GetTransactionResponse, error)
	GetRateLimit(ctx context.Context, _ *pb.GetRateLimitRequest) (*pb.GetRateLimitResponse, error)
	GetRaydiumPoolReserve(ctx context.Context, req *pb.GetRaydiumPoolReserveRequest) (*pb.GetRaydiumPoolReserveResponse, error)
	GetRaydiumPools(ctx context.Context, _ *pb.GetRaydiumPoolsRequest) (*pb.GetRaydiumPoolsResponse, error)
	GetRaydiumQuotes(ctx context.Context, request *pb.GetRaydiumQuotesRequest) (*pb.GetRaydiumQuotesResponse, error)
	GetRaydiumQuotesCPMM(ctx context.Context, request *pb.GetRaydiumCPMMQuotesRequest) (*pb.GetRaydiumCPMMQuotesResponse, error)
	GetPumpFunQuotes(ctx context.Context, request *pb.GetPumpFunQuotesRequest) (*pb.GetPumpFunQuotesResponse, error)
	GetRaydiumPrices(ctx context.Context, request *pb.GetRaydiumPricesRequest) (*pb.GetRaydiumPricesResponse, error)
	SubmitRaydiumCLMMSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitRaydiumCLMMRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	PostRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest) (*pb.PostRaydiumSwapResponse, error)
	PostRaydiumCPMMSwap(ctx context.Context, request *pb.PostRaydiumCPMMSwapRequest) (*pb.PostRaydiumCPMMSwapResponse, error)
	PostPumpFunSwap(ctx context.Context, request *pb.PostPumpFunSwapRequest) (*pb.PostPumpFunSwapResponse, error)
	PostRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest) (*pb.PostRaydiumRouteSwapResponse, error)
	GetJupiterQuotes(ctx context.Context, request *pb.GetJupiterQuotesRequest) (*pb.GetJupiterQuotesResponse, error)
	GetJupiterPrices(ctx context.Context, request *pb.GetJupiterPricesRequest) (*pb.GetJupiterPricesResponse, error)
	PostJupiterSwap(ctx context.Context, request *pb.PostJupiterSwapRequest) (*pb.PostJupiterSwapResponse, error)
	PostJupiterSwapInstructions(ctx context.Context, request *pb.PostJupiterSwapInstructionsRequest) (*pb.PostJupiterSwapInstructionsResponse, error)
	PostRaydiumSwapInstructions(ctx context.Context, request *pb.PostRaydiumSwapInstructionsRequest) (*pb.PostRaydiumSwapInstructionsResponse, error)
	PostJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest) (*pb.PostJupiterRouteSwapResponse, error)
	GetPools(ctx context.Context, projects []pb.Project) (*pb.GetPoolsResponse, error)
	GetTokenAccounts(ctx context.Context, req *pb.GetTokenAccountsRequest) (*pb.GetTokenAccountsResponse, error)
	PostSubmit(ctx context.Context, txBase64 string, opts PostSubmitOpts) (*pb.PostSubmitResponse, error)
	PostSubmitSnipeV2(ctx context.Context, request *pb.PostSubmitSnipeRequest) (*pb.PostSubmitSnipeResponse, error)
	PostSubmitBatch(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error)
	PostSubmitV2(ctx context.Context, txBase64 string, opts PostSubmitOpts) (*pb.PostSubmitResponse, error)
	PostSubmitBatchV2(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error)
	SignAndSubmit(ctx context.Context, tx *pb.TransactionMessage,
		skipPreFlight bool, frontRunningProtection bool, useStakedRPCs bool) (string, error)
	SignAndSubmitSnipe(ctx context.Context, transactions []*pb.TransactionMessage, useStakedRPCs bool) ([]string, error)
	SignAndSubmitPaladin(ctx context.Context, tx *pb.TransactionMessage, revertProtection *bool) (string, error)
	SignAndSubmitBatch(ctx context.Context, transactions []*pb.TransactionMessage, useBundle bool,
		opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitRaydiumSwapCPMM(ctx context.Context, request *pb.PostRaydiumCPMMSwapRequest) (string, error)
	SubmitPostPumpFunSwap(ctx context.Context, request *pb.PostPumpFunSwapRequest) (string, error)
	SubmitRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitJupiterSwap(ctx context.Context, request *pb.PostJupiterSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitJupiterSwapInstructions(ctx context.Context, request *pb.PostJupiterSwapInstructionsRequest, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitRaydiumSwapInstructions(ctx context.Context, request *pb.PostRaydiumSwapInstructionsRequest, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	GetRecentBlockHash(ctx context.Context) (*pb.GetRecentBlockHashResponse, error)
	GetRecentBlockHashV2(ctx context.Context, offset uint64) (*pb.GetRecentBlockHashResponseV2, error)
	GetPriorityFee(ctx context.Context, project pb.Project, percentile *float64) (*pb.GetPriorityFeeResponse, error)
	GetPriorityFeeByProgram(ctx context.Context, programs []string) (*pb.GetPriorityFeeByProgramResponse, error)
	GetLeaderSchedule(ctx context.Context, maxSlots uint) (*pb.GetLeaderScheduleResponse, error)
}

type HTTPClientTraderAPISubmitOnly interface {
	PostSubmit(ctx context.Context, txBase64 string, opts PostSubmitOpts) (*pb.PostSubmitResponse, error)
	PostSubmitSnipeV2(ctx context.Context, request *pb.PostSubmitSnipeRequest) (*pb.PostSubmitSnipeResponse, error)
	PostSubmitBatch(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error)
	PostSubmitV2(ctx context.Context, txBase64 string, opts PostSubmitOpts) (*pb.PostSubmitResponse, error)
	PostSubmitBatchV2(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error)
	SignAndSubmit(ctx context.Context, tx *pb.TransactionMessage,
		skipPreFlight bool, frontRunningProtection bool, useStakedRPCs bool) (string, error)
	SignAndSubmitSnipe(ctx context.Context, transactions []*pb.TransactionMessage, useStakedRPCs bool) ([]string, error)
	SignAndSubmitPaladin(ctx context.Context, tx *pb.TransactionMessage, revertProtection *bool) (string, error)
	SignAndSubmitBatch(ctx context.Context, transactions []*pb.TransactionMessage, useBundle bool,
		opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
}

type HTTPClient struct {
	pb.UnimplementedApiServer

	baseURL    string
	httpClient *http.Client
	requestID  utils.RequestID
	privateKey *solana.PrivateKey
	authHeader string
}

// NewHTTPClientFullService connects to Mainnet Trader pb with full service pb
func NewHTTPClientFullService(bloxrouteEndpoint string) (HTTPClientTraderAPI, error) {
	if bloxrouteEndpoint != MainnetNYHTTP && bloxrouteEndpoint != MainnetUKHTTP {
		return nil, fmt.Errorf("not a valid endpoint for a full service trader pb")
	}

	opts := DefaultRPCOpts(bloxrouteEndpoint)
	opts.UseTLS = true

	httpClient := NewHTTPClientWithOpts(nil, opts)
	return httpClient, nil
}

// NewHTTPClientSubmitOnly connects to Mainnet Trader pb with full service pb
func NewHTTPClientSubmitOnly(bloxrouteEndpoint string) (HTTPClientTraderAPISubmitOnly, error) {
	if bloxrouteEndpoint != MainnetAmsterdamHTTP &&
		bloxrouteEndpoint != MainnetFrankfurtHTTP &&
		bloxrouteEndpoint != MainnetLAHTTP &&
		bloxrouteEndpoint != MainnetTokyoHTTP {
		return nil, fmt.Errorf("not a valid endpoint for submit only trader pb")
	}

	opts := DefaultRPCOpts(bloxrouteEndpoint)
	opts.UseTLS = true

	httpClient := NewHTTPClientWithOpts(nil, opts)
	return httpClient, nil
}

// NewHTTPClient connects to Mainnet Trader pb
func NewHTTPClient() HTTPClientTraderAPI {
	opts := DefaultRPCOpts(MainnetNYHTTP)
	return NewHTTPClientWithOpts(nil, opts)
}

// NewHTTPClientPumpNY connects to Mainnet Trader pb
func NewHTTPClientPumpNY() HTTPClientTraderAPI {
	opts := DefaultRPCOpts(MainnetPumpNYHTTP)
	return NewHTTPClientWithOpts(nil, opts)
}

// NewHTTPTestnet connects to Testnet Trader pb
func NewHTTPTestnet() HTTPClientTraderAPI {
	opts := DefaultRPCOpts(TestnetHTTP)
	opts.UseTLS = true
	return NewHTTPClientWithOpts(nil, opts)
}

// NewHTTPDevnet connects to Devnet Trader pb
func NewHTTPDevnet() HTTPClientTraderAPI {
	opts := DefaultRPCOpts(DevnetHTTP)
	return NewHTTPClientWithOpts(nil, opts)
}

// NewHTTPLocal connects to local Trader pb
func NewHTTPLocal() HTTPClientTraderAPI {
	opts := DefaultRPCOpts(LocalHTTP)
	return NewHTTPClientWithOpts(nil, opts)
}

// NewHTTPClientWithOpts connects to custom Trader pb (set client to nil to use default client)
func NewHTTPClientWithOpts(client *http.Client, opts RPCOpts) *HTTPClient {
	if client == nil {
		client = &http.Client{}
	}

	return &HTTPClient{
		baseURL:    opts.Endpoint,
		httpClient: client,
		privateKey: opts.PrivateKey,
		authHeader: opts.AuthHeader,
	}
}
