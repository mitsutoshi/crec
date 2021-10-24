package gmo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mitsutoshi/crec/base"
	"github.com/mitsutoshi/crec/utils"
)

const (
	OriginUrl    = "https://api.coin.z.com"            // rest api server url
	WssUrl       = "wss://api.coin.z.com/ws/public/v1" // public websocket url
	TradeHeaders = "receive_time,symbol,side,size,price,timestamp"
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

type GmoWebsocketCallback struct {
	tradeFile *os.File
	writer    *bufio.Writer
}

func NewWebsocketCallback(f *os.File, w *bufio.Writer) *GmoWebsocketCallback {
	return &GmoWebsocketCallback{
		tradeFile: f,
		writer:    w,
	}
}

func (c *GmoWebsocketCallback) GetSubscribeTradesParam(symbol string) string {
	return fmt.Sprintf(`{
		"command": "subscribe",
		"channel": "trades",
		"symbol": "%s",
		"option": "TAKER_ONLY"}`, symbol)
}

func (c *GmoWebsocketCallback) OnReceiveTrade(msg string, errCh chan<- string) {

	var t Trade
	err := json.Unmarshal([]byte(msg), &t)
	if err != nil {
		errCh <- fmt.Sprintf("failed to unmarchal json to gmo.Trade: %v", err)
	}
	fmt.Printf("gmo: %v\n", t)

	now := time.Now().UTC()
	names := strings.Split(c.tradeFile.Name(), ".")
	date := names[0][strings.LastIndex(names[0], "_")+1:]

	if date != now.Format("20060102") {

		// get new file name
		newName := fmt.Sprintf("%s_%s.csv", names[0][:strings.LastIndex(names[0], "_")], date)
		fmt.Printf("Close current file(%s) and open new file(%s).", c.tradeFile.Name(), newName)

		// close and open file
		c.writer.Flush()
		c.tradeFile.Close()
		f, err := utils.OpenNewFile(newName)
		if err != nil {
			panic(err) // TODO
		}

		// replace file
		c.tradeFile = f
		c.writer = bufio.NewWriter(f)
	}

	c.writer.WriteString(fmt.Sprintf("%v,%s,%s,%.8f,%.3f,%v\n",
		now.Format(base.TimeFormat),
		t.Symbol,
		t.Side,
		t.Size,
		t.Price,
		t.Timestamp.Format(base.TimeFormat)))
}
