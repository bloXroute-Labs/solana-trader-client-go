package provider

import (
	"errors"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"os"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
)

const (
	mainnetNY     = "ny.solana.dex.blxrbdn.com"
	mainnetPumpNY = "pump-ny.solana.dex.blxrbdn.com"
	mainnetUK     = "uk.solana.dex.blxrbdn.com"
	testnet       = "solana.dex.bxrtest.com"
	devnet        = "solana-trader-api-nlb-6b0f765f2fc759e1.elb.us-east-1.amazonaws.com"
)

var (
	MainnetNYHTTP     = httpEndpoint(mainnetNY, true)
	MainnetUKHTTP     = httpEndpoint(mainnetUK, true)
	MainnetNYWS       = wsEndpoint(mainnetNY, true)
	MainnetPumpNYWS   = wsEndpoint(mainnetPumpNY, true)
	MainnetUKWS       = wsEndpoint(mainnetUK, true)
	MainnetNYGRPC     = grpcEndpoint(mainnetNY, true)
	MainnetPumpNYGRPC = grpcEndpoint(mainnetPumpNY, true)
	MainnetUKGRPC     = grpcEndpoint(mainnetUK, true)

	TestnetHTTP = httpEndpoint(testnet, true)
	TestnetWS   = wsEndpoint(testnet, true)
	TestnetGRPC = grpcEndpoint(testnet, true)

	DevnetHTTP = httpEndpoint(devnet, false)
	DevnetWS   = wsEndpoint(devnet, false)
	DevnetGRPC = grpcEndpoint(devnet, false)

	LocalHTTP = "http://localhost:9000"
	LocalWS   = "ws://localhost:9000/ws"
	LocalGRPC = "localhost:9000"
)

func httpEndpoint(baseUrl string, secure bool) string {
	prefix := "http"
	if secure {
		prefix = "https"
	}
	return fmt.Sprintf("%v://%v", prefix, baseUrl)
}

func wsEndpoint(baseUrl string, secure bool) string {
	prefix := "ws"
	if secure {
		prefix = "wss"
	}
	return fmt.Sprintf("%v://%v/ws", prefix, baseUrl)
}

func grpcEndpoint(baseUrl string, secure bool) string {
	port := "80"
	if secure {
		port = "443"
	}
	return fmt.Sprintf("%v:%v", baseUrl, port)
}

var ErrPrivateKeyNotFound = errors.New("private key not provided for signing transaction")

type PostOrderOpts struct {
	OpenOrdersAddress string
	ClientOrderID     uint64
	SkipPreFlight     *bool
}

type SubmitOpts struct {
	SubmitStrategy pb.SubmitStrategy
	SkipPreFlight  *bool
}

type RPCOpts struct {
	Endpoint       string
	DisableAuth    bool
	UseTLS         bool
	PrivateKey     *solana.PrivateKey
	AuthHeader     string
	CacheBlockHash bool
	BlockHashTtl   time.Duration
}

func DefaultRPCOpts(endpoint string) RPCOpts {
	var spk *solana.PrivateKey
	privateKey, err := transaction.LoadPrivateKeyFromEnv()
	if err == nil {
		spk = &privateKey
	}
	return RPCOpts{
		Endpoint:   endpoint,
		PrivateKey: spk,
		AuthHeader: os.Getenv("AUTH_HEADER"),
	}
}

var stringToAmm = map[string]pb.Project{
	"unknown": pb.Project_P_UNKNOWN,
	"jupiter": pb.Project_P_JUPITER,
	"raydium": pb.Project_P_RAYDIUM,
	"all":     pb.Project_P_ALL,
}

func ProjectFromString(project string) (pb.Project, error) {
	if apiProject, ok := stringToAmm[strings.ToLower(project)]; ok {
		return apiProject, nil
	}

	return pb.Project_P_UNKNOWN, fmt.Errorf("could not find project %s", project)
}

func buildBatchRequest(transactions []*pb.TransactionMessage, privateKey solana.PrivateKey, useBundle bool, opts SubmitOpts) (*pb.PostSubmitBatchRequest, error) {
	batchRequest := pb.PostSubmitBatchRequest{}
	batchRequest.SubmitStrategy = opts.SubmitStrategy

	for _, tx := range transactions {
		request, err := createBatchRequestEntry(opts, tx.Content, privateKey)
		if err != nil {
			return nil, err
		}

		batchRequest.Entries = append(batchRequest.Entries, request)

	}

	batchRequest.UseBundle = &useBundle

	return &batchRequest, nil
}

func createBatchRequestEntry(opts SubmitOpts, txBase64 string, privateKey solana.PrivateKey) (*pb.PostSubmitRequestEntry, error) {
	oneRequest := pb.PostSubmitRequestEntry{}
	if opts.SkipPreFlight == nil {
		oneRequest.SkipPreFlight = true
	} else {
		oneRequest.SkipPreFlight = *opts.SkipPreFlight
	}

	signedTxBase64, err := transaction.SignTxWithPrivateKey(txBase64, privateKey)
	if err != nil {
		return nil, err
	}
	oneRequest.Transaction = &pb.TransactionMessage{
		Content: signedTxBase64,
	}
	return &oneRequest, nil
}
