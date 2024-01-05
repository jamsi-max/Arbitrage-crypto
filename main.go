package main

import (
	"log"
	"time"

	"github.com/jamsi-max/arbitrage/orderbook"
	"github.com/jamsi-max/arbitrage/providers"
	"github.com/jamsi-max/arbitrage/socket"
	// "github.com/jamsi-max/arbitrage//util"
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
	},
	"DOGEUSD": {
		"Binance":  "DOGEUSDT",
		"Kraken":   "XDG/USD",
		"Coinbase": "DOGE-USD", 
	},
	"BTCUSD": {
		"Binance":  "BTCUSDT",
		"Kraken":   "XBT/USD",
		"Coinbase": "BTC-USD", 
	},
	"ETHUSD": {
		"Binance":  "ETHUSDT",
		"Kraken":   "ETH/USD",
		"Coinbase": "ETH-USD",
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
	datach := make(chan orderbook.DataFeed, 1024)
	pvrs := []orderbook.Provider{
		providers.NewKrakenProvider(datach, mapSymbolsFor("Kraken")),
		providers.NewCoinbaseProvider(datach, mapSymbolsFor("Coinbase")),
		providers.NewBinanceOrderbooks(datach, mapSymbolsFor("Binance")),
	}

	for _, provider := range pvrs {
		if err := provider.Start(); err != nil {
			log.Fatal(err)
		}
	}

	bestSpreadch := make(chan map[string][]orderbook.BestSpread, 1024)
	go func() {
		ticker := time.NewTicker(time.Millisecond * 100)
		for {
			calcBestSpreads(bestSpreadch, pvrs)
			<-ticker.C
		}
  }()

	socketServer := socket.NewServer(bestSpreadch)
	socketServer.Start()
	// go func() {
	// 	ticker := time.NewTicker(time.Millisecond * 100)
	// 	for {
	// 		for _, p := range pvrs {
	// 			for _, book := range p.GetOrderbooks() {
	// 				var (
	// 					spread = book.Spread()
	// 					bestAsk = book.BestAsk()
	// 					bestBid = book.BestBid()
	// 				)
	// 				if bestAsk == nil || bestBid == nil {
	// 					continue
	// 				}
	// 				datach <- orderbook.DataFeed{
	// 					Provider: p.Name(),
	// 					Symbol: book.Symbol,
	// 					BestAsk: bestAsk.Price,
	// 					BestBid: bestBid.Price,
	// 					Spread: spread,
	// 				}
	// 			}
	// 		}
	// 		<-ticker.C
	// 	}
  // }()

	

	// for data := range datach {
	// 	fmt.Println(data)
	// }

	// for data := range datach {
	// 		fmt.Printf("[%s | %s] ASK %f %f BID [%f]\n", data.Provider, data.Symbol, data.BestAsk, data.BestBid, data.Spread)
	// }

	// select {}
}

func calcBestSpreads(datach chan map[string][]orderbook.BestSpread, pvrs []orderbook.Provider) {
	data := map[string][]orderbook.BestSpread{}

	for _, symbol := range symbols{
		bestSpreads := []orderbook.BestSpread{}
		for i := 0; i < len(pvrs); i++ {
			a := pvrs[i]
			var b orderbook.Provider
			if len(pvrs)-1 == i {
				b = pvrs[0]
			} else {
				b = pvrs[i+1]
			}
			bookA := a.GetOrderbooks()[getSymbolForProvider(a.Name(), symbol)]
			bookB := b.GetOrderbooks()[getSymbolForProvider(b.Name(), symbol)]

			best := orderbook.BestSpread{
				Symbol: symbol,
			}

			bestBidA := bookA.BestBid()
			bestBidB := bookB.BestBid()
			if bestBidA == nil || bestBidB == nil {
				continue
			}

			if bestBidA.Price < bestBidB.Price {
				best.A = a.Name()
				best.B = b.Name()
				best.BestBid = bestBidA.Price
				best.BestAsk = bookB.BestAsk().Price
			} else {
				best.A = b.Name()
				best.B = a.Name()
				best.BestBid = bestBidB.Price
				best.BestAsk = bookA.BestAsk().Price
			}

			best.Spread = util.Round(best.BestAsk-best.BestBid, 10_000)
			bestSpreads = append(bestSpreads, best)
		}
		data[symbol] = bestSpreads
	}
	datach <- data
}



// package main

// import (
// 	"fmt"
// 	"log"

// 	"github.com/gorilla/websocket"
// 	"github.com/jamsi-max/arbitrage/providers"
// )

// type BybitMessage struct {
// 	Op   string   `json:"op"`
// 	Args []string `json:"args"`
// }

// func main() {

// 	bybit := providers.NewBybitProvider([]string{"orderbook.1.BTCUSDT"})
// 	if err := bybit.Start(); err != nil {
// 		log.Fatal(err)
// 	}

// 	return
// 	ws, _, err := websocket.DefaultDialer.Dial("wss://stream.bybit.com/v5/public/linear", nil)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	go func() {
// 		for {
// 			_, b, err := ws.ReadMessage()
// 			if err != nil {
// 				log.Fatalln(err)
// 			}

// 			fmt.Println(string(b))
// 		}
// 	}()

// 	msg := BybitMessage{
// 		Op:   "subscribe",
// 		Args: []string{"orderbook.1.BTCUSDT"},
// 	}

// 	if err = ws.WriteJSON(msg); err != nil {
// 		log.Fatalln(err)
// 	}
// 	select {}
// }
