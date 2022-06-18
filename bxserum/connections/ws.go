package connections

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sourcegraph/jsonrpc2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func WSRequest[T proto.Message](conn *websocket.Conn, request []byte, response T) error {
	err := sendWS(conn, request)
	if err != nil {
		return err
	}

	return recvWS[T](conn, response)
}

func WSStream[T proto.Message](ctx context.Context, conn *websocket.Conn, request []byte, responseChan chan T, response T) error {
	err := sendWS(conn, request)
	if err != nil {
		return err
	}

	err = recvWS[T](conn, response)
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
				err = recvWS[T](conn, response)
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
		return fmt.Errorf("error with sending message: %w", err)
	}
	return nil
}

func recvWS[T proto.Message](conn *websocket.Conn, result T) error {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("error reading WS response: %w", err)
	}

	// extract the result
	var resp jsonrpc2.Response
	if err = json.Unmarshal(msg, &resp); err != nil {
		return fmt.Errorf("error unmarshalling JSON response: %w", err)
	}
	if resp.Error != nil {
		m, err := json.Marshal(resp.Error.Data)
		if err != nil {
			return err
		}

		return errors.New(string(m))
	}

	if err = protojson.Unmarshal(*resp.Result, result); err != nil {
		fmt.Println(resp.Result)
		return fmt.Errorf("error unmarshalling message of type %T: %w", result, err)
	}

	return nil
}
