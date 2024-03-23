package main

import (
	"jon-richards.com/stock-app/internal/providers"
)

func setTickerDescription(provider providers.ProviderName, tickerId string) error {
	var err error

	// 1. fetch the ticker details (based on the above)
	description, err := providers.FetchTickerDescription(provider, tickerId)

	if err != nil {
		return err
	}

	// 2. insert this ^ data into the ticker table
	err = dbService.SetTickerDescription(tickerId, description)

	return err
}

func setTickerHistoricalPrices(provider providers.ProviderName, tickerId string) error {
	var err error

	prices, err := providers.FetchTickerHistoricalPrices(provider, tickerId)

	if err != nil {
		return err
	}

	err = dbService.SetTickerHistoricalPrices(tickerId, prices)

	return err
}

func updateTickerPrices(provider providers.ProviderName, tickerIds []string) error {
	var err error

	prices, err := providers.FetchTickerDailyPrices(provider, tickerIds)

	if err != nil {
		return err
	}

	for tickerId, price := range *prices {
		err = dbService.UpdateTickerDailyPrices(tickerId, &price)
		if err != nil {
			break
		}
	}

	return err
}
