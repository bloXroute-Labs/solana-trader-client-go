package ws_test

import (
	"github.com/bloXroute-Labs/serum-api/bxserum/provider"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOrderbook(t *testing.T) {
	w, err := provider.NewWSClient()
	require.Nil(t, err)

	orderbook, err := w.GetOrderbook("ETH-USDT")
	require.Nil(t, err)
	require.Equal(t, "ETH-USDT", orderbook.Market)
	require.Equal(t, "9xQeWvG816bUx9EPjHmaT23yvVM2ZWbrrpZb9PusVFin", orderbook.MarketAddress)

	orderbook, err = w.GetOrderbook("FRONT/USDC")
	require.Nil(t, err)
	require.Equal(t, "ETH/USDT", orderbook.Market)
	require.Equal(t, "B95oZN5HCLGmFAhbzReWBA9cuSGPFQAXeuhm2FfpdrML", orderbook.MarketAddress)

	orderbook, err = w.GetOrderbook("soFRONTUSDC")
	require.Nil(t, err)
	require.Equal(t, "soFRONTUSDC", orderbook.Market)
	require.Equal(t, "7oKqJhnz9b8af8Mw47dieTiuxeaHnRYYGBiqCrRpzTRD", orderbook.MarketAddress)

	orderbook, err = w.GetOrderbook("scnSOL:USDC")
	require.Nil(t, err)
	require.Equal(t, "scnSOL:USDC", orderbook.Market)
	require.Equal(t, "D52sefGCWho2nd5UGxWd7wCftAzeNEMNYZkdEPGEdQTb", orderbook.MarketAddress)
}

func TestMarkets(t *testing.T) {
	w, err := provider.NewWSClient()
	require.Nil(t, err)

	markets, err := w.GetMarkets()
	require.Nil(t, err)

	/*
		market, ok := markets.Markets["ETH-USDT"] // Is this
		assert.True(t, ok)

		market, ok = markets.Markets["SOL/USDT"]
		assert.True(t, ok)
	*/
}
