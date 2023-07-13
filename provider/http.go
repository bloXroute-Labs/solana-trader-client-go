package provider

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/connections"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/bloXroute-Labs/solana-trader-proto/common"
	"github.com/gagliardetto/solana-go"
	"net/http"
)

type HTTPClient struct {
	pb.UnimplementedApiServer

	baseURL    string
	httpClient *http.Client
	requestID  utils.RequestID
	privateKey *solana.PrivateKey
	authHeader string
}

// NewHTTPClient connects to Mainnet Trader API
func NewHTTPClient() *HTTPClient {
	opts := DefaultRPCOpts(MainnetVirginiaHTTP)
	return NewHTTPClientWithOpts(nil, opts)
}

// NewHTTPTestnet connects to Testnet Trader API
func NewHTTPTestnet() *HTTPClient {
	opts := DefaultRPCOpts(TestnetHTTP)
	return NewHTTPClientWithOpts(nil, opts)
}

// NewHTTPDevnet connects to Devnet Trader API
func NewHTTPDevnet() *HTTPClient {
	opts := DefaultRPCOpts(DevnetHTTP)
	return NewHTTPClientWithOpts(nil, opts)
}

// NewHTTPLocal connects to local Trader API
func NewHTTPLocal() *HTTPClient {
	opts := DefaultRPCOpts(LocalHTTP)
	return NewHTTPClientWithOpts(nil, opts)
}

// NewHTTPClientWithOpts connects to custom Trader API (set client to nil to use default client)
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

// GetRaydiumPools returns pools on Raydium
func (h *HTTPClient) GetRaydiumPools(ctx context.Context, request *pb.GetRaydiumPoolsRequest) (*pb.GetRaydiumPoolsResponse, error) {
	url := fmt.Sprintf("%s/api/v2/raydium/pools?pairOrAddress=%s", h.baseURL, request.PairOrAddress)
	pools := new(pb.GetRaydiumPoolsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetRaydiumPoolsResponse](ctx, url, h.httpClient, pools, h.authHeader); err != nil {
		return nil, err
	}

	return pools, nil
}

// GetRaydiumQuotes returns the possible amount(s) of outToken for an inToken and the route to achieve it on Raydium
func (h *HTTPClient) GetRaydiumQuotes(ctx context.Context, request *pb.GetRaydiumQuotesRequest) (*pb.GetRaydiumQuotesResponse, error) {
	url := fmt.Sprintf("%s/api/v2/raydium/quotes?inToken=%s&outToken=%s&inAmount=%v&slippage=%v&limit=%v",
		h.baseURL, request.InToken, request.OutToken, request.InAmount, request.Slippage, request.Limit)
	response := new(pb.GetRaydiumQuotesResponse)
	if err := connections.HTTPGetWithClient[*pb.GetRaydiumQuotesResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// GetRaydiumPrices returns the USDC price of requested tokens on Raydium
func (h *HTTPClient) GetRaydiumPrices(ctx context.Context, request *pb.GetRaydiumPricesRequest) (*pb.GetRaydiumPricesResponse, error) {
	tokensArg := convertStrSliceArgument("tokens", true, request.Tokens)
	url := fmt.Sprintf("%s/api/v2/raydium/prices%s", h.baseURL, tokensArg)
	respons := new(pb.GetRaydiumPricesResponse)
	if err := connections.HTTPGetWithClient[*pb.GetRaydiumPricesResponse](ctx, url, h.httpClient, respons, h.authHeader); err != nil {
		return nil, err
	}

	return respons, nil
}

// PostRaydiumSwap returns a partially signed transaction(s) for submitting a swap request on Raydium
func (h *HTTPClient) PostRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest) (*pb.PostRaydiumSwapResponse, error) {
	url := fmt.Sprintf("%s/api/v2/raydium/swap", h.baseURL)
	var response pb.PostRaydiumSwapResponse
	err := connections.HTTPPostWithClient[*pb.PostRaydiumSwapResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// PostRaydiumRouteSwap returns a partially signed transaction(s) for submitting a swap request on Raydium
func (h *HTTPClient) PostRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest) (*pb.PostRaydiumRouteSwapResponse, error) {
	url := fmt.Sprintf("%s/api/v2/raydium/route-swap", h.baseURL)
	var response pb.PostRaydiumRouteSwapResponse
	err := connections.HTTPPostWithClient[*pb.PostRaydiumRouteSwapResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetJupiterQuotes returns the possible amount(s) of outToken for an inToken and the route to achieve it on Jupiter
func (h *HTTPClient) GetJupiterQuotes(ctx context.Context, request *pb.GetJupiterQuotesRequest) (*pb.GetJupiterQuotesResponse, error) {
	url := fmt.Sprintf("%s/api/v2/jupiter/quotes?inToken=%s&outToken=%s&inAmount=%v&slippage=%v&limit=%v",
		h.baseURL, request.InToken, request.OutToken, request.InAmount, request.Slippage, request.Limit)
	response := new(pb.GetJupiterQuotesResponse)
	if err := connections.HTTPGetWithClient[*pb.GetJupiterQuotesResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// GetJupiterPrices returns the USDC price of requested tokens on Jupiter
func (h *HTTPClient) GetJupiterPrices(ctx context.Context, request *pb.GetJupiterPricesRequest) (*pb.GetJupiterPricesResponse, error) {
	tokensArg := convertStrSliceArgument("tokens", true, request.Tokens)
	url := fmt.Sprintf("%s/api/v2/jupiter/prices%s", h.baseURL, tokensArg)
	response := new(pb.GetJupiterPricesResponse)
	if err := connections.HTTPGetWithClient[*pb.GetJupiterPricesResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// PostJupiterSwap returns a partially signed transaction(s) for submitting a swap request on Jupiter
func (h *HTTPClient) PostJupiterSwap(ctx context.Context, request *pb.PostJupiterSwapRequest) (*pb.PostJupiterSwapResponse, error) {
	url := fmt.Sprintf("%s/api/v2/jupiter/swap", h.baseURL)
	var response pb.PostJupiterSwapResponse
	err := connections.HTTPPostWithClient[*pb.PostJupiterSwapResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// PostJupiterRouteSwap returns a partially signed transaction(s) for submitting a swap request on Jupiter
func (h *HTTPClient) PostJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest) (*pb.PostJupiterRouteSwapResponse, error) {
	url := fmt.Sprintf("%s/api/v2/jupiter/route-swap", h.baseURL)
	var response pb.PostJupiterRouteSwapResponse
	err := connections.HTTPPostWithClient[*pb.PostJupiterRouteSwapResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// PostCloseDriftPerpPositions returns a partially signed transaction for canceling perp positions on Drift
func (h *HTTPClient) PostCloseDriftPerpPositions(ctx context.Context, request *pb.PostCloseDriftPerpPositionsRequest) (*pb.PostCloseDriftPerpPositionsResponse, error) {
	url := fmt.Sprintf("%s/api/v2/drift/perp/close", h.baseURL)
	var response pb.PostCloseDriftPerpPositionsResponse
	err := connections.HTTPPostWithClient[*pb.PostCloseDriftPerpPositionsResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetDriftPerpOrderbook returns the current state of perpetual contract orderbook on Drift
func (h *HTTPClient) GetDriftPerpOrderbook(ctx context.Context, request *pb.GetDriftPerpOrderbookRequest) (*pb.GetDriftPerpOrderbookResponse, error) {
	url := fmt.Sprintf("%s/api/v2/drift/perp/orderbook/%s?limit=%d", h.baseURL, request.Contract, request.Limit)
	orderbook := new(pb.GetDriftPerpOrderbookResponse)
	if err := connections.HTTPGetWithClient[*pb.GetDriftPerpOrderbookResponse](ctx, url, h.httpClient, orderbook, h.authHeader); err != nil {
		return nil, err
	}

	return orderbook, nil
}

// PostCreateDriftUser returns a partially signed transaction for creating a user on Drift
func (h *HTTPClient) PostCreateDriftUser(ctx context.Context, request *pb.PostCreateDriftUserRequest) (*pb.PostCreateDriftUserResponse, error) {
	url := fmt.Sprintf("%s/api/v2/drift/user", h.baseURL)
	response := new(pb.PostCreateDriftUserResponse)
	if err := connections.HTTPPostWithClient[*pb.PostCreateDriftUserResponse](ctx, url, h.httpClient, request, response, h.authHeader); err != nil {
		return nil, err
	}
	return response, nil
}

// GetDriftUser returns a user's info on Drift
func (h *HTTPClient) GetDriftUser(ctx context.Context, request *pb.GetDriftUserRequest) (*pb.GetDriftUserResponse, error) {
	url := fmt.Sprintf("%s/api/v2/drift/user?ownerAddress=%s&accountAddress=%s", h.baseURL,
		request.OwnerAddress, request.AccountAddress)
	resp := new(pb.GetDriftUserResponse)
	if err := connections.HTTPGetWithClient[*pb.GetDriftUserResponse](ctx, url, h.httpClient, resp, h.authHeader); err != nil {
		return nil, err
	}

	return resp, nil
}

// PostDriftManageCollateral returns a partially signed transaction for managing collateral on Drift
func (h *HTTPClient) PostDriftManageCollateral(ctx context.Context, request *pb.PostDriftManageCollateralRequest) (*pb.PostDriftManageCollateralResponse, error) {
	url := fmt.Sprintf("%s/api/v2/drift/manage-collateral", h.baseURL)
	response := new(pb.PostDriftManageCollateralResponse)
	if err := connections.HTTPPostWithClient[*pb.PostDriftManageCollateralResponse](ctx, url, h.httpClient, request, response, h.authHeader); err != nil {
		return nil, err
	}
	return response, nil
}

// PostDriftSettlePNL returns partially signed transactions for settling PNL on Drift
func (h *HTTPClient) PostDriftSettlePNL(ctx context.Context, request *pb.PostDriftSettlePNLRequest) (*pb.PostDriftSettlePNLResponse, error) {
	url := fmt.Sprintf("%s/api/v2/drift/perp/settle-pnl", h.baseURL)
	response := new(pb.PostDriftSettlePNLResponse)
	if err := connections.HTTPPostWithClient[*pb.PostDriftSettlePNLResponse](ctx, url, h.httpClient, request, response, h.authHeader); err != nil {
		return nil, err
	}
	return response, nil
}

// PostDriftSettlePNLs returns partially signed transactions for settling PNLs on Drift
func (h *HTTPClient) PostDriftSettlePNLs(ctx context.Context, request *pb.PostDriftSettlePNLsRequest) (*pb.PostDriftSettlePNLsResponse, error) {
	url := fmt.Sprintf("%s/api/v2/drift/perp/settle-pnls", h.baseURL)
	response := new(pb.PostDriftSettlePNLsResponse)
	if err := connections.HTTPPostWithClient[*pb.PostDriftSettlePNLsResponse](ctx, url, h.httpClient, request, response, h.authHeader); err != nil {
		return nil, err
	}
	return response, nil
}

// GetDriftAssets returns list of assets for user on Drift
func (h *HTTPClient) GetDriftAssets(ctx context.Context, request *pb.GetDriftAssetsRequest) (*pb.GetDriftAssetsResponse, error) {
	url := fmt.Sprintf("%s/api/v2/drift/assets?ownerAddress=%s&accountAddress=%s", h.baseURL,
		request.OwnerAddress, request.AccountAddress)
	assets := new(pb.GetDriftAssetsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetDriftAssetsResponse](ctx, url, h.httpClient, assets, h.authHeader); err != nil {
		return nil, err
	}

	return assets, nil
}

// GetDriftPerpContracts returns list of available perp contracts on Drift
func (h *HTTPClient) GetDriftPerpContracts(ctx context.Context, _ *pb.GetDriftPerpContractsRequest) (*pb.GetDriftPerpContractsResponse, error) {
	url := fmt.Sprintf("%s/api/v2/drift/perp/contracts", h.baseURL)
	positions := new(pb.GetDriftPerpContractsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetDriftPerpContractsResponse](ctx, url, h.httpClient, positions, h.authHeader); err != nil {
		return nil, err
	}

	return positions, nil
}

// PostLiquidateDriftPerp returns a partially signed transaction for liquidating perp position on Drift
func (h *HTTPClient) PostLiquidateDriftPerp(ctx context.Context, request *pb.PostLiquidateDriftPerpRequest) (*pb.PostLiquidateDriftPerpResponse, error) {
	url := fmt.Sprintf("%s/api/v2/drift/perp/liquidate", h.baseURL)
	response := new(pb.PostLiquidateDriftPerpResponse)
	if err := connections.HTTPPostWithClient[*pb.PostLiquidateDriftPerpResponse](ctx, url, h.httpClient, request, response, h.authHeader); err != nil {
		return nil, err
	}
	return response, nil
}

// GetDriftOpenPerpOrder returns an open perp order on Drift
func (h *HTTPClient) GetDriftOpenPerpOrder(ctx context.Context, request *pb.GetDriftOpenPerpOrderRequest) (*pb.GetDriftOpenPerpOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v2/drift/perp/open-order?ownerAddress=%s&accountAddress=%s&clientOrderID=%d&orderID=%d", h.baseURL,
		request.OwnerAddress, request.AccountAddress, request.ClientOrderID, request.OrderID)
	order := new(pb.GetDriftOpenPerpOrderResponse)
	if err := connections.HTTPGetWithClient[*pb.GetDriftOpenPerpOrderResponse](ctx, url, h.httpClient, order, h.authHeader); err != nil {
		return nil, err
	}

	return order, nil
}

// GetDriftOpenMarginOrder return a open margin order on Drift
func (h *HTTPClient) GetDriftOpenMarginOrder(ctx context.Context, request *pb.GetDriftOpenMarginOrderRequest) (*pb.GetDriftOpenMarginOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v2/drift/margin/open-order?ownerAddress=%s&accountAddress=%s&clientOrderID=%d&orderID=%d",
		h.baseURL, request.OwnerAddress, request.AccountAddress, request.GetClientOrderID(), request.GetOrderID())
	order := new(pb.GetDriftOpenMarginOrderResponse)
	if err := connections.HTTPGetWithClient[*pb.GetDriftOpenMarginOrderResponse](ctx, url, h.httpClient, order, h.authHeader); err != nil {
		return nil, err
	}

	return order, nil
}

// GetDriftPerpPositions returns all perp positions on Drift
func (h *HTTPClient) GetDriftPerpPositions(ctx context.Context, request *pb.GetDriftPerpPositionsRequest) (*pb.GetDriftPerpPositionsResponse, error) {
	contractsString := convertStrSliceArgument("contracts", false, request.Contracts)
	url := fmt.Sprintf("%s/api/v2/drift/perp/positions?ownerAddress=%s&accountAddress=%s%s",
		h.baseURL, request.OwnerAddress, request.AccountAddress, contractsString)
	response := new(pb.GetDriftPerpPositionsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetDriftPerpPositionsResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// GetDriftOpenPerpOrders returns all open perp orders on Drift
func (h *HTTPClient) GetDriftOpenPerpOrders(ctx context.Context, request *pb.GetDriftOpenPerpOrdersRequest) (*pb.GetDriftOpenPerpOrdersResponse, error) {
	contractsString := convertStrSliceArgument("contracts", false, request.Contracts)
	url := fmt.Sprintf("%s/api/v2/drift/perp/open-orders?ownerAddress=%s&accountAddress=%s%s",
		h.baseURL, request.OwnerAddress, request.AccountAddress, contractsString)
	response := new(pb.GetDriftOpenPerpOrdersResponse)
	if err := connections.HTTPGetWithClient[*pb.GetDriftOpenPerpOrdersResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// PostDriftCancelPerpOrder returns a partially signed transaction for canceling Drift perp order(s)
func (h *HTTPClient) PostDriftCancelPerpOrder(ctx context.Context, request *pb.PostDriftCancelPerpOrderRequest) (*pb.PostDriftCancelPerpOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v2/drift/perp/cancel", h.baseURL)
	response := new(pb.PostDriftCancelPerpOrderResponse)
	if err := connections.HTTPPostWithClient[*pb.PostDriftCancelPerpOrderResponse](ctx, url, h.httpClient, request, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// GetOrderbook returns the requested market's orderbook (e.h. asks and bids). Set limit to 0 for all bids / asks.
func (h *HTTPClient) GetOrderbook(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetOrderbookResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/orderbooks/%s?limit=%v&project=%v", h.baseURL, market, limit, project)
	orderbook := new(pb.GetOrderbookResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOrderbookResponse](ctx, url, h.httpClient, orderbook, h.authHeader); err != nil {
		return nil, err
	}

	return orderbook, nil
}

// GetMarketDepth returns the requested market's coalesced price data (e.h. asks and bids). Set limit to 0 for all bids / asks.
func (h *HTTPClient) GetMarketDepth(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetMarketDepthResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/depth/%s?limit=%v&project=%v", h.baseURL, market, limit, project)
	mktDepth := new(pb.GetMarketDepthResponse)
	if err := connections.HTTPGetWithClient[*pb.GetMarketDepthResponse](ctx, url, h.httpClient, mktDepth, h.authHeader); err != nil {
		return nil, err
	}

	return mktDepth, nil
}

// GetTrades returns the requested market's currently executing trades. Set limit to 0 for all trades.
func (h *HTTPClient) GetTrades(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetTradesResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/trades/%s?limit=%v&project=%v", h.baseURL, market, limit, project)
	marketTrades := new(pb.GetTradesResponse)
	if err := connections.HTTPGetWithClient[*pb.GetTradesResponse](ctx, url, h.httpClient, marketTrades, h.authHeader); err != nil {
		return nil, err
	}

	return marketTrades, nil
}

// GetPools returns pools for given projects.
func (h *HTTPClient) GetPools(ctx context.Context, projects []pb.Project) (*pb.GetPoolsResponse, error) {
	projectsArg := convertSliceArgument("projects", true, projects)
	url := fmt.Sprintf("%s/api/v1/market/pools%s", h.baseURL, projectsArg)
	pools := new(pb.GetPoolsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetPoolsResponse](ctx, url, h.httpClient, pools, h.authHeader); err != nil {
		return nil, err
	}

	return pools, nil
}

// GetTickers returns the requested market tickets. Set market to "" for all markets.
func (h *HTTPClient) GetTickers(ctx context.Context, market string, project pb.Project) (*pb.GetTickersResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/tickers/%s?project=%v", h.baseURL, market, project)
	tickers := new(pb.GetTickersResponse)
	if err := connections.HTTPGetWithClient[*pb.GetTickersResponse](ctx, url, h.httpClient, tickers, h.authHeader); err != nil {
		return nil, err
	}

	return tickers, nil
}

// GetOpenOrders returns all open orders by owner address and market
func (h *HTTPClient) GetOpenOrders(ctx context.Context, market string, owner string, openOrdersAddress string, project pb.Project) (*pb.GetOpenOrdersResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/openorders/%s?address=%s&openOrdersAddress=%s&project=%s", h.baseURL, market, owner, openOrdersAddress, project)
	orders := new(pb.GetOpenOrdersResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOpenOrdersResponse](ctx, url, h.httpClient, orders, h.authHeader); err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOrderByID returns an order by id
func (h *HTTPClient) GetOrderByID(ctx context.Context, in *pb.GetOrderByIDRequest) (*pb.GetOrderByIDResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/orderbyid/%s?market=%s&project=%s", h.baseURL, in.OrderID, in.Market, in.Project)
	orders := new(pb.GetOrderByIDResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOrderByIDResponse](ctx, url, h.httpClient, orders, h.authHeader); err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOpenPerpOrders returns all open perp orders
func (h *HTTPClient) GetOpenPerpOrders(ctx context.Context, request *pb.GetOpenPerpOrdersRequest) (*pb.GetOpenPerpOrdersResponse, error) {
	contractsString := convertSliceArgument("contracts", false, request.Contracts)
	url := fmt.Sprintf("%s/api/v1/trade/perp/open-orders?ownerAddress=%s&accountAddress=%s&project=%s%s",
		h.baseURL, request.OwnerAddress, request.AccountAddress, request.Project, contractsString)
	orders := new(pb.GetOpenPerpOrdersResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOpenPerpOrdersResponse](ctx, url, h.httpClient, orders, h.authHeader); err != nil {
		return nil, err
	}

	return orders, nil
}

// GetDriftOpenMarginOrders returns all open margin orders on Drift
func (h *HTTPClient) GetDriftOpenMarginOrders(ctx context.Context, request *pb.GetDriftOpenMarginOrdersRequest) (*pb.GetDriftOpenMarginOrdersResponse, error) {
	marketsString := convertStrSliceArgument("markets", false, request.Markets)
	url := fmt.Sprintf("%s/api/v2/drift/margin/open-orders?ownerAddress=%s&accountAddress=%s%s",
		h.baseURL, request.OwnerAddress, request.AccountAddress, marketsString)
	orders := new(pb.GetDriftOpenMarginOrdersResponse)
	if err := connections.HTTPGetWithClient[*pb.GetDriftOpenMarginOrdersResponse](ctx, url, h.httpClient, orders, h.authHeader); err != nil {
		return nil, err
	}

	return orders, nil
}

// PostCancelPerpOrder returns a partially signed transaction for canceling perp order
func (h *HTTPClient) PostCancelPerpOrder(ctx context.Context, request *pb.PostCancelPerpOrderRequest) (*pb.PostCancelPerpOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/perp/cancelbyid", h.baseURL)
	response := new(pb.PostCancelPerpOrderResponse)
	if err := connections.HTTPPostWithClient[*pb.PostCancelPerpOrderResponse](ctx, url, h.httpClient, request, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// PostCancelDriftMarginOrder returns a partially signed transaction for canceling margin orders on Drift
func (h *HTTPClient) PostCancelDriftMarginOrder(ctx context.Context, request *pb.PostCancelDriftMarginOrderRequest) (*pb.PostCancelDriftMarginOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v2/drift/margin/cancel", h.baseURL)
	response := new(pb.PostCancelDriftMarginOrderResponse)
	if err := connections.HTTPPostWithClient[*pb.PostCancelDriftMarginOrderResponse](ctx, url, h.httpClient, request, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// PostCancelPerpOrders returns a partially signed transaction for canceling all perp orders of a user
func (h *HTTPClient) PostCancelPerpOrders(ctx context.Context, request *pb.PostCancelPerpOrdersRequest) (*pb.PostCancelPerpOrdersResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/perp/cancel", h.baseURL)
	response := new(pb.PostCancelPerpOrdersResponse)
	if err := connections.HTTPPostWithClient[*pb.PostCancelPerpOrdersResponse](ctx, url, h.httpClient, request, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// PostCreateUser returns a partially signed transaction for creating a user
func (h *HTTPClient) PostCreateUser(ctx context.Context, request *pb.PostCreateUserRequest) (*pb.PostCreateUserResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/user", h.baseURL)
	response := new(pb.PostCreateUserResponse)
	if err := connections.HTTPPostWithClient[*pb.PostCreateUserResponse](ctx, url, h.httpClient, request, response, h.authHeader); err != nil {
		return nil, err
	}
	return response, nil
}

// GetUser returns a user's info
func (h *HTTPClient) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/user?ownerAddress=%s&accountAddress=%s&project=%s", h.baseURL,
		request.OwnerAddress, request.AccountAddress, request.Project)
	resp := new(pb.GetUserResponse)
	if err := connections.HTTPGetWithClient[*pb.GetUserResponse](ctx, url, h.httpClient, resp, h.authHeader); err != nil {
		return nil, err
	}

	return resp, nil
}

// PostManageCollateral returns a partially signed transaction for managing collateral
func (h *HTTPClient) PostManageCollateral(ctx context.Context, request *pb.PostManageCollateralRequest) (*pb.PostManageCollateralResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/perp/managecollateral", h.baseURL)
	response := new(pb.PostManageCollateralResponse)
	if err := connections.HTTPPostWithClient[*pb.PostManageCollateralResponse](ctx, url, h.httpClient, request, response, h.authHeader); err != nil {
		return nil, err
	}
	return response, nil
}

// PostSettlePNL returns a partially signed transaction for settling PNL
func (h *HTTPClient) PostSettlePNL(ctx context.Context, request *pb.PostSettlePNLRequest) (*pb.PostSettlePNLResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/perp/settle-pnl", h.baseURL)
	response := new(pb.PostSettlePNLResponse)
	if err := connections.HTTPPostWithClient[*pb.PostSettlePNLResponse](ctx, url, h.httpClient, request, response, h.authHeader); err != nil {
		return nil, err
	}
	return response, nil
}

// PostSettlePNLs returns partially signed transactions for settling PNLs
func (h *HTTPClient) PostSettlePNLs(ctx context.Context, request *pb.PostSettlePNLsRequest) (*pb.PostSettlePNLsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/perp/settle-pnls", h.baseURL)
	response := new(pb.PostSettlePNLsResponse)
	if err := connections.HTTPPostWithClient[*pb.PostSettlePNLsResponse](ctx, url, h.httpClient, request, response, h.authHeader); err != nil {
		return nil, err
	}
	return response, nil
}

// GetAssets returns list of assets for user
func (h *HTTPClient) GetAssets(ctx context.Context, request *pb.GetAssetsRequest) (*pb.GetAssetsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/perp/assets?ownerAddress=%s&accountAddress=%s&project=%s", h.baseURL,
		request.OwnerAddress, request.AccountAddress, request.Project)
	assets := new(pb.GetAssetsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetAssetsResponse](ctx, url, h.httpClient, assets, h.authHeader); err != nil {
		return nil, err
	}

	return assets, nil
}

// GetPerpContracts returns list of available perp contracts
func (h *HTTPClient) GetPerpContracts(ctx context.Context, request *pb.GetPerpContractsRequest) (*pb.GetPerpContractsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/perp/contracts?project=%s", h.baseURL, request.Project)
	positions := new(pb.GetPerpContractsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetPerpContractsResponse](ctx, url, h.httpClient, positions, h.authHeader); err != nil {
		return nil, err
	}

	return positions, nil
}

// PostLiquidatePerp returns a partially signed transaction for liquidating perp position
func (h *HTTPClient) PostLiquidatePerp(ctx context.Context, request *pb.PostLiquidatePerpRequest) (*pb.PostLiquidatePerpResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/perp/liquidate", h.baseURL)
	response := new(pb.PostLiquidatePerpResponse)
	if err := connections.HTTPPostWithClient[*pb.PostLiquidatePerpResponse](ctx, url, h.httpClient, request, response, h.authHeader); err != nil {
		return nil, err
	}
	return response, nil
}

// GetOpenPerpOrder returns an open perp order
func (h *HTTPClient) GetOpenPerpOrder(ctx context.Context, request *pb.GetOpenPerpOrderRequest) (*pb.GetOpenPerpOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/perp/open-order-by-id?ownerAddress=%s&accountAddress=%s&project=%s&clientOrderID=%d&orderID=%d", h.baseURL,
		request.OwnerAddress, request.AccountAddress, request.Project, request.ClientOrderID, request.OrderID)
	order := new(pb.GetOpenPerpOrderResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOpenPerpOrderResponse](ctx, url, h.httpClient, order, h.authHeader); err != nil {
		return nil, err
	}

	return order, nil
}

// GetPerpPositions returns all perp positions by owner address and market
func (h *HTTPClient) GetPerpPositions(ctx context.Context, request *pb.GetPerpPositionsRequest) (*pb.GetPerpPositionsResponse, error) {
	var strs []string
	for _, c := range request.Contracts {
		strs = append(strs, fmt.Sprint(c))
	}

	contractsArg := convertStrSliceArgument("contracts", false, strs)
	url := fmt.Sprintf("%s/api/v1/trade/perp/positions?ownerAddress=%s&accountAddress=%s&project=%s%s", h.baseURL,
		request.OwnerAddress, request.AccountAddress, request.Project, contractsArg)
	positions := new(pb.GetPerpPositionsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetPerpPositionsResponse](ctx, url, h.httpClient, positions, h.authHeader); err != nil {
		return nil, err
	}

	return positions, nil
}

// GetMarkets returns the list of all available named markets
func (h *HTTPClient) GetMarkets(ctx context.Context) (*pb.GetMarketsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/markets", h.baseURL)
	markets := new(pb.GetMarketsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetMarketsResponse](ctx, url, h.httpClient, markets, h.authHeader); err != nil {
		return nil, err
	}

	return markets, nil
}

// GetDriftMarkets returns the list of all available named markets
func (h *HTTPClient) GetDriftMarkets(ctx context.Context, request *pb.GetDriftMarketsRequest) (*pb.GetDriftMarketsResponse, error) {
	url := fmt.Sprintf("%s/api/v2/drift/markets?metadata=%v", h.baseURL, request.Metadata)
	markets := new(pb.GetDriftMarketsResponse)
	if err := connections.HTTPGetWithClient[*pb.GetDriftMarketsResponse](ctx, url, h.httpClient, markets, h.authHeader); err != nil {
		return nil, err
	}

	return markets, nil
}

// GetUnsettled returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (h *HTTPClient) GetUnsettled(ctx context.Context, market string, owner string, project pb.Project) (*pb.GetUnsettledResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/unsettled/%s?ownerAddress=%s&project=%s", h.baseURL, market, owner, project)
	result := new(pb.GetUnsettledResponse)
	if err := connections.HTTPGetWithClient[*pb.GetUnsettledResponse](ctx, url, h.httpClient, result, h.authHeader); err != nil {
		return nil, err
	}

	return result, nil
}

// GetAccountBalance returns all OpenOrders accounts for a given market with the amountsctx context.Context,  of unsettled funds
func (h *HTTPClient) GetAccountBalance(ctx context.Context, owner string) (*pb.GetAccountBalanceResponse, error) {
	url := fmt.Sprintf("%s/api/v1/account/balance?ownerAddress=%s", h.baseURL, owner)
	result := new(pb.GetAccountBalanceResponse)
	if err := connections.HTTPGetWithClient[*pb.GetAccountBalanceResponse](ctx, url, h.httpClient, result, h.authHeader); err != nil {
		return nil, err
	}

	return result, nil
}

// GetPrice returns the USDC price of requested tokens
func (h *HTTPClient) GetPrice(ctx context.Context, tokens []string) (*pb.GetPriceResponse, error) {
	tokensArg := convertStrSliceArgument("tokens", true, tokens)
	url := fmt.Sprintf("%s/api/v1/market/price%s", h.baseURL, tokensArg)
	pools := new(pb.GetPriceResponse)
	if err := connections.HTTPGetWithClient[*pb.GetPriceResponse](ctx, url, h.httpClient, pools, h.authHeader); err != nil {
		return nil, err
	}

	return pools, nil
}

// GetQuotes returns the possible amount(s) of outToken for an inToken and the route to achieve it
func (h *HTTPClient) GetQuotes(ctx context.Context, inToken, outToken string, inAmount, slippage float64, limit int32, projects []pb.Project) (*pb.GetQuotesResponse, error) {
	projectString := convertSliceArgument("projects", false, projects)

	url := fmt.Sprintf("%s/api/v1/market/quote?inToken=%s&outToken=%s&inAmount=%v&slippage=%v&limit=%v%s",
		h.baseURL, inToken, outToken, inAmount, slippage, limit, projectString)
	result := new(pb.GetQuotesResponse)
	if err := connections.HTTPGetWithClient[*pb.GetQuotesResponse](ctx, url, h.httpClient, result, h.authHeader); err != nil {
		return nil, err
	}

	return result, nil
}

// PostSubmit posts the transaction string to the Solana network.
func (h *HTTPClient) PostSubmit(ctx context.Context, txBase64 string, skipPreFlight bool) (*pb.PostSubmitResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/submit", h.baseURL)
	request := &pb.PostSubmitRequest{Transaction: &pb.TransactionMessage{Content: txBase64}, SkipPreFlight: skipPreFlight}

	var response pb.PostSubmitResponse
	err := connections.HTTPPostWithClient[*pb.PostSubmitResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostSubmitBatch posts a bundle of transactions string based on a specific SubmitStrategy to the Solana network.
func (h *HTTPClient) PostSubmitBatch(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/submit-batch", h.baseURL)

	var response pb.PostSubmitBatchResponse
	err := connections.HTTPPostWithClient[*pb.PostSubmitBatchResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SignAndSubmit signs the given transaction and submits it.
func (h *HTTPClient) SignAndSubmit(ctx context.Context, tx *pb.TransactionMessage, skipPreFlight bool) (string, error) {
	if h.privateKey == nil {
		return "", ErrPrivateKeyNotFound
	}
	txBase64, err := transaction.SignTxWithPrivateKey(tx.Content, *h.privateKey)
	if err != nil {
		return "", err
	}

	response, err := h.PostSubmit(ctx, txBase64, skipPreFlight)
	if err != nil {
		return "", err
	}

	return response.Signature, nil
}

// SignAndSubmitBatch signs the given transactions and submits them.
func (h *HTTPClient) SignAndSubmitBatch(ctx context.Context, transactions []*pb.TransactionMessage, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	if h.privateKey == nil {
		return nil, ErrPrivateKeyNotFound
	}
	batchRequest, err := buildBatchRequest(transactions, *h.privateKey, opts)
	if err != nil {
		return nil, err
	}
	return h.PostSubmitBatch(ctx, batchRequest)
}

// PostTradeSwap returns a partially signed transaction for submitting a swap request
func (h *HTTPClient) PostTradeSwap(ctx context.Context, ownerAddress, inToken, outToken string, inAmount, slippage float64, project pb.Project) (*pb.TradeSwapResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/swap", h.baseURL)
	request := &pb.TradeSwapRequest{
		OwnerAddress: ownerAddress,
		InToken:      inToken,
		OutToken:     outToken,
		InAmount:     inAmount,
		Slippage:     slippage,
		Project:      project,
	}

	var response pb.TradeSwapResponse
	err := connections.HTTPPostWithClient[*pb.TradeSwapResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitTradeSwap builds a TradeSwap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitTradeSwap(ctx context.Context, owner, inToken, outToken string, inAmount, slippage float64, project pb.Project, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := h.PostTradeSwap(ctx, owner, inToken, outToken, inAmount, slippage, project)
	if err != nil {
		return nil, err
	}
	return h.SignAndSubmitBatch(ctx, resp.Transactions, opts)
}

// PostRouteTradeSwap returns a partially signed transaction(s) for submitting a route swap request
func (h *HTTPClient) PostRouteTradeSwap(ctx context.Context, request *pb.RouteTradeSwapRequest) (*pb.TradeSwapResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/route-swap", h.baseURL)

	var response pb.TradeSwapResponse
	err := connections.HTTPPostWithClient[*pb.TradeSwapResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitRouteTradeSwap builds a RouteTradeSwap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitRouteTradeSwap(ctx context.Context, request *pb.RouteTradeSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := h.PostRouteTradeSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return h.SignAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitRaydiumSwap builds a Raydium Swap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := h.PostRaydiumSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return h.SignAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitRaydiumRouteSwap builds a Raydium RouteSwap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := h.PostRaydiumRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return h.SignAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitJupiterSwap builds a Jupiter Swap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitJupiterSwap(ctx context.Context, request *pb.PostJupiterSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := h.PostJupiterSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return h.SignAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitJupiterRouteSwap builds a Jupiter RouteSwap transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := h.PostJupiterRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return h.SignAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitPostSettlePNL builds a settle-pnl txn, signs and submits it to the network.
func (h *HTTPClient) SubmitPostSettlePNL(ctx context.Context, request *pb.PostSettlePNLRequest, skipPreFlight bool) (string, error) {
	resp, err := h.PostSettlePNL(ctx, request)
	if err != nil {
		return "", err
	}
	return h.SignAndSubmit(ctx, resp.Transaction, skipPreFlight)
}

// SubmitPostSettlePNLs builds one or many settle-pnl txn, signs and submits them to the network.
func (h *HTTPClient) SubmitPostSettlePNLs(ctx context.Context, request *pb.PostSettlePNLsRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := h.PostSettlePNLs(ctx, request)
	if err != nil {
		return nil, err
	}
	return h.SignAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitPostLiquidatePerp builds a liquidate-perp txn, signs and submits it to the network.
func (h *HTTPClient) SubmitPostLiquidatePerp(ctx context.Context, request *pb.PostLiquidatePerpRequest, skipPreFlight bool) (string, error) {
	resp, err := h.PostLiquidatePerp(ctx, request)
	if err != nil {
		return "", err
	}
	return h.SignAndSubmit(ctx, resp.Transaction, skipPreFlight)
}

// SubmitManageCollateral builds a deposit collateral transaction then signs it, and submits to the network.
func (h *HTTPClient) SubmitManageCollateral(ctx context.Context, request *pb.PostManageCollateralRequest, skipPreFlight bool) (string, error) {
	resp, err := h.PostManageCollateral(ctx, request)
	if err != nil {
		return "", err
	}
	return h.SignAndSubmit(ctx, resp.Transaction, skipPreFlight)
}

// PostOrder returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (h *HTTPClient) PostOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/place", h.baseURL)
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
	err := connections.HTTPPostWithClient[*pb.PostOrderResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostModifyDriftOrder returns a partially signed transaction for modifying a Drift order. Typically, you want to use SubmitPostModifyDriftOrder instead of this.
func (h *HTTPClient) PostModifyDriftOrder(ctx context.Context, request *pb.PostModifyDriftOrderRequest) (*pb.PostModifyDriftOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v2/drift/modify-order", h.baseURL)

	var response pb.PostModifyDriftOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostModifyDriftOrderResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostPerpOrder returns a partially signed transaction for placing a perp order. Typically, you want to use SubmitPerpOrder instead of this.
func (h *HTTPClient) PostPerpOrder(ctx context.Context, request *pb.PostPerpOrderRequest) (*pb.PostPerpOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/perp/order", h.baseURL)

	var response pb.PostPerpOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostPerpOrderResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitPerpOrder builds a perp order, signs it, and submits to the network.
func (h *HTTPClient) SubmitPerpOrder(ctx context.Context, request *pb.PostPerpOrderRequest, opts PostOrderOpts) (string, error) {
	order, err := h.PostPerpOrder(ctx, request)
	if err != nil {
		return "", err
	}

	sig, err := h.SignAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
	return sig, err
}

// SubmitDriftMarginOrder builds a margin order, signs it, and submits to the network.
func (h *HTTPClient) SubmitDriftMarginOrder(ctx context.Context, request *pb.PostDriftMarginOrderRequest, opts PostOrderOpts) (string, error) {
	order, err := h.PostDriftMarginOrder(ctx, request)
	if err != nil {
		return "", err
	}

	sig, err := h.SignAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
	return sig, err
}

// PostDriftMarginOrder returns a partially signed transaction for placing a margin order. Typically, you want to use SubmitDriftMarginOrder instead of this.
func (h *HTTPClient) PostDriftMarginOrder(ctx context.Context, request *pb.PostDriftMarginOrderRequest) (*pb.PostDriftMarginOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v2/drift/margin/place", h.baseURL)

	var response pb.PostDriftMarginOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostDriftMarginOrderResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostDriftEnableMarginTrading returns a partially signed transaction for enabling/disabling margin trading.
func (h *HTTPClient) PostDriftEnableMarginTrading(ctx context.Context, request *pb.PostDriftEnableMarginTradingRequest) (*pb.PostDriftEnableMarginTradingResponse, error) {
	url := fmt.Sprintf("%s/api/v2/drift/margin-enable", h.baseURL)

	var response pb.PostDriftEnableMarginTradingResponse
	err := connections.HTTPPostWithClient[*pb.PostDriftEnableMarginTradingResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitDriftEnableMarginTrading builds a perp order, signs it, and submits to the network.
func (h *HTTPClient) SubmitDriftEnableMarginTrading(ctx context.Context, request *pb.PostDriftEnableMarginTradingRequest, skipPreFlight bool) (string, error) {
	tx, err := h.PostDriftEnableMarginTrading(ctx, request)
	if err != nil {
		return "", err
	}

	return h.SignAndSubmit(ctx, tx.Transaction, skipPreFlight)
}

// SubmitOrder builds a Serum market order, signs it, and submits to the network.
func (h *HTTPClient) SubmitOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error) {
	order, err := h.PostOrder(ctx, owner, payer, market, side, types, amount, price, project, opts)
	if err != nil {
		return "", err
	}

	sig, err := h.SignAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
	return sig, err
}

// PostCancelOrder builds a Serum cancel order.
func (h *HTTPClient) PostCancelOrder(
	ctx context.Context,
	orderID string,
	side pb.Side,
	owner,
	market,
	openOrders string,
	project pb.Project,
) (*pb.PostCancelOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/cancel", h.baseURL)
	request := &pb.PostCancelOrderRequest{
		OrderID:           orderID,
		Side:              side,
		OwnerAddress:      owner,
		MarketAddress:     market,
		Project:           project,
		OpenOrdersAddress: openOrders,
	}

	var response pb.PostCancelOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostCancelOrderResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// SubmitCancelOrder builds a Serum cancel order, signs and submits it to the network.
func (h *HTTPClient) SubmitCancelOrder(
	ctx context.Context,
	orderID string,
	side pb.Side,
	owner,
	market,
	openOrders string,
	project pb.Project,
	skipPreFlight bool,
) (string, error) {
	order, err := h.PostCancelOrder(ctx, orderID, side, owner, market, openOrders, project)
	if err != nil {
		return "", err
	}

	return h.SignAndSubmit(ctx, order.Transaction, skipPreFlight)
}

// PostClosePerpPositions builds cancel perp positions txn.
func (h *HTTPClient) PostClosePerpPositions(ctx context.Context, request *pb.PostClosePerpPositionsRequest) (*pb.PostClosePerpPositionsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/perp/close", h.baseURL)
	var response pb.PostClosePerpPositionsResponse
	err := connections.HTTPPostWithClient[*pb.PostClosePerpPositionsResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// SubmitClosePerpPositions builds a close perp positions txn, signs and submits it to the network.
func (h *HTTPClient) SubmitClosePerpPositions(ctx context.Context, request *pb.PostClosePerpPositionsRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	order, err := h.PostClosePerpPositions(ctx, request)
	if err != nil {
		return nil, err
	}

	var msgs []*pb.TransactionMessage
	for _, txn := range order.Transactions {
		msgs = append(msgs, txn)
	}

	return h.SignAndSubmitBatch(ctx, msgs, opts)
}

// SubmitCancelPerpOrder builds a cancel perp order txn, signs and submits it to the network.
func (h *HTTPClient) SubmitCancelPerpOrder(ctx context.Context, request *pb.PostCancelPerpOrderRequest, skipPreFlight bool) (string, error) {
	order, err := h.PostCancelPerpOrder(ctx, request)
	if err != nil {
		return "", err
	}

	return h.SignAndSubmit(ctx, order.Transaction, skipPreFlight)
}

// SubmitCancelDriftMarginOrder builds a cancel Drift margin order txn, signs and submits it to the network.
func (h *HTTPClient) SubmitCancelDriftMarginOrder(ctx context.Context, request *pb.PostCancelDriftMarginOrderRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	order, err := h.PostCancelDriftMarginOrder(ctx, request)
	if err != nil {
		return nil, err
	}

	var msgs []*pb.TransactionMessage
	for _, txn := range order.Transactions {
		msgs = append(msgs, txn)
	}

	return h.SignAndSubmitBatch(ctx, msgs, opts)
}

// SubmitCancelPerpOrders builds a cancel perp orders txn, signs and submits it to the network.
func (h *HTTPClient) SubmitCancelPerpOrders(ctx context.Context, request *pb.PostCancelPerpOrdersRequest, skipPreFlight bool) (*pb.PostSubmitBatchResponse, error) {
	resp, err := h.PostCancelPerpOrders(ctx, request)
	if err != nil {
		return nil, err
	}
	return h.SignAndSubmitBatch(ctx, resp.Transactions, SubmitOpts{
		SkipPreFlight: skipPreFlight,
	})
}

// SubmitCreateUser builds a create-user txn, signs and submits it to the network.
func (h *HTTPClient) SubmitCreateUser(ctx context.Context, request *pb.PostCreateUserRequest, skipPreFlight bool) (string, error) {
	order, err := h.PostCreateUser(ctx, request)
	if err != nil {
		return "", err
	}

	return h.SignAndSubmit(ctx, order.Transaction, skipPreFlight)
}

// SubmitPostPerpOrder builds a post order txn, signs and submits it to the network.
func (h *HTTPClient) SubmitPostPerpOrder(ctx context.Context, request *pb.PostPerpOrderRequest, skipPreFlight bool) (string, error) {
	order, err := h.PostPerpOrder(ctx, request)
	if err != nil {
		return "", err
	}

	return h.SignAndSubmit(ctx, order.Transaction, skipPreFlight)
}

// SubmitPostModifyDriftOrder builds a Drift modify-order txn, signs and submits it to the network.
func (h *HTTPClient) SubmitPostModifyDriftOrder(ctx context.Context, request *pb.PostModifyDriftOrderRequest, skipPreFlight bool) (string, error) {
	order, err := h.PostModifyDriftOrder(ctx, request)
	if err != nil {
		return "", err
	}

	return h.SignAndSubmit(ctx, order.Transaction, skipPreFlight)
}

// SubmitPostDriftMarginOrder builds a margin order txn, signs and submits it to the network.
func (h *HTTPClient) SubmitPostDriftMarginOrder(ctx context.Context, request *pb.PostDriftMarginOrderRequest, skipPreFlight bool) (string, error) {
	order, err := h.PostDriftMarginOrder(ctx, request)
	if err != nil {
		return "", err
	}

	return h.SignAndSubmit(ctx, order.Transaction, skipPreFlight)
}

// PostCancelByClientOrderID builds a Serum cancel order by client ID.
func (h *HTTPClient) PostCancelByClientOrderID(
	ctx context.Context,
	clientOrderID uint64,
	owner,
	market,
	openOrders string,
	project pb.Project,
) (*pb.PostCancelOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/cancelbyid", h.baseURL)
	request := &pb.PostCancelByClientOrderIDRequest{
		ClientOrderID:     clientOrderID,
		OwnerAddress:      owner,
		MarketAddress:     market,
		Project:           project,
		OpenOrdersAddress: openOrders,
	}

	var response pb.PostCancelOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostCancelOrderResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// SubmitCancelByClientOrderID builds a Serum cancel order by client ID, signs and submits it to the network.
func (h *HTTPClient) SubmitCancelByClientOrderID(
	ctx context.Context,
	clientOrderID uint64,
	owner,
	market,
	openOrders string,
	project pb.Project,
	skipPreFlight bool,
) (string, error) {
	order, err := h.PostCancelByClientOrderID(ctx, clientOrderID, owner, market, openOrders, project)
	if err != nil {
		return "", err
	}

	return h.SignAndSubmit(ctx, order.Transaction, skipPreFlight)
}

func (h *HTTPClient) PostCancelAll(ctx context.Context, market, owner string, openOrdersAddresses []string, project pb.Project) (*pb.PostCancelAllResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/cancelall", h.baseURL)
	request := &pb.PostCancelAllRequest{
		Market:              market,
		OwnerAddress:        owner,
		OpenOrdersAddresses: openOrdersAddresses,
		Project:             project,
	}

	var response pb.PostCancelAllResponse
	err := connections.HTTPPostWithClient[*pb.PostCancelAllResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (h *HTTPClient) SubmitCancelAll(ctx context.Context, market, owner string, openOrders []string, project pb.Project, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	orders, err := h.PostCancelAll(ctx, market, owner, openOrders, project)
	if err != nil {
		return nil, err
	}
	return h.SignAndSubmitBatch(ctx, orders.Transactions, opts)
}

// PostSettle returns a partially signed transaction for settling market funds. Typically, you want to use SubmitSettle instead of this.
func (h *HTTPClient) PostSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, project pb.Project) (*pb.PostSettleResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/settle", h.baseURL)
	request := &pb.PostSettleRequest{
		OwnerAddress:      owner,
		Market:            market,
		BaseTokenWallet:   baseTokenWallet,
		QuoteTokenWallet:  quoteTokenWallet,
		OpenOrdersAddress: openOrdersAccount,
		Project:           project,
	}

	var response pb.PostSettleResponse
	err := connections.HTTPPostWithClient[*pb.PostSettleResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitSettle builds a market SubmitSettle transaction, signs it, and submits to the network.
func (h *HTTPClient) SubmitSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, project pb.Project, skipPreflight bool) (string, error) {
	order, err := h.PostSettle(ctx, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount, project)
	if err != nil {
		return "", err
	}

	return h.SignAndSubmit(ctx, order.Transaction, skipPreflight)
}

func (h *HTTPClient) PostReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/replacebyclientid", h.baseURL)
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
	err := connections.HTTPPostWithClient[*pb.PostOrderResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (h *HTTPClient) SubmitReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error) {
	order, err := h.PostReplaceByClientOrderID(ctx, owner, payer, market, side, types, amount, price, project, opts)
	if err != nil {
		return "", err
	}

	return h.SignAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

func (h *HTTPClient) PostReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v1/trade/replace", h.baseURL)
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
	err := connections.HTTPPostWithClient[*pb.PostOrderResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (h *HTTPClient) SubmitReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error) {
	order, err := h.PostReplaceOrder(ctx, orderID, owner, payer, market, side, types, amount, price, project, opts)
	if err != nil {
		return "", err
	}

	return h.SignAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

// GetRecentBlockHash subscribes to a stream for getting recent block hash.
func (h *HTTPClient) GetRecentBlockHash(ctx context.Context) (*pb.GetRecentBlockHashResponse, error) {
	url := fmt.Sprintf("%s/api/v1/system/blockhash", h.baseURL)
	response := new(pb.GetRecentBlockHashResponse)
	if err := connections.HTTPGetWithClient[*pb.GetRecentBlockHashResponse](ctx, url, h.httpClient, response, h.authHeader); err != nil {
		return nil, err
	}

	return response, nil
}

// GetPerpOrderbook returns the current state of perpetual contract orderbook.
func (h *HTTPClient) GetPerpOrderbook(ctx context.Context, request *pb.GetPerpOrderbookRequest) (*pb.GetPerpOrderbookResponse, error) {
	url := fmt.Sprintf("%s/api/v1/market/perp/orderbook/%s?limit=%d&project=%v", h.baseURL, request.Contract, request.Limit, request.Project)
	orderbook := new(pb.GetPerpOrderbookResponse)
	if err := connections.HTTPGetWithClient[*pb.GetPerpOrderbookResponse](ctx, url, h.httpClient, orderbook, h.authHeader); err != nil {
		return nil, err
	}

	return orderbook, nil
}

// GetDriftMarginOrderbook returns the current state of margin contract orderbook.
func (h *HTTPClient) GetDriftMarginOrderbook(ctx context.Context, request *pb.GetDriftMarginOrderbookRequest) (*pb.GetDriftMarginOrderbookResponse, error) {
	url := fmt.Sprintf("%s/api/v2/drift/margin/orderbooks/%s?limit=%d&metadata=%v", h.baseURL, request.Market, request.Limit, request.Metadata)
	orderbook := new(pb.GetDriftMarginOrderbookResponse)
	if err := connections.HTTPGetWithClient[*pb.GetDriftMarginOrderbookResponse](ctx, url, h.httpClient, orderbook, h.authHeader); err != nil {
		return nil, err
	}
	return orderbook, nil
}

// GetDriftMarketDepth returns the current state of Drift market depth.
func (h *HTTPClient) GetDriftMarketDepth(ctx context.Context, request *pb.GetDriftMarketDepthRequest) (*pb.GetDriftMarketDepthResponse, error) {
	url := fmt.Sprintf("%s/api/v2/drift/perp/market-depth/%s?limit=%d", h.baseURL, request.Contract, request.Limit)
	maarketDepthData := new(pb.GetDriftMarketDepthResponse)
	if err := connections.HTTPGetWithClient[*pb.GetDriftMarketDepthResponse](ctx, url, h.httpClient, maarketDepthData, h.authHeader); err != nil {
		return nil, err
	}

	return maarketDepthData, nil
}

//V2 Openbook

// GetMarketsV2 returns the list of all available named markets
func (h *HTTPClient) GetMarketsV2(ctx context.Context) (*pb.GetMarketsResponseV2, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/markets", h.baseURL)
	markets := new(pb.GetMarketsResponseV2)
	if err := connections.HTTPGetWithClient[*pb.GetMarketsResponseV2](ctx, url, h.httpClient, markets, h.authHeader); err != nil {
		return nil, err
	}

	return markets, nil
}

// GetOrderbookV2 returns the requested market's orderbook (e.h. asks and bids). Set limit to 0 for all bids / asks.
func (h *HTTPClient) GetOrderbookV2(ctx context.Context, market string, limit uint32) (*pb.GetOrderbookResponseV2, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/orderbooks/%s?limit=%v", h.baseURL, market, limit)
	orderbook := new(pb.GetOrderbookResponseV2)
	if err := connections.HTTPGetWithClient[*pb.GetOrderbookResponseV2](ctx, url, h.httpClient, orderbook, h.authHeader); err != nil {
		return nil, err
	}

	return orderbook, nil
}

// GetMarketDepthV2 returns the requested market's coalesced price data (e.h. asks and bids). Set limit to 0 for all bids / asks.
func (h *HTTPClient) GetMarketDepthV2(ctx context.Context, market string, limit uint32) (*pb.GetMarketDepthResponseV2, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/depth/%s?limit=%v", h.baseURL, market, limit)
	mktDepth := new(pb.GetMarketDepthResponseV2)
	if err := connections.HTTPGetWithClient[*pb.GetMarketDepthResponseV2](ctx, url, h.httpClient, mktDepth, h.authHeader); err != nil {
		return nil, err
	}

	return mktDepth, nil
}

// GetTickersV2 returns the requested market tickets. Set market to "" for all markets.
func (h *HTTPClient) GetTickersV2(ctx context.Context, market string) (*pb.GetTickersResponseV2, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/tickers/%s", h.baseURL, market)
	tickers := new(pb.GetTickersResponseV2)
	if err := connections.HTTPGetWithClient[*pb.GetTickersResponseV2](ctx, url, h.httpClient, tickers, h.authHeader); err != nil {
		return nil, err
	}

	return tickers, nil
}

// GetOpenOrdersV2 returns all open orders by owner address and market
func (h *HTTPClient) GetOpenOrdersV2(ctx context.Context, market string, owner string, openOrdersAddress string, orderID string, clientOrderID uint64) (*pb.GetOpenOrdersResponse, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/open-orders/%s?address=%s&openOrdersAddress=%s&orderID=%s&clientOrderID=%v",
		h.baseURL, market, owner, openOrdersAddress, orderID, clientOrderID)

	orders := new(pb.GetOpenOrdersResponse)
	if err := connections.HTTPGetWithClient[*pb.GetOpenOrdersResponse](ctx, url, h.httpClient, orders, h.authHeader); err != nil {
		return nil, err
	}

	return orders, nil
}

// GetUnsettledV2 returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (h *HTTPClient) GetUnsettledV2(ctx context.Context, market string, owner string) (*pb.GetUnsettledResponse, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/unsettled/%s?ownerAddress=%s", h.baseURL, market, owner)
	result := new(pb.GetUnsettledResponse)
	if err := connections.HTTPGetWithClient[*pb.GetUnsettledResponse](ctx, url, h.httpClient, result, h.authHeader); err != nil {
		return nil, err
	}

	return result, nil
}

// PostOrderV2 returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (h *HTTPClient) PostOrderV2(ctx context.Context, owner, payer, market string, side pb.Side, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/place", h.baseURL)
	request := &pb.PostOrderRequestV2{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Amount:            amount,
		Price:             price,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
	}

	var response pb.PostOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostOrderResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitOrderV2 builds a Serum market order, signs it, and submits to the network.
func (h *HTTPClient) SubmitOrderV2(ctx context.Context, owner, payer, market string, side pb.Side, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := h.PostOrderV2(ctx, owner, payer, market, side, amount, price, opts)
	if err != nil {
		return "", err
	}

	sig, err := h.SignAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
	return sig, err
}

// PostCancelOrderV2 builds a Serum cancel order.
func (h *HTTPClient) PostCancelOrderV2(
	ctx context.Context,
	orderID string,
	clientOrderID uint64,
	side pb.Side,
	owner,
	market,
	openOrders string,
) (*pb.PostCancelOrderResponseV2, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/cancel", h.baseURL)
	request := &pb.PostCancelOrderRequestV2{
		OrderID:           orderID,
		ClientOrderID:     clientOrderID,
		Side:              side,
		OwnerAddress:      owner,
		MarketAddress:     market,
		OpenOrdersAddress: openOrders,
	}

	var response pb.PostCancelOrderResponseV2
	err := connections.HTTPPostWithClient[*pb.PostCancelOrderResponseV2](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// SubmitCancelOrderV2 builds a Serum cancel order, signs and submits it to the network.
func (h *HTTPClient) SubmitCancelOrderV2(
	ctx context.Context,
	orderID string,
	clientOrderID uint64,
	side pb.Side,
	owner,
	market,
	openOrders string,
	opts SubmitOpts,
) (*pb.PostSubmitBatchResponse, error) {
	order, err := h.PostCancelOrderV2(ctx, orderID, clientOrderID, side, owner, market, openOrders)
	if err != nil {
		return nil, err
	}

	return h.SignAndSubmitBatch(ctx, order.Transactions, opts)
}

// PostSettleV2 returns a partially signed transaction for settling market funds. Typically, you want to use SubmitSettle instead of this.
func (h *HTTPClient) PostSettleV2(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string) (*pb.PostSettleResponse, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/settle", h.baseURL)
	request := &pb.PostSettleRequestV2{
		OwnerAddress:      owner,
		Market:            market,
		BaseTokenWallet:   baseTokenWallet,
		QuoteTokenWallet:  quoteTokenWallet,
		OpenOrdersAddress: openOrdersAccount,
	}

	var response pb.PostSettleResponse
	err := connections.HTTPPostWithClient[*pb.PostSettleResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitSettleV2 builds a market SubmitSettle transaction, signs it, and submits to the network.
func (h *HTTPClient) SubmitSettleV2(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, skipPreflight bool) (string, error) {
	order, err := h.PostSettleV2(ctx, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount)
	if err != nil {
		return "", err
	}

	return h.SignAndSubmit(ctx, order.Transaction, skipPreflight)
}

func (h *HTTPClient) PostReplaceOrderV2(ctx context.Context, orderID, owner, payer, market string, side pb.Side, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	url := fmt.Sprintf("%s/api/v2/openbook/replace", h.baseURL)
	request := &pb.PostReplaceOrderRequestV2{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Amount:            amount,
		Price:             price,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
		OrderID:           orderID,
	}

	var response pb.PostOrderResponse
	err := connections.HTTPPostWithClient[*pb.PostOrderResponse](ctx, url, h.httpClient, request, &response, h.authHeader)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (h *HTTPClient) SubmitReplaceOrderV2(ctx context.Context, orderID, owner, payer, market string, side pb.Side, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := h.PostReplaceOrderV2(ctx, orderID, owner, payer, market, side, amount, price, opts)
	if err != nil {
		return "", err
	}

	return h.SignAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

type stringable interface {
	String() string
}

func convertSliceArgument[T stringable](argName string, isFirst bool, s []T) string {
	r := ""
	for i, v := range s {
		arg := fmt.Sprintf("%v=%v", argName, v.String())
		if i == 0 && isFirst {
			r += "?" + arg
		} else {
			r += "&" + arg
		}
	}
	return r
}

func convertStrSliceArgument(argName string, isFirst bool, s []string) string {
	r := ""
	for i, v := range s {
		arg := fmt.Sprintf("%v=%v", argName, v)
		if i == 0 && isFirst {
			r += "?" + arg
		} else {
			r += "&" + arg
		}
	}
	return r
}
