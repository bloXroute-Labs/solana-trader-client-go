package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/bloXroute-Labs/serum-api/bxserum/provider"
	api "github.com/bloXroute-Labs/serum-api/proto"
	log "github.com/sirupsen/logrus"
)

func main() {
	callOrderbookHTTP()
	callOpenOrdersHTTP()
	callTradesHTTP()
	callTickersHTTP()
	callUnsettledHTTP()

	ownerAddr, ok := os.LookupEnv("PUBLIC_KEY")
	if !ok {
		log.Infof("PUBLIC_KEY environment variable not set")
		log.Infof("Skipping Place and Cancel Order examples")
		return
	}

	ooAddr, _ := os.LookupEnv("OPEN_ORDERS")
	if !ok {
		log.Infof("OPEN_ORDERS environment variable not set")
		log.Infof("Skipping Place and Cancel Order examples")
		return
	}

	clientID := callPlaceOrderGRPC(ownerAddr, ooAddr)
	callCancelOrderByClientID(ownerAddr, ooAddr, clientID)
}

func callOrderbookHTTP() {
	h, _ := provider.NewHTTPClient()

	// Unary response
	orderbook, err := h.GetOrderbook("ETH-USDT", 0)
	if err != nil {
		log.Errorf("error with GetOrderbook request for ETH-USDT: %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()

	orderbook, err = h.GetOrderbook("SOLUSDT", 2)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOLUSDT: %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()

	orderbook, err = h.GetOrderbook("SOL:USDC", 3)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL:USDC: %v", err)
	} else {
		fmt.Println(orderbook)
	}
}

func callOpenOrdersHTTP() {
	client := &http.Client{Timeout: time.Second * 60}
	opts, err := provider.DefaultRPCOpts(provider.MainnetSerumAPIHTTP)
	h := provider.NewHTTPClientWithOpts(client, opts)

	orders, err := h.GetOpenOrders("SOLUSDT", "HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc")
	if err != nil {
		log.Errorf("error with GetOrders request for SOLUSDT: %v", err)
	} else {
		fmt.Println(orders)
	}

	fmt.Println()
}

func callUnsettledHTTP() {
	client := &http.Client{Timeout: time.Second * 60}
	opts, err := provider.DefaultRPCOpts(provider.MainnetSerumAPIHTTP)
	h := provider.NewHTTPClientWithOpts(client, opts)

	response, err := h.GetUnsettled("SOLUSDT", "HxFLKUAmAMLz1jtT3hbvCMELwH5H9tpM2QugP8sKyfhc")
	if err != nil {
		log.Errorf("error with GetOrders request for SOLUSDT: %v", err)
	} else {
		fmt.Println(response)
	}

	fmt.Println()
}

func callTradesHTTP() {
	h, _ := provider.NewHTTPClient()

	trades, err := h.GetTrades("SOLUSDT", 5)
	if err != nil {
		log.Errorf("error with GetTrades request for SOLUSDT: %v", err)
	} else {
		fmt.Println(trades)
	}

	fmt.Println()
}

func callTickersHTTP() {
	h, _ := provider.NewHTTPClient()

	tickers, err := h.GetTickers("SOLUSDT")
	if err != nil {
		log.Errorf("error with GetTickers request for SOLUSDT: %v", err)
	} else {
		fmt.Println(tickers)
	}

	fmt.Println()
}

const (
	// SOL/USDC market
	marketAddr = "9wFFyRfZBsuAha4YcuxcXLKwMxJR43S7fPfQLusDBzvT"

	orderSide   = api.Side_S_ASK
	orderType   = api.OrderType_OT_LIMIT
	orderPrice  = float64(170200)
	orderAmount = float64(0.1)
)

func callPlaceOrderGRPC(ownerAddr, ooAddr string) uint64 {
	client := &http.Client{Timeout: time.Second * 30}
	opts, err := provider.DefaultRPCOpts(provider.MainnetSerumAPIHTTP)
	h := provider.NewHTTPClientWithOpts(client, opts)

	// generate a random clientId for this order
	rand.Seed(time.Now().UnixNano())
	clientID := rand.Uint64()

	orderOpts := provider.PostOrderOpts{
		ClientOrderID:     clientID,
		OpenOrdersAddress: ooAddr,
	}

	response, err := h.SubmitOrder(ownerAddr, ownerAddr, marketAddr,
		orderSide, []api.OrderType{orderType}, orderAmount,
		orderPrice, orderOpts)
	if err != nil {
		log.Fatalf("failed to submit order (%w)", err)
	}

	fmt.Printf("placed order %v with clientID %x\n", response, clientID)

	return clientID
}

func callCancelOrderByClientID(ownerAddr, ooAddr string, clientID uint64) {
	client := &http.Client{Timeout: time.Second * 30}
	opts, err := provider.DefaultRPCOpts(provider.MainnetSerumAPIHTTP)
	h := provider.NewHTTPClientWithOpts(client, opts)

	_, err = h.SubmitCancelOrderByClientID(clientID, ownerAddr,
		marketAddr, ooAddr)
	if err != nil {
		log.Fatalf("failed to cancel order by client ID (%w)", err)
	}

	fmt.Printf("canceled order for clientID %x\n", clientID)
}
