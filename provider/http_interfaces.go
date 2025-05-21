package provider

import (
	"context"
	"fmt"
	"net/http"

	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/bloXroute-Labs/solana-trader-proto/common"
	"github.com/gagliardetto/solana-go"
)

type HTTPClientTraderAPI interface {
	GetServerTime(ctx context.Context, request *pb.GetServerTimeRequest) (*pb.GetServerTimeResponse, error)
	PostSubmitPaladinV2(ctx context.Context, request *pb.PostSubmitPaladinRequest) (*pb.PostSubmitResponse, error)
	PostPumpFunSwapSol(ctx context.Context, request *pb.PostPumpFunSwapRequestSol) (*pb.PostPumpFunSwapResponse, error)
	GetRaydiumCPMMQuotes(ctx context.Context, request *pb.GetRaydiumCPMMQuotesRequest) (*pb.GetRaydiumCPMMQuotesResponse, error)
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
	GetOrderbook(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetOrderbookResponse, error)
	GetMarketDepth(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetMarketDepthResponse, error)
	GetTrades(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetTradesResponse, error)
	GetPools(ctx context.Context, projects []pb.Project) (*pb.GetPoolsResponse, error)
	GetTickers(ctx context.Context, market string, project pb.Project) (*pb.GetTickersResponse, error)
	GetOpenOrders(ctx context.Context, market string, owner string, openOrdersAddress string, project pb.Project) (*pb.GetOpenOrdersResponse, error)
	GetOrderByID(ctx context.Context, in *pb.GetOrderByIDRequest) (*pb.GetOrderByIDResponse, error)
	GetMarkets(ctx context.Context) (*pb.GetMarketsResponse, error)
	GetUnsettled(ctx context.Context, market string, owner string, project pb.Project) (*pb.GetUnsettledResponse, error)
	GetAccountBalance(ctx context.Context, owner string) (*pb.GetAccountBalanceResponse, error)
	GetTokenAccounts(ctx context.Context, req *pb.GetTokenAccountsRequest) (*pb.GetTokenAccountsResponse, error)
	GetPrice(ctx context.Context, tokens []string) (*pb.GetPriceResponse, error)
	GetQuotes(ctx context.Context, inToken, outToken string, inAmount, slippage float64, limit int32, projects []pb.Project) (*pb.GetQuotesResponse, error)
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
	PostTradeSwap(ctx context.Context, ownerAddress, inToken, outToken string, inAmount, slippage float64, project pb.Project) (*pb.TradeSwapResponse, error)
	SubmitTradeSwap(ctx context.Context, owner, inToken, outToken string, inAmount, slippage float64, project pb.Project, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	PostRouteTradeSwap(ctx context.Context, request *pb.RouteTradeSwapRequest) (*pb.TradeSwapResponse, error)
	SubmitRouteTradeSwap(ctx context.Context, request *pb.RouteTradeSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitRaydiumSwapCPMM(ctx context.Context, request *pb.PostRaydiumCPMMSwapRequest) (string, error)
	SubmitPostPumpFunSwap(ctx context.Context, request *pb.PostPumpFunSwapRequest) (string, error)
	SubmitRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitJupiterSwap(ctx context.Context, request *pb.PostJupiterSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitJupiterSwapInstructions(ctx context.Context, request *pb.PostJupiterSwapInstructionsRequest, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitRaydiumSwapInstructions(ctx context.Context, request *pb.PostRaydiumSwapInstructionsRequest, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	PostOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error)
	SubmitOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error)
	PostCancelOrder(
		ctx context.Context,
		orderID string,
		side pb.Side,
		owner,
		market,
		openOrders string,
		project pb.Project,
	) (*pb.PostCancelOrderResponse, error)
	SubmitCancelOrder(
		ctx context.Context,
		orderID string,
		side pb.Side,
		owner,
		market,
		openOrders string,
		project pb.Project,
		skipPreFlight bool,
	) (string, error)
	PostCancelByClientOrderID(
		ctx context.Context,
		clientOrderID uint64,
		owner,
		market,
		openOrders string,
		project pb.Project,
	) (*pb.PostCancelOrderResponse, error)
	SubmitCancelByClientOrderID(
		ctx context.Context,
		clientOrderID uint64,
		owner,
		market,
		openOrders string,
		project pb.Project,
		skipPreFlight bool,
	) (string, error)
	PostCancelAll(ctx context.Context, market, owner string, openOrdersAddresses []string, project pb.Project) (*pb.PostCancelAllResponse, error)
	SubmitCancelAll(ctx context.Context, market, owner string, openOrders []string, project pb.Project, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	PostSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, project pb.Project) (*pb.PostSettleResponse, error)
	SubmitSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, project pb.Project, skipPreflight bool) (string, error)
	PostReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error)
	SubmitReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error)
	PostReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error)
	SubmitReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error)
	GetRecentBlockHash(ctx context.Context) (*pb.GetRecentBlockHashResponse, error)
	GetRecentBlockHashV2(ctx context.Context, offset uint64) (*pb.GetRecentBlockHashResponseV2, error)
	GetPriorityFee(ctx context.Context, project pb.Project, percentile *float64) (*pb.GetPriorityFeeResponse, error)
	GetPriorityFeeByProgram(ctx context.Context, programs []string) (*pb.GetPriorityFeeByProgramResponse, error)
	GetMarketsV2(ctx context.Context) (*pb.GetMarketsResponseV2, error)
	GetOrderbookV2(ctx context.Context, market string, limit uint32) (*pb.GetOrderbookResponseV2, error)
	GetMarketDepthV2(ctx context.Context, market string, limit uint32) (*pb.GetMarketDepthResponseV2, error)
	GetTickersV2(ctx context.Context, market string) (*pb.GetTickersResponseV2, error)
	GetOpenOrdersV2(ctx context.Context, market string, owner string, openOrdersAddress string, orderID string, clientOrderID uint64) (*pb.GetOpenOrdersResponse, error)
	GetUnsettledV2(ctx context.Context, market string, owner string) (*pb.GetUnsettledResponse, error)
	PostOrderV2(ctx context.Context, owner, payer, market string, side string, orderType string, amount,
		price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error)
	PostOrderV2WithPriorityFee(ctx context.Context, owner, payer, market string, side string,
		orderType string, amount, price float64, computeLimit uint32, computePrice uint64, opts PostOrderOpts) (*pb.PostOrderResponse, error)
	SubmitOrderV2(ctx context.Context, owner, payer, market string, side string, orderType string,
		amount, price float64, opts PostOrderOpts) (string, error)
	SubmitOrderV2WithPriorityFee(ctx context.Context, owner, payer, market string, side string,
		orderType string, amount, price float64, computeLimit uint32, computePrice uint64, opts PostOrderOpts) (string, error)
	PostCancelOrderV2(
		ctx context.Context,
		orderID string,
		clientOrderID uint64,
		side string,
		owner,
		market,
		openOrders string,
	) (*pb.PostCancelOrderResponseV2, error)
	SubmitCancelOrderV2(
		ctx context.Context,
		orderID string,
		clientOrderID uint64,
		side string,
		owner,
		market,
		openOrders string,
		opts SubmitOpts,
	) (*pb.PostSubmitBatchResponse, error)
	PostSettleV2(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string) (*pb.PostSettleResponse, error)
	SubmitSettleV2(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, skipPreflight bool) (string, error)
	PostReplaceOrderV2(ctx context.Context, orderID, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error)
	SubmitReplaceOrderV2(ctx context.Context, orderID, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (string, error)
}

type HTTPClientTraderAPISubmitOnly interface {
	PostSubmit(ctx context.Context, txBase64 string, opts PostSubmitOpts) (*pb.PostSubmitResponse, error)
	PostSubmitSnipeV2(ctx context.Context, request *pb.PostSubmitSnipeRequest) (*pb.PostSubmitSnipeResponse, error)
	PostSubmitBatch(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error)
	PostSubmitV2(ctx context.Context, txBase64 string, opts PostSubmitOpts) (*pb.PostSubmitResponse, error)
	PostSubmitBatchV2(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error)
	SignAndSubmit(ctx context.Context, tx *pb.TransactionMessage, skipPreFlight bool, frontRunningProtection bool, useStakedRPCs bool) (string, error)
	SignAndSubmitSnipe(ctx context.Context, transactions []*pb.TransactionMessage, useStakedRPCs bool) ([]string, error)
	SignAndSubmitPaladin(ctx context.Context, tx *pb.TransactionMessage, revertProtection *bool) (string, error)
	SignAndSubmitBatch(ctx context.Context, transactions []*pb.TransactionMessage, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
}

type HTTPClient struct {
	pb.UnimplementedApiServer

	baseURL    string
	httpClient *http.Client
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
