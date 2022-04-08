package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bloXroute-Labs/serum-api/borsh/serumborsh"
	"github.com/bloXroute-Labs/serum-api/logger"
	"github.com/gorilla/websocket"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"reflect"
)

// want to be able to connect/disconnect to it using the address
// client should be able to give me results for api calls?
// maybe like
// g := NewGRPCClient(addr)
// g.GetOrderBookStream(market, channel)
// get a provider?

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
	command := fmt.Sprintf(`{"jsonrpc": "2.0", "id": 1, "method": "GetOrderbook", "params": {"market":"%s"}"`, market)
	msg, err := w.unaryWSRequest(command)
	if err != nil {
		return nil, err
	}

	var orderBook serumborsh.Orderbook
	if err := json.Unmarshal(msg, &orderBook); err != nil {
		return nil, fmt.Errorf("error with unmarshalling message - %v", err)
	}

	return &orderBook, nil
}

func (w *WSClient) GetOrderbookStream(ctx context.Context, market string, orderbookChan chan *serumborsh.Orderbook) {
	command := fmt.Sprintf(`{"jsonrpc": "2.0", "id": 1, "method": "GetOrderbookStream", "params": {"market":"%s"}"`, market)
	if err := w.conn.WriteMessage(websocket.TextMessage, []byte(command)); err != nil {
		logger.Log().Errorw("error with sending message", "err", err)
	}

	b := make(chan serumborsh.Orderbook)
	w.unaryWSStream(command, b)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case resp := <-b:
				orderbookChan <- &resp
			}
		}
	}()

	return nil
}

func (w *WSClient) unaryWSRequest[T any](request string) (interface{}, error) {
	if err := w.conn.WriteMessage(websocket.TextMessage, []byte(request)); err != nil {
		return nil, fmt.Errorf("error with sending message - %v", err)
	}

	_, msg, err := w.conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("error with reading message - %v", err)
	}

	var response T
	if err = json.Unmarshal(msg, &response); err != nil {
		return nil, fmt.Errorf("error with unmarshalling message of type %T - %v", response, err) // TODO check that response type is actually printed
	}

	return &response, nil
}

func (w *WSClient) unaryWSStream[T any](request string, respChannel T) {
	if err := w.conn.WriteMessage(websocket.TextMessage, []byte(request)); err != nil {
		logger.Log().Errorf("error with sending message - %v", err)
	}

	go func() {
		for {
			_, msg, err := w.conn.ReadMessage()
			if err != nil {
				logger.Log().Errorw("error with reading message - %v", "err", err)
				return
			}

			var response T
			if err = json.Unmarshal(msg, &response); err != nil {
				logger.Log().Errorw("error with unmarshalling message", "type", reflect.TypeOf(response), "err", err) // TODO check that response type is actually printed
				continue
			}

			respChannel <- response
		}
	}()
}

func (w *WSClient) CloseConn() error {
	return w.conn.Close()
}