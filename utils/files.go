package utils

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func NewFileName(exchange string, symbol string) string {
	dt := time.Now().Format("20060102")
	return fmt.Sprintf(
		"%s_%s_%s.csv", exchange, strings.ReplaceAll(symbol, "/", ""), dt)
}

func OpenNewFile(fileName string) (*os.File, error) {

	// make file name
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
