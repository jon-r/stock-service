package provider

import (
	"time"
)

type Name string

const (
	PolygonIo Name = "POLYGON_IO"
)

func GetRequestsPerMin() map[Name]time.Duration {
	return map[Name]time.Duration{
		PolygonIo: time.Minute / 5,
	}
}
