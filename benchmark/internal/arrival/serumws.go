package arrival

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/bloXroute-Labs/serum-client-go/benchmark/internal/logger"
	pb "github.com/bloXroute-Labs/serum-client-go/proto"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sourcegraph/jsonrpc2"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"net/http"
)

type SerumUpdate struct {
	Slot     int
	Asks     []*pb.OrderbookItem
	Bids     []*pb.OrderbookItem
	previous *SerumUpdate
}

func (s SerumUpdate) IsRedundant() bool {
	if s.previous == nil {
		return false
	}
	return orderbookEqual(s.previous.Bids, s.Bids) && orderbookEqual(s.previous.Asks, s.Asks)
}

type serumOrderbookStream struct {
	wsConn  *websocket.Conn
	address string
	market  string
}

func NewSerumOrderbookStream(address, market, authHeader string) (Source[[]byte, SerumUpdate], error) {
	s := serumOrderbookStream{
		address: address,
		market:  market,
	}

	// ws provider not used to delay message deserialization until all complete
	dialer := websocket.Dialer{TLSClientConfig: &tls.Config{}}
	header := http.Header{}
	header.Set("Authorization", authHeader)
	wsConn, resp, err := dialer.Dial(s.address, header)
	if err != nil {
		return nil, err
	}

	s.log().Debugw("connection established")

	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(resp.Body)
	s.wsConn = wsConn
	return s, nil
}

func (s serumOrderbookStream) log() *zap.SugaredLogger {
	return logger.Log().With("source", "serum", "address", s.address, "market", s.market)
}

func (s serumOrderbookStream) Name() string {
	return fmt.Sprintf("serum[%v]", s.address)
}

// Run stops when parent ctx is canceled
func (s serumOrderbookStream) Run(parent context.Context) ([]StreamUpdate[[]byte], error) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	subscribeRequest := fmt.Sprintf(`{"jsonrpc": "2.0", "id": 1, "method": "subscribe", "params": ["GetOrderbooksStream", {"markets": ["%v"]}]}`, s.market)
	err := s.wsConn.WriteMessage(websocket.TextMessage, []byte(subscribeRequest))
	if err != nil {
		return nil, err
	}

	s.log().Debugw("subscription created")

	wsMessages := make(chan StreamUpdate[[]byte], 100)
	go func() {
		for {
			if ctx.Err() != nil {
				return
			}

			_, b, err := s.wsConn.ReadMessage()
			if err != nil {
				s.log().Debugw("closing connection", "err", err)
				return
			}

			wsMessages <- NewStreamUpdate(b)
		}
	}()

	messages := make([]StreamUpdate[[]byte], 0)
	for {
		select {
		case msg := <-wsMessages:
			messages = append(messages, msg)
		case <-ctx.Done():
			err = s.wsConn.Close()
			if err != nil {
				s.log().Errorw("could not close connection", "err", err)
			}
			return messages, nil
		}
	}
}

type subscriptionUpdate struct {
	SubscriptionID string          `json:"subscriptionId"`
	Result         json.RawMessage `json:"result"`
}

func (s serumOrderbookStream) Process(updates []StreamUpdate[[]byte], removeDuplicates bool) (map[int][]ProcessedUpdate[SerumUpdate], map[int][]ProcessedUpdate[SerumUpdate], error) {
	var (
		err      error
		previous *SerumUpdate
	)

	results := make(map[int][]ProcessedUpdate[SerumUpdate])
	duplicates := make(map[int][]ProcessedUpdate[SerumUpdate])
	allowedFailures := 1 // allowed to skip processing of subscription confirmation message

	for _, update := range updates {
		var rpcUpdate jsonrpc2.Request
		err = json.Unmarshal(update.Data, &rpcUpdate)
		if err != nil || rpcUpdate.Params == nil {
			allowedFailures--
			if allowedFailures < 0 {
				return nil, nil, errors.Wrap(err, "too many response errors")
			}
			continue
		}

		var subU subscriptionUpdate
		err = json.Unmarshal(*rpcUpdate.Params, &subU)
		if err != nil {
			allowedFailures--
			if allowedFailures < 0 {
				return nil, nil, errors.Wrap(err, "too many response errors")
			}
			continue
		}

		var orderbookInc pb.GetOrderbooksStreamResponse
		err = protojson.Unmarshal(subU.Result, &orderbookInc)
		if err != nil {
			return nil, nil, err
		}

		slot := int(orderbookInc.Slot)
		_, ok := results[slot]
		if !ok {
			results[slot] = make([]ProcessedUpdate[SerumUpdate], 0)
		}

		su := SerumUpdate{
			Slot:     slot,
			Asks:     orderbookInc.Orderbook.Asks,
			Bids:     orderbookInc.Orderbook.Bids,
			previous: previous,
		}
		pu := ProcessedUpdate[SerumUpdate]{
			Timestamp: update.Timestamp,
			Slot:      slot,
			Data:      su,
		}

		if su.IsRedundant() {
			duplicates[slot] = append(results[slot], pu)
		} else {
			previous = &su
		}

		if !removeDuplicates || !su.IsRedundant() {
			results[slot] = append(results[slot], pu)
		}
	}

	return results, duplicates, nil
}
