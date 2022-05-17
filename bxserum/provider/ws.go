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
	return NewWSClientWithEndpoint("ws://174.129.154.164:1810/ws")
}

// NewWSClientTestnet connects to Testnet Serum API
func NewWSClientTestnet() (*WSClient, error) {
	panic("implement me")
}

// NewWSClientWithEndpoint onnects to custom Serum API
func NewWSClientWithEndpoint(addr string) (*WSClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}

	privateKey, err := transaction.LoadPrivateKeyFromEnv()
	if err != nil {
		return nil, err
	}

	return &WSClient{
		addr:       addr,
		conn:       conn,
		requestID:  utils.NewRequestID(),
		privateKey: privateKey,
	}, nil
}

// GetOrderbook returns the requested market's orderbook (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (w *WSClient) GetOrderbook(market string, limit uint32) (*pb.GetOrderbookResponse, error) {
	request, err := w.jsonRPCRequest("GetOrderbook", map[string]interface{}{"market": market, "limit": limit})
	if err != nil {
		return nil, err
	}
	return connections.WSRequest[pb.GetOrderbookResponse](w.conn, request)
}

func (w *WSClient) GetOrderbookStream(ctx context.Context, market string, limit uint32, orderbookChan chan *pb.GetOrderbookStreamResponse) error {
	request, err := w.jsonRPCRequest("GetOrderbookStream", map[string]interface{}{"market": market, "limit": limit})
	if err != nil {
		return err
	}
	return connections.WSStream[pb.GetOrderbookStreamResponse](ctx, w.conn, request, orderbookChan)
}

// GetTrades returns the requested market's currently executing trades. Set limit to 0 for all trades.
func (w *WSClient) GetTrades(market string, limit uint32) (*pb.GetTradesResponse, error) {
	request, err := w.jsonRPCRequest("GetTrades", map[string]interface{}{"market": market, "limit": limit})
	if err != nil {
		return nil, err
	}
	return connections.WSRequest[pb.GetTradesResponse](w.conn, request)
}

func (w *WSClient) GetTradesStream(ctx context.Context, market string, limit uint32, tradesChan chan *pb.GetTradesStreamResponse) error {
	request, err := w.jsonRPCRequest("GetTradeStream", map[string]interface{}{"market": market, "limit": limit})
	if err != nil {
		return err
	}
	return connections.WSStream[pb.GetTradesStreamResponse](ctx, w.conn, request, tradesChan)
}

// GetTickers returns the requested market tickets. Set market to "" for all markets.
func (w *WSClient) GetTickers(market string) (*pb.GetTickersResponse, error) {
	request, err := w.jsonRPCRequest("GetTickers", map[string]interface{}{"market": market})
	if err != nil {
		return nil, err
	}
	return connections.WSRequest[pb.GetTickersResponse](w.conn, request)
}

// GetOrders returns all opened orders by owner address and market
func (w *WSClient) GetOrders(market string, owner string) (*pb.GetOrdersResponse, error) {
	request, err := w.jsonRPCRequest("GetOrders", map[string]interface{}{"market": market, "address": owner})
	if err != nil {
		return nil, err
	}
	return connections.WSRequest[pb.GetOrdersResponse](w.conn, request)
}

// GetMarkets returns the list of all available named markets
func (w *WSClient) GetMarkets() (*pb.GetMarketsResponse, error) {
	request, err := w.jsonRPCRequest("GetMarkets", nil)
	if err != nil {
		return nil, err
	}
	return connections.WSRequest[pb.GetMarketsResponse](w.conn, request)
}

// PostOrder returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (w *WSClient) PostOrder(owner, payer, market string, side pb.Side, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	request, err := w.jsonRPCRequest("PostOrder", &pb.PostOrderRequest{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              []pb.OrderType{pb.OrderType_OT_LIMIT},
		Amount:            amount,
		Price:             price,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
	})
	if err != nil {
		return nil, err
	}

	return connections.WSRequest[pb.PostOrderResponse](w.conn, request)
}

// PostSubmit posts the transaction string to the Solana network.
func (w *WSClient) PostSubmit(txBase64 string) (*pb.PostSubmitResponse, error) {
	request, err := w.jsonRPCRequest("PostSubmit", &pb.PostSubmitRequest{Transaction: txBase64})
	if err != nil {
		return nil, err
	}

	return connections.WSRequest[pb.PostSubmitResponse](w.conn, request)
}

// SubmitOrder builds a Serum market order, signs it, and submits to the network.
func (w *WSClient) SubmitOrder(owner, payer, market string, side pb.Side, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := w.PostOrder(owner, payer, market, side, amount, price, opts)
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
