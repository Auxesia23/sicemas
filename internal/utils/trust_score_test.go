package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sicemas/internal/dto"
)

func TestGetIPDifferenceLevel(t *testing.T) {
	tests := []struct {
		name       string
		current    string
		trusted    string
		expected   string
	}{
		{"same IP", "192.168.1.1", "192.168.1.1", "None"},
		{"minor change - last octet different", "192.168.1.10", "192.168.1.20", "Minor"},
		{"major change - third octet different", "192.168.1.1", "192.168.2.1", "Major"},
		{"total change - first two octets differ", "192.168.1.1", "169.168.1.1", "Total"},
		{"total change - all different", "10.0.0.1", "192.168.1.1", "Total"},
		{"total change - invalid IP", "not-an-ip", "192.168.1.1", "Total"},
		{"total change - both invalid", "abc", "def", "Total"},
		{"total change - empty trusted IP", "192.168.1.1", "", "Total"},
		{"total change - first octet same rest different", "10.0.0.1", "10.1.2.3", "Total"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getIPDifferenceLevel(tt.current, tt.trusted)
			assert.Equal(t, tt.expected, result)
		})
	}
}

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

	t.Run("different device ID deducts 10", func(t *testing.T) {
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
		assert.Equal(t, 90.0, score)
	})

	t.Run("different geo location deducts 20", func(t *testing.T) {
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
		assert.Equal(t, 80.0, score)
	})

	t.Run("minor IP change deducts 10 (WeightIPAddress * 0.25)", func(t *testing.T) {
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

	t.Run("major IP change deducts 24 (WeightIPAddress * 0.60)", func(t *testing.T) {
		current := &dto.SessionRequest{
			UserAgent:   "Mozilla/5.0",
			IPAddress:   "192.168.5.100",
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
		assert.Equal(t, 76.0, score)
	})

	t.Run("total IP change deducts 40 (WeightIPAddress * 1.00)", func(t *testing.T) {
		current := &dto.SessionRequest{
			UserAgent:   "Mozilla/5.0",
			IPAddress:   "10.0.0.1",
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
		assert.Equal(t, 60.0, score)
	})

	t.Run("all factors different - score = 100 - 40(IP) - 30(UA) - 20(Geo) - 10(Dev) = 0", func(t *testing.T) {
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
			UserAgent:   "FIREFOX", // case-insensitive match, same
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
		assert.Equal(t, 90.0, score)
	})
}
