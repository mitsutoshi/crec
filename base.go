package base

import "golang.org/x/net/websocket"

type WebsocketCallback interface {
	OnReceiveTrade(t interface{})
	SubscribeTrades(symbol string) error
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
