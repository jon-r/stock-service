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

	err = dbService.SetTickerHistoricalPrices(log, tickerId, prices)

	return err
}

//func updateTickerPrices(ctx context.Context, provider providers.ProviderName, tickerIds []string) error {
//	// todo CPT-95 add logger to all these worker functions. maybe pass logger around instead of context?
//	log := logging.NewLogger(ctx)
//	defer log.Sync()
//
//	var err error
//
//	prices, err := providers.FetchTickerDailyPrices(provider, tickerIds)
//
//	if err != nil {
//		return err
//	}
//
//	if prices == nil {
//		log.Warnw("No prices for today",
//			"provider", provider,
//		)
//		return nil
//	}
//
//	for tickerId, price := range *prices {
//		err = dbService.UpdateTickerDailyPrices(tickerId, []providers.TickerPrices{price})
//
//		if err != nil {
//			break
//		}
//	}
//
//	return err
//}
