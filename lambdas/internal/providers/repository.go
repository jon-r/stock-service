package providers

import (
	"fmt"
	"log"
)

type ProviderName string

const (
	Slow ProviderName = "SLOW"
	Fast ProviderName = "FAST"

	PolygonIo ProviderName = "POLYGON_IO"
)

type Settings struct {
	Delay int32
	Url   string
}

type IndexDetails struct {
	TickerId string
	Currency string
	FullName string
	// icon
}

var SettingsList = map[ProviderName]Settings{
	Slow: {Delay: 7, Url: "https://dog.ceo/api/breeds/image/random"},
	Fast: {Delay: 4, Url: "https://dog.ceo/api/breeds/image/random"},

	PolygonIo: {Delay: 12, Url: "https://api.polygon.io/v2/"},
	// todo url may not be needed for the settings.
	//  - config in general may be skippable since all providers are going to work differently anyway
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
