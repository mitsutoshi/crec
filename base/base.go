package base

import (
	"log"

	"golang.org/x/net/websocket"
)

type WebsocketCallback interface {
	OnReceiveTrade(msg string, errCh chan<- string)
	GetSubscribeTradesParam(symbol string) string
}

type Websocket struct {
	Conn     *websocket.Conn
	Callback WebsocketCallback
}

func (ws *Websocket) Connect(wssUrl string, originUrl string) error {
	conn, err := websocket.Dial(wssUrl, "", originUrl)
	if err != nil {
		return err
	}
	ws.Conn = conn
	return nil
}

func (c *Websocket) SubscribeTrades(symbol string) error {
	param := c.Callback.GetSubscribeTradesParam(symbol)
	err := websocket.Message.Send(c.Conn, param)
	log.Printf("Sent subscribe reuqest for [trade.%s]\n", symbol)
	return err
}

func (ws *Websocket) Receive(errCh chan<- string) {
	var msg string
	for {
		websocket.Message.Receive(ws.Conn, &msg)
		ws.Callback.OnReceiveTrade(msg, errCh)
	}
}
