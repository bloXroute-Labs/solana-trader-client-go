package provider

import (
	"context"
	"fmt"

	"github.com/bloXroute-Labs/serum-api/bxserum/connections"
	"github.com/bloXroute-Labs/serum-api/bxserum/transaction"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/bloXroute-Labs/serum-api/utils"
	"github.com/gagliardetto/solana-go"
	"github.com/gorilla/websocket"
	"github.com/sourcegraph/jsonrpc2"
)

type WSClient struct {
	pb.UnimplementedApiServer

	addr       string
	conn       *websocket.Conn
	requestID  utils.RequestID
	privateKey solana.PrivateKey
}

// NewWSClient connects to Mainnet Serum API
func NewWSClient() (*WSClient, error) {
	opts, err := DefaultRPCOpts(MainnetSerumAPIWS)
	if err != nil {
		return nil, err
	}
	return NewWSClientWithOpts(opts)
}

// NewWSClientTestnet connects to Testnet Serum API
func NewWSClientTestnet() (*WSClient, error) {
	opts, err := DefaultRPCOpts(TestnetSerumAPIWS)
	if err != nil {
		return nil, err
	}
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
	request, err := w.jsonRPCRequest("GetOrderbookStream", map[string]interface{}{"market": market, "limit": limit})
	if err != nil {
		return err
	}
	return connections.WSStream(ctx, w.conn, request, orderbookChan)
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
	request, err := w.jsonRPCRequest("GetTradeStream", map[string]interface{}{"market": market, "limit": limit})
	if err != nil {
		return err
	}
	return connections.WSStream(ctx, w.conn, request, tradesChan)
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
func (w *WSClient) PostSubmit(txBase64 string) (*pb.PostSubmitResponse, error) {
	request, err := w.jsonRPCRequest("PostSubmit", &pb.PostSubmitRequest{Transaction: txBase64})
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
func (w *WSClient) signAndSubmit(tx string) (string, error) {
	txBase64, err := transaction.SignTxWithPrivateKey(tx, w.privateKey)
	if err != nil {
		return "", err
	}

	response, err := w.PostSubmit(txBase64)
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

	return w.signAndSubmit(order.Transaction)
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
) (string, error) {
	order, err := w.PostCancelOrder(orderID, side, owner, market, openOrders)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(order.Transaction)
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
) (string, error) {
	order, err := w.PostCancelByClientOrderID(clientOrderID, owner, market, openOrders)
	if err != nil {
		return "", err
	}

	return w.signAndSubmit(order.Transaction)
}

// PostSettle returns a partially signed transaction for settling market funds. Typically, you want to use SettleFunds instead of this.
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

// SettleFunds builds a market SettleFunds transaction, signs it, and submits to the network.
func (w *WSClient) SettleFunds(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string) (string, error) {
	order, err := w.PostSettle(ctx, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount)
	if err != nil {
		return "", err
	}

	txBase64, err := transaction.SignTxWithPrivateKey(order.Transaction, w.privateKey)
	if err != nil {
		return "", err
	}

	response, err := w.PostSubmit(txBase64)
	if err != nil {
		return "", err
	}
	return response.Signature, nil
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
