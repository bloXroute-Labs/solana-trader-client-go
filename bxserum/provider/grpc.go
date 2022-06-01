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
func (g *GRPCClient) GetOrderbookStream(ctx context.Context, market string, limit uint32, outputChan chan *pb.GetOrderbooksStreamResponse) error {
	stream, err := g.apiClient.GetOrderbooksStream(ctx, &pb.GetOrderBookRequest{Market: market, Limit: limit})
	if err != nil {
		return err
	}

	return connections.GRPCStream[pb.GetOrderbooksStreamResponse](stream, market, outputChan)
}

// GetTrades returns the requested market's currently executing trades. Set limit to 0 for all trades.
func (g *GRPCClient) GetTrades(ctx context.Context, market string, limit uint32) (*pb.GetTradesResponse, error) {
	return g.apiClient.GetTrades(ctx, &pb.GetTradesRequest{Market: market, Limit: limit})
}

// GetTradesStream subscribes to a stream for trades as they execute. Set limit to 0 for all trades.
func (g *GRPCClient) GetTradesStream(ctx context.Context, market string, limit uint32, outputChan chan *pb.GetTradesStreamResponse) error {
	stream, err := g.apiClient.GetTradesStream(ctx, &pb.GetTradesRequest{Market: market, Limit: limit})
	if err != nil {
		return err
	}

	return connections.GRPCStream[pb.GetTradesStreamResponse](stream, market, outputChan)
}

// GetTickers returns the requested market tickets. Set market to "" for all markets.
func (g *GRPCClient) GetTickers(ctx context.Context, market string) (*pb.GetTickersResponse, error) {
	return g.apiClient.GetTickers(ctx, &pb.GetTickersRequest{Market: market})
}

// GetOpenOrders returns all opened orders by owner address and market
func (g *GRPCClient) GetOpenOrders(ctx context.Context, market string, owner string) (*pb.GetOpenOrdersResponse, error) {
	return g.apiClient.GetOpenOrders(ctx, &pb.GetOpenOrdersRequest{Market: market, Address: owner})
}

// GetUnsettled returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (g *GRPCClient) GetUnsettled(ctx context.Context, market string, owner string) (*pb.GetUnsettledResponse, error) {
	return g.apiClient.GetUnsettled(ctx, &pb.GetUnsettledRequest{Market: market, Owner: owner})
}

// GetMarkets returns the list of all available named markets
func (g *GRPCClient) GetMarkets(ctx context.Context) (*pb.GetMarketsResponse, error) {
	return g.apiClient.GetMarkets(ctx, &pb.GetMarketsRequest{})
}

// signAndSubmit signs the given transaction and submits it.
func (g *GRPCClient) signAndSubmit(ctx context.Context, tx string) (string, error) {
	txBase64, err := transaction.SignTxWithPrivateKey(tx, g.privateKey)
	if err != nil {
		return "", err
	}

	response, err := g.PostSubmit(ctx, txBase64)
	if err != nil {
		return "", err
	}

	return response.Signature, nil
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
func (g *GRPCClient) SubmitOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (string, string, error) {
	order, err := g.PostOrder(ctx, owner, payer, market, side, types, amount, price, opts)
	if err != nil {
		return "", "", err
	}

	sig, err := g.signAndSubmit(ctx, order.Transaction)
	return sig, order.OpenOrdersAddress, err
}

// PostCancelOrder builds a Serum cancel order.
func (g *GRPCClient) PostCancelOrder(
	ctx context.Context,
	orderID string,
	side pb.Side,
	owner,
	market,
	openOrders string,
) (*pb.PostCancelOrderResponse, error) {
	return g.apiClient.PostCancelOrder(ctx, &pb.PostCancelOrderRequest{
		OrderID:           orderID,
		Side:              side,
		OwnerAddress:      owner,
		MarketAddress:     market,
		OpenOrdersAddress: openOrders,
	})
}

// SubmitCancelOrder builds a Serum cancel order, signs and submits it to the network.
func (g *GRPCClient) SubmitCancelOrder(
	ctx context.Context,
	orderID string,
	side pb.Side,
	owner,
	market,
	openOrders string,
) (string, error) {
	order, err := g.PostCancelOrder(ctx, orderID, side, owner, market, openOrders)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction)
}

// PostCancelByClientOrderID builds a Serum cancel order by client ID.
func (g *GRPCClient) PostCancelByClientOrderID(
	ctx context.Context,
	clientOrderID uint64,
	owner,
	market,
	openOrders string,
) (*pb.PostCancelOrderResponse, error) {
	return g.apiClient.PostCancelByClientOrderID(ctx, &pb.PostCancelByClientOrderIDRequest{
		ClientOrderID:     clientOrderID,
		OwnerAddress:      owner,
		MarketAddress:     market,
		OpenOrdersAddress: openOrders,
	})
}

// SubmitCancelByClientOrderID builds a Serum cancel order by client ID, signs and submits it to the network.
func (g *GRPCClient) SubmitCancelByClientOrderID(
	ctx context.Context,
	clientOrderID uint64,
	owner,
	market,
	openOrders string,
) (string, error) {
	order, err := g.PostCancelByClientOrderID(ctx, clientOrderID, owner, market, openOrders)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction)
}

// PostSettle returns a partially signed transaction for settling market funds. Typically, you want to use SettleFunds instead of this.
func (g *GRPCClient) PostSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string) (*pb.PostSettleResponse, error) {
	return g.apiClient.PostSettle(ctx, &pb.PostSettleRequest{
		OwnerAddress:      owner,
		Market:            market,
		BaseTokenWallet:   baseTokenWallet,
		QuoteTokenWallet:  quoteTokenWallet,
		OpenOrdersAddress: openOrdersAccount,
	})
}

// SettleFunds builds a market SettleFunds transaction, signs it, and submits to the network.
func (g *GRPCClient) SettleFunds(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string) (string, error) {
	order, err := g.PostSettle(ctx, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount)
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
