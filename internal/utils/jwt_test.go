package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"sicemas/internal/dto"
	"sicemas/internal/entity"
)

func TestJWTFunctions(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-jwt-secret-key-minimum-32-chars!")
	defer func() {
		os.Unsetenv("JWT_SECRET")
	}()

	user := &entity.User{
		ID:          uuid.New(),
		NamaLengkap: "Test User",
	}
	jti := uuid.New()
	sid := uuid.New()

	t.Run("GenerateAccessToken and ParseAccessToken roundtrip", func(t *testing.T) {
		tokenString, err := GenerateAccessToken(user, jti, sid)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		claims, err := ParseAccessToken(tokenString)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, user.NamaLengkap, claims.NamaLengkap)
		assert.Equal(t, sid.String(), claims.SID)
		assert.Equal(t, user.ID.String(), claims.Subject)
		assert.Equal(t, jti.String(), claims.ID)
	})

	t.Run("GenerateRefreshToken and ParseRefreshToken roundtrip", func(t *testing.T) {
		tokenString, err := GenerateRefreshToken(user, jti)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		claims, err := ParseRefreshToken(tokenString)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, user.ID.String(), claims.Subject)
		assert.Equal(t, jti.String(), claims.ID)
	})

	t.Run("ParseAccessToken fails with empty string", func(t *testing.T) {
		_, err := ParseAccessToken("")
		assert.Error(t, err)
	})

	t.Run("ParseAccessToken fails with malformed token", func(t *testing.T) {
		_, err := ParseAccessToken("this.is.not.a.valid.jwt")
		assert.Error(t, err)
	})

	t.Run("ParseAccessToken fails with wrong signing method", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodNone, dto.AccessToken{
			NamaLengkap: "test",
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: user.ID.String(),
			},
		})
		tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
		assert.NoError(t, err)

		_, err = ParseAccessToken(tokenString)
		assert.Error(t, err)
	})

	t.Run("ParseAccessToken fails with expired token", func(t *testing.T) {
		claims := dto.AccessToken{
			NamaLengkap: user.NamaLengkap,
			SID:         sid.String(),
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   user.ID.String(),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
				Issuer:    "KUA Ciemas",
				ID:        jti.String(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		assert.NoError(t, err)

		_, err = ParseAccessToken(tokenString)
		assert.Error(t, err)
	})

	t.Run("ParseAccessToken fails with wrong secret", func(t *testing.T) {
		tokenString, err := GenerateAccessToken(user, jti, sid)
		assert.NoError(t, err)

		os.Setenv("JWT_SECRET", "different-secret-key-that-is-32-bytes!!")
		defer os.Setenv("JWT_SECRET", "test-jwt-secret-key-minimum-32-chars!")

		_, err = ParseAccessToken(tokenString)
		assert.Error(t, err)
	})

	t.Run("ParseRefreshToken fails with empty string", func(t *testing.T) {
		_, err := ParseRefreshToken("")
		assert.Error(t, err)
	})

	t.Run("ParseRefreshToken fails with malformed token", func(t *testing.T) {
		_, err := ParseRefreshToken("invalid-token-format")
		assert.Error(t, err)
	})

	t.Run("ParseRefreshToken fails with wrong signing method", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.RegisteredClaims{
			Subject: user.ID.String(),
		})
		tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
		assert.NoError(t, err)

		_, err = ParseRefreshToken(tokenString)
		assert.Error(t, err)
	})

	t.Run("ParseRefreshToken fails with wrong secret", func(t *testing.T) {
		tokenString, err := GenerateRefreshToken(user, jti)
		assert.NoError(t, err)

		os.Setenv("JWT_SECRET", "different-secret-key-that-is-32-bytes!!")
		defer os.Setenv("JWT_SECRET", "test-jwt-secret-key-minimum-32-chars!")

		_, err = ParseRefreshToken(tokenString)
		assert.Error(t, err)
	})
}
