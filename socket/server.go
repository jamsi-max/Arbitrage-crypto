package socket

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jamsi-max/arbitrage/orderbook"
)

var symbols = []string{
	"BTCUSD",
	"ETHUSD",
	"DOGEUSD",
	"ADAUSD",
}

type Message struct {
	Type    string   `json:"type"`
	Topic   string   `json:"topic"`
	Symbols []string `json:"symbols"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Server struct {
	bsch chan map[string][]orderbook.BestSpread
	lock sync.RWMutex
	// conns map[string][]*websocket.Conn
	// conns map[*websocket.Conn]bool
	conns map[string]map[*websocket.Conn]bool
}

func NewServer(bsch chan map[string][]orderbook.BestSpread) *Server {
	s :=  &Server{
		bsch: bsch,
		conns: make(map[string]map[*websocket.Conn]bool),
	}

	for _, symbol := range symbols {
		s.conns[symbol] = map[*websocket.Conn]bool{}
	}

	return s
}

func (s *Server) Start() error {
	http.HandleFunc("/bestspreads", s.handleBestSpreads)
	http.HandleFunc("/", s.handleWS)
	go s.writeLoop()
	return http.ListenAndServe(":4000", nil)
}

func (s *Server) unregisterConn(ws *websocket.Conn) {
	s.lock.Lock()
	// delete(s.conns, ws)
	s.lock.Unlock()

	fmt.Printf("unregistered connection %s\n", ws.RemoteAddr())

	ws.Close()
}

func (s *Server) registerConn(symbol string, ws *websocket.Conn) {
	s.conns[symbol][ws] = true

	fmt.Printf("registered connection to symbol %s %s\n", symbol, ws.RemoteAddr())
}

func (s *Server) writeLoop() {
	for data := range s.bsch {
		for symbol, spreads := range data {
			for ws := range s.conns[symbol] {
				ws.WriteJSON(spreads)
			}
		}
	}
}

func (s *Server) readLoop(ws *websocket.Conn) {
	defer func() {
		s.unregisterConn(ws)	
	}()

	for {
		msg := Message{}
		if err := ws.ReadJSON(&msg); err != nil {
			fmt.Println("socket read error:", err)
			break
		}
		if err := s.handleSocketMessage(ws, msg); err != nil{
			fmt.Println("handle msg error:", err)
			break
		}
	}
}

func (s *Server) handleSocketMessage(ws *websocket.Conn, msg Message) error {
	for _, symbol := range msg.Symbols {
		s.registerConn(symbol, ws)
	}
	return nil
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket upgrade error:", err)
	}

	ws.WriteJSON(map[string]string{"version": "0.1"})

	go s.readLoop(ws)

}

func (s *Server) handleBestSpreads(w http.ResponseWriter, r *http.Request) {
	// ws, err := upgrader.Upgrade(w, r, nil)
	// if err != nil {
	// 	log.Println("websocket upgrade error:", err)
	// }

	// s.registerConn(ws)
}