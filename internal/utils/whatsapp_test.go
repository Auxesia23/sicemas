package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendWhatsAppOTP(t *testing.T) {
	os.Setenv("GOWA_DEVICE_ID", "test-device-id")
	defer func() {
		os.Unsetenv("GOWA_DEVICE_ID")
		os.Unsetenv("GOWA_URL")
	}()

	t.Run("success - returns nil", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "POST", r.Method)
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.Equal(t, "test-device-id", r.Header.Get("X-Devic-Id"))

			user, pass, ok := r.BasicAuth()
			assert.True(t, ok)
			assert.Equal(t, "kuaci_emas", user)
			assert.Equal(t, "kuaci_emas", pass)

			body, _ := io.ReadAll(r.Body)
			var payload map[string]any
			err := json.Unmarshal(body, &payload)
			assert.NoError(t, err)
			assert.Equal(t, "+6281234567890@s.whatsapp.net", payload["phone"])
			assert.Contains(t, payload["message"].(string), "123456")
			assert.Equal(t, false, payload["is_forwarded"])

			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()

		os.Setenv("GOWA_URL", mockServer.URL)
		defer os.Unsetenv("GOWA_URL")

		err := SendWhatsAppOTP("081234567890", "123456")
		assert.NoError(t, err)
	})

	t.Run("failure - server returns non-200", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer mockServer.Close()

		os.Setenv("GOWA_URL", mockServer.URL)

		err := SendWhatsAppOTP("081234567890", "654321")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "500")
	})

	t.Run("failure - server returns bad request", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer mockServer.Close()

		os.Setenv("GOWA_URL", mockServer.URL)

		err := SendWhatsAppOTP("081234567890", "654321")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "400")
	})

	t.Run("failure - invalid URL causes error", func(t *testing.T) {
		os.Setenv("GOWA_URL", "http://127.0.0.1:1")

		err := SendWhatsAppOTP("081234567890", "654321")
		assert.Error(t, err)
	})

	t.Run("failure - empty GOWA_URL", func(t *testing.T) {
		os.Setenv("GOWA_URL", "")

		err := SendWhatsAppOTP("081234567890", "654321")
		assert.Error(t, err)
	})
}

func TestFormatPhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"starts with 62", "6281234567890", "+6281234567890"},
		{"starts with 0", "081234567890", "+6281234567890"},
		{"already has +62", "+6281234567890", "+6281234567890"},
		{"no prefix", "81234567890", "81234567890"},
		{"with spaces", "  081234567890  ", "+6281234567890"},
		{"empty string", "", ""},
		{"just 62", "62", "+62"},
		{"starts with 62 then 0", "62081234567890", "+62081234567890"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatPhoneNumber(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSendWhatsAppOTP_RequestPayload(t *testing.T) {
	os.Setenv("GOWA_DEVICE_ID", "test-device-id")
	os.Setenv("GOWA_URL", "http://example.com")
	defer func() {
		os.Unsetenv("GOWA_DEVICE_ID")
		os.Unsetenv("GOWA_URL")
	}()

	t.Run("formats phone with 62 prefix properly", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var payload map[string]any
			json.Unmarshal(body, &payload)

			assert.Equal(t, "+62@s.whatsapp.net", payload["phone"])
			assert.Contains(t, payload["message"].(string), "000000")

			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()

		os.Setenv("GOWA_URL", mockServer.URL)

		err := SendWhatsAppOTP("62", "000000")
		assert.NoError(t, err)
	})

	t.Run("formats phone with plus prefix", func(t *testing.T) {
		innerMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			var payload map[string]any
			json.Unmarshal(body, &payload)
			assert.Equal(t, "+621234567890@s.whatsapp.net", payload["phone"])
			w.WriteHeader(http.StatusOK)
		}))
		defer innerMock.Close()

		os.Setenv("GOWA_URL", innerMock.URL)
		err := SendWhatsAppOTP("+621234567890", "000000")
		assert.NoError(t, err)
	})
}
