package provider

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/serum-api/bxserum/helpers"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/gorilla/websocket"
)

type WSClient struct {
	pb.UnsafeApiServer

	addr string
	conn *websocket.Conn
}

// Connects to Mainnet Serum API
func NewWSClient() (*WSClient, error) {
	return NewWSClientWithEndpoint("ws://174.129.154.164:1810/ws")
}

// TODO Connects to Testnet Serum API
func NewWSClientTestnet() (*WSClient, error) {
	panic("implement me")
}

// Connects to specified Serum API
func NewWSClientWithEndpoint(addr string) (*WSClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}

	return &WSClient{addr: addr, conn: conn}, nil
}

func (w *WSClient) GetOrderbook(market string) (*pb.GetOrderbookResponse, error) {
	request := fmt.Sprintf(`{"jsonrpc": "2.0", "id": 1, "method": "GetOrderbook", "params": {"market":"%s"}}`, market)
	return helpers.UnaryWSRequest[pb.GetOrderbookResponse](w.conn, request)
}

func (w *WSClient) GetOrderbookStream(ctx context.Context, market string, orderbookChan chan *pb.GetOrderbookStreamResponse) error {
	request := fmt.Sprintf(`{"jsonrpc": "2.0", "id": 1, "method": "GetOrderbookStream", "params": {"market":"%s"}}`, market)
	return helpers.UnaryWSStream[pb.GetOrderbookStreamResponse](ctx, w.conn, request, orderbookChan)
}

func (w *WSClient) Close() error {
	err := w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		// TODO close conn harshly?
		return fmt.Errorf("error writing close msg -  %v", err)
	}
	return nil
}
