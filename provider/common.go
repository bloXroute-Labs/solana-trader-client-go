package provider

import (
	"errors"
	"fmt"
	api "github.com/bloXroute-Labs/solana-trader-client-go/proto"
	pb "github.com/bloXroute-Labs/solana-trader-client-go/proto"
	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	"os"
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

	for _, tx := range transactions.([]interface{}) {

		oneRequest := pb.PostSubmitRequestEntry{}
		oneRequest.SkipPreFlight = opts.SkipPreFlight

		if txStr, ok := tx.(string); ok {
			signedTxBase64, err := transaction.SignTxWithPrivateKey(txStr, privateKey)
			if err != nil {
				return &pb.PostSubmitBatchRequest{}, err
			}
			oneRequest.Transaction = &pb.TransactionMessage{
				Content: signedTxBase64,
			}
		} else if txMsg, ok := tx.(*pb.TransactionMessage); ok {
			signedTxBase64, err := transaction.SignTxWithPrivateKey(txMsg.Content, privateKey)
			if err != nil {
				return &pb.PostSubmitBatchRequest{}, err
			}
			oneRequest.Transaction = &pb.TransactionMessage{
				Content:   signedTxBase64,
				IsCleanup: txMsg.IsCleanup,
			}
		}

		batchRequest.Entries = append(batchRequest.Entries, &oneRequest)
	}
	return &batchRequest, nil
}
