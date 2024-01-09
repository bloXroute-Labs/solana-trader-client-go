package provider

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/connections"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/bloXroute-Labs/solana-trader-proto/common"
	"github.com/gagliardetto/solana-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	pb.UnimplementedApiServer

	apiClient pb.ApiClient

	privateKey           *solana.PrivateKey
	recentBlockHashStore *recentBlockHashStore
}

// NewGRPCClient connects to Mainnet Trader API
func NewGRPCClient() (*GRPCClient, error) {
	opts := DefaultRPCOpts(MainnetNYGRPC)
	opts.UseTLS = true
	return NewGRPCClientWithOpts(opts)
}

// NewGRPCTestnet connects to Testnet Trader API
func NewGRPCTestnet() (*GRPCClient, error) {
	opts := DefaultRPCOpts(TestnetGRPC)
	return NewGRPCClientWithOpts(opts)
}

// NewGRPCDevnet connects to Devnet Trader API
func NewGRPCDevnet() (*GRPCClient, error) {
	opts := DefaultRPCOpts(DevnetGRPC)
	return NewGRPCClientWithOpts(opts)
}

// NewGRPCLocal connects to local Trader API
func NewGRPCLocal() (*GRPCClient, error) {
	opts := DefaultRPCOpts(LocalGRPC)
	return NewGRPCClientWithOpts(opts)
}

type blxrCredentials struct {
	authorization string
}

func (bc blxrCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": bc.authorization,
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

	grpcOpts = append(grpcOpts, dialOpts...)
	conn, err = grpc.Dial(opts.Endpoint, grpcOpts...)
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

func (g *GRPCClient) RecentBlockHash(ctx context.Context) (*pb.GetRecentBlockHashResponse, error) {
	return g.recentBlockHashStore.get(ctx)
}

func (g *GRPCClient) GetRecentBlockHash(ctx context.Context) (*pb.GetRecentBlockHashResponse, error) {
	return g.apiClient.GetRecentBlockHash(ctx, &pb.GetRecentBlockHashRequest{})
}

// GetRaydiumPools returns pools on Raydium
func (g *GRPCClient) GetRaydiumPools(ctx context.Context, request *pb.GetRaydiumPoolsRequest) (*pb.GetRaydiumPoolsResponse, error) {
	return g.apiClient.GetRaydiumPools(ctx, request)
}

// GetRaydiumQuotes returns the possible amount(s) of outToken for an inToken and the route to achieve it on Raydium
func (g *GRPCClient) GetRaydiumQuotes(ctx context.Context, request *pb.GetRaydiumQuotesRequest) (*pb.GetRaydiumQuotesResponse, error) {
	return g.apiClient.GetRaydiumQuotes(ctx, request)
}

// GetRaydiumPrices returns the USDC price of requested tokens on Raydium
func (g *GRPCClient) GetRaydiumPrices(ctx context.Context, request *pb.GetRaydiumPricesRequest) (*pb.GetRaydiumPricesResponse, error) {
	return g.apiClient.GetRaydiumPrices(ctx, request)
}

// PostRaydiumSwap returns a partially signed transaction(s) for submitting a swap request on Raydium
func (g *GRPCClient) PostRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest) (*pb.PostRaydiumSwapResponse, error) {
	return g.apiClient.PostRaydiumSwap(ctx, request)
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

// PostJupiterRouteSwap returns a partially signed transaction(s) for submitting a swap request on Jupiter
func (g *GRPCClient) PostJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest) (*pb.PostJupiterRouteSwapResponse, error) {
	return g.apiClient.PostJupiterRouteSwap(ctx, request)
}

// GetOrderbook returns the requested market's orderbook (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (g *GRPCClient) GetOrderbook(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetOrderbookResponse, error) {
	return g.apiClient.GetOrderbook(ctx, &pb.GetOrderbookRequest{Market: market, Limit: limit, Project: project})
}

// GetMarketDepth returns the requested market's coalesced price data (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (g *GRPCClient) GetMarketDepth(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetMarketDepthResponse, error) {
	return g.apiClient.GetMarketDepth(ctx, &pb.GetMarketDepthRequest{Market: market, Limit: limit, Project: project})
}

// GetPools returns pools for given projects.
func (g *GRPCClient) GetPools(ctx context.Context, projects []pb.Project) (*pb.GetPoolsResponse, error) {
	return g.apiClient.GetPools(ctx, &pb.GetPoolsRequest{Projects: projects})
}

// GetTrades returns the requested market's currently executing trades. Set limit to 0 for all trades.
func (g *GRPCClient) GetTrades(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetTradesResponse, error) {
	return g.apiClient.GetTrades(ctx, &pb.GetTradesRequest{Market: market, Limit: limit, Project: project})
}

// GetTickers returns the requested market tickets. Set market to "" for all markets.
func (g *GRPCClient) GetTickers(ctx context.Context, market string, project pb.Project) (*pb.GetTickersResponse, error) {
	return g.apiClient.GetTickers(ctx, &pb.GetTickersRequest{Market: market, Project: project})
}

// GetOpenOrders returns all open orders by owner address and market
func (g *GRPCClient) GetOpenOrders(ctx context.Context, market string, owner string, openOrdersAddress string, project pb.Project) (*pb.GetOpenOrdersResponse, error) {
	return g.apiClient.GetOpenOrders(ctx, &pb.GetOpenOrdersRequest{Market: market, Address: owner, OpenOrdersAddress: openOrdersAddress, Project: project})
}

// GetOrderByID returns an order by id
func (g *GRPCClient) GetOrderByID(ctx context.Context, in *pb.GetOrderByIDRequest) (*pb.GetOrderByIDResponse, error) {
	return g.apiClient.GetOrderByID(ctx, in)
}

// GetOpenPerpOrders returns all open perp orders
func (g *GRPCClient) GetOpenPerpOrders(ctx context.Context, request *pb.GetOpenPerpOrdersRequest) (*pb.GetOpenPerpOrdersResponse, error) {
	return g.apiClient.GetOpenPerpOrders(ctx, request)
}

// PostCloseDriftPerpPositions returns a partially signed transaction for canceling perp positions on Drift
func (g *GRPCClient) PostCloseDriftPerpPositions(ctx context.Context, request *pb.PostCloseDriftPerpPositionsRequest) (*pb.PostCloseDriftPerpPositionsResponse, error) {
	return g.apiClient.PostCloseDriftPerpPositions(ctx, request)
}

// GetDriftPerpOrderbook returns the current state of perpetual contract orderbook on Drift
func (g *GRPCClient) GetDriftPerpOrderbook(ctx context.Context, request *pb.GetDriftPerpOrderbookRequest) (*pb.GetDriftPerpOrderbookResponse, error) {
	return g.apiClient.GetDriftPerpOrderbook(ctx, request)
}

// PostCreateDriftUser returns a partially signed transaction for creating a user on Drift
func (g *GRPCClient) PostCreateDriftUser(ctx context.Context, request *pb.PostCreateDriftUserRequest) (*pb.PostCreateDriftUserResponse, error) {
	return g.apiClient.PostCreateDriftUser(ctx, request)
}

// GetDriftUser returns a user's info on Drift
func (g *GRPCClient) GetDriftUser(ctx context.Context, request *pb.GetDriftUserRequest) (*pb.GetDriftUserResponse, error) {
	return g.apiClient.GetDriftUser(ctx, request)
}

// PostDriftManageCollateral returns a partially signed transaction for managing collateral on Drift
func (g *GRPCClient) PostDriftManageCollateral(ctx context.Context, request *pb.PostDriftManageCollateralRequest) (*pb.PostDriftManageCollateralResponse, error) {
	return g.apiClient.PostDriftManageCollateral(ctx, request)
}

// PostDriftSettlePNL returns partially signed transactions for settling PNL on Drift
func (g *GRPCClient) PostDriftSettlePNL(ctx context.Context, request *pb.PostDriftSettlePNLRequest) (*pb.PostDriftSettlePNLResponse, error) {
	return g.apiClient.PostDriftSettlePNL(ctx, request)
}

// PostDriftSettlePNLs returns partially signed transactions for settling PNLs on Drift
func (g *GRPCClient) PostDriftSettlePNLs(ctx context.Context, request *pb.PostDriftSettlePNLsRequest) (*pb.PostDriftSettlePNLsResponse, error) {
	return g.apiClient.PostDriftSettlePNLs(ctx, request)
}

// GetDriftAssets returns list of assets for user on Drift
func (g *GRPCClient) GetDriftAssets(ctx context.Context, request *pb.GetDriftAssetsRequest) (*pb.GetDriftAssetsResponse, error) {
	return g.apiClient.GetDriftAssets(ctx, request)
}

// GetDriftPerpContracts returns list of available perp contracts on Drift
func (g *GRPCClient) GetDriftPerpContracts(ctx context.Context, request *pb.GetDriftPerpContractsRequest) (*pb.GetDriftPerpContractsResponse, error) {
	return g.apiClient.GetDriftPerpContracts(ctx, request)
}

// PostLiquidateDriftPerp returns a partially signed transaction for liquidating perp position on Drift
func (g *GRPCClient) PostLiquidateDriftPerp(ctx context.Context, request *pb.PostLiquidateDriftPerpRequest) (*pb.PostLiquidateDriftPerpResponse, error) {
	return g.apiClient.PostLiquidateDriftPerp(ctx, request)
}

// GetDriftOpenPerpOrder returns an open perp order on Drift
func (g *GRPCClient) GetDriftOpenPerpOrder(ctx context.Context, request *pb.GetDriftOpenPerpOrderRequest) (*pb.GetDriftOpenPerpOrderResponse, error) {
	return g.apiClient.GetDriftOpenPerpOrder(ctx, request)
}

// GetDriftOpenMarginOrder return a open margin order on Drift
func (g *GRPCClient) GetDriftOpenMarginOrder(ctx context.Context, request *pb.GetDriftOpenMarginOrderRequest) (*pb.GetDriftOpenMarginOrderResponse, error) {
	return g.apiClient.GetDriftOpenMarginOrder(ctx, request)
}

// GetDriftOpenMarginOrders returns all open margin orders on Drift
func (g *GRPCClient) GetDriftOpenMarginOrders(ctx context.Context, request *pb.GetDriftOpenMarginOrdersRequest) (*pb.GetDriftOpenMarginOrdersResponse, error) {
	return g.apiClient.GetDriftOpenMarginOrders(ctx, request)
}

// PostCancelPerpOrder returns a partially signed transaction for canceling perp order
func (g *GRPCClient) PostCancelPerpOrder(ctx context.Context, request *pb.PostCancelPerpOrderRequest) (*pb.PostCancelPerpOrderResponse, error) {
	return g.apiClient.PostCancelPerpOrder(ctx, request)
}

// PostCancelDriftMarginOrder returns a partially signed transaction for canceling margin orders on Drift
func (g *GRPCClient) PostCancelDriftMarginOrder(ctx context.Context, request *pb.PostCancelDriftMarginOrderRequest) (*pb.PostCancelDriftMarginOrderResponse, error) {
	return g.apiClient.PostCancelDriftMarginOrder(ctx, request)
}

// PostCancelPerpOrders returns a partially signed transaction for canceling all perp orders of a user
func (g *GRPCClient) PostCancelPerpOrders(ctx context.Context, request *pb.PostCancelPerpOrdersRequest) (*pb.PostCancelPerpOrdersResponse, error) {
	return g.apiClient.PostCancelPerpOrders(ctx, request)
}

// PostCreateUser returns a partially signed transaction for creating a user
func (g *GRPCClient) PostCreateUser(ctx context.Context, request *pb.PostCreateUserRequest) (*pb.PostCreateUserResponse, error) {
	return g.apiClient.PostCreateUser(ctx, request)
}

// GetUser returns a user's info
func (g *GRPCClient) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return g.apiClient.GetUser(ctx, request)
}

// PostManageCollateral returns a partially signed transaction for managing collateral
func (g *GRPCClient) PostManageCollateral(ctx context.Context, request *pb.PostManageCollateralRequest) (*pb.PostManageCollateralResponse, error) {
	return g.apiClient.PostManageCollateral(ctx, request)
}

// PostSettlePNL returns a partially signed transaction for settling PNL
func (g *GRPCClient) PostSettlePNL(ctx context.Context, request *pb.PostSettlePNLRequest) (*pb.PostSettlePNLResponse, error) {
	return g.apiClient.PostSettlePNL(ctx, request)
}

// PostSettlePNLs returns partially signed transactions for settling PNLs
func (g *GRPCClient) PostSettlePNLs(ctx context.Context, request *pb.PostSettlePNLsRequest) (*pb.PostSettlePNLsResponse, error) {
	return g.apiClient.PostSettlePNLs(ctx, request)
}

// GetAssets returns list of assets for user
func (g *GRPCClient) GetAssets(ctx context.Context, request *pb.GetAssetsRequest) (*pb.GetAssetsResponse, error) {
	return g.apiClient.GetAssets(ctx, request)
}

// GetPerpContracts returns list of available perp contracts
func (g *GRPCClient) GetPerpContracts(ctx context.Context, request *pb.GetPerpContractsRequest) (*pb.GetPerpContractsResponse, error) {
	return g.apiClient.GetPerpContracts(ctx, request)
}

// PostLiquidatePerp returns a partially signed transaction for liquidating perp position
func (g *GRPCClient) PostLiquidatePerp(ctx context.Context, request *pb.PostLiquidatePerpRequest) (*pb.PostLiquidatePerpResponse, error) {
	return g.apiClient.PostLiquidatePerp(ctx, request)
}

// GetOpenPerpOrder returns an open perp order
func (g *GRPCClient) GetOpenPerpOrder(ctx context.Context, request *pb.GetOpenPerpOrderRequest) (*pb.GetOpenPerpOrderResponse, error) {
	return g.apiClient.GetOpenPerpOrder(ctx, request)
}

// GetPerpPositions returns all perp positions by owner address and market
func (g *GRPCClient) GetPerpPositions(ctx context.Context, request *pb.GetPerpPositionsRequest) (*pb.GetPerpPositionsResponse, error) {
	return g.apiClient.GetPerpPositions(ctx, request)
}

// GetUnsettled returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (g *GRPCClient) GetUnsettled(ctx context.Context, market string, ownerAddress string, project pb.Project) (*pb.GetUnsettledResponse, error) {
	return g.apiClient.GetUnsettled(ctx, &pb.GetUnsettledRequest{Market: market, OwnerAddress: ownerAddress, Project: project})
}

// GetMarkets returns the list of all available named markets
func (g *GRPCClient) GetMarkets(ctx context.Context) (*pb.GetMarketsResponse, error) {
	return g.apiClient.GetMarkets(ctx, &pb.GetMarketsRequest{})
}

// GetDriftMarkets returns the list of all available named markets
func (g *GRPCClient) GetDriftMarkets(ctx context.Context, request *pb.GetDriftMarketsRequest) (*pb.GetDriftMarketsResponse, error) {
	return g.apiClient.GetDriftMarkets(ctx, request)
}

// GetAccountBalance returns all tokens associated with the owner address including Serum unsettled amounts
func (g *GRPCClient) GetAccountBalance(ctx context.Context, owner string) (*pb.GetAccountBalanceResponse, error) {
	return g.apiClient.GetAccountBalance(ctx, &pb.GetAccountBalanceRequest{OwnerAddress: owner})
}

// GetPrice returns the USDC price of requested tokens
func (g *GRPCClient) GetPrice(ctx context.Context, tokens []string) (*pb.GetPriceResponse, error) {
	return g.apiClient.GetPrice(ctx, &pb.GetPriceRequest{Tokens: tokens})
}

// GetQuotes returns the possible amount(s) of outToken for an inToken and the route to achieve it
func (g *GRPCClient) GetQuotes(ctx context.Context, inToken, outToken string, inAmount, slippage float64, limit int32, projects []pb.Project) (*pb.GetQuotesResponse, error) {
	return g.apiClient.GetQuotes(ctx, &pb.GetQuotesRequest{InToken: inToken, OutToken: outToken, InAmount: inAmount, Slippage: slippage, Limit: limit, Projects: projects})
}

// signAndSubmit signs the given transaction and submits it.
func (g *GRPCClient) signAndSubmit(ctx context.Context, tx *pb.TransactionMessage, skipPreFlight bool) (string, error) {
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
	}, skipPreFlight)
	if err != nil {
		return "", err
	}

	return response.Signature, nil
}

// signAndSubmitBatch signs the given transactions and submits them.
func (g *GRPCClient) signAndSubmitBatch(ctx context.Context, transactions []*pb.TransactionMessage, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	if g.privateKey == nil {
		return nil, ErrPrivateKeyNotFound
	}
	batchRequest, err := buildBatchRequest(transactions, *g.privateKey, opts)
	if err != nil {
		return nil, err
	}

	return g.PostSubmitBatch(ctx, batchRequest)
}

// PostTradeSwap returns a partially signed transaction for submitting a swap request
func (g *GRPCClient) PostTradeSwap(ctx context.Context, ownerAddress, inToken, outToken string, inAmount, slippage float64, project pb.Project) (*pb.TradeSwapResponse, error) {
	return g.apiClient.PostTradeSwap(ctx, &pb.TradeSwapRequest{
		OwnerAddress: ownerAddress,
		InToken:      inToken,
		OutToken:     outToken,
		InAmount:     inAmount,
		Slippage:     slippage,
		Project:      project,
	})
}

// PostRouteTradeSwap returns a partially signed transaction(s) for submitting a swap request
func (g *GRPCClient) PostRouteTradeSwap(ctx context.Context, request *pb.RouteTradeSwapRequest) (*pb.TradeSwapResponse, error) {
	return g.apiClient.PostRouteTradeSwap(ctx, request)
}

// PostOrder returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (g *GRPCClient) PostOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	return g.apiClient.PostOrder(ctx, &pb.PostOrderRequest{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              types,
		Amount:            amount,
		Price:             price,
		Project:           project,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
	})
}

// GetDriftPerpPositions returns all perp positions on Drift
func (g *GRPCClient) GetDriftPerpPositions(ctx context.Context, request *pb.GetDriftPerpPositionsRequest) (*pb.GetDriftPerpPositionsResponse, error) {
	return g.apiClient.GetDriftPerpPositions(ctx, request)
}

// GetDriftOpenPerpOrders returns all open perp orders on Drift
func (g *GRPCClient) GetDriftOpenPerpOrders(ctx context.Context, request *pb.GetDriftOpenPerpOrdersRequest) (*pb.GetDriftOpenPerpOrdersResponse, error) {
	return g.apiClient.GetDriftOpenPerpOrders(ctx, request)
}

// PostDriftCancelPerpOrder returns a partially signed transaction for canceling Drift perp order(s)
func (g *GRPCClient) PostDriftCancelPerpOrder(ctx context.Context, request *pb.PostDriftCancelPerpOrderRequest) (*pb.PostDriftCancelPerpOrderResponse, error) {
	return g.apiClient.PostDriftCancelPerpOrder(ctx, request)
}

// PostPerpOrder returns a partially signed transaction for placing a perp order. Typically, you want to use SubmitPerpOrder instead of this.
func (g *GRPCClient) PostPerpOrder(ctx context.Context, request *pb.PostPerpOrderRequest) (*pb.PostPerpOrderResponse, error) {
	return g.apiClient.PostPerpOrder(ctx, request)
}

// PostDriftPerpOrder returns a partially signed transaction for placing a Drift perp order. Typically, you want to use SubmitDriftPerpOrder instead of this.
func (g *GRPCClient) PostDriftPerpOrder(ctx context.Context, request *pb.PostDriftPerpOrderRequest) (*pb.PostDriftPerpOrderResponse, error) {
	return g.apiClient.PostDriftPerpOrder(ctx, request)
}

// PostModifyDriftOrder returns a partially signed transaction for modifying a Drift order. Typically, you want to use SubmitPostModifyDriftOrder instead of this.
func (g *GRPCClient) PostModifyDriftOrder(ctx context.Context, request *pb.PostModifyDriftOrderRequest) (*pb.PostModifyDriftOrderResponse, error) {
	return g.apiClient.PostModifyDriftOrder(ctx, request)
}

// PostDriftMarginOrder returns a partially signed transaction for placing a Margin order. Typically, you want to use SubmitDriftMarginOrder instead of this.
func (g *GRPCClient) PostDriftMarginOrder(ctx context.Context, request *pb.PostDriftMarginOrderRequest) (*pb.PostDriftMarginOrderResponse, error) {
	return g.apiClient.PostDriftMarginOrder(ctx, request)
}

// PostDriftEnableMarginTrading returns a partially signed transaction for enabling/disabling margin trading.
func (g *GRPCClient) PostDriftEnableMarginTrading(ctx context.Context, request *pb.PostDriftEnableMarginTradingRequest) (*pb.PostDriftEnableMarginTradingResponse, error) {
	return g.apiClient.PostDriftEnableMarginTrading(ctx, request)
}

// SubmitDriftEnableMarginTrading builds a perp order, signs it, and submits to the network.
func (g *GRPCClient) SubmitDriftEnableMarginTrading(ctx context.Context, request *pb.PostDriftEnableMarginTradingRequest, skipPreFlight bool) (string, error) {
	tx, err := g.PostDriftEnableMarginTrading(ctx, request)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, tx.Transaction, skipPreFlight)
}

// PostSubmit posts the transaction string to the Solana network.
func (g *GRPCClient) PostSubmit(ctx context.Context, tx *pb.TransactionMessage, skipPreFlight bool) (*pb.PostSubmitResponse, error) {
	return g.apiClient.PostSubmit(ctx, &pb.PostSubmitRequest{Transaction: tx,
		SkipPreFlight: skipPreFlight})
}

// PostSubmitBatch posts a bundle of transactions string based on a specific SubmitStrategy to the Solana network.
func (g *GRPCClient) PostSubmitBatch(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error) {
	return g.apiClient.PostSubmitBatch(ctx, request)
}

// PostSubmitV2 posts the transaction string to the Solana network.
func (g *GRPCClient) PostSubmitV2(ctx context.Context, tx *pb.TransactionMessage, skipPreFlight bool) (*pb.PostSubmitResponse, error) {
	return g.apiClient.PostSubmitV2(ctx, &pb.PostSubmitRequest{Transaction: tx,
		SkipPreFlight: skipPreFlight})
}

// PostSubmitBatchV2 posts a bundle of transactions string based on a specific SubmitStrategy to the Solana network.
func (g *GRPCClient) PostSubmitBatchV2(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error) {
	return g.apiClient.PostSubmitBatchV2(ctx, request)
}

// SubmitTradeSwap builds a TradeSwap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitTradeSwap(ctx context.Context, ownerAddress, inToken, outToken string, inAmount, slippage float64, project pb.Project, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.apiClient.PostTradeSwap(ctx, &pb.TradeSwapRequest{
		OwnerAddress: ownerAddress,
		InToken:      inToken,
		OutToken:     outToken,
		InAmount:     inAmount,
		Slippage:     slippage,
		Project:      project,
	})
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitRouteTradeSwap builds a RouteTradeSwap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitRouteTradeSwap(ctx context.Context, request *pb.RouteTradeSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostRouteTradeSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitRaydiumSwap builds a Raydium Swap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostRaydiumSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitRaydiumRouteSwap builds a Raydium RouteSwap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostRaydiumRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitJupiterSwap builds a Jupiter Swap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitJupiterSwap(ctx context.Context, request *pb.PostJupiterSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostJupiterSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitJupiterRouteSwap builds a Jupiter RouteSwap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostJupiterRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitOrder builds a Serum market order, signs it, and submits to the network.
func (g *GRPCClient) SubmitOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error) {
	order, err := g.PostOrder(ctx, owner, payer, market, side, types, amount, price, project, opts)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

// SubmitPerpOrder builds a perp order, signs it, and submits to the network.
func (g *GRPCClient) SubmitPerpOrder(ctx context.Context, request *pb.PostPerpOrderRequest, skipPreFlight bool) (string, error) {
	order, err := g.PostPerpOrder(ctx, request)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, skipPreFlight)
}

// SubmitDriftPerpOrder builds a Drift perp order, signs it, and submits to the network.
func (g *GRPCClient) SubmitDriftPerpOrder(ctx context.Context, request *pb.PostDriftPerpOrderRequest, skipPreFlight bool) (string, error) {
	order, err := g.PostDriftPerpOrder(ctx, request)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, skipPreFlight)
}

// SubmitCancelPerpOrder builds a cancel perp order txn, signs and submits it to the network.
func (g *GRPCClient) SubmitCancelPerpOrder(ctx context.Context, request *pb.PostCancelPerpOrderRequest, skipPreFlight bool) (string, error) {
	resp, err := g.PostCancelPerpOrder(ctx, request)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, resp.Transaction, skipPreFlight)
}

// SubmitCancelDriftMarginOrder builds a cancel Drift margin order txn, signs and submits it to the network.
func (g *GRPCClient) SubmitCancelDriftMarginOrder(ctx context.Context, request *pb.PostCancelDriftMarginOrderRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostCancelDriftMarginOrder(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitCancelPerpOrders builds a cancel perp orders txn, signs and submits it to the network.
func (g *GRPCClient) SubmitCancelPerpOrders(ctx context.Context, request *pb.PostCancelPerpOrdersRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostCancelPerpOrders(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitDriftCancelPerpOrder builds a Drift cancel perp orders txn, signs and submits it to the network.
func (g *GRPCClient) SubmitDriftCancelPerpOrder(ctx context.Context, request *pb.PostDriftCancelPerpOrderRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostDriftCancelPerpOrder(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Data.Transactions, opts)
}

// PostCancelOrder builds a Serum cancel order.
func (g *GRPCClient) PostCancelOrder(
	ctx context.Context,
	orderID string,
	side pb.Side,
	owner,
	market,
	openOrders string,
	project pb.Project,
) (*pb.PostCancelOrderResponse, error) {
	return g.apiClient.PostCancelOrder(ctx, &pb.PostCancelOrderRequest{
		OrderID:           orderID,
		Side:              side,
		OwnerAddress:      owner,
		MarketAddress:     market,
		OpenOrdersAddress: openOrders,
		Project:           project,
	})
}

// SubmitCancelOrder builds a Serum cancel order, signs and submits it to the network.
func (g *GRPCClient) SubmitCancelOrder(
	ctx context.Context,
	orderID string,
	side pb.Side,
	owner,
	market,
	openOrders string,
	project pb.Project,
	skipPreFlight bool,
) (string, error) {
	order, err := g.PostCancelOrder(ctx, orderID, side, owner, market, openOrders, project)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, skipPreFlight)
}

// PostClosePerpPositions builds cancel perp positions txn.
func (g *GRPCClient) PostClosePerpPositions(ctx context.Context, request *pb.PostClosePerpPositionsRequest) (*pb.PostClosePerpPositionsResponse, error) {
	return g.apiClient.PostClosePerpPositions(ctx, request)
}

// SubmitClosePerpPositions builds a close perp positions txn, signs and submits it to the network.
func (g *GRPCClient) SubmitClosePerpPositions(ctx context.Context, request *pb.PostClosePerpPositionsRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	order, err := g.PostClosePerpPositions(ctx, request)
	if err != nil {
		return nil, err
	}
	var msgs []*pb.TransactionMessage
	for _, txn := range order.Transactions {
		msgs = append(msgs, txn)
	}

	return g.signAndSubmitBatch(ctx, msgs, opts)
}

// SubmitCreateUser builds a create-user txn, signs and submits it to the network.
func (g *GRPCClient) SubmitCreateUser(ctx context.Context, request *pb.PostCreateUserRequest, skipPreFlight bool) (string, error) {
	resp, err := g.PostCreateUser(ctx, request)
	if err != nil {
		return "", err
	}
	return g.signAndSubmit(ctx, resp.Transaction, skipPreFlight)
}

// SubmitPostPerpOrder builds a perp order txn, signs and submits it to the network.
func (g *GRPCClient) SubmitPostPerpOrder(ctx context.Context, request *pb.PostPerpOrderRequest, skipPreFlight bool) (string, error) {
	resp, err := g.PostPerpOrder(ctx, request)
	if err != nil {
		return "", err
	}
	return g.signAndSubmit(ctx, resp.Transaction, skipPreFlight)
}

// SubmitPostDriftPerpOrder builds a Drift perp order txn, signs and submits it to the network.
func (g *GRPCClient) SubmitPostDriftPerpOrder(ctx context.Context, request *pb.PostDriftPerpOrderRequest, skipPreFlight bool) (string, error) {
	resp, err := g.PostDriftPerpOrder(ctx, request)
	if err != nil {
		return "", err
	}
	return g.signAndSubmit(ctx, resp.Transaction, skipPreFlight)
}

// SubmitPostModifyDriftOrder builds a Drift modify-order txn, signs and submits it to the network.
func (g *GRPCClient) SubmitPostModifyDriftOrder(ctx context.Context, request *pb.PostModifyDriftOrderRequest, skipPreFlight bool) (string, error) {
	resp, err := g.PostModifyDriftOrder(ctx, request)
	if err != nil {
		return "", err
	}
	return g.signAndSubmit(ctx, resp.Transaction, skipPreFlight)
}

// SubmitPostMarginOrder builds a create-user txn, signs and submits it to the network.
func (g *GRPCClient) SubmitPostMarginOrder(ctx context.Context, request *pb.PostDriftMarginOrderRequest, skipPreFlight bool) (string, error) {
	resp, err := g.PostDriftMarginOrder(ctx, request)
	if err != nil {
		return "", err
	}
	return g.signAndSubmit(ctx, resp.Transaction, skipPreFlight)
}

// SubmitPostSettlePNL builds a settle-pnl txn, signs and submits it to the network.
func (g *GRPCClient) SubmitPostSettlePNL(ctx context.Context, request *pb.PostSettlePNLRequest, skipPreFlight bool) (string, error) {
	resp, err := g.PostSettlePNL(ctx, request)
	if err != nil {
		return "", err
	}
	return g.signAndSubmit(ctx, resp.Transaction, skipPreFlight)
}

// SubmitPostSettlePNLs builds one or many settle-pnl txn, signs and submits them to the network.
func (g *GRPCClient) SubmitPostSettlePNLs(ctx context.Context, request *pb.PostSettlePNLsRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostSettlePNLs(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitPostLiquidatePerp builds a liquidate-perp txn, signs and submits it to the network.
func (g *GRPCClient) SubmitPostLiquidatePerp(ctx context.Context, request *pb.PostLiquidatePerpRequest, skipPreFlight bool) (string, error) {
	resp, err := g.PostLiquidatePerp(ctx, request)
	if err != nil {
		return "", err
	}
	return g.signAndSubmit(ctx, resp.Transaction, skipPreFlight)
}

// SubmitManageCollateral builds a deposit collateral transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitManageCollateral(ctx context.Context, request *pb.PostManageCollateralRequest, skipPreFlight bool) (string, error) {
	resp, err := g.PostManageCollateral(ctx, request)
	if err != nil {
		return "", err
	}
	return g.signAndSubmit(ctx, resp.Transaction, skipPreFlight)
}

// PostCancelByClientOrderID builds a Serum cancel order by client ID.
func (g *GRPCClient) PostCancelByClientOrderID(
	ctx context.Context,
	clientOrderID uint64,
	owner,
	market,
	openOrders string,
	project pb.Project,

) (*pb.PostCancelOrderResponse, error) {
	return g.apiClient.PostCancelByClientOrderID(ctx, &pb.PostCancelByClientOrderIDRequest{
		ClientOrderID:     clientOrderID,
		OwnerAddress:      owner,
		MarketAddress:     market,
		OpenOrdersAddress: openOrders,
		Project:           project,
	})
}

// SubmitCancelByClientOrderID builds a Serum cancel order by client ID, signs and submits it to the network.
func (g *GRPCClient) SubmitCancelByClientOrderID(
	ctx context.Context,
	clientOrderID uint64,
	owner,
	market,
	openOrders string,
	project pb.Project,
	skipPreFlight bool,
) (string, error) {
	order, err := g.PostCancelByClientOrderID(ctx, clientOrderID, owner, market, openOrders, project)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, skipPreFlight)
}

func (g *GRPCClient) PostCancelAll(ctx context.Context, market, owner string, openOrders []string, project pb.Project) (*pb.PostCancelAllResponse, error) {
	return g.apiClient.PostCancelAll(ctx, &pb.PostCancelAllRequest{
		Market:              market,
		OwnerAddress:        owner,
		OpenOrdersAddresses: openOrders,
		Project:             project,
	})
}

func (g *GRPCClient) SubmitCancelAll(ctx context.Context, market, owner string, openOrdersAddresses []string, project pb.Project, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	orders, err := g.PostCancelAll(ctx, market, owner, openOrdersAddresses, project)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, orders.Transactions, opts)
}

// PostSettle returns a partially signed transaction for settling market funds. Typically, you want to use SubmitSettle instead of this.
func (g *GRPCClient) PostSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, project pb.Project) (*pb.PostSettleResponse, error) {
	return g.apiClient.PostSettle(ctx, &pb.PostSettleRequest{
		OwnerAddress:      owner,
		Market:            market,
		BaseTokenWallet:   baseTokenWallet,
		QuoteTokenWallet:  quoteTokenWallet,
		OpenOrdersAddress: openOrdersAccount,
		Project:           project,
	})
}

// SubmitSettle builds a market SubmitSettle transaction, signs it, and submits to the network.
func (g *GRPCClient) SubmitSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, project pb.Project, skipPreflight bool) (string, error) {
	order, err := g.PostSettle(ctx, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount, project)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, skipPreflight)
}

func (g *GRPCClient) PostReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	return g.apiClient.PostReplaceByClientOrderID(ctx, &pb.PostOrderRequest{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              types,
		Amount:            amount,
		Price:             price,
		Project:           project,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
	})
}

func (g *GRPCClient) SubmitReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error) {
	order, err := g.PostReplaceByClientOrderID(ctx, owner, payer, market, side, types, amount, price, project, opts)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

func (g *GRPCClient) PostReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	return g.apiClient.PostReplaceOrder(ctx, &pb.PostReplaceOrderRequest{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              types,
		Amount:            amount,
		Price:             price,
		Project:           project,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
		OrderID:           orderID,
	})
}

func (g *GRPCClient) SubmitReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error) {
	order, err := g.PostReplaceOrder(ctx, orderID, owner, payer, market, side, types, amount, price, project, opts)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

// GetOrderbookStream subscribes to a stream for changes to the requested market updates (e.g. asks and bids. Set limit to 0 for all bids/ asks).
func (g *GRPCClient) GetOrderbookStream(ctx context.Context, markets []string, limit uint32, project pb.Project) (connections.Streamer[*pb.GetOrderbooksStreamResponse], error) {
	stream, err := g.apiClient.GetOrderbooksStream(ctx, &pb.GetOrderbooksRequest{
		Markets: markets, Limit: limit,
		Project: project})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetOrderbooksStreamResponse](stream, fmt.Sprint(markets)), nil
}

// GetDriftMarginOrderbooksStream subscribes to a stream for changes to the requested market updates (e.g. asks and bids. Set limit to 0 for all bids/ asks).
func (g *GRPCClient) GetDriftMarginOrderbooksStream(ctx context.Context, request *pb.GetDriftMarginOrderbooksRequest) (connections.Streamer[*pb.GetDriftMarginOrderbooksStreamResponse], error) {
	stream, err := g.apiClient.GetDriftMarginOrderbooksStream(ctx, request)
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetDriftMarginOrderbooksStreamResponse](stream, fmt.Sprint(request.GetMarkets())), nil
}

// GetMarketDepthsStream subscribes to a stream for changes to the requested market data updates (e.g. asks and bids. Set limit to 0 for all bids/ asks).
func (g *GRPCClient) GetMarketDepthsStream(ctx context.Context, markets []string, limit uint32, project pb.Project) (connections.Streamer[*pb.GetMarketDepthsStreamResponse], error) {
	stream, err := g.apiClient.GetMarketDepthsStream(ctx, &pb.GetMarketDepthsRequest{Markets: markets, Limit: limit, Project: project})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetMarketDepthsStreamResponse](stream, fmt.Sprint(markets)), nil
}

// GetTradesStream subscribes to a stream for trades as they execute. Set limit to 0 for all trades.
func (g *GRPCClient) GetTradesStream(ctx context.Context, market string, limit uint32, project pb.Project) (connections.Streamer[*pb.GetTradesStreamResponse], error) {
	stream, err := g.apiClient.GetTradesStream(ctx, &pb.GetTradesRequest{Market: market, Limit: limit, Project: project})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetTradesStreamResponse](stream, market), nil
}

// GetOrderStatusStream subscribes to a stream that shows updates to the owner's orders
func (g *GRPCClient) GetOrderStatusStream(ctx context.Context, market, ownerAddress string, project pb.Project) (connections.Streamer[*pb.GetOrderStatusStreamResponse], error) {
	stream, err := g.apiClient.GetOrderStatusStream(ctx, &pb.GetOrderStatusStreamRequest{Market: market, OwnerAddress: ownerAddress, Project: project})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetOrderStatusStreamResponse](stream, market), nil
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
func (g *GRPCClient) GetPoolReservesStream(ctx context.Context, projects []pb.Project) (connections.Streamer[*pb.GetPoolReservesStreamResponse], error) {
	stream, err := g.apiClient.GetPoolReservesStream(ctx, &pb.GetPoolReservesStreamRequest{
		Projects: projects,
	})
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

// GetNewRaydiumPoolsStream subscribes to a stream for getting recent swaps on projects & markets of interest.
func (g *GRPCClient) GetNewRaydiumPoolsStream(
	ctx context.Context,
) (connections.Streamer[*pb.GetNewRaydiumPoolsResponse], error) {
	stream, err := g.apiClient.GetNewRaydiumPoolsStream(ctx, &pb.GetNewRaydiumPoolsRequest{})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetNewRaydiumPoolsResponse](stream, ""), nil
}

// GetBlockStream subscribes to a stream for getting recent blocks.
func (g *GRPCClient) GetBlockStream(ctx context.Context) (connections.Streamer[*pb.GetBlockStreamResponse], error) {
	stream, err := g.apiClient.GetBlockStream(ctx, &pb.GetBlockStreamRequest{})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetBlockStreamResponse](stream, ""), nil
}

// ------- Drift ----

// GetPerpOrderbook returns the current state of perpetual contract orderbook.
func (g *GRPCClient) GetPerpOrderbook(ctx context.Context, request *pb.GetPerpOrderbookRequest) (*pb.GetPerpOrderbookResponse, error) {
	return g.apiClient.GetPerpOrderbook(ctx, request)
}

// GetDriftMarginOrderbook returns the current state of margin contract orderbook.
func (g *GRPCClient) GetDriftMarginOrderbook(ctx context.Context, request *pb.GetDriftMarginOrderbookRequest) (*pb.GetDriftMarginOrderbookResponse, error) {
	return g.apiClient.GetDriftMarginOrderbook(ctx, request)
}

// GetPerpOrderbooksStream subscribes to a stream for perpetual orderbook updates.
func (g *GRPCClient) GetPerpOrderbooksStream(ctx context.Context, request *pb.GetPerpOrderbooksRequest) (connections.Streamer[*pb.GetPerpOrderbooksStreamResponse], error) {
	stream, err := g.apiClient.GetPerpOrderbooksStream(ctx, request)
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetPerpOrderbooksStreamResponse](stream, ""), nil
}

// GetPerpTradesStream subscribes to a stream for trades to the requested contracts
func (g *GRPCClient) GetPerpTradesStream(ctx context.Context, request *pb.GetPerpTradesStreamRequest) (connections.Streamer[*pb.GetPerpTradesStreamResponse], error) {
	stream, err := g.apiClient.GetPerpTradesStream(ctx, request)
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetPerpTradesStreamResponse](stream, ""), nil
}

// GetDriftMarketDepth returns Drift market depth data.
func (g *GRPCClient) GetDriftMarketDepth(ctx context.Context, request *pb.GetDriftMarketDepthRequest) (*pb.GetDriftMarketDepthResponse, error) {
	return g.apiClient.GetDriftMarketDepth(ctx, request)
}

// GetDriftMarketDepthsStream subscribes to a stream for drift market depth data updates.
func (g *GRPCClient) GetDriftMarketDepthsStream(ctx context.Context, request *pb.GetDriftMarketDepthsStreamRequest) (connections.Streamer[*pb.GetDriftMarketDepthStreamResponse], error) {
	stream, err := g.apiClient.GetDriftMarketDepthsStream(ctx, request)
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetDriftMarketDepthStreamResponse](stream, ""), nil
}

// V2 Openbook

// GetOrderbookV2 returns the requested market's orderbook (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (g *GRPCClient) GetOrderbookV2(ctx context.Context, market string, limit uint32) (*pb.GetOrderbookResponseV2, error) {
	return g.apiClient.GetOrderbookV2(ctx, &pb.GetOrderbookRequestV2{Market: market, Limit: limit})
}

// GetMarketDepthV2 returns the requested market's coalesced price data (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (g *GRPCClient) GetMarketDepthV2(ctx context.Context, market string, limit uint32) (*pb.GetMarketDepthResponseV2, error) {
	return g.apiClient.GetMarketDepthV2(ctx, &pb.GetMarketDepthRequestV2{Market: market, Limit: limit})
}

// GetTickersV2 returns the requested market tickets. Set market to "" for all markets.
func (g *GRPCClient) GetTickersV2(ctx context.Context, market string) (*pb.GetTickersResponseV2, error) {
	return g.apiClient.GetTickersV2(ctx, &pb.GetTickersRequestV2{Market: market})
}

// GetOpenOrdersV2 returns all open orders by owner address and market
func (g *GRPCClient) GetOpenOrdersV2(ctx context.Context, market string, owner string, openOrdersAddress string, orderID string, clientOrderID uint64) (*pb.GetOpenOrdersResponseV2, error) {
	return g.apiClient.GetOpenOrdersV2(ctx, &pb.GetOpenOrdersRequestV2{Market: market, Address: owner, OpenOrdersAddress: openOrdersAddress, OrderID: orderID, ClientOrderID: clientOrderID})
}

// GetUnsettledV2 returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (g *GRPCClient) GetUnsettledV2(ctx context.Context, market string, ownerAddress string) (*pb.GetUnsettledResponse, error) {
	return g.apiClient.GetUnsettledV2(ctx, &pb.GetUnsettledRequestV2{Market: market, OwnerAddress: ownerAddress})
}

// GetMarketsV2 returns the list of all available named markets
func (g *GRPCClient) GetMarketsV2(ctx context.Context) (*pb.GetMarketsResponseV2, error) {
	return g.apiClient.GetMarketsV2(ctx, &pb.GetMarketsRequestV2{})
}

// PostOrderV2 returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (g *GRPCClient) PostOrderV2(ctx context.Context, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	return g.apiClient.PostOrderV2(ctx, &pb.PostOrderRequestV2{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              orderType,
		Amount:            amount,
		Price:             price,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
	})
}

// PostOrderV2WithPriorityFee returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (g *GRPCClient) PostOrderV2WithPriorityFee(ctx context.Context, owner, payer, market string, side string,
	orderType string, amount, price float64, computeLimit uint32, computePrice uint64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	return g.apiClient.PostOrderV2(ctx, &pb.PostOrderRequestV2{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              orderType,
		Amount:            amount,
		Price:             price,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ComputeLimit:      computeLimit,
		ComputePrice:      computePrice,
		ClientOrderID:     opts.ClientOrderID,
	})
}

// SubmitOrderV2 builds a Serum market order, signs it, and submits to the network.
func (g *GRPCClient) SubmitOrderV2(ctx context.Context, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := g.PostOrderV2(ctx, owner, payer, market, side, orderType, amount, price, opts)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

// SubmitOrderV2WithPriorityFee builds a Serum market order, signs it, and submits to the network with specified computeLimit and computePrice
func (g *GRPCClient) SubmitOrderV2WithPriorityFee(ctx context.Context, owner, payer, market string, side string,
	orderType string, amount, price float64, computeLimit uint32, computePrice uint64, opts PostOrderOpts) (string, error) {
	order, err := g.PostOrderV2WithPriorityFee(ctx, owner, payer, market, side, orderType, amount, price, computeLimit, computePrice, opts)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

// PostCancelOrderV2 builds a Serum cancel order.
func (g *GRPCClient) PostCancelOrderV2(
	ctx context.Context,
	orderID string,
	clientOrderID uint64,
	side string,
	owner,
	market,
	openOrders string,
) (*pb.PostCancelOrderResponseV2, error) {
	return g.apiClient.PostCancelOrderV2(ctx, &pb.PostCancelOrderRequestV2{
		OrderID:           orderID,
		Side:              side,
		OwnerAddress:      owner,
		MarketAddress:     market,
		OpenOrdersAddress: openOrders,
		ClientOrderID:     clientOrderID,
	})
}

// SubmitCancelOrderV2 builds a Serum cancel order, signs and submits it to the network.
func (g *GRPCClient) SubmitCancelOrderV2(
	ctx context.Context,
	orderID string,
	clientOrderID uint64,
	side string,
	owner,
	market,
	openOrders string,
	opts SubmitOpts,
) (*pb.PostSubmitBatchResponse, error) {
	order, err := g.PostCancelOrderV2(ctx, orderID, clientOrderID, side, owner, market, openOrders)
	if err != nil {
		return nil, err
	}

	return g.signAndSubmitBatch(ctx, order.Transactions, opts)
}

// PostSettleV2 returns a partially signed transaction for settling market funds. Typically, you want to use SubmitSettle instead of this.
func (g *GRPCClient) PostSettleV2(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string) (*pb.PostSettleResponse, error) {
	return g.apiClient.PostSettleV2(ctx, &pb.PostSettleRequestV2{
		OwnerAddress:      owner,
		Market:            market,
		BaseTokenWallet:   baseTokenWallet,
		QuoteTokenWallet:  quoteTokenWallet,
		OpenOrdersAddress: openOrdersAccount,
	})
}

// SubmitSettleV2 builds a market SubmitSettle transaction, signs it, and submits to the network.
func (g *GRPCClient) SubmitSettleV2(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, skipPreflight bool) (string, error) {
	order, err := g.PostSettleV2(ctx, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, skipPreflight)
}

func (g *GRPCClient) PostReplaceOrderV2(ctx context.Context, orderID, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	return g.apiClient.PostReplaceOrderV2(ctx, &pb.PostReplaceOrderRequestV2{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              orderType,
		Amount:            amount,
		Price:             price,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
		OrderID:           orderID,
	})
}

func (g *GRPCClient) SubmitReplaceOrderV2(ctx context.Context, orderID, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := g.PostReplaceOrderV2(ctx, orderID, owner, payer, market, side, orderType, amount, price, opts)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}
