package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bloXroute-Labs/serum-api/logger"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/gorilla/websocket"
	"github.com/sourcegraph/jsonrpc2"
)

type WSClient struct {
	pb.UnsafeApiServer

	addr string
	conn *websocket.Conn
}

func NewWSClient(addr string) (*WSClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}

	return &WSClient{addr: addr, conn: conn}, nil
}

func (w *WSClient) GetOrderbook(market string) (*pb.GetOrderbookResponse, error) {
	request := fmt.Sprintf(`{"jsonrpc": "2.0", "id": 1, "method": "GetOrderbook", "params": {"market":"%s"}}`, market)
	return unaryWSRequest[pb.GetOrderbookResponse](w.conn, request)
}

func (w *WSClient) GetOrderbookStream(ctx context.Context, market string, orderbookChan chan *pb.GetOrderbookStreamResponse) error {
	request := fmt.Sprintf(`{"jsonrpc": "2.0", "id": 1, "method": "GetOrderbookStream", "params": {"market":"%s"}}`, market)
	return unaryWSStream[pb.GetOrderbookStreamResponse](ctx, w.conn, request, orderbookChan)
}

func unaryWSRequest[T any](conn *websocket.Conn, request string) (*T, error) {
	err := sendWSRequest(conn, request)
	if err != nil {
		return nil, err
	}

	return recvWSResponse[T](conn)
}

func unaryWSStream[T any](ctx context.Context, conn *websocket.Conn, request string, responseChan chan *T) error {
	err := sendWSRequest(conn, request)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				response, err := recvWSResponse[T](conn)
				if err != nil {
					logger.Log().Error(err)
					continue
				}

				responseChan <- response
			}
		}
	}()

	return nil
}

func sendWSRequest(conn *websocket.Conn, request string) error {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(request)); err != nil {
		return fmt.Errorf("error with sending message - %v", err)
	}
	return nil
}

func recvWSResponse[T any](conn *websocket.Conn) (*T, error) {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("error reading WS response - %v", err)
	}

	// extract the result
	var r jsonrpc2.Response
	err = r.UnmarshalJSON(msg)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON response - %v", err)
	}
	bytes, err := r.Result.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON data - %v", err)
	}

	var response T
	if err = json.Unmarshal(bytes, &response); err != nil {
		return nil, fmt.Errorf("error with unmarshalling message of type %T - %v", response, err)
	}
	return &response, nil
}

func (w *WSClient) CloseConn() error {
	err := w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		// TODO close conn harshly?
		return fmt.Errorf("error writing close msg -  %v", err)
	}
	return nil
}
