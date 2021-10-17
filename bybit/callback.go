package bybit

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mitsutoshi/crec/utils"
)

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

func (c *BybitWebsocketCallback) OnReceiveTrade(t Trade) {
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
			now.Format(timeFormat),
			d.Timestamp.Format(timeFormat),
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
