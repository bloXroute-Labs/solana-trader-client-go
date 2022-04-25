package connections

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/sourcegraph/jsonrpc2"
)

// TODO Handle sending responses to their correct locations
type response struct {
	Result json.RawMessage
	Error  jsonrpc2.Error
}

func WSRequest[T any](connectionManager *ConnectionManager, request []byte) (*T, error) {
	conn, id, err := connectionManager.Next()
	if err != nil {
		return nil, err
	}
	defer func(connectionManager *ConnectionManager, id int) {
		err := connectionManager.CloseConn(id)
		if err != nil {
			log.Errorf("connection with id %v not closed correctly - %s", id, err.Error())
		}
	}(connectionManager, id)

	err = sendWSRequest(conn, request)
	if err != nil {
		return nil, err
	}

	return recvWSResult[T](conn)
}

func WSStream[T any](ctx context.Context, connectionManager *ConnectionManager, request []byte, responseChan chan *T) error {
	conn, id, err := connectionManager.Next()
	if err != nil {
		return err
	}

	err = sendWSRequest(conn, request)
	log.Infof("WS Stream Request: %v Error: %v", string(request), err)
	if err != nil {
		return err
	}

	response, err := recvWSResult[T](conn)
	if err != nil {
		log.Errorf("error in ws stream %v", err)
		return err
	}
	responseChan <- response

	go func(responseChan chan *T, conn *websocket.Conn, id int) {
		defer func() {
			err := connectionManager.CloseConn(id)
			if err != nil {
				log.Errorf("connection with id %v not closed correctly - %s", id, err.Error())
			}
		}()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				response, err := recvWSResult[T](conn)
				if err != nil {
					log.Error(err)
					break
				}

				responseChan <- response
			}
		}
	}(responseChan, conn, id)

	return nil
}

func sendWSRequest(conn *websocket.Conn, request []byte) error {
	if err := conn.WriteMessage(websocket.TextMessage, request); err != nil {
		return fmt.Errorf("error with sending message - %v", err)
	}
	return nil
}

func recvWSResult[T any](conn *websocket.Conn) (*T, error) {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("error reading WS response - %v", err)
	}
	// extract the HTTP Response Result
	var resp response
	if err = json.Unmarshal(msg, &resp); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON response - %v", err)
	}
	if resp.Error.Data != nil {
		m, err := json.Marshal(resp.Error.Data)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(string(m))
	}

	var result T
	if err = json.Unmarshal(resp.Result, &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling message of type %T - %v", result, err)
	}
	return &result, nil
}
