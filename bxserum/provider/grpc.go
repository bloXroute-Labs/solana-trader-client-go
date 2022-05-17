package provider

import (
	"context"
	"github.com/bloXroute-Labs/serum-api/bxserum/connections"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/gagliardetto/solana-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	pb.UnimplementedApiServer

	apiClient  pb.ApiClient
	privateKey solana.PrivateKey
}

// NewGRPCClient onnects to Mainnet Serum API
func NewGRPCClient() (*GRPCClient, error) {
	return NewGRPCClientWithEndpoint("174.129.154.164:1811")
}

// NewGRPCTestnet onnects to Testnet Serum API
func NewGRPCTestnet() (*GRPCClient, error) {
	panic("implement me")
}

// NewGRPCClientWithEndpoint connects to custom Serum API
func NewGRPCClientWithEndpoint(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &GRPCClient{
		apiClient: pb.NewApiClient(conn),
	}, nil
}

// Set limit to 0 to get all bids/asks
func (g *GRPCClient) GetOrderbook(ctx context.Context, market string, limit uint32) (*pb.GetOrderbookResponse, error) {
	return g.apiClient.GetOrderbook(ctx, &pb.GetOrderBookRequest{Market: market, Limit: limit})
}

func (g *GRPCClient) GetOrderbookStream(ctx context.Context, market string, limit uint32, outputChan chan *pb.GetOrderbookStreamResponse) error {
	stream, err := g.apiClient.GetOrderbookStream(ctx, &pb.GetOrderBookRequest{Market: market, Limit: limit})
	if err != nil {
		return err
	}

	return connections.GRPCStream[pb.GetOrderbookStreamResponse](stream, market, outputChan)
}

// Set limit to 0 to get all trades
func (g *GRPCClient) GetTrades(ctx context.Context, market string, limit uint32) (*pb.GetTradesResponse, error) {
	return g.apiClient.GetTrades(ctx, &pb.GetTradesRequest{Market: market, Limit: limit})
}

func (g *GRPCClient) GetTradesStream(ctx context.Context, market string, limit uint32, outputChan chan *pb.GetTradesStreamResponse) error {
	stream, err := g.apiClient.GetTradeStream(ctx, &pb.GetTradesRequest{Market: market, Limit: limit})
	if err != nil {
		return err
	}

	return connections.GRPCStream[pb.GetTradesStreamResponse](stream, market, outputChan)
}

// GetOrders returns all opened orders by owner address and market
func (g *GRPCClient) GetOrders(ctx context.Context, market string, owner string) (*pb.GetOrdersResponse, error) {
	return g.apiClient.GetOrders(ctx, &pb.GetOrdersRequest{Market: market, Address: owner})
}

// Set market to empty string to 0 to get all tickers
func (g *GRPCClient) GetTickers(ctx context.Context, market string) (*pb.GetTickersResponse, error) {
	return g.apiClient.GetTickers(ctx, &pb.GetTickersRequest{Market: market})
}

func (g *GRPCClient) GetMarkets(ctx context.Context) (*pb.GetMarketsResponse, error) {
	return g.apiClient.GetMarkets(ctx, &pb.GetMarketsRequest{})
}
