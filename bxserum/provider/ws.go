package provider

import (
	"context"
	"encoding/json"
	"github.com/bloXroute-Labs/serum-api/bxserum/connections"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/gorilla/websocket"
)

type WsRPC interface {
	Request(methodName string, message json.RawMessage) (json.RawMessage, error)
	Stream(methodName string, message json.RawMessage, responses chan json.RawMessage) error
	Close() error
}

type WSClient struct {
	pb.UnimplementedApiServer

	connection WsRPC
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
		connection: connections.NewWSRPC(conn),
	}

	return client, nil
}

// Set limit to 0 to get all bids/asks
func (w *WSClient) GetOrderbook(market string, limit uint32) (*pb.GetOrderbookResponse, error) {
	message, err := json.Marshal(map[string]interface{}{"market": market, "limit": limit})
	if err != nil {
		return nil, err
	}

	result, err := w.connection.Request("GetOrderbook", message)
	if err != nil {
		return nil, err
	}

	return unmarshalJSON[pb.GetOrderbookResponse](result)
}

func (w *WSClient) GetOrderbookStream(ctx context.Context, market string, limit uint32, orderbookChan chan *pb.GetOrderbookStreamResponse) error {
	message, err := json.Marshal(map[string]interface{}{"market": market, "limit": limit})
	if err != nil {
		return err
	}

	response := make(chan json.RawMessage)

	err = w.connection.Stream("GetOrderbookStream", message, response)
	if err != nil {
		return err
	}

	go readStreamResponse[pb.GetOrderbookStreamResponse](ctx, response, orderbookChan)

	return nil
}

func (w *WSClient) GetMarkets() (*pb.GetMarketsResponse, error) {
	result, err := w.connection.Request("GetMarkets", nil)
	if err != nil {
		return nil, err
	}

	return unmarshalJSON[pb.GetMarketsResponse](result)
}

// TODO: Add Unsubscribe for each subscription
func (w *WSClient) Close() error {
	return w.connection.Close()
}

func readStreamResponse[T any](ctx context.Context, streamChan chan json.RawMessage, clientChan chan *T) {
	for {
		select {
		case <-ctx.Done():
			return
		case rawMessage := <-streamChan:
			result, err := unmarshalJSON[T](rawMessage)
			if err != nil {
				continue
			}
			clientChan <- result
		}
	}
}

func unmarshalJSON[T any](message json.RawMessage) (*T, error) {
	var result T
	err := json.Unmarshal(message, &result)
	return &result, err
}
