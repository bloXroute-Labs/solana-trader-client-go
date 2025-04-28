package provider

import (
	"context"
	"crypto/tls"
	"fmt"

	package_info "github.com/bloXroute-Labs/solana-trader-client-go"
	"github.com/bloXroute-Labs/solana-trader-client-go/connections"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/bloXroute-Labs/solana-trader-proto/common"
	"github.com/gagliardetto/solana-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClientTraderAPISubmitOnly interface {
	PostSubmit(ctx context.Context, tx *pb.TransactionMessage, opts PostSubmitOpts) (*pb.PostSubmitResponse, error)
	SignAndSubmit(ctx context.Context, tx *pb.TransactionMessage, skipPreFlight bool, frontRunningProtection bool, useStakedRPCs bool) (string, error)
	SignAndSubmitSnipe(ctx context.Context, transactions []*pb.TransactionMessage, useStakedRPCs bool) ([]string, error)
	SignAndSubmitPaladin(ctx context.Context, tx *pb.TransactionMessage, revertProtection *bool) (string, error)
	PostSubmitPaladinV2(ctx context.Context, request *pb.PostSubmitPaladinRequest) (*pb.PostSubmitResponse, error)
	PostSubmitBatch(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error)
	PostSubmitV2(ctx context.Context, tx *pb.TransactionMessage, opts PostSubmitOpts) (*pb.PostSubmitResponse, error)
	PostSubmitBatchV2(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error)
	PostSubmitSnipeV2(ctx context.Context, request *pb.PostSubmitSnipeRequest) (*pb.PostSubmitSnipeResponse, error)
}

type GRPCClientTraderAPI interface {
	GetServerTime(ctx context.Context, request *pb.GetServerTimeRequest) (*pb.GetServerTimeResponse, error)
	PostSubmitPaladinV2(ctx context.Context, request *pb.PostSubmitPaladinRequest) (*pb.PostSubmitResponse, error)
	PostPumpFunSwapSol(ctx context.Context, request *pb.PostPumpFunSwapRequestSol) (*pb.PostPumpFunSwapResponse, error)
	GetRaydiumCPMMQuotes(ctx context.Context, request *pb.GetRaydiumCPMMQuotesRequest) (*pb.GetRaydiumCPMMQuotesResponse, error)
	PostRaydiumCLMMRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest) (*pb.PostRaydiumRouteSwapResponse, error)
	PostRaydiumCLMMSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest) (*pb.PostRaydiumSwapResponse, error)
	PostRaydiumCPMMSwap(ctx context.Context, request *pb.PostRaydiumCPMMSwapRequest) (*pb.PostRaydiumCPMMSwapResponse, error)
	RecentBlockHash(ctx context.Context) (*pb.GetRecentBlockHashResponse, error)
	GetRecentBlockHash(ctx context.Context) (*pb.GetRecentBlockHashResponse, error)
	GetRecentBlockHashV2(ctx context.Context, offset uint64) (*pb.GetRecentBlockHashResponseV2, error)
	GetPriorityFee(ctx context.Context, request *pb.GetPriorityFeeRequest) (*pb.GetPriorityFeeResponse, error)
	GetPriorityFeeByProgram(ctx context.Context, request *pb.GetPriorityFeeByProgramRequest) (*pb.GetPriorityFeeByProgramResponse, error)
	GetLeaderSchedule(ctx context.Context, request *pb.GetLeaderScheduleRequest) (*pb.GetLeaderScheduleResponse, error)
	GetRateLimit(ctx context.Context, request *pb.GetRateLimitRequest) (*pb.GetRateLimitResponse, error)
	GetTransaction(ctx context.Context, request *pb.GetTransactionRequest) (*pb.GetTransactionResponse, error)
	GetRaydiumPoolReserve(ctx context.Context, req *pb.GetRaydiumPoolReserveRequest) (*pb.GetRaydiumPoolReserveResponse, error)
	GetRaydiumPools(ctx context.Context, request *pb.GetRaydiumPoolsRequest) (*pb.GetRaydiumPoolsResponse, error)
	GetRaydiumQuotes(ctx context.Context, request *pb.GetRaydiumQuotesRequest) (*pb.GetRaydiumQuotesResponse, error)
	GetRaydiumQuotesCPMM(ctx context.Context, request *pb.GetRaydiumCPMMQuotesRequest) (*pb.GetRaydiumCPMMQuotesResponse, error)
	GetPumpFunQuotes(ctx context.Context, request *pb.GetPumpFunQuotesRequest) (*pb.GetPumpFunQuotesResponse, error)
	GetRaydiumPrices(ctx context.Context, request *pb.GetRaydiumPricesRequest) (*pb.GetRaydiumPricesResponse, error)
	SubmitRaydiumCLMMSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitRaydiumCLMMRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	PostRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest) (*pb.PostRaydiumSwapResponse, error)
	PostPumpFunSwap(ctx context.Context, request *pb.PostPumpFunSwapRequest) (*pb.PostPumpFunSwapResponse, error)
	PostRaydiumSwapCPMM(ctx context.Context, request *pb.PostRaydiumCPMMSwapRequest) (*pb.PostRaydiumCPMMSwapResponse, error)
	PostRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest) (*pb.PostRaydiumRouteSwapResponse, error)
	GetJupiterQuotes(ctx context.Context, request *pb.GetJupiterQuotesRequest) (*pb.GetJupiterQuotesResponse, error)
	GetJupiterPrices(ctx context.Context, request *pb.GetJupiterPricesRequest) (*pb.GetJupiterPricesResponse, error)
	PostJupiterSwap(ctx context.Context, request *pb.PostJupiterSwapRequest) (*pb.PostJupiterSwapResponse, error)
	PostJupiterSwapInstructions(ctx context.Context, request *pb.PostJupiterSwapInstructionsRequest) (*pb.PostJupiterSwapInstructionsResponse, error)
	PostRaydiumSwapInstructions(ctx context.Context, request *pb.PostRaydiumSwapInstructionsRequest) (*pb.PostRaydiumSwapInstructionsResponse, error)
	PostJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest) (*pb.PostJupiterRouteSwapResponse, error)
	GetOrderbook(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetOrderbookResponse, error)
	GetMarketDepth(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetMarketDepthResponse, error)
	GetPools(ctx context.Context, projects []pb.Project) (*pb.GetPoolsResponse, error)
	GetTrades(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetTradesResponse, error)
	GetTickers(ctx context.Context, market string, project pb.Project) (*pb.GetTickersResponse, error)
	GetOpenOrders(ctx context.Context, market string, owner string, openOrdersAddress string, project pb.Project) (*pb.GetOpenOrdersResponse, error)
	GetOrderByID(ctx context.Context, in *pb.GetOrderByIDRequest) (*pb.GetOrderByIDResponse, error)
	GetUnsettled(ctx context.Context, market string, ownerAddress string, project pb.Project) (*pb.GetUnsettledResponse, error)
	GetMarkets(ctx context.Context) (*pb.GetMarketsResponse, error)
	GetAccountBalance(ctx context.Context, owner string) (*pb.GetAccountBalanceResponse, error)
	GetTokenAccounts(ctx context.Context, req *pb.GetTokenAccountsRequest) (*pb.GetTokenAccountsResponse, error)
	GetPrice(ctx context.Context, tokens []string) (*pb.GetPriceResponse, error)
	GetQuotes(ctx context.Context, inToken, outToken string, inAmount, slippage float64, limit int32, projects []pb.Project) (*pb.GetQuotesResponse, error)
	SignAndSubmit(ctx context.Context, tx *pb.TransactionMessage,
		skipPreFlight bool, frontRunningProtection bool, useStakedRPCs bool) (string, error)
	SignAndSubmitSnipe(ctx context.Context, transactions []*pb.TransactionMessage, useStakedRPCs bool) ([]string, error)
	SignAndSubmitPaladin(ctx context.Context, tx *pb.TransactionMessage, revertProtection *bool) (string, error)
	PostTradeSwap(ctx context.Context, ownerAddress, inToken, outToken string, inAmount, slippage float64, project pb.Project) (*pb.TradeSwapResponse, error)
	PostRouteTradeSwap(ctx context.Context, request *pb.RouteTradeSwapRequest) (*pb.TradeSwapResponse, error)
	PostOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error)
	PostSubmit(ctx context.Context, tx *pb.TransactionMessage, opts PostSubmitOpts) (*pb.PostSubmitResponse, error)
	PostSubmitBatch(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error)
	PostSubmitV2(ctx context.Context, tx *pb.TransactionMessage, opts PostSubmitOpts) (*pb.PostSubmitResponse, error)
	PostSubmitSnipeV2(ctx context.Context, request *pb.PostSubmitSnipeRequest) (*pb.PostSubmitSnipeResponse, error)
	PostSubmitBatchV2(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error)
	SubmitTradeSwap(ctx context.Context, ownerAddress, inToken, outToken string, inAmount, slippage float64, project pb.Project, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitRouteTradeSwap(ctx context.Context, request *pb.RouteTradeSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitPostPumpFunSwap(ctx context.Context, request *pb.PostPumpFunSwapRequest) (string, error)
	SubmitRaydiumSwapCPMM(ctx context.Context, request *pb.PostRaydiumCPMMSwapRequest) (string, error)
	SubmitRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitJupiterSwap(ctx context.Context, request *pb.PostJupiterSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitJupiterSwapInstructions(ctx context.Context, request *pb.PostJupiterSwapInstructionsRequest, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitRaydiumSwapInstructions(ctx context.Context, request *pb.PostRaydiumSwapInstructionsRequest, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	SubmitJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
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
	PostCancelAll(ctx context.Context, market, owner string, openOrders []string, project pb.Project) (*pb.PostCancelAllResponse, error)
	SubmitCancelAll(ctx context.Context, market, owner string, openOrdersAddresses []string, project pb.Project, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error)
	PostSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, project pb.Project) (*pb.PostSettleResponse, error)
	SubmitSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, project pb.Project, skipPreflight bool) (string, error)
	PostReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error)
	SubmitReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error)
	PostReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error)
	SubmitReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error)
	GetOrderbookStream(ctx context.Context, markets []string, limit uint32, project pb.Project) (connections.Streamer[*pb.GetOrderbooksStreamResponse], error)
	GetPumpFunSwapsStream(ctx context.Context, req *pb.GetPumpFunSwapsStreamRequest) (connections.Streamer[*pb.GetPumpFunSwapsStreamResponse], error)
	GetPumpFunAmmSwapStream(ctx context.Context, req *pb.GetPumpFunAMMSwapStreamRequest) (connections.Streamer[*pb.GetPumpFunAMMSwapStreamResponse], error)
	GetPumpFunNewTokensStream(ctx context.Context, req *pb.GetPumpFunNewTokensStreamRequest) (connections.Streamer[*pb.GetPumpFunNewTokensStreamResponse], error)
	GetMarketDepthsStream(ctx context.Context, markets []string, limit uint32, project pb.Project) (connections.Streamer[*pb.GetMarketDepthsStreamResponse], error)
	GetTradesStream(ctx context.Context, market string, limit uint32, project pb.Project) (connections.Streamer[*pb.GetTradesStreamResponse], error)
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
	GetNewRaydiumPoolsStream(
		ctx context.Context, includeCPMM bool,
	) (connections.Streamer[*pb.GetNewRaydiumPoolsResponse], error)
	GetNewRaydiumPoolsByTransactionStream(
		ctx context.Context, includeCPMM bool,
	) (connections.Streamer[*pb.GetNewRaydiumPoolsByTransactionResponse], error)
	GetBlockStream(ctx context.Context) (connections.Streamer[*pb.GetBlockStreamResponse], error)
	GetPumpFunNewAmmPoolStream(ctx context.Context, req *pb.GetPumpFunNewAmmPoolStreamRequest) (connections.Streamer[*pb.GetPumpFunNewAmmPoolStreamResponse], error)
	GetPriorityFeeStream(ctx context.Context, project pb.Project, percentile *float64) (connections.Streamer[*pb.GetPriorityFeeResponse], error)
	GetPriorityFeeByProgramStream(ctx context.Context, programs []string) (connections.Streamer[*pb.GetPriorityFeeByProgramResponse], error)
	GetBundleTipStream(ctx context.Context) (connections.Streamer[*pb.GetBundleTipResponse], error)
	GetOrderbookV2(ctx context.Context, market string, limit uint32) (*pb.GetOrderbookResponseV2, error)
	GetMarketDepthV2(ctx context.Context, market string, limit uint32) (*pb.GetMarketDepthResponseV2, error)
	GetTickersV2(ctx context.Context, market string) (*pb.GetTickersResponseV2, error)
	GetOpenOrdersV2(ctx context.Context, market string, owner string, openOrdersAddress string, orderID string, clientOrderID uint64) (*pb.GetOpenOrdersResponseV2, error)
	GetUnsettledV2(ctx context.Context, market string, ownerAddress string) (*pb.GetUnsettledResponse, error)
	GetMarketsV2(ctx context.Context) (*pb.GetMarketsResponseV2, error)
	PostOrderV2(ctx context.Context, owner, payer, market string, side string, orderType string, amount, price float64, bundleTip *uint64, opts PostOrderOpts) (*pb.PostOrderResponse, error)
	PostOrderV2WithPriorityFee(ctx context.Context, owner, payer, market string, side string,
		orderType string, amount, price float64, computeLimit uint32, computePrice uint64, bundleTip *uint64, opts PostOrderOpts) (*pb.PostOrderResponse, error)
	SubmitOrderV2(ctx context.Context, owner, payer, market string, side string, orderType string, amount, price float64, bundleTip *uint64, opts PostOrderOpts) (string, error)
	SubmitOrderV2WithPriorityFee(ctx context.Context, owner, payer, market string, side string,
		orderType string, amount, price float64, computeLimit uint32, computePrice uint64, bundleTip *uint64, opts PostOrderOpts) (string, error)
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
	GetRaydiumCLMMPools(ctx context.Context, request *pb.GetRaydiumCLMMPoolsRequest) (*pb.GetRaydiumCLMMPoolsResponse, error)
	GetRaydiumCLMMQuotes(ctx context.Context, request *pb.GetRaydiumCLMMQuotesRequest) (*pb.GetRaydiumCLMMQuotesResponse, error)
}

type GRPCClient struct {
	pb.UnimplementedApiServer

	apiClient pb.ApiClient

	privateKey           *solana.PrivateKey
	recentBlockHashStore *recentBlockHashStore
}

// NewGRPCClientFullService connects to Mainnet Trader API with full service API
func NewGRPCClientFullService(bloxrouteEndpoint string) (GRPCClientTraderAPI, error) {
	if bloxrouteEndpoint != MainnetNYGRPC && bloxrouteEndpoint != MainnetUKGRPC {
		return nil, fmt.Errorf("not a valid endpoint for a full service trader api")
	}

	opts := DefaultRPCOpts(bloxrouteEndpoint)
	opts.UseTLS = true

	return NewGRPCClientWithOpts(opts)
}

// NewGRPCClientSubmitOnly connects to Mainnet Trader API with full service API
func NewGRPCClientSubmitOnly(bloxrouteEndpoint string) (GRPCClientTraderAPISubmitOnly, error) {
	if bloxrouteEndpoint != MainnetAmsterdamGRPC &&
		bloxrouteEndpoint != MainnetFrankfurtGRPC &&
		bloxrouteEndpoint != MainnetLAGRPC &&
		bloxrouteEndpoint != MainnetTokyoGRPC {
		return nil, fmt.Errorf("not a valid endpoint for submit only trader api")
	}

	opts := DefaultRPCOpts(bloxrouteEndpoint)
	opts.UseTLS = true

	return NewGRPCClientWithOpts(opts)
}

// NewGRPCClientPumpNY connects to Mainnet NY Pump Trader API
func NewGRPCClientPumpNY() (GRPCClientTraderAPI, error) {
	opts := DefaultRPCOpts(MainnetPumpNYGRPC)
	opts.UseTLS = true
	return NewGRPCClientWithOpts(opts)
}

// NewGRPCTestnet connects to Testnet Trader API
func NewGRPCTestnet() (GRPCClientTraderAPI, error) {
	opts := DefaultRPCOpts(TestnetGRPC)
	opts.UseTLS = true
	return NewGRPCClientWithOpts(opts)
}

// NewGRPCDevnet connects to Devnet Trader API
func NewGRPCDevnet() (GRPCClientTraderAPI, error) {
	opts := DefaultRPCOpts(DevnetGRPC)
	return NewGRPCClientWithOpts(opts)
}

// NewGRPCLocal connects to local Trader API
func NewGRPCLocal() (GRPCClientTraderAPI, error) {
	opts := DefaultRPCOpts(LocalGRPC)
	return NewGRPCClientWithOpts(opts)
}

type blxrCredentials struct {
	authorization string
}

func (bc blxrCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": bc.authorization,
		"x-sdk":         package_info.Name,
		"x-sdk-version": package_info.Version,
	}, nil
}

func (bc blxrCredentials) RequireTransportSecurity() bool {
	return false
}

// NewGRPCClientWithOpts connects to custom Trader API
func NewGRPCClientWithOpts(opts RPCOpts, dialOpts ...grpc.DialOption) (*GRPCClient, error) {
	var (
		conn     grpc.ClientConnInterface
		err      error
		grpcOpts = make([]grpc.DialOption, 0)
	)

	transportOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	if opts.UseTLS {
		transportOption = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))
	}
	grpcOpts = append(grpcOpts, transportOption)

	if !opts.DisableAuth {
		grpcOpts = append(grpcOpts, grpc.WithPerRPCCredentials(blxrCredentials{authorization: opts.AuthHeader}))
	}
	grpcOpts = append(grpcOpts, grpc.WithDefaultCallOptions(&grpc.MaxRecvMsgSizeCallOption{MaxRecvMsgSize: 1024 * 1024 * 16}))
	grpcOpts = append(grpcOpts, dialOpts...)
	conn, err = grpc.NewClient(opts.Endpoint, grpcOpts...)
	if err != nil {
		return nil, err
	}

	client := &GRPCClient{
		apiClient:  pb.NewApiClient(conn),
		privateKey: opts.PrivateKey,
	}

	client.recentBlockHashStore = newRecentBlockHashStore(
		client.GetRecentBlockHash,
		client.GetRecentBlockHashStream,
		opts,
	)
	if opts.CacheBlockHash {
		go client.recentBlockHashStore.run(context.Background())
	}
	return client, nil
}
