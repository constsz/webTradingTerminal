package system

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"log"
	"net/http"
	"strconv"
	"ttk/data"
)

func (app *App) symbolSettingsGet(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(app.DB.GetSymbolSettings())
	if err != nil {
		fmt.Println("API symbolSettingsGet:", err)
	}

	w.Write(jsonResp)
}

func (app *App) symbolSettingsAdd(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == http.MethodPost {
		req := make(map[string]string)
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Println("API symbolSettingsAdd:", err)
		}

		fmt.Println("ADDED new symbol: ", req["Symbol"])

		app.DB.NewBlankSymbolSettings(req["Symbol"])

		listOfSymbols := app.DB.GetSymbolsList()
		app.chPriceFeedSymbolList <- listOfSymbols

	}
}

func (app *App) symbolSettingsDelete(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == http.MethodPost {
		req := make(map[string]string)
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Println("API symbolSettingsDelete:", err)
		}

		// Cancel all orders
		// Create command for OMS
		omsCommand := data.OMSCommand{data.CancelAllOrders, data.OrderDetail{
			Symbol: req["Symbol"],
			SL:     0.3,
		}}

		// Send command to OMS
		app.chOmsCommands <- omsCommand

		fmt.Println("DELETED symbol: ", req["Symbol"])

		app.DB.DeleteSymbolSettings(req["Symbol"])

		listOfSymbols := app.DB.GetSymbolsList()
		app.chPriceFeedSymbolList <- listOfSymbols
	}

}

// ----------------------------------------
// ACTIONS

/*
	jsonReq = {
		ActionType: "SetOrder", (int)
		SymbolSettings: SymbolSettings,
	}

	function:
		save new SymbolSettings to DB

		switch by Mode --
			case "Manual":
				switch by ActionType --
					case SetOrder:
						send order "SetOrder" to OMS (symbol, side, entry, tp, sl)
					case CancelAllOrders:
						send order "CancelAllOrders" to OMS (symbol)
					case MarketEntry:
						send order "MarketEntry" to OMS (symbol, side)
					case MarketExit:
						send order "MarketExit" to OMS (symbol)

			case "Auto":
				switch by ActionType --
					case ...

*/

func (app *App) symbolAction(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == http.MethodPost {
		var s symbolActionReq
		err := json.NewDecoder(r.Body).Decode(&s)
		if err != nil {
			log.Println("API symbolAction json Decode:", err)
		}

		app.DB.UpdateSymbolSettings(s.SymbolSettings)

		switch s.SymbolSettings.Mode {
		case data.Manual:
			switch s.ActionType {
			case data.SetOrder:
				//fmt.Println("SymbolAction: Manual | SetOrder")
				side := s.SymbolSettings.Side
				entry := s.SymbolSettings.Entry

				// Create command for OMS
				omsCommand := data.OMSCommand{
					data.SetOrder,
					data.OrderDetail{
						Symbol:    s.SymbolSettings.Symbol,
						Side:      side,
						OrderType: data.Limit,
						Entry:     entry,
						EntryType: s.SymbolSettings.EntryType,
						TP:        s.SymbolSettings.TP,
						SL:        s.SymbolSettings.SL,
					},
				}

				fmt.Println("\n>> SetOrder")
				fmt.Println("Symbol:", omsCommand.OrderDetail.Symbol)
				fmt.Println("Side:", omsCommand.OrderDetail.Side)
				fmt.Println("Entry:", omsCommand.OrderDetail.Entry)
				fmt.Println("EntryType:", omsCommand.OrderDetail.EntryType)
				fmt.Println("SL:", omsCommand.OrderDetail.SL)
				fmt.Println("TP:", omsCommand.OrderDetail.TP)

				// Send command to OMS
				app.chOmsCommands <- omsCommand

			case data.CancelAllOrders:
				//fmt.Println("SymbolAction: Manual | CancelAllOrders")
				// Send command to OMS telling to cancel all orders

				// Create command for OMS
				omsCommand := data.OMSCommand{data.CancelAllOrders, data.OrderDetail{
					Symbol: s.SymbolSettings.Symbol,
					SL:     s.SymbolSettings.SL,
				}}

				fmt.Println("\n>> Cancel All Orders")
				fmt.Println("Symbol:", omsCommand.OrderDetail.Symbol)

				// Send command to OMS
				app.chOmsCommands <- omsCommand

			case data.MarketEntry:
				//fmt.Println("SymbolAction: Manual | MarketEntry")
				// Get last price to set correct SL and TP orders
				lastPrice := app.DB.GetLastPrice(s.SymbolSettings.Symbol)

				// Create command for OMS
				omsCommand := data.OMSCommand{
					data.MarketEntry,
					data.OrderDetail{
						Symbol:    s.SymbolSettings.Symbol,
						Side:      s.SymbolSettings.Side,
						OrderType: data.Market,
						Entry:     lastPrice,
						EntryType: 0,
						TP:        s.SymbolSettings.TP,
						SL:        s.SymbolSettings.SL,
					},
				}

				fmt.Println("\n>> Market Entry")
				fmt.Println("Symbol:", omsCommand.OrderDetail.Symbol)
				fmt.Println("Side:", omsCommand.OrderDetail.Side)
				fmt.Println("Entry:", omsCommand.OrderDetail.Entry)
				fmt.Println("SL:", omsCommand.OrderDetail.SL)
				fmt.Println("TP:", omsCommand.OrderDetail.TP)

				// Send command to OMS
				app.chOmsCommands <- omsCommand

			case data.MarketExit:
				//fmt.Println("SymbolAction: Manual | MarketExit")

				// Create command for OMS
				omsCommand := data.OMSCommand{data.MarketExit, data.OrderDetail{
					Symbol: s.SymbolSettings.Symbol,
					SL:     s.SymbolSettings.SL,
				}}

				fmt.Println("\n>> Market Exit")
				fmt.Println("Symbol:", omsCommand.OrderDetail.Symbol)

				// Send command to OMS
				app.chOmsCommands <- omsCommand

			}
		case data.Auto:
			//fmt.Println("Incoming symbol action | MODE: AUTO")
		}
	}
}

type symbolActionReq struct {
	ActionType     data.ActionType
	SymbolSettings SymbolSettings
}

// ----------------------------------------
// BASE SIZE & OTHER SETTINGS

func (app *App) baseSizeGet(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	baseSize := app.DB.GetBaseSize()

	w.Header().Set("Content-Type", "application/json")
	resp := map[string]int{
		"BaseSize": baseSize,
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Println("API GetBaseSize:", err)
	}

	w.Write(jsonResp)

}

func (app *App) baseSizeSet(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == http.MethodPost {
		bs := make(map[string]int)
		err := json.NewDecoder(r.Body).Decode(&bs)
		if err != nil {
			log.Println("API SetBaseSize:", err)
		}

		err = app.DB.SetBaseSize(bs["BaseSize"])

		if err != nil {
			log.Println("API -> DB SetBaseSize:", err)
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResp, err := json.Marshal(map[string]bool{
			"Result": true,
		})
		if err != nil {
			log.Println("API baseSizeSet-> json.Marshal:", err)
		}

		w.Write(jsonResp)
	}
}

func (app *App) symbolsStatus(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	// DISPLAY:
	// SIDE: Size     ∙  AvgPrice  ∙  PnL%  ∙  PnL$
	// LONG: 3($200)  ∙  0.011312  ∙  $1.3  ∙  +0.54%

	// Get status info from db
	allSymbolsStatus := app.DB.GetSymbolsStatus()

	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(allSymbolsStatus)
	if err != nil {
		fmt.Println("API symbolsStatus:", err)
	}

	w.Write(jsonResp)

}

func (app *App) tradeLogsGet(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	rows, err := strconv.Atoi(r.URL.Query().Get("rows"))
	if err != nil || rows < 0 {
		return
	}

	fmt.Fprintf(w, "GET: get Trade Logs for last %d rows", rows)

}

func (app *App) listFuturesSymbols(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	prices, err := app.Client.NewListPricesService().Do(context.Background())
	if err != nil {
		fmt.Println("API listFuturesSymbols:", err)
	}

	var futuresList []string
	for _, p := range prices {
		futuresList = append(futuresList, p.Symbol)
	}

	futuresList = data.SymbolNoUSDTList(futuresList)

	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(futuresList)
	if err != nil {
		log.Println("API listFuturesSymbols:", err)
	}

	w.Write(jsonResp)

}

// futuresListByVolume returns a list of futures sorted by volume and a TradingView addSymbol Prompt
// Can receive filter params:
// - minVolume | default = 0 (in USDT, filter out all symbols with smaller volume)
// - maxCount  | default = 0 (int, max number of symbols in a slice)
// if default: filter is ignored
func (app *App) futuresListByVolume(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	// TODO Backend:
	// - Read incoming request and get params: minVolume, maxCount
	// - Request Binance for the list of futures
	// - Create new slice with symbols, where volume is converted to USDT
	// - Sort this slice by volume
	// - Check filter params; if filterParam is not ignored (> 0) - filter out list
	// - Create object: {listOfSymbols, tradingViewPrompt}
	// - Marshall the object; pass to ResponseWriter

	// TODO Frontend:
	// - Button that copies tradingView Prompt into memory on click
	// - Create table with the symbols
	// - Add to menu

}

// tradingViewAddSymbolsPrompt Creates a long string used to add multiple symbols in TradingView
func tradingViewAddSymbolsPrompt(listOfSymbolsNames []string) string {
	var promptString string
	prefix := "BINANCE:"
	suffix := "USDT.P,"

	for _, symbol := range listOfSymbolsNames {
		newSymbolString := prefix + symbol + suffix
		promptString += newSymbolString
	}

	return ""
}

// App is a core element that holds DB (and maybe smth else in future),
// that should be shared across various interconnected sections (like API and DB)
type App struct {
	Port                  string
	DB                    *DB
	Client                *futures.Client
	chOmsCommands         chan<- data.OMSCommand
	chPriceFeedSymbolList chan<- []string
}

func StartAPI(db *DB, client *futures.Client,
	chOmsCommands chan<- data.OMSCommand, chPriceFeedSymbolList chan<- []string) {

	mux := http.NewServeMux()

	app := &App{
		Port:                  "7351",
		DB:                    db,
		Client:                client,
		chOmsCommands:         chOmsCommands,
		chPriceFeedSymbolList: chPriceFeedSymbolList,
	}

	mux.HandleFunc("/symbolSettings/get", app.symbolSettingsGet)       // GET
	mux.HandleFunc("/symbolSettings/add", app.symbolSettingsAdd)       // POST
	mux.HandleFunc("/symbolSettings/delete", app.symbolSettingsDelete) // POST
	mux.HandleFunc("/symbolAction", app.symbolAction)                  // POST
	mux.HandleFunc("/symbolsStatus", app.symbolsStatus)                // GET

	mux.HandleFunc("/baseSize/get", app.baseSizeGet) // GET
	mux.HandleFunc("/baseSize/set", app.baseSizeSet) // POST

	mux.HandleFunc("/tradeLogs/get", app.tradeLogsGet)            // GET
	mux.HandleFunc("/listFuturesSymbols", app.listFuturesSymbols) // GET

	fmt.Println("Starting server on :7351")

	go func() {
		err := http.ListenAndServe(":"+app.Port, mux)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST")
	(*w).Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
}
