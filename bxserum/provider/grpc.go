package provider

import (
	"context"
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
func NewGRPCClientWithEndpoint(baseURL string) (*GRPCClient, error) {
	conn, err := grpc.Dial(baseURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &GRPCClient{
		apiClient: pb.NewApiClient(conn),
		requestID: utils.NewRequestID(),
	}, nil
}

func (g *GRPCClient) GetOrderbook(ctx context.Context, market string) (*pb.GetOrderbookResponse, error) {
	return g.apiClient.GetOrderbook(ctx, &pb.GetOrderBookRequest{Market: market})
}

func (g *GRPCClient) GetOrderbookStream(ctx context.Context, market string, outputChan chan *pb.GetOrderbookStreamResponse) error {
	stream, err := g.apiClient.GetOrderbookStream(ctx, &pb.GetOrderBookRequest{Market: market})
	if err != nil {
		return err
	}
	go streamResponse[pb.GetOrderbookStreamResponse](stream, market, outputChan)
	return nil
}

func streamResponse[T any](stream grpc.ClientStream, input string, outputChan chan *T) {
	for {
		output := new(T)
		err := stream.RecvMsg(output)
		if err == io.EOF {
			log.Errorf("stream for input %s ended successfully", input)
			return
		} else if err != nil {
			log.Errorf("error when receiving message for input %s - %v", input, err)
			return
		} else {
			outputChan <- output
		}
	}
}
