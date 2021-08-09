package ftx

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/net/websocket"
)

type Trade struct {
	Channel string `json:"channel"`
	Market  string `json:"market"`
	Type    string `json:"type"`
	Data    []Data `json:"data"`
}

type Data struct {
	Id          int       `json:"id"`
	Price       float64   `json:"price"`
	Size        float64   `json:"size"`
	Side        string    `json:"side"`
	Liquidation bool      `json:"liquidation"`
	Time        time.Time `json:"time"`
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

func (pws *Websocket) SubscribeTrades(symbol string) error {

	param := fmt.Sprintf(`{
		"op": "subscribe",
		"channel": "%s",
		"market": "%s"}`, channelTrades, symbol)

	err := websocket.Message.Send(pws.ws, param)
	fmt.Printf("Sent subscribe reuqest for [%s]\n", channelTrades)
	return err
}

func (pws *Websocket) Receive(errCh chan<- string) {

	var msg string
	for {
		websocket.Message.Receive(pws.ws, &msg)

		// extract a channel name from received message
		var m map[string]interface{}
		err := json.Unmarshal([]byte(msg), &m)
		if err != nil {
			errCh <- fmt.Sprintf("failed to unmarchal json to map[string]interface{}: %v", err)
		}

		channel := m["channel"].(string)
		if channel == channelTrades {
			var tr Trade
			err := json.Unmarshal([]byte(msg), &tr)
			if err != nil {
				errCh <- fmt.Sprintf("failed to unmarchal json to ftx.Trade: %v", err)
			}
			pws.Callback.OnReceiveTrade(tr)
		}
	}
}
