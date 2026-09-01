package utils

import (
	"sicemas/internal/dto"
	"strings"
)

const (
	WeightGeoLocation float64 = 40.0
	WeightUserAgent   float64 = 30.0
	WeightDeviceID    float64 = 20.0
	WeightIPAddress   float64 = 10.0
)

func CalculateTrustScore(current *dto.SessionRequest, trusted *dto.SessionValue) float64 {
	score := 100.0

	if current.IPAddress != trusted.IPAddress {
		score -= WeightIPAddress
	}

	if !strings.EqualFold(current.UserAgent, trusted.UserAgent) {
		score -= WeightUserAgent
	}

	if current.DeviceID != trusted.DeviceID {
		score -= WeightDeviceID
	}

	if current.GeoLocation != trusted.GeoLocation {
		score -= WeightGeoLocation
	}

	if score < 0.0 {
		return 0.0
	}

	return score
}
