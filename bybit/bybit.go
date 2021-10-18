package bybit

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
	TradeHeaders = "receive_time,timestamp,trade_time_ms,price,side,size,tick_direction,trade_id,cross_seq"
	WssUrl       = "wss://stream.bytick.com/realtime"
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

type BybitWebsocketCallback struct {
	tradeFile *os.File
	writer    *bufio.Writer
}

func NewWebsocketCallback(f *os.File, w *bufio.Writer) *BybitWebsocketCallback {
	return &BybitWebsocketCallback{
		tradeFile: f,
		writer:    w,
	}
}

func (c *BybitWebsocketCallback) GetSubscribeTradesParam(symbol string) string {
	return fmt.Sprintf(`{
		"op": "subscribe",
		"args": ["trade.%s"]}`, symbol)
}

func (c *BybitWebsocketCallback) OnReceiveTrade(msg string, errCh chan<- string) {

	var t Trade
	err := json.Unmarshal([]byte(msg), &t)
	if err != nil {
		errCh <- fmt.Sprintf("failed to unmarchal json to bybit.Trade: %v", err)
	}

	fmt.Printf("bybit(%v): %v\n", len(t.Data), t)

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
		c.writer.WriteString(fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n",
			now.Format(base.TimeFormat),
			d.Timestamp.Format(base.TimeFormat),
			d.TradeTimeMs,
			d.Price,
			d.Symbol,
			d.Side,
			d.Size,
			d.TickDirection,
			d.TradeId,
			d.CrossSeq))
	}
}
