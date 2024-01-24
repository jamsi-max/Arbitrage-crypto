package main

import (
	"log"
	"time"

	"github.com/jamsi-max/arbitrage/orderbook"
	"github.com/jamsi-max/arbitrage/providers"
	"github.com/jamsi-max/arbitrage/socket"
	"github.com/jamsi-max/arbitrage/util"
)

var symbols = []string{
	"BTCUSD",
	"ETHUSD",
	"DOGEUSD",
	"ADAUSD",
}

var pairs = map[string]map[string]string{
	"ADAUSD": {
		"Binance":  "ADAUSDT",
		"Kraken":   "ADA/USD",
		"Coinbase": "ADA-USD",
		"Bybit":    "orderbook.1.ADAUSDT",
	},
	"DOGEUSD": {
		"Binance":  "DOGEUSDT",
		"Kraken":   "XDG/USD",
		"Coinbase": "DOGE-USD",
		"Bybit":    "orderbook.1.DOGEUSDT", 
	},
	"BTCUSD": {
		"Binance":  "BTCUSDT",
		"Kraken":   "XBT/USD",
		"Coinbase": "BTC-USD",
		"Bybit":    "orderbook.1.BTCUSDT", 
	},
	"ETHUSD": {
		"Binance":  "ETHUSDT",
		"Kraken":   "ETH/USD",
		"Coinbase": "ETH-USD",
		"Bybit":    "orderbook.1.ETHUSDT", 
	},
}

func getSymbolForProvider(p string, symbol string) string {
	return pairs[symbol][p]
}

func mapSymbolsFor(provider string) []string {
	out := make([]string, len(symbols))
	for i, symbol := range symbols {
		out[i] = pairs[symbol][provider]
	}
	return out
}

func main() {
	pvrs := []orderbook.Provider{
		providers.NewKrakenProvider(mapSymbolsFor("Kraken")),
		providers.NewCoinbaseProvider(mapSymbolsFor("Coinbase")),
		providers.NewBinanceOrderbooks(mapSymbolsFor("Binance")),
		providers.NewBybitProvider(mapSymbolsFor("Bybit")),
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
			calcCrossSpreads(crossSpreadch, pvrs)
			<-ticker.C
		}
  }()

	socketServer := socket.NewServer(crossSpreadch)
	socketServer.Start()
}

func calcCrossSpreads(datach chan map[string][]orderbook.CrossSpread, pvrs []orderbook.Provider) {
	data := map[string][]orderbook.CrossSpread{}

	for _, symbol := range symbols{
		crossSpreads := []orderbook.CrossSpread{}
		for i := 0; i < len(pvrs); i++ {
			a := pvrs[i]
			var b orderbook.Provider
			if len(pvrs)-1 == i {
				b = pvrs[0]
			} else {
				b = pvrs[i+1]
			}

			var (
				crossSpread = orderbook.CrossSpread{
					Symbol: symbol,
				}
				bestAsk  = orderbook.BestPrice{}
				bestBid  = orderbook.BestPrice{}
				bookA    = a.GetOrderbooks()[getSymbolForProvider(a.Name(), symbol)]
			  bookB    = b.GetOrderbooks()[getSymbolForProvider(b.Name(), symbol)]
				bestBidA = bookA.BestBid()
			  bestBidB = bookB.BestBid()
			)
			if bestBidA == nil || bestBidB == nil {
				continue
			}
			if bestBidA.Price < bestBidB.Price {
				bestBid.Provider = a.Name()
				bestAsk.Provider = b.Name()
				bestBid.Price = bestBidA.Price
				bestBid.Size = bestBidA.TotalVolume
				bestAsk.Price = bookB.BestAsk().Price
				bestAsk.Size = bookB.BestAsk().TotalVolume
			} else {
				bestBid.Provider = b.Name()
				bestAsk.Provider = a.Name()
				bestBid.Price = bestBidB.Price
				bestBid.Size = bestBidB.TotalVolume
				bestAsk.Price = bookA.BestAsk().Price
				bestAsk.Size = bookA.BestAsk().TotalVolume
			}

			crossSpread.Spread = util.Round(bestAsk.Price-bestBid.Price, 10_000)
			crossSpread.BestAsk = bestAsk
			crossSpread.BestBid = bestBid
			crossSpreads =append(crossSpreads, crossSpread)
		}
		data[symbol] = crossSpreads
	}
	datach <- data
}