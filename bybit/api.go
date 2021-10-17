package bybit

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/net/websocket"
)

type Trade struct {
	Topic string `json:"topic"`
	Data  []Data `json:"data"`
}

type Data struct {
	Timestamp     time.Time `json:"timestamp"`
	TradeTimeMs   int       `json:"trade_time_ms"`
	Price         float64   `json:"price"`
	Symbol        string    `json:"symbol"`
	Side          string    `json:"side"`
	Size          int       `json:"size"`
	TickDirection string    `json:"tick_direction"`
	TradeId       string    `json:"trade_id"`
	CrossSeq      int       `json:"cross_seq"`
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
		"args": ["trade.%s"]}`, symbol)

	err := websocket.Message.Send(pws.ws, param)
	fmt.Printf("Sent subscribe reuqest for [trade.%s]\n", symbol)
	return err
}

func (pws *Websocket) Receive(errCh chan<- string) {

	var msg string
	for {
		websocket.Message.Receive(pws.ws, &msg)

		var tr Trade
		err := json.Unmarshal([]byte(msg), &tr)
		if err != nil {
			errCh <- fmt.Sprintf("failed to unmarchal json to bybit.Trade: %v", err)
		}
		pws.Callback.OnReceiveTrade(tr)
	}
}
