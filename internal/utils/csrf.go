package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GenerateCSRFToken(userId uuid.UUID) string {
	var csrfSecret = []byte(os.Getenv("CSRF_SECRET"))
	expiresAt := time.Now().Add(15 * time.Minute).Unix()

	payload := fmt.Sprintf("%s|%d", userId.String(), expiresAt)

	encodedPayload := base64.RawURLEncoding.EncodeToString([]byte(payload))

	mac := hmac.New(sha256.New, csrfSecret)
	mac.Write([]byte(encodedPayload))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return fmt.Sprintf("%s.%s", encodedPayload, signature)
}

func VerifyCSRFToken(token string, expectedUserId uuid.UUID) bool {
	var csrfSecret = []byte(os.Getenv("CSRF_SECRET"))
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return false
	}

	encodedPayload := parts[0]
	providedSignature := parts[1]

	mac := hmac.New(sha256.New, csrfSecret)
	mac.Write([]byte(encodedPayload))
	expectedSignature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(providedSignature), []byte(expectedSignature)) {
		return false
	}

	decodedPayload, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return false
	}

	payloadParts := strings.Split(string(decodedPayload), "|")
	if len(payloadParts) != 2 {
		return false
	}

	tokenIdStr := payloadParts[0]
	expiresAtStr := payloadParts[1]

	if tokenIdStr != expectedUserId.String() {
		return false
	}

	expiresAt, err := strconv.ParseInt(expiresAtStr, 10, 64)
	if err != nil {
		return false
	}

	if time.Now().Unix() > expiresAt {
		return false
	}

	return true
}
