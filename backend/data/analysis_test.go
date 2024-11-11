package data

import (
	"fmt"
	"testing"
)

func TestBinanceFuturesListByVolume(t *testing.T) {
	symbolsList := binanceFuturesListByVolume()

	if len(symbolsList) == 0 {
		fmt.Println("ERROR, empty list")
	}

	for _, symbol := range symbolsList {
		fmt.Println(symbol)
	}

}
