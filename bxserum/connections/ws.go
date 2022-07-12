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
	"time"
)

const (
	subscriptionBuffer     = 1000
	unsubscribeGracePeriod = 3 * time.Second
)

type WS struct {
	m         sync.Mutex
	requestID utils.RequestID
	conn      *websocket.Conn
	ctx       context.Context
	cancel    context.CancelFunc
	err       error

	requestMap      map[uint64]requestTracker
	subscriptionMap map[string]subscriptionEntry
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
	defer w.cancel()

	for {
		// if lock is held by another processing routine, then wait for release before processing any further messages
		// typically, messages can be dispatched independently, but in some cases it may be important to finish processing a single message before allowing the socket to read any more (e.g. subscription response before processing potential updates)
		w.m.Lock()
		w.m.Unlock()

		_, msg, err := w.conn.ReadMessage()
		if err != nil {
			_ = w.Close(err)
			return
		}

		// try response format first
		var response jsonrpc2.Response
		err = json.Unmarshal(msg, &response)
		if err == nil {
			w.processRPCResponse(response)
			continue
		}

		// if not, try subscription format
		var update jsonrpc2.Request
		err = json.Unmarshal(msg, &update)
		if err == nil {
			w.processSubscriptionUpdate(update)
			continue
		}

		// no message works: exit loop and cancel connection
		_ = w.Close(fmt.Errorf("unknown jsonrpc message format: %v", string(msg)))
		return
	}
}

func (w *WS) processRPCResponse(response jsonrpc2.Response) {
	requestID := response.ID.Num
	rt, ok := w.requestMap[requestID]
	if !ok {
		_ = w.Close(fmt.Errorf("unknown request ID: got %v, most recent %v", requestID, w.requestID.Current()))
		return
	}

	ru := responseUpdate{
		v:        response,
		lockHeld: false,
	}

	// hold lock: now it's the responsibility of the listening channel to release the lock for the next loop
	if rt.lockRequired {
		w.m.Lock()
		ru.lockHeld = true
	}

	rt.ch <- ru
}

func (w *WS) processSubscriptionUpdate(update jsonrpc2.Request) {
	var f feedUpdate
	err := json.Unmarshal(*update.Params, &f)
	if err != nil {
		_ = w.Close(fmt.Errorf("could not deserialize feed update: %w", err))
		return
	}

	sub, ok := w.subscriptionMap[f.SubscriptionID]
	if !ok {
		_ = w.Close(fmt.Errorf("unknown subscription ID: %v", f.SubscriptionID))
		return
	}
	// skip message for inactive subscription: will be closed soon
	if !sub.active {
		return
	}

	sub.ch <- f.Result
}

func (w *WS) Request(ctx context.Context, method string, request proto.Message, response proto.Message) error {
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

	rpcResponse, err := w.request(ctx, rpcRequest, false)
	if err != nil {
		return err
	}

	if rpcResponse.Error != nil {
		var rpcErr string
		err = json.Unmarshal(*rpcResponse.Error.Data, &rpcErr)
		if err != nil {
			return err
		}
		return errors.New(rpcErr)
	}

	if err = protojson.Unmarshal(*rpcRequest.Params, response); err != nil {
		return fmt.Errorf("error unmarshalling message of type %T: %w", response, err)
	}
	return nil
}

func (w *WS) request(ctx context.Context, request jsonrpc2.Request, lockRequired bool) (jsonrpc2.Response, error) {
	b, err := json.Marshal(request)
	if err != nil {
		return jsonrpc2.Response{}, err
	}

	// setup listener for next request ID that matches response
	responseCh := make(chan responseUpdate)
	w.requestMap[request.ID.Num] = requestTracker{
		ch:           responseCh,
		lockRequired: lockRequired,
	}
	defer func() {
		delete(w.requestMap, request.ID.Num)
	}()

	err = w.conn.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return jsonrpc2.Response{}, fmt.Errorf("error with sending message: %w", err)
	}

	select {
	case response := <-responseCh:
		return response.v, nil
	case <-ctx.Done():
		return jsonrpc2.Response{}, ctx.Err()
	case <-w.ctx.Done():
		// connection closed
		return jsonrpc2.Response{}, fmt.Errorf("websocket connection was closed: %w", w.err)
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

	// requires lock held on subscription mutex, otherwise a subscription message could be processed before the map entry is created
	rpcResponse, err := w.request(ctx, rpcRequest, true)
	if err != nil {
		return nil, err
	}
	defer w.m.Unlock()

	var subscriptionID string
	err = json.Unmarshal(*rpcResponse.Result, &subscriptionID)
	if err != nil {
		return nil, err
	}

	ch := make(chan json.RawMessage, subscriptionBuffer)
	streamCtx, streamCancel := context.WithCancel(ctx)
	w.subscriptionMap[subscriptionID] = subscriptionEntry{
		active: true,
		ch:     ch,
		cancel: streamCancel,
	}

	// set goroutine to unsubscribe when ctx is canceled
	go func() {
		<-streamCtx.Done()

		// immediately mark as inactive
		w.subscriptionMap[subscriptionID] = subscriptionEntry{active: false}

		up := unsubscribeParams{SubscriptionID: subscriptionID}
		b, _ := json.Marshal(up)
		rm := json.RawMessage(b)

		unsubscribeMessage := jsonrpc2.Request{
			ID:     jsonrpc2.ID{Num: w.requestID.Next()},
			Method: unsubscribeMethod,
			Params: &rm,
		}

		_, _ = w.request(w.ctx, unsubscribeMessage, false)

		// wait for server to process message before forcing errors from unknown subscription IDs
		time.Sleep(unsubscribeGracePeriod)
		delete(w.subscriptionMap, subscriptionID)
	}()

	return func() (T, error) {
		select {
		case b := <-ch:
			v := resultInitFn()
			err := proto.Unmarshal(b, v)
			if err != nil {
				return nil, err
			}
			return v, nil
		case <-w.ctx.Done():
			return nil, fmt.Errorf("connection has been closed: %w", w.err)
		case <-streamCtx.Done():
			return nil, errors.New("stream context has been closed")
		}
	}, nil
}

func (w *WS) Close(reason error) error {
	w.err = reason

	// cancel main connection ctx
	w.cancel()

	// cancel  all subscriptions
	for _, sub := range w.subscriptionMap {
		sub.close()
	}

	// close underlying connection
	return w.conn.Close()
}
