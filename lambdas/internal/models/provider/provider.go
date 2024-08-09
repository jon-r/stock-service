package provider

import (
	"fmt"
	"time"
)

type Name string

const (
	PolygonIo Name = "POLYGON_IO"
)

func GetRequestsPerMin(providerName Name) (time.Duration, error) {
	switch providerName {
	case PolygonIo:
		return time.Minute / 5, nil
	default:
		return time.Hour, fmt.Errorf("incorrect provider name: %v", providerName)
	}
}
