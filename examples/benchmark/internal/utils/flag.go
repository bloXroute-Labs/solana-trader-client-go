package utils

import "github.com/urfave/cli/v2"

var (
	SolanaHTTPRPCEndpointFlag = &cli.StringFlag{
		Name:  "solana-http-endpoint",
		Usage: "HTTP RPC server endpoint to make blockchain queries against",
		Value: "http://34.203.186.197:8899",
	}
	SolanaWSRPCEndpointFlag = &cli.StringFlag{
		Name:  "solana-ws-endpoint",
		Usage: "WS RPC server endpoint to make blockchain pub/sub queries against",
		Value: "ws://34.203.186.197:8900/ws",
	}
)
