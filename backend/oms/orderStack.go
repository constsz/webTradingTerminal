package oms

import (
	"fmt"
	"sync"
	"ttk/data"
)

type OrderStack struct {
	mu     sync.Mutex
	Symbol map[string]*OpenOrder
}

type OpenOrder struct {
	Active     bool
	Symbol     string
	BinanceId  int64
	Side       data.Side
	EntryPrice float64
	Quantity   float64
	TP         float64
	SL         float64
}

func newOrderStack(symbolsList []string) *OrderStack {
	orderStack := &OrderStack{
		Symbol: make(map[string]*OpenOrder, len(symbolsList)),
	}

	if len(symbolsList) > 0 {
		for _, s := range symbolsList {
			orderStack.Symbol[s] = &OpenOrder{}
		}
	}

	return orderStack
}

func (oStack *OrderStack) setOrder(oo OpenOrder) {
	oStack.mu.Lock()
	defer oStack.mu.Unlock()

	oStack.Symbol[oo.Symbol] = &oo
}

func (oStack *OrderStack) getOrder(symbol string) (*OpenOrder, bool) {
	oStack.mu.Lock()
	defer oStack.mu.Unlock()

	oo := oStack.Symbol[symbol]
	nilOO := &OpenOrder{}

	if oo != nil && oo != nilOO {
		return oStack.Symbol[symbol], true
	} else {
		fmt.Println("oo.getOrder: Open Order == nil")
		return oStack.Symbol[symbol], false
	}
}

func (oStack *OrderStack) removeOpenOrder(symbol string) {
	oStack.mu.Lock()
	defer oStack.mu.Unlock()

	_, exists := oStack.Symbol[symbol]

	if exists {
		oStack.Symbol[symbol] = &OpenOrder{}
	}

}

func (oStack *OrderStack) symbolHasOpenOrders(symbol string) bool {
	oStack.mu.Lock()
	defer oStack.mu.Unlock()

	hasOpenOrders := false

	openOrder, exists := oStack.Symbol[symbol]

	if exists {
		if openOrder.Active {
			hasOpenOrders = true
		}
	}
	//else {
	//	fmt.Println("symbolHasOpenOrders: (!) Symbol NOT FOUND!")
	//}

	return hasOpenOrders

}

func (oo *OpenOrder) Print() {
	/*type OpenOrder struct {
		Active       bool
		Symbol       string
		BinanceId    int64
		Side         data.Side
		EntryPrice   float64
		Quantity     float64
		TP           float64
		SL           float64
	}*/

	fmt.Println()
	fmt.Println("OpenOrder Print:")
	fmt.Println(oo.Active)
	fmt.Println(oo.Symbol)
	fmt.Println(oo.Side)
	fmt.Println(oo.Side)
}
