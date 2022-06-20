# Serum Golang Client

## Objective
This SDK is designed to make it easy for you to use the [bloXroute Labs Serum API](https://github.com/bloXroute-Labs/serum-api)
in Go. 

## Installation
```
go get github.com/bloXroute-Labs/serum-api/bxserum/provider
go get github.com/bloXroute-Labs/serum-api/proto
```

## Usage

This library supports HTTP, websockets, and GRPC interfaces. You must use websockets or GRPC for any streaming methods, 
but any simple request/response calls are universally supported.

#### Request:
```go
package main

import (
    "github.com/bloXroute-Labs/serum-api/bxserum/provider"
    pb "github.com/bloXroute-Labs/serum-api/proto"
    "context"
)

func main() {
    // GPRC
    g, err := provider.NewGRPCClient()
    if err != nil {
        // ...
    }

    orderbook, err := g.GetOrderbook(context.Background(), "ETH/USDT", 5) // in this case limit to 5 bids and asks. 0 for no limit
    if err != nil {
        // ...
    }

    trades, err := g.GetTrades(context.Background(), "ETH/USDT", 5) // in this case limit to 5 trades. 0 for no limit
    if err != nil {
        // ...
    }

    tickers, err := g.GetTickers(context.Background(), "ETH/USDT") 
    if err != nil {
        // ...
    }
    
    openOrders, err := g.GetOpenOrders(context.Background(), "ETH/USDT", "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv") 
    if err != nil {
        // ...
    }
	
    unsettledFunds, err := g.GetUnsettled(context.Background(), "ETH/USDT", "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv")
    if err != nil {
        // ...
    }
	
    supportedMarkets, err := g.GetMarkets(context.Background()) 
    if err != nil {
        // ...
    }
	
    newOrderUnsignedTransaction, err := g.SubmitOrder(
        context.Background(), 
        "BraJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgC", //owner solana wallet address
        "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv", // SPL token wallet address 
        "SOL/USDC", // market 
        pb.Side_S_BID, // trade side (Bid/Ask)
        []pb.OrderType{pb.OrderType_OT_LIMIT}, // OrderType
        20, // order size
        float64(0.124), // order price 
        provider.PostOrderOpts{
            ClientOrderID: 5000, // Client controlled OrderID
        })
	// The SubmitOrder relies on the PRIVATE_KEY env variable holding your wallet's private key, to sign the transaction
    
    if err != nil {
        // ...
    }

    settleTransaction, err := g.SettleFunds(
        context.Background(),
        "BraJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgC", //owner solana wallet address
        "SOL/USDC", // market 
        "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv", // base SPL token wallet address 
        "CbafjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjFun", // quote SPL token wallet address 
        "neePfCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSJOKE", // open orders account address for this market
        )
    // Settle relies on the PRIVATE_KEY env variable holding your wallet's private key, to sign the transaction

	if err != nil {
        // ...
    }


    // HTTP
    h := provider.NewHTTPClient()
    orderbook, err := h.GetOrderbook("ETH-USDT") // do not use forward slashes for the HTTP market parameter
    if err != nil {
        // ...
    }

    trades, err := h.GetTrades(context.Background(), "ETH/USDT", 5) // in this case limit to 5 trades. 0 for no limit
    if err != nil {
        // ...
    }

    tickers, err := h.GetTickers(context.Background(), "ETH/USDT")
    if err != nil {
        // ...
    }

    openOrders, err := h.GetOpenOrders(context.Background(), "ETH/USDT", "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv")
    if err != nil {
        // ...
    }
	
    unsettledFunds, err := h.GetUnsettled(context.Background(), "ETH/USDT", "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv")
    if err != nil {
        // ...
    }

    supportedMarkets, err := h.GetMarkets(context.Background())
    if err != nil {
        // ...
    }

    newOrderUnsignedTransaction, err := h.SubmitOrder(
        context.Background(),
        "BraJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgC", //owner solana wallet address
        "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv", // SPL token wallet address 
        "SOL/USDC", // market 
        pb.Side_S_BID, // trade side (Bid/Ask)
        []pb.OrderType{pb.OrderType_OT_LIMIT}, // OrderType
        20, // order size
        float64(0.124), // order price 
        provider.PostOrderOpts{
            ClientOrderID: 5000, // Client controlled OrderID
        }) 
    // The SubmitOrder relies on the PRIVATE_KEY env variable holding your wallet's private key, to sign the transaction
    if err != nil {
        // ...
    }

    settleTransaction, err := h.SettleFunds(
        context.Background(),
        "BraJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgC", //owner solana wallet address
        "SOL/USDC", // market 
        "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv", // base SPL token wallet address 
        "CbafjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjFun", // quote SPL token wallet address 
        "neePfCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSJOKE", // open orders account address for this market
    )
    // Settle relies on the PRIVATE_KEY env variable holding your wallet's private key, to sign the transaction

    if err != nil {
    // ...
    }
    
    
    // WS
    w, err := provider.NewWSClient()
    if err != nil {
        // ...
    }

    orderbook, err := w.GetOrderbook("ETH/USDT")
    if err != nil {
        // ...
    }

    trades, err := w.GetTrades(context.Background(), "ETH/USDT", 5) // in this case limit to 5 trades. 0 for no limit
    if err != nil {
        // ...
    }

    tickers, err := w.GetTickers(context.Background(), "ETH/USDT")
    if err != nil {
        // ...
    }

    openOrders, err := w.GetOpenOrders(context.Background(), "ETH/USDT", "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv")
    if err != nil {
        // ...
    }

    unsettledFunds, err := w.GetUnsettled(context.Background(), "ETH/USDT", "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv")
    if err != nil {
        // ...
    }

    supportedMarkets, err := w.GetMarkets(context.Background())
    if err != nil {
        // ...
    }

    newOrderUnsignedTransaction, err := w.SubmitOrder(
        context.Background(),
        "BraJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgC", //owner solana wallet address
        "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv", // SPL token wallet address 
        "SOL/USDC", // market 
        pb.Side_S_BID, // trade side (Bid/Ask)
        []pb.OrderType{pb.OrderType_OT_LIMIT}, // OrderType
        20, // order size
        float64(0.124), // order price 
        provider.PostOrderOpts{
            ClientOrderID: 5000, // Client controlled OrderID
        })
        // The SubmitOrder relies on the PRIVATE_KEY env variable holding your wallet's private key, to sign the transaction

    if err != nil {
        // ...
    }

    settleTransaction, err := w.SettleFunds(
        context.Background(),
        "BraJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgC", //owner solana wallet address
        "SOL/USDC", // market 
        "4raJjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjcgv", // base SPL token wallet address 
        "CbafjCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSjFun", // quote SPL token wallet address 
        "neePfCwLLqw8TciQXYruDEF4YhDkGwoEnwnAdwJSJOKE", // open orders account address for this market
        )
    // Settle relies on the PRIVATE_KEY env variable holding your wallet's private key, to sign the transaction

    if err != nil {
        // ...
    }
}

```
#### Stream (only in GRPC/WS):
```go
import (
    "github.com/bloXroute-Labs/serum-api/bxserum/provider"
    pb "github.com/bloXroute-Labs/serum-api/proto"
    "context"
)

func main() {
    
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    g, err := provider.NewGRPCClient() // replace this with `NewWSClient()` to use WebSockets
    if err != nil {
        // ...
    }
	
    orderbookChan := make(chan *pb.GetOrderbookStreamResponse)
    err = g.GetOrderbookStream(ctx, "SOL/USDT", 5,  orderbookChan)
    if err != nil {
        // ...
    }
    for {
        orderbook := <-orderbookChan
    }

    tradesChan := make(chan *pb.GetTradesStreamResponse)
    err = g.GetTradesStream(ctx, "SOL/USDT", 5, tradesChan)
    if err != nil {
        // ...
    }
    for {
        trades := <-tradesChan
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
