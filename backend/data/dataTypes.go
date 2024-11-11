package data

// Candles
// HeikenAshi
// Side (long, short)
// Event (Entry, TP, SL)

type Side int

const (
	SideLong Side = iota
	SideShort
)

type Mode int

const (
	Manual Mode = iota
	Auto
)

type PositionStatus int

const (
	StatusNotInPositionNoOpenOrders PositionStatus = iota
	StatusNotInPositionOrdersPlaced                // 1
	StatusInPositionNoOpenOrders                   // 2
	StatusInPositionAddOrderPlaced                 // 3
)

type ActionType int

const (
	SetOrder ActionType = iota
	CancelAllOrders
	MarketEntry
	MarketExit
)

type OMSCommand struct {
	Type        ActionType
	OrderDetail OrderDetail
}

type OrderDetail struct {
	Symbol    string
	Side      Side
	OrderType OrderType
	Entry     float64
	EntryType EntryType
	TP        float64
	SL        float64
}

type OrderType int

const (
	Limit OrderType = iota
	Market
)

func SideInvert(side Side) Side {
	if side == SideLong {
		return SideShort
	} else {
		return SideLong
	}
}

type EntryType int

const (
	EntryByPrice EntryType = iota
	EntryByPercent
)

type Event int

const (
	EventEntry Event = iota
	EventTP
	EventSL
)

type BotStatus int

const (
	BotStandby BotStatus = iota
	BotWorking
)

type Execution int

const (
	ExecOnce Execution = iota
	ExecLoop
)
