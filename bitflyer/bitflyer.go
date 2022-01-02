package bitflyer

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
	TradeHeaders = "receive_time,id,exec_date,price,size,side,buy_child_order_acceptance_id,sell_child_order_acceptance_id"
	WssUrl       = "wss://ws.lightstream.bitflyer.com/json-rpc"
)

type jsonRPC2 struct {
	Version string      `json:"version"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Id      int         `json:"id"`
}

type subscribeParams struct {
	Channel string `json:"channel"`
}

type Trade struct {
	Params Param `json:"params"`
}

type Param struct {
	Message []Data `json:"message"`
}

type Data struct {
	Id                         int64     `json:"id"`                             // id
	ExecDate                   time.Time `json:"exec_date"`                      // exec_date
	Price                      float64   `json:"price"`                          // price
	Size                       float64   `json:"size"`                           // size
	Side                       string    `json:"side"`                           // side
	BuyChildOrderAcceptanceId  string    `json:"buy_child_order_acceptance_id"`  // buy_child_order_acceptance_id
	SellChildOrderAcceptanceId string    `json:"sell_child_order_acceptance_id"` // sell_child_order_acceptance_id
}

type BitflyerWebsocketCallback struct {
	tradeFile *os.File
	writer    *bufio.Writer
}

func NewWebsocketCallback(f *os.File, w *bufio.Writer) *BitflyerWebsocketCallback {
	return &BitflyerWebsocketCallback{
		tradeFile: f,
		writer:    w,
	}
}

func (c *BitflyerWebsocketCallback) GetSubscribeTradesParam(symbol string) string {
	jsonrpc := &jsonRPC2{
		Version: "2.0",
		Method:  "subscribe",
		Params:  &subscribeParams{fmt.Sprintf("lightning_executions_%s", symbol)},
	}
	jsonValue, err := json.Marshal(jsonrpc)
	if err != nil {
		panic(err)
	}
	return string(jsonValue)
}

func (c *BitflyerWebsocketCallback) OnReceiveTrade(msg string, errCh chan<- string) {

	var t Trade
	err := json.Unmarshal([]byte(msg), &t)
	if err != nil {
		errCh <- fmt.Sprintf("failed to unmarchal json to bitflyer.Trade: %v", err)
	}

	fmt.Printf("bitflyer(%v): %v\n", len(t.Params.Message), t)

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

	for _, d := range t.Params.Message {
		c.writer.WriteString(fmt.Sprintf("%v,%v,%v,%.8f,%v,%v,%v,%v\n",
			now.Format(base.TimeFormat),
			d.Id,
			d.ExecDate.Format(base.TimeFormat),
			d.Price,
			d.Size,
			d.Side,
			d.BuyChildOrderAcceptanceId,
			d.SellChildOrderAcceptanceId))
	}
}
