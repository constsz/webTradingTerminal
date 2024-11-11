package main

import (
	"fmt"
	"ttk/data"
	"ttk/oms"
	"ttk/system"

	"github.com/adshao/go-binance/v2"
)

// Launcher starts all elements of the system: DB, API, OMS
func Launcher() {
	forever := make(chan struct{})

	// Get ApiKey
	credentials, err := system.ReadJsonCredentials()
	if err != nil {
		fmt.Println("\n(!) Failed to load API Credentials! Check `credentials.json` file.")
		fmt.Println(err)
		return
	}

	client := binance.NewFuturesClient(credentials.ApiKey, credentials.SecretKey)

	db := system.NewDB()

	// Channels
	chOmsCommands := make(chan data.OMSCommand)
	chPriceFeedSymbolList := make(chan []string)

	// Serve API Entry points
	system.StartAPI(db, client, chOmsCommands, chPriceFeedSymbolList)

	// Send requests to DataProvider, Bot, OMS
	oms.Start(db, client, chOmsCommands)

	system.PnlUpdater(db, client, chPriceFeedSymbolList)

	<-forever
}

// ----------------------------------------
