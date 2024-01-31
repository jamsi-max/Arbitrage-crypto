package providers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"time"

	"github.com/Kucoin/kucoin-go-sdk"
	"github.com/gorilla/websocket"
	"github.com/jamsi-max/arbitrage/orderbook"
)

type TokenModel struct {
	Data *DataTokenModel
}

type DataTokenModel struct {
	Token           string `json:"token"`
	InstanceServers []InstanceServers
}

type InstanceServers struct {
	Endpoint string `json:"endpoint"`
	Encrypt  bool   `json:"encrypt"`
	Protocol string `json:"protocol"`
}

type CucoinProvider struct {
	Orderbooks orderbook.Orderbooks
	symbols    []string
}

func NewCucoinProvider(symbols []string) *CucoinProvider {
	books := orderbook.Orderbooks{}
	for _, symbol := range symbols {
		books[symbol] = orderbook.NewBook(symbol)
	}

	return &CucoinProvider{
		Orderbooks: books,
		symbols:    symbols,
	}
}

func (c *CucoinProvider) GetOrderbooks() orderbook.Orderbooks {
	return c.Orderbooks
}

func (c *CucoinProvider) Name() string {
	return "Cucoin"
}

const ApiGetPublickToken = "https://api.kucoin.com/api/v1/bullet-public"

func GetTokenAndEndpoint() (string, string) {
	data := []byte{}

	r := bytes.NewReader(data)
	resp, err := http.Post(ApiGetPublickToken, "application/json", r)
	if err != nil {
		log.Printf("Error get Token kukoin: %v", err)
	}
	defer resp.Body.Close()

	var t TokenModel
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		log.Printf("Error response Unmarshal kukoin: %v", err)
	}

	if t.Data.InstanceServers != nil && t.Data.InstanceServers[0].Endpoint != "" {
		return t.Data.Token, t.Data.InstanceServers[0].Endpoint
	}

	return t.Data.Token, "wss://ws-api-spot.kucoin.com/"
}

type MessageSubscribe struct {
	Id             int64  `json:"id"`
	Type           string `json:"type"`
	Topic          string `json:"topic"`
	PrivateChannel bool   `json:"privateChannel"`
	Response       bool   `json:"response"`
}

type CucoinMessage struct {
	Id    string `json:"id,omitempty"`
	Topic string `json:"topic"`
	Type  string `json:"type"`
	Data  *CucoinMessageData
}

type CucoinMessageData struct {
	BestAsk     string `json:"bestAsk"`
	BestAskSize string `json:"bestAskSize"`
	BestBid     string `json:"bestBid"`
	BestBidSize string `json:"bestBidSize"`
	Price       string `json:"price"`
}

func (c *CucoinProvider) Start() error {
	token, endpoint := GetTokenAndEndpoint()
	connectId := time.Now().UnixNano()
	wsURL := endpoint + "?token=" + token + "&[connectId=" + kucoin.IntToString(connectId) + "]"

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return err
	}
	ws.ReadMessage()

	for _, symbol := range c.symbols{
		ws.WriteJSON(MessageSubscribe{
			Id:             connectId,
			Type:           "subscribe",
			Topic:          symbol,
			PrivateChannel: false,
			Response:       true,
		})
	}
	
	go func() {
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("->", err)
				break
			}

			msg := CucoinMessage{}
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Println("-->", err)
				break
			}

			if msg.Type == "message" {
				c.cucoinHandleUpdate(msg.Topic, msg.Data)
				// log.Printf("%+v", msg)
			}
		}
	}()

	return nil
}

func (c *CucoinProvider) cucoinHandleUpdate(symbol string, data *CucoinMessageData) error {
	book := c.Orderbooks[symbol]

	priceAsk, _ := strconv.ParseFloat(data.BestAsk, 64)
	sizeAsk, _ := strconv.ParseFloat(data.BestAskSize, 64)
	book.Asks.Update(priceAsk, sizeAsk)

	priceBid, _ := strconv.ParseFloat(data.BestBid, 64)
	sizeBid, _ := strconv.ParseFloat(data.BestBidSize, 64)
	book.Bids.Update(priceBid, sizeBid)

	log.Println(symbol, priceAsk, sizeAsk, priceBid, sizeBid)

	return nil
}