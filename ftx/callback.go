package ftx

import (
	"fmt"
	"os"
	"time"
)

type FtxWebsocketCallback struct {
	TradeFile *os.File
}

func (c *FtxWebsocketCallback) OnReceiveTrade(t Trade) {
	fmt.Printf("ftx(%v): %v\n", len(t.Data), t)

	if t.Type == msgTypeSubscribed {
		fmt.Println("channel has being subscribed.")
		return
	}

	now := time.Now().UTC().Format(timeFormat)
	for _, d := range t.Data {
		c.TradeFile.WriteString(fmt.Sprintf("%v,%v,%v,%.8f,%s,%v,%v\n",
			now,
			d.Id,
			d.Price,
			d.Size,
			d.Side,
			d.Liquidation,
			d.Time.Format(timeFormat)))
	}
}
