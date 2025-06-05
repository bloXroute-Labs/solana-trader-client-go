package provider

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gagliardetto/solana-go"

	"github.com/bloXroute-Labs/solana-trader-client-go/transaction"
	"github.com/bloXroute-Labs/solana-trader-client-go/utils"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
)

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

type PostSubmitOpts struct {
	SkipPreFlight          bool
	FrontRunningProtection bool
	UseStakedRPCs          bool
	AllowBackRun           bool
	RevenueAddress         string
	Sniping                bool
	AllowRevert            bool
}

type RPCOpts struct {
	Endpoint        string
	DisableAuth     bool
	UseTLS          bool
	PrivateKey      *solana.PrivateKey
	AuthHeader      string
	DisablePingLoop bool
	CacheBlockHash  bool
	BlockHashTtl    time.Duration
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
	batchRequest.Timestamp = utils.GetTimestamp()

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

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
