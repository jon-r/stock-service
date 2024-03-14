package providers

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

var SettingsList = map[ProviderName]Settings{
	Slow: {Delay: 7, Url: "https://dog.ceo/api/breeds/image/random"},
	Fast: {Delay: 4, Url: "https://dog.ceo/api/breeds/image/random"},

	PolygonIo: {Delay: 12, Url: "https://api.polygon.io/v2/"},
	// todo url may not be needed for the settings.
	//  - config in general may be skippable since all providers are going to work differently anyway
}

type SettingsDelay int

const (
	PolygonIoDelay SettingsDelay = 12
)
