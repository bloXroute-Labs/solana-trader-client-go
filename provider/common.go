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

const WarningTLSSlowDown = "Performance Notice: Secure (TLS) endpoints may introduce latency due to handshake overhead. For optimal trading speed, consider using non-secure endpoints when appropriate."

type Region string

// for information about submit only regions, see documentation: https://docs.bloxroute.com/solana/trader-api/introduction/regions
const (
	mainnetNY        Region = "ny.solana.dex.blxrbdn.com"
	mainnetPumpNY    Region = "pump-ny.solana.dex.blxrbdn.com"
	mainnetUK        Region = "uk.solana.dex.blxrbdn.com"
	mainnetPumpUK    Region = "pump-uk.solana.dex.blxrbdn.com"
	mainnetFrankfurt Region = "germany.solana.dex.blxrbdn.com"
	mainnetLA        Region = "la.solana.dex.blxrbdn.com"
	mainnetAmsterdam Region = "amsterdam.solana.dex.blxrbdn.com"
	mainnetTokyo     Region = "tokyo.solana.dex.blxrbdn.com"
	testnet          Region = "160.202.128.145"
	devnet           Region = "solana-trader-api-nlb-6b0f765f2fc759e1.elb.us-east-1.amazonaws.com"
)

var (

	// FULL SERVICE

	// http
	MainnetNYHTTP     = httpEndpoint(mainnetNY, false)
	MainnetPumpNYHTTP = httpEndpoint(mainnetPumpNY, false)
	MainnetUKHTTP     = httpEndpoint(mainnetUK, false)
	MainnetPumpUKHTTP = httpEndpoint(mainnetPumpUK, false)

	// grpc
	MainnetNYGRPC     = grpcEndpoint(mainnetNY, false)
	MainnetPumpNYGRPC = grpcEndpoint(mainnetPumpNY, false)
	MainnetUKGRPC     = grpcEndpoint(mainnetUK, false)
	MainnetPumpUKGRPC = grpcEndpoint(mainnetPumpUK, false)

	// ws
	MainnetNYWS     = wsEndpoint(mainnetNY, false)
	MainnetPumpNYWS = wsEndpoint(mainnetPumpNY, false)
	MainnetUKWS     = wsEndpoint(mainnetUK, false)
	MainnetPumpUKWS = wsEndpoint(mainnetPumpUK, false)

	// SUBMIT ONLY

	// http
	MainnetFrankfurtHTTP = httpEndpoint(mainnetFrankfurt, false)
	MainnetLAHTTP        = httpEndpoint(mainnetLA, false)
	MainnetAmsterdamHTTP = httpEndpoint(mainnetAmsterdam, false)
	MainnetTokyoHTTP     = httpEndpoint(mainnetTokyo, false)

	// grpc
	MainnetFrankfurtGRPC = grpcEndpoint(mainnetFrankfurt, false)
	MainnetLAGRPC        = grpcEndpoint(mainnetLA, false)
	MainnetAmsterdamGRPC = grpcEndpoint(mainnetAmsterdam, false)
	MainnetTokyoGRPC     = grpcEndpoint(mainnetTokyo, false)

	// ws
	MainnetFrankfurtWS = wsEndpoint(mainnetFrankfurt, false)
	MainnetLAWS        = wsEndpoint(mainnetLA, false)
	MainnetAmsterdamWS = wsEndpoint(mainnetAmsterdam, false)
	MainnetTokyoWS     = wsEndpoint(mainnetTokyo, false)

	// TESTING

	// testnet
	TestnetHTTP = httpEndpoint(testnet, false)
	TestnetWS   = wsEndpoint(testnet, false)
	TestnetGRPC = grpcEndpoint(testnet, false)

	// devnet
	DevnetHTTP = httpEndpoint(devnet, false)
	DevnetWS   = wsEndpoint(devnet, false)
	DevnetGRPC = grpcEndpoint(devnet, false)

	// local
	LocalHTTP = "http://localhost:9001"
	LocalWS   = "ws://localhost:9001/ws"
	LocalGRPC = "localhost:9001"
)

var (

	// FULL SERVICE

	// http
	MainnetNYHTTPSecure     = httpEndpoint(mainnetNY, true)
	MainnetPumpNYHTTPSecure = httpEndpoint(mainnetPumpNY, true)
	MainnetUKHTTPSecure     = httpEndpoint(mainnetUK, true)
	MainnetPumpUKHTTPSecure = httpEndpoint(mainnetPumpUK, true)

	// grpc
	MainnetNYGRPCSecure     = grpcEndpoint(mainnetNY, true)
	MainnetPumpNYGRPCSecure = grpcEndpoint(mainnetPumpNY, true)
	MainnetUKGRPCSecure     = grpcEndpoint(mainnetUK, true)
	MainnetPumpUKGRPCSecure = grpcEndpoint(mainnetPumpUK, true)

	// ws
	MainnetNYWSSecure     = wsEndpoint(mainnetNY, true)
	MainnetPumpNYWSSecure = wsEndpoint(mainnetPumpNY, true)
	MainnetUKWSSecure     = wsEndpoint(mainnetUK, true)
	MainnetPumpUKWSSecure = wsEndpoint(mainnetPumpUK, true)

	// SUBMIT ONLY

	// http
	MainnetFrankfurtHTTPSecure = httpEndpoint(mainnetFrankfurt, true)
	MainnetLAHTTPSecure        = httpEndpoint(mainnetLA, true)
	MainnetAmsterdamHTTPSecure = httpEndpoint(mainnetAmsterdam, true)
	MainnetTokyoHTTPSecure     = httpEndpoint(mainnetTokyo, true)

	// grpc
	MainnetFrankfurtGRPCSecure = grpcEndpoint(mainnetFrankfurt, true)
	MainnetLAGRPCSecure        = grpcEndpoint(mainnetLA, true)
	MainnetAmsterdamGRPCSecure = grpcEndpoint(mainnetAmsterdam, true)
	MainnetTokyoGRPCSecure     = grpcEndpoint(mainnetTokyo, true)

	// ws
	MainnetFrankfurtWSSecure = wsEndpoint(mainnetFrankfurt, true)
	MainnetLAWSSecure        = wsEndpoint(mainnetLA, true)
	MainnetAmsterdamWSSecure = wsEndpoint(mainnetAmsterdam, true)
	MainnetTokyoWSSecure     = wsEndpoint(mainnetTokyo, true)

	// TESTING

	// testnet
	TestnetHTTPSecure = httpEndpoint(testnet, true)
	TestnetWSSecure   = wsEndpoint(testnet, true)
	TestnetGRPCSecure = grpcEndpoint(testnet, true)

	// devnet
	DevnetHTTPSecure = httpEndpoint(devnet, true)
	DevnetWSSecure   = wsEndpoint(devnet, true)
	DevnetGRPCSecure = grpcEndpoint(devnet, true)
)

func httpEndpoint(baseUrl Region, secure bool) string {
	prefix := "http"
	if secure {
		prefix = "https"
	}
	return fmt.Sprintf("%v://%v", prefix, baseUrl)
}

func wsEndpoint(baseUrl Region, secure bool) string {
	prefix := "ws"
	if secure {
		prefix = "wss"
	}
	return fmt.Sprintf("%v://%v/ws", prefix, baseUrl)
}

func grpcEndpoint(baseUrl Region, secure bool) string {
	port := "1809"
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
	SkipPreFlight          bool
	FrontRunningProtection bool
	UseStakedRPCs          bool
	AllowBackRun           bool
	RevenueAddress         string
	AllowRevert            bool
	FastBestEffort         bool
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

func buildBatchRequest(transactions []*pb.TransactionMessage, privateKey solana.PrivateKey, frp bool, opts SubmitOpts) (*pb.PostSubmitBatchRequest, error) {
	batchRequest := pb.PostSubmitBatchRequest{}

	for _, tx := range transactions {
		request, err := createBatchRequestEntry(opts, tx.Content, privateKey)
		if err != nil {
			return nil, err
		}

		batchRequest.Entries = append(batchRequest.Entries, request)

	}

	batchRequest.FrontRunningProtection = &frp
	batchRequest.Timestamp = utils.GetTimestamp()

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

func isSecureGRPC(bloxrouteEndpoint string) bool {
	return strings.HasSuffix(bloxrouteEndpoint, "443")
}
func isSecureHTTP(bloxrouteEndpoint string) bool {
	return strings.HasPrefix(bloxrouteEndpoint, "https://")
}
func isSecureWS(bloxrouteEndpoint string) bool {
	return strings.HasPrefix(bloxrouteEndpoint, "ws://")
}
