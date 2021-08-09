package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/mitsutoshi/bookkeeper/ftx"
)

const (
	configFileName = "config.yaml"
)

var (
	verOpt  = flag.Bool("v", false, "Show version.")
	version = "v0.0.1"
)

func main() {
	flag.Parse()

	// show version and exit
	if *verOpt {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	err := run()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}

func loadConfig(path string) (*config, error) {

	// open config file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// load config
	var c config
	if err = yaml.NewDecoder(f).Decode(&c); err != nil {
		return nil, err
	}
	return &c, err
}

func openSaveFile(exchange string, symbol string) (*os.File, error) {

	// make file name
	dt := time.Now().Format("20060102")
	fileName := fmt.Sprintf(
		"%s_%s_%s.csv", exchange, strings.ReplaceAll(symbol, "/", ""), dt)

	if _, err := os.Stat(fileName); err == nil {
		os.Remove(fileName)
	}

	// open file
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func run() error {

	// signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	c, err := loadConfig(configFileName)
	if err != nil {
		return err
	}

	errCh := make(chan string)
	for _, coin := range c.Coins {
		fmt.Println(coin)

		// ftx
		if coin.Exchange == "ftx" {

			f, err := openSaveFile(coin.Exchange, coin.Symbol)
			if err != nil {
				return err
			}
			defer f.Close()

			ws := ftx.Websocket{
				Callback: &ftx.FtxWebsocketCallback{TradeFile: f},
			}
			if err := ws.Connect(); err != nil {
				return err
			}
			if err := ws.SubscribeTrades(coin.Symbol); err != nil {
				return err
			}

			// write headers
			f.WriteString(ftx.TradeHeaders + "\n")

			// start receiving
			go ws.Receive(errCh)

		} else {
			return fmt.Errorf("Unknow exchange: %v\n", coin.Exchange)
		}
	}

	for {
		select {
		case msg := <-errCh:
			return fmt.Errorf("Received erorr message: %s", msg)
		case s := <-sig:
			fmt.Printf("Receved signal: %v\n", s)
			return nil
		default:
		}
	}
}

type Coin struct {
	Exchange string `yaml:"exchange"`
	Symbol   string `yaml:"symbol"`
}

type config struct {
	Coins []Coin `yaml:"coins"`
}
