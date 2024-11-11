package system

import (
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"log"
	"strconv"
	"time"
	"ttk/data"
	"ttk/ta"
)

func PnlUpdater(db *DB, client *futures.Client, chPriceFeedSymbolList <-chan []string) {
	listOfSymbols := db.GetSymbolsList()

	// Start WS AggTrade service
	// Event handler to send each new tick with aggTrade
	wsAllMiniTickerMarketHandler := func(event futures.WsAllMiniMarketTickerEvent) {
		if len(listOfSymbols) > 0 {
			for _, r := range event {
				// For all symbols in SymbolsSettings apply function to calculate PnLs
				if symbolExists(r.Symbol, listOfSymbols) {
					//t := time.Now()
					//fmt.Println(t.Format("15:4:5"))
					//fmt.Println(r.Symbol, r.ClosePrice)

					// Parse variables for symbol and lastPrice
					symbol := data.SymbolNoUSDT(r.Symbol)
					lastPrice, err := strconv.ParseFloat(r.ClosePrice, 64)
					if err != nil {
						fmt.Println("PnlUpdater: ParseFloat of r.ClosePrice: ", err)
					}

					// Save last price for symbol
					db.SetLastPrice(symbol, lastPrice)

					// Apply function to calculate PnL and save it to DB
					updateSymbolPnl(db, symbol, lastPrice)

				}
			}
		}

	}

	errHandler := func(err error) {
		log.Println(err)
	}

	doneC, _, err := futures.WsAllMiniMarketTickerServe(wsAllMiniTickerMarketHandler, errHandler)
	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		for {
			listOfSymbols = <-chPriceFeedSymbolList
			fmt.Println("PnlUpdater: Symbol change")
			fmt.Println(listOfSymbols)
		}
	}()

	// Save these info into

	<-doneC

}

func updateSymbolPnl(db *DB, symbol string, lastPrice float64) {
	i, symbolFound := db.GetSymbolSettingsIndex(symbol)

	if symbolFound {
		db.mu.Lock()
		defer db.mu.Unlock()

		if db.data.Symbols[i].AvgPrice != "" {
			avgPrice, err := strconv.ParseFloat(db.data.Symbols[i].AvgPrice, 64)
			if err != nil {
				fmt.Println("PnlUpdater | updateSymbolPnl : ParseFloat of avgPrice: ", err)
			}

			var pctChangeMult float64

			side := db.GetSymbolSideNonSync(symbol)

			// Calculate price percent change from AvgPrice
			if side == data.SideLong {
				pctChangeMult = ta.PctChangeMult(avgPrice, lastPrice)
			} else {
				pctChangeMult = ta.PctChangeMult(lastPrice, avgPrice)
			}

			// Save Price Pct Change to PnlPct
			pnlPct := ((pctChangeMult - 1) * 100) - 0.024
			db.data.Symbols[i].PnlPct = pnlPct

			// Calculate total USDT position size
			totalUsdtPositionSize := float64(db.data.Symbols[i].NumOfEntries) * db.data.Symbols[i].BaseQuantityUsdt

			// Calculate dollar difference: (UsdtPositionSize * PctChange) - UsdtPositionSize
			pnlUsd := (totalUsdtPositionSize * pctChangeMult) - totalUsdtPositionSize

			// Save it to PnlUsd
			db.data.Symbols[i].PnlUsd = pnlUsd

			//printPnlUpdate(pnlPct, pnlUsd)
		}

	} else {
		fmt.Println("(!) updateSymbolPnl : Symbol WAS NOT FOUND in DB")
	}

}

func symbolExists(symbolUsdt string, listOfSymbols []string) bool {
	for _, s := range listOfSymbols {
		if data.SymbolNoUSDT(symbolUsdt) == s {
			return true
		}
	}

	return false
}

func printPnlUpdate(pnlPct, pnlUsd float64) {
	fmt.Println()
	t := time.Now()
	fmt.Printf("%v  |  PNL: %.3f%v  |  $%.3f \n", t.Format("15:4:5"), pnlPct, "%", pnlUsd)
}
