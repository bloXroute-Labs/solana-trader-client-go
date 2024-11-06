package utils

import "github.com/urfave/cli/v2"

var (
	OutputFileFlag = &cli.StringFlag{
		Name:     "output",
		Usage:    "file to output CSV results to",
		Required: true,
	}
	SolanaHTTPRPCEndpointFlag = &cli.StringFlag{
		Name:  "solana-http-endpoint",
		Usage: "HTTP RPC server endpoint to make blockchain queries against",
	}
	SolanaWSRPCEndpointFlag = &cli.StringFlag{
		Name:  "solana-ws-endpoint",
		Usage: "WS RPC server endpoint to make blockchain pub/sub queries against",
	}
	APIWSEndpoint = &cli.StringFlag{
		Name:  "solana-trader-ws-endpoint",
		Usage: "Solana Trader API API websocket connection endpoint",
		Value: "wss://virginia.solana.dex.blxrbdn.com/ws",
	}
)
