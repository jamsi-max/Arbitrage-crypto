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
	// "BTCUSD",
	// "ETHUSD",
	// "DOGEUSD",
	// "ADAUSD",
	// "LTCUSD",
	"SOLUSD",
	"XLMUSD",
	"WAXPUSD",
}

var pairs = map[string]map[string]string{
	"XLMUSD": {
		"Binance": "XLMUSDT",
		// "Kraken":   "XLM/USD",
		"Coinbase": "XLM-USD",
		"Bybit":    "orderbook.1.XLMUSDT",
		"Cucoin":   "/market/level2:XLM-USDT",
		"OKX":      "XLM-USDT",
	},
	"SOLUSD": {
		"Binance":  "SOLUSDT",
		"Kraken":   "SOL/USD",
		"Coinbase": "SOL-USD",
		"Bybit":    "orderbook.1.SOLUSDT",
		"Cucoin":   "/market/level2:SOL-USDT",
		"OKX":      "SOL-USDT",
	},
	"WAXPUSD": {
		"Binance": "WAXPUSDT",
		// "Kraken":   "WAX/USD",
		// "Coinbase": "WAXT-USD",
		"Bybit":  "orderbook.1.WAXPUSDT",
		"Cucoin": "/market/level2:WAX-USDT",
		"OKX":      "WAXP-USDT",
	},
	// "ADAUSD": {
	// 	"Binance":  "ADAUSDT",
	// 	"Kraken":   "ADA/USD",
	// 	"Coinbase": "ADA-USD",
	// 	"Bybit":    "orderbook.1.ADAUSDT",
	// 	"Cucoin":   "/market/ticker:ADA-USDT",
	// },
	// "DOGEUSD": {
	// 	"Binance":  "DOGEUSDT",
	// 	"Kraken":   "XDG/USD",
	// 	"Coinbase": "DOGE-USD",
	// 	"Bybit":    "orderbook.1.DOGEUSDT",
	// 	"Cucoin":   "/market/ticker:DOGE-USDT",
	// },
	// "BTCUSD": {
	// 	"Binance":  "BTCUSDT",
	// 	// "Kraken":   "XBT/USD",
	// 	"Coinbase": "BTC-USD",
	// 	"Bybit":    "orderbook.1.BTCUSDT",
	// 	"Cucoin":   "/market/ticker:BTC-USDT",
	// },
	// "ETHUSD": {
	// 	"Binance":  "ETHUSDT",
	// 	"Kraken":   "ETH/USD",
	// 	"Coinbase": "ETH-USD",
	// 	"Bybit":    "orderbook.1.ETHUSDT",
	// 	"Cucoin":   "/market/ticker:ETH-USDT",
	// },
	// "LTCUSD": {
	// 	"Binance":  "LTCUSDT",
	// 	"Kraken":   "LTC/USD",
	// 	"Coinbase": "LTC-USD",
	// 	"Bybit":    "orderbook.1.LTCUSDT",
	// 	"Cucoin":   "/market/ticker:LTC-USDT",
	// },
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
		// providers.NewKrakenProvider(mapSymbolsFor("Kraken")),
		providers.NewCoinbaseProvider(mapSymbolsFor("Coinbase")),
		providers.NewBinanceProvider(mapSymbolsFor("Binance")),
		providers.NewBybitProvider(mapSymbolsFor("Bybit")),
		providers.NewCucoinProvider(mapSymbolsFor("Cucoin")),
		providers.NewOKXProvider(mapSymbolsFor("OKX")),
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

	for _, symbol := range symbols {
		crossSpreads := []orderbook.CrossSpread{}
		// for i := 0; i < len(pvrs); i++ {
		for i, j := 0, 1; i < len(pvrs)-1; {
			a := pvrs[i]
			b := pvrs[j]

			// a := pvrs[i]
			// var b orderbook.Provider
			// if len(pvrs)-1 == i {
			// 	b = pvrs[0]
			// } else {
			// 	b = pvrs[i+1]
			// }

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
			// fmt.Println(i, j, a.Name(), bestBidA, b.Name(), bestBidB)
			if bestBidA == nil || bestBidB == nil {
				if j < len(pvrs)-1 {
					j++
				} else {
					i++
					j = i + 1
				}
				continue
			}

			//DEBUG
			// fmt.Printf("%+v", len(pvrs))
			// fmt.Println(i, j, a.Name(), bestBidA.Price, b.Name(), bestBidB.Price)
			// if b.Name() == "Cucoin"  {
			// 	log.Println(b.Name(), bookB.BestAsk().Price)
			// }
			// DEBUG END

			if bestBidA.Price < bestBidB.Price {
				// log.Println("a<b", a.Name(),bestBidA.Price, b.Name(), bestBidB.Price, "b-a:", bestBidB.Price-bestBidA.Price)
				bestBid.Provider = a.Name()
				bestAsk.Provider = b.Name()
				bestBid.Price = bestBidA.Price
				bestBid.Size = bestBidA.TotalVolume
				if bookB.BestAsk() != nil {
					bestAsk.Price = bookB.BestAsk().Price
					bestAsk.Size = bookB.BestAsk().TotalVolume
				}
				// if symbol == "SOLUSD" && (a.Name() == "Cucoin" || b.Name() == "Cucoin"){
				// 	log.Println("a<b",symbol, a.Name(), bestBid.Price, b.Name(), bestAsk.Price)
				// }
			} else {
				bestBid.Provider = b.Name()
				bestAsk.Provider = a.Name()
				bestBid.Price = bestBidB.Price
				bestBid.Size = bestBidB.TotalVolume
				if bookA.BestAsk() != nil {
					bestAsk.Price = bookA.BestAsk().Price
					bestAsk.Size = bookA.BestAsk().TotalVolume
				}
				// if symbol == "SOLUSD" && (a.Name() == "Cucoin" || b.Name() == "Cucoin") {
				// 	log.Println("a>b",symbol, b.Name(), bestBid.Price, a.Name(), bestAsk.Price)
				// }
			}

			crossSpread.Spread = util.Round(bestAsk.Price-bestBid.Price, 10_0000)
			crossSpread.BestAsk = bestAsk
			crossSpread.BestBid = bestBid
			crossSpreads = append(crossSpreads, crossSpread)

			// if len(pvrs) == 2 {
			// 	continue
			// } else if j < len(pvrs)-2 {
			// 	j++
			// } else {
			// 	i++
			// 	j = i
			// }
			// if len(pvrs) == 2 {
			// 		break
			// }
			if j < len(pvrs)-1 {
				j++
			} else {
				i++
				j = i + 1
			}

		}
		// fmt.Printf("%v -> %+v",symbol, data[symbol])
		data[symbol] = crossSpreads
	}
	datach <- data
}
