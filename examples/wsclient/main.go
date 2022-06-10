package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/bloXroute-Labs/serum-api/bxserum/provider"
	api "github.com/bloXroute-Labs/serum-api/proto"
	pb "github.com/bloXroute-Labs/serum-api/proto"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

func main() {
	callOrderbookWS()
	callOpenOrdersWS()
	callTickersWS()
	callUnsettledWS()

	callOrderbookWSStream()
	callTradesWSStream()

	ownerAddr, ok := os.LookupEnv("PUBLIC_KEY")
	if !ok {
		log.Infof("PUBLIC_KEY environment variable not set")
		log.Infof("Skipping Place and Cancel Order examples")
		return
	}

	ooAddr, ok := os.LookupEnv("OPEN_ORDERS")
	if !ok {
		log.Infof("OPEN_ORDERS environment variable not set")
		log.Infof("Skipping Place and Cancel Order examples")
		return
	}

	clientOrderID := callPlaceOrderWS(ownerAddr, ooAddr)
	callCancelByClientOrderIDWS(ownerAddr, ooAddr, clientOrderID)
}

// Unary response
func callOrderbookWS() {
	fmt.Println("fetching orderbooks...")

	w, err := provider.NewWSClient()
	if err != nil {
		log.Fatalf("error dialing WS client: %v", err)
		return
	}
	defer w.Close()

	orderbook, err := w.GetOrderbook("ETH-USDT", 0)
	if err != nil {
		log.Errorf("error with GetOrderbook request for ETH-USDT: %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()

	orderbook, err = w.GetOrderbook("SOLUSDT", 2)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL-USDT: %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()

	orderbook, err = w.GetOrderbook("SOL:USDC", 3)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL:USDC: %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()
}

func callOpenOrdersWS() {
	fmt.Println("fetching open orders...")

	w, err := provider.NewWSClient()
	if err != nil {
		log.Fatalf("error dialing WS client: %v", err)
		return
	}
	defer w.Close()

	orders, err := w.GetOpenOrders("SOLUSDC", "AFT8VayE7qr8MoQsW3wHsDS83HhEvhGWdbNSHRKeUDfQ")
	if err != nil {
		log.Errorf("error with GetOrders request for SOL-USDT: %v", err)
	} else {
		fmt.Println(orders)
	}

	fmt.Println()
}

func callUnsettledWS() {
	fmt.Println("fetching unsettled...")

	w, err := provider.NewWSClient()
	if err != nil {
		log.Fatalf("error dialing WS client: %v", err)
		return
	}
	defer w.Close()

	response, err := w.GetUnsettled("SOLUSDC", "AFT8VayE7qr8MoQsW3wHsDS83HhEvhGWdbNSHRKeUDfQ")
	if err != nil {
		log.Errorf("error with GetOrders request for SOL-USDT: %v", err)
	} else {
		fmt.Println(response)
	}

	fmt.Println()
}

func callTickersWS() {
	fmt.Println("fetching tickers...")

	w, err := provider.NewWSClient()
	if err != nil {
		log.Fatalf("error dialing WS client: %v", err)
		return
	}
	defer w.Close()

	tickers, err := w.GetTickers("SOLUSDC")
	if err != nil {
		log.Errorf("error with GetTickers request for SOL-USDT: %v", err)
	} else {
		fmt.Println(tickers)
	}

	fmt.Println()
}

// Stream response
func callOrderbookWSStream() {
	fmt.Println("starting orderbook stream")

	w, err := provider.NewWSClient()
	if err != nil {
		log.Fatalf("error dialing WS client: %v", err)
		return
	}
	defer w.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	orderbookChan := make(chan *pb.GetOrderbooksStreamResponse)

	err = w.GetOrderbookStream(ctx, "SOL/USDC", 3, orderbookChan)
	if err != nil {
		log.Errorf("error with GetOrderbookStream request for SOL/USDC: %v", err)
	} else {
		for i := 1; i <= 5; i++ {
			<-orderbookChan
			fmt.Printf("response %v received\n", i)
		}
	}
}

func callTradesWSStream() {
	fmt.Println("starting trades stream")

	w, err := provider.NewWSClient()
	if err != nil {
		log.Fatalf("error dialing WS client: %v", err)
		return
	}
	defer w.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tradesChan := make(chan *pb.GetTradesStreamResponse)

	err = w.GetTradesStream(ctx, "SOL/USDC", 3, tradesChan)
	if err != nil {
		log.Errorf("error with GetTradesStream request for SOL/USDC: %v", err)
	} else {
		for i := 1; i <= 5; i++ {
			<-tradesChan
			fmt.Printf("response %v received\n", i)
		}
	}
}

const (
	// SOL/USDC market
	marketAddr = "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT"

	orderSide   = api.Side_S_ASK
	orderType   = api.OrderType_OT_LIMIT
	orderPrice  = float64(170200)
	orderAmount = float64(0.1)
)

func callPlaceOrderWS(ownerAddr, ooAddr string) uint64 {
	fmt.Println("trying to place an order")

	w, err := provider.NewWSClient()
	if err != nil {
		log.Fatalf("error dialing WS client: %v", err)
	}
	defer w.Close()

	// generate a random clientOrderId for this order
	rand.Seed(time.Now().UnixNano())
	clientOrderID := rand.Uint64()

	opts := provider.PostOrderOpts{
		ClientOrderID:     clientOrderID,
		OpenOrdersAddress: ooAddr,
	}

	sig, err := w.SubmitOrder(ownerAddr, ownerAddr, marketAddr,
		orderSide, []api.OrderType{orderType}, orderAmount,
		orderPrice, opts)
	if err != nil {
		log.Fatalf("failed to submit order (%v)", err)
	}

	fmt.Printf("placed order %v with clientOrderID %v\n", sig, clientOrderID)

	return clientOrderID
}

func callCancelByClientOrderIDWS(ownerAddr, ooAddr string, clientOrderID uint64) {
	fmt.Println("trying to cancel order")

	w, err := provider.NewWSClient()
	if err != nil {
		log.Fatalf("error dialing WS client: %v", err)
		return
	}
	defer w.Close()

	_, err = w.SubmitCancelByClientOrderID(clientOrderID, ownerAddr,
		marketAddr, ooAddr, true)
	if err != nil {
		log.Fatalf("failed to cancel order by client ID (%v)", err)
	}

	fmt.Printf("canceled order for clientOrderID %v\n", clientOrderID)
}
