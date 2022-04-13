# Golang-Serum-Client

## Objective
This SDK is designed to make it easy for you to use the [Bloxroute Labs Serum API](https://github.com/bloXroute-Labs/serum-api)
in Go. As we continue to develop the Serum API, we will update this code to have more methods as well. Currently the methods supported are:

#### GRPC/WS:
```azure
GetOrderbook
GetOrderbookStream
```

#### HTTP:
```azure
GetOrderbook
```

## Usage
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

    orderbook, err := g.GetOrderbook(context.Background(), "ETH-USDT")
    if err != nil {
        // ...
    }

    // HTTP
    h := provider.NewHTTPClient()
    orderbook, err := h.GetOrderbook("ETH-USDT") // do not put slashes in the HTTP request, you can do A-B (`ETH-USDT`) or AB (i.e. `ETHUSDT`)
    if err != nil {
        // ...
    }

    // WS
    w, err := provider.NewWSClient()
    if err != nil {
        // ...
    }

    orderbook, err := w.GetOrderbook("ETH-USDT")
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

    g, err := provider.NewGRPCClient() // you can replace this with `NewWSClient()`
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

To run some working code samples, please visit the `cmd` where you can see examples of code for each client.