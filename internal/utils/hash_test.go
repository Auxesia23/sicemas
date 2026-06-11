package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashFunctions(t *testing.T) {
	os.Setenv("PEPPER", "test-pepper-value-for-hashing")
	defer func() {
		os.Unsetenv("PEPPER")
	}()

	t.Run("HashIndex returns consistent value", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
		}{
			{"simple nip", "1234567890"},
			{"empty string", ""},
			{"with special chars", "abc!@#"},
			{"unicode", "测试"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r1 := HashIndex(tt.input)
				r2 := HashIndex(tt.input)
				assert.Equal(t, r1, r2)
				assert.NotEmpty(t, r1)
			})
		}
	})

	t.Run("HashIndex produces different hash for different inputs", func(t *testing.T) {
		h1 := HashIndex("input-a")
		h2 := HashIndex("input-b")
		assert.NotEqual(t, h1, h2)
	})

	t.Run("HashPassword returns non-empty hash", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
		}{
			{"simple password", "password123"},
			{"complex password", "P@ssw0rd!2024#Secure"},
			{"empty string", ""},
			{"short password", "ab"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				hash, err := HashPassword(tt.input)
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
			})
		}
	})

	t.Run("HashPassword returns different hashes for same input", func(t *testing.T) {
		h1, err := HashPassword("same-password")
		assert.NoError(t, err)
		h2, err := HashPassword("same-password")
		assert.NoError(t, err)
		assert.NotEqual(t, h1, h2)
	})

	t.Run("ComparePassword returns true for correct password", func(t *testing.T) {
		password := "my-secure-password"
		hash, err := HashPassword(password)
		assert.NoError(t, err)

		ok := ComparePassword(hash, password)
		assert.True(t, ok)
	})

	t.Run("ComparePassword returns false for wrong password", func(t *testing.T) {
		hash, err := HashPassword("correct-password")
		assert.NoError(t, err)

		ok := ComparePassword(hash, "wrong-password")
		assert.False(t, ok)
	})

	t.Run("ComparePassword returns false for empty hash", func(t *testing.T) {
		ok := ComparePassword("", "password")
		assert.False(t, ok)
	})

	t.Run("ComparePassword returns false for malformed hash", func(t *testing.T) {
		ok := ComparePassword("not-a-valid-bcrypt-hash", "password")
		assert.False(t, ok)
	})
}

// HashIndex calls log.Fatal (os.Exit) when PEPPER is unset, which cannot be
// recovered from in a test. We only test the happy path above.
