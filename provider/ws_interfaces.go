package provider

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/connections"
	"github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/bloXroute-Labs/solana-trader-proto/common"
	"github.com/gagliardetto/solana-go"
)

type WSClientTraderAPISubmitOnly interface {
	PostSubmit(ctx context.Context, txBase64 string, opts PostSubmitOpts) (*api.PostSubmitResponse, error)
	PostSubmitSnipeV2(ctx context.Context, request *api.PostSubmitSnipeRequest) (*api.PostSubmitSnipeResponse, error)
	PostSubmitBatch(ctx context.Context, request *api.PostSubmitBatchRequest) (*api.PostSubmitBatchResponse, error)
	PostSubmitV2(ctx context.Context, txBase64 string, opts PostSubmitOpts) (*api.PostSubmitResponse, error)
	PostSubmitBatchV2(ctx context.Context, request *api.PostSubmitBatchRequest) (*api.PostSubmitBatchResponse, error)
	SignAndSubmit(ctx context.Context, tx *api.TransactionMessage,
		skipPreFlight bool, frontRunningProtection bool, useStakedRPCs bool) (string, error)
	SignAndSubmitSnipe(ctx context.Context, transactions []*api.TransactionMessage, useStakedRPCs bool) ([]string, error)
	SignAndSubmitPaladin(ctx context.Context, tx *api.TransactionMessage, revertProtection *bool) (string, error)
	SignAndSubmitBatch(ctx context.Context, transactions []*api.TransactionMessage, useBundle bool, opts SubmitOpts) (*api.PostSubmitBatchResponse, error)
}

type WSClientTraderAPI interface {
	RecentBlockHash(ctx context.Context) (*api.GetRecentBlockHashResponse, error)
	GetTransaction(ctx context.Context, request *api.GetTransactionRequest) (*api.GetTransactionResponse, error)
	GetRateLimit(ctx context.Context, request *api.GetRateLimitRequest) (*api.GetRateLimitResponse, error)
	GetRaydiumPoolReserve(ctx context.Context, req *api.GetRaydiumPoolReserveRequest) (*api.GetRaydiumPoolReserveResponse, error)
	GetRaydiumPools(ctx context.Context, request *api.GetRaydiumPoolsRequest) (*api.GetRaydiumPoolsResponse, error)
	GetRaydiumQuotes(ctx context.Context, request *api.GetRaydiumQuotesRequest) (*api.GetRaydiumQuotesResponse, error)
	GetRaydiumQuotesCPMM(ctx context.Context, request *api.GetRaydiumCPMMQuotesRequest) (*api.GetRaydiumCPMMQuotesResponse, error)
	GetPumpFunQuotes(ctx context.Context, request *api.GetPumpFunQuotesRequest) (*api.GetPumpFunQuotesResponse, error)
	GetRaydiumPrices(ctx context.Context, request *api.GetRaydiumPricesRequest) (*api.GetRaydiumPricesResponse, error)
	GetRaydiumCLMMQuotes(ctx context.Context, request *api.GetRaydiumCLMMQuotesRequest) (*api.GetRaydiumCLMMQuotesResponse, error)
	GetRaydiumCLMMPools(ctx context.Context, request *api.GetRaydiumCLMMPoolsRequest) (*api.GetRaydiumCLMMPoolsResponse, error)
	PostRaydiumCLMMSwap(ctx context.Context, request *api.PostRaydiumSwapRequest) (*api.PostRaydiumSwapResponse, error)
	PostRaydiumCLMMRouteSwap(ctx context.Context, request *api.PostRaydiumRouteSwapRequest) (*api.PostRaydiumRouteSwapResponse, error)
	PostRaydiumSwap(ctx context.Context, request *api.PostRaydiumSwapRequest) (*api.PostRaydiumSwapResponse, error)
	PostRaydiumSwapCPMM(ctx context.Context, request *api.PostRaydiumCPMMSwapRequest) (*api.PostRaydiumCPMMSwapResponse, error)
	PostPumpFunSwap(ctx context.Context, request *api.PostPumpFunSwapRequest) (*api.PostPumpFunSwapResponse, error)
	PostRaydiumRouteSwap(ctx context.Context, request *api.PostRaydiumRouteSwapRequest) (*api.PostRaydiumRouteSwapResponse, error)
	SubmitRaydiumCLMMSwap(ctx context.Context, request *api.PostRaydiumSwapRequest, opts SubmitOpts) (*api.PostSubmitBatchResponse, error)
	SubmitRaydiumCLMMRouteSwap(ctx context.Context, request *api.PostRaydiumRouteSwapRequest, opts SubmitOpts) (*api.PostSubmitBatchResponse, error)
	GetJupiterQuotes(ctx context.Context, request *api.GetJupiterQuotesRequest) (*api.GetJupiterQuotesResponse, error)
	GetJupiterPrices(ctx context.Context, request *api.GetJupiterPricesRequest) (*api.GetJupiterPricesResponse, error)
	PostJupiterSwap(ctx context.Context, request *api.PostJupiterSwapRequest) (*api.PostJupiterSwapResponse, error)
	PostJupiterSwapInstructions(ctx context.Context, request *api.PostJupiterSwapInstructionsRequest) (*api.PostJupiterSwapInstructionsResponse, error)
	PostRaydiumSwapInstructions(ctx context.Context, request *api.PostRaydiumSwapInstructionsRequest) (*api.PostRaydiumSwapInstructionsResponse, error)
	PostJupiterRouteSwap(ctx context.Context, request *api.PostJupiterRouteSwapRequest) (*api.PostJupiterRouteSwapResponse, error)
	GetOrderbook(ctx context.Context, market string, limit uint32, project api.Project) (*api.GetOrderbookResponse, error)
	GetMarketDepth(ctx context.Context, market string, limit uint32, project api.Project) (*api.GetMarketDepthResponse, error)
	GetTrades(ctx context.Context, market string, limit uint32, project api.Project) (*api.GetTradesResponse, error)
	GetPools(ctx context.Context, projects []api.Project) (*api.GetPoolsResponse, error)
	GetTickers(ctx context.Context, market string, project api.Project) (*api.GetTickersResponse, error)
	GetOpenOrders(ctx context.Context, market string, owner string, openOrdersAddress string, project api.Project) (*api.GetOpenOrdersResponse, error)
	GetOrderByID(ctx context.Context, in *api.GetOrderByIDRequest) (*api.GetOrderByIDResponse, error)
	GetUnsettled(ctx context.Context, market string, ownerAddress string, project api.Project) (*api.GetUnsettledResponse, error)
	GetAccountBalance(ctx context.Context, owner string) (*api.GetAccountBalanceResponse, error)
	GetTokenAccounts(ctx context.Context, req *api.GetTokenAccountsRequest) (*api.GetTokenAccountsResponse, error)
	GetMarkets(ctx context.Context) (*api.GetMarketsResponse, error)
	GetPrice(ctx context.Context, tokens []string) (*api.GetPriceResponse, error)
	GetQuotes(ctx context.Context, inToken, outToken string, inAmount, slippage float64, limit int32, projects []api.Project) (*api.GetQuotesResponse, error)
	GetPriorityFee(ctx context.Context, project api.Project, percentile *float64) (*api.GetPriorityFeeResponse, error)
	GetPriorityFeeByProgram(ctx context.Context, programs []string) (*api.GetPriorityFeeByProgramResponse, error)
	GetLeaderSchedule(ctx context.Context, maxSlots uint64) (*api.GetLeaderScheduleResponse, error)
	PostTradeSwap(ctx context.Context, ownerAddress, inToken, outToken string, inAmount, slippage float64, projectStr string) (*api.TradeSwapResponse, error)
	PostTradeSwapWithPriorityFee(ctx context.Context, ownerAddress, inToken, outToken string, inAmount,
		slippage float64, computeLimit uint32, computePrice uint64, projectStr string) (*api.TradeSwapResponse, error)
	PostRouteTradeSwap(ctx context.Context, request *api.RouteTradeSwapRequest) (*api.TradeSwapResponse, error)
	PostOrder(ctx context.Context, owner, payer, market string, side api.Side, types []common.OrderType, amount, price float64, project api.Project, opts PostOrderOpts) (*api.PostOrderResponse, error)
	PostSubmit(ctx context.Context, txBase64 string, opts PostSubmitOpts) (*api.PostSubmitResponse, error)
	PostSubmitSnipeV2(ctx context.Context, request *api.PostSubmitSnipeRequest) (*api.PostSubmitSnipeResponse, error)
	PostSubmitBatch(ctx context.Context, request *api.PostSubmitBatchRequest) (*api.PostSubmitBatchResponse, error)
	PostSubmitV2(ctx context.Context, txBase64 string, opts PostSubmitOpts) (*api.PostSubmitResponse, error)
	PostSubmitBatchV2(ctx context.Context, request *api.PostSubmitBatchRequest) (*api.PostSubmitBatchResponse, error)
	SignAndSubmit(ctx context.Context, tx *api.TransactionMessage,
		skipPreFlight bool, frontRunningProtection bool, useStakedRPCs bool) (string, error)
	SignAndSubmitSnipe(ctx context.Context, transactions []*api.TransactionMessage, useStakedRPCs bool) ([]string, error)
	SignAndSubmitPaladin(ctx context.Context, tx *api.TransactionMessage, revertProtection *bool) (string, error)
	SignAndSubmitBatch(ctx context.Context, transactions []*api.TransactionMessage, useBundle bool, opts SubmitOpts) (*api.PostSubmitBatchResponse, error)
	SubmitTradeSwap(ctx context.Context, owner, inToken, outToken string, inAmount, slippage float64, project string, opts SubmitOpts) (*api.PostSubmitBatchResponse, error)
	SubmitTradeSwapWithPriorityFee(ctx context.Context, owner, inToken, outToken string,
		inAmount, slippage float64, project string, computeLimit uint32, computePrice uint64, opts SubmitOpts) (*api.PostSubmitBatchResponse, error)
	SubmitRouteTradeSwap(ctx context.Context, request *api.RouteTradeSwapRequest, opts SubmitOpts) (*api.PostSubmitBatchResponse, error)
	SubmitRaydiumSwap(ctx context.Context, request *api.PostRaydiumSwapRequest, opts SubmitOpts) (*api.PostSubmitBatchResponse, error)
	SubmitRaydiumSwapCPMM(ctx context.Context, request *api.PostRaydiumCPMMSwapRequest) (string, error)
	SubmitPostPumpFunSwap(ctx context.Context, request *api.PostPumpFunSwapRequest) (string, error)
	SubmitRaydiumRouteSwap(ctx context.Context, request *api.PostRaydiumRouteSwapRequest, opts SubmitOpts) (*api.PostSubmitBatchResponse, error)
	SubmitJupiterSwap(ctx context.Context, request *api.PostJupiterSwapRequest, opts SubmitOpts) (*api.PostSubmitBatchResponse, error)
	SubmitJupiterSwapInstructions(ctx context.Context, request *api.PostJupiterSwapInstructionsRequest, useBundle bool, opts SubmitOpts) (*api.PostSubmitBatchResponse, error)
	SubmitRaydiumSwapInstructions(ctx context.Context, request *api.PostRaydiumSwapInstructionsRequest, useBundle bool, opts SubmitOpts) (*api.PostSubmitBatchResponse, error)
	SubmitJupiterRouteSwap(ctx context.Context, request *api.PostJupiterRouteSwapRequest, opts SubmitOpts) (*api.PostSubmitBatchResponse, error)
	SubmitOrder(ctx context.Context, owner, payer, market string, side api.Side,
		types []common.OrderType, amount, price float64, project api.Project, opts PostOrderOpts) (string, error)
	PostCancelOrder(ctx context.Context, request *api.PostCancelOrderRequest) (*api.PostCancelOrderResponse, error)
	SubmitCancelOrder(ctx context.Context, request *api.PostCancelOrderRequest, skipPreFlight bool) (string, error)
	PostCancelByClientOrderID(
		ctx context.Context,
		clientOrderID uint64,
		owner,
		market,
		openOrders string,
		project api.Project,
	) (*api.PostCancelOrderResponse, error)
	SubmitCancelByClientOrderID(
		ctx context.Context,
		clientOrderID uint64,
		owner,
		market,
		openOrders string,
		project api.Project,
		skipPreFlight bool,
	) (string, error)
	PostCancelAll(
		ctx context.Context,
		market,
		owner string,
		openOrdersAddresses []string,
		project api.Project,
	) (*api.PostCancelAllResponse, error)
	SubmitCancelAll(ctx context.Context, market, owner string, openOrdersAddresses []string, project api.Project, opts SubmitOpts) (*api.PostSubmitBatchResponse, error)
	PostSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, project api.Project) (*api.PostSettleResponse, error)
	SubmitSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, project api.Project, skipPreflight bool) (string, error)
	PostReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side api.Side, types []common.OrderType, amount, price float64, project api.Project, opts PostOrderOpts) (*api.PostOrderResponse, error)
	SubmitReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side api.Side, types []common.OrderType, amount, price float64, project api.Project, opts PostOrderOpts) (string, error)
	PostReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side api.Side, types []common.OrderType, amount, price float64, project api.Project, opts PostOrderOpts) (*api.PostOrderResponse, error)
	SubmitReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side api.Side, types []common.OrderType, amount, price float64, project api.Project, opts PostOrderOpts) (string, error)
	Close() error
	GetOrderbooksStream(ctx context.Context, markets []string, limit uint32, project api.Project) (connections.Streamer[*api.GetOrderbooksStreamResponse], error)
	GetPumpFunSwapsStream(ctx context.Context, req *api.GetPumpFunSwapsStreamRequest) (connections.Streamer[*api.GetPumpFunSwapsStreamResponse], error)
	GetPumpFunNewTokensStream(ctx context.Context, req *api.GetPumpFunNewTokensStreamRequest) (connections.Streamer[*api.GetPumpFunNewTokensStreamResponse], error)
	GetMarketDepthsStream(ctx context.Context, markets []string, limit uint32, project api.Project) (connections.Streamer[*api.GetMarketDepthsStreamResponse], error)
	GetTradesStream(ctx context.Context, market string, limit uint32, project api.Project) (connections.Streamer[*api.GetTradesStreamResponse], error)
	GetNewRaydiumPoolsStream(ctx context.Context, includeCPMM bool) (connections.Streamer[*api.GetNewRaydiumPoolsResponse], error)
	GetNewRaydiumPoolsByTransactionStream(ctx context.Context) (connections.Streamer[*api.GetNewRaydiumPoolsByTransactionResponse], error)
	GetOrderStatusStream(ctx context.Context, market, ownerAddress string, project api.Project) (connections.Streamer[*api.GetOrderStatusStreamResponse], error)
	GetRecentBlockHashStream(ctx context.Context) (connections.Streamer[*api.GetRecentBlockHashResponse], error)
	GetQuotesStream(ctx context.Context, projects []api.Project, tokenPairs []*api.TokenPair) (connections.Streamer[*api.GetQuotesStreamResponse], error)
	GetPoolReservesStream(ctx context.Context, request *api.GetPoolReservesStreamRequest) (connections.Streamer[*api.GetPoolReservesStreamResponse], error)
	GetPricesStream(ctx context.Context, projects []api.Project, tokens []string) (connections.Streamer[*api.GetPricesStreamResponse], error)
	GetTickersStream(ctx context.Context, request *api.GetTickersStreamRequest) (connections.Streamer[*api.GetTickersStreamResponse], error)
	GetSwapsStream(
		ctx context.Context,
		projects []api.Project,
		markets []string,
		includeFailed bool,
	) (connections.Streamer[*api.GetSwapsStreamResponse], error)
	GetBlockStream(ctx context.Context) (connections.Streamer[*api.GetBlockStreamResponse], error)
	GetPriorityFeeStream(ctx context.Context, project api.Project, percentile *float64) (connections.Streamer[*api.GetPriorityFeeResponse], error)
	GetPriorityFeeByProgramStream(ctx context.Context, programs []string) (connections.Streamer[*api.GetPriorityFeeByProgramResponse], error)
	GetBundleTipStream(ctx context.Context) (connections.Streamer[*api.GetBundleTipResponse], error)
	GetMarketsV2(ctx context.Context) (*api.GetMarketsResponse, error)
	GetOrderbookV2(ctx context.Context, market string, limit uint32) (*api.GetOrderbookResponseV2, error)
	GetMarketDepthV2(ctx context.Context, market string, limit uint32) (*api.GetMarketDepthResponseV2, error)
	GetTickersV2(ctx context.Context, market string) (*api.GetTickersResponseV2, error)
	GetOpenOrdersV2(ctx context.Context, market string, owner string, openOrdersAddress string, orderID string, clientOrderID uint64) (*api.GetOpenOrdersResponse, error)
	GetUnsettledV2(ctx context.Context, market string, ownerAddress string) (*api.GetUnsettledResponse, error)
	PostOrderV2(ctx context.Context, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (*api.PostOrderResponse, error)
	SubmitOrderV2(ctx context.Context, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (string, error)
	PostCancelOrderV2(ctx context.Context, request *api.PostCancelOrderRequestV2) (*api.PostCancelOrderResponseV2, error)
	SubmitCancelOrderV2(ctx context.Context, request *api.PostCancelOrderRequestV2, skipPreFlight bool) (*api.PostSubmitBatchResponse, error)
	PostSettleV2(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string) (*api.PostSettleResponse, error)
	SubmitSettleV2(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, skipPreflight bool) (string, error)
	PostReplaceOrderV2(ctx context.Context, orderID, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (*api.PostOrderResponse, error)
	SubmitReplaceOrderV2(ctx context.Context, orderID, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (string, error)
	GetRecentBlockHash(ctx context.Context, request *api.GetRecentBlockHashRequest) (*api.GetRecentBlockHashResponse, error)
	GetRecentBlockHashV2(ctx context.Context, request *api.GetRecentBlockHashRequestV2) (*api.GetRecentBlockHashResponseV2, error)
}

type WSClient struct {
	api.UnimplementedApiServer

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
		func(ctx context.Context) (*api.GetRecentBlockHashResponse, error) {
			return client.GetRecentBlockHash(ctx, &api.GetRecentBlockHashRequest{})
		},
		client.GetRecentBlockHashStream,
		opts,
	)
	if opts.CacheBlockHash {
		go client.recentBlockHashStore.run(context.Background())
	}
	return client, nil
}
