package main

import (
	"go.uber.org/zap"
	"jon-richards.com/stock-app/internal/providers"
)

func setTickerDescription(log *zap.SugaredLogger, provider providers.ProviderName, tickerId string) error {
	var err error

	// 1. fetch the ticker details (based on the above)
	description, err := providers.FetchTickerDescription(provider, tickerId)

	if err != nil {
		return err
	}

	// 2. insert this ^ data into the ticker table
	err = dbService.SetTickerDescription(log, tickerId, description)

	return err
}

func setTickerHistoricalPrices(log *zap.SugaredLogger, provider providers.ProviderName, tickerId string) error {
	var err error

	prices, err := providers.FetchTickerHistoricalPrices(provider, tickerId)

	if err != nil {
		return err
	}

	err = dbService.AddTickerPrices(log, prices)

	return err
}

func updateTickerPrices(log *zap.SugaredLogger, provider providers.ProviderName, tickerIds []string) error {
	var err error

	prices, err := providers.FetchTickerDailyPrices(provider, tickerIds)

	if err != nil {
		return err
	}

	if prices == nil {
		log.Warnw("No prices for today",
			"provider", provider,
		)
		return nil
	}

	err = dbService.AddTickerPrices(log, prices)

	return err
}
