package provider

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"

	"github.com/bloXroute-Labs/solana-trader-client-go/connections"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/gagliardetto/solana-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	opts := DefaultRPCOpts(MainnetNYGRPC)
	opts.UseTLS = true
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
func NewGRPCClientWithOpts(opts RPCOpts, dialOpts ...grpc.DialOption) (*GRPCClient, error) {
	var (
		conn     grpc.ClientConnInterface
		err      error
		grpcOpts = make([]grpc.DialOption, 0)
	)

	transportOption := grpc.WithTransportCredentials(insecure.NewCredentials())
	if opts.UseTLS {
		transportOption = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))
	}
	grpcOpts = append(grpcOpts, transportOption)

	if !opts.DisableAuth {
		grpcOpts = append(grpcOpts, grpc.WithPerRPCCredentials(blxrCredentials{authorization: opts.AuthHeader}))
	}

	grpcOpts = append(grpcOpts, dialOpts...)
	conn, err = grpc.Dial(opts.Endpoint, grpcOpts...)
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
	if opts.CacheBlockHash {
		go client.recentBlockHashStore.run(context.Background())
	}
	return client, nil
}

func (g *GRPCClient) RecentBlockHash(ctx context.Context) (*pb.GetRecentBlockHashResponse, error) {
	return g.recentBlockHashStore.get(ctx)
}

func (g *GRPCClient) GetRecentBlockHash(ctx context.Context) (*pb.GetRecentBlockHashResponse, error) {
	return g.apiClient.GetRecentBlockHash(ctx, &pb.GetRecentBlockHashRequest{})
}

// GetTransaction returns details of a recent transaction
func (g *GRPCClient) GetTransaction(ctx context.Context, request *pb.GetTransactionRequest) (*pb.GetTransactionResponse, error) {
	return g.apiClient.GetTransaction(ctx, request)
}

// GetRaydiumPools returns pools on Raydium
func (g *GRPCClient) GetRaydiumPools(ctx context.Context, request *pb.GetRaydiumPoolsRequest) (*pb.GetRaydiumPoolsResponse, error) {
	return g.apiClient.GetRaydiumPools(ctx, request)
}

// GetRaydiumQuotes returns the possible amount(s) of outToken for an inToken and the route to achieve it on Raydium
func (g *GRPCClient) GetRaydiumQuotes(ctx context.Context, request *pb.GetRaydiumQuotesRequest) (*pb.GetRaydiumQuotesResponse, error) {
	return g.apiClient.GetRaydiumQuotes(ctx, request)
}

// GetRaydiumPrices returns the USDC price of requested tokens on Raydium
func (g *GRPCClient) GetRaydiumPrices(ctx context.Context, request *pb.GetRaydiumPricesRequest) (*pb.GetRaydiumPricesResponse, error) {
	return g.apiClient.GetRaydiumPrices(ctx, request)
}

// PostRaydiumSwap returns a partially signed transaction(s) for submitting a swap request on Raydium
func (g *GRPCClient) PostRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest) (*pb.PostRaydiumSwapResponse, error) {
	return g.apiClient.PostRaydiumSwap(ctx, request)
}

// PostRaydiumRouteSwap returns a partially signed transaction(s) for submitting a swap request on Raydium
func (g *GRPCClient) PostRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest) (*pb.PostRaydiumRouteSwapResponse, error) {
	return g.apiClient.PostRaydiumRouteSwap(ctx, request)
}

// GetJupiterQuotes returns the possible amount(s) of outToken for an inToken and the route to achieve it on Jupiter
func (g *GRPCClient) GetJupiterQuotes(ctx context.Context, request *pb.GetJupiterQuotesRequest) (*pb.GetJupiterQuotesResponse, error) {
	return g.apiClient.GetJupiterQuotes(ctx, request)
}

// GetJupiterPrices returns the USDC price of requested tokens on Jupiter
func (g *GRPCClient) GetJupiterPrices(ctx context.Context, request *pb.GetJupiterPricesRequest) (*pb.GetJupiterPricesResponse, error) {
	return g.apiClient.GetJupiterPrices(ctx, request)
}

// PostJupiterSwap returns a partially signed transaction(s) for submitting a swap request on Jupiter
func (g *GRPCClient) PostJupiterSwap(ctx context.Context, request *pb.PostJupiterSwapRequest) (*pb.PostJupiterSwapResponse, error) {
	return g.apiClient.PostJupiterSwap(ctx, request)
}

// PostJupiterRouteSwap returns a partially signed transaction(s) for submitting a swap request on Jupiter
func (g *GRPCClient) PostJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest) (*pb.PostJupiterRouteSwapResponse, error) {
	return g.apiClient.PostJupiterRouteSwap(ctx, request)
}

// GetTrades returns the requested market's currently executing trades. Set limit to 0 for all trades.
func (g *GRPCClient) GetTrades(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetTradesResponse, error) {
	return g.apiClient.GetTrades(ctx, &pb.GetTradesRequest{Market: market, Limit: limit, Project: project})
}

// GetOrderByID returns an order by id
func (g *GRPCClient) GetOrderByID(ctx context.Context, in *pb.GetOrderByIDRequest) (*pb.GetOrderByIDResponse, error) {
	return g.apiClient.GetOrderByID(ctx, in)
}

// GetAccountBalance returns all tokens associated with the owner address including Serum unsettled amounts
func (g *GRPCClient) GetAccountBalance(ctx context.Context, owner string) (*pb.GetAccountBalanceResponse, error) {
	return g.apiClient.GetAccountBalance(ctx, &pb.GetAccountBalanceRequest{OwnerAddress: owner})
}

// PostSubmitV2 posts the transaction string to the Solana network.
func (g *GRPCClient) PostSubmitV2(ctx context.Context, tx *pb.TransactionMessage, skipPreFlight bool) (*pb.PostSubmitResponse, error) {
	return g.apiClient.PostSubmitV2(ctx, &pb.PostSubmitRequest{Transaction: tx,
		SkipPreFlight: skipPreFlight})
}

// PostSubmitBatchV2 posts a bundle of transactions string based on a specific SubmitStrategy to the Solana network.
func (g *GRPCClient) PostSubmitBatchV2(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error) {
	return g.apiClient.PostSubmitBatchV2(ctx, request)
}

// SubmitRaydiumSwap builds a Raydium Swap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostRaydiumSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitRaydiumRouteSwap builds a Raydium RouteSwap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostRaydiumRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitJupiterSwap builds a Jupiter Swap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitJupiterSwap(ctx context.Context, request *pb.PostJupiterSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostJupiterSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitJupiterRouteSwap builds a Jupiter RouteSwap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostJupiterRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// signAndSubmit signs the given transaction and submits it.
func (g *GRPCClient) signAndSubmit(ctx context.Context, tx *pb.TransactionMessage, skipPreFlight bool) (string, error) {
	if g.privateKey == nil {
		return "", ErrPrivateKeyNotFound
	}
	txBase64, err := transaction.SignTxWithPrivateKey(tx.Content, *g.privateKey)
	if err != nil {
		return "", err
	}

	response, err := g.PostSubmit(ctx, &pb.TransactionMessage{
		Content:   txBase64,
		IsCleanup: tx.IsCleanup,
	}, skipPreFlight)
	if err != nil {
		return "", err
	}

	return response.Signature, nil
}

// PostSubmit posts the transaction string to the Solana network.
func (g *GRPCClient) PostSubmit(ctx context.Context, tx *pb.TransactionMessage, skipPreFlight bool) (*pb.PostSubmitResponse, error) {
	return g.apiClient.PostSubmitV2(ctx, &pb.PostSubmitRequest{Transaction: tx,
		SkipPreFlight: skipPreFlight})
}

// PostSubmitBatch posts a bundle of transactions string based on a specific SubmitStrategy to the Solana network.
func (g *GRPCClient) PostSubmitBatch(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error) {
	return g.apiClient.PostSubmitBatchV2(ctx, request)
}

// signAndSubmitBatch signs the given transactions and submits them.
func (g *GRPCClient) signAndSubmitBatch(ctx context.Context, transactions []*pb.TransactionMessage, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	if g.privateKey == nil {
		return nil, ErrPrivateKeyNotFound
	}
	batchRequest, err := buildBatchRequest(transactions, *g.privateKey, opts)
	if err != nil {
		return nil, err
	}

	return g.PostSubmitBatch(ctx, batchRequest)
}

// GetOrderbookStream subscribes to a stream for changes to the requested market updates (e.g. asks and bids. Set limit to 0 for all bids/ asks).
func (g *GRPCClient) GetOrderbookStream(ctx context.Context, markets []string, limit uint32, project pb.Project) (connections.Streamer[*pb.GetOrderbooksStreamResponse], error) {
	stream, err := g.apiClient.GetOrderbooksStream(ctx, &pb.GetOrderbooksRequest{
		Markets: markets, Limit: limit,
		Project: project})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetOrderbooksStreamResponse](stream, fmt.Sprint(markets)), nil
}

// GetMarketDepthsStream subscribes to a stream for changes to the requested market data updates (e.g. asks and bids. Set limit to 0 for all bids/ asks).
func (g *GRPCClient) GetMarketDepthsStream(ctx context.Context, markets []string, limit uint32, project pb.Project) (connections.Streamer[*pb.GetMarketDepthsStreamResponse], error) {
	stream, err := g.apiClient.GetMarketDepthsStream(ctx, &pb.GetMarketDepthsRequest{Markets: markets, Limit: limit, Project: project})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetMarketDepthsStreamResponse](stream, fmt.Sprint(markets)), nil
}

// GetTradesStream subscribes to a stream for trades as they execute. Set limit to 0 for all trades.
func (g *GRPCClient) GetTradesStream(ctx context.Context, market string, limit uint32, project pb.Project) (connections.Streamer[*pb.GetTradesStreamResponse], error) {
	stream, err := g.apiClient.GetTradesStream(ctx, &pb.GetTradesRequest{Market: market, Limit: limit, Project: project})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetTradesStreamResponse](stream, market), nil
}

// GetOrderStatusStream subscribes to a stream that shows updates to the owner's orders
func (g *GRPCClient) GetOrderStatusStream(ctx context.Context, market, ownerAddress string, project pb.Project) (connections.Streamer[*pb.GetOrderStatusStreamResponse], error) {
	stream, err := g.apiClient.GetOrderStatusStream(ctx, &pb.GetOrderStatusStreamRequest{Market: market, OwnerAddress: ownerAddress, Project: project})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetOrderStatusStreamResponse](stream, market), nil
}

// GetRecentBlockHashStream subscribes to a stream for getting recent block hash.
func (g *GRPCClient) GetRecentBlockHashStream(ctx context.Context) (connections.Streamer[*pb.GetRecentBlockHashResponse], error) {
	stream, err := g.apiClient.GetRecentBlockHashStream(ctx, &pb.GetRecentBlockHashRequest{})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetRecentBlockHashResponse](stream, ""), nil
}

// GetQuotesStream subscribes to a stream for getting recent quotes of tokens of interest.
func (g *GRPCClient) GetQuotesStream(ctx context.Context, projects []pb.Project, tokenPairs []*pb.TokenPair) (connections.Streamer[*pb.GetQuotesStreamResponse], error) {
	stream, err := g.apiClient.GetQuotesStream(ctx, &pb.GetQuotesStreamRequest{
		Projects:   projects,
		TokenPairs: tokenPairs,
	})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetQuotesStreamResponse](stream, ""), nil
}

// GetPoolReservesStream subscribes to a stream for getting recent quotes of tokens of interest.
func (g *GRPCClient) GetPoolReservesStream(ctx context.Context, projects []pb.Project) (connections.Streamer[*pb.GetPoolReservesStreamResponse], error) {
	stream, err := g.apiClient.GetPoolReservesStream(ctx, &pb.GetPoolReservesStreamRequest{
		Projects: projects,
	})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetPoolReservesStreamResponse](stream, ""), nil
}

// GetPricesStream subscribes to a stream for getting recent prices of tokens of interest.
func (g *GRPCClient) GetPricesStream(ctx context.Context, projects []pb.Project, tokens []string) (connections.Streamer[*pb.GetPricesStreamResponse], error) {
	stream, err := g.apiClient.GetPricesStream(ctx, &pb.GetPricesStreamRequest{
		Projects: projects,
		Tokens:   tokens,
	})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetPricesStreamResponse](stream, ""), nil
}

// GetSwapsStream subscribes to a stream for getting recent swaps on projects & markets of interest.
func (g *GRPCClient) GetSwapsStream(
	ctx context.Context,
	projects []pb.Project,
	markets []string,
	includeFailed bool,
) (connections.Streamer[*pb.GetSwapsStreamResponse], error) {
	stream, err := g.apiClient.GetSwapsStream(ctx, &pb.GetSwapsStreamRequest{
		Projects:      projects,
		Pools:         markets,
		IncludeFailed: includeFailed,
	})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetSwapsStreamResponse](stream, ""), nil
}

// GetNewRaydiumPoolsStream subscribes to a stream for getting recent swaps on projects & markets of interest.
func (g *GRPCClient) GetNewRaydiumPoolsStream(
	ctx context.Context,
) (connections.Streamer[*pb.GetNewRaydiumPoolsResponse], error) {
	stream, err := g.apiClient.GetNewRaydiumPoolsStream(ctx, &pb.GetNewRaydiumPoolsRequest{})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetNewRaydiumPoolsResponse](stream, ""), nil
}

// GetBlockStream subscribes to a stream for getting recent blocks.
func (g *GRPCClient) GetBlockStream(ctx context.Context) (connections.Streamer[*pb.GetBlockStreamResponse], error) {
	stream, err := g.apiClient.GetBlockStream(ctx, &pb.GetBlockStreamRequest{})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetBlockStreamResponse](stream, ""), nil
}

// V2 Openbook

// GetOrderbookV2 returns the requested market's orderbook (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (g *GRPCClient) GetOrderbookV2(ctx context.Context, market string, limit uint32) (*pb.GetOrderbookResponseV2, error) {
	return g.apiClient.GetOrderbookV2(ctx, &pb.GetOrderbookRequestV2{Market: market, Limit: limit})
}

// GetMarketDepthV2 returns the requested market's coalesced price data (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (g *GRPCClient) GetMarketDepthV2(ctx context.Context, market string, limit uint32) (*pb.GetMarketDepthResponseV2, error) {
	return g.apiClient.GetMarketDepthV2(ctx, &pb.GetMarketDepthRequestV2{Market: market, Limit: limit})
}

// GetTickersV2 returns the requested market tickets. Set market to "" for all markets.
func (g *GRPCClient) GetTickersV2(ctx context.Context, market string) (*pb.GetTickersResponseV2, error) {
	return g.apiClient.GetTickersV2(ctx, &pb.GetTickersRequestV2{Market: market})
}

// GetOpenOrdersV2 returns all open orders by owner address and market
func (g *GRPCClient) GetOpenOrdersV2(ctx context.Context, market string, owner string, openOrdersAddress string, orderID string, clientOrderID uint64) (*pb.GetOpenOrdersResponseV2, error) {
	return g.apiClient.GetOpenOrdersV2(ctx, &pb.GetOpenOrdersRequestV2{Market: market, Address: owner, OpenOrdersAddress: openOrdersAddress, OrderID: orderID, ClientOrderID: clientOrderID})
}

// GetUnsettledV2 returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (g *GRPCClient) GetUnsettledV2(ctx context.Context, market string, ownerAddress string) (*pb.GetUnsettledResponse, error) {
	return g.apiClient.GetUnsettledV2(ctx, &pb.GetUnsettledRequestV2{Market: market, OwnerAddress: ownerAddress})
}

// GetMarketsV2 returns the list of all available named markets
func (g *GRPCClient) GetMarketsV2(ctx context.Context) (*pb.GetMarketsResponseV2, error) {
	return g.apiClient.GetMarketsV2(ctx, &pb.GetMarketsRequestV2{})
}

// PostOrderV2 returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (g *GRPCClient) PostOrderV2(ctx context.Context, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	return g.apiClient.PostOrderV2(ctx, &pb.PostOrderRequestV2{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              orderType,
		Amount:            amount,
		Price:             price,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
	})
}

// PostOrderV2WithPriorityFee returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (g *GRPCClient) PostOrderV2WithPriorityFee(ctx context.Context, owner, payer, market string, side string,
	orderType string, amount, price float64, computeLimit uint32, computePrice uint64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	return g.apiClient.PostOrderV2(ctx, &pb.PostOrderRequestV2{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              orderType,
		Amount:            amount,
		Price:             price,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ComputeLimit:      computeLimit,
		ComputePrice:      computePrice,
		ClientOrderID:     opts.ClientOrderID,
	})
}

// SubmitOrderV2 builds a Serum market order, signs it, and submits to the network.
func (g *GRPCClient) SubmitOrderV2(ctx context.Context, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := g.PostOrderV2(ctx, owner, payer, market, side, orderType, amount, price, opts)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

// SubmitOrderV2WithPriorityFee builds a Serum market order, signs it, and submits to the network with specified computeLimit and computePrice
func (g *GRPCClient) SubmitOrderV2WithPriorityFee(ctx context.Context, owner, payer, market string, side string,
	orderType string, amount, price float64, computeLimit uint32, computePrice uint64, opts PostOrderOpts) (string, error) {
	order, err := g.PostOrderV2WithPriorityFee(ctx, owner, payer, market, side, orderType, amount, price, computeLimit, computePrice, opts)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}

// PostCancelOrderV2 builds a Serum cancel order.
func (g *GRPCClient) PostCancelOrderV2(
	ctx context.Context,
	orderID string,
	clientOrderID uint64,
	side string,
	owner,
	market,
	openOrders string,
) (*pb.PostCancelOrderResponseV2, error) {
	return g.apiClient.PostCancelOrderV2(ctx, &pb.PostCancelOrderRequestV2{
		OrderID:           orderID,
		Side:              side,
		OwnerAddress:      owner,
		MarketAddress:     market,
		OpenOrdersAddress: openOrders,
		ClientOrderID:     clientOrderID,
	})
}

// SubmitCancelOrderV2 builds a Serum cancel order, signs and submits it to the network.
func (g *GRPCClient) SubmitCancelOrderV2(
	ctx context.Context,
	orderID string,
	clientOrderID uint64,
	side string,
	owner,
	market,
	openOrders string,
	opts SubmitOpts,
) (*pb.PostSubmitBatchResponse, error) {
	order, err := g.PostCancelOrderV2(ctx, orderID, clientOrderID, side, owner, market, openOrders)
	if err != nil {
		return nil, err
	}

	return g.signAndSubmitBatch(ctx, order.Transactions, opts)
}

// PostSettleV2 returns a partially signed transaction for settling market funds. Typically, you want to use SubmitSettle instead of this.
func (g *GRPCClient) PostSettleV2(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string) (*pb.PostSettleResponse, error) {
	return g.apiClient.PostSettleV2(ctx, &pb.PostSettleRequestV2{
		OwnerAddress:      owner,
		Market:            market,
		BaseTokenWallet:   baseTokenWallet,
		QuoteTokenWallet:  quoteTokenWallet,
		OpenOrdersAddress: openOrdersAccount,
	})
}

// SubmitSettleV2 builds a market SubmitSettle transaction, signs it, and submits to the network.
func (g *GRPCClient) SubmitSettleV2(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, skipPreflight bool) (string, error) {
	order, err := g.PostSettleV2(ctx, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, skipPreflight)
}

func (g *GRPCClient) PostReplaceOrderV2(ctx context.Context, orderID, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	return g.apiClient.PostReplaceOrderV2(ctx, &pb.PostReplaceOrderRequestV2{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              orderType,
		Amount:            amount,
		Price:             price,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
		OrderID:           orderID,
	})
}

func (g *GRPCClient) SubmitReplaceOrderV2(ctx context.Context, orderID, owner, payer, market string, side string, orderType string, amount, price float64, opts PostOrderOpts) (string, error) {
	order, err := g.PostReplaceOrderV2(ctx, orderID, owner, payer, market, side, orderType, amount, price, opts)
	if err != nil {
		return "", err
	}

	return g.signAndSubmit(ctx, order.Transaction, opts.SkipPreFlight)
}
