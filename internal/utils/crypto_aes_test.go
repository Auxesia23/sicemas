package utils

import (
	"encoding/hex"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	os.Setenv("AES_256_KEY", "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=")
	defer func() {
		os.Unsetenv("AES_256_KEY")
	}()

	t.Run("EncryptDecrypt - roundtrip succeeds", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
		}{
			{"plain text", "Hello, World!"},
			{"with special chars", "abc123!@#$%^&*()_+"},
			{"unicode", "你好世界"},
			{"long text", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."},
			{"numeric", "1234567890"},
			{"single char", "a"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				encrypted, err := Encrypt(tt.input)
				assert.NoError(t, err)
				assert.NotEmpty(t, encrypted)

				decrypted, err := Decrypt(encrypted)
				assert.NoError(t, err)
				assert.Equal(t, tt.input, decrypted)
			})
		}
	})

	t.Run("Encrypt returns empty for empty input", func(t *testing.T) {
		result, err := Encrypt("")
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Decrypt returns empty for empty input", func(t *testing.T) {
		result, err := Decrypt("")
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Encrypt produces unique ciphertext each time", func(t *testing.T) {
		r1, err := Encrypt("same-value")
		assert.NoError(t, err)
		r2, err := Encrypt("same-value")
		assert.NoError(t, err)
		assert.NotEqual(t, r1, r2)
	})

	t.Run("Decrypt fails with malformed hex", func(t *testing.T) {
		_, err := Decrypt("not-hex-string")
		assert.Error(t, err)
	})

	t.Run("Decrypt fails with short ciphertext (>= nonce but garbage)", func(t *testing.T) {
		_, err := Decrypt("abcdef0123456789abcdef0123456789")
		assert.Error(t, err)
	})

	t.Run("Decrypt fails with tampered ciphertext", func(t *testing.T) {
		encrypted, err := Encrypt("secret-data")
		assert.NoError(t, err)

		ct, _ := hex.DecodeString(encrypted)
		ct[len(ct)-1] ^= 0xFF
		tampered := hex.EncodeToString(ct)

		_, err = Decrypt(tampered)
		assert.Error(t, err)
	})

	t.Run("Encrypt fails with invalid AES key", func(t *testing.T) {
		os.Setenv("AES_256_KEY", "aW52YWxpZA==")
		defer os.Setenv("AES_256_KEY", "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=")

		_, err := Encrypt("data")
		assert.Error(t, err)
	})

	t.Run("Decrypt fails with invalid AES key", func(t *testing.T) {
		encrypted, err := Encrypt("data")
		assert.NoError(t, err)

		os.Setenv("AES_256_KEY", "aW52YWxpZA==")
		defer os.Setenv("AES_256_KEY", "MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTI=")

		_, err = Decrypt(encrypted)
		assert.Error(t, err)
	})
}
