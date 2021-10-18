package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/mitsutoshi/crec/base"
	"github.com/mitsutoshi/crec/bybit"
	"github.com/mitsutoshi/crec/ftx"
	"github.com/mitsutoshi/crec/gmo"
	"github.com/mitsutoshi/crec/liquid"
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
	dt := time.Now().UTC().Format("20060102")
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

		// check and return error
		if coin.Exchange != "ftx" &&
			coin.Exchange != "bybit" &&
			coin.Exchange != "gmo" &&
			coin.Exchange != "liquid" {
			return fmt.Errorf("Unknow exchange: %v\n", coin.Exchange)
		}

		// create a file
		f, err := openSaveFile(coin.Exchange, coin.Symbol)
		if err != nil {
			return err
		}
		defer f.Close()

		// create writer
		writer := bufio.NewWriter(f)
		defer writer.Flush()

		var callback base.WebsocketCallback
		var wssUrl string
		var originUrl string
		var headers string

		if coin.Exchange == "ftx" {
			headers = ftx.TradeHeaders
			callback = ftx.NewWebsocketCallback(f, writer)
			wssUrl = ftx.WssUrl
			originUrl = ftx.WssUrl
		} else if coin.Exchange == "bybit" {
			headers = bybit.TradeHeaders
			callback = bybit.NewWebsocketCallback(f, writer)
			wssUrl = bybit.WssUrl
			originUrl = bybit.WssUrl
		} else if coin.Exchange == "gmo" {
			headers = gmo.TradeHeaders
			callback = gmo.NewWebsocketCallback(f, writer)
			wssUrl = gmo.WssUrl
			originUrl = gmo.OriginUrl
		} else if coin.Exchange == "liquid" {
			headers = liquid.TradeHeaders
		}

		// write file headers
		f.WriteString(headers + "\n")

		// connect exchange's websocket
		ws := base.Websocket{Callback: callback}
		if err := ws.Connect(wssUrl, originUrl); err != nil {
			return err
		}
		if err := ws.SubscribeTrades(coin.Symbol); err != nil {
			return err
		}

		// start to receive data
		go ws.Receive(errCh)
	}

	for {
		select {
		case msg := <-errCh:
			return fmt.Errorf("Received erorr message: %s", msg)
		case s := <-sig:
			fmt.Printf("Receved signal: %v\n", s)
			return nil
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
