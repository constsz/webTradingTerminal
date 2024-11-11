package oms

import (
	"context"
	"github.com/adshao/go-binance/v2/futures"
	"log"
	"ttk/data"
)

type binanceSymbolRules map[string]*binanceSymbolRulesDetail

type binanceSymbolRulesDetail struct {
	TickSize          string
	QuantityPrecision int
}

// NewBinanceSymbolRules is launched one time at the start of this app.
// It stores rules for all futures available on Binance.
func NewBinanceSymbolRules(client *futures.Client) *binanceSymbolRules {
	// client.ExchangeInfoService ... чтобы получить все Precision и TickSize.
	res, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Initialize main variable
	bsr := make(binanceSymbolRules, len(res.Symbols))

	for _, s := range res.Symbols {
		bsr[data.SymbolNoUSDT(s.Symbol)] = &binanceSymbolRulesDetail{
			TickSize:          s.PriceFilter().TickSize,
			QuantityPrecision: s.QuantityPrecision,
		}
	}

	return &bsr
}

func (bsr *binanceSymbolRules) GetTickSize(symbol string) string {
	return (*bsr)[symbol].TickSize
}
