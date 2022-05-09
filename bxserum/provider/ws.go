package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bloXroute-Labs/serum-api/bxserum/connections"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/bloXroute-Labs/serum-api/utils"
	"github.com/gorilla/websocket"
	"github.com/sourcegraph/jsonrpc2"
	"log"
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
	subscriptionsReader chan connections.SubscriptionResponse
}

func (s *subscriptions) read() {
	for response := range s.subscriptionsReader {
		log.Println(response)
		_, ok := s.subscriptions[response.Result]
		if ok {
			continue
		}
		s.subscriptions[response.Result] = s.dataReceivers[response.ID]
		s.delete(response.ID)
	}
}

type WSClient struct {
	pb.UnimplementedApiServer

	addr                string
	conn                *websocket.Conn
	requestID           utils.RequestID
	requests            requests
	streamSubscriptions subscriptions
}

// NewWSClient connects to Mainnet Serum API
func NewWSClient() (*WSClient, error) {
	return NewWSClientWithEndpoint("ws://174.129.154.164:1810/ws")
}

// NewWSClientTestnet connects to Testnet Serum API
func NewWSClientTestnet() (*WSClient, error) {
	panic("implement me")
}

// NewWSClientWithEndpoint connects to custom Serum API
func NewWSClientWithEndpoint(addr string) (*WSClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}

	client := &WSClient{
		addr:      addr,
		conn:      conn,
		requestID: utils.NewRequestID(),
		requests: requests{
			dataReceivers: map[uint64]chan json.RawMessage{},
		},
		streamSubscriptions: subscriptions{
			requests:            &requests{dataReceivers: map[uint64]chan json.RawMessage{}},
			subscriptions:       map[string]chan json.RawMessage{},
			subscriptionsReader: make(chan connections.SubscriptionResponse, 1000),
		},
	}

	go client.streamSubscriptions.read()
	go client.read()

	return client, nil
}

// Set limit to 0 to get all bids/asks
func (w *WSClient) GetOrderbook(market string, limit uint32) (*pb.GetOrderbookResponse, error) {
	request, requestID, err := w.jsonRPCRequest("GetOrderbook", map[string]interface{}{"market": market, "limit": limit})
	if err != nil {
		return nil, err
	}

	responseChan := make(chan json.RawMessage)
	w.requests.add(requestID, responseChan)
	defer w.requests.delete(requestID)

	return connections.WSRequest[pb.GetOrderbookResponse](w.conn, request, responseChan)
}

func (w *WSClient) GetOrderbookStream(ctx context.Context, market string, limit uint32, orderbookChan chan *pb.GetOrderbookStreamResponse) error {
	request, requestID, err := w.jsonRPCRequest("GetOrderbookStream", map[string]interface{}{"market": market, "limit": limit})
	if err != nil {
		return err
	}

	responseChan := make(chan json.RawMessage)
	w.streamSubscriptions.add(requestID, responseChan)

	return connections.WSStream[pb.GetOrderbookStreamResponse](ctx, w.conn, request, responseChan, orderbookChan)
}

func (w *WSClient) GetMarkets() (*pb.GetMarketsResponse, error) {
	request, requestID, err := w.jsonRPCRequest("GetMarkets", nil)
	if err != nil {
		return nil, err
	}

	responseChan := make(chan json.RawMessage)
	w.requests.add(requestID, responseChan)
	defer w.requests.delete(requestID)

	return connections.WSRequest[pb.GetMarketsResponse](w.conn, request, responseChan)
}

// TODO: Add Unsubscribe for each subscription
func (w *WSClient) Close() error {
	err := w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return fmt.Errorf("error writing close msg -  %v", err)
	}
	return nil
}

func (w *WSClient) jsonRPCRequest(method string, params map[string]interface{}) ([]byte, uint64, error) {
	id := w.requestID.Next()
	req := jsonrpc2.Request{
		Method: method,
		ID:     jsonrpc2.ID{Num: id},
	}
	if err := req.SetParams(params); err != nil {
		return nil, 0, err
	}
	payload, err := req.MarshalJSON()
	if err != nil {
		return nil, 0, err
	}

	return payload, id, nil
}

func (w *WSClient) read() {
	for {
		err := connections.RecvWS(w.conn, w.requests.dataReceivers, w.streamSubscriptions.subscriptions, w.streamSubscriptions.subscriptionsReader)
		if err != nil {
			break
		}
	}
}
