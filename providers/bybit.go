package providers

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/jamsi-max/arbitrage/orderbook"
)

type BybitMessage struct{
	Op string `json:"op"`
	Args []string `json:"args"`
}

type BybitProvider struct{
	Orderbooks orderbook.Orderbook
	symbols []string
}

func NewBybitProvider(feedch chan orderbook.DataFeed, symbols []string) *BinanceProvider {
	books := orderbook.Orderbooks{}
	for _, symbol := range symbols {
		books[symbol] = orderbook.NewBook(symbol)
	}

	return &BinanceProvider{
		Orderbooks: books,
		symbols:    symbols,
		feedch:     feedch,
	}
}
// func NewBybitProvider(symbols []string) *BybitProvider  {
// 	books := orderbook.Orderbook{}
// 	for _, symbol := range symbols{
// 		books[symbol] = orderbook.NewBook(symbol)
// 	}
	
// 	return &BybitProvider{
// 		Orderbooks: books,
// 		symbols: symbols,
// 	}
// }

func (p *BybitProvider) Start() error{
	ws, _, err := websocket.DefaultDialer.Dial("wss://stream.bybit.com/v5/public/linear", nil)
	if err != nil{
		log.Fatal(err)
	}

	msg := BybitMessage{
		Op: "subscribe",
		Args: []string{"orderbook.1.BTCUSDT"},
	}
	
	if err = ws.WriteJSON(msg); err!=nil{
		log.Fatal(err)
	}

	go func(){
		for{
			msg := BybitSocketResponse{}
			if err := ws.ReadJSON(&msg); err != nil {
				log.Fatal(err)
				break
			}

			fmt.Println(msg)
		}
	}()

	return nil
}

type BybitSocketResponse struct{
	Topic string        `json:"topic"`
	Type string         `json:"type"`
	Data BybitOrderbook `json:"data"`
}

type BybitOrderbook struct {
	S string  `json:"s"`
	B []Entry `json:"b"`
	A []Entry `json:"a"`
}

type Entry [2]string