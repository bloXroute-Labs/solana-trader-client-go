# Serum Golang Client

## Objective
This SDK is designed to make it easy for you to use the [bloXroute Labs Serum API](https://github.com/bloXroute-Labs/serum-api)
in Go. As we continue to develop the Serum API, we will update the methods here as well. Currently, the methods supported are:

#### GRPC/WS:
```
GetOrderbook
GetOrderbookStream
```

#### HTTP:
```
GetOrderbook
```
Methods that end in `Stream` continuously stream responses through a channel, while other methods return a one and done response.

Furthermore, the HTTP client only supports unary/non-streaming methods.

## Installation
```
go get github.com/bloXroute-Labs/serum-api/bxserum/provider
go get github.com/bloXroute-Labs/serum-api/proto
```

## Usage
Note: Markets can be provided in different formats:
1. `A/B` (only for GRPC/WS clients) --> `ETH/USDT`
2. `A:B` --> `ETH:USDT`
3. `A-B` --> `ETH-USDT`
4. `AB` --> `ETHUSDT`

#### Unary Response:
```go
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

    orderbook, err := g.GetOrderbook(context.Background(), "ETH/USDT")
    if err != nil {
        // ...
    }


    // HTTP
    h := provider.NewHTTPClient()
    orderbook, err := h.GetOrderbook("ETH-USDT") // do not use forward slashes for the HTTP market parameter
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
}

```
#### Stream response (only in GRPC/WS):
```go
import (
    "github.com/bloXroute-Labs/serum-api/bxserum/provider"
    pb "github.com/bloXroute-Labs/serum-api/proto"
    "context"
)

func main() {
    orderbookChan := make(chan *pb.GetOrderbookStreamResponse)
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    g, err := provider.NewGRPCClient() // replace this with `NewWSClient()` to use WebSockets
    if err != nil {
        // ...
    }

    err = g.GetOrderbookStream(ctx, "SOL/USDT", orderbookChan)
    if err != nil {
        // ...
    }
    for {
        orderbook := <-orderbookChan
    }
}
```

To run some working code samples, please visit the `examples` directory where you can see examples of code for each client.