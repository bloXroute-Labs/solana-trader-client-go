package provider

import (
	"fmt"
)

// Warning to display when creating TLS client
const WarningTLSSlowDown = "Performance Notice: Secure (TLS) endpoints may introduce latency due to handshake overhead. For optimal trading speed, consider using non-secure endpoints when appropriate."

// Endpoint type as before
type Endpoint string

// Region DNS struct
type Region struct {
	HTTP, HTTPSecure         string
	WS, WSSecure             string
	GRPC, GRPCSecure         string
	PumpHTTP, PumpHTTPSecure string
	PumpWS, PumpWSSecure     string
	PumpGRPC, PumpGRPCSecure string
}

// Remote endpoints
const (
	NY         Endpoint = "ny.solana.dex.blxrbdn.com"
	NYPump     Endpoint = "pump-ny.solana.dex.blxrbdn.com"
	UK         Endpoint = "uk.solana.dex.blxrbdn.com"
	UKPump     Endpoint = "pump-uk.solana.dex.blxrbdn.com"
	Frankfurt  Endpoint = "germany.solana.dex.blxrbdn.com"
	LosAngeles Endpoint = "la.solana.dex.blxrbdn.com"
	Amsterdam  Endpoint = "amsterdam.solana.dex.blxrbdn.com"
	Tokyo      Endpoint = "tokyo.solana.dex.blxrbdn.com"
)

// Developer endpoints
const (
	Testnet Endpoint = "solana.dex.bxrtest.com"
	Devnet  Endpoint = "solana-trader-api-nlb-6b0f765f2fc759e1.elb.us-east-1.amazonaws.com"
)

// Local endpoints
const (
	LocalHTTP string = "http://localhost:9000"
	LocalWS   string = "ws://localhost:9000/ws"
	LocalGRPC string = "localhost:9000"
)

var endpointsByRegion = map[string]Region{
	"NY": {
		HTTP:       httpEndpoint(NY, false),
		HTTPSecure: httpEndpoint(NY, true),
		WS:         wsEndpoint(NY, false),
		WSSecure:   wsEndpoint(NY, true),
		GRPC:       grpcEndpoint(NY, false),
		GRPCSecure: grpcEndpoint(NY, true),

		PumpHTTP:       httpEndpoint(NYPump, false),
		PumpHTTPSecure: httpEndpoint(NYPump, true),
		PumpWS:         wsEndpoint(NYPump, false),
		PumpWSSecure:   wsEndpoint(NYPump, true),
		PumpGRPC:       grpcEndpoint(NYPump, false),
		PumpGRPCSecure: grpcEndpoint(NYPump, true),
	},
	"UK": {
		HTTP:       httpEndpoint(UK, false),
		HTTPSecure: httpEndpoint(UK, true),
		WS:         wsEndpoint(UK, false),
		WSSecure:   wsEndpoint(UK, true),
		GRPC:       grpcEndpoint(UK, false),
		GRPCSecure: grpcEndpoint(UK, true),

		PumpHTTP:       httpEndpoint(UKPump, false),
		PumpHTTPSecure: httpEndpoint(UKPump, true),
		PumpWS:         wsEndpoint(UKPump, false),
		PumpWSSecure:   wsEndpoint(UKPump, true),
		PumpGRPC:       grpcEndpoint(UKPump, false),
		PumpGRPCSecure: grpcEndpoint(UKPump, true),
	},
	"Frankfurt": {
		HTTP:       httpEndpoint(Frankfurt, false),
		HTTPSecure: httpEndpoint(Frankfurt, true),
		WS:         wsEndpoint(Frankfurt, false),
		WSSecure:   wsEndpoint(Frankfurt, true),
		GRPC:       grpcEndpoint(Frankfurt, false),
		GRPCSecure: grpcEndpoint(Frankfurt, true),
	},
	"LosAngeles": {
		HTTP:       httpEndpoint(LosAngeles, false),
		HTTPSecure: httpEndpoint(LosAngeles, true),
		WS:         wsEndpoint(LosAngeles, false),
		WSSecure:   wsEndpoint(LosAngeles, true),
		GRPC:       grpcEndpoint(LosAngeles, false),
		GRPCSecure: grpcEndpoint(LosAngeles, true),
	},
	"Amsterdam": {
		HTTP:       httpEndpoint(Amsterdam, false),
		HTTPSecure: httpEndpoint(Amsterdam, true),
		WS:         wsEndpoint(Amsterdam, false),
		WSSecure:   wsEndpoint(Amsterdam, true),
		GRPC:       grpcEndpoint(Amsterdam, false),
		GRPCSecure: grpcEndpoint(Amsterdam, true),
	},
	"Tokyo": {
		HTTP:       httpEndpoint(Tokyo, false),
		HTTPSecure: httpEndpoint(Tokyo, true),
		WS:         wsEndpoint(Tokyo, false),
		WSSecure:   wsEndpoint(Tokyo, true),
		GRPC:       grpcEndpoint(Tokyo, false),
		GRPCSecure: grpcEndpoint(Tokyo, true),
	},
	"Testnet": {
		HTTP:       httpEndpoint(Testnet, false),
		HTTPSecure: httpEndpoint(Testnet, true),
		WS:         wsEndpoint(Testnet, false),
		WSSecure:   wsEndpoint(Testnet, true),
		GRPC:       grpcEndpoint(Testnet, false),
		GRPCSecure: grpcEndpoint(Testnet, true),
	},
	"Devnet": {
		HTTP:       httpEndpoint(Devnet, false),
		HTTPSecure: httpEndpoint(Devnet, true),
		WS:         wsEndpoint(Devnet, false),
		WSSecure:   wsEndpoint(Devnet, true),
		GRPC:       grpcEndpoint(Devnet, false),
		GRPCSecure: grpcEndpoint(Devnet, true),
	},
}

var (
	MainnetNYHTTP       = endpointsByRegion["NY"].HTTP
	MainnetNYHTTPSecure = endpointsByRegion["NY"].HTTPSecure
	MainnetNYWS         = endpointsByRegion["NY"].WS
	MainnetNYWSSecure   = endpointsByRegion["NY"].WSSecure
	MainnetNYGRPC       = endpointsByRegion["NY"].GRPC
	MainnetNYGRPCSecure = endpointsByRegion["NY"].GRPCSecure

	MainnetPumpNYHTTP       = endpointsByRegion["NY"].PumpHTTP
	MainnetPumpNYHTTPSecure = endpointsByRegion["NY"].PumpHTTPSecure
	MainnetPumpNYWS         = endpointsByRegion["NY"].PumpWS
	MainnetPumpNYWSSecure   = endpointsByRegion["NY"].PumpWSSecure
	MainnetPumpNYGRPC       = endpointsByRegion["NY"].PumpGRPC
	MainnetPumpNYGRPCSecure = endpointsByRegion["NY"].PumpGRPCSecure

	MainnetUKHTTP       = endpointsByRegion["UK"].HTTP
	MainnetUKHTTPSecure = endpointsByRegion["UK"].HTTPSecure
	MainnetUKWS         = endpointsByRegion["UK"].WS
	MainnetUKWSSecure   = endpointsByRegion["UK"].WSSecure
	MainnetUKGRPC       = endpointsByRegion["UK"].GRPC
	MainnetUKGRPCSecure = endpointsByRegion["UK"].GRPCSecure

	MainnetPumpUKHTTP       = endpointsByRegion["UK"].PumpHTTP
	MainnetPumpUKHTTPSecure = endpointsByRegion["UK"].PumpHTTPSecure
	MainnetPumpUKWS         = endpointsByRegion["UK"].PumpWS
	MainnetPumpUKWSSecure   = endpointsByRegion["UK"].PumpWSSecure
	MainnetPumpUKGRPC       = endpointsByRegion["UK"].PumpGRPC
	MainnetPumpUKGRPCSecure = endpointsByRegion["UK"].PumpGRPCSecure

	MainnetFrankfurtHTTP       = endpointsByRegion["Frankfurt"].HTTP
	MainnetFrankfurtHTTPSecure = endpointsByRegion["Frankfurt"].HTTPSecure
	MainnetFrankfurtWS         = endpointsByRegion["Frankfurt"].WS
	MainnetFrankfurtWSSecure   = endpointsByRegion["Frankfurt"].WSSecure
	MainnetFrankfurtGRPC       = endpointsByRegion["Frankfurt"].GRPC
	MainnetFrankfurtGRPCSecure = endpointsByRegion["Frankfurt"].GRPCSecure

	MainnetLAHTTP       = endpointsByRegion["LosAngeles"].HTTP
	MainnetLAHTTPSecure = endpointsByRegion["LosAngeles"].HTTPSecure
	MainnetLAWS         = endpointsByRegion["LosAngeles"].WS
	MainnetLAWSSecure   = endpointsByRegion["LosAngeles"].WSSecure
	MainnetLAGRPC       = endpointsByRegion["LosAngeles"].GRPC
	MainnetLAGRPCSecure = endpointsByRegion["LosAngeles"].GRPCSecure

	MainnetAmsterdamHTTP       = endpointsByRegion["Amsterdam"].HTTP
	MainnetAmsterdamHTTPSecure = endpointsByRegion["Amsterdam"].HTTPSecure
	MainnetAmsterdamWS         = endpointsByRegion["Amsterdam"].WS
	MainnetAmsterdamWSSecure   = endpointsByRegion["Amsterdam"].WSSecure
	MainnetAmsterdamGRPC       = endpointsByRegion["Amsterdam"].GRPC
	MainnetAmsterdamGRPCSecure = endpointsByRegion["Amsterdam"].GRPCSecure

	MainnetTokyoHTTP       = endpointsByRegion["Tokyo"].HTTP
	MainnetTokyoHTTPSecure = endpointsByRegion["Tokyo"].HTTPSecure
	MainnetTokyoWS         = endpointsByRegion["Tokyo"].WS
	MainnetTokyoWSSecure   = endpointsByRegion["Tokyo"].WSSecure
	MainnetTokyoGRPC       = endpointsByRegion["Tokyo"].GRPC
	MainnetTokyoGRPCSecure = endpointsByRegion["Tokyo"].GRPCSecure

	TestnetHTTP       = endpointsByRegion["Testnet"].HTTP
	TestnetHTTPSecure = endpointsByRegion["Testnet"].HTTPSecure
	TestnetWS         = endpointsByRegion["Testnet"].WS
	TestnetWSSecure   = endpointsByRegion["Testnet"].WSSecure
	TestnetGRPC       = endpointsByRegion["Testnet"].GRPC
	TestnetGRPCSecure = endpointsByRegion["Testnet"].GRPCSecure

	DevnetHTTP       = endpointsByRegion["Devnet"].HTTP
	DevnetHTTPSecure = endpointsByRegion["Devnet"].HTTPSecure
	DevnetWS         = endpointsByRegion["Devnet"].WS
	DevnetWSSecure   = endpointsByRegion["Devnet"].WSSecure
	DevnetGRPC       = endpointsByRegion["Devnet"].GRPC
	DevnetGRPCSecure = endpointsByRegion["Devnet"].GRPCSecure
)

func GetFullServiceGRPCAllEndpoints() []string {
	return []string{
		MainnetNYGRPC,
		MainnetNYGRPCSecure,
		MainnetUKGRPC,
		MainnetUKGRPCSecure,
	}
}

func GetFullServiceHTTPAllEndpoints() []string {
	return []string{
		MainnetNYHTTP,
		MainnetNYHTTPSecure,
		MainnetUKHTTP,
		MainnetUKHTTPSecure,
	}
}

func GetFullServiceWSAllEndpoints() []string {
	return []string{
		MainnetNYWS,
		MainnetNYWSSecure,
		MainnetUKWS,
		MainnetUKWSSecure,
	}
}

func GetSubmitOnlyGRPCAllEndpoints() []string {
	return []string{
		MainnetAmsterdamGRPC,
		MainnetAmsterdamGRPCSecure,
		MainnetLAGRPC,
		MainnetLAGRPCSecure,
		MainnetFrankfurtGRPC,
		MainnetFrankfurtGRPCSecure,
		MainnetTokyoGRPC,
		MainnetTokyoGRPCSecure,
	}
}

func GetSubmitOnlyHTTPAllEndpoints() []string {
	return []string{
		MainnetAmsterdamHTTP,
		MainnetAmsterdamHTTPSecure,
		MainnetLAHTTP,
		MainnetLAHTTPSecure,
		MainnetFrankfurtHTTP,
		MainnetFrankfurtHTTPSecure,
		MainnetTokyoHTTP,
		MainnetTokyoHTTPSecure,
	}
}

func GetSubmitOnlyWSAllEndpoints() []string {
	return []string{
		MainnetAmsterdamWS,
		MainnetAmsterdamWSSecure,
		MainnetLAWS,
		MainnetLAWSSecure,
		MainnetFrankfurtWS,
		MainnetFrankfurtWSSecure,
		MainnetTokyoWS,
		MainnetTokyoWSSecure,
	}
}

func GetPumpGRPCAllEndpoints() []string {
	return []string{
		MainnetPumpNYGRPC,
		MainnetPumpNYGRPCSecure,
		MainnetPumpUKGRPC,
		MainnetPumpUKGRPCSecure,
	}
}

func GetPumpHTTPAllEndpoints() []string {
	return []string{
		MainnetPumpNYHTTP,
		MainnetPumpNYHTTPSecure,
		MainnetPumpUKHTTP,
		MainnetPumpUKHTTPSecure,
	}
}

func GetPumpWSAllEndpoints() []string {
	return []string{
		MainnetPumpNYWS,
		MainnetPumpNYWSSecure,
		MainnetPumpUKWS,
		MainnetPumpUKWSSecure,
	}
}

// Endpoint string generators
func httpEndpoint(e Endpoint, secure bool) string {
	proto := "http"
	if secure {
		proto = "https"
	}
	return fmt.Sprintf("%s://%s", proto, e)
}

func wsEndpoint(e Endpoint, secure bool) string {
	proto := "ws"
	if secure {
		proto = "wss"
	}
	return fmt.Sprintf("%s://%s/ws", proto, e)
}

func grpcEndpoint(e Endpoint, secure bool) string {
	port := "80"
	if secure {
		port = "443"
	}
	return fmt.Sprintf("%s:%s", e, port)
}
