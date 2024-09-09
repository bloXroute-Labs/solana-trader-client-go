package provider

import (
	"context"
	"crypto/tls"
	"fmt"

	package_info "github.com/bloXroute-Labs/solana-trader-client-go"
	"github.com/bloXroute-Labs/solana-trader-client-go/connections"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"github.com/bloXroute-Labs/solana-trader-proto/common"
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
	opts.UseTLS = true
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
		"x-sdk":         package_info.Name,
		"x-sdk-version": package_info.Version,
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
	grpcOpts = append(grpcOpts, grpc.WithDefaultCallOptions(&grpc.MaxRecvMsgSizeCallOption{MaxRecvMsgSize: 1024 * 1024 * 16}))
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

// GetRaydiumCLMMQuotes returns the CLMM quotes on Raydium
func (g *GRPCClient) GetRaydiumCLMMQuotes(ctx context.Context, request *pb.GetRaydium) (*pb.GetRaydiumCLMMQuotesResponse, error) {
	return g.apiClient.GetRaydiumCLMMQuotes(ctx, request)
}

// GetRaydiumCLMMPools returns the CLMM pools on Raydium
func (g *GRPCClient) GetRaydiumCLMMPools(ctx context.Context, request *pb.GetRaydiumCLMMPoolsRequest) (*pb.GetRaydiumCLMMPoolsResponse, error) {
	return g.apiClient.GetRaydiumCLMMPools(ctx, request)
}

// PostRaydiumCLMMSwap returns a partially signed transaction(s) for submitting a swap request on Raydium
func (g *GRPCClient) PostRaydiumCLMMSwap(ctx context.Context, request *pb.PostRaydiumCLMMSwapRequest) (*pb.PostRaydiumCLMMSwapResponse, error) {
	return g.apiClient.PostRaydiumCLMMSwap(ctx, request)
}

// PostRaydiumCLMMRouteSwap returns a partially signed transaction(s) for submitting a route swap request on Raydium
func (g *GRPCClient) PostRaydiumCLMMRouteSwap(ctx context.Context, request *pb.PostRaydiumCLMMRouteSwapRequest) (*pb.PostRaydiumCLMMRouteSwapResponse, error) {
	return g.apiClient.PostRaydiumCLMMRouteSwap(ctx, request)
}

// GetPriorityFee returns a priority fee estimate for a given percentile
func (g *GRPCClient) GetPriorityFee(ctx context.Context, request *pb.GetPriorityFeeRequest) (*pb.GetPriorityFeeResponse, error) {
	return g.apiClient.GetPriorityFee(ctx, request)
}

// GetRateLimit returns details of an account rate-limits
func (g *GRPCClient) GetRateLimit(ctx context.Context, request *pb.GetRateLimitRequest) (*pb.GetRateLimitResponse, error) {
	return g.apiClient.GetRateLimit(ctx, request)
}

// GetTransaction returns details of a recent transaction
func (g *GRPCClient) GetTransaction(ctx context.Context, request *pb.GetTransactionRequest) (*pb.GetTransactionResponse, error) {
	return g.apiClient.GetTransaction(ctx, request)
}

// GetRaydiumPoolReserve returns pools details for a given set of pairs or addresses on Raydium
func (g *GRPCClient) GetRaydiumPoolReserve(ctx context.Context, req *pb.GetRaydiumPoolReserveRequest) (*pb.GetRaydiumPoolReserveResponse, error) {
	return g.apiClient.GetRaydiumPoolReserve(ctx, req)
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

// SubmitRaydiumCLMMSwap builds a Raydium Swap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitRaydiumCLMMSwap(ctx context.Context, request *pb.PostRaydiumCLMMSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostRaydiumCLMMSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, opts)
}

// SubmitRaydiumCLMMRouteSwap builds a Raydium RouteSwap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitRaydiumCLMMRouteSwap(ctx context.Context, request *pb.PostRaydiumCLMMRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostRaydiumCLMMRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, opts)
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

// PostJupiterSwapInstructions returns instructions to build a transaction and submit it on jupiter
func (g *GRPCClient) PostJupiterSwapInstructions(ctx context.Context, request *pb.PostJupiterSwapInstructionsRequest) (*pb.PostJupiterSwapInstructionsResponse, error) {
	return g.apiClient.PostJupiterSwapInstructions(ctx, request)
}

// PostRaydiumSwapInstructions returns instructions to build a transaction and submit it on raydium
func (g *GRPCClient) PostRaydiumSwapInstructions(ctx context.Context, request *pb.PostRaydiumSwapInstructionsRequest) (*pb.PostRaydiumSwapInstructionsResponse, error) {
	return g.apiClient.PostRaydiumSwapInstructions(ctx, request)
}

// PostJupiterRouteSwap returns a partially signed transaction(s) for submitting a swap request on Jupiter
func (g *GRPCClient) PostJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest) (*pb.PostJupiterRouteSwapResponse, error) {
	return g.apiClient.PostJupiterRouteSwap(ctx, request)
}

// GetOrderbook returns the requested market's orderbook (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (g *GRPCClient) GetOrderbook(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetOrderbookResponse, error) {
	return g.apiClient.GetOrderbook(ctx, &pb.GetOrderbookRequest{Market: market, Limit: limit, Project: project})
}

// GetMarketDepth returns the requested market's coalesced price data (e.g. asks and bids). Set limit to 0 for all bids / asks.
func (g *GRPCClient) GetMarketDepth(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetMarketDepthResponse, error) {
	return g.apiClient.GetMarketDepth(ctx, &pb.GetMarketDepthRequest{Market: market, Limit: limit, Project: project})
}

// GetPools returns pools for given projects.
func (g *GRPCClient) GetPools(ctx context.Context, projects []pb.Project) (*pb.GetPoolsResponse, error) {
	return g.apiClient.GetPools(ctx, &pb.GetPoolsRequest{Projects: projects})
}

// GetTrades returns the requested market's currently executing trades. Set limit to 0 for all trades.
func (g *GRPCClient) GetTrades(ctx context.Context, market string, limit uint32, project pb.Project) (*pb.GetTradesResponse, error) {
	return g.apiClient.GetTrades(ctx, &pb.GetTradesRequest{Market: market, Limit: limit, Project: project})
}

// GetTickers returns the requested market tickets. Set market to "" for all markets.
func (g *GRPCClient) GetTickers(ctx context.Context, market string, project pb.Project) (*pb.GetTickersResponse, error) {
	return g.apiClient.GetTickers(ctx, &pb.GetTickersRequest{Market: market, Project: project})
}

// GetOpenOrders returns all open orders by owner address and market
func (g *GRPCClient) GetOpenOrders(ctx context.Context, market string, owner string, openOrdersAddress string, project pb.Project) (*pb.GetOpenOrdersResponse, error) {
	return g.apiClient.GetOpenOrders(ctx, &pb.GetOpenOrdersRequest{Market: market, Address: owner, OpenOrdersAddress: openOrdersAddress, Project: project})
}

// GetOrderByID returns an order by id
func (g *GRPCClient) GetOrderByID(ctx context.Context, in *pb.GetOrderByIDRequest) (*pb.GetOrderByIDResponse, error) {
	return g.apiClient.GetOrderByID(ctx, in)
}

// GetUnsettled returns all OpenOrders accounts for a given market with the amounts of unsettled funds
func (g *GRPCClient) GetUnsettled(ctx context.Context, market string, ownerAddress string, project pb.Project) (*pb.GetUnsettledResponse, error) {
	return g.apiClient.GetUnsettled(ctx, &pb.GetUnsettledRequest{Market: market, OwnerAddress: ownerAddress, Project: project})
}

// GetMarkets returns the list of all available named markets
func (g *GRPCClient) GetMarkets(ctx context.Context) (*pb.GetMarketsResponse, error) {
	return g.apiClient.GetMarkets(ctx, &pb.GetMarketsRequest{})
}

// GetAccountBalance returns all tokens associated with the owner address including Serum unsettled amounts
func (g *GRPCClient) GetAccountBalance(ctx context.Context, owner string) (*pb.GetAccountBalanceResponse, error) {
	return g.apiClient.GetAccountBalanceV2(ctx, &pb.GetAccountBalanceRequest{OwnerAddress: owner})
}

// GetTokenAccounts returns all tokens associated with the owner address
func (g *GRPCClient) GetTokenAccounts(ctx context.Context, req *pb.GetTokenAccountsRequest) (*pb.GetTokenAccountsResponse, error) {
	return g.apiClient.GetTokenAccounts(ctx, req)
}

// GetPrice returns the USDC price of requested tokens
func (g *GRPCClient) GetPrice(ctx context.Context, tokens []string) (*pb.GetPriceResponse, error) {
	return g.apiClient.GetPrice(ctx, &pb.GetPriceRequest{Tokens: tokens})
}

// GetQuotes returns the possible amount(s) of outToken for an inToken and the route to achieve it
func (g *GRPCClient) GetQuotes(ctx context.Context, inToken, outToken string, inAmount, slippage float64, limit int32, projects []pb.Project) (*pb.GetQuotesResponse, error) {
	return g.apiClient.GetQuotes(ctx, &pb.GetQuotesRequest{InToken: inToken, OutToken: outToken, InAmount: inAmount, Slippage: slippage, Limit: limit, Projects: projects})
}

// SignAndSubmit signs the given transaction and submits it.
func (g *GRPCClient) SignAndSubmit(ctx context.Context, tx *pb.TransactionMessage,
	skipPreFlight bool, frontRunningProtection bool, useStakedRPCs bool) (string, error) {
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
	}, skipPreFlight, frontRunningProtection, useStakedRPCs)
	if err != nil {
		return "", err
	}

	return response.Signature, nil
}

// signAndSubmitBatch signs the given transactions and submits them.
func (g *GRPCClient) signAndSubmitBatch(ctx context.Context, transactions []*pb.TransactionMessage, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	if g.privateKey == nil {
		return nil, ErrPrivateKeyNotFound
	}

	if len(transactions) == 1 {
		signature, err := g.SignAndSubmit(ctx, transactions[0], *opts.SkipPreFlight, false, false)
		if err != nil {
			return nil, err
		}
		return &pb.PostSubmitBatchResponse{
			Transactions: []*pb.PostSubmitBatchResponseEntry{
				{
					Signature: signature,
					Error:     "",
					Submitted: true,
				},
			},
		}, nil
	}

	batchRequest, err := buildBatchRequest(transactions, *g.privateKey, useBundle, opts)
	if err != nil {
		return nil, err
	}

	return g.PostSubmitBatch(ctx, batchRequest)
}

// PostTradeSwap returns a partially signed transaction for submitting a swap request
func (g *GRPCClient) PostTradeSwap(ctx context.Context, ownerAddress, inToken, outToken string, inAmount, slippage float64, project pb.Project) (*pb.TradeSwapResponse, error) {
	return g.apiClient.PostTradeSwap(ctx, &pb.TradeSwapRequest{
		OwnerAddress: ownerAddress,
		InToken:      inToken,
		OutToken:     outToken,
		InAmount:     inAmount,
		Slippage:     slippage,
		Project:      project,
	})
}

// PostRouteTradeSwap returns a partially signed transaction(s) for submitting a swap request
func (g *GRPCClient) PostRouteTradeSwap(ctx context.Context, request *pb.RouteTradeSwapRequest) (*pb.TradeSwapResponse, error) {
	return g.apiClient.PostRouteTradeSwap(ctx, request)
}

// PostOrder returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (g *GRPCClient) PostOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	return g.apiClient.PostOrder(ctx, &pb.PostOrderRequest{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              types,
		Amount:            amount,
		Price:             price,
		Project:           project,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
	})
}

// PostSubmit posts the transaction string to the Solana network.
func (g *GRPCClient) PostSubmit(ctx context.Context, tx *pb.TransactionMessage, skipPreFlight bool,
	frontRunningProtection bool, useStakedRPCs bool) (*pb.PostSubmitResponse, error) {
	return g.apiClient.PostSubmit(ctx, &pb.PostSubmitRequest{Transaction: tx,
		SkipPreFlight:          skipPreFlight,
		FrontRunningProtection: &frontRunningProtection,
		UseStakedRPCs:          &useStakedRPCs,
	})
}

// PostSubmitBatch posts a bundle of transactions string based on a specific SubmitStrategy to the Solana network.
func (g *GRPCClient) PostSubmitBatch(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error) {
	return g.apiClient.PostSubmitBatch(ctx, request)
}

// PostSubmitV2 posts the transaction string to the Solana network.
func (g *GRPCClient) PostSubmitV2(ctx context.Context, tx *pb.TransactionMessage,
	skipPreFlight bool, frontRunningProtection bool, tpu uint32) (*pb.PostSubmitResponse, error) {
	return g.apiClient.PostSubmitV2(ctx, &pb.PostSubmitRequest{Transaction: tx,
		SkipPreFlight:          skipPreFlight,
		FrontRunningProtection: &frontRunningProtection,
	})
}

// PostSubmitBatchV2 posts a bundle of transactions string based on a specific SubmitStrategy to the Solana network.
func (g *GRPCClient) PostSubmitBatchV2(ctx context.Context, request *pb.PostSubmitBatchRequest) (*pb.PostSubmitBatchResponse, error) {
	return g.apiClient.PostSubmitBatchV2(ctx, request)
}

// SubmitTradeSwap builds a TradeSwap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitTradeSwap(ctx context.Context, ownerAddress, inToken, outToken string, inAmount, slippage float64, project pb.Project, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.apiClient.PostTradeSwap(ctx, &pb.TradeSwapRequest{
		OwnerAddress: ownerAddress,
		InToken:      inToken,
		OutToken:     outToken,
		InAmount:     inAmount,
		Slippage:     slippage,
		Project:      project,
	})
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// SubmitRouteTradeSwap builds a RouteTradeSwap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitRouteTradeSwap(ctx context.Context, request *pb.RouteTradeSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostRouteTradeSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// SubmitRaydiumSwap builds a Raydium Swap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitRaydiumSwap(ctx context.Context, request *pb.PostRaydiumSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostRaydiumSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// SubmitRaydiumRouteSwap builds a Raydium RouteSwap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitRaydiumRouteSwap(ctx context.Context, request *pb.PostRaydiumRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostRaydiumRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// SubmitJupiterSwap builds a Jupiter Swap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitJupiterSwap(ctx context.Context, request *pb.PostJupiterSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostJupiterSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// SubmitJupiterSwapInstructions builds a Jupiter Swap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitJupiterSwapInstructions(ctx context.Context, request *pb.PostJupiterSwapInstructionsRequest, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	swapInstructions, err := g.PostJupiterSwapInstructions(ctx, request)
	if err != nil {
		return nil, err
	}

	txBuilder := solana.NewTransactionBuilder()

	addressLookupTable, err := utils.ConvertProtoAddressLookupTable(swapInstructions.AddressLookupTableAddresses)
	if err != nil {
		return nil, err
	}

	txBuilder.WithOpt(solana.TransactionAddressTables(addressLookupTable))

	instructions, err := utils.ConvertJupiterInstructions(swapInstructions.Instructions)
	if err != nil {
		return nil, err
	}

	for _, inst := range instructions {
		txBuilder.AddInstruction(inst)
	}

	txBuilder.SetFeePayer(g.privateKey.PublicKey())
	blockHash, err := g.RecentBlockHash(ctx)

	if err != nil {
		panic(fmt.Errorf("server error: could not retrieve block hash: %w", err))
	}

	hash, err := solana.HashFromBase58(blockHash.BlockHash)
	if err != nil {
		return nil, err
	}

	txBuilder.SetRecentBlockHash(hash)
	tx, err := txBuilder.Build()
	if err != nil {
		return nil, err
	}

	err = transaction.PartialSign(tx, g.privateKey.PublicKey(), make(map[solana.PublicKey]solana.PrivateKey))
	if err != nil {
		return nil, err
	}

	var txToBeSigned []*pb.TransactionMessage

	txBase64, err := tx.ToBase64()
	if err != nil {
		return nil, err
	}

	txToBeSigned = append(txToBeSigned, &pb.TransactionMessage{
		Content:   txBase64,
		IsCleanup: false,
	})

	return g.signAndSubmitBatch(ctx, txToBeSigned, useBundle, opts)
}

// SubmitRaydiumSwapInstructions builds a Raydium Swap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitRaydiumSwapInstructions(ctx context.Context, request *pb.PostRaydiumSwapInstructionsRequest, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	swapInstructions, err := g.PostRaydiumSwapInstructions(ctx, request)
	if err != nil {
		return nil, err
	}

	instructions, err := utils.ConvertRaydiumInstructions(swapInstructions.Instructions)
	if err != nil {
		return nil, err
	}
	txBuilder := solana.NewTransactionBuilder()

	for _, inst := range instructions {
		txBuilder.AddInstruction(inst)
	}

	txBuilder.SetFeePayer(g.privateKey.PublicKey())
	blockHash, err := g.RecentBlockHash(ctx)

	if err != nil {
		panic(fmt.Errorf("server error: could not retrieve block hash: %w", err))
	}

	hash, err := solana.HashFromBase58(blockHash.BlockHash)
	if err != nil {
		return nil, err
	}

	txBuilder.SetRecentBlockHash(hash)
	tx, err := txBuilder.Build()
	if err != nil {
		return nil, err
	}

	err = transaction.PartialSign(tx, g.privateKey.PublicKey(), make(map[solana.PublicKey]solana.PrivateKey))
	if err != nil {
		return nil, err
	}

	var txToBeSigned []*pb.TransactionMessage

	txBase64, err := tx.ToBase64()
	if err != nil {
		return nil, err
	}

	txToBeSigned = append(txToBeSigned, &pb.TransactionMessage{
		Content:   txBase64,
		IsCleanup: false,
	})

	return g.signAndSubmitBatch(ctx, txToBeSigned, useBundle, opts)
}

// SubmitJupiterRouteSwap builds a Jupiter RouteSwap transaction then signs it, and submits to the network.
func (g *GRPCClient) SubmitJupiterRouteSwap(ctx context.Context, request *pb.PostJupiterRouteSwapRequest, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	resp, err := g.PostJupiterRouteSwap(ctx, request)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, resp.Transactions, false, opts)
}

// SubmitOrder builds a Serum market order, signs it, and submits to the network.
func (g *GRPCClient) SubmitOrder(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error) {
	order, err := g.PostOrder(ctx, owner, payer, market, side, types, amount, price, project, opts)
	if err != nil {
		return "", err
	}
	skipPreFlight := true
	if opts.SkipPreFlight != nil {
		skipPreFlight = *opts.SkipPreFlight
	}
	return g.SignAndSubmit(ctx, order.Transaction, skipPreFlight, false, false)
}

// PostCancelOrder builds a Serum cancel order.
func (g *GRPCClient) PostCancelOrder(
	ctx context.Context,
	orderID string,
	side pb.Side,
	owner,
	market,
	openOrders string,
	project pb.Project,
) (*pb.PostCancelOrderResponse, error) {
	return g.apiClient.PostCancelOrder(ctx, &pb.PostCancelOrderRequest{
		OrderID:           orderID,
		Side:              side,
		OwnerAddress:      owner,
		MarketAddress:     market,
		OpenOrdersAddress: openOrders,
		Project:           project,
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
	project pb.Project,
	skipPreFlight bool,
) (string, error) {
	order, err := g.PostCancelOrder(ctx, orderID, side, owner, market, openOrders, project)
	if err != nil {
		return "", err
	}

	return g.SignAndSubmit(ctx, order.Transaction, skipPreFlight, false, false)
}

// PostCancelByClientOrderID builds a Serum cancel order by client ID.
func (g *GRPCClient) PostCancelByClientOrderID(
	ctx context.Context,
	clientOrderID uint64,
	owner,
	market,
	openOrders string,
	project pb.Project,

) (*pb.PostCancelOrderResponse, error) {
	return g.apiClient.PostCancelByClientOrderID(ctx, &pb.PostCancelByClientOrderIDRequest{
		ClientOrderID:     clientOrderID,
		OwnerAddress:      owner,
		MarketAddress:     market,
		OpenOrdersAddress: openOrders,
		Project:           project,
	})
}

// SubmitCancelByClientOrderID builds a Serum cancel order by client ID, signs and submits it to the network.
func (g *GRPCClient) SubmitCancelByClientOrderID(
	ctx context.Context,
	clientOrderID uint64,
	owner,
	market,
	openOrders string,
	project pb.Project,
	skipPreFlight bool,
) (string, error) {
	order, err := g.PostCancelByClientOrderID(ctx, clientOrderID, owner, market, openOrders, project)
	if err != nil {
		return "", err
	}

	return g.SignAndSubmit(ctx, order.Transaction, skipPreFlight, false, false)
}

func (g *GRPCClient) PostCancelAll(ctx context.Context, market, owner string, openOrders []string, project pb.Project) (*pb.PostCancelAllResponse, error) {
	return g.apiClient.PostCancelAll(ctx, &pb.PostCancelAllRequest{
		Market:              market,
		OwnerAddress:        owner,
		OpenOrdersAddresses: openOrders,
		Project:             project,
	})
}

func (g *GRPCClient) SubmitCancelAll(ctx context.Context, market, owner string, openOrdersAddresses []string, project pb.Project, opts SubmitOpts) (*pb.PostSubmitBatchResponse, error) {
	orders, err := g.PostCancelAll(ctx, market, owner, openOrdersAddresses, project)
	if err != nil {
		return nil, err
	}
	return g.signAndSubmitBatch(ctx, orders.Transactions, false, opts)
}

// PostSettle returns a partially signed transaction for settling market funds. Typically, you want to use SubmitSettle instead of this.
func (g *GRPCClient) PostSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, project pb.Project) (*pb.PostSettleResponse, error) {
	return g.apiClient.PostSettle(ctx, &pb.PostSettleRequest{
		OwnerAddress:      owner,
		Market:            market,
		BaseTokenWallet:   baseTokenWallet,
		QuoteTokenWallet:  quoteTokenWallet,
		OpenOrdersAddress: openOrdersAccount,
		Project:           project,
	})
}

// SubmitSettle builds a market SubmitSettle transaction, signs it, and submits to the network.
func (g *GRPCClient) SubmitSettle(ctx context.Context, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount string, project pb.Project, skipPreflight bool) (string, error) {
	order, err := g.PostSettle(ctx, owner, market, baseTokenWallet, quoteTokenWallet, openOrdersAccount, project)
	if err != nil {
		return "", err
	}

	return g.SignAndSubmit(ctx, order.Transaction, skipPreflight, false, false)
}

func (g *GRPCClient) PostReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	return g.apiClient.PostReplaceByClientOrderID(ctx, &pb.PostOrderRequest{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              types,
		Amount:            amount,
		Price:             price,
		Project:           project,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
	})
}

func (g *GRPCClient) SubmitReplaceByClientOrderID(ctx context.Context, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error) {
	order, err := g.PostReplaceByClientOrderID(ctx, owner, payer, market, side, types, amount, price, project, opts)
	if err != nil {
		return "", err
	}
	skipPreFlight := true
	if opts.SkipPreFlight != nil {
		skipPreFlight = *opts.SkipPreFlight
	}
	return g.SignAndSubmit(ctx, order.Transaction, skipPreFlight, false, false)
}

func (g *GRPCClient) PostReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	return g.apiClient.PostReplaceOrder(ctx, &pb.PostReplaceOrderRequest{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              types,
		Amount:            amount,
		Price:             price,
		Project:           project,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
		OrderID:           orderID,
	})
}

func (g *GRPCClient) SubmitReplaceOrder(ctx context.Context, orderID, owner, payer, market string, side pb.Side, types []common.OrderType, amount, price float64, project pb.Project, opts PostOrderOpts) (string, error) {
	order, err := g.PostReplaceOrder(ctx, orderID, owner, payer, market, side, types, amount, price, project, opts)
	if err != nil {
		return "", err
	}
	skipPreFlight := true
	if opts.SkipPreFlight != nil {
		skipPreFlight = *opts.SkipPreFlight
	}
	return g.SignAndSubmit(ctx, order.Transaction, skipPreFlight, false, false)
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
func (g *GRPCClient) GetPoolReservesStream(ctx context.Context, request *pb.GetPoolReservesStreamRequest) (connections.Streamer[*pb.GetPoolReservesStreamResponse], error) {
	stream, err := g.apiClient.GetPoolReservesStream(ctx, request)
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

// GetTickersStream subscribes to a stream for getting recent tickers of specified markets.
func (g *GRPCClient) GetTickersStream(ctx context.Context, request *pb.GetTickersStreamRequest) (connections.Streamer[*pb.GetTickersStreamResponse], error) {
	stream, err := g.apiClient.GetTickersStream(ctx, request)
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetTickersStreamResponse](stream, ""), nil
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

// GetNewRaydiumPoolsStream subscribes to a stream for getting recent swaps on projects & markets of interest with
// option to include Raydium cpmm amm.
func (g *GRPCClient) GetNewRaydiumPoolsStream(
	ctx context.Context, includeCPMM bool,
) (connections.Streamer[*pb.GetNewRaydiumPoolsResponse], error) {
	stream, err := g.apiClient.GetNewRaydiumPoolsStream(ctx, &pb.GetNewRaydiumPoolsRequest{
		IncludeCPMM: &includeCPMM,
	})
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

// GetPriorityFeeStream subscribes to a stream of priority fees for a given percentile
func (g *GRPCClient) GetPriorityFeeStream(ctx context.Context, project pb.Project, percentile *float64) (connections.Streamer[*pb.GetPriorityFeeResponse], error) {
	request := &pb.GetPriorityFeeRequest{
		Project: project,
	}
	if percentile != nil {
		request.Percentile = percentile
	}
	stream, err := g.apiClient.GetPriorityFeeStream(ctx, request)
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetPriorityFeeResponse](stream, fmt.Sprint(percentile)), nil
}

// GetBundleTipStream subscribes to a stream of bundle tip percentiles
func (g *GRPCClient) GetBundleTipStream(ctx context.Context) (connections.Streamer[*pb.GetBundleTipResponse], error) {
	stream, err := g.apiClient.GetBundleTipStream(ctx, &pb.GetBundleTipRequest{})
	if err != nil {
		return nil, err
	}

	return connections.GRPCStream[pb.GetBundleTipResponse](stream, ""), nil
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
func (g *GRPCClient) PostOrderV2(ctx context.Context, owner, payer, market string, side string, orderType string, amount, price float64, bundleTip *uint64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
	return g.apiClient.PostOrderV2(ctx, &pb.PostOrderRequestV2{
		OwnerAddress:      owner,
		PayerAddress:      payer,
		Market:            market,
		Side:              side,
		Type:              orderType,
		Amount:            amount,
		Price:             price,
		Tip:               bundleTip,
		OpenOrdersAddress: opts.OpenOrdersAddress,
		ClientOrderID:     opts.ClientOrderID,
	})
}

// PostOrderV2WithPriorityFee returns a partially signed transaction for placing a Serum market order. Typically, you want to use SubmitOrder instead of this.
func (g *GRPCClient) PostOrderV2WithPriorityFee(ctx context.Context, owner, payer, market string, side string,
	orderType string, amount, price float64, computeLimit uint32, computePrice uint64, bundleTip *uint64, opts PostOrderOpts) (*pb.PostOrderResponse, error) {
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
		Tip:               bundleTip,
		ClientOrderID:     opts.ClientOrderID,
	})
}

// SubmitOrderV2 builds a Serum market order, signs it, and submits to the network.
func (g *GRPCClient) SubmitOrderV2(ctx context.Context, owner, payer, market string, side string, orderType string, amount, price float64, bundleTip *uint64, opts PostOrderOpts) (string, error) {
	order, err := g.PostOrderV2(ctx, owner, payer, market, side, orderType, amount, price, bundleTip, opts)
	if err != nil {
		return "", err
	}

	skipPreFlight := true
	if opts.SkipPreFlight != nil {
		skipPreFlight = *opts.SkipPreFlight
	}
	return g.SignAndSubmit(ctx, order.Transaction, skipPreFlight, false, false)
}

// SubmitOrderV2WithPriorityFee builds a Serum market order, signs it, and submits to the network with specified computeLimit and computePrice
func (g *GRPCClient) SubmitOrderV2WithPriorityFee(ctx context.Context, owner, payer, market string, side string,
	orderType string, amount, price float64, computeLimit uint32, computePrice uint64, bundleTip *uint64, opts PostOrderOpts) (string, error) {
	order, err := g.PostOrderV2WithPriorityFee(ctx, owner, payer, market, side, orderType, amount, price, computeLimit, computePrice, bundleTip, opts)
	if err != nil {
		return "", err
	}

	skipPreFlight := true
	if opts.SkipPreFlight != nil {
		skipPreFlight = *opts.SkipPreFlight
	}
	return g.SignAndSubmit(ctx, order.Transaction, skipPreFlight, false, false)
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

	return g.signAndSubmitBatch(ctx, order.Transactions, false, opts)
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

	return g.SignAndSubmit(ctx, order.Transaction, skipPreflight, false, false)
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
	skipPreFlight := true
	if opts.SkipPreFlight != nil {
		skipPreFlight = *opts.SkipPreFlight
	}
	return g.SignAndSubmit(ctx, order.Transaction, skipPreFlight, false, false)
}
