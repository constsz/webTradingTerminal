package data

// This is the data structure that API
// returns to frontend for "/symbolsStatus"
// (it's not stored in DB, as all fields are
// duplicates of SymbolSettings fields)

type AllSymbolsStatus map[string]*SymbolStatus

// DISPLAY:
// SIDE: Size     ∙  AvgPrice  ∙  PnL$  ∙  PnL%
// LONG: 3($200)  ∙  0.011312  ∙  $1.3  ∙  +0.54%

type SymbolStatus struct {
	PositionStatus PositionStatus
	Side           Side

	// inventory
	BaseQuantityUsdt float64
	NumOfEntries     int
	AvgPrice         string

	PnlUsd float64
	PnlPct float64
}
