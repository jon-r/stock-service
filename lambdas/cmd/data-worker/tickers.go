package main

import "github.com/jon-r/stock-service/lambdas/internal/providers"

func (handler DataWorkerHandler) setTickerDescription(provider providers.ProviderName, tickerId string) error {
	var err error

	// 1. fetch the ticker details (based on the above)
	description, err := providers.FetchTickerDescription(provider, tickerId)

	if err != nil {
		return err
	}

	// 2. insert this ^ data into the ticker table
	err = handler.dbService.SetTickerDescription(handler.logService, tickerId, description)

	return err
}

func (handler DataWorkerHandler) setTickerHistoricalPrices(provider providers.ProviderName, tickerId string) error {
	var err error

	prices, err := providers.FetchTickerHistoricalPrices(provider, tickerId)

	if err != nil {
		return err
	}

	err = handler.dbService.AddTickerPrices(handler.logService, prices)

	return err
}

func (handler DataWorkerHandler) updateTickerPrices(provider providers.ProviderName, tickerIds []string) error {
	var err error

	prices, err := providers.FetchTickerDailyPrices(provider, tickerIds)

	if err != nil {
		return err
	}

	if prices == nil {
		handler.logService.Warnw("No prices for today",
			"provider", provider,
		)
		return nil
	}

	err = handler.dbService.AddTickerPrices(handler.logService, prices)

	return err
}
