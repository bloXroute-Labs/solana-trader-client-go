package examples

import (
	"context"
	"errors"

	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	log "github.com/sirupsen/logrus"
)

func GetPumpFunNewTokenHelper() (*pb.GetPumpFunNewTokensStreamResponse, error) {
	grpcClient, err := provider.NewGRPCClientPumpNY()
	if err != nil {
		panic(err)
	}
	log.Info("starting GetPumpFunNewTokens stream")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := grpcClient.GetPumpFunNewTokensStream(ctx, &pb.GetPumpFunNewTokensStreamRequest{})
	if err != nil {
		log.Errorf("error with GetPumpFunNewTokens stream request: %v", err)
		return nil, err
	}

	ch := stream.Channel(0)

	// Wait for a single response
	v, ok := <-ch
	if !ok {
		return nil, errors.New("Token searcher failed")
	}

	log.Infof("response %v received", v)
	return v, nil
}
