package provider

import (
	"context"
	"fmt"
	connections2 "github.com/bloXroute-Labs/solana-trader-client-go/connections"
	pb "github.com/bloXroute-Labs/solana-trader-client-go/proto"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	"github.com/gagliardetto/solana-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	pb.UnimplementedApiServer

	apiClient pb.ApiClient

	privateKey           *solana.PrivateKey
	recentBlockHashStore *recentBlockHashStore
}

// NewGRPCClient connects to Mainnet Trader API
func NewGRPCClient() (*GRPCClient, error) {
	opts := DefaultRPCOpts(MainnetGRPC)
	return NewGRPCClientWithOpts(opts)
}

// NewGRPCTestnet connects to Testnet Trader API
func NewGRPCTestnet() (*GRPCClient, error) {
	opts := DefaultRPCOpts(TestnetGRPC)
	return NewGRPCClientWithOpts(opts)
}

// NewGRPCDevnet connects to Devnet Trader API
func NewGRPCDevnet() (*GRPCClient, error) {
	opts := DefaultRPCOpts(DevnetGRPC)
	return NewGRPCClientWithOpts(opts)
}

// NewGRPCLocal connects to local Trader API
func NewGRPCLocal() (*GRPCClient, error) {
	opts := DefaultRPCOpts(LocalGRPC)
	return NewGRPCClientWithOpts(opts)
}

type blxrCredentials struct {
	authorization string
}

func (bc blxrCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": bc.authorization,
	}, nil
}

func (bc blxrCredentials) RequireTransportSecurity() bool {
	return false
}

// NewGRPCClientWithOpts connects to custom Trader API
func NewGRPCClientWithOpts(opts RPCOpts) (*GRPCClient, error) {
	authOption := grpc.WithPerRPCCredentials(blxrCredentials{authorization: opts.AuthHeader})
	conn, err := grpc.Dial(opts.Endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), authOption)
	if err != nil {
		return nil, err
	}
	client := &GRPCClient{
		apiClient:  pb.NewApiClient(conn),
		privateKey: opts.PrivateKey,
	}
	client.recentBlockHashStore = newRecentBlockHashStore(
		client.GetRecentBlockHash,
		client.GetRecentBlockHashStream,
		opts,
	)
	go client.recentBlockHashStore.run()
	return client, nil
}

func (g *GRPCClient) RecentBlockHash(ctx context.Context) (*pb.GetRecentBlockHashResponse, error) {
	return g.recentBlockHashStore.recentBlockHash(ctx)
}

func (g *GRPCClient) GetRecentBlockHash(ctx context.Context) (*pb.GetRecentBlockHashResponse, error) {
	return g.apiClient.GetRecentBlockHash(ctx, &pb.GetRecentBlockHashRequest{})
}

// GetOrderbook returns the requested market's orderbook (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (g *GRPCClient) GetOrderbook(ctx context.Context, market string, limit uint32) (*pb.GetOrderbookResponse, error) {
	return g.apiClient.GetOrderbook(ctx, &pb.GetOrderbookRequest{Market: market, Limit: limit})
}

// GetTrades returns the requested market's currently executing trades. Set limit to 0 for all trades.
func (g *GRPCClient) GetTrades(ctx context.Context, market string, limit uint32) (*pb.GetTradesResponse, error) {
	return g.apiClient.GetTrades(ctx, &pb.GetTradesRequest{Market: market, Limit: limit})
}

// GetTickers returns the requested market tickets. Set market to "" for all markets.
func (g *GRPCClient) GetTickers(ctx context.Context, market string) (*pb.GetTickersResponse, error) {
	return g.apiClient.GetTickers(ctx, &pb.GetTickersRequest{Market: market})
}

// GetOpenOrders returns all opened orders by owner address and market
func (g *GRPCClient) GetOpenOrders(ctx context.Context, market string, owner string, openOrdersAddress string) (*pb.GetOpenOrdersResponse, error) {
	return g.apiClient.GetOpenOrders(ctx, &pb.GetOpenOrdersRequest{Market: market, Address: owner, OpenOrdersAddress: openOrdersAddress})
}

// GetUnsettled returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (g *GRPCClient) GetUnsettled(ctx context.Context, market string, owner string) (*pb.GetUnsettledResponse, error) {
	return g.apiClient.GetUnsettled(ctx, &pb.GetUnsettledRequest{Market: market, Owner: owner})
}

// GetMarkets returns the list of all available named markets
func (g *GRPCClient) GetMarkets(ctx context.Context) (*pb.GetMarketsResponse, error) {
	return g.apiClient.GetMarkets(ctx, &pb.GetMarketsRequest{})
}

// GetAccountBalance returns all tokens associated with the owner address including Serum unsettled amounts
func (g *GRPCClient) GetAccountBalance(ctx context.Context, owner string) (*pb.GetAccountBalanceResponse, error) {
	return g.apiClient.GetAccountBalance(ctx, &pb.GetAccountBalanceRequest{OwnerAddress: owner})
}

// signAndSubmit signs the given transaction and submits it.
func (g *GRPCClient) signAndSubmit(ctx context.Context, tx string, skipPreFlight bool) (string, error) {
	if g.privateKey == nil {
		return "", ErrPrivateKeyNotFound
	}
	txBase64, err := transaction.SignTxWithPrivateKey(tx, *g.privateKey)
	if err != nil {
		return "", err
	}

	response, err := g.PostSubmit(ctx, txBase64, skipPreFlight)
	if err != nil {
		return "", err
	}

	return response.Signature, nil
}

// PostTradeSwap returns a partially signed transaction for submitting a swap request
func (g *GRPCClient) PostTradeSwap(ctx context.Context, owner, inToken, outToken string, inAmount, slippage float64, project pb.Project) (*pb.TradeSwapResponse, error) {
	return g.apiClient.PostTradeSwap(ctx, &pb.TradeSwapRequest{
		Owner:    owner,
		InToken:  inToken,
		OutToken: outToken,
		InAmount: inAmount,
		Slippage: slippage,
		Project:  project,
	})
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
func (g *GRPCClient) PostSubmit(ctx context.Context, txBase64 string, skipPreFlight bool) (*pb.PostSubmitResponse, error) {
	return g.apiClient.PostSubmit(ctx, &pb.PostSubmitRequest{Transaction: txBase64,
		SkipPreFlight: skipPreFlight})
}

// SubmitTradeSwap builds a TradeSwap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitTradeSwap(ctx context.Context, owner, inToken, outToken string, inAmount, slippage float64, project pb.Project, skipPreFlight bool) ([]string, error) {
	resp, err := g.apiClient.PostTradeSwap(ctx, &pb.TradeSwapRequest{
		Owner:    owner,
		InToken:  inToken,
		OutToken: outToken,
		InAmount: inAmount,
		Slippage: slippage,
		Project:  project,
	})
	if err != nil {
		return []string{}, err
	}

	var signatures []string
	for _, tx := range resp.Transactions {
		signature, err := g.signAndSubmit(ctx, tx, skipPreFlight)
		if err != nil {
			return signatures, err
		}

		signatures = append(signatures, signature)
	}

	return signatures, nil
}

// SubmitOrder builds a Serum market order, signs it, and submits to the network.
func (g *GRPCClient) SubmitOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := g.PostOrder(ctx, owner, payer, market, side, types, amount, price, opts)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
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
	skipPreFlight bool,
) (string, error) {
	order, err := g.PostCancelOrder(ctx, orderID, side, owner, market, openOrders)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, skipPreFlight)
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
	skipPreFlight bool,
) (string, error) {
	order, err := g.PostCancelByClientOrderID(ctx, clientOrderID, owner, market, openOrders)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, skipPreFlight)
}

func (g *GRPCClient) PostCancelAll(ctx context.Context, market, owner string, openOrders []string) (*pb.PostCancelAllResponse, error) {
	return g.apiClient.PostCancelAll(ctx, &pb.PostCancelAllRequest{
		Market:              market,
		OwnerAddress:        owner,
		OpenOrdersAddresses: openOrders,
	})
}

func (g *GRPCClient) SubmitCancelAll(ctx context.Context, market, owner string, openOrdersAddresses []string, skipPreFlight bool) ([]string, error) {
	orders, err := g.PostCancelAll(ctx, market, owner, openOrdersAddresses)
	if err != nil {
		return nil, err
	}

	var signatures []string
	for _, tx := range orders.Transactions {
		signature, err := g.signAndSubmit(ctx, tx, skipPreFlight)
		if err != nil {
			return signatures, err
		}

		signatures = append(signatures, signature)
	}

	return signatures, nil
}

// PostSettle returns a partially signed transaction for settling market funds. Typically, you want to use SubmitSettle instead of this.
func (g *GRPCClient) PostSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string) (*pb.PostSettleResponse, error) {
	return g.apiClient.PostSettle(ctx, &pb.PostSettleRequest{
		OwnerAddress:      owner,
		Market:            market,
		BaseTokenWallet:   baseTokenWallet,
		QuoteTokenWallet:  quoteTokenWallet,
		OpenOrdersAddress: openOrdersAccount,
	})
}

// SubmitSettle builds a market SubmitSettle transaction, signs it, and submits to the network.
func (g *GRPCClient) SubmitSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, skipPreflight bool) (string, error) {
	order, err := g.PostSettle(ctx, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, skipPreflight)
}

func (g *GRPCClient) PostReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	return g.apiClient.PostReplaceByClientOrderID(ctx, &pb.PostOrderRequest{
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

func (g *GRPCClient) SubmitReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := g.PostReplaceByClientOrderID(ctx, owner, payer, market, side, types, amount, price, opts)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

func (g *GRPCClient) PostReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	return g.apiClient.PostReplaceOrder(ctx, &pb.PostReplaceOrderRequest{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              types,
		Amount:            amount,
		Price:             price,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
		OrderID:           orderID,
	})
}

func (g *GRPCClient) SubmitReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side pb.Side, types []pb.OrderType, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := g.PostReplaceOrder(ctx, orderID, owner, payer, market, side, types, amount, price, opts)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

// GetOrderbookStream subscribes to a stream for changes to the requested market updates (e.g. asks and bids. Set limit to 0 for all bids/ asks).
func (g *GRPCClient) GetOrderbookStream(ctx context.Context, markets []string, limit uint32) (connections2.Streamer[*pb.GetOrderbooksStreamResponse], error) {
	stream, err := g.apiClient.GetOrderbooksStream(ctx, &pb.GetOrderbooksRequest{Markets: markets, Limit: limit})
	if err != nil {
		return nil, err
	}

	return connections2.GRPCStream[pb.GetOrderbooksStreamResponse](stream, fmt.Sprint(markets)), nil
}

// GetTradesStream subscribes to a stream for trades as they execute. Set limit to 0 for all trades.
func (g *GRPCClient) GetTradesStream(ctx context.Context, market string, limit uint32) (connections2.Streamer[*pb.GetTradesStreamResponse], error) {
	stream, err := g.apiClient.GetTradesStream(ctx, &pb.GetTradesRequest{Market: market, Limit: limit})
	if err != nil {
		return nil, err
	}

	return connections2.GRPCStream[pb.GetTradesStreamResponse](stream, market), nil
}

// GetOrderStatusStream subscribes to a stream that shows updates to the owner's orders
func (g *GRPCClient) GetOrderStatusStream(ctx context.Context, market, ownerAddress string) (connections2.Streamer[*pb.GetOrderStatusStreamResponse], error) {
	stream, err := g.apiClient.GetOrderStatusStream(ctx, &pb.GetOrderStatusStreamRequest{Market: market, OwnerAddress: ownerAddress})
	if err != nil {
		return nil, err
	}

	return connections2.GRPCStream[pb.GetOrderStatusStreamResponse](stream, market), nil
}

// GetRecentBlockHashStream subscribes to a stream for getting recent block hash.
func (g *GRPCClient) GetRecentBlockHashStream(ctx context.Context) (connections2.Streamer[*pb.GetRecentBlockHashResponse], error) {
	stream, err := g.apiClient.GetRecentBlockHashStream(ctx, &pb.GetRecentBlockHashRequest{})
	if err != nil {
		return nil, err
	}

	return connections2.GRPCStream[pb.GetRecentBlockHashResponse](stream, ""), nil
}