package connections

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bloXroute-Labs/serum-api/utils"
	"github.com/gorilla/websocket"
	"github.com/sourcegraph/jsonrpc2"
	"sync"
)

type requests struct {
	sync.Mutex
	dataReceivers map[uint64]chan json.RawMessage
}

func (r *requests) add(id uint64, receiver chan json.RawMessage) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.dataReceivers[id] = receiver
}

func (r *requests) delete(id uint64) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	delete(r.dataReceivers, id)
}

type subscriptions struct {
	*requests
	subscriptions       map[string]chan json.RawMessage
	subscriptionsReader chan subscriptionResponse
}

func (s *subscriptions) read() {
	for response := range s.subscriptionsReader {
		_, ok := s.subscriptions[response.Result]
		if ok {
			continue
		}
		s.subscriptions[response.Result] = s.dataReceivers[response.ID]
		s.delete(response.ID)
	}
}

type subscriptionResponse struct {
	ID     uint64 `json:"id"`
	Result string `json:"result"`
}

type streamResponse struct {
	SubscriptionID string          `json:"subscriptionID"`
	Params         json.RawMessage `json:"params"`
}

type WSRPC struct {
	requestID utils.RequestID
	conn      *websocket.Conn

	streamSubscriptions subscriptions
	nonStreamRequests   requests
}

func NewWSRPC(conn *websocket.Conn) *WSRPC {
	wsRPC := &WSRPC{
		requestID: utils.NewRequestID(),
		conn:      conn,
		nonStreamRequests: requests{
			dataReceivers: map[uint64]chan json.RawMessage{},
		},
		streamSubscriptions: subscriptions{
			requests:            &requests{dataReceivers: map[uint64]chan json.RawMessage{}},
			subscriptions:       map[string]chan json.RawMessage{},
			subscriptionsReader: make(chan subscriptionResponse, 1000),
		},
	}

	go wsRPC.read()
	go wsRPC.streamSubscriptions.read()

	return wsRPC
}

func (w *WSRPC) Request(methodName string, message json.RawMessage) (json.RawMessage, error) {
	id := w.requestID.Next()

	payload, err := w.makeRPCRequest(id, methodName, &message)
	if err != nil {
		return nil, err
	}

	receiver := make(chan json.RawMessage)
	w.nonStreamRequests.add(id, receiver)
	defer w.nonStreamRequests.delete(id)

	err = w.sendWS(payload)
	if err != nil {
		return nil, err
	}

	response := <-receiver

	return response, nil
}

func (w *WSRPC) Stream(methodName string, message json.RawMessage, responses chan json.RawMessage) error {
	id := w.requestID.Next()

	payload, err := w.makeRPCRequest(id, methodName, &message)
	if err != nil {
		return err
	}
	w.streamSubscriptions.add(id, responses)
	return w.sendWS(payload)
}

func (w *WSRPC) Close() error {
	err := w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return fmt.Errorf("error writing close msg -  %v", err)
	}
	return nil
}

func (w *WSRPC) makeRPCRequest(requestID uint64, methodName string, message *json.RawMessage) ([]byte, error) {
	rpcMessage := jsonrpc2.Request{ID: jsonrpc2.ID{Num: requestID}, Method: methodName, Params: message}
	return rpcMessage.MarshalJSON()
}

func (w *WSRPC) sendWS(request []byte) error {
	if err := w.conn.WriteMessage(websocket.TextMessage, request); err != nil {
		return fmt.Errorf("error with sending message: %v", err)
	}
	return nil
}

func (w *WSRPC) read() error {
	for {
		_, msg, err := w.conn.ReadMessage()
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

			responseReader, ok := w.streamSubscriptions.subscriptions[subscription.SubscriptionID]
			if !ok {
				return nil
			}

			responseReader <- subscription.Params

			continue
		}

		if resp.ID.IsString {
			return errors.New("received invalid request ID")
		}

		// new subscription for feed
		subscriptionResponse := subscriptionResponse{}
		err = json.Unmarshal(msg, &subscriptionResponse)
		if err == nil {
			subscriptionResponse.ID = resp.ID.Num
			w.streamSubscriptions.subscriptionsReader <- subscriptionResponse
			continue
		}

		responseReader, ok := w.nonStreamRequests.dataReceivers[resp.ID.Num]
		if !ok {
			continue
		}

		responseReader <- *resp.Result
	}
}
