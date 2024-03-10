package providers

import (
	"fmt"
	"log"
)

type IndexDetails struct {
	TickerId string
	Currency string
	FullName string
	// icon
}

func GetSettings(provider ProviderName) Settings {
	settings, ok := SettingsList[provider]

	if !ok {
		log.Printf("Missing settings for provider = %v", provider)
	}

	return settings
}

func GetStockIndexDetails(provider ProviderName, ticker string) (error, *IndexDetails) {
	switch provider {
	case PolygonIo:
		return GetPolygonIndexDetails(ticker)
	default:
		return fmt.Errorf("incorrect provider name: %v", provider), nil
	}
}
