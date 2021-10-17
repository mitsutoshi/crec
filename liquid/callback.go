package liquid

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mitsutoshi/crec/utils"
)

type LiquidWebsocketCallback struct {
	TradeFile *os.File
}

func (c *LiquidWebsocketCallback) OnReceiveTrade(t Trade) {
	fmt.Printf("liquid: %v\n", t)

	now := time.Now().UTC()
	names := strings.Split(c.TradeFile.Name(), ".")
	date := names[0][strings.LastIndex(names[0], "_")+1:]

	if date != now.Format("20060102") {

		// get new file name
		newName := fmt.Sprintf("%s_%s.csv", names[0][:strings.LastIndex(names[0], "_")], date)
		fmt.Printf("Close current file(%s) and open new file(%s).", c.TradeFile.Name(), newName)

		// close and open file
		c.TradeFile.Close()
		f, err := utils.OpenNewFile(newName)
		if err != nil {
			panic(err) // TODO
		}

		// replace file
		c.TradeFile = f
	}

	c.TradeFile.WriteString(fmt.Sprintf("%v,%s,%s,%.8f,%v,%v\n",
		now.Format(timeFormat),
		t.Symbol,
		t.Side,
		t.Size,
		t.Price,
		t.Timestamp.Format(timeFormat)))
}
