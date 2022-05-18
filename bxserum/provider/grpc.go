package provider

import (
	"context"
	"github.com/bloXroute-Labs/serum-api/bxserum/connections"
	"github.com/bloXroute-Labs/serum-api/bxserum/transaction"
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

// NewGRPCClient connects to Mainnet Serum API
func NewGRPCClient() (*GRPCClient, error) {
	opts, err := DefaultRPCOpts(MainnetSerumAPIGRPC)
	if err != nil {
		return nil, err
	}
	return NewGRPCClientWithOpts(opts)
}

// NewGRPCTestnet connects to Testnet Serum API
func NewGRPCTestnet() (*GRPCClient, error) {
	panic("implement me")
}

// NewGRPCClientWithOpts connects to custom Serum API
func NewGRPCClientWithOpts(opts RPCOpts) (*GRPCClient, error) {
	conn, err := grpc.Dial(opts.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &GRPCClient{
		apiClient:  pb.NewApiClient(conn),
		privateKey: opts.PrivateKey,
	}, nil
}

// GetOrderbook returns the requested market's orderbook (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (g *GRPCClient) GetOrderbook(ctx context.Context, market string, limit uint32) (*pb.GetOrderbookResponse, error) {
	return g.apiClient.GetOrderbook(ctx, &pb.GetOrderBookRequest{Market: market, Limit: limit})
}

// GetOrderbookStream subscribes to a stream for changes to the requested market updates (e.g. asks and bids. Set limit to 0 for all bids/ asks).
func (g *GRPCClient) GetOrderbookStream(ctx context.Context, market string, limit uint32, outputChan chan *pb.GetOrderbookStreamResponse) error {
	stream, err := g.apiClient.GetOrderbookStream(ctx, &pb.GetOrderBookRequest{Market: market, Limit: limit})
	if err != nil {
		return err
	}

	return connections.GRPCStream[pb.GetOrderbookStreamResponse](stream, market, outputChan)
}

// GetTrades returns the requested market's currently executing trades. Set limit to 0 for all trades.
func (g *GRPCClient) GetTrades(ctx context.Context, market string, limit uint32) (*pb.GetTradesResponse, error) {
	return g.apiClient.GetTrades(ctx, &pb.GetTradesRequest{Market: market, Limit: limit})
}

// GetTradesStream subscribes to a stream for trades as they execute. Set limit to 0 for all trades.
func (g *GRPCClient) GetTradesStream(ctx context.Context, market string, limit uint32, outputChan chan *pb.GetTradesStreamResponse) error {
	stream, err := g.apiClient.GetTradeStream(ctx, &pb.GetTradesRequest{Market: market, Limit: limit})
	if err != nil {
		return err
	}

	return connections.GRPCStream[pb.GetTradesStreamResponse](stream, market, outputChan)
}

// GetTickers returns the requested market tickets. Set market to "" for all markets.
func (g *GRPCClient) GetTickers(ctx context.Context, market string) (*pb.GetTickersResponse, error) {
	return g.apiClient.GetTickers(ctx, &pb.GetTickersRequest{Market: market})
}

// GetOrders returns all opened orders by owner address and market
func (g *GRPCClient) GetOrders(ctx context.Context, market string, owner string) (*pb.GetOrdersResponse, error) {
	return g.apiClient.GetOrders(ctx, &pb.GetOrdersRequest{Market: market, Address: owner})
}

// GetMarkets returns the list of all available named markets
func (g *GRPCClient) GetMarkets(ctx context.Context) (*pb.GetMarketsResponse, error) {
	return g.apiClient.GetMarkets(ctx, &pb.GetMarketsRequest{})
}

// PostOrder returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (g *GRPCClient) PostOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	return g.apiClient.PostOrder(ctx, &pb.PostOrderRequest{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              types,
		Amount:            amount,
		Price:             price,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
	})
}

// PostSubmit posts the transaction string to the Solana network.
func (g *GRPCClient) PostSubmit(ctx context.Context, txBase64 string) (*pb.PostSubmitResponse, error) {
	return g.apiClient.PostSubmit(ctx, &pb.PostSubmitRequest{Transaction: txBase64})
}

// SubmitOrder builds a Serum market order, signs it, and submits to the network.
func (g *GRPCClient) SubmitOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := g.PostOrder(ctx, owner, payer, market, side, types, amount, price, opts)
	if err != nil {
		return "", err
	}

	txBase64, err := transaction.SignTxWithPrivateKey(order.Transaction, g.privateKey)
	if err != nil {
		return "", err
	}

	response, err := g.PostSubmit(ctx, txBase64)
	if err != nil {
		return "", err
	}
	return response.Signature, nil
}
