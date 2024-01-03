package provider

import (
	"context"
	"errors"
	"github.com/bloXroute-Labs/solana-trader-client-go/connections"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/bloXroute-Labs/solana-trader-proto/common"
	"github.com/gagliardetto/solana-go"
)

type WSClient struct {
	pb.UnimplementedApiServer

	addr                 string
	conn                 *connections.WS
	privateKey           *solana.PrivateKey
	recentBlockHashStore *recentBlockHashStore
}

// NewWSClient connects to Mainnet Trader API
func NewWSClient() (*WSClient, error) {
	opts := DefaultRPCOpts(MainnetNYWS)
	return NewWSClientWithOpts(opts)
}

// NewWSClientTestnet connects to Testnet Trader API
func NewWSClientTestnet() (*WSClient, error) {
	opts := DefaultRPCOpts(TestnetWS)
	return NewWSClientWithOpts(opts)
}

// NewWSClientDevnet connects to Devnet Trader API
func NewWSClientDevnet() (*WSClient, error) {
	opts := DefaultRPCOpts(DevnetWS)
	return NewWSClientWithOpts(opts)
}

// NewWSClientLocal connects to local Trader API
func NewWSClientLocal() (*WSClient, error) {
	opts := DefaultRPCOpts(LocalWS)
	return NewWSClientWithOpts(opts)
}

// NewWSClientWithOpts connects to custom Trader API
func NewWSClientWithOpts(opts RPCOpts) (*WSClient, error) {
	conn, err := connections.NewWS(opts.Endpoint, opts.AuthHeader)
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

func (w *WSClient) RecentBlockHash(ctx context.Context) (*pb.GetRecentBlockHashResponse, error) {
	return w.recentBlockHashStore.get(ctx)
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

// GetRaydiumPrices returns the USDC price of requested tokens on Raydium
func (w *WSClient) GetRaydiumPrices(ctx context.Context, request *pb.GetRaydiumPricesRequest) (*pb.GetRaydiumPricesResponse, error) {
	var response pb.GetRaydiumPricesResponse
	err := w.conn.Request(ctx, "GetRaydiumPrices", request, &response)
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

// PostRaydiumRouteSwap returns a partially signed transaction(s) for submitting a swap request on Raydium
func (w *WSClient) PostRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest) (*pb.PostRaydiumRouteSwapResponse, error) {
	var response pb.PostRaydiumRouteSwapResponse
	err := w.conn.Request(ctx, "PostRaydiumRouteSwap", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
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

// PostJupiterRouteSwap returns a partially signed transaction(s) for submitting a swap request on Jupiter
func (w *WSClient) PostJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest) (*pb.PostJupiterRouteSwapResponse, error) {
	var response pb.PostJupiterRouteSwapResponse
	err := w.conn.Request(ctx, "PostJupiterRouteSwap", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostCloseDriftPerpPositions returns a partially signed transaction for canceling perp positions on Drift
func (w *WSClient) PostCloseDriftPerpPositions(ctx context.Context, request *pb.PostCloseDriftPerpPositionsRequest) (*pb.PostCloseDriftPerpPositionsResponse, error) {
	var response pb.PostCloseDriftPerpPositionsResponse
	err := w.conn.Request(ctx, "PostCloseDriftPerpPositions", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetDriftPerpOrderbook returns the current state of perpetual contract orderbook on Drift
func (w *WSClient) GetDriftPerpOrderbook(ctx context.Context, request *pb.GetDriftPerpOrderbookRequest) (*pb.GetDriftPerpOrderbookResponse, error) {
	var response pb.GetDriftPerpOrderbookResponse
	err := w.conn.Request(ctx, "GetDriftPerpOrderbook", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostCreateDriftUser returns a partially signed transaction for creating a user on Drift
func (w *WSClient) PostCreateDriftUser(ctx context.Context, request *pb.PostCreateDriftUserRequest) (*pb.PostCreateDriftUserResponse, error) {
	var response pb.PostCreateDriftUserResponse
	err := w.conn.Request(ctx, "PostCreateDriftUser", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetDriftUser returns a user's info on Drift
func (w *WSClient) GetDriftUser(ctx context.Context, request *pb.GetDriftUserRequest) (*pb.GetDriftUserResponse, error) {
	var response pb.GetDriftUserResponse
	err := w.conn.Request(ctx, "GetDriftUser", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostDriftManageCollateral returns a partially signed transaction for managing collateral on Drift
func (w *WSClient) PostDriftManageCollateral(ctx context.Context, request *pb.PostDriftManageCollateralRequest) (*pb.PostDriftManageCollateralResponse, error) {
	var response pb.PostDriftManageCollateralResponse
	err := w.conn.Request(ctx, "PostDriftManageCollateral", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostDriftSettlePNL returns partially signed transactions for settling PNL on Drift
func (w *WSClient) PostDriftSettlePNL(ctx context.Context, request *pb.PostDriftSettlePNLRequest) (*pb.PostDriftSettlePNLResponse, error) {
	var response pb.PostDriftSettlePNLResponse
	err := w.conn.Request(ctx, "PostDriftSettlePNL", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostDriftSettlePNLs returns partially signed transactions for settling PNLs on Drift
func (w *WSClient) PostDriftSettlePNLs(ctx context.Context, request *pb.PostDriftSettlePNLsRequest) (*pb.PostDriftSettlePNLsResponse, error) {
	var response pb.PostDriftSettlePNLsResponse
	err := w.conn.Request(ctx, "PostDriftSettlePNLs", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetDriftAssets returns list of assets for user on Drift
func (w *WSClient) GetDriftAssets(ctx context.Context, request *pb.GetDriftAssetsRequest) (*pb.GetDriftAssetsResponse, error) {
	var response pb.GetDriftAssetsResponse
	err := w.conn.Request(ctx, "GetDriftAssets", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetDriftPerpContracts returns list of available perp contracts on Drift
func (w *WSClient) GetDriftPerpContracts(ctx context.Context, request *pb.GetDriftPerpContractsRequest) (*pb.GetDriftPerpContractsResponse, error) {
	var response pb.GetDriftPerpContractsResponse
	err := w.conn.Request(ctx, "GetDriftPerpContracts", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostLiquidateDriftPerp returns a partially signed transaction for liquidating perp position on Drift
func (w *WSClient) PostLiquidateDriftPerp(ctx context.Context, request *pb.PostLiquidateDriftPerpRequest) (*pb.PostLiquidateDriftPerpResponse, error) {
	var response pb.PostLiquidateDriftPerpResponse
	err := w.conn.Request(ctx, "PostLiquidateDriftPerp", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetDriftOpenPerpOrder returns an open perp order on Drift
func (w *WSClient) GetDriftOpenPerpOrder(ctx context.Context, request *pb.GetDriftOpenPerpOrderRequest) (*pb.GetDriftOpenPerpOrderResponse, error) {
	var response pb.GetDriftOpenPerpOrderResponse
	err := w.conn.Request(ctx, "GetDriftOpenPerpOrder", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetDriftOpenMarginOrder return a open margin order on Drift
func (w *WSClient) GetDriftOpenMarginOrder(ctx context.Context, request *pb.GetDriftOpenMarginOrderRequest) (*pb.GetDriftOpenMarginOrderResponse, error) {
	var response pb.GetDriftOpenMarginOrderResponse
	err := w.conn.Request(ctx, "GetDriftOpenMarginOrder", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetDriftPerpPositions returns all perp positions on Drift
func (w *WSClient) GetDriftPerpPositions(ctx context.Context, request *pb.GetDriftPerpPositionsRequest) (*pb.GetDriftPerpPositionsResponse, error) {
	var response pb.GetDriftPerpPositionsResponse
	err := w.conn.Request(ctx, "GetDriftPerpPositions", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetDriftOpenPerpOrders returns all open perp orders on Drift
func (w *WSClient) GetDriftOpenPerpOrders(ctx context.Context, request *pb.GetDriftOpenPerpOrdersRequest) (*pb.GetDriftOpenPerpOrdersResponse, error) {
	var response pb.GetDriftOpenPerpOrdersResponse
	err := w.conn.Request(ctx, "GetDriftOpenPerpOrders", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostDriftCancelPerpOrder returns a partially signed transaction for canceling Drift perp order(s)
func (w *WSClient) PostDriftCancelPerpOrder(ctx context.Context, request *pb.PostDriftCancelPerpOrderRequest) (*pb.PostDriftCancelPerpOrderResponse, error) {
	var response pb.PostDriftCancelPerpOrderResponse
	err := w.conn.Request(ctx, "PostDriftCancelPerpOrder", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetOrderbook returns the requested market's orderbook (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (w *WSClient) GetOrderbook(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetOrderbookResponse, error) {
	var response pb.GetOrderbookResponse
	err := w.conn.Request(ctx, "GetOrderbook", &pb.GetOrderbookRequest{Market: market, Limit: limit, Project: project}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetMarketDepth returns the requested market's coalesced price data (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (w *WSClient) GetMarketDepth(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetMarketDepthResponse, error) {
	var response pb.GetMarketDepthResponse
	err := w.conn.Request(ctx, "GetMarketDepth", &pb.GetMarketDepthRequest{Market: market, Limit: limit, Project: project}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetTrades returns the requested market's currently executing trades. Set limit to 0 for all trades.
func (w *WSClient) GetTrades(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetTradesResponse, error) {
	var response pb.GetTradesResponse
	err := w.conn.Request(ctx, "GetTrades", &pb.GetTradesRequest{Market: market, Limit: limit, Project: project}, &response)
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

// GetTickers returns the requested market tickets. Set market to "" for all markets.
func (w *WSClient) GetTickers(ctx context.Context, market string, project pb.Project) (*pb.GetTickersResponse, error) {
	var response pb.GetTickersResponse
	err := w.conn.Request(ctx, "GetTickers", &pb.GetTickersRequest{Market: market, Project: project}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetOpenOrders returns all open orders by owner address and market
func (w *WSClient) GetOpenOrders(ctx context.Context, market string, owner string, openOrdersAddress string, project pb.Project) (*pb.GetOpenOrdersResponse, error) {
	var response pb.GetOpenOrdersResponse
	err := w.conn.Request(ctx, "GetOpenOrders", &pb.GetOpenOrdersRequest{Market: market, Address: owner, OpenOrdersAddress: openOrdersAddress, Project: project}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetOrderByID returns an order by id
func (w *WSClient) GetOrderByID(ctx context.Context, in *pb.GetOrderByIDRequest) (*pb.GetOrderByIDResponse, error) {
	var response pb.GetOrderByIDResponse
	err := w.conn.Request(ctx, "GetOrderByID", &pb.GetOrderByIDRequest{OrderID: in.OrderID, Market: in.Market, Project: in.Project}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetOpenPerpOrders returns all open perp orders
func (w *WSClient) GetOpenPerpOrders(ctx context.Context, request *pb.GetOpenPerpOrdersRequest) (*pb.GetOpenPerpOrdersResponse, error) {
	var response pb.GetOpenPerpOrdersResponse
	err := w.conn.Request(ctx, "GetOpenPerpOrders", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetDriftOpenMarginOrders returns all open margin orders on Drift
func (w *WSClient) GetDriftOpenMarginOrders(ctx context.Context, request *pb.GetDriftOpenMarginOrdersRequest) (*pb.GetDriftOpenMarginOrdersResponse, error) {
	var response pb.GetDriftOpenMarginOrdersResponse
	err := w.conn.Request(ctx, "GetDriftOpenMarginOrders", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostCancelPerpOrder returns a partially signed transaction for canceling perp order
func (w *WSClient) PostCancelPerpOrder(ctx context.Context, request *pb.PostCancelPerpOrderRequest) (*pb.PostCancelPerpOrderResponse, error) {
	var response pb.PostCancelPerpOrderResponse
	err := w.conn.Request(ctx, "PostCancelPerpOrder", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostCancelDriftMarginOrder returns a partially signed transaction for canceling margin orders on Drift
func (w *WSClient) PostCancelDriftMarginOrder(ctx context.Context, request *pb.PostCancelDriftMarginOrderRequest) (*pb.PostCancelDriftMarginOrderResponse, error) {
	var response pb.PostCancelDriftMarginOrderResponse
	err := w.conn.Request(ctx, "PostCancelDriftMarginOrder", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostCancelPerpOrders returns a partially signed transaction for canceling all perp orders of a user
func (w *WSClient) PostCancelPerpOrders(ctx context.Context, request *pb.PostCancelPerpOrdersRequest) (*pb.PostCancelPerpOrdersResponse, error) {
	var response pb.PostCancelPerpOrdersResponse
	err := w.conn.Request(ctx, "PostCancelPerpOrders", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostCreateUser returns a partially signed transaction for creating a user
func (w *WSClient) PostCreateUser(ctx context.Context, request *pb.PostCreateUserRequest) (*pb.PostCreateUserResponse, error) {
	var response pb.PostCreateUserResponse
	err := w.conn.Request(ctx, "PostCreateUser", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetUser returns a user's info
func (w *WSClient) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	var response pb.GetUserResponse
	err := w.conn.Request(ctx, "GetUser", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostManageCollateral returns a partially signed transaction for managing collateral
func (w *WSClient) PostManageCollateral(ctx context.Context, req *pb.PostManageCollateralRequest) (*pb.PostManageCollateralResponse, error) {
	var response pb.PostManageCollateralResponse
	err := w.conn.Request(ctx, "PostManageCollateral", req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostSettlePNL returns a partially signed transaction for settling PNL
func (w *WSClient) PostSettlePNL(ctx context.Context, request *pb.PostSettlePNLRequest) (*pb.PostSettlePNLResponse, error) {
	var response pb.PostSettlePNLResponse
	err := w.conn.Request(ctx, "PostSettlePNL", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostSettlePNLs returns partially signed transactions for settling PNLs
func (w *WSClient) PostSettlePNLs(ctx context.Context, request *pb.PostSettlePNLsRequest) (*pb.PostSettlePNLsResponse, error) {
	var response pb.PostSettlePNLsResponse
	err := w.conn.Request(ctx, "PostSettlePNLs", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetAssets returns list of assets for user
func (w *WSClient) GetAssets(ctx context.Context, request *pb.GetAssetsRequest) (*pb.GetAssetsResponse, error) {
	var response pb.GetAssetsResponse
	err := w.conn.Request(ctx, "GetAssets", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetPerpContracts returns list of available perp contracts
func (w *WSClient) GetPerpContracts(ctx context.Context, request *pb.GetPerpContractsRequest) (*pb.GetPerpContractsResponse, error) {
	var response pb.GetPerpContractsResponse
	err := w.conn.Request(ctx, "GetPerpContracts", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostLiquidatePerp returns a partially signed transaction for liquidating perp position
func (w *WSClient) PostLiquidatePerp(ctx context.Context, request *pb.PostLiquidatePerpRequest) (*pb.PostLiquidatePerpResponse, error) {
	var response pb.PostLiquidatePerpResponse
	err := w.conn.Request(ctx, "PostLiquidatePerp", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetOpenPerpOrder returns an open perp order
func (w *WSClient) GetOpenPerpOrder(ctx context.Context, request *pb.GetOpenPerpOrderRequest) (*pb.GetOpenPerpOrderResponse, error) {
	var response pb.GetOpenPerpOrderResponse
	err := w.conn.Request(ctx, "GetOpenPerpOrder", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetPerpPositions returns all perp positions by owner address and market
func (w *WSClient) GetPerpPositions(ctx context.Context, request *pb.GetPerpPositionsRequest) (*pb.GetPerpPositionsResponse, error) {
	var response pb.GetPerpPositionsResponse
	err := w.conn.Request(ctx, "GetPerpPositions", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetUnsettled returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (w *WSClient) GetUnsettled(ctx context.Context, market string, ownerAddress string, project pb.Project) (*pb.GetUnsettledResponse, error) {
	var response pb.GetUnsettledResponse
	err := w.conn.Request(ctx, "GetUnsettled", &pb.GetUnsettledRequest{Market: market, OwnerAddress: ownerAddress, Project: project}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetAccountBalance returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (w *WSClient) GetAccountBalance(ctx context.Context, owner string) (*pb.GetAccountBalanceResponse, error) {
	var response pb.GetAccountBalanceResponse
	err := w.conn.Request(ctx, "GetAccountBalance", &pb.GetAccountBalanceRequest{OwnerAddress: owner}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetMarkets returns the list of all available named markets
func (w *WSClient) GetMarkets(ctx context.Context) (*pb.GetMarketsResponse, error) {
	var response pb.GetMarketsResponse
	err := w.conn.Request(ctx, "GetMarkets", &pb.GetMarketsRequest{}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetDriftMarkets returns the list of all available named markets
func (w *WSClient) GetDriftMarkets(ctx context.Context, request *pb.GetDriftMarketsRequest) (*pb.GetDriftMarketsResponse, error) {
	var response pb.GetDriftMarketsResponse
	err := w.conn.Request(ctx, "GetDriftMarkets", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetPrice returns the USDC price of requested tokens
func (w *WSClient) GetPrice(ctx context.Context, tokens []string) (*pb.GetPriceResponse, error) {
	var response pb.GetPriceResponse
	err := w.conn.Request(ctx, "GetPrice", &pb.GetPriceRequest{Tokens: tokens}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetQuotes returns the possible amount(s) of outToken for an inToken and the route to achieve it
func (w *WSClient) GetQuotes(ctx context.Context, inToken, outToken string, inAmount, slippage float64, limit int32, projects []pb.Project) (*pb.GetQuotesResponse, error) {
	var response pb.GetQuotesResponse
	request := &pb.GetQuotesRequest{InToken: inToken, OutToken: outToken, InAmount: inAmount, Slippage: slippage, Limit: limit, Projects: projects}

	err := w.conn.Request(ctx, "GetQuotes", request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// PostTradeSwap returns a partially signed transaction for submitting a swap request
func (w *WSClient) PostTradeSwap(ctx context.Context, ownerAddress, inToken, outToken string, inAmount, slippage float64, projectStr string) (*pb.TradeSwapResponse, error) {
	project, err := ProjectFromString(projectStr)
	if err != nil {
		return nil, err
	}
	request := &pb.TradeSwapRequest{
		OwnerAddress: ownerAddress,
		InToken:      inToken,
		OutToken:     outToken,
		InAmount:     inAmount,
		Slippage:     slippage,
		Project:      project,
	}

	var response pb.TradeSwapResponse
	err = w.conn.Request(ctx, "PostTradeSwap", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostRouteTradeSwap returns a partially signed transaction(s) for submitting a swap request
func (w *WSClient) PostRouteTradeSwap(ctx context.Context, request *pb.RouteTradeSwapRequest) (*pb.TradeSwapResponse, error) {
	var response pb.TradeSwapResponse
	err := w.conn.Request(ctx, "PostRouteTradeSwap", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostOrder returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (w *WSClient) PostOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	request := &pb.PostOrderRequest{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              types,
		Amount:            amount,
		Price:             price,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
		Project:           project,
	}
	var response pb.PostOrderResponse
	err := w.conn.Request(ctx, "PostOrder", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostPerpOrder returns a partially signed transaction for placing a perp order. Typically, you want to use SubmitPerpOrder instead of this.
func (w *WSClient) PostPerpOrder(ctx context.Context, request *pb.PostPerpOrderRequest) (*pb.PostPerpOrderResponse, error) {
	var response pb.PostPerpOrderResponse
	err := w.conn.Request(ctx, "PostPerpOrder", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostDriftPerpOrder returns a partially signed transaction for placing a Drift perp order. Typically, you want to use SubmitDriftPerpOrder instead of this.
func (w *WSClient) PostDriftPerpOrder(ctx context.Context, request *pb.PostDriftPerpOrderRequest) (*pb.PostDriftPerpOrderResponse, error) {
	var response pb.PostDriftPerpOrderResponse
	err := w.conn.Request(ctx, "PostDriftPerpOrder", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostModifyDriftOrder returns a partially signed transaction for modifying a Drift order. Typically, you want to use SubmitPostModifyDriftOrder instead of this.
func (w *WSClient) PostModifyDriftOrder(ctx context.Context, request *pb.PostModifyDriftOrderRequest) (*pb.PostModifyDriftOrderResponse, error) {
	var response pb.PostModifyDriftOrderResponse
	err := w.conn.Request(ctx, "PostModifyDriftOrder", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostDriftMarginOrder returns a partially signed transaction for placing a margin order. Typically, you want to use SubmitPostMarginOrder instead of this.
func (w *WSClient) PostDriftMarginOrder(ctx context.Context, request *pb.PostDriftMarginOrderRequest) (*pb.PostDriftMarginOrderResponse, error) {
	var response pb.PostDriftMarginOrderResponse
	err := w.conn.Request(ctx, "PostDriftMarginOrder", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostDriftEnableMarginTrading returns a partially signed transaction for enabling/disabling margin trading.
func (w *WSClient) PostDriftEnableMarginTrading(ctx context.Context, request *pb.PostDriftEnableMarginTradingRequest) (*pb.PostDriftEnableMarginTradingResponse, error) {
	var response pb.PostDriftEnableMarginTradingResponse
	err := w.conn.Request(ctx, "PostDriftEnableMarginTrading", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitDriftEnableMarginTrading builds a perp order, signs it, and submits to the network.
func (w *WSClient) SubmitDriftEnableMarginTrading(ctx context.Context, request *pb.PostDriftEnableMarginTradingRequest, skipPreFlight bool) (string, error) {
	tx, err := w.PostDriftEnableMarginTrading(ctx, request)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, tx.Transaction, skipPreFlight)
}

// PostSubmit posts the transaction string to the Solana network.
func (w *WSClient) PostSubmit(ctx context.Context, txBase64 string, skipPreFlight bool) (*pb.PostSubmitResponse, error) {
	request := &pb.PostSubmitRequest{
		Transaction: &pb.TransactionMessage{
			Content: txBase64,
		},
		SkipPreFlight: skipPreFlight,
	}
	var response pb.PostSubmitResponse
	err := w.conn.Request(ctx, "PostSubmit", request, &response)
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
func (w *WSClient) PostSubmitV2(ctx context.Context, txBase64 string, skipPreFlight bool) (*pb.PostSubmitResponse, error) {
	request := &pb.PostSubmitRequest{
		Transaction: &pb.TransactionMessage{
			Content: txBase64,
		},
		SkipPreFlight: skipPreFlight,
	}
	var response pb.PostSubmitResponse
	err := w.conn.Request(ctx, "PostSubmitV2", request, &response)
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

// signAndSubmit signs the given transaction and submits it.
func (w *WSClient) signAndSubmit(ctx context.Context, tx *pb.TransactionMessage, skipPreFlight bool) (string, error) {
	if w.privateKey == nil {
		return "", ErrPrivateKeyNotFound
	}

	txBase64, err := transaction.SignTxWithPrivateKey(tx.Content, *w.privateKey)
	if err != nil {
		return "", err
	}

	response, err := w.PostSubmit(ctx, txBase64, skipPreFlight)
	if err != nil {
		return "", err
	}

	return response.Signature, nil
}

// signAndSubmitBatch signs the given transactions and submits them.
func (w *WSClient) signAndSubmitBatch(ctx context.Context, transactions []*pb.TransactionMessage, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	if w.privateKey == nil {
		return nil, ErrPrivateKeyNotFound
	}
	batchRequest, err := buildBatchRequest(transactions, *w.privateKey, opts)
	if err != nil {
		return nil, err
	}
	return w.PostSubmitBatch(ctx, batchRequest)
}

// SubmitTradeSwap builds a TradeSwap transaction then signs it, and submits to the network.
func (w *WSClient) SubmitTradeSwap(ctx context.Context, owner, inToken, outToken string, inAmount, slippage float64, project string, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := w.PostTradeSwap(ctx, owner, inToken, outToken, inAmount, slippage, project)
	if err != nil {
		return nil, err
	}
	return w.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitRouteTradeSwap builds a RouteTradeSwap transaction then signs it, and submits to the network.
func (w *WSClient) SubmitRouteTradeSwap(ctx context.Context, request *pb.RouteTradeSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := w.PostRouteTradeSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return w.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitRaydiumSwap builds a Raydium Swap transaction then signs it, and submits to the network.
func (w *WSClient) SubmitRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := w.PostRaydiumSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return w.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitRaydiumRouteSwap builds a Raydium RouteSwap transaction then signs it, and submits to the network.
func (w *WSClient) SubmitRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := w.PostRaydiumRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return w.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitJupiterSwap builds a Jupiter Swap transaction then signs it, and submits to the network.
func (w *WSClient) SubmitJupiterSwap(ctx context.Context, request *pb.PostJupiterSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := w.PostJupiterSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return w.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitJupiterRouteSwap builds a Jupiter RouteSwap transaction then signs it, and submits to the network.
func (w *WSClient) SubmitJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := w.PostJupiterRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return w.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitPerpOrder builds a perp order, signs it, and submits to the network.
func (w *WSClient) SubmitPerpOrder(ctx context.Context, request *pb.PostPerpOrderRequest, opts PostOrderOpts) (string, error) {

	order, err := w.PostPerpOrder(ctx, request)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

// SubmitPostDriftMarginOrder builds a margin order, signs it, and submits to the network.
func (w *WSClient) SubmitPostDriftMarginOrder(ctx context.Context, request *pb.PostDriftMarginOrderRequest, opts PostOrderOpts) (string, error) {

	order, err := w.PostDriftMarginOrder(ctx, request)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

// SubmitOrder builds a Serum market order, signs it, and submits to the network.
func (w *WSClient) SubmitOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error) {
	order, err := w.PostOrder(ctx, owner, payer, market, side, types, amount, price, project, opts)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

// PostCancelOrder builds a Serum cancel order.
func (w *WSClient) PostCancelOrder(ctx context.Context, request *pb.PostCancelOrderRequest) (*pb.PostCancelOrderResponse, error) {
	var response pb.PostCancelOrderResponse
	err := w.conn.Request(ctx, "PostCancelOrder", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitCancelOrder builds a Serum cancel order, signs and submits it to the network.
func (w *WSClient) SubmitCancelOrder(ctx context.Context, request *pb.PostCancelOrderRequest, skipPreFlight bool) (string, error) {
	order, err := w.PostCancelOrder(ctx, request)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, order.Transaction, skipPreFlight)
}

// PostClosePerpPositions builds cancel perp positions txn.
func (w *WSClient) PostClosePerpPositions(ctx context.Context, request *pb.PostClosePerpPositionsRequest) (*pb.PostClosePerpPositionsResponse, error) {
	var response pb.PostClosePerpPositionsResponse
	err := w.conn.Request(ctx, "PostClosePerpPositions", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitClosePerpPositions builds a close perp positions txn, signs and submits it to the network.
func (w *WSClient) SubmitClosePerpPositions(ctx context.Context, request *pb.PostClosePerpPositionsRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := w.PostClosePerpPositions(ctx, request)
	if err != nil {
		return nil, err
	}
	var msgs []*pb.TransactionMessage
	for _, txn := range resp.Transactions {
		msgs = append(msgs, txn)
	}

	return w.signAndSubmitBatch(ctx, msgs, opts)
}

// SubmitCancelPerpOrder builds a cancel perp order txn, signs and submits it to the network.
func (w *WSClient) SubmitCancelPerpOrder(ctx context.Context, request *pb.PostCancelPerpOrderRequest, skipPreFlight bool,
) (string, error) {
	resp, err := w.PostCancelPerpOrder(ctx, request)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, resp.Transaction, skipPreFlight)
}

// SubmitCancelDriftMarginOrder builds a cancel Drift margin order txn , signs and submits it to the network.
func (w *WSClient) SubmitCancelDriftMarginOrder(ctx context.Context, request *pb.PostCancelDriftMarginOrderRequest, skipPreFlight bool) (*pb.PostSubmitBatchResponse, error) {
	resp, err := w.PostCancelDriftMarginOrder(ctx, request)
	if err != nil {
		return nil, err
	}
	return w.signAndSubmitBatch(ctx, resp.Transactions, SubmitOpts{
		SkipPreFlight: skipPreFlight,
	})
}

// SubmitCancelPerpOrders builds a cancel perp orders txn, signs and submits it to the network.
func (w *WSClient) SubmitCancelPerpOrders(ctx context.Context, request *pb.PostCancelPerpOrdersRequest, skipPreFlight bool) (*pb.PostSubmitBatchResponse, error) {
	resp, err := w.PostCancelPerpOrders(ctx, request)
	if err != nil {
		return nil, err
	}
	return w.signAndSubmitBatch(ctx, resp.Transactions, SubmitOpts{
		SkipPreFlight: skipPreFlight,
	})
}

// PostCancelByClientOrderID builds a Serum cancel order by client ID.
func (w *WSClient) PostCancelByClientOrderID(
	ctx context.Context,
	clientOrderID uint64,
	owner,
	market,
	openOrders string,
	project pb.Project,
) (*pb.PostCancelOrderResponse, error) {
	request := &pb.PostCancelByClientOrderIDRequest{
		ClientOrderID:     clientOrderID,
		OwnerAddress:      owner,
		MarketAddress:     market,
		OpenOrdersAddress: openOrders,
		Project:           project,
	}
	var response pb.PostCancelOrderResponse
	err := w.conn.Request(ctx, "PostCancelByClientOrderID", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitCancelByClientOrderID builds a Serum cancel order by client ID, signs and submits it to the network.
func (w *WSClient) SubmitCancelByClientOrderID(
	ctx context.Context,
	clientOrderID uint64,
	owner,
	market,
	openOrders string,
	project pb.Project,
	skipPreFlight bool,
) (string, error) {
	order, err := w.PostCancelByClientOrderID(ctx, clientOrderID, owner, market, openOrders, project)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, order.Transaction, skipPreFlight)
}

func (w *WSClient) PostCancelAll(
	ctx context.Context,
	market,
	owner string,
	openOrdersAddresses []string,
	project pb.Project,
) (*pb.PostCancelAllResponse, error) {
	request := &pb.PostCancelAllRequest{
		Market:              market,
		OwnerAddress:        owner,
		OpenOrdersAddresses: openOrdersAddresses,
		Project:             project,
	}
	var response pb.PostCancelAllResponse
	err := w.conn.Request(ctx, "PostCancelAll", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (w *WSClient) SubmitCancelAll(ctx context.Context, market, owner string, openOrdersAddresses []string, project pb.Project, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	orders, err := w.PostCancelAll(ctx, market, owner, openOrdersAddresses, project)
	if err != nil {
		return nil, err
	}
	return w.signAndSubmitBatch(ctx, orders.Transactions, opts)
}

// PostSettle returns a partially signed transaction for settling market funds. Typically, you want to use SubmitSettle instead of this.
func (w *WSClient) PostSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, project pb.Project) (*pb.PostSettleResponse, error) {
	request := &pb.PostSettleRequest{
		OwnerAddress:      owner,
		Market:            market,
		BaseTokenWallet:   baseTokenWallet,
		QuoteTokenWallet:  quoteTokenWallet,
		OpenOrdersAddress: openOrdersAccount,
		Project:           project,
	}
	var response pb.PostSettleResponse
	err := w.conn.Request(ctx, "PostSettle", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitSettle builds a market SubmitSettle transaction, signs it, and submits to the network.
func (w *WSClient) SubmitSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, project pb.Project, skipPreflight bool) (string, error) {
	order, err := w.PostSettle(ctx, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount, project)
	if err != nil {
		return "", err
	}
	return w.signAndSubmit(ctx, order.Transaction, skipPreflight)
}

func (w *WSClient) PostReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	request := &pb.PostOrderRequest{
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
	}
	var response pb.PostOrderResponse
	err := w.conn.Request(ctx, "PostReplaceByClientOrderID", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (w *WSClient) SubmitReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error) {
	order, err := w.PostReplaceByClientOrderID(ctx, owner, payer, market, side, types, amount, price, project, opts)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

func (w *WSClient) PostReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	request := &pb.PostReplaceOrderRequest{
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
	}
	var response pb.PostOrderResponse
	err := w.conn.Request(ctx, "PostReplaceOrder", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (w *WSClient) SubmitReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error) {
	order, err := w.PostReplaceOrder(ctx, orderID, owner, payer, market, side, types, amount, price, project, opts)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

// SubmitPostSettlePNL builds a settle-pnl txn, signs and submits it to the network.
func (w *WSClient) SubmitPostSettlePNL(ctx context.Context, request *pb.PostSettlePNLRequest, skipPreFlight bool) (string, error) {
	resp, err := w.PostSettlePNL(ctx, request)
	if err != nil {
		return "", err
	}
	return w.signAndSubmit(ctx, resp.Transaction, skipPreFlight)
}

// SubmitPostSettlePNLs builds one or many settle-pnl txn, signs and submits them to the network.
func (w *WSClient) SubmitPostSettlePNLs(ctx context.Context, request *pb.PostSettlePNLsRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := w.PostSettlePNLs(ctx, request)
	if err != nil {
		return nil, err
	}
	return w.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitPostLiquidatePerp builds a liquidate-perp txn, signs and submits it to the network.
func (w *WSClient) SubmitPostLiquidatePerp(ctx context.Context, request *pb.PostLiquidatePerpRequest, skipPreFlight bool) (string, error) {
	resp, err := w.PostLiquidatePerp(ctx, request)
	if err != nil {
		return "", err
	}
	return w.signAndSubmit(ctx, resp.Transaction, skipPreFlight)
}

func (w *WSClient) SubmitManageCollateral(ctx context.Context, req *pb.PostManageCollateralRequest, skipPreFlight bool) (string, error) {
	resp, err := w.PostManageCollateral(ctx, req)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, resp.Transaction, skipPreFlight)
}

func (w *WSClient) SubmitCreateUser(ctx context.Context, req *pb.PostCreateUserRequest, skipPreFlight bool) (interface{}, interface{}) {
	resp, err := w.PostCreateUser(ctx, req)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, resp.Transaction, skipPreFlight)
}

func (w *WSClient) SubmitPostPerpOrder(ctx context.Context, req *pb.PostPerpOrderRequest, skipPreFlight bool) (interface{}, interface{}) {
	resp, err := w.PostPerpOrder(ctx, req)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, resp.Transaction, skipPreFlight)
}

// SubmitPostModifyDriftOrder builds a Drift modify-order txn, signs and submits it to the network.
func (w *WSClient) SubmitPostModifyDriftOrder(ctx context.Context, req *pb.PostModifyDriftOrderRequest, skipPreFlight bool) (string, error) {
	resp, err := w.PostModifyDriftOrder(ctx, req)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, resp.Transaction, skipPreFlight)
}

func (w *WSClient) Close() error {
	return w.conn.Close(errors.New("shutdown requested"))
}

// GetOrderbooksStream subscribes to a stream for changes to the requested market updates (e.g. asks and bids. Set limit to 0 for all bids/ asks).
func (w *WSClient) GetOrderbooksStream(ctx context.Context, markets []string, limit uint32, project pb.Project) (connections.Streamer[*pb.GetOrderbooksStreamResponse], error) {
	return connections.WSStreamProto(w.conn, ctx, "GetOrderbooksStream", &pb.GetOrderbooksRequest{
		Markets: markets,
		Limit:   limit,
		Project: project,
	}, func() *pb.GetOrderbooksStreamResponse {
		var v pb.GetOrderbooksStreamResponse
		return &v
	})
}

// GetDriftMarginOrderbooksStream subscribes to a stream for changes to the requested market updates (e.g. asks and bids. Set limit to 0 for all bids/ asks).
func (w *WSClient) GetDriftMarginOrderbooksStream(ctx context.Context, request *pb.GetDriftMarginOrderbooksRequest) (connections.Streamer[*pb.GetDriftMarginOrderbooksStreamResponse], error) {
	return connections.WSStreamProto(w.conn, ctx, "GetDriftMarginOrderbooksStream", request, func() *pb.GetDriftMarginOrderbooksStreamResponse {
		var v pb.GetDriftMarginOrderbooksStreamResponse
		return &v
	})
}

// GetMarketDepthsStream subscribes to a stream for changes to the requested market data updates (e.g. asks and bids. Set limit to 0 for all bids/ asks).
func (w *WSClient) GetMarketDepthsStream(ctx context.Context, markets []string, limit uint32, project pb.Project) (connections.Streamer[*pb.GetMarketDepthsStreamResponse], error) {
	return connections.WSStreamProto(w.conn, ctx, "GetMarketDepthsStream", &pb.GetMarketDepthsRequest{
		Markets: markets,
		Limit:   limit,
		Project: project,
	}, func() *pb.GetMarketDepthsStreamResponse {
		var v pb.GetMarketDepthsStreamResponse
		return &v
	})
}

// GetTradesStream subscribes to a stream for trades as they execute. Set limit to 0 for all trades.
func (w *WSClient) GetTradesStream(ctx context.Context, market string, limit uint32, project pb.Project) (connections.Streamer[*pb.GetTradesStreamResponse], error) {
	return connections.WSStreamProto(w.conn, ctx, "GetTradesStream", &pb.GetTradesRequest{
		Market:  market,
		Limit:   limit,
		Project: project,
	}, func() *pb.GetTradesStreamResponse {
		var v pb.GetTradesStreamResponse
		return &v
	})
}

// GetNewRaydiumPoolsStream subscribes to a stream for new Raydium Pools when they are created.
func (w *WSClient) GetNewRaydiumPoolsStream(ctx context.Context) (connections.Streamer[*pb.GetNewRaydiumPoolsResponse], error) {
	return connections.WSStreamProto(w.conn, ctx, "GetNewRaydiumPoolsStream", &pb.GetNewRaydiumPoolsRequest{}, func() *pb.GetNewRaydiumPoolsResponse {
		var v pb.GetNewRaydiumPoolsResponse
		return &v
	})
}

// GetOrderStatusStream subscribes to a stream that shows updates to the owner's orders
func (w *WSClient) GetOrderStatusStream(ctx context.Context, market, ownerAddress string, project pb.Project) (connections.Streamer[*pb.GetOrderStatusStreamResponse], error) {
	return connections.WSStreamProto(w.conn, ctx, "GetOrderStatusStream", &pb.GetOrderStatusStreamRequest{
		Market:       market,
		OwnerAddress: ownerAddress,
		Project:      project,
	}, func() *pb.GetOrderStatusStreamResponse {
		var v pb.GetOrderStatusStreamResponse
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
func (w *WSClient) GetPoolReservesStream(ctx context.Context, projects []pb.Project) (connections.Streamer[*pb.GetPoolReservesStreamResponse], error) {
	return connections.WSStreamProto(w.conn, ctx, "GetPoolReservesStream", &pb.GetPoolReservesStreamRequest{
		Projects: projects,
	}, func() *pb.GetPoolReservesStreamResponse {
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

// GetPerpOrderbook returns the current state of perpetual contract orderbook.
func (w *WSClient) GetPerpOrderbook(ctx context.Context, request *pb.GetPerpOrderbookRequest) (*pb.GetPerpOrderbookResponse, error) {
	var response pb.GetPerpOrderbookResponse
	err := w.conn.Request(ctx, "GetPerpOrderbook", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetDriftMarginOrderbook returns the current state of margin contract orderbook.
func (w *WSClient) GetDriftMarginOrderbook(ctx context.Context, request *pb.GetDriftMarginOrderbookRequest) (*pb.GetDriftMarginOrderbookResponse, error) {
	var response pb.GetDriftMarginOrderbookResponse
	err := w.conn.Request(ctx, "GetDriftMarginOrderbook", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetPerpOrderbooksStream subscribes to a stream for perpetual orderbook updates.
func (w *WSClient) GetPerpOrderbooksStream(ctx context.Context, request *pb.GetPerpOrderbooksRequest) (connections.Streamer[*pb.GetPerpOrderbooksStreamResponse], error) {
	newResponse := func() *pb.GetPerpOrderbooksStreamResponse {
		return &pb.GetPerpOrderbooksStreamResponse{}
	}
	return connections.WSStreamProto(w.conn, ctx, "GetPerpOrderbooksStream", request, newResponse)
}

// GetPerpTradesStream subscribes to a stream for trades to the requested contracts
func (w *WSClient) GetPerpTradesStream(ctx context.Context, request *pb.GetPerpTradesStreamRequest) (connections.Streamer[*pb.GetPerpTradesStreamResponse], error) {
	newResponse := func() *pb.GetPerpTradesStreamResponse {
		return &pb.GetPerpTradesStreamResponse{}
	}
	return connections.WSStreamProto(w.conn, ctx, "GetPerpTradesStream", request, newResponse)
}

// GetDriftMarketDepth returns market depth data.
func (w *WSClient) GetDriftMarketDepth(ctx context.Context, request *pb.GetDriftMarketDepthRequest) (*pb.GetDriftMarketDepthResponse, error) {
	var response pb.GetDriftMarketDepthResponse
	err := w.conn.Request(ctx, "GetDriftMarketDepth", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetDriftMarketDepthsStream subscribes to a stream for drift market depth data updates.
func (w *WSClient) GetDriftMarketDepthsStream(ctx context.Context, request *pb.GetDriftMarketDepthsStreamRequest) (connections.Streamer[*pb.GetDriftMarketDepthStreamResponse], error) {
	newResponse := func() *pb.GetDriftMarketDepthStreamResponse {
		return &pb.GetDriftMarketDepthStreamResponse{}
	}
	return connections.WSStreamProto(w.conn, ctx, "GetDriftMarketDepthsStream", request, newResponse)
}

// V2 Openbook

// GetMarketsV2 returns the list of all available named markets
func (w *WSClient) GetMarketsV2(ctx context.Context) (*pb.GetMarketsResponse, error) {
	var response pb.GetMarketsResponse
	err := w.conn.Request(ctx, "GetMarketsV2", &pb.GetMarketsRequestV2{}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetOrderbookV2 returns the requested market's orderbook (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (w *WSClient) GetOrderbookV2(ctx context.Context, market string, limit uint32) (*pb.GetOrderbookResponseV2, error) {
	var response pb.GetOrderbookResponseV2
	err := w.conn.Request(ctx, "GetOrderbookV2", &pb.GetOrderbookRequestV2{Market: market, Limit: limit}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetMarketDepthV2 returns the requested market's coalesced price data (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (w *WSClient) GetMarketDepthV2(ctx context.Context, market string, limit uint32) (*pb.GetMarketDepthResponseV2, error) {
	var response pb.GetMarketDepthResponseV2
	err := w.conn.Request(ctx, "GetMarketDepthV2", &pb.GetMarketDepthRequestV2{Market: market, Limit: limit}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetTickersV2 returns the requested market tickets. Set market to "" for all markets.
func (w *WSClient) GetTickersV2(ctx context.Context, market string) (*pb.GetTickersResponseV2, error) {
	var response pb.GetTickersResponseV2
	err := w.conn.Request(ctx, "GetTickersV2", &pb.GetTickersRequestV2{Market: market}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetOpenOrdersV2 returns all open orders by owner address and market
func (w *WSClient) GetOpenOrdersV2(ctx context.Context, market string, owner string, openOrdersAddress string, orderID string, clientOrderID uint64) (*pb.GetOpenOrdersResponse, error) {
	var response pb.GetOpenOrdersResponse
	err := w.conn.Request(ctx, "GetOpenOrdersV2", &pb.GetOpenOrdersRequestV2{Market: market, Address: owner, OpenOrdersAddress: openOrdersAddress, OrderID: orderID, ClientOrderID: clientOrderID}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetUnsettledV2 returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (w *WSClient) GetUnsettledV2(ctx context.Context, market string, ownerAddress string) (*pb.GetUnsettledResponse, error) {
	var response pb.GetUnsettledResponse
	err := w.conn.Request(ctx, "GetUnsettledV2", &pb.GetUnsettledRequestV2{Market: market, OwnerAddress: ownerAddress}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostOrderV2 returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (w *WSClient) PostOrderV2(ctx context.Context, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	request := &pb.PostOrderRequestV2{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              orderType,
		Amount:            amount,
		Price:             price,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
	}
	var response pb.PostOrderResponse
	err := w.conn.Request(ctx, "PostOrderV2", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitOrderV2 builds a Serum market order, signs it, and submits to the network.
func (w *WSClient) SubmitOrderV2(ctx context.Context, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := w.PostOrderV2(ctx, owner, payer, market, side, orderType, amount, price, opts)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

// PostCancelOrderV2 builds a Serum cancel order.
func (w *WSClient) PostCancelOrderV2(ctx context.Context, request *pb.PostCancelOrderRequestV2) (*pb.PostCancelOrderResponseV2, error) {
	var response pb.PostCancelOrderResponseV2
	err := w.conn.Request(ctx, "PostCancelOrderV2", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitCancelOrderV2 builds a Serum cancel order, signs and submits it to the network.
func (w *WSClient) SubmitCancelOrderV2(ctx context.Context, request *pb.PostCancelOrderRequestV2, skipPreFlight bool) (*pb.PostSubmitBatchResponse, error) {
	order, err := w.PostCancelOrderV2(ctx, request)
	if err != nil {
		return nil, err
	}

	return w.signAndSubmitBatch(ctx, order.Transactions, SubmitOpts{
		SubmitStrategy: pb.SubmitStrategy_P_SUBMIT_ALL,
		SkipPreFlight:  skipPreFlight,
	})
}

// PostSettleV2 returns a partially signed transaction for settling market funds. Typically, you want to use SubmitSettle instead of this.
func (w *WSClient) PostSettleV2(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string) (*pb.PostSettleResponse, error) {
	request := &pb.PostSettleRequestV2{
		OwnerAddress:      owner,
		Market:            market,
		BaseTokenWallet:   baseTokenWallet,
		QuoteTokenWallet:  quoteTokenWallet,
		OpenOrdersAddress: openOrdersAccount,
	}
	var response pb.PostSettleResponse
	err := w.conn.Request(ctx, "PostSettleV2", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitSettleV2 builds a market SubmitSettle transaction, signs it, and submits to the network.
func (w *WSClient) SubmitSettleV2(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, skipPreflight bool) (string, error) {
	order, err := w.PostSettleV2(ctx, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount)
	if err != nil {
		return "", err
	}
	return w.signAndSubmit(ctx, order.Transaction, skipPreflight)
}

func (w *WSClient) PostReplaceOrderV2(ctx context.Context, orderID, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	request := &pb.PostReplaceOrderRequestV2{
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
	}
	var response pb.PostOrderResponse
	err := w.conn.Request(ctx, "PostReplaceOrderV2", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (w *WSClient) SubmitReplaceOrderV2(ctx context.Context, orderID, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := w.PostReplaceOrderV2(ctx, orderID, owner, payer, market, side, orderType, amount, price, opts)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}
