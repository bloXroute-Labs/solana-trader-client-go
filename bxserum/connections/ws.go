package connections

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/sourcegraph/jsonrpc2"
	"google.golang.org/protobuf/encoding/protojson"
)

// TODO Handle sending responses to their correct locations
type response struct {
	Result json.RawMessage
	Error  jsonrpc2.Error
}

func WSRequest[T proto.Message](conn *websocket.Conn, request []byte) (T, error) {
	err := sendWS(conn, request)
	if err != nil {
		return nil, err
	}

	return recvWS[T](conn)
}

func WSStream[T proto.Message](ctx context.Context, conn *websocket.Conn, request []byte, responseChan chan T) error {
	err := sendWS(conn, request)
	if err != nil {
		return err
	}

	response, err := recvWS[T](conn)
	if err != nil {
		return err
	}

	go func(response T, responseChan chan T, conn *websocket.Conn) {
		responseChan <- response

		for {
			select {
			case <-ctx.Done():
				return
			default:
				response, err = recvWS[T](conn)
				if err != nil {
					break
				}

				responseChan <- response
			}
		}
	}(response, responseChan, conn)

	return nil
}

func sendWS(conn *websocket.Conn, request []byte) error {
	if err := conn.WriteMessage(websocket.TextMessage, request); err != nil {
		return fmt.Errorf("error with sending message: %v", err)
	}
	return nil
}

func recvWS[T proto.Message](conn *websocket.Conn) (T, error) {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("error reading WS response: %v", err)
	}

	// extract the result
	var resp response
	if err = json.Unmarshal(msg, &resp); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON response: %v", err)
	}
	if resp.Error.Data != nil {
		m, err := json.Marshal(resp.Error.Data)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(string(m))
	}

	var result T
	if err = protojson.Unmarshal(resp.Result, result); err != nil {
		return nil, fmt.Errorf("error unmarshalling message of type %T: %v", result, err)
	}

	return result, nil
}
