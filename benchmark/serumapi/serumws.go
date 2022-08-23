package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/bloXroute-Labs/serum-client-go/benchmark/internal/arrival"
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

type serumUpdate struct {
	asks []*pb.OrderbookItem
	bids []*pb.OrderbookItem
}

type serumOrderbookStream struct {
	wsConn  *websocket.Conn
	address string
	market  string
}

func newSerumOrderbookStream(address, market, authHeader string) (arrival.Source[[]byte, serumUpdate], error) {
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

// Run stops when parent ctx is canceled
func (s serumOrderbookStream) Run(parent context.Context) ([]arrival.StreamUpdate[[]byte], error) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	subscribeRequest := fmt.Sprintf(`{"jsonrpc": "2.0", "id": 1, "method": "subscribe", "params": ["GetOrderbooksStream", {"markets": ["%v"]}]}`, s.market)
	err := s.wsConn.WriteMessage(websocket.TextMessage, []byte(subscribeRequest))
	if err != nil {
		return nil, err
	}

	s.log().Debugw("subscription created")

	wsMessages := make(chan arrival.StreamUpdate[[]byte], 100)
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

			wsMessages <- arrival.NewStreamUpdate(b)
		}
	}()

	messages := make([]arrival.StreamUpdate[[]byte], 0)
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

func (s serumOrderbookStream) Process(updates []arrival.StreamUpdate[[]byte]) (map[int][]arrival.ProcessedUpdate[serumUpdate], error) {
	var err error

	results := make(map[int][]arrival.ProcessedUpdate[serumUpdate])
	allowedFailures := 1 // allowed to skip processing of subscription confirmation message

	for _, update := range updates {
		var rpcUpdate jsonrpc2.Request
		err = json.Unmarshal(update.Data, &rpcUpdate)
		if err != nil || rpcUpdate.Params == nil {
			allowedFailures--
			if allowedFailures < 0 {
				return nil, errors.Wrap(err, "too many response errors")
			}
			continue
		}

		var subU subscriptionUpdate
		err = json.Unmarshal(*rpcUpdate.Params, &subU)
		if err != nil {
			allowedFailures--
			if allowedFailures < 0 {
				return nil, errors.Wrap(err, "too many response errors")
			}
			continue
		}

		// note for future: when WS stream follows RPC spec will need to discard subscribe message
		var orderbookInc pb.GetOrderbooksStreamResponse
		err = protojson.Unmarshal(subU.Result, &orderbookInc)
		if err != nil {
			return nil, err
		}

		slot := int(orderbookInc.Slot)
		_, ok := results[slot]
		if !ok {
			results[slot] = make([]arrival.ProcessedUpdate[serumUpdate], 0)
		}

		su := serumUpdate{
			asks: orderbookInc.Orderbook.Asks,
			bids: orderbookInc.Orderbook.Bids,
		}
		results[slot] = append(results[slot], arrival.ProcessedUpdate[serumUpdate]{
			Timestamp: update.Timestamp,
			Slot:      slot,
			Data:      su,
		})
	}

	return results, nil
}
