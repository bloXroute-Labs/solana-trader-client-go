package provider

import (
	"context"
	"fmt"

	"github.com/bloXroute-Labs/solana-trader-client-go/connections"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/bloXroute-Labs/solana-trader-proto/common"
	"github.com/gagliardetto/solana-go"
)

type WSClientTraderAPISubmitOnly interface {
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

type WSClientTraderAPI interface {
	GetServerTime(ctx context.Context, request *pb.GetServerTimeRequest) (*pb.GetServerTimeResponse, error)
	PostSubmitPaladinV2(ctx context.Context, request *pb.PostSubmitPaladinRequest) (*pb.PostSubmitResponse, error)
	PostPumpFunSwapSol(ctx context.Context, request *pb.PostPumpFunSwapRequestSol) (*pb.PostPumpFunSwapResponse, error)
	GetRaydiumCPMMQuotes(ctx context.Context, request *pb.GetRaydiumCPMMQuotesRequest) (*pb.GetRaydiumCPMMQuotesResponse, error)
	PostRaydiumCPMMSwap(ctx context.Context, request *pb.PostRaydiumCPMMSwapRequest) (*pb.PostRaydiumCPMMSwapResponse, error)
	RecentBlockHash(ctx context.Context) (*pb.GetRecentBlockHashResponse, error)
	GetTransaction(ctx context.Context, request *pb.GetTransactionRequest) (*pb.GetTransactionResponse, error)
	GetRateLimit(ctx context.Context, request *pb.GetRateLimitRequest) (*pb.GetRateLimitResponse, error)
	GetRaydiumPoolReserve(ctx context.Context, req *pb.GetRaydiumPoolReserveRequest) (*pb.GetRaydiumPoolReserveResponse, error)
	GetRaydiumPools(ctx context.Context, request *pb.GetRaydiumPoolsRequest) (*pb.GetRaydiumPoolsResponse, error)
	GetRaydiumQuotes(ctx context.Context, request *pb.GetRaydiumQuotesRequest) (*pb.GetRaydiumQuotesResponse, error)
	GetRaydiumQuotesCPMM(ctx context.Context, request *pb.GetRaydiumCPMMQuotesRequest) (*pb.GetRaydiumCPMMQuotesResponse, error)
	GetPumpFunQuotes(ctx context.Context, request *pb.GetPumpFunQuotesRequest) (*pb.GetPumpFunQuotesResponse, error)
	GetRaydiumPrices(ctx context.Context, request *pb.GetRaydiumPricesRequest) (*pb.GetRaydiumPricesResponse, error)
	GetRaydiumCLMMQuotes(ctx context.Context, request *pb.GetRaydiumCLMMQuotesRequest) (*pb.GetRaydiumCLMMQuotesResponse, error)
	GetRaydiumCLMMPools(ctx context.Context, request *pb.GetRaydiumCLMMPoolsRequest) (*pb.GetRaydiumCLMMPoolsResponse, error)
	PostRaydiumCLMMSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest) (*pb.PostRaydiumSwapResponse, error)
	PostRaydiumCLMMRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest) (*pb.PostRaydiumRouteSwapResponse, error)
	PostRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest) (*pb.PostRaydiumSwapResponse, error)
	PostRaydiumSwapCPMM(ctx context.Context, request *pb.PostRaydiumCPMMSwapRequest) (*pb.PostRaydiumCPMMSwapResponse, error)
	PostPumpFunSwap(ctx context.Context, request *pb.PostPumpFunSwapRequest) (*pb.PostPumpFunSwapResponse, error)
	PostRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest) (*pb.PostRaydiumRouteSwapResponse, error)
	SubmitRaydiumCLMMSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitRaydiumCLMMRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
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
	GetUnsettled(ctx context.Context, market string, ownerAddress string, project pb.Project) (*pb.GetUnsettledResponse, error)
	GetAccountBalance(ctx context.Context, owner string) (*pb.GetAccountBalanceResponse, error)
	GetTokenAccounts(ctx context.Context, req *pb.GetTokenAccountsRequest) (*pb.GetTokenAccountsResponse, error)
	GetMarkets(ctx context.Context) (*pb.GetMarketsResponse, error)
	GetPrice(ctx context.Context, tokens []string) (*pb.GetPriceResponse, error)
	GetQuotes(ctx context.Context, inToken, outToken string, inAmount, slippage float64, limit int32, projects []pb.Project) (*pb.GetQuotesResponse, error)
	GetPriorityFee(ctx context.Context, project pb.Project, percentile *float64) (*pb.GetPriorityFeeResponse, error)
	GetPriorityFeeByProgram(ctx context.Context, programs []string) (*pb.GetPriorityFeeByProgramResponse, error)
	GetLeaderSchedule(ctx context.Context, maxSlots uint64) (*pb.GetLeaderScheduleResponse, error)
	PostTradeSwap(ctx context.Context, ownerAddress, inToken, outToken string, inAmount, slippage float64, projectStr string) (*pb.TradeSwapResponse, error)
	PostTradeSwapWithPriorityFee(ctx context.Context, ownerAddress, inToken, outToken string, inAmount,
		slippage float64, computeLimit uint32, computePrice uint64, projectStr string) (*pb.TradeSwapResponse, error)
	PostRouteTradeSwap(ctx context.Context, request *pb.RouteTradeSwapRequest) (*pb.TradeSwapResponse, error)
	PostOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error)
	PostSubmit(ctx context.Context, txBase64 string, opts PostSubmitOpts) (*pb.PostSubmitResponse, error)
	PostSubmitSnipeV2(ctx context.Context, request *pb.PostSubmitSnipeRequest) (*pb.PostSubmitSnipeResponse, error)
	PostSubmitBatch(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error)
	PostSubmitV2(ctx context.Context, txBase64 string, opts PostSubmitOpts) (*pb.PostSubmitResponse, error)
	PostSubmitBatchV2(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error)
	SignAndSubmit(ctx context.Context, tx *pb.TransactionMessage,
		skipPreFlight bool, frontRunningProtection bool, useStakedRPCs bool) (string, error)
	SignAndSubmitSnipe(ctx context.Context, transactions []*pb.TransactionMessage, useStakedRPCs bool) ([]string, error)
	SignAndSubmitPaladin(ctx context.Context, tx *pb.TransactionMessage, revertProtection *bool) (string, error)
	SignAndSubmitBatch(ctx context.Context, transactions []*pb.TransactionMessage, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitTradeSwap(ctx context.Context, owner, inToken, outToken string, inAmount, slippage float64, project string, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitTradeSwapWithPriorityFee(ctx context.Context, owner, inToken, outToken string,
		inAmount, slippage float64, project string, computeLimit uint32, computePrice uint64, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitRouteTradeSwap(ctx context.Context, request *pb.RouteTradeSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitRaydiumSwapCPMM(ctx context.Context, request *pb.PostRaydiumCPMMSwapRequest) (string, error)
	SubmitPostPumpFunSwap(ctx context.Context, request *pb.PostPumpFunSwapRequest) (string, error)
	SubmitRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitJupiterSwap(ctx context.Context, request *pb.PostJupiterSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitJupiterSwapInstructions(ctx context.Context, request *pb.PostJupiterSwapInstructionsRequest, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitRaydiumSwapInstructions(ctx context.Context, request *pb.PostRaydiumSwapInstructionsRequest, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitOrder(ctx context.Context, owner, payer, market string, side pb.Side,
		types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error)
	PostCancelOrder(ctx context.Context, request *pb.PostCancelOrderRequest) (*pb.PostCancelOrderResponse, error)
	SubmitCancelOrder(ctx context.Context, request *pb.PostCancelOrderRequest, skipPreFlight bool) (string, error)
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
	PostCancelAll(
		ctx context.Context,
		market,
		owner string,
		openOrdersAddresses []string,
		project pb.Project,
	) (*pb.PostCancelAllResponse, error)
	SubmitCancelAll(ctx context.Context, market, owner string, openOrdersAddresses []string, project pb.Project, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	PostSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, project pb.Project) (*pb.PostSettleResponse, error)
	SubmitSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, project pb.Project, skipPreflight bool) (string, error)
	PostReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error)
	SubmitReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error)
	PostReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error)
	SubmitReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error)
	Close() error
	GetOrderbooksStream(ctx context.Context, markets []string, limit uint32, project pb.Project) (connections.Streamer[*pb.GetOrderbooksStreamResponse], error)
	GetPumpFunSwapsStream(ctx context.Context, req *pb.GetPumpFunSwapsStreamRequest) (connections.Streamer[*pb.GetPumpFunSwapsStreamResponse], error)
	GetPumpFunAmmSwapStream(ctx context.Context, req *pb.GetPumpFunAMMSwapStreamRequest) (connections.Streamer[*pb.GetPumpFunAMMSwapStreamResponse], error)
	GetPumpFunNewTokensStream(ctx context.Context, req *pb.GetPumpFunNewTokensStreamRequest) (connections.Streamer[*pb.GetPumpFunNewTokensStreamResponse], error)
	GetMarketDepthsStream(ctx context.Context, markets []string, limit uint32, project pb.Project) (connections.Streamer[*pb.GetMarketDepthsStreamResponse], error)
	GetTradesStream(ctx context.Context, market string, limit uint32, project pb.Project) (connections.Streamer[*pb.GetTradesStreamResponse], error)
	GetNewRaydiumPoolsStream(ctx context.Context, includeCPMM bool) (connections.Streamer[*pb.GetNewRaydiumPoolsResponse], error)
	GetNewRaydiumPoolsByTransactionStream(ctx context.Context) (connections.Streamer[*pb.GetNewRaydiumPoolsByTransactionResponse], error)
	GetOrderStatusStream(ctx context.Context, market, ownerAddress string, project pb.Project) (connections.Streamer[*pb.GetOrderStatusStreamResponse], error)
	GetRecentBlockHashStream(ctx context.Context) (connections.Streamer[*pb.GetRecentBlockHashResponse], error)
	GetQuotesStream(ctx context.Context, projects []pb.Project, tokenPairs []*pb.TokenPair) (connections.Streamer[*pb.GetQuotesStreamResponse], error)
	GetPoolReservesStream(ctx context.Context, request *pb.GetPoolReservesStreamRequest) (connections.Streamer[*pb.GetPoolReservesStreamResponse], error)
	GetPricesStream(ctx context.Context, projects []pb.Project, tokens []string) (connections.Streamer[*pb.GetPricesStreamResponse], error)
	GetTickersStream(ctx context.Context, request *pb.GetTickersStreamRequest) (connections.Streamer[*pb.GetTickersStreamResponse], error)
	GetSwapsStream(
		ctx context.Context,
		projects []pb.Project,
		markets []string,
		includeFailed bool,
	) (connections.Streamer[*pb.GetSwapsStreamResponse], error)
	GetBlockStream(ctx context.Context) (connections.Streamer[*pb.GetBlockStreamResponse], error)
	GetPumpFunNewAmmPoolStream(ctx context.Context, req *pb.GetPumpFunNewAmmPoolStreamRequest) (connections.Streamer[*pb.GetPumpFunNewAmmPoolStreamResponse], error)
	GetPriorityFeeStream(ctx context.Context, project pb.Project, percentile *float64) (connections.Streamer[*pb.GetPriorityFeeResponse], error)
	GetPriorityFeeByProgramStream(ctx context.Context, programs []string) (connections.Streamer[*pb.GetPriorityFeeByProgramResponse], error)
	GetBundleTipStream(ctx context.Context) (connections.Streamer[*pb.GetBundleTipResponse], error)
	GetMarketsV2(ctx context.Context) (*pb.GetMarketsResponse, error)
	GetOrderbookV2(ctx context.Context, market string, limit uint32) (*pb.GetOrderbookResponseV2, error)
	GetMarketDepthV2(ctx context.Context, market string, limit uint32) (*pb.GetMarketDepthResponseV2, error)
	GetTickersV2(ctx context.Context, market string) (*pb.GetTickersResponseV2, error)
	GetOpenOrdersV2(ctx context.Context, market string, owner string, openOrdersAddress string, orderID string, clientOrderID uint64) (*pb.GetOpenOrdersResponse, error)
	GetUnsettledV2(ctx context.Context, market string, ownerAddress string) (*pb.GetUnsettledResponse, error)
	PostOrderV2(ctx context.Context, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error)
	SubmitOrderV2(ctx context.Context, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (string, error)
	PostCancelOrderV2(ctx context.Context, request *pb.PostCancelOrderRequestV2) (*pb.PostCancelOrderResponseV2, error)
	SubmitCancelOrderV2(ctx context.Context, request *pb.PostCancelOrderRequestV2, skipPreFlight bool) (*pb.PostSubmitBatchResponse, error)
	PostSettleV2(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string) (*pb.PostSettleResponse, error)
	SubmitSettleV2(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, skipPreflight bool) (string, error)
	PostReplaceOrderV2(ctx context.Context, orderID, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error)
	SubmitReplaceOrderV2(ctx context.Context, orderID, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (string, error)
	GetRecentBlockHash(ctx context.Context, request *pb.GetRecentBlockHashRequest) (*pb.GetRecentBlockHashResponse, error)
	GetRecentBlockHashV2(ctx context.Context, request *pb.GetRecentBlockHashRequestV2) (*pb.GetRecentBlockHashResponseV2, error)
}

type WSClient struct {
	pb.UnimplementedApiServer

	addr                 string
	conn                 *connections.WS
	privateKey           *solana.PrivateKey
	recentBlockHashStore *recentBlockHashStore
}

// NewWSClientFullService connects to Mainnet Trader API with full service
func NewWSClientFullService(bloxrouteEndpoint string) (WSClientTraderAPI, error) {
	if bloxrouteEndpoint != MainnetNYWS && bloxrouteEndpoint != MainnetUKWS {
		return nil, fmt.Errorf("not a valid endpoint for a full service trader api")
	}

	opts := DefaultRPCOpts(MainnetNYWS)
	return NewWSClientWithOpts(opts)
}

// NewWSClientSubmitOnly connects to Mainnet Trader API with submission only service
func NewWSClientSubmitOnly(bloxrouteEndpoint string) (WSClientTraderAPISubmitOnly, error) {
	if bloxrouteEndpoint != MainnetAmsterdamWS &&
		bloxrouteEndpoint != MainnetFrankfurtWS &&
		bloxrouteEndpoint != MainnetLAWS &&
		bloxrouteEndpoint != MainnetTokyoWS {
		return nil, fmt.Errorf("not a valid endpoint for submit only trader api")
	}

	opts := DefaultRPCOpts(bloxrouteEndpoint)
	opts.UseTLS = true

	return NewWSClientWithOpts(opts)
}

// NewWSClientPumpNY connects to Mainnet NY Pump Trader API
func NewWSClientPumpNY() (WSClientTraderAPI, error) {
	opts := DefaultRPCOpts(MainnetPumpNYWS)
	return NewWSClientWithOpts(opts)
}

// NewWSClientTestnet connects to Testnet Trader API
func NewWSClientTestnet() (WSClientTraderAPI, error) {
	opts := DefaultRPCOpts(TestnetWS)
	opts.UseTLS = true
	return NewWSClientWithOpts(opts)
}

// NewWSClientDevnet connects to Devnet Trader API
func NewWSClientDevnet() (WSClientTraderAPI, error) {
	opts := DefaultRPCOpts(DevnetWS)
	return NewWSClientWithOpts(opts)
}

// NewWSClientLocal connects to local Trader API
func NewWSClientLocal() (WSClientTraderAPI, error) {
	opts := DefaultRPCOpts(LocalWS)
	return NewWSClientWithOpts(opts)
}

// NewWSClientWithOpts connects to custom Trader API
func NewWSClientWithOpts(opts RPCOpts) (*WSClient, error) {
	conn, err := connections.NewWS(opts.Endpoint, opts.AuthHeader, opts.DisablePingLoop)
	if err != nil {
		return nil, err
	}

	client := &WSClient{
		addr:       opts.Endpoint,
		conn:       conn,
		privateKey: opts.PrivateKey,
	}
	client.recentBlockHashStore = newRecentBlockHashStore(
		func(ctx context.Context) (*pb.GetRecentBlockHashResponse, error) {
			return client.GetRecentBlockHash(ctx, &pb.GetRecentBlockHashRequest{})
		},
		client.GetRecentBlockHashStream,
		opts,
	)
	if opts.CacheBlockHash {
		go client.recentBlockHashStore.run(context.Background())
	}
	return client, nil
}
