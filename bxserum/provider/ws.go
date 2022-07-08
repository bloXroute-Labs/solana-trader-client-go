package provider

import (
	"context"
	"github.com/bloXroute-Labs/serum-client-go/bxserum/connections"
	"github.com/bloXroute-Labs/serum-client-go/bxserum/transaction"
	pb "github.com/bloXroute-Labs/serum-client-go/proto"
	"github.com/bloXroute-Labs/serum-client-go/utils"
	"github.com/gagliardetto/solana-go"
)

type WSClient struct {
	pb.UnimplementedApiServer

	addr       string
	conn       *connections.WS
	requestID  utils.RequestID
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

// NewWSClientWithOpts connects to custom Serum API
func NewWSClientWithOpts(opts RPCOpts) (*WSClient, error) {
	conn, err := connections.NewWS(opts.Endpoint)
	if err != nil {
		return nil, err
	}

	return &WSClient{
		addr:       opts.Endpoint,
		conn:       conn,
		requestID:  utils.NewRequestID(),
		privateKey: opts.PrivateKey,
	}, nil
}

// GetOrderbook returns the requested market's orderbook (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (w *WSClient) GetOrderbook(market string, limit uint32) (*pb.GetOrderbookResponse, error) {
	var response pb.GetOrderbookResponse
	err := w.conn.Request("GetOrderbook", &pb.GetOrderBookRequest{Market: market, Limit: limit}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetOrderbooksStream subscribes to a stream for changes to the requested market updates (e.g. asks and bids. Set limit to 0 for all bids/ asks).
func (w *WSClient) GetOrderbooksStream(ctx context.Context, market string, limit uint32, orderbookChan chan *pb.GetOrderbooksStreamResponse) error {
	generator, err := connections.WSStream(w.conn, ctx, "GetOrderbooksStream", &pb.GetOrderBookRequest{
		Market: market,
		Limit:  limit,
	}, func() *pb.GetOrderbooksStreamResponse {
		var v pb.GetOrderbooksStreamResponse
		return &v
	})
	if err != nil {
		return err
	}

	go func() {
		result, err := generator()
		if err != nil {
			close(orderbookChan)
			return
		}
		orderbookChan <- result
	}()

	return nil
}

// GetFilteredOrderbooksStream subscribes to a stream for changes to the requested market updates (e.g. asks and bids. Set limit to 0 for all bids/ asks).
func (w *WSClient) GetFilteredOrderbooksStream(ctx context.Context, markets []string, limit uint32, orderbookChan chan *pb.GetOrderbooksStreamResponse) error {
	request, err := w.jsonRPCRequest("GetFilteredOrderbooksStream", map[string]interface{}{"markets": markets, "limit": limit})
	if err != nil {
		return err
	}

	var response pb.GetOrderbooksStreamResponse
	return connections.WSStream(ctx, w.conn, request, orderbookChan, &response)
}

// GetTrades returns the requested market's currently executing trades. Set limit to 0 for all trades.
func (w *WSClient) GetTrades(market string, limit uint32) (*pb.GetTradesResponse, error) {
	var response pb.GetTradesResponse
	err := w.conn.Request("GetTrades", &pb.GetTradesRequest{Market: market, Limit: limit}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetTradesStream subscribes to a stream for trades as they execute. Set limit to 0 for all trades.
func (w *WSClient) GetTradesStream(ctx context.Context, market string, limit uint32, tradesChan chan *pb.GetTradesStreamResponse) error {
	generator, err := connections.WSStream(w.conn, ctx, "GetTradesStream", &pb.GetTradesRequest{
		Market: market,
		Limit:  limit,
	}, func() *pb.GetTradesStreamResponse {
		var v pb.GetTradesStreamResponse
		return &v
	})
	if err != nil {
		return err
	}

	go func() {
		result, err := generator()
		if err != nil {
			close(tradesChan)
			return
		}
		tradesChan <- result
	}()

	return nil
}

// GetOrderStatusStream subscribes to a stream that shows updates to the owner's orders
func (w *WSClient) GetOrderStatusStream(ctx context.Context, market, ownerAddress string, statusUpdateChan chan *pb.GetOrderStatusStreamResponse) error {
	request, err := w.jsonRPCRequest("GetOrderStatusStream", map[string]interface{}{"market": market, "ownerAddress": ownerAddress})
	if err != nil {
		return err
	}

	var response pb.GetOrderStatusStreamResponse
	return connections.WSStream(ctx, w.conn, request, statusUpdateChan, &response)
}

// GetTickers returns the requested market tickets. Set market to "" for all markets.
func (w *WSClient) GetTickers(market string) (*pb.GetTickersResponse, error) {
	var response pb.GetTickersResponse
	err := w.conn.Request("GetTickers", &pb.GetTickersRequest{Market: market}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetOpenOrders returns all opened orders by owner address and market
func (w *WSClient) GetOpenOrders(market string, owner string) (*pb.GetOpenOrdersResponse, error) {
	var response pb.GetOpenOrdersResponse
	err := w.conn.Request("GetOpenOrders", &pb.GetOpenOrdersRequest{Market: market, Address: owner}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetUnsettled returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (w *WSClient) GetUnsettled(market string, owner string) (*pb.GetUnsettledResponse, error) {
	var response pb.GetUnsettledResponse
	err := w.conn.Request("GetUnsettled", &pb.GetUnsettledRequest{Market: market, Owner: owner}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetAccountBalance returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (w *WSClient) GetAccountBalance(owner string) (*pb.GetAccountBalanceResponse, error) {
	request, err := w.jsonRPCRequest("GetAccountBalance", map[string]interface{}{"ownerAddress": owner})
	if err != nil {
		return nil, err
	}

	var response pb.GetAccountBalanceResponse
	err = connections.WSRequest(w.conn, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetMarkets returns the list of all available named markets
func (w *WSClient) GetMarkets() (*pb.GetMarketsResponse, error) {
	var response pb.GetMarketsResponse
	err := w.conn.Request("GetMarkets", &pb.GetMarketsRequest{}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostOrder returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (w *WSClient) PostOrder(owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
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
	err := w.conn.Request("PostOrder", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostSubmit posts the transaction string to the Solana network.
func (w *WSClient) PostSubmit(txBase64 string, skipPreFlight bool) (*pb.PostSubmitResponse, error) {
	request := &pb.PostSubmitRequest{
		Transaction:   txBase64,
		SkipPreFlight: skipPreFlight,
	}
	var response pb.PostSubmitResponse
	err := w.conn.Request("PostSubmit", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// signAndSubmit signs the given transaction and submits it.
func (w *WSClient) signAndSubmit(tx string, skipPreFlight bool) (string, error) {
	if w.privateKey == nil {
		return "", ErrPrivateKeyNotFound
	}

	txBase64, err := transaction.SignTxWithPrivateKey(tx, *w.privateKey)
	if err != nil {
		return "", err
	}

	response, err := w.PostSubmit(txBase64, skipPreFlight)
	if err != nil {
		return "", err
	}

	return response.Signature, nil
}

// SubmitOrder builds a Serum market order, signs it, and submits to the network.
func (w *WSClient) SubmitOrder(owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := w.PostOrder(owner, payer, market, side, types, amount, price, opts)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(order.Transaction, opts.SkipPreFlight)
}

// PostCancelOrder builds a Serum cancel order.
func (w *WSClient) PostCancelOrder(
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
	err := w.conn.Request("PostCancelOrder", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitCancelOrder builds a Serum cancel order, signs and submits it to the network.
func (w *WSClient) SubmitCancelOrder(
	orderID string,
	side pb.Side,
	owner,
	market,
	openOrders string,
	skipPreFlight bool,
) (string, error) {
	order, err := w.PostCancelOrder(orderID, side, owner, market, openOrders)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(order.Transaction, skipPreFlight)
}

// PostCancelByClientOrderID builds a Serum cancel order by client ID.
func (w *WSClient) PostCancelByClientOrderID(
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
	err := w.conn.Request("PostCancelByClientOrderID", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitCancelByClientOrderID builds a Serum cancel order by client ID, signs and submits it to the network.
func (w *WSClient) SubmitCancelByClientOrderID(
	clientOrderID uint64,
	owner,
	market,
	openOrders string,
	skipPreFlight bool,
) (string, error) {
	order, err := w.PostCancelByClientOrderID(clientOrderID, owner, market, openOrders)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(order.Transaction, skipPreFlight)
}

// PostSettle returns a partially signed transaction for settling market funds. Typically, you want to use SubmitSettle instead of this.
func (w *WSClient) PostSettle(owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string) (*pb.PostSettleResponse, error) {
	request := &pb.PostSettleRequest{
		OwnerAddress:      owner,
		Market:            market,
		BaseTokenWallet:   baseTokenWallet,
		QuoteTokenWallet:  quoteTokenWallet,
		OpenOrdersAddress: openOrdersAccount,
	}
	var response pb.PostSettleResponse
	err := w.conn.Request("PostSettle", request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// SubmitSettle builds a market SubmitSettle transaction, signs it, and submits to the network.
func (w *WSClient) SubmitSettle(owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, skipPreflight bool) (string, error) {
	order, err := w.PostSettle(owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount)
	if err != nil {
		return "", err
	}
	return w.signAndSubmit(order.Transaction, skipPreflight)
}
