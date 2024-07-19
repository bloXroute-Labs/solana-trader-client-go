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
		Value: "http://185.209.178.215:8899",
	}
	SolanaWSRPCEndpointFlag = &cli.StringFlag{
		Name:  "solana-ws-endpoint",
		Usage: "WS RPC server endpoint to make blockchain pub/sub queries against",
		Value: "ws://185.209.178.215:6677/ws",
	}
	APIWSEndpoint = &cli.StringFlag{
		Name:  "solana-trader-ws-endpoint",
		Usage: "Solana Trader API API websocket connection endpoint",
		Value: "wss://ny.solana.dex.blxrbdn.com/ws",
	}
)
