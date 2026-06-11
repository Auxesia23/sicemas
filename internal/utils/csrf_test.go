package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateVerifyCSRFToken(t *testing.T) {
	os.Setenv("CSRF_SECRET", "test-csrf-secret-key-for-testing")
	defer func() {
		os.Unsetenv("CSRF_SECRET")
	}()

	userID := uuid.New()

	t.Run("Generate returns valid format", func(t *testing.T) {
		token := GenerateCSRFToken(userID)
		parts := strings.Split(token, ".")
		assert.Len(t, parts, 2)
		assert.NotEmpty(t, parts[0])
		assert.NotEmpty(t, parts[1])
	})

	t.Run("Verify passes for valid token", func(t *testing.T) {
		token := GenerateCSRFToken(userID)
		ok := VerifyCSRFToken(token, userID)
		assert.True(t, ok)
	})

	t.Run("Verify fails for wrong user", func(t *testing.T) {
		token := GenerateCSRFToken(userID)
		otherID := uuid.New()
		ok := VerifyCSRFToken(token, otherID)
		assert.False(t, ok)
	})

	t.Run("Verify fails for tampered payload", func(t *testing.T) {
		token := GenerateCSRFToken(userID)
		parts := strings.Split(token, ".")
		tampered := parts[0] + "a." + parts[1]
		ok := VerifyCSRFToken(tampered, userID)
		assert.False(t, ok)
	})

	t.Run("Verify fails for tampered signature", func(t *testing.T) {
		token := GenerateCSRFToken(userID)
		parts := strings.Split(token, ".")
		tampered := parts[0] + "." + parts[1][:len(parts[1])-1] + "X"
		ok := VerifyCSRFToken(tampered, userID)
		assert.False(t, ok)
	})

	t.Run("Verify fails for malformed token - no dot", func(t *testing.T) {
		ok := VerifyCSRFToken("malformed-no-dot", userID)
		assert.False(t, ok)
	})

	t.Run("Verify fails for malformed token - too many parts", func(t *testing.T) {
		ok := VerifyCSRFToken("part1.part2.part3", userID)
		assert.False(t, ok)
	})

	t.Run("Verify fails for expired token", func(t *testing.T) {
		var csrfSecret = []byte(os.Getenv("CSRF_SECRET"))
		expiresAt := time.Now().Add(-1 * time.Minute).Unix()
		payload := userID.String() + "|" + func() string {
			return strings.TrimRight(strings.Replace(
				time.Unix(expiresAt, 0).String(), "-", "", -1), "0")
		}()
		_ = payload

		payload = userID.String() + "|" + "0"
		encodedPayload := strings.TrimRight(
			func() string {
				mac := hmac.New(sha256.New, csrfSecret)
				mac.Write([]byte(base64.RawURLEncoding.EncodeToString([]byte(payload))))
				return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
			}(), "=")

		_ = encodedPayload

		token := func() string {
			encoded := base64.RawURLEncoding.EncodeToString([]byte(payload))
			mac := hmac.New(sha256.New, csrfSecret)
			mac.Write([]byte(encoded))
			sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
			return encoded + "." + sig
		}()

		ok := VerifyCSRFToken(token, userID)
		assert.False(t, ok)
	})

	t.Run("Verify fails for empty token", func(t *testing.T) {
		ok := VerifyCSRFToken("", userID)
		assert.False(t, ok)
	})

	t.Run("Verify fails with tampered payload content", func(t *testing.T) {
		token := GenerateCSRFToken(userID)
		parts := strings.Split(token, ".")
		oldPayload, _ := base64.RawURLEncoding.DecodeString(parts[0])

		newPayload := uuid.New().String() + "|" + strings.Split(string(oldPayload), "|")[1]
		newEncoded := base64.RawURLEncoding.EncodeToString([]byte(newPayload))

		tampered := newEncoded + "." + parts[1]
		ok := VerifyCSRFToken(tampered, userID)
		assert.False(t, ok)
	})

	t.Run("Generate produces same token deterministically in same second", func(t *testing.T) {
		t1 := GenerateCSRFToken(userID)
		t2 := GenerateCSRFToken(userID)
		assert.Equal(t, t1, t2)
	})
}
