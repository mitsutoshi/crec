package liquid

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

const (
	publicWssUrl  = WssUrl // public websocket url
	channelTrades = "price_ladders_cash_btcjpy_buy"
)

// Trade is type of "trades" channel received data from websocket
type Trade struct {
	Channel   string    `json:"channel"`
	Price     float64   `json:"price,string"`
	Side      string    `json:"side"`
	Size      float64   `json:"size,string"`
	Timestamp time.Time `json:"timestamp,string"`
	Symbol    string    `json:"symbol"`
}

type WebsocketCallback interface {
	OnReceiveTrade(t Trade)
}

type Websocket struct {
	ws       *websocket.Conn
	Callback WebsocketCallback
}

func (pws *Websocket) Connect() error {
	ws, err := websocket.Dial(publicWssUrl, "", publicWssUrl)
	if err != nil {
		return err
	}
	pws.ws = ws
	return nil
}

func (pws *Websocket) SubscribeTrades(symbol string, takerOnly bool) error {

	param := fmt.Sprintf(`{
		"command": "subscribe",
		"channel": "%s",
		"symbol": "%s",
		}`, channelTrades, symbol)

	err := websocket.Message.Send(pws.ws, param)
	log.Printf("Sent subscribe reuqest for [%s]\n", channelTrades)
	return err
}

func (pws *Websocket) Receive(errCh chan<- string) {

	var msg string
	for {
		websocket.Message.Receive(pws.ws, &msg)
		//fmt.Println(msg)
		if msg == "" {
			continue
		}

		if msg == "d payload" {
			//fmt.Println("ignore")
			continue
		}
		fmt.Println(msg)

		// extract a channel name from received message
		var m map[string]string
		err := json.Unmarshal([]byte(msg), &m)
		if err != nil {
			errCh <- fmt.Sprintf("failed to unmarchal json to map[string]interface{}: %v", err)
		}

		fmt.Println(m["event"])
		if m["event"] == "pusher:connection_established" {

			//data := map[string]string{
			//	"channel": "executions_cash_btcjpy",
			//}
			//param := map[string]interface{}{
			//	"event": "pusher:subscribe",
			//	"data":  data,
			//}

			param := `{
				"event": "pusher:subscribe",
				"data": {
					"channel": "executions_cash_btcjpy"
				}
			}`

			//s, _ := json.Marshal(param)
			//fmt.Println(string(s))
			//err := websocket.Message.Send(pws.ws, string(s))
			err := websocket.Message.Send(pws.ws, param)
			if err != nil {
				fmt.Println(err)
			}
			log.Printf("Sent subscribe reuqest for [%s]\n", param)
		}

		// received error
		//if val, ok := m["error"]; ok {
		//	errCh <- fmt.Sprintf("error liquid.Trade: %v", val)
		//}

		//channel := m["channel"].(string)
		//if channel == channelTrades {
		//	var tr Trade
		//	err := json.Unmarshal([]byte(msg), &tr)
		//	if err != nil {
		//		errCh <- fmt.Sprintf("failed to unmarchal json to liquid.Trade: %v", err)
		//	}
		//	pws.Callback.OnReceiveTrade(tr)
		//}
	}
}

type ReceiveData struct {
	Data  WsData `json:"data"`
	Event string `json:"event"`
}

type WsData struct {
	ActivityTimeout int    `json:"activity_timeout"`
	SocketId        string `json:"socket_id"`
}
