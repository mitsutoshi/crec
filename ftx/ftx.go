package ftx

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
	TradeHeaders      = "receive_time,id,price,size,side,liquidation,time"
	WssUrl            = "wss://ftx.com/ws" // public websocket url
	msgTypeSubscribed = "subscribed"
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

type FtxWebsocketCallback struct {
	tradeFile *os.File
	writer    *bufio.Writer
}

func NewWebsocketCallback(f *os.File, w *bufio.Writer) *FtxWebsocketCallback {
	return &FtxWebsocketCallback{
		tradeFile: f,
		writer:    w,
	}
}

func (c *FtxWebsocketCallback) GetSubscribeTradesParam(symbol string) string {
	return fmt.Sprintf(`{
		"op": "subscribe",
		"channel": "trades",
		"market": "%s"}`, symbol)
}

func (c *FtxWebsocketCallback) OnReceiveTrade(msg string, errCh chan<- string) {

	var t Trade
	err := json.Unmarshal([]byte(msg), &t)
	if err != nil {
		errCh <- fmt.Sprintf("failed to unmarchal json to ftx.Trade: %v", err)
	}

	fmt.Printf("ftx(%v): %v\n", len(t.Data), t)

	if t.Type == msgTypeSubscribed {
		fmt.Println("channel has being subscribed.")
		return
	}

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

	for _, d := range t.Data {
		c.writer.WriteString(fmt.Sprintf("%v,%v,%v,%.8f,%s,%v,%v\n",
			now.Format(base.TimeFormat),
			d.Id,
			d.Price,
			d.Size,
			d.Side,
			d.Liquidation,
			d.Time.Format(base.TimeFormat)))
	}
}
