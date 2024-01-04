package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jamsi-max/arbitrage/orderbook"
	"github.com/jamsi-max/arbitrage/providers"
)

var symbols = []string{
	"BTCUSD",
	"ETHUSD",
	"DOGEUSD",
}

var pairs = map[string]map[string]string{
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

func mapSymbolsFor(provider string) []string {
	out := make([]string,len(symbols))
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

	ticker := time.NewTicker(time.Millisecond * 100)
	go func() {
		for {
			for _, p := range pvrs {
				for _, book := range p.GetOrderbooks() {
					var (
						spread = book.Spread()
						bestAsk = book.BestAsk()
						bestBid = book.BestBid()
					)
					if bestAsk == nil || bestBid == nil {
						continue
					}
					datach <- orderbook.DataFeed{
						Provider: p.Name(),
						Symbol: book.Symbol,
						BestAsk: bestAsk.Price,
						BestBid: bestBid.Price,
						Spread: spread,
					}
				}
			}
			<-ticker.C
		}
  }()

	

	// for data := range datach {
	// 	fmt.Println(data)
	// }
	for data := range datach {
			fmt.Printf("[%s | %s] ASK %f %f BID [%f]\n", data.Provider, data.Symbol, data.BestAsk, data.BestBid, data.Spread)
	}

	select {}
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
