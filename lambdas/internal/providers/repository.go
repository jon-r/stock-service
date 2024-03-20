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
		log.Fatalf("Missing settings for provider = %v", provider)
	}

	return settings
}

func FetchTickerDescription(provider ProviderName, ticker string) (error, *IndexDetails) {
	switch provider {
	case PolygonIo:
		return FetchPolygonTickerDescription(ticker)
	default:
		return fmt.Errorf("incorrect provider name: %v", provider), nil
	}
}
