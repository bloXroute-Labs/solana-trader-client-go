# Solana Trader Golang Client

## Objective
This SDK is designed to make it easy for you to use the bloXroute Labs Solana Trader API
in Go. 

## Installation
```
go get github.com/bloXroute-Labs/solana-trader-client-go
```

## Usage

This library supports HTTP, websockets, and GRPC interfaces. You must use websockets or GRPC for any streaming methods, 
but any simple request/response calls are universally supported.

For any methods involving transaction creation you will need to provide your Solana private key. You can provide this 
via the environment variable `PRIVATE_KEY`, or specify it via the provider configuration if you want to load it with
some other mechanism. See samples for more information. As a general note on this: methods named `Post*` (e.g. 
`PostOrder`) typically do not sign/submit the transaction, only return the raw unsigned transaction. This isn't 
very useful to most users (unless you want to write a signer in a different language), and you'll typically want the 
similarly named `Submit*` methods (e.g. `SubmitOrder`). These methods generate, sign, and submit the
transaction all at once.

You will also need your bloXroute authorization header to use these endpoints. By default, this is loaded from the 
`AUTH_HEADER` environment variable.

## Quickstart

### Request sample:

```go
package main

import (
	"context"
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
)

func main() {
	// GRPC
	g, err := provider.NewGRPCClient()
	if err != nil {
		panic(err)
	}

	// Get Raydium pools
	pools, err := g.GetRaydiumPools(context.Background(), &pb.GetRaydiumPoolsRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Raydium pools:", pools)

	// HTTP
	h := provider.NewHTTPClient()
	
	// Get token prices from Jupiter
	prices, err := h.GetJupiterPrices(context.Background(), &pb.GetJupiterPricesRequest{
		Tokens: []string{"So11111111111111111111111111111111111111112", "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Jupiter prices:", prices)
	
	// WS
	w, err := provider.NewWSClient()
	if err != nil {
		panic(err)
	}
	
	// Get swap quotes from Raydium
	quotes, err := w.GetRaydiumQuotes(context.Background(), &pb.GetRaydiumQuotesRequest{
		InToken:  "So11111111111111111111111111111111111111112", // SOL
		OutToken: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", // USDC
		InAmount: 0.01,
		Slippage: 0.5,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Raydium quotes:", quotes)
}

```
#### Stream (only in GRPC/WS):

```go
package main

import (
	"fmt"
	"github.com/bloXroute-Labs/solana-trader-client-go/provider"
	pb "github.com/bloXroute-Labs/solana-trader-proto/api"
	"context"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, err := provider.NewGRPCClient() // replace this with `NewWSClient()` to use WebSockets
	if err != nil {
		panic(err)
	}

	// Stream prices updates for SOL token
	stream, err := g.GetPricesStream(ctx, []pb.Project{pb.Project_P_RAYDIUM}, 
		[]string{"So11111111111111111111111111111111111111112"})
	if err != nil {
		panic(err)
	}
	
	// Wrap result in channel for ease of use
	pricesCh := make(chan *pb.GetPricesStreamResponse)
	stream.Into(pricesCh)
	for i := 0; i < 3; i++ {
		prices := <-pricesCh
		fmt.Println("Price update:", prices)
	}
	
	// Example of pool reserves stream
	poolsStream, err := g.GetPoolReservesStream(ctx, &pb.GetPoolReservesStreamRequest{
		Projects: []pb.Project{pb.Project_P_RAYDIUM},
		Pools: []string{
			"58oQChx4yWmvKdwLLZzBi4ChoCc2fqCUWBkwMihLYQo2", // SOL-USDC pool
		},
	})
	if err != nil {
		panic(err)
	}
	
	poolsCh := make(chan *pb.GetPoolReservesStreamResponse)
	poolsStream.Into(poolsCh)
	for i := 0; i < 3; i++ {
		poolUpdate := <-poolsCh
		fmt.Println("Pool reserves update:", poolUpdate)
	}
}
```

More code samples are provided in the `examples/` directory.

**A quick note on market names:**
You can use a couple of different formats, with restrictions: 
1. `A/B` (only for GRPC/WS clients) --> `ETH/USDT`
2. `A:B` --> `ETH:USDT`
3. `A-B` --> `ETH-USDT`
4. `AB` --> `ETHUSDT`


## Development

Unit tests:

```
$ make unit
```

Integration tests per provider:
```
$ make grpc-examples

$ make http-examples

$ make ws-examples
```