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
	handshakeTimeout       = 5 * time.Second
	subscriptionBuffer     = 1000
	unsubscribeGracePeriod = 3 * time.Second
)

type WS struct {
	messageM      sync.Mutex
	subscriptionM sync.RWMutex
	requestID     *utils.RequestID
	conn          *websocket.Conn
	ctx           context.Context
	cancel        context.CancelFunc
	err           error
	writeCh       chan []byte

	requestMap      map[uint64]requestTracker
	subscriptionMap map[string]subscriptionEntry
}

func NewWS(endpoint string) (*WS, error) {
	dialer := websocket.Dialer{HandshakeTimeout: handshakeTimeout}
	conn, _, err := dialer.Dial(endpoint, nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	ws := &WS{
		requestID:       utils.NewRequestID(),
		conn:            conn,
		ctx:             ctx,
		cancel:          cancel,
		writeCh:         make(chan []byte, 100),
		requestMap:      make(map[uint64]requestTracker),
		subscriptionMap: make(map[string]subscriptionEntry),
	}
	go ws.readLoop()
	go ws.writeLoop()
	return ws, nil
}

func (w *WS) readLoop() {
	defer w.cancel()

	for {
		// if lock is held by another processing routine, then wait for release before processing any further messages
		// typically, messages can be dispatched independently, but in some cases it may be important to finish processing a single message before allowing the socket to read any more (e.g. subscription response before processing potential updates)
		w.messageM.Lock()
		w.messageM.Unlock()

		_, msg, err := w.conn.ReadMessage()
		if err != nil {
			_ = w.Close(err)
			return
		}

		// try response format first
		var response jsonrpc2.Response
		err = json.Unmarshal(msg, &response)
		if err == nil && (response.Result != nil || response.Error != nil) {
			w.processRPCResponse(response)
			continue
		}

		// if not, try subscription format
		var update jsonrpc2.Request
		err = json.Unmarshal(msg, &update)
		if err == nil && update.Params != nil {
			w.processSubscriptionUpdate(update)
			continue
		}

		// no message works: exit loop and cancel connection
		_ = w.Close(fmt.Errorf("unknown jsonrpc message format: %v", string(msg)))
		return
	}
}

func (w *WS) writeLoop() {
	for {
		m := <-w.writeCh
		err := w.conn.WriteMessage(websocket.TextMessage, m)
		if err != nil {
			_ = w.Close(fmt.Errorf("error sending message: %w", err))
			return
		}
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
		w.messageM.Lock()
		ru.lockHeld = true
	}

	rt.ch <- ru
}

func (w *WS) processSubscriptionUpdate(update jsonrpc2.Request) {
	var f FeedUpdate
	err := json.Unmarshal(*update.Params, &f)
	if err != nil {
		_ = w.Close(fmt.Errorf("could not deserialize feed update: %w", err))
		return
	}

	w.subscriptionM.RLock()
	defer w.subscriptionM.RUnlock()
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
	params, err := protojson.Marshal(request)
	if err != nil {
		return err
	}
	rawParams := json.RawMessage(params)
	rpcRequest.Params = &rawParams

	rpcResponse, err := w.request(ctx, rpcRequest, false)
	if err != nil {
		return err
	}

	if err = protojson.Unmarshal(*rpcResponse.Result, response); err != nil {
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

	w.writeCh <- b

	select {
	case response := <-responseCh:
		rpcResponse := response.v
		if rpcResponse.Error != nil {
			var rpcErr string
			err = json.Unmarshal(*rpcResponse.Error.Data, &rpcErr)
			if err != nil {
				return jsonrpc2.Response{}, err
			}
			return rpcResponse, errors.New(rpcErr)
		}
		return response.v, nil
	case <-ctx.Done():
		return jsonrpc2.Response{}, ctx.Err()
	case <-w.ctx.Done():
		// connection closed
		return jsonrpc2.Response{}, fmt.Errorf("websocket connection was closed: %w", w.err)
	}
}

func WSStream[T proto.Message](w *WS, ctx context.Context, streamName string, streamParams proto.Message, resultInitFn func() T) (func() (T, error), error) {
	streamParamsB, err := protojson.Marshal(streamParams)
	if err != nil {
		return nil, err
	}
	params := SubscribeParams{
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
	defer w.messageM.Unlock()

	var subscriptionID string
	err = json.Unmarshal(*rpcResponse.Result, &subscriptionID)
	if err != nil {
		return nil, err
	}

	ch := make(chan json.RawMessage, subscriptionBuffer)
	streamCtx, streamCancel := context.WithCancel(ctx)

	w.subscriptionM.Lock()
	w.subscriptionMap[subscriptionID] = subscriptionEntry{
		active: true,
		ch:     ch,
		cancel: streamCancel,
	}
	w.subscriptionM.Unlock()

	// set goroutine to unsubscribe when ctx is canceled
	go func() {
		<-streamCtx.Done()

		// immediately mark as inactive
		w.subscriptionM.Lock()
		w.subscriptionMap[subscriptionID] = subscriptionEntry{active: false}
		w.subscriptionM.Unlock()

		up := UnsubscribeParams{SubscriptionID: subscriptionID}
		b, _ := json.Marshal(up)
		rm := json.RawMessage(b)

		unsubscribeMessage := jsonrpc2.Request{
			ID:     jsonrpc2.ID{Num: w.requestID.Next()},
			Method: unsubscribeMethod,
			Params: &rm,
		}

		_, err = w.request(w.ctx, unsubscribeMessage, false)
		if err != nil {
			_ = w.Close(fmt.Errorf("unsubscribe requested rejected: %w", err))
		}

		// wait for server to process message before forcing errors from unknown subscription IDs
		time.Sleep(unsubscribeGracePeriod)
		w.subscriptionM.Lock()
		delete(w.subscriptionMap, subscriptionID)
		w.subscriptionM.Unlock()
	}()

	return func() (T, error) {
		var zero T
		select {
		case b := <-ch:
			v := resultInitFn()
			err := protojson.Unmarshal(b, v)
			if err != nil {
				return zero, err
			}
			return v, nil
		case <-w.ctx.Done():
			return zero, fmt.Errorf("connection has been closed: %w", w.err)
		case <-streamCtx.Done():
			return zero, errors.New("stream context has been closed")
		}
	}, nil
}

func (w *WS) Close(reason error) error {
	w.messageM.Lock()
	defer w.messageM.Unlock()
	w.subscriptionM.Lock()
	defer w.subscriptionM.Unlock()

	if w.ctx.Err() != nil {
		return nil
	}

	w.err = reason

	// cancel main connection ctx
	w.cancel()

	// cancel all subscriptions
	for _, sub := range w.subscriptionMap {
		sub.close()
	}

	// close underlying connection
	return w.conn.Close()
}
