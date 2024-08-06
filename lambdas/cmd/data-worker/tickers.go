package main

import "github.com/jon-r/stock-service/lambdas/internal/providers"

func (handler DataWorkerHandler) setTickerDescription(provider providers.ProviderName, tickerId string) error {
	var err error

	// 1. fetch the ticker details (based on the above)
	description, err := handler.ProviderService.FetchTickerDescription(provider, tickerId)

	if err != nil {
		return err
	}

	// 2. insert this ^ data into the ticker table
	err = handler.DbService.SetTickerDescription(handler.LogService, tickerId, description)

	return err
}

func (handler DataWorkerHandler) setTickerHistoricalPrices(provider providers.ProviderName, tickerId string) error {
	var err error

	prices, err := handler.ProviderService.FetchTickerHistoricalPrices(provider, tickerId)

	if err != nil {
		return err
	}

	err = handler.DbService.AddTickerPrices(handler.LogService, prices)

	return err
}

func (handler DataWorkerHandler) updateTickerPrices(provider providers.ProviderName, tickerIds []string) error {
	var err error

	prices, err := handler.ProviderService.FetchTickerDailyPrices(provider, tickerIds)

	if err != nil {
		return err
	}

	if prices == nil {
		handler.LogService.Warnw("No prices for today",
			"provider", provider,
		)
		return nil
	}

	err = handler.DbService.AddTickerPrices(handler.LogService, prices)

	return err
}
