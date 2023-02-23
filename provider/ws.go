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
	opts := DefaultRPCOpts(MainnetWS)
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

// GetOpenOrders returns all opened orders by owner address and market
func (w *WSClient) GetOpenOrders(ctx context.Context, market string, owner string, openOrdersAddress string, project pb.Project) (*pb.GetOpenOrdersResponse, error) {
	var response pb.GetOpenOrdersResponse
	err := w.conn.Request(ctx, "GetOpenOrders", &pb.GetOpenOrdersRequest{Market: market, Address: owner, OpenOrdersAddress: openOrdersAddress, Project: project}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetOpenPerpOrders returns all opened perp orders
func (w *WSClient) GetOpenPerpOrders(ctx context.Context, request *pb.GetOpenPerpOrdersRequest) (*pb.GetOpenPerpOrdersResponse, error) {
	var response pb.GetOpenPerpOrdersResponse
	err := w.conn.Request(ctx, "GetOpenPerpOrders", request, &response)
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

// SubmitPerpOrder builds a perp order, signs it, and submits to the network.
func (w *WSClient) SubmitPerpOrder(ctx context.Context, request *pb.PostPerpOrderRequest, opts PostOrderOpts) (string, error) {

	order, err := w.PostPerpOrder(ctx, request)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, &pb.TransactionMessage{
		Content: order.Transaction,
	}, opts.SkipPreFlight)
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
		msgs = append(msgs, &pb.TransactionMessage{Content: txn})
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

	return w.signAndSubmit(ctx, &pb.TransactionMessage{Content: resp.Transaction}, skipPreFlight)
}

// SubmitCancelPerpOrders builds a cancel perp orders txn, signs and submits it to the network.
func (w *WSClient) SubmitCancelPerpOrders(ctx context.Context, request *pb.PostCancelPerpOrdersRequest, skipPreFlight bool) (string, error) {
	resp, err := w.PostCancelPerpOrders(ctx, request)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, &pb.TransactionMessage{Content: resp.Transaction}, skipPreFlight)
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

func (w *WSClient) SubmitManageCollateral(ctx context.Context, req *pb.PostManageCollateralRequest, skipPreFlight bool) (string, error) {
	resp, err := w.PostManageCollateral(ctx, req)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, &pb.TransactionMessage{Content: resp.Transaction}, skipPreFlight)
}

func (w *WSClient) SubmitCreateUser(ctx context.Context, req *pb.PostCreateUserRequest, skipPreFlight bool) (interface{}, interface{}) {
	resp, err := w.PostCreateUser(ctx, req)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, &pb.TransactionMessage{Content: resp.Transaction}, skipPreFlight)
}

func (w *WSClient) SubmitPostPerpOrder(ctx context.Context, req *pb.PostPerpOrderRequest, skipPreFlight bool) (interface{}, interface{}) {
	resp, err := w.PostPerpOrder(ctx, req)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, &pb.TransactionMessage{Content: resp.Transaction}, skipPreFlight)
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

// GetPerpOrderbooksStream subscribes to a stream for perpetual orderbook updates.
func (w *WSClient) GetPerpOrderbooksStream(ctx context.Context, request *pb.GetPerpOrderbooksRequest) (connections.Streamer[*pb.GetPerpOrderbooksStreamResponse], error) {
	newResponse := func() *pb.GetPerpOrderbooksStreamResponse {
		return &pb.GetPerpOrderbooksStreamResponse{}
	}
	return connections.WSStreamProto(w.conn, ctx, "GetPerpOrderbooksStream", request, newResponse)
}
