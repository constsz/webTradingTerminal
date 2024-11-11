package oms

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"log"
	"math"
	"strconv"
	"time"
	"ttk/data"
)

func (oms *OMS) userDataListener() {
	// Get ListenKey from Binance to use for WsUserData
	listenKey, err := oms.client.NewStartUserStreamService().Do(context.Background())
	if err != nil {
		fmt.Println("OMS userDataListener:", err)
		return
	}

	// WS User Data Listener
	errHandler := func(err error) {
		log.Println("WsUserDataService: ", err)
	}

	wsUserHandler := func(event *futures.WsUserDataEvent) {
		oms.chUserData <- event
	}

	doneC, _, err := futures.WsUserDataServe(listenKey, wsUserHandler, errHandler)
	if err != nil {
		fmt.Println("(!) futures.WsUserDataServe:", err)
		return
	}

	tickerKeepAlive := time.NewTicker(10 * time.Minute)

	go func() {
		for {
			<-tickerKeepAlive.C
			//fmt.Println("<><><><> KeepAlive UserStreamService")
			oms.client.NewKeepaliveUserStreamService().ListenKey(listenKey).Do(context.Background())
		}
	}()

	<-doneC
}

// To get last price use either upd.LastFilledPrice or upd.AveragePrice
// Checking steps:
// 		- lastFilledPrice > 0 && originalQty == lastFilledQty
// 			- if !inPosition: new position
// 			- if inPosition:
//				- if position.Side == upd.Side: add to position
//				- if sides are opposite: exit position

func (oms *OMS) processUserDataEvent(ude *futures.WsUserDataEvent) {

	// Check if Event type is OrderTradeUpdate
	if ude.Event == futures.UserDataEventTypeAccountUpdate {
		fmt.Println("\n--- --- ---\nBINANCE POSITIONS:")
		for _, wsPosition := range ude.AccountUpdate.Positions {

			fmt.Println(" - ", wsPosition.Symbol, "| SIDE:", wsPosition.Side, "| AMOUNT:", wsPosition.Amount)

			symbol := data.SymbolNoUSDT(wsPosition.Symbol)

			var correctSymbol bool
			for _, s := range oms.db.GetSymbolsList() {
				if symbol == s {
					correctSymbol = true
				}
			}

			if correctSymbol {
				inPosition := oms.inventory.symbolInPosition(symbol)
				oo, openOrderExists := oms.orderStack.getOrder(symbol) // oo can be nil like OpenOrder{}
				wsPositionAmount, _ := strconv.ParseFloat(wsPosition.Amount, 64)
				wsPositionAmountAbs := math.Abs(wsPositionAmount)

				var wsPositionSide data.Side
				if wsPositionAmount > 0 {
					wsPositionSide = data.SideLong
				} else if wsPositionAmount < 0 {
					wsPositionSide = data.SideShort
				}

				if !inPosition {
					if openOrderExists {
						if oo.Quantity == wsPositionAmountAbs {
							fmt.Println()
							fmt.Println(symbol, " | ENTER new Position")

							position := oms.inventory.enterPosition(*oo)
							numOfEntries := len(position.Entries)

							avgPriceString := oms.calcAverageEntryPriceString(symbol)
							oms.db.SetSymbolPositionStatus(symbol, data.StatusInPositionNoOpenOrders, numOfEntries, avgPriceString)

							oms.CreateSlTpOrders(symbol)
							oms.orderStack.removeOpenOrder(symbol)

						}
					} else {
						fmt.Println("OMS ERROR: WS Position update, but no OpenOrder exists!")
					}
				} else if inPosition {
					fmt.Println()
					fmt.Printf("%s | IN POSITION", symbol)
					position := oms.inventory.getPositionCopy(symbol)

					// ENTRY: ADD TO EXISTING POSITION
					if wsPositionAmount > 0 || wsPositionAmount < 0 {
						fmt.Printf(" | ADD to existing position\n")
						if position.Side == wsPositionSide {
							fmt.Println("--- good: position sides match ---")

							qtyPrecision := (*oms.symbolRules)[symbol].QuantityPrecision

							fullQty := data.RoundToPrecision(oo.Quantity+position.PositionQuantity, qtyPrecision)

							fmt.Println("--- | Comparing QTY: OO:", oo.Quantity, "INV:", position.PositionQuantity, "FULL:", fullQty, "WS:", wsPositionAmountAbs)

							if wsPositionAmountAbs == fullQty {
								fmt.Println("--- good: position QTY match ---")

								newPosition := oms.inventory.enterPosition(*oo)
								numOfEntries := len(newPosition.Entries)
								avgPriceString := oms.calcAverageEntryPriceString(symbol)
								oms.db.SetSymbolPositionStatus(symbol, data.StatusInPositionNoOpenOrders, numOfEntries, avgPriceString)

								oms.CreateSlTpOrders(symbol)
								oms.orderStack.removeOpenOrder(symbol)

							}

						} else {
							fmt.Println("(!) OMS userDataListener ERRROR:")
							fmt.Println("Binance Position SIDE and Inventory position side NOT MATCH!")
						}

					} else if wsPositionAmount == 0 {
						fmt.Printf(" | EXIT position\n")

						oms.CancelAllOrders(symbol)
						oms.inventory.resetPosition(symbol)
						oms.db.SetSymbolPositionStatus(symbol, data.StatusNotInPositionNoOpenOrders, 0, "")

					}

					// EXIT POSITION
				}

			}

		}
	}

	alwaysFalse := false
	if alwaysFalse && ude.Event == futures.UserDataEventTypeOrderTradeUpdate {
		upd := ude.OrderTradeUpdate
		symbol := data.SymbolNoUSDT(upd.Symbol)

		var correctSymbol bool
		for _, s := range oms.db.GetSymbolsList() {
			if symbol == s {
				correctSymbol = true
			}
		}

		if correctSymbol {
			oo, _ := oms.orderStack.getOrder(symbol)

			originalQty, _ := strconv.ParseFloat(upd.OriginalQty, 64)
			lastFilledQty, _ := strconv.ParseFloat(upd.LastFilledQty, 64)
			//lastFilledPrice, _ := strconv.ParseFloat(upd.LastFilledPrice, 64)

			inPosition := oms.inventory.symbolInPosition(symbol)

			// Check if WHOLE order was finally filled
			if /*lastFilledPrice > 0 && */ originalQty == lastFilledQty {

				oo.EntryPrice, _ = strconv.ParseFloat(upd.LastFilledPrice, 64)

				if !inPosition {
					// ENTER NEW POSITION
					fmt.Println()
					fmt.Println(symbol, " | ENTER new Position")
					printBinanceUpdate(&upd)

					position := oms.inventory.enterPosition(*oo)
					numOfEntries := len(position.Entries)

					avgPriceString := oms.calcAverageEntryPriceString(symbol)
					oms.db.SetSymbolPositionStatus(symbol, data.StatusInPositionNoOpenOrders, numOfEntries, avgPriceString)
					/*
						To enter the position:
						- inventory
						- db records:  set PositionStatus
						- db from inv: BaseQuantityUsdt, NumOfEntries, AvgPrice

					*/

					oms.CreateSlTpOrders(symbol)
					oms.orderStack.removeOpenOrder(symbol)
				} else {
					// IF IN POSITION
					fmt.Println()
					fmt.Println(symbol, " | inPosition")
					positionSide := oms.inventory.positionSide(symbol)
					binanceUpdateSide := binanceSideToInternal(upd.Side)

					if positionSide == binanceUpdateSide {
						// ADD TO EXISTING POSITION
						fmt.Println("---- ADD to Position")
						printBinanceUpdate(&upd)

						position := oms.inventory.enterPosition(*oo)
						numOfEntries := len(position.Entries)
						avgPriceString := oms.calcAverageEntryPriceString(symbol)
						oms.db.SetSymbolPositionStatus(symbol, data.StatusInPositionNoOpenOrders, numOfEntries, avgPriceString)

						oms.CreateSlTpOrders(symbol)
						oms.orderStack.removeOpenOrder(symbol)
					} else {
						// EXIT POSITION
						fmt.Println("---- EXIT")
						printBinanceUpdate(&upd)

						oms.CancelAllOrders(symbol)
						oms.inventory.resetPosition(symbol)
						oms.db.SetSymbolPositionStatus(symbol, data.StatusNotInPositionNoOpenOrders, 0, "")

					}
				}

			}

		}
	}

}

func (oms *OMS) calcAverageEntryPriceString(symbol string) string {
	entries := oms.inventory.getEntries(symbol)
	avgPriceFloat := data.Average(entries)

	avgPrice := BinancePrice(avgPriceFloat, oms.symbolRules.GetTickSize(symbol))

	return avgPrice
}

func printBinanceUpdate(upd *futures.WsOrderTradeUpdate) {
	//fmt.Println("TradeID:        ", upd.TradeID)
	t := time.Now()
	fmt.Println(t.Format("Mon 15:4:5"))
	fmt.Println("Type:           ", upd.Type)
	fmt.Println("Side:           ", upd.Side)
	fmt.Println("LastFilledPrice:", upd.LastFilledPrice)
	fmt.Println("OriginalQty:    ", upd.OriginalQty)
	fmt.Println("LastFilledQty:  ", upd.LastFilledQty)
	fmt.Println()
}
