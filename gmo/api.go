package gmo

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

const (
	publicWssUrl  = WssUrl + "/public/v1" // public websocket url
	channelTrades = "trades"              // websocket channel: trades
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
	ws, err := websocket.Dial(publicWssUrl, "", OriginUrl)
	if err != nil {
		return err
	}
	pws.ws = ws
	return nil
}

func (pws *Websocket) SubscribeTrades(symbol string, takerOnly bool) error {

	// if TAKER_ONLY option is enable
	to := ""
	if takerOnly {
		to = "TAKER_ONLY"
	}

	param := fmt.Sprintf(`{
		"command": "subscribe",
		"channel": "%s",
		"symbol": "%s",
		"option": "%s"}`, channelTrades, symbol, to)

	err := websocket.Message.Send(pws.ws, param)
	log.Printf("Sent subscribe reuqest for [%s]\n", channelTrades)
	return err
}

func (pws *Websocket) Receive(errCh chan<- string) {

	var msg string
	for {
		websocket.Message.Receive(pws.ws, &msg)

		if msg == "" {
			continue
		}

		var tr Trade
		err := json.Unmarshal([]byte(msg), &tr)
		if err != nil {
			errCh <- fmt.Sprintf("failed to unmarchal json to gmo.Trade: %v", err)
		}
		pws.Callback.OnReceiveTrade(tr)
	}
}
