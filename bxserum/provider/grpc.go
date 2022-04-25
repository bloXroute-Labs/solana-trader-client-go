package provider

import (
	"context"
	"fmt"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/bloXroute-Labs/serum-api/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
)

type GRPCClient struct {
	pb.UnimplementedApiServer

	apiClient pb.ApiClient
	requestID utils.RequestID
}

// Connects to Mainnet Serum API
func NewGRPCClient() (*GRPCClient, error) {
	return NewGRPCClientWithEndpoint("174.129.154.164:1811")
}

// Connects to Testnet Serum API
func NewGRPCTestnet() (*GRPCClient, error) {
	panic("implement me")
}

// Connects to custom Serum API
func NewGRPCClientWithEndpoint(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &GRPCClient{
		apiClient: pb.NewApiClient(conn),
		requestID: utils.NewRequestID(),
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
	return streamResponse[pb.GetOrderbookStreamResponse](stream, market, outputChan)

}

func (g *GRPCClient) GetMarkets(ctx context.Context) (*pb.GetMarketsResponse, error) {
	return g.apiClient.GetMarkets(ctx, &pb.GetMarketsRequest{})
}

func streamResponse[T any](stream grpc.ClientStream, input string, outputChan chan *T) error {
	val, err := readGRPCStream[T](stream, input)
	if err != nil {
		return err
	}
	outputChan <- val

	go func(stream grpc.ClientStream, input string) {
		for {
			val, err = readGRPCStream[T](stream, input)
			if err != nil {
				log.Errorf(err.Error())
				return
			} else {
				outputChan <- val
			}
		}
	}(stream, input)

	return nil
}

func readGRPCStream[T any](stream grpc.ClientStream, input string) (*T, error) {
	m := new(T)
	err := stream.RecvMsg(m)
	if err == io.EOF {
		return nil, fmt.Errorf("stream for input %s ended successfully", input)
	} else if err != nil {
		return nil, err
	}

	return m, nil
}
