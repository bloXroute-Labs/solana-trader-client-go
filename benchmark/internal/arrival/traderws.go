package arrival

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/benchmark/internal/logger"
	pb "github.com/bloXroute-Labs/solana-trader-client-go/proto"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sourcegraph/jsonrpc2"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"net/http"
)

type TraderAPIUpdate struct {
	Asks     []*pb.OrderbookItem
	Bids     []*pb.OrderbookItem
	previous *TraderAPIUpdate
}

func (s TraderAPIUpdate) IsRedundant() bool {
	if s.previous == nil {
		return false
	}
	return orderbookEqual(s.previous.Bids, s.Bids) && orderbookEqual(s.previous.Asks, s.Asks)
}

type apiOrderbookStream struct {
	wsConn  *websocket.Conn
	address string
	market  string
}

func NewAPIOrderbookStream(address, market, authHeader string) (Source[[]byte, TraderAPIUpdate], error) {
	s := apiOrderbookStream{
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

func (s apiOrderbookStream) log() *zap.SugaredLogger {
	return logger.Log().With("source", "traderapi", "address", s.address, "market", s.market)
}

func (s apiOrderbookStream) Name() string {
	return fmt.Sprintf("traderapi[%v]", s.address)
}

// Run stops when parent ctx is canceled
func (s apiOrderbookStream) Run(parent context.Context) ([]StreamUpdate[[]byte], error) {
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

func (s apiOrderbookStream) Process(updates []StreamUpdate[[]byte], removeDuplicates bool) (map[int][]ProcessedUpdate[TraderAPIUpdate], map[int][]ProcessedUpdate[TraderAPIUpdate], error) {
	var (
		err      error
		previous *TraderAPIUpdate
	)

	results := make(map[int][]ProcessedUpdate[TraderAPIUpdate])
	duplicates := make(map[int][]ProcessedUpdate[TraderAPIUpdate])
	allowedFailures := 1 // allowed to skip processing of subscription confirmation message

	for _, update := range updates {
		var rpcUpdate jsonrpc2.Request
		err = json.Unmarshal(update.Data, &rpcUpdate)
		if err != nil || rpcUpdate.Params == nil {
			allowedFailures--
			if allowedFailures < 0 {
				return nil, nil, errors.Wrap(err, "too many response errors")
			}

			var subscriptionResponse jsonrpc2.Response
			err := json.Unmarshal(update.Data, &subscriptionResponse)
			if err != nil {
				return nil, nil, fmt.Errorf("did not receive proper subscription response: %v", string(update.Data))
			}

			if subscriptionResponse.Error != nil {
				return nil, nil, fmt.Errorf("did not receive proper subscription response: %v", string(update.Data))
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
		su := TraderAPIUpdate{
			Asks:     orderbookInc.Orderbook.Asks,
			Bids:     orderbookInc.Orderbook.Bids,
			previous: previous,
		}
		pu := ProcessedUpdate[TraderAPIUpdate]{
			Timestamp: update.Timestamp,
			Slot:      slot,
			Data:      su,
		}

		redundant := su.IsRedundant()
		if redundant {
			duplicates[slot] = append(results[slot], pu)
		} else {
			previous = &su
		}

		// skip redundant updates if duplicate updates flag is set
		if !(removeDuplicates && redundant) {
			_, ok := results[slot]
			if !ok {
				results[slot] = make([]ProcessedUpdate[TraderAPIUpdate], 0)
			}
			results[slot] = append(results[slot], pu)
		}
	}

	return results, duplicates, nil
}
