package provider

import (
	"context"
	"github.com/bloXroute-Labs/serum-api/bxserum/connections"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/bloXroute-Labs/serum-api/utils"
	"github.com/sourcegraph/jsonrpc2"
)

type WSClient struct {
	pb.UnimplementedApiServer

	addr              string
	connectionManager *connections.ConnectionManager
	requestID         utils.RequestID
}

// Connects to Mainnet Serum API
func NewWSClient() *WSClient {
	return NewWSClientWithEndpoint("ws://174.129.154.164:1810/ws")
}

// Connects to Testnet Serum API
func NewWSClientTestnet() (*WSClient, error) {
	panic("implement me")
}

// Connects to custom Serum API
func NewWSClientWithEndpoint(addr string) *WSClient {
	connectionManager := connections.NewConnectionManager(addr)

	return &WSClient{
		addr:              addr,
		connectionManager: &connectionManager,
		requestID:         utils.NewRequestID(),
	}
}

// Set limit to 0 to get all bids/asks
func (w *WSClient) GetOrderbook(market string, limit uint32) (*pb.GetOrderbookResponse, error) {
	request, err := w.jsonRPCRequest("GetOrderbook", map[string]interface{}{"market": market, "limit": limit})
	if err != nil {
		return nil, err
	}

	return connections.WSResponse[pb.GetOrderbookResponse](w.connectionManager, request)
}

func (w *WSClient) GetOrderbookStream(ctx context.Context, market string, limit uint32, orderbookChan chan *pb.GetOrderbookStreamResponse) error {
	request, err := w.jsonRPCRequest("GetOrderbookStream", map[string]interface{}{"market": market, "limit": limit})
	if err != nil {
		return err
	}
	return connections.WSStream[pb.GetOrderbookStreamResponse](ctx, w.connectionManager, request, orderbookChan)
}

func (w *WSClient) GetMarkets() (*pb.GetMarketsResponse, error) {
	request, err := w.jsonRPCRequest("GetMarkets", nil)
	if err != nil {
		return nil, err
	}
	return connections.WSResponse[pb.GetMarketsResponse](w.connectionManager, request)
}

func (w *WSClient) jsonRPCRequest(method string, params map[string]interface{}) ([]byte, error) {
	id := w.requestID.Next()
	req := jsonrpc2.Request{
		Method: method,
		ID:     jsonrpc2.ID{Num: id},
	}
	if err := req.SetParams(params); err != nil {
		return nil, err
	}

	return req.MarshalJSON()
}
