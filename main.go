package main

import (
	"log"
	"time"

	"github.com/jamsi-max/arbitrage/orderbook"
	"github.com/jamsi-max/arbitrage/providers"
	sttg "github.com/jamsi-max/arbitrage/settings"
	"github.com/jamsi-max/arbitrage/socket"
	"github.com/jamsi-max/arbitrage/spread"
)

func mapSymbolsFor(provider string) []string {
	out := make([]string, len(sttg.Symbols))
	for i, symbol := range sttg.Symbols {
		out[i] = sttg.Pairs[symbol][provider]
	}
	return out
}

func main() {
	pvrs := []orderbook.Provider{
		// providers.NewKrakenProvider(mapSymbolsFor("Kraken")),
		providers.NewCoinbaseProvider(mapSymbolsFor("Coinbase")),
		providers.NewBinanceProvider(mapSymbolsFor("Binance")),
		providers.NewBybitProvider(mapSymbolsFor("Bybit")),
		providers.NewCucoinProvider(mapSymbolsFor("Cucoin")),
		providers.NewOKXProvider(mapSymbolsFor("OKX")),
		providers.NewMEXCProvider(mapSymbolsFor("MEXC")),
		providers.NewFinexPovider(mapSymbolsFor("Finex")),
	}

	for _, provider := range pvrs {
		if err := provider.Start(); err != nil {
			log.Fatal(err)
		}
	}

	crossSpreadch := make(chan map[string][]orderbook.CrossSpread, 1024)
	go func() {
		ticker := time.NewTicker(time.Millisecond * 100)
		for {
			spread.CalcCrossSpreads(crossSpreadch, pvrs)
			<-ticker.C
		}
	}()

	socketServer := socket.NewServer(crossSpreadch)
	socketServer.Start()
}
