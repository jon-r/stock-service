package main

import (
	"jon-richards.com/stock-app/internal/providers"
)

func newStock(provider providers.ProviderName, tickerId string) error {
	var err error

	// 1. fetch the ticker details (based on the above)
	//err, details := providers.GetStockIndexDetails(provider, tickerId)

	if err != nil {
		return err
	}

	// 2. insert this ^ data into the ticker table
	//err = dbService.NewStockItem(provider, tickerId, db.StockItemProperties{
	//	FullName: details.FullName,
	//	Currency: details.Currency,
	//})

	if err != nil {
		return err
	}

	// 3. return any new actions to the jobs queue, OR return a 'reattempt' OR a dead-letter-queue
	// todo

	return err
}
