package provider

import (
	"fmt"
)

// Warning to display when creating TLS client
const WarningTLSSlowDown = "Performance Notice: Secure (TLS) endpoints may introduce latency due to handshake overhead. For optimal trading speed, consider using non-secure endpoints when appropriate."

type Protocol string

const (
	HTTP Protocol = "http"
	WS   Protocol = "ws"
	GRPC Protocol = "grpc"
)

type ServiceType string

const (
	FullService ServiceType = "full"
	SubmitOnly  ServiceType = "submit"
	Pump        ServiceType = "pump"
	Development ServiceType = "dev"
)

type EndpointConfig struct {
	Host        string
	ServiceType ServiceType
	PumpHost    string
}

type Region struct {
	HTTP, HTTPSecure         string
	WS, WSSecure             string
	GRPC, GRPCSecure         string
	PumpHTTP, PumpHTTPSecure string
	PumpWS, PumpWSSecure     string
	PumpGRPC, PumpGRPCSecure string
}

// Regional endpoint configurations
var endpointConfigs = map[string]EndpointConfig{
	"NY": {
		Host:        "ny.solana.dex.blxrbdn.com",
		ServiceType: FullService,
		PumpHost:    "pump-ny.solana.dex.blxrbdn.com",
	},
	"UK": {
		Host:        "uk.solana.dex.blxrbdn.com",
		ServiceType: FullService,
		PumpHost:    "pump-uk.solana.dex.blxrbdn.com",
	},
	"Frankfurt": {
		Host:        "germany.solana.dex.blxrbdn.com",
		ServiceType: SubmitOnly,
	},
	"LosAngeles": {
		Host:        "la.solana.dex.blxrbdn.com",
		ServiceType: SubmitOnly,
	},
	"Amsterdam": {
		Host:        "amsterdam.solana.dex.blxrbdn.com",
		ServiceType: SubmitOnly,
	},
	"Tokyo": {
		Host:        "tokyo.solana.dex.blxrbdn.com",
		ServiceType: SubmitOnly,
	},
	"Testnet": {
		Host:        "solana.dex.bxrtest.com",
		ServiceType: Development,
	},
	"Devnet": {
		Host:        "solana-trader-api-nlb-6b0f765f2fc759e1.elb.us-east-1.amazonaws.com",
		ServiceType: Development,
	},
}

// Local endpoints
const (
	LocalHTTP string = "http://localhost:9000"
	LocalWS   string = "ws://localhost:9000/ws"
	LocalGRPC string = "localhost:9000"
)

var endpointsByRegion = generateRegions()

var (
	// NY Endpoints
	MainnetNYHTTP, MainnetNYHTTPSecure         = getEndpoints("NY", HTTP)
	MainnetNYWS, MainnetNYWSSecure             = getEndpoints("NY", WS)
	MainnetNYGRPC, MainnetNYGRPCSecure         = getEndpoints("NY", GRPC)
	MainnetPumpNYHTTP, MainnetPumpNYHTTPSecure = getPumpEndpoints("NY", HTTP)
	MainnetPumpNYWS, MainnetPumpNYWSSecure     = getPumpEndpoints("NY", WS)
	MainnetPumpNYGRPC, MainnetPumpNYGRPCSecure = getPumpEndpoints("NY", GRPC)

	// UK Endpoints
	MainnetUKHTTP, MainnetUKHTTPSecure         = getEndpoints("UK", HTTP)
	MainnetUKWS, MainnetUKWSSecure             = getEndpoints("UK", WS)
	MainnetUKGRPC, MainnetUKGRPCSecure         = getEndpoints("UK", GRPC)
	MainnetPumpUKHTTP, MainnetPumpUKHTTPSecure = getPumpEndpoints("UK", HTTP)
	MainnetPumpUKWS, MainnetPumpUKWSSecure     = getPumpEndpoints("UK", WS)
	MainnetPumpUKGRPC, MainnetPumpUKGRPCSecure = getPumpEndpoints("UK", GRPC)

	// Other regions (Submit-only)
	MainnetFrankfurtHTTP, MainnetFrankfurtHTTPSecure = getEndpoints("Frankfurt", HTTP)
	MainnetFrankfurtWS, MainnetFrankfurtWSSecure     = getEndpoints("Frankfurt", WS)
	MainnetFrankfurtGRPC, MainnetFrankfurtGRPCSecure = getEndpoints("Frankfurt", GRPC)
	MainnetLAHTTP, MainnetLAHTTPSecure               = getEndpoints("LosAngeles", HTTP)
	MainnetLAWS, MainnetLAWSSecure                   = getEndpoints("LosAngeles", WS)
	MainnetLAGRPC, MainnetLAGRPCSecure               = getEndpoints("LosAngeles", GRPC)
	MainnetAmsterdamHTTP, MainnetAmsterdamHTTPSecure = getEndpoints("Amsterdam", HTTP)
	MainnetAmsterdamWS, MainnetAmsterdamWSSecure     = getEndpoints("Amsterdam", WS)
	MainnetAmsterdamGRPC, MainnetAmsterdamGRPCSecure = getEndpoints("Amsterdam", GRPC)
	MainnetTokyoHTTP, MainnetTokyoHTTPSecure         = getEndpoints("Tokyo", HTTP)
	MainnetTokyoWS, MainnetTokyoWSSecure             = getEndpoints("Tokyo", WS)
	MainnetTokyoGRPC, MainnetTokyoGRPCSecure         = getEndpoints("Tokyo", GRPC)

	// Development endpoints
	TestnetHTTP, TestnetHTTPSecure = getEndpoints("Testnet", HTTP)
	TestnetWS, TestnetWSSecure     = getEndpoints("Testnet", WS)
	TestnetGRPC, TestnetGRPCSecure = getEndpoints("Testnet", GRPC)
	DevnetHTTP, DevnetHTTPSecure   = getEndpoints("Devnet", HTTP)
	DevnetWS, DevnetWSSecure       = getEndpoints("Devnet", WS)
	DevnetGRPC, DevnetGRPCSecure   = getEndpoints("Devnet", GRPC)
)

func generateRegions() map[string]Region {
	regions := make(map[string]Region)

	for name, config := range endpointConfigs {
		region := Region{
			HTTP:       buildEndpoint(config.Host, HTTP, false),
			HTTPSecure: buildEndpoint(config.Host, HTTP, true),
			WS:         buildEndpoint(config.Host, WS, false),
			WSSecure:   buildEndpoint(config.Host, WS, true),
			GRPC:       buildEndpoint(config.Host, GRPC, false),
			GRPCSecure: buildEndpoint(config.Host, GRPC, true),
		}

		// Add pump endpoints if available
		if config.PumpHost != "" {
			region.PumpHTTP = buildEndpoint(config.PumpHost, HTTP, false)
			region.PumpHTTPSecure = buildEndpoint(config.PumpHost, HTTP, true)
			region.PumpWS = buildEndpoint(config.PumpHost, WS, false)
			region.PumpWSSecure = buildEndpoint(config.PumpHost, WS, true)
			region.PumpGRPC = buildEndpoint(config.PumpHost, GRPC, false)
			region.PumpGRPCSecure = buildEndpoint(config.PumpHost, GRPC, true)
		}

		regions[name] = region
	}

	return regions
}

func buildEndpoint(host string, protocol Protocol, secure bool) string {
	switch protocol {
	case HTTP:
		scheme := "http"
		if secure {
			scheme = "https"
		}
		return fmt.Sprintf("%s://%s", scheme, host)
	case WS:
		scheme := "ws"
		if secure {
			scheme = "wss"
		}
		return fmt.Sprintf("%s://%s/ws", scheme, host)
	case GRPC:
		port := "80"
		if secure {
			port = "443"
		}
		return fmt.Sprintf("%s:%s", host, port)
	default:
		return ""
	}
}

func getEndpoints(region string, protocol Protocol) (string, string) {
	config := endpointConfigs[region]
	return buildEndpoint(config.Host, protocol, false), buildEndpoint(config.Host, protocol, true)
}

func getPumpEndpoints(region string, protocol Protocol) (string, string) {
	config := endpointConfigs[region]
	if config.PumpHost == "" {
		return "", ""
	}
	return buildEndpoint(config.PumpHost, protocol, false), buildEndpoint(config.PumpHost, protocol, true)
}

func GetEndpointsByService(serviceType ServiceType, protocol Protocol) []string {
	var endpoints []string

	for name, config := range endpointConfigs {
		if config.ServiceType == serviceType {
			region := endpointsByRegion[name]
			switch protocol {
			case HTTP:
				endpoints = append(endpoints, region.HTTP, region.HTTPSecure)
			case WS:
				endpoints = append(endpoints, region.WS, region.WSSecure)
			case GRPC:
				endpoints = append(endpoints, region.GRPC, region.GRPCSecure)
			}
		}
	}

	return endpoints
}

func GetPumpEndpoints(protocol Protocol) []string {
	var endpoints []string

	for name, config := range endpointConfigs {
		if config.PumpHost != "" {
			region := endpointsByRegion[name]
			switch protocol {
			case HTTP:
				endpoints = append(endpoints, region.PumpHTTP, region.PumpHTTPSecure)
			case WS:
				endpoints = append(endpoints, region.PumpWS, region.PumpWSSecure)
			case GRPC:
				endpoints = append(endpoints, region.PumpGRPC, region.PumpGRPCSecure)
			}
		}
	}

	return endpoints
}

func GetFullServiceGRPCAllEndpoints() []string {
	return GetEndpointsByService(FullService, GRPC)
}

func GetFullServiceHTTPAllEndpoints() []string {
	return GetEndpointsByService(FullService, HTTP)
}

func GetFullServiceWSAllEndpoints() []string {
	return GetEndpointsByService(FullService, WS)
}

func GetSubmitOnlyGRPCAllEndpoints() []string {
	return GetEndpointsByService(SubmitOnly, GRPC)
}

func GetSubmitOnlyHTTPAllEndpoints() []string {
	return GetEndpointsByService(SubmitOnly, HTTP)
}

func GetSubmitOnlyWSAllEndpoints() []string {
	return GetEndpointsByService(SubmitOnly, WS)
}

func GetPumpGRPCAllEndpoints() []string {
	return GetPumpEndpoints(GRPC)
}

func GetPumpHTTPAllEndpoints() []string {
	return GetPumpEndpoints(HTTP)
}

func GetPumpWSAllEndpoints() []string {
	return GetPumpEndpoints(WS)
}
