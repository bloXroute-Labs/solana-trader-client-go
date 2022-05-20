package main

import (
	"fmt"
	"github.com/bloXroute-Labs/serum-api/bxserum/provider"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	callOrderbookHTTP()
	callOpenOrdersHTTP()
	callTradesHTTP()
	callTickersHTTP()
	callUnsettledHTTP()
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
