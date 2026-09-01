package utils

import (
	"testing"

	"sicemas/internal/dto"

	"github.com/stretchr/testify/assert"
)

func TestCalculateTrustScore(t *testing.T) {
	t.Run("all factors match - score 100", func(t *testing.T) {
		current := &dto.SessionRequest{
			UserAgent:   "Mozilla/5.0",
			IPAddress:   "192.168.1.100",
			GeoLocation: "Jakarta",
			DeviceID:    "device-123",
		}
		trusted := &dto.SessionValue{
			UserAgent:   "Mozilla/5.0",
			IPAddress:   "192.168.1.100",
			GeoLocation: "Jakarta",
			DeviceID:    "device-123",
		}

		score := CalculateTrustScore(current, trusted)
		assert.Equal(t, 100.0, score)
	})

	t.Run("different user agent deducts 30", func(t *testing.T) {
		current := &dto.SessionRequest{
			UserAgent:   "Chrome",
			IPAddress:   "192.168.1.100",
			GeoLocation: "Jakarta",
			DeviceID:    "device-123",
		}
		trusted := &dto.SessionValue{
			UserAgent:   "Firefox",
			IPAddress:   "192.168.1.100",
			GeoLocation: "Jakarta",
			DeviceID:    "device-123",
		}

		score := CalculateTrustScore(current, trusted)
		assert.Equal(t, 70.0, score)
	})

	t.Run("different device ID deducts 20", func(t *testing.T) {
		current := &dto.SessionRequest{
			UserAgent:   "Mozilla/5.0",
			IPAddress:   "192.168.1.100",
			GeoLocation: "Jakarta",
			DeviceID:    "device-456",
		}
		trusted := &dto.SessionValue{
			UserAgent:   "Mozilla/5.0",
			IPAddress:   "192.168.1.100",
			GeoLocation: "Jakarta",
			DeviceID:    "device-123",
		}

		score := CalculateTrustScore(current, trusted)
		assert.Equal(t, 80.0, score)
	})

	t.Run("different geo location deducts 40", func(t *testing.T) {
		current := &dto.SessionRequest{
			UserAgent:   "Mozilla/5.0",
			IPAddress:   "192.168.1.100",
			GeoLocation: "Bandung",
			DeviceID:    "device-123",
		}
		trusted := &dto.SessionValue{
			UserAgent:   "Mozilla/5.0",
			IPAddress:   "192.168.1.100",
			GeoLocation: "Jakarta",
			DeviceID:    "device-123",
		}

		score := CalculateTrustScore(current, trusted)
		assert.Equal(t, 60.0, score)
	})

	t.Run("IP change deducts 10", func(t *testing.T) {
		current := &dto.SessionRequest{
			UserAgent:   "Mozilla/5.0",
			IPAddress:   "192.168.1.200",
			GeoLocation: "Jakarta",
			DeviceID:    "device-123",
		}
		trusted := &dto.SessionValue{
			UserAgent:   "Mozilla/5.0",
			IPAddress:   "192.168.1.100",
			GeoLocation: "Jakarta",
			DeviceID:    "device-123",
		}

		score := CalculateTrustScore(current, trusted)
		assert.Equal(t, 90.0, score)
	})

	t.Run("all factors different - score = 100 - 40(Geo) - 30(UA) - 20(Dev) - 10(IP) = 0", func(t *testing.T) {
		current := &dto.SessionRequest{
			UserAgent:   "Chrome",
			IPAddress:   "10.0.0.1",
			GeoLocation: "Bandung",
			DeviceID:    "device-456",
		}
		trusted := &dto.SessionValue{
			UserAgent:   "Firefox",
			IPAddress:   "192.168.1.100",
			GeoLocation: "Jakarta",
			DeviceID:    "device-123",
		}

		score := CalculateTrustScore(current, trusted)
		assert.Equal(t, 0.0, score)
	})

	t.Run("score does not go below 0", func(t *testing.T) {
		current := &dto.SessionRequest{
			UserAgent:   "Chrome",
			IPAddress:   "10.0.0.1",
			GeoLocation: "Bandung",
			DeviceID:    "device-456",
		}
		trusted := &dto.SessionValue{
			UserAgent:   "FIREFOX",
			IPAddress:   "192.168.1.100",
			GeoLocation: "Jakarta",
			DeviceID:    "device-123",
		}

		score := CalculateTrustScore(current, trusted)
		assert.GreaterOrEqual(t, score, 0.0)
	})

	t.Run("user agent comparison is case-insensitive", func(t *testing.T) {
		current := &dto.SessionRequest{
			UserAgent:   "MOZILLA/5.0",
			IPAddress:   "192.168.1.100",
			GeoLocation: "Jakarta",
			DeviceID:    "device-123",
		}
		trusted := &dto.SessionValue{
			UserAgent:   "mozilla/5.0",
			IPAddress:   "192.168.1.100",
			GeoLocation: "Jakarta",
			DeviceID:    "device-123",
		}

		score := CalculateTrustScore(current, trusted)
		assert.Equal(t, 100.0, score)
	})

	t.Run("device ID comparison is case-sensitive", func(t *testing.T) {
		current := &dto.SessionRequest{
			UserAgent:   "Mozilla/5.0",
			IPAddress:   "192.168.1.100",
			GeoLocation: "Jakarta",
			DeviceID:    "Device-123",
		}
		trusted := &dto.SessionValue{
			UserAgent:   "Mozilla/5.0",
			IPAddress:   "192.168.1.100",
			GeoLocation: "Jakarta",
			DeviceID:    "device-123",
		}

		score := CalculateTrustScore(current, trusted)
		assert.Equal(t, 80.0, score)
	})
}
