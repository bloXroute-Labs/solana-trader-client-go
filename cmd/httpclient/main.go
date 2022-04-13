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

	// One time request
	orderbook, err := h.GetOrderbook("ETH-USDT")
	if err != nil {
		log.Errorf("error with GetOrderbook request for ETH/USDT - %v", err)
	} else {
		fmt.Println(orderbook)
	}
	fmt.Println()

	orderbook, err = h.GetOrderbook("SOL-USDT")
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL/USDT - %v", err)
	} else {
		fmt.Println(orderbook)
	}
	fmt.Println()

	orderbook, err = h.GetOrderbook("ETH-USDC")
	if err != nil {
		log.Errorf("error with GetOrderbook request for ETH/USDC - %v", err)
	} else {
		fmt.Println(orderbook)
	}
	fmt.Println()
}
