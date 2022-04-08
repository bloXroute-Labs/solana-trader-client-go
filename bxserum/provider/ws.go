package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bloXroute-Labs/serum-api/borsh/serumborsh"
	"github.com/bloXroute-Labs/serum-api/logger"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/gorilla/websocket"
	"reflect"
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

func (w *WSClient) GetOrderbook(market string) (*serumborsh.Orderbook, error) {
	request := fmt.Sprintf(`{"jsonrpc": "2.0", "id": 1, "method": "GetOrderbook", "params": {"market":"%s"}"`, market)
	return unaryWSRequest[serumborsh.Orderbook](w.conn, request)
}

func (w *WSClient) GetOrderbookStream(ctx context.Context, market string, orderbookChan chan serumborsh.Orderbook) {
	request := fmt.Sprintf(`{"jsonrpc": "2.0", "id": 1, "method": "GetOrderbookStream", "params": {"market":"%s"}"`, market)
	unaryWSStream(ctx, w.conn, request, orderbookChan)
}

func unaryWSRequest[T any](conn *websocket.Conn, request string) (*T, error) {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(request)); err != nil {
		return nil, fmt.Errorf("error with sending message - %v", err)
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("error with reading message - %v", err)
	}

	var response T
	if err = json.Unmarshal(msg, &response); err != nil {
		return nil, fmt.Errorf("error with unmarshalling message of type %T - %v", response, err) // TODO check that response type is actually printed
	}

	return &response, nil
}

func unaryWSStream[T any](ctx context.Context, conn *websocket.Conn, request string, respChannel chan T) error {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(request)); err != nil {
		return fmt.Errorf("error with sending message - %v", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Log().Debug("stream closed", "remote addr", conn.RemoteAddr())
				return
			default:
				_, msg, err := conn.ReadMessage()
				if err != nil {
					logger.Log().Errorw("error with reading message - %v", "err", err)
					return
				}

				var response T
				if err = json.Unmarshal(msg, &response); err != nil {
					logger.Log().Errorw("error with unmarshalling message", "type", reflect.TypeOf(response), "err", err)
					continue
				}

				respChannel <- response
			}
		}
	}()

	return nil
}

func (w *WSClient) CloseConn() error {
	return w.conn.Close()
}
