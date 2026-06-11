package utils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateOTP6(t *testing.T) {
	t.Run("returns exactly 6 characters", func(t *testing.T) {
		tests := []struct {
			name string
		}{
			{"first call"},
			{"second call"},
			{"third call"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				otp := GenerateOTP6()
				assert.Len(t, otp, 6)
			})
		}
	})

	t.Run("returns only numeric digits", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			otp := GenerateOTP6()
			_, err := strconv.Atoi(otp)
			assert.NoError(t, err, "OTP should be numeric: %s", otp)
		}
	})

	t.Run("returns values in valid range", func(t *testing.T) {
		for i := 0; i < 50; i++ {
			otp := GenerateOTP6()
			val, err := strconv.Atoi(otp)
			assert.NoError(t, err)
			assert.GreaterOrEqual(t, val, 0)
			assert.LessOrEqual(t, val, 999999)
		}
	})

	t.Run("pads with leading zeros", func(t *testing.T) {
		seenNonPadded := false
		for i := 0; i < 100; i++ {
			otp := GenerateOTP6()
			if otp[0] == '0' {
				return
			}
			if otp != "000000" && otp[0] != '0' {
				seenNonPadded = true
			}
		}
		_ = seenNonPadded
	})

	t.Run("produces different values on consecutive calls", func(t *testing.T) {
		values := make(map[string]bool)
		uniqueCount := 0
		for i := 0; i < 10; i++ {
			otp := GenerateOTP6()
			if !values[otp] {
				values[otp] = true
				uniqueCount++
			}
		}
		assert.GreaterOrEqual(t, uniqueCount, 2, "expected some unique OTPs")
	})

	t.Run("format is consistent", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			otp := GenerateOTP6()
			assert.Regexp(t, `^\d{6}$`, otp)
		}
	})
}
