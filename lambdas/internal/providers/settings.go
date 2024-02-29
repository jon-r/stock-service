package providers

import "log"

type ProviderName string

const (
	Slow ProviderName = "SLOW"
	Fast ProviderName = "FAST"
)

type Settings struct {
	Delay int32
	Url   string
}

var SettingsList = map[ProviderName]Settings{
	Slow: {Delay: 7, Url: "https://dog.ceo/api/breeds/image/random"},
	Fast: {Delay: 4, Url: "https://dog.ceo/api/breeds/image/random"},
}

func GetSettings(provider ProviderName) Settings {
	settings, ok := SettingsList[provider]

	if !ok {
		log.Printf("Missing settings for provider = %v", provider)
	}

	return settings
}
