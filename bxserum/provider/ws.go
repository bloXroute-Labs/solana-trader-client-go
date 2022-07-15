package provider

import (
	"context"
	"fmt"

	"github.com/bloXroute-Labs/serum-client-go/bxserum/connections"
	"github.com/bloXroute-Labs/serum-client-go/bxserum/transaction"
	pb "github.com/bloXroute-Labs/serum-client-go/proto"
	"github.com/bloXroute-Labs/serum-client-go/utils"
	"github.com/gagliardetto/solana-go"
	"github.com/gorilla/websocket"
	"github.com/sourcegraph/jsonrpc2"
)

type WSClient struct {
	pb.UnimplementedApiServer

	addr       string
	conn       *websocket.Conn
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
	conn, _, err := websocket.DefaultDialer.Dial(opts.Endpoint, nil)
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
	request, err := w.jsonRPCRequest("GetOrderbook", map[string]interface{}{"market": market, "limit": limit})
	if err != nil {
		return nil, err
	}

	var response pb.GetOrderbookResponse
	err = connections.WSRequest(w.conn, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetOrderbookStream subscribes to a stream for changes to the requested market updates (e.g. asks and bids. Set limit to 0 for all bids/ asks).
func (w *WSClient) GetOrderbookStream(ctx context.Context, market string, limit uint32, orderbookChan chan *pb.GetOrderbooksStreamResponse) error {
	request, err := w.jsonRPCRequest("GetOrderbooksStream", map[string]interface{}{"market": market, "limit": limit})
	if err != nil {
		return err
	}

	var response pb.GetOrderbooksStreamResponse
	return connections.WSStream(ctx, w.conn, request, orderbookChan, &response)
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
	request, err := w.jsonRPCRequest("GetTrades", map[string]interface{}{"market": market, "limit": limit})
	if err != nil {
		return nil, err
	}

	var response pb.GetTradesResponse
	err = connections.WSRequest(w.conn, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetTradesStream subscribes to a stream for trades as they execute. Set limit to 0 for all trades.
func (w *WSClient) GetTradesStream(ctx context.Context, market string, limit uint32, tradesChan chan *pb.GetTradesStreamResponse) error {
	request, err := w.jsonRPCRequest("GetTradesStream", map[string]interface{}{"market": market, "limit": limit})
	if err != nil {
		return err
	}

	var response pb.GetTradesStreamResponse
	return connections.WSStream(ctx, w.conn, request, tradesChan, &response)
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
	request, err := w.jsonRPCRequest("GetTickers", map[string]interface{}{"market": market})
	if err != nil {
		return nil, err
	}

	var response pb.GetTickersResponse
	err = connections.WSRequest(w.conn, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetOpenOrders returns all opened orders by owner address and market
func (w *WSClient) GetOpenOrders(market string, owner string) (*pb.GetOpenOrdersResponse, error) {
	request, err := w.jsonRPCRequest("GetOpenOrders", map[string]interface{}{"market": market, "address": owner})
	if err != nil {
		return nil, err
	}

	var response pb.GetOpenOrdersResponse
	err = connections.WSRequest(w.conn, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetUnsettled returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (w *WSClient) GetUnsettled(market string, owner string) (*pb.GetUnsettledResponse, error) {
	request, err := w.jsonRPCRequest("GetUnsettled", map[string]interface{}{"market": market, "owner": owner})
	if err != nil {
		return nil, err
	}

	var response pb.GetUnsettledResponse
	err = connections.WSRequest(w.conn, request, &response)
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
	var params struct{}
	request, err := w.jsonRPCRequest("GetMarkets", params)
	if err != nil {
		return nil, err
	}

	var response pb.GetMarketsResponse
	err = connections.WSRequest(w.conn, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostOrder returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (w *WSClient) PostOrder(owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	request, err := w.jsonRPCRequest("PostOrder", &pb.PostOrderRequest{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              types,
		Amount:            amount,
		Price:             price,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
	})
	if err != nil {
		return nil, err
	}

	var response pb.PostOrderResponse
	err = connections.WSRequest(w.conn, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// PostSubmit posts the transaction string to the Solana network.
func (w *WSClient) PostSubmit(txBase64 string, skipPreFlight bool) (*pb.PostSubmitResponse, error) {
	request, err := w.jsonRPCRequest("PostSubmit", &pb.PostSubmitRequest{
		Transaction:   txBase64,
		SkipPreFlight: skipPreFlight,
	})
	if err != nil {
		return nil, err
	}

	var response pb.PostSubmitResponse
	err = connections.WSRequest(w.conn, request, &response)
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
	request, err := w.jsonRPCRequest("PostCancelOrder", &pb.PostCancelOrderRequest{
		OrderID:           orderID,
		Side:              side,
		OwnerAddress:      owner,
		MarketAddress:     market,
		OpenOrdersAddress: openOrders,
	})
	if err != nil {
		return nil, err
	}

	var response pb.PostCancelOrderResponse
	err = connections.WSRequest(w.conn, request, &response)
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
	request, err := w.jsonRPCRequest("PostCancelByClientOrderID", &pb.PostCancelByClientOrderIDRequest{
		ClientOrderID:     clientOrderID,
		OwnerAddress:      owner,
		MarketAddress:     market,
		OpenOrdersAddress: openOrders,
	})
	if err != nil {
		return nil, err
	}

	var response pb.PostCancelOrderResponse
	err = connections.WSRequest(w.conn, request, &response)
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

func (w *WSClient) PostCancelAll(market, owner, openOrders string) (*pb.PostCancelAllResponse, error) {
	request, err := w.jsonRPCRequest("PostCancelAll", &pb.PostCancelAllRequest{
		Market:           market,
		OwnerAddress:     owner,
		OpenOrderAddress: openOrders,
	})
	if err != nil {
		return nil, err
	}

	var response pb.PostCancelAllResponse
	err = connections.WSRequest(w.conn, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (w *WSClient) SubmitCancelAll(market, owner, openOrders string, skipPreFlight bool) ([]string, error) {
	orders, err := w.PostCancelAll(market, owner, openOrders)
	if err != nil {
		return nil, err
	}

	var signatures []string
	for _, tx := range orders.Transactions {
		signature, err := w.signAndSubmit(tx, skipPreFlight)
		if err != nil {
			return signatures, err
		}

		signatures = append(signatures, signature)
	}

	return signatures, nil
}

// PostSettle returns a partially signed transaction for settling market funds. Typically, you want to use SubmitSettle instead of this.
func (w *WSClient) PostSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string) (*pb.PostSettleResponse, error) {
	request, err := w.jsonRPCRequest("PostSettle", &pb.PostSettleRequest{
		OwnerAddress:      owner,
		Market:            market,
		BaseTokenWallet:   baseTokenWallet,
		QuoteTokenWallet:  quoteTokenWallet,
		OpenOrdersAddress: openOrdersAccount,
	})
	if err != nil {
		return nil, err
	}

	var response pb.PostSettleResponse
	err = connections.WSRequest(w.conn, request, &response)
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
	return w.signAndSubmit(order.Transaction, skipPreflight)
}

func (w *WSClient) Close() error {
	err := w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return fmt.Errorf("error writing close msg -  %v", err)
	}
	return nil
}

func (w *WSClient) jsonRPCRequest(method string, params interface{}) ([]byte, error) {
	id := w.requestID.Next()
	req := jsonrpc2.Request{
		Method: method,
		ID:     jsonrpc2.ID{Num: id},
	}
	if err := req.SetParams(params); err != nil {
		return nil, err
	}

	return req.MarshalJSON()
}
