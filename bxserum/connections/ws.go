package connections

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bloXroute-Labs/serum-client-go/utils"
	"github.com/gorilla/websocket"
	"github.com/sourcegraph/jsonrpc2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"sync"
)

const subscriptionBuffer = 1000

type subscription struct {
	ch     chan json.RawMessage
	cancel context.CancelFunc
}

type WS struct {
	requestID utils.RequestID
	conn      *websocket.Conn
	ctx       context.Context
	cancel    context.CancelFunc

	requestMap map[uint64]chan jsonrpc2.Response

	subscriptionM   sync.Mutex
	subscriptionMap map[string]subscription
}

func NewWS(endpoint string) (*WS, error) {
	conn, _, err := websocket.DefaultDialer.Dial(endpoint, nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	ws := &WS{
		requestID: utils.NewRequestID(),
		conn:      conn,
		ctx:       ctx,
		cancel:    cancel,
	}
	go ws.readLoop()
	return ws, nil
}

func (w *WS) readLoop() {
	for {
		_, msg, err := w.conn.ReadMessage()
		if err != nil {
			w.cancel()
			return
		}

		// try JSONRPC response format first
		var response jsonrpc2.Response
		err = json.Unmarshal(msg, &response)
		if err == nil {
			w.processRPCResponse(response)
			continue
		}

		// if not, try subscription update
		var update jsonrpc2.Request
		err = json.Unmarshal(msg, &update)
		if err == nil {
			w.processSubscriptionUpdate(update)
			continue
		}

		// TODO: no message format works, not good!
		w.cancel()
		return
	}
}

func (w *WS) processRPCResponse(response jsonrpc2.Response) {
	requestID := response.ID.Num
	responseCh, ok := w.requestMap[requestID]
	if !ok {
		// TODO what do?
		return
	}
	// TODO special handling for subscription responses?
	responseCh <- response
}

func (w *WS) processSubscriptionUpdate(update jsonrpc2.Request) {
	w.subscriptionM.Lock()
	defer w.subscriptionM.Unlock()

	var f feedUpdate
	err := json.Unmarshal(*update.Params, &f)
	if err != nil {
		return
	}

	subscription, ok := w.subscriptionMap[f.SubscriptionID]
	if !ok {
		// TODO
		return
	}

	subscription.ch <- f.Result
}

func (w *WS) Request(method string, request proto.Message, response proto.Message) error {
	requestID := w.requestID.Next()
	rpcRequest := jsonrpc2.Request{
		Method: method,
		ID:     jsonrpc2.ID{Num: requestID},
	}
	params, err := proto.Marshal(request)
	if err != nil {
		return err
	}
	rawParams := json.RawMessage(params)
	rpcRequest.Params = &rawParams

	rpcResponse, err := w.request(rpcRequest)
	if err != nil {
		return err
	}

	if rpcResponse.Error != nil {
		m, err := json.Marshal(rpcResponse.Error.Data)
		if err != nil {
			return err
		}
		return errors.New(string(m))
	}

	if err = protojson.Unmarshal(*rpcRequest.Params, response); err != nil {
		return fmt.Errorf("error unmarshalling message of type %T: %w", response, err)
	}
	return nil
}

func (w *WS) request(request jsonrpc2.Request) (jsonrpc2.Response, error) {
	b, err := json.Marshal(request)
	if err != nil {
		return jsonrpc2.Response{}, err
	}

	// setup listener for next request ID that matches response
	responseCh := make(chan jsonrpc2.Response)
	w.requestMap[request.ID.Num] = responseCh
	defer func() {
		delete(w.requestMap, request.ID.Num)
	}()

	err = w.conn.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return jsonrpc2.Response{}, fmt.Errorf("error with sending message: %w", err)
	}

	select {
	case response := <-responseCh:
		return response, nil
	case <-w.ctx.Done():
		// connection closed
		return jsonrpc2.Response{}, errors.New("websocket connection was closed before response received")
	}
}

func WSStream[T proto.Message](w *WS, ctx context.Context, streamName string, streamParams proto.Message, resultInitFn func() T) (func() (T, error), error) {
	streamParamsB, err := proto.Marshal(streamParams)
	if err != nil {
		return nil, err
	}
	params := subscribeParams{
		StreamName: streamName,
		StreamOpts: streamParamsB,
	}
	paramsB, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	rawParams := json.RawMessage(paramsB)
	rpcRequest := jsonrpc2.Request{
		Method: subscribeMethod,
		ID:     jsonrpc2.ID{Num: w.requestID.Next()},
		Params: &rawParams,
	}

	rpcResponse, err := w.request(rpcRequest)
	if err != nil {
		return nil, err
	}

	var subscriptionID string
	err = json.Unmarshal(*rpcResponse.Result, &subscriptionID)
	if err != nil {
		return nil, err
	}

	// todo probably causes a deadlock...
	w.subscriptionM.Lock()
	defer w.subscriptionM.Unlock()

	ch := make(chan json.RawMessage, subscriptionBuffer)
	streamCtx, streamCancel := context.WithCancel(ctx)
	w.subscriptionMap[subscriptionID] = subscription{
		ch:     ch,
		cancel: streamCancel,
	}

	return func() (T, error) {
		select {
		case b := <-ch:
			v := resultInitFn()
			err := proto.Unmarshal(b, v)
			if err != nil {
				return nil, err
			}
			return v, nil
		case <-streamCtx.Done():
			return nil, errors.New("channel closed")
		}
	}, nil
}

func (w *WS) Close() error {
	w.cancel()
	return w.conn.Close()
}
