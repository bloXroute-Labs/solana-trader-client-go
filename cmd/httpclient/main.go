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
	w := provider.NewHTTPClient()

	// One time request
	orderbook, err := w.GetOrderbook("ETH-USDT")
	if err != nil {
		log.Errorf("error with GetOrderbook request for ETH/USDT - %v", err)
	} else {
		fmt.Println(orderbook)
	}

	orderbook, err = w.GetOrderbook("SOL-USDT")
	if err != nil {
		log.Errorf("error with GetOrderbook request for SOL/USDT - %v", err)
	} else {
		fmt.Println(orderbook)
	}

	orderbook, err = w.GetOrderbook("ETH-USDC")
	if err != nil {
		log.Errorf("error with GetOrderbook request for ETH/USDC - %v", err)
	} else {
		fmt.Println(orderbook)
	}

	w = provider.NewHTTPClientWithEndpoint("http://1.1.1.1:1809")

	orderbook, err = w.GetOrderbook("ETH-USDC")
	if err != nil {
		log.Errorf("error with GetOrderbook request for ETH/USDC - %v", err)
	} else {
		fmt.Println(orderbook)
	}
}
