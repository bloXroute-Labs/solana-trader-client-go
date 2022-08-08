package provider

import (
	"context"
	"errors"
	"github.com/bloXroute-Labs/serum-client-go/bxserum/connections"
	"github.com/bloXroute-Labs/serum-client-go/bxserum/transaction"
	pb "github.com/bloXroute-Labs/serum-client-go/proto"
	"github.com/gagliardetto/solana-go"
)

type WSClient struct {
	pb.UnimplementedApiServer

	addr       string
	conn       *connections.WS
	privateKey *solana.PrivateKey
}

// NewWSClient connects to Mainnet Serum API
func NewWSClient() (*WSClient, error) {
	opts := DefaultRPCOpts(MainnetSerumAPIWS)
	return NewWSClientWithOpts(opts)
}

// NewWSClientTestnet connects to Testnet Serum API
func NewWSClientTestnet() (*WSClient, error) {
	opts := DefaultRPCOpts(TestnetSerumAPIWS)
	return NewWSClientWithOpts(opts)
}

// NewWSClientLocal connects to Local Serum API
func NewWSClientLocal() (*WSClient, error) {
	opts := DefaultRPCOpts(LocalSerumAPIWES)
	return NewWSClientWithOpts(opts)
}

// NewWSClientWithOpts connects to custom Serum API
func NewWSClientWithOpts(opts RPCOpts) (*WSClient, error) {
	conn, err := connections.NewWS(opts.Endpoint, opts.AuthHeader)
	if err != nil {
		return nil, err
	}

	return &WSClient{
		addr:       opts.Endpoint,
		conn:       conn,
		privateKey: opts.PrivateKey,
	}, nil
}

// GetOrderbook returns the requested market's orderbook (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (w *WSClient) GetOrderbook(ctx context.Context, market string, limit uint32) (*pb.GetOrderbookResponse, error) {
	var response pb.GetOrderbookResponse
	err := w.conn.Request(ctx, "GetOrderbook", &pb.GetOrderbookRequest{Market: market, Limit: limit}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetTrades returns the requested market's currently executing trades. Set limit to 0 for all trades.
func (w *WSClient) GetTrades(ctx context.Context, market string, limit uint32) (*pb.GetTradesResponse, error) {
	var response pb.GetTradesResponse
	err := w.conn.Request(ctx, "GetTrades", &pb.GetTradesRequest{Market: market, Limit: limit}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetTickers returns the requested market tickets. Set market to "" for all markets.
func (w *WSClient) GetTickers(ctx context.Context, market string) (*pb.GetTickersResponse, error) {
	var response pb.GetTickersResponse
	err := w.conn.Request(ctx, "GetTickers", &pb.GetTickersRequest{Market: market}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetOpenOrders returns all opened orders by owner address and market
func (w *WSClient) GetOpenOrders(ctx context.Context, market string, owner string) (*pb.GetOpenOrdersResponse, error) {
	var response pb.GetOpenOrdersResponse
	err := w.conn.Request(ctx, "GetOpenOrders", &pb.GetOpenOrdersRequest{Market: market, Address: owner}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetUnsettled returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (w *WSClient) GetUnsettled(ctx context.Context, market string, owner string) (*pb.GetUnsettledResponse, error) {
	var response pb.GetUnsettledResponse
	err := w.conn.Request(ctx, "GetUnsettled", &pb.GetUnsettledRequest{Market: market, Owner: owner}, &response)
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

// PostOrder returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (w *WSClient) PostOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
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
	}
	var response pb.PostOrderResponse
	err := w.conn.Request(ctx, "PostOrder", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostSubmit posts the transaction string to the Solana network.
func (w *WSClient) PostSubmit(ctx context.Context, txBase64 string, skipPreFlight bool) (*pb.PostSubmitResponse, error) {
	request := &pb.PostSubmitRequest{
		Transaction:   txBase64,
		SkipPreFlight: skipPreFlight,
	}
	var response pb.PostSubmitResponse
	err := w.conn.Request(ctx, "PostSubmit", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// signAndSubmit signs the given transaction and submits it.
func (w *WSClient) signAndSubmit(ctx context.Context, tx string, skipPreFlight bool) (string, error) {
	if w.privateKey == nil {
		return "", ErrPrivateKeyNotFound
	}

	txBase64, err := transaction.SignTxWithPrivateKey(tx, *w.privateKey)
	if err != nil {
		return "", err
	}

	response, err := w.PostSubmit(ctx, txBase64, skipPreFlight)
	if err != nil {
		return "", err
	}

	return response.Signature, nil
}

// SubmitOrder builds a Serum market order, signs it, and submits to the network.
func (w *WSClient) SubmitOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := w.PostOrder(ctx, owner, payer, market, side, types, amount, price, opts)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

// PostCancelOrder builds a Serum cancel order.
func (w *WSClient) PostCancelOrder(
	ctx context.Context,
	orderID string,
	side pb.Side,
	owner,
	market,
	openOrders string,
) (*pb.PostCancelOrderResponse, error) {
	request := &pb.PostCancelOrderRequest{
		OrderID:           orderID,
		Side:              side,
		OwnerAddress:      owner,
		MarketAddress:     market,
		OpenOrdersAddress: openOrders,
	}

	var response pb.PostCancelOrderResponse
	err := w.conn.Request(ctx, "PostCancelOrder", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitCancelOrder builds a Serum cancel order, signs and submits it to the network.
func (w *WSClient) SubmitCancelOrder(
	ctx context.Context,
	orderID string,
	side pb.Side,
	owner,
	market,
	openOrders string,
	skipPreFlight bool,
) (string, error) {
	order, err := w.PostCancelOrder(ctx, orderID, side, owner, market, openOrders)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, order.Transaction, skipPreFlight)
}

// PostCancelByClientOrderID builds a Serum cancel order by client ID.
func (w *WSClient) PostCancelByClientOrderID(
	ctx context.Context,
	clientOrderID uint64,
	owner,
	market,
	openOrders string,
) (*pb.PostCancelOrderResponse, error) {
	request := &pb.PostCancelByClientOrderIDRequest{
		ClientOrderID:     clientOrderID,
		OwnerAddress:      owner,
		MarketAddress:     market,
		OpenOrdersAddress: openOrders,
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
	skipPreFlight bool,
) (string, error) {
	order, err := w.PostCancelByClientOrderID(ctx, clientOrderID, owner, market, openOrders)
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
) (*pb.PostCancelAllResponse, error) {
	request := &pb.PostCancelAllRequest{
		Market:              market,
		OwnerAddress:        owner,
		OpenOrdersAddresses: openOrdersAddresses,
	}
	var response pb.PostCancelAllResponse
	err := w.conn.Request(ctx, "PostCancelAll", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (w *WSClient) SubmitCancelAll(
	ctx context.Context,
	market,
	owner string,
	openOrdersAddresses []string,
	skipPreFlight bool,
) ([]string, error) {
	orders, err := w.PostCancelAll(ctx, market, owner, openOrdersAddresses)
	if err != nil {
		return nil, err
	}

	var signatures []string
	for _, tx := range orders.Transactions {
		signature, err := w.signAndSubmit(ctx, tx, skipPreFlight)
		if err != nil {
			return signatures, err
		}

		signatures = append(signatures, signature)
	}

	return signatures, nil
}

// PostSettle returns a partially signed transaction for settling market funds. Typically, you want to use SubmitSettle instead of this.
func (w *WSClient) PostSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string) (*pb.PostSettleResponse, error) {
	request := &pb.PostSettleRequest{
		OwnerAddress:      owner,
		Market:            market,
		BaseTokenWallet:   baseTokenWallet,
		QuoteTokenWallet:  quoteTokenWallet,
		OpenOrdersAddress: openOrdersAccount,
	}
	var response pb.PostSettleResponse
	err := w.conn.Request(ctx, "PostSettle", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitSettle builds a market SubmitSettle transaction, signs it, and submits to the network.
func (w *WSClient) SubmitSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, skipPreflight bool) (string, error) {
	order, err := w.PostSettle(ctx, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount)
	if err != nil {
		return "", err
	}
	return w.signAndSubmit(ctx, order.Transaction, skipPreflight)
}

func (w *WSClient) PostReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
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
	}
	var response pb.PostOrderResponse
	err := w.conn.Request(ctx, "PostReplaceByClientOrderID", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (w *WSClient) SubmitReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := w.PostReplaceByClientOrderID(ctx, owner, payer, market, side, types, amount, price, opts)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

func (w *WSClient) PostReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	request := &pb.PostReplaceOrderRequest{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              types,
		Amount:            amount,
		Price:             price,
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

func (w *WSClient) SubmitReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := w.PostReplaceOrder(ctx, orderID, owner, payer, market, side, types, amount, price, opts)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

func (w *WSClient) Close() error {
	return w.conn.Close(errors.New("shutdown requested"))
}

// GetOrderbooksStream subscribes to a stream for changes to the requested market updates (e.g. asks and bids. Set limit to 0 for all bids/ asks).
func (w *WSClient) GetOrderbooksStream(ctx context.Context, markets []string, limit uint32) (connections.Streamer[*pb.GetOrderbooksStreamResponse], error) {
	return connections.WSStream(w.conn, ctx, "GetOrderbooksStream", &pb.GetOrderbooksRequest{
		Markets: markets,
		Limit:   limit,
	}, func() *pb.GetOrderbooksStreamResponse {
		var v pb.GetOrderbooksStreamResponse
		return &v
	})
}

// GetTradesStream subscribes to a stream for trades as they execute. Set limit to 0 for all trades.
func (w *WSClient) GetTradesStream(ctx context.Context, market string, limit uint32) (connections.Streamer[*pb.GetTradesStreamResponse], error) {
	return connections.WSStream(w.conn, ctx, "GetTradesStream", &pb.GetTradesRequest{
		Market: market,
		Limit:  limit,
	}, func() *pb.GetTradesStreamResponse {
		var v pb.GetTradesStreamResponse
		return &v
	})
}

// GetOrderStatusStream subscribes to a stream that shows updates to the owner's orders
func (w *WSClient) GetOrderStatusStream(ctx context.Context, market, ownerAddress string) (connections.Streamer[*pb.GetOrderStatusStreamResponse], error) {
	return connections.WSStream(w.conn, ctx, "GetOrderStatusStream", &pb.GetOrderStatusStreamRequest{
		Market:       market,
		OwnerAddress: ownerAddress,
	}, func() *pb.GetOrderStatusStreamResponse {
		var v pb.GetOrderStatusStreamResponse
		return &v
	})
}
