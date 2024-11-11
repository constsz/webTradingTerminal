package oms

import (
	"sync"
	"time"
	"ttk/data"
)

// TYPES

type Inventory struct {
	mu                sync.Mutex
	PositionMaxTrades int
	Symbol            map[string]*Position
}

type Position struct {
	Active    bool
	Side      data.Side
	EntryTime int64

	BaseQuantity     float64
	PositionQuantity float64

	Entries []float64
	TP      float64
	SL      float64

	BinanceIdEntry int64
	BinanceIdTP    int64
	BinanceIdSL    int64
}

// METHODS

// newInventory - for default capacity use value 0
func newInventory(symbolsList []string, positionMaxTrades int) *Inventory {
	defaultPositionMaxTrades := 1
	if positionMaxTrades == 0 {
		positionMaxTrades = defaultPositionMaxTrades
	}

	inventory := Inventory{
		PositionMaxTrades: positionMaxTrades,
		Symbol:            make(map[string]*Position, len(symbolsList)),
	}

	if len(symbolsList) > 0 {
		for _, s := range symbolsList {
			inventory.Symbol[s] = &Position{}
		}
	}

	return &inventory
}

func (inv *Inventory) symbolInPosition(symbol string) bool {
	inv.mu.Lock()
	defer inv.mu.Unlock()

	inPosition := false

	// Check if inventory of symbol exists
	position, exists := inv.Symbol[symbol]

	if exists {
		// If exists - check if it has position
		if position.Active {
			inPosition = true
		}
	}

	return inPosition

}

func (inv *Inventory) enterPosition(oo OpenOrder) Position {
	inPosition := inv.symbolInPosition(oo.Symbol)

	inv.mu.Lock()
	defer inv.mu.Unlock()

	// IF NOT IN POSITION
	if !inPosition {
		inv.Symbol[oo.Symbol] = &Position{
			Active:           true,
			Side:             oo.Side,
			EntryTime:        time.Now().UnixMilli(),
			BaseQuantity:     oo.Quantity,
			PositionQuantity: oo.Quantity,
			Entries:          []float64{oo.EntryPrice},
			TP:               oo.TP,
			SL:               oo.SL,
			BinanceIdEntry:   oo.BinanceId,
			BinanceIdTP:      0,
			BinanceIdSL:      0,
		}
	} else {
		// IF IN POSITION
		inv.Symbol[oo.Symbol].PositionQuantity = inv.Symbol[oo.Symbol].PositionQuantity + oo.Quantity
		inv.Symbol[oo.Symbol].Entries = append(inv.Symbol[oo.Symbol].Entries, oo.EntryPrice)
		inv.Symbol[oo.Symbol].TP = oo.TP
	}

	return *inv.Symbol[oo.Symbol]

}

func (inv *Inventory) getEntries(symbol string) []float64 {
	inv.mu.Lock()
	defer inv.mu.Unlock()

	return inv.Symbol[symbol].Entries
}

func (inv *Inventory) getPositionCopy(symbol string) Position {
	inv.mu.Lock()
	defer inv.mu.Unlock()

	return *inv.Symbol[symbol]
}

//func (inv *Inventory) calculatePositionStatusUpdates(symbol string) {
//	inv.mu.Lock()
//	defer inv.mu.Unlock()
//
//	pos := inv.Symbol[symbol]
//
//}

// RemovePosition is a Backtester function, can be ignored for now.
func (inv *Inventory) resetPosition(symbol string) {
	inv.mu.Lock()
	defer inv.mu.Unlock()

	inv.Symbol[symbol] = &Position{}
}

func (inv *Inventory) positionSide(symbol string) data.Side {
	inv.mu.Lock()
	defer inv.mu.Unlock()

	return inv.Symbol[symbol].Side
}

func (inv *Inventory) updateSLTP(symbol string, newSL, newTP float64) *Position {
	inv.mu.Lock()
	defer inv.mu.Unlock()

	inv.Symbol[symbol].SL = newSL
	inv.Symbol[symbol].TP = newTP

	return inv.Symbol[symbol]
}
