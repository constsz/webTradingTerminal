package system

import (
	"fmt"
	"sync"
	"ttk/data"
)

var DBJsonFilePath = "db.json"

type DB struct {
	mu   sync.Mutex
	data DbData
}

// ----------------------------------------
// DATA SCHEMAS

// DbData is what data is stored in DB
// - Symbols and BaseSize are persisted on drive,
// - TradeLog is in-memory and flushed after restart
type DbData struct {
	BaseSize    int
	Symbols     []SymbolSettings
	SymbolsList []string
	TradeLog    TradeLog
}

type SymbolSettings struct {
	Symbol string
	Mode   data.Mode
	Side   data.Side

	Entry     float64
	EntryType data.EntryType
	TP        float64
	SL        float64

	// STATUS
	PositionStatus   data.PositionStatus
	BaseQuantityUsdt float64
	NumOfEntries     int
	AvgPrice         string
	PnlUsd           float64
	PnlPct           float64

	LastPrice float64

	// AUTO-strategy settings
	TF        int
	ATRM      float64
	Execution data.Execution
	BotStatus data.BotStatus
}

type TradeLog []TradeLogEntry

type TradeLogEntry struct {
	Time   int64
	Symbol string
	Side   data.Side
	Event  data.Event
}

// ----------------------------------------
// DB METHODS

func NewDB() *DB {
	return &DB{
		data: DbData{
			BaseSize:    6,
			Symbols:     []SymbolSettings{},
			SymbolsList: []string{},
			TradeLog:    TradeLog{},
		},
	}
}

func (db *DB) GetBaseSize() int {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.data.BaseSize
}

func (db *DB) SetBaseSize(newSize int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data.BaseSize = newSize
	return nil
}

func (db *DB) GetSymbolSettings() []SymbolSettings {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.data.Symbols
}

func (db *DB) GetSymbolSideNonSync(symbol string) data.Side {
	for _, ss := range db.data.Symbols {
		if symbol == ss.Symbol {
			return ss.Side
		}
	}

	return 0

}

func (db *DB) GetSymbolSettingsIndex(symbol string) (int, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for i, ss := range db.data.Symbols {
		if ss.Symbol == symbol {
			return i, true
		}
	}

	return 0, false
}

func (db *DB) GetSymbolsList() []string {
	db.mu.Lock()
	defer db.mu.Unlock()

	return db.data.SymbolsList
}

func (db *DB) NewBlankSymbolSettings(symbol string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Add symbol to SymbolSettings
	db.data.Symbols = append(db.data.Symbols,
		SymbolSettings{
			Symbol: symbol,
			TP:     1,
			SL:     1,
		})

	// Add symbol to SymbolsList
	db.data.SymbolsList = append(db.data.SymbolsList, symbol)

}

func (db *DB) DeleteSymbolSettings(symbolName string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Delete symbol from SymbolSettings
	indexToDelete := 0

	for i, s := range db.data.Symbols {
		if symbolName == s.Symbol {
			indexToDelete = i
		}
	}

	db.data.Symbols[indexToDelete] = db.data.Symbols[len(db.data.Symbols)-1]
	db.data.Symbols = db.data.Symbols[:len(db.data.Symbols)-1]

	// Get index of slice element
	indexOfString := 0

	for i, s := range db.data.SymbolsList {
		if symbolName == s {
			indexOfString = i
		}
	}

	// Delete symbol from SymbolsList
	// Delete element by index
	db.data.SymbolsList[indexOfString] = db.data.SymbolsList[len(db.data.SymbolsList)-1]
	db.data.SymbolsList = db.data.SymbolsList[:len(db.data.SymbolsList)-1]

	fmt.Println("DELETING FROM SymbolsList:")
	fmt.Println("SLICE:", db.data.SymbolsList)
	fmt.Println("LEN:  ", len(db.data.SymbolsList))

}

func (db *DB) UpdateSymbolSettings(ss SymbolSettings) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for i, s := range db.data.Symbols {
		if ss.Symbol == s.Symbol {

			db.data.Symbols[i].Mode = ss.Mode
			db.data.Symbols[i].Side = ss.Side

			db.data.Symbols[i].Entry = ss.Entry
			db.data.Symbols[i].TP = ss.TP
			db.data.Symbols[i].SL = ss.SL

			db.data.Symbols[i].TF = ss.TF
			db.data.Symbols[i].ATRM = ss.ATRM
			db.data.Symbols[i].Execution = ss.Execution

			break
		}
	}

}

func (db *DB) GetSymbolsStatus() *data.AllSymbolsStatus {
	db.mu.Lock()
	defer db.mu.Unlock()

	numOfSymbols := len(db.data.Symbols)

	// Initialize AllSymbolsStatus
	allSymbolsStatus := make(data.AllSymbolsStatus, numOfSymbols)

	// Populate it with Status data
	for _, ss := range db.data.Symbols {
		allSymbolsStatus[ss.Symbol] = &data.SymbolStatus{
			PositionStatus: ss.PositionStatus,
			Side:           ss.Side,

			BaseQuantityUsdt: ss.BaseQuantityUsdt,
			NumOfEntries:     ss.NumOfEntries,
			AvgPrice:         ss.AvgPrice,

			PnlUsd: ss.PnlUsd,
			PnlPct: ss.PnlPct,
		}
	}

	return &allSymbolsStatus

}

func (db *DB) SetSymbolPositionStatus(symbol string, positionStatus data.PositionStatus, numberOfEntries int, avgPriceString string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for i := range db.data.Symbols {
		if db.data.Symbols[i].Symbol == symbol {

			// Set the PositionStatus
			db.data.Symbols[i].PositionStatus = positionStatus

			// Record changes of SymbolSettings according to
			// PositionStatus changes made in inventory.
			switch positionStatus {
			case data.StatusNotInPositionNoOpenOrders:
				db.data.Symbols[i].NumOfEntries = 0
				db.data.Symbols[i].AvgPrice = ""

			case data.StatusNotInPositionOrdersPlaced:

			case data.StatusInPositionNoOpenOrders:
				db.data.Symbols[i].NumOfEntries = numberOfEntries
				db.data.Symbols[i].AvgPrice = avgPriceString

			case data.StatusInPositionAddOrderPlaced:

			}

		}
	}

}

func (db *DB) SetSymbolBaseQuantityUsdt(symbol string, baseQuantityUsdt float64) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for i := range db.data.Symbols {
		if db.data.Symbols[i].Symbol == symbol {
			db.data.Symbols[i].BaseQuantityUsdt = baseQuantityUsdt
		}

	}

}

func (db *DB) SetLastPrice(symbol string, lastPrice float64) {
	db.mu.Lock()
	defer db.mu.Unlock()

	for i := range db.data.Symbols {
		if db.data.Symbols[i].Symbol == symbol {
			db.data.Symbols[i].LastPrice = lastPrice
		}
	}

}

func (db *DB) GetLastPrice(symbol string) float64 {
	db.mu.Lock()
	defer db.mu.Unlock()

	for i := range db.data.Symbols {
		if db.data.Symbols[i].Symbol == symbol {
			return db.data.Symbols[i].LastPrice
		}
	}

	return 0

}

func (db *DB) GetTradeLogs() {

}
