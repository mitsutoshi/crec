package ftx

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mitsutoshi/crec/utils"
)

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

func (c *FtxWebsocketCallback) OnReceiveTrade(t Trade) {
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
			now.Format(timeFormat),
			d.Id,
			d.Price,
			d.Size,
			d.Side,
			d.Liquidation,
			d.Time.Format(timeFormat)))
	}
}
