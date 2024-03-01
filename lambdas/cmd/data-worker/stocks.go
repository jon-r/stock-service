package main

import (
	"jon-richards.com/stock-app/db"
	"jon-richards.com/stock-app/providers"
)

func newStock(provider providers.ProviderName, tickerId string) error {
	var err error

	// 1. fetch the ticker details (based on the above)
	err, details := providers.GetStockIndexDetails(provider, tickerId)

	if err != nil {
		return err
	}

	// 2. insert this ^ data into the stockIndex table
	err = dbService.NewStockItem(provider, tickerId, db.StockItemProperties{
		FullName: details.FullName,
		Currency: details.Currency,
	})

	if err != nil {
		return err
	}

	// 3. insert a 'prepopulate history' action to the jobs table
	// todo

	return err
}
