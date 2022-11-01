package provider

import (
	"errors"
	"fmt"
	api "github.com/bloXroute-Labs/solana-trader-client-go/proto"
	pb "github.com/bloXroute-Labs/solana-trader-client-go/proto"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	log "github.com/sirupsen/logrus"

	"os"
	"reflect"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"
)

const (
	MainnetHTTP       = "https://virginia.solana.dex.blxrbdn.com"
	MainnetWS         = "wss://virginia.solana.dex.blxrbdn.com/ws"
	MainnetGRPC       = "virginia.solana.dex.blxrbdn.com:443"
	TestnetHTTP       = "http://serum-nlb-5a2c3912804344a3.elb.us-east-1.amazonaws.com"
	TestnetWS         = "ws://serum-nlb-5a2c3912804344a3.elb.us-east-1.amazonaws.com/ws"
	TestnetGRPC       = "serum-nlb-5a2c3912804344a3.elb.us-east-1.amazonaws.com:80"
	DevnetHTTP        = "http://serum-nlb-53baf45ef9775263.elb.us-east-1.amazonaws.com"
	DevnetWS          = "ws://serum-nlb-53baf45ef9775263.elb.us-east-1.amazonaws.com/ws"
	DevnetGRPC        = "serum-nlb-53baf45ef9775263.elb.us-east-1.amazonaws.com:80"
	LocalWS           = "ws://localhost:9000/ws"
	LocalGRPC         = "localhost:9000"
	LocalHTTP         = "http://127.0.0.1:9000"
	defaultRPCTimeout = 7 * time.Second
)

var ErrPrivateKeyNotFound = errors.New("private key not provided for signing transaction")

type PostOrderOpts struct {
	OpenOrdersAddress string
	ClientOrderID     uint64
	SkipPreFlight     bool
}

type SubmitOpts struct {
	SubmitStrategy pb.SubmitStrategy
	SkipPreFlight  bool
}

type RPCOpts struct {
	Endpoint       string
	UseTLS         bool
	Timeout        time.Duration
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
		Timeout:    defaultRPCTimeout,
		PrivateKey: spk,
		AuthHeader: os.Getenv("AUTH_HEADER"),
	}
}

var stringToAmm = map[string]api.Project{
	"unknown": api.Project_P_UNKNOWN,
	"jupiter": api.Project_P_JUPITER,
	"raydium": api.Project_P_RAYDIUM,
	"all":     api.Project_P_ALL,
}

func ProjectFromString(project string) (api.Project, error) {
	if apiProject, ok := stringToAmm[strings.ToLower(project)]; ok {
		return apiProject, nil
	}

	return api.Project_P_UNKNOWN, fmt.Errorf("could not find project %s", project)
}

func buildBatchRequest(transactions interface{}, privateKey solana.PrivateKey, opts SubmitOpts) (*pb.PostSubmitBatchRequest, error) {
	batchRequest := pb.PostSubmitBatchRequest{}
	batchRequest.SubmitStrategy = opts.SubmitStrategy

	switch transactions.(type) {
	case []*api.TransactionMessage:
		for _, txMsg := range transactions.([]*api.TransactionMessage) {
			oneRequest, err := createBatchRequestEntry(opts, txMsg.Content, privateKey)
			if err != nil {
				return nil, err
			}

			batchRequest.Entries = append(batchRequest.Entries, oneRequest)
		}
	case []string:
		for _, txStr := range transactions.([]string) {
			oneRequest, err := createBatchRequestEntry(opts, txStr, privateKey)
			if err != nil {
				return nil, err
			}

			batchRequest.Entries = append(batchRequest.Entries, oneRequest)
		}
	default:
		log.Fatalf("transactions has unknown type : %v", reflect.TypeOf(transactions))
	}

	return &batchRequest, nil
}

func createBatchRequestEntry(opts SubmitOpts, txBase64 string, privateKey solana.PrivateKey) (*pb.PostSubmitRequestEntry, error) {
	oneRequest := pb.PostSubmitRequestEntry{}
	oneRequest.SkipPreFlight = opts.SkipPreFlight
	signedTxBase64, err := transaction.SignTxWithPrivateKey(txBase64, privateKey)
	if err != nil {
		return nil, err
	}
	oneRequest.Transaction = &pb.TransactionMessage{
		Content: signedTxBase64,
	}
	return &oneRequest, nil
}
