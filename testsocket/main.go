package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/jamsi-max/arbitrage/socket"
)

func main()  {
	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:4000", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer ws.Close()

	msg := socket.Message{
		Type: "subscribe",
		Topic: "spreads",
		Symbols: []string{"BTCUSD"},
	}

	if err := ws.WriteJSON(msg); err != nil {
		log.Fatal(err)
	}

	for {
		_, msg, err := ws.ReadMessage()
		if  err != nil {
			return
		}
		fmt.Println(msg)
	}
}