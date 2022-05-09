package connections

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sourcegraph/jsonrpc2"
	"log"
)

type SubscriptionResponse struct {
	ID     uint64 `json:"id"`
	Result string `json:"result"`
}

type streamResponse struct {
	SubscriptionID string          `json:"subscriptionID"`
	Params         json.RawMessage `json:"params"`
}

func WSRequest[T any](conn *websocket.Conn, request []byte, responseChan chan json.RawMessage) (*T, error) {
	err := sendWS(conn, request)
	if err != nil {
		return nil, err
	}
	resp := <-responseChan

	var result *T
	if err = json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling message of type %T: %v", result, err)
	}

	return result, nil
}

func WSStream[T any](ctx context.Context, conn *websocket.Conn, request []byte, subscriptionChan chan json.RawMessage, responseChan chan *T) error {
	err := sendWS(conn, request)
	if err != nil {
		return err
	}

	go func(responseChan chan *T, conn *websocket.Conn) {
		for {
			select {
			case <-ctx.Done():
				return
			case resp := <-subscriptionChan:
				var result T
				if err = json.Unmarshal(resp, &result); err != nil {
					return
				}
				log.Println(result)
				responseChan <- &result
			}
		}
	}(responseChan, conn)

	return nil
}

func sendWS(conn *websocket.Conn, request []byte) error {
	if err := conn.WriteMessage(websocket.TextMessage, request); err != nil {
		return fmt.Errorf("error with sending message: %v", err)
	}
	return nil
}

func RecvWS(conn *websocket.Conn, nonStreamResponsesReaders map[uint64]chan json.RawMessage, streamResponsesReaders map[string]chan json.RawMessage, subscriptionsReader chan SubscriptionResponse) error {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return fmt.Errorf("error reading WS response: %v", err)
	}

	resp := jsonrpc2.Response{}

	if err = json.Unmarshal(msg, &resp); err != nil {
		return fmt.Errorf("error unmarshalling JSON response: %v", err)
	}

	if resp.Error != nil {
		if resp.Error.Data != nil {
			m, err := json.Marshal(resp.Error.Data)
			if err != nil {
				return err
			}

			return errors.New(string(m))
		}

		return errors.New("failed to check ")
	}

	if resp.ID.Str == "null" {
		subscription := streamResponse{}
		err := json.Unmarshal(*resp.Result, &subscription)
		if err != nil {
			return err
		}

		if subscription.SubscriptionID == "" {
			return errors.New("received invalid subscription ID")
		}

		responseReader, ok := streamResponsesReaders[subscription.SubscriptionID]
		if !ok {
			return nil
		}

		responseReader <- subscription.Params

		return nil
	}

	if resp.ID.IsString {
		return errors.New("received invalid request ID")
	}

	// new subscription for feed
	subscriptionResponse := SubscriptionResponse{}
	err = json.Unmarshal(msg, &subscriptionResponse)
	if err == nil {
		subscriptionResponse.ID = resp.ID.Num
		subscriptionsReader <- subscriptionResponse
		return nil
	}

	responseReader, ok := nonStreamResponsesReaders[resp.ID.Num]
	if !ok {
		return nil
	}

	responseReader <- *resp.Result

	return nil
}
