package main

//
//func (handlers DataWorkerHandler) setTickerDescription(provider providers_old.ProviderName, tickerId string) error {
//	var err error
//
//	// 1. fetch the ticker details (based on the above)
//	description, err := handlers.ProviderService.FetchTickerDescription(provider, tickerId)
//
//	if err != nil {
//		return err
//	}
//
//	// 2. insert this ^ data into the ticker table
//	err = handlers.DbService.SetTickerDescription(handlers.LogService, tickerId, description)
//
//	return err
//}
//
//func (handlers DataWorkerHandler) setTickerHistoricalPrices(provider providers_old.ProviderName, tickerId string) error {
//	var err error
//
//	prices, err := handlers.ProviderService.FetchTickerHistoricalPrices(provider, tickerId)
//
//	if err != nil {
//		return err
//	}
//
//	err = handlers.DbService.AddTickerPrices(handlers.LogService, prices)
//
//	return err
//}
//
//func (handlers DataWorkerHandler) updateTickerPrices(provider providers_old.ProviderName, tickerIds []string) error {
//	var err error
//
//	prices, err := handlers.ProviderService.FetchTickerDailyPrices(provider, tickerIds)
//
//	if err != nil {
//		return err
//	}
//
//	if prices == nil {
//		handlers.LogService.Warnw("No prices for today",
//			"provider", provider,
//		)
//		return nil
//	}
//
//	err = handlers.DbService.AddTickerPrices(handlers.LogService, prices)
//
//	return err
//}
