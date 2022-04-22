package main

import (
	"fmt"
	"github.com/bloXroute-Labs/serum-api/bxserum/provider"
	log "github.com/sirupsen/logrus"
)

func main() {
	callHTTP()
}

func callHTTP() {
	h := provider.NewHTTPClient()

	// Unary response
	orderbook, err := h.GetOrderbook("ETH-USDT", 0)
	if err != nil {
		log.Errorf("error with GetOrderbook request for ETH-USDT - %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()

	orderbook, err = h.GetOrderbook("SOLUSDT", 2)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOLUSDT - %v", err)
	} else {
		fmt.Println(orderbook)
	}

	fmt.Println()

	orderbook, err = h.GetOrderbook("SOL:USDC", 3)
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL:USDC - %v", err)
	} else {
		fmt.Println(orderbook)
	}
}
