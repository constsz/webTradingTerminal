package oms

import (
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"ttk/data"
	"ttk/system"
	"ttk/ta"
)

type OMS struct {
	db     *system.DB
	client *futures.Client

	chOmsCommands chan data.OMSCommand
	chUserData    chan *futures.WsUserDataEvent

	symbolRules *binanceSymbolRules
	inventory   *Inventory
	orderStack  *OrderStack
}

// MAIN

func Start(db *system.DB, client *futures.Client, chOmsCommands chan data.OMSCommand) {
	oms := newOMS(db, client, chOmsCommands)
	oms.Listen()
	fmt.Println("OMS started")
}

func newOMS(db *system.DB, client *futures.Client, chOmsCommands chan data.OMSCommand) *OMS {
	symbolsList := db.GetSymbolsList()

	symbolRules := NewBinanceSymbolRules(client)

	return &OMS{
		db:     db,
		client: client,

		chOmsCommands: chOmsCommands,
		chUserData:    make(chan *futures.WsUserDataEvent),

		symbolRules: symbolRules,
		inventory:   newInventory(symbolsList, 1),
		orderStack:  newOrderStack(symbolsList),
	}
}

func (oms *OMS) Listen() {
	go oms.userDataListener()

	go func() {
		for {
			select {
			case ude := <-oms.chUserData:
				oms.processUserDataEvent(ude)
			case command := <-oms.chOmsCommands:
				oms.processCommand(&command)
			}
		}
	}()

}

// METHODS

func (oms *OMS) processCommand(command *data.OMSCommand) {

	symbolName := command.OrderDetail.Symbol
	// Check if symbol is in Position
	inPosition := oms.inventory.symbolInPosition(symbolName)
	// Check if symbol has OpenOrders
	hasPostedOrders := oms.orderStack.symbolHasOpenOrders(symbolName)

	// Convert OMS Command to OpenOrder
	openOrder := oms.commandToOpenOrder(command, inPosition)

	switch command.Type {

	case data.SetOrder:
		// Command can be treated in different ways depending on situation:
		// 	- "Place Order": if symbol is not in position and has no OpenOrders, it's a simple place order.
		// 	- "Move/Adjust Order": if symbol is not in position, but already has OpenOrder - move/adjust order.
		//	- "Position Editing": if symbol is in position, SetOrder command is treated as editing
		fmt.Println("OMS: SetOrder")

		// if inPosition - this Command is Position Editing
		if inPosition {
			// Parse Command:
			// if EntryPrice != 0 : it means adding new order to add to current Position
			// 		in this case check if len(inv.Entries) < inv.PositionMaxTrades, if yes: proceed.
			// else: check if TP != Position.TP, if true: adjust TP
			// adjusting SL is not allowed

			// Case 1: New OpenOrder to add to currently active position
			if openOrder.EntryPrice > 0 {
				oms.CancelAllOrders(symbolName)
				orderPosted := oms.CreateEntryLimitOrder(openOrder)

				if orderPosted {
					oms.db.SetSymbolPositionStatus(symbolName, data.StatusInPositionAddOrderPlaced, 0, "")
				}

			} else if openOrder.EntryPrice == 0 {
				// Case 2: Adjust TP
				if openOrder.TP > 0 {
					// Set new TP in Position (inv)
					oms.inventory.updateSLTP(symbolName, openOrder.SL, openOrder.TP)

					// Cancel old SL/TP and set new orders on Binance
					oms.CancelAllOrders(symbolName)
					oms.CreateSlTpOrders(symbolName)

					position := oms.inventory.getPositionCopy(symbolName)
					numOfEntries := len(position.Entries)

					avgPriceString := oms.calcAverageEntryPriceString(symbolName)
					oms.db.SetSymbolPositionStatus(symbolName, data.StatusInPositionNoOpenOrders, numOfEntries, avgPriceString)

				}
			}
		} else if !inPosition {
			if openOrder.EntryPrice > 0 && openOrder.SL > 0 && openOrder.TP > 0 {
				if hasPostedOrders {
					// Cancel openOrders on Binance
					oms.CancelAllOrders(symbolName)
				}

				// Create new order on Binance
				orderPosted := oms.CreateEntryLimitOrder(openOrder)

				if orderPosted {
					oms.db.SetSymbolPositionStatus(symbolName, data.StatusNotInPositionOrdersPlaced, 0, "")
				}

			}

		}

	case data.CancelAllOrders:
		fmt.Println("OMS: CancelAllOrders")

		if !inPosition && hasPostedOrders {
			// Cancel all orders
			oms.CancelAllOrders(symbolName)
			oms.db.SetSymbolPositionStatus(symbolName, data.StatusNotInPositionNoOpenOrders, 0, "")
		} else if inPosition && hasPostedOrders {
			oms.CancelAllOrders(symbolName)
			oms.CreateSlTpOrders(symbolName)

			position := oms.inventory.getPositionCopy(symbolName)
			numOfEntries := len(position.Entries)

			avgPriceString := oms.calcAverageEntryPriceString(symbolName)
			oms.db.SetSymbolPositionStatus(symbolName, data.StatusInPositionNoOpenOrders, numOfEntries, avgPriceString)

		}

	case data.MarketEntry:
		fmt.Println("OMS: MarketEntry")
		if openOrder.Quantity > 0 {
			if hasPostedOrders {
				oms.CancelAllOrders(symbolName)
			}
			oms.MarketEntry(openOrder)
		} else {
			fmt.Println("(!) Market order QUANTITY == 0")
		}

	case data.MarketExit:
		fmt.Println("OMS: MarketExit")
		if inPosition {
			oms.MarketExit(openOrder)
		}

	}

}

func (oms *OMS) commandToOpenOrder(c *data.OMSCommand, inPosition bool) *OpenOrder {
	symbolName := c.OrderDetail.Symbol

	// Calculate EntryPrice if it's in Percents
	entryPrice := c.OrderDetail.Entry

	if c.OrderDetail.EntryType == data.EntryByPercent {
		// Get last price
		lastPrice := oms.db.GetLastPrice(c.OrderDetail.Symbol)

		// Multiply/divide by percentage (depending on the Side)
		pct := ta.Pct(entryPrice)

		switch c.OrderDetail.Side {
		case data.SideLong:
			entryPrice = lastPrice / pct
		case data.SideShort:
			entryPrice = lastPrice * pct
		}
	}

	qtyPrecision := (*oms.symbolRules)[symbolName].QuantityPrecision

	// Calculate position quantity
	baseSize := oms.db.GetBaseSize()

	// Calculate Quantity based on StopLoss size
	// or if inPosition use previously saved qtyUsdt

	var qtySymbol float64

	if !inPosition {
		//qtyUsdt := data.RoundToPrecision(float64(baseSize)/c.OrderDetail.SL, 1)
		qtyUsdt := data.RoundToPrecision(float64(baseSize), 1)

		oms.db.SetSymbolBaseQuantityUsdt(symbolName, qtyUsdt)
		qtySymbol = data.ConvertQuantityFromUSDT(qtyUsdt, entryPrice, qtyPrecision)
	} else {
		qtySymbol = oms.inventory.Symbol[symbolName].BaseQuantity
	}

	openOrder := &OpenOrder{
		Active:     true,
		Symbol:     symbolName,
		BinanceId:  0,
		Side:       c.OrderDetail.Side,
		EntryPrice: entryPrice,
		Quantity:   qtySymbol,
		TP:         c.OrderDetail.TP,
		SL:         c.OrderDetail.SL,
	}

	return openOrder
}
