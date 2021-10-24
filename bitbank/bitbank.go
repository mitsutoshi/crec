package bitbank

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mitsutoshi/crec/base"
	"github.com/mitsutoshi/crec/utils"
)

const (
	TradeHeaders = "receive_time,transaction_id,side,price,amount,executed_at"
	WssUrl       = "wss://stream.bitbank.cc/socket.io/?transport=websocket&EIO=4"
	keyword      = `"message",`
)

var lastTransactionId int = 0

type Trade struct {
	RoomName string  `json:"room_name"`
	Message  message `json:"message"`
}

type message struct {
	Data data `json:"data"`
}

type data struct {
	Transactions []transaction `json:"transactions"`
}

type transaction struct {
	TransactionId int     `json:"transaction_id"`
	Side          string  `json:"side"`
	Price         float64 `json:"price,string"`
	Amount        float64 `json:"amount,string"`
	ExecutedAt    int     `json:"executed_at"`
}

type BitbankWebsocketCallback struct {
	tradeFile *os.File
	writer    *bufio.Writer
}

func NewWebsocketCallback(f *os.File, w *bufio.Writer) *BitbankWebsocketCallback {
	return &BitbankWebsocketCallback{
		tradeFile: f,
		writer:    w,
	}
}

func (c *BitbankWebsocketCallback) GetSubscribeTradesParam(symbol string) string {
	return fmt.Sprintf(`42["join-room", "transactions_%s"]`, symbol)
}

func (c *BitbankWebsocketCallback) OnReceiveTrade(msg string, errCh chan<- string) {
	log.Println(msg)
	if msg[0:2] == "42" {

		// exclude `["message",`
		body := msg[strings.Index(msg, keyword)+len(keyword) : len(msg)-1]

		var t Trade
		err := json.Unmarshal([]byte(body), &t)
		if err != nil {
			errCh <- fmt.Sprintf("failed to unmarchal json to bitbank.Trade: %v", err)
		}
		//log.Printf("bitbank(%v): %v\n", len(t.Message.Data.Transactions), t)

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

		for _, d := range t.Message.Data.Transactions {
			if d.TransactionId > lastTransactionId {
				c.writer.WriteString(fmt.Sprintf("%v,%v,%v,%.3f,%v,%v\n",
					now.Format(base.TimeFormat),
					d.TransactionId,
					d.Side,
					d.Price,
					d.Amount,
					d.ExecutedAt))
				lastTransactionId = d.TransactionId
			}
		}

	} else {
		log.Printf("bitbank: %v\n", msg)
	}
}
