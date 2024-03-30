package main

import (
	"context"

	"jon-richards.com/stock-app/internal/logging"
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
	err = dbService.SetTickerDescription(tickerId, *description)

	return err
}

func setTickerHistoricalPrices(provider providers.ProviderName, tickerId string) error {
	var err error

	prices, err := providers.FetchTickerHistoricalPrices(provider, tickerId)

	if err != nil {
		return err
	}

	err = dbService.SetTickerHistoricalPrices(tickerId, *prices)

	return err
}

func updateTickerPrices(ctx context.Context, provider providers.ProviderName, tickerIds []string) error {
	// todo add logger to all these worker functions. maybe pass logger around instead of context?
	log := logging.NewLogger(ctx)
	defer log.Sync()

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

	// todo delete
	log.Infow("Have prices",
		"prices", prices,
	)

	for tickerId, price := range *prices {
		var input interface{}
		err, input = dbService.UpdateTickerDailyPrices(tickerId, []providers.TickerPrices{price})

		// todo delete
		log.Infow("Input Check",
			"input", input,
		)

		if err != nil {
			break
		}
	}

	return err
}
