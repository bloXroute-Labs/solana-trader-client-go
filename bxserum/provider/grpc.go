package provider

import (
	"context"
	"errors"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	"github.com/bloXroute-Labs/serum-api/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
)

type GRPCClient struct {
	pb.UnsafeApiServer // TODO Regular API Server?
	pb.UnimplementedApiServer

	grpcConn  *grpc.ClientConn
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
		grpcConn:  conn,
		requestID: utils.NewRequestID(),
	}, nil
}

func (g *GRPCClient) GetOrderbook(ctx context.Context, market string) (*pb.GetOrderbookResponse, error) {
	method := "/api.Api/GetOrderbook"
	in := &pb.GetOrderBookRequest{Market: market}
	out := new(pb.GetOrderbookResponse)
	return grpcResponse[pb.GetOrderbookResponse](ctx, g.grpcConn, method, in, out)
}

func grpcResponse[T any](ctx context.Context, client *grpc.ClientConn, method string, in interface{}, out *T) (*T, error) {
	if client == nil {
		return nil, errors.New("client is nil, please create one using a `NewGRPCClient` function")
	}
	err := client.Invoke(ctx, method, in, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func grpcStream[T any](ctx context.Context, client *grpc.ClientConn, method string, in interface{}, out chan<- *T) error {
	if client == nil {
		return errors.New("client is nil, please create one using a `NewGRPCClient` function")
	}
	stream, err := client.NewStream(ctx, nil, method)
	if err != nil {
		return err
	}
	if err := stream.SendMsg(in); err != nil {
		return err
	}

	go func() {
		for {
			output := new(T)
			err := stream.RecvMsg(&output)
			if err == io.EOF {
				log.Errorf("stream for method %s ended successfully", method)
			} else if err != nil {
				log.Errorf("error when receiving message for method %s - %v", method, err)
			} else {
				out <- output
			}
		}
	}()

	return nil
}
