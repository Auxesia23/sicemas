package middlewares

import (
	"encoding/base64"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"sicemas/internal/dto"
	"sicemas/internal/geoip"
	"sicemas/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func freshAuthMiddleware(t *testing.T) (*MockCache, *MockLocator, *MockAuthService, AuthMiddleware) {
	t.Helper()

	dummyKeyBytes := []byte("12345678901234567890123456789012")
	os.Setenv("AES_256_KEY", base64.StdEncoding.EncodeToString(dummyKeyBytes))
	os.Setenv("JWT_SECRET", "test-jwt-secret-key-minimum-32-chars!")
	os.Setenv("CSRF_SECRET", "test-csrf-secret-for-testing-12345")

	mockCache := new(MockCache)
	mockLocator := new(MockLocator)
	mockAuthSvc := new(MockAuthService)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	mw := NewAuthMiddleware(mockAuthSvc, nil, mockLocator, mockCache, logger)

	return mockCache, mockLocator, mockAuthSvc, mw
}

func clearAuthEnv() {
	os.Unsetenv("AES_256_KEY")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("CSRF_SECRET")
}

// ==========================================
// JWTAuthenticator Tests
// ==========================================

func TestJWTAuthenticator(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-jwt-secret-key-minimum-32-chars!")
	defer os.Unsetenv("JWT_SECRET")

	user := testUserEntity()
	jti := uuid.New()
	sid := uuid.New()

	validToken, err := utils.GenerateAccessToken(user, jti, sid)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		setupMock      func(*MockCache)
		expectedStatus int
	}{
		{
			name: "success - valid token passes",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: validToken})
				return req
			},
			setupMock: func(mc *MockCache) {
				mc.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "failure - no cookie returns 401",
			setupRequest: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/test", nil)
			},
			setupMock:      func(mc *MockCache) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "failure - empty cookie returns 401",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: ""})
				return req
			},
			setupMock:      func(mc *MockCache) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "failure - blocked token returns 401",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: validToken})
				return req
			},
			setupMock: func(mc *MockCache) {
				mc.On("Get", mock.Anything, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						ptr := args.Get(2).(*bool)
						*ptr = true
					}).Return(nil).Once()
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "failure - invalid token returns 401",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: "this-is-not-a-valid-jwt"})
				return req
			},
			setupMock: func(mc *MockCache) {
				mc.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(redis.Nil).Once()
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "failure - expired token returns 401",
			setupRequest: func() *http.Request {
				expClaims := dto.AccessToken{
					NamaLengkap: "Test User",
					SID:         sid.String(),
					RegisteredClaims: jwt.RegisteredClaims{
						Subject:   user.ID.String(),
						IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
						Issuer:    "KUA Ciemas",
						ID:        jti.String(),
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, expClaims)
				expToken, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: expToken})
				return req
			},
			setupMock: func(mc *MockCache) {
				mc.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(redis.Nil).Once()
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache, _, _, mw := freshAuthMiddleware(t)

			tt.setupMock(mockCache)

			app := fiber.New()
			app.Get("/test", mw.JWTAuthenticator, func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			})

			req := tt.setupRequest()
			resp, err := app.Test(req, 1000)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockCache.AssertExpectations(t)
		})
	}
}

// ==========================================
// ZeroTrustValidator Tests
// ==========================================

type zeroTrustCase struct {
	name           string
	setupRequest   func(string) *http.Request
	setupMock      func(*MockCache, *MockLocator, *MockAuthService, uuid.UUID, uuid.UUID, string)
	expectedStatus int
}

func TestZeroTrustValidator(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-jwt-secret-key-minimum-32-chars!")
	defer os.Unsetenv("JWT_SECRET")

	user := testUserEntity()
	jti := uuid.New()
	sid := uuid.New()

	validToken, err := utils.GenerateAccessToken(user, jti, sid)
	assert.NoError(t, err)

	tests := []zeroTrustCase{
		{
			name: "success - high trust score passes",
			setupRequest: func(token string) *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: token})
				req.Header.Set("User-Agent", "test-agent")
				req.Header.Set("X-Device-Id", "device-123")
				return req
			},
			setupMock: func(mc *MockCache, ml *MockLocator, ma *MockAuthService, uid uuid.UUID, sId uuid.UUID, token string) {
				ml.On("Lookup", mock.Anything).Return(&geoip.Location{City: "Jakarta"}, nil).Once()
				mc.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
				mc.On("Get", mock.Anything, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						ptr := args.Get(2).(*dto.SessionValue)
						*ptr = dto.SessionValue{
							UserID:        uid,
							SID:           sId,
							UserAgent:     "test-agent",
							IPAddress:     "0.0.0.0",
							GeoLocation:   "Jakarta",
							DeviceID:      "device-123",
							IsMFAVerified: true,
						}
					}).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "failure - session not found (redis.Nil) returns 401",
			setupRequest: func(token string) *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: token})
				req.Header.Set("User-Agent", "test-agent")
				return req
			},
			setupMock: func(mc *MockCache, ml *MockLocator, ma *MockAuthService, uid uuid.UUID, sId uuid.UUID, token string) {
				ml.On("Lookup", mock.Anything).Return(&geoip.Location{City: "Jakarta"}, nil).Once()
				mc.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
				mc.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(redis.Nil).Once()
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "failure - geoip lookup error returns 500",
			setupRequest: func(token string) *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: token})
				return req
			},
			setupMock: func(mc *MockCache, ml *MockLocator, ma *MockAuthService, uid uuid.UUID, sId uuid.UUID, token string) {
				mc.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
				ml.On("Lookup", mock.Anything).Return(nil, fmt.Errorf("geoip database error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "failure - MFA not verified returns 403",
			setupRequest: func(token string) *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: token})
				req.Header.Set("User-Agent", "test-agent")
				req.Header.Set("X-Device-Id", "device-123")
				return req
			},
			setupMock: func(mc *MockCache, ml *MockLocator, ma *MockAuthService, uid uuid.UUID, sId uuid.UUID, token string) {
				ml.On("Lookup", mock.Anything).Return(&geoip.Location{City: "Jakarta"}, nil).Once()
				mc.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
				mc.On("Get", mock.Anything, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						ptr := args.Get(2).(*dto.SessionValue)
						*ptr = dto.SessionValue{
							UserID:        uid,
							SID:           sId,
							UserAgent:     "test-agent",
							IPAddress:     "0.0.0.0",
							GeoLocation:   "Jakarta",
							DeviceID:      "device-123",
							IsMFAVerified: false,
						}
					}).Return(nil).Once()
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "failure - medium trust score triggers step up and returns 403",
			setupRequest: func(token string) *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: token})
				req.Header.Set("User-Agent", "different-agent")
				req.Header.Set("X-Device-Id", "device-123")
				return req
			},
			setupMock: func(mc *MockCache, ml *MockLocator, ma *MockAuthService, uid uuid.UUID, sId uuid.UUID, token string) {
				ml.On("Lookup", mock.Anything).Return(&geoip.Location{City: "Jakarta"}, nil).Once()
				mc.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
				mc.On("Get", mock.Anything, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						ptr := args.Get(2).(*dto.SessionValue)
						*ptr = dto.SessionValue{
							UserID:        uid,
							SID:           sId,
							UserAgent:     "test-agent",
							IPAddress:     "0.0.0.0",
							GeoLocation:   "Jakarta",
							DeviceID:      "device-123",
							IsMFAVerified: true,
						}
					}).Return(nil).Once()
				mc.On("Set", mock.Anything, mock.Anything, mock.Anything, 7*24*time.Hour).Return(nil).Once()
				ma.On("TriggerStepUpOTP", mock.Anything, uid.String()).Return(nil).Once()
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "failure - low trust score invalidates session and returns 401",
			setupRequest: func(token string) *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: token})
				req.Header.Set("User-Agent", "totally-different-agent")
				req.Header.Set("X-Device-Id", "totally-different-device")
				return req
			},
			setupMock: func(mc *MockCache, ml *MockLocator, ma *MockAuthService, uid uuid.UUID, sId uuid.UUID, token string) {
				ml.On("Lookup", mock.Anything).Return(&geoip.Location{City: "Surabaya"}, nil).Once()
				mc.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
				mc.On("Get", mock.Anything, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						ptr := args.Get(2).(*dto.SessionValue)
						*ptr = dto.SessionValue{
							UserID:        uid,
							SID:           sId,
							UserAgent:     "test-agent",
							IPAddress:     "10.0.0.1",
							GeoLocation:   "Jakarta",
							DeviceID:      "device-123",
							IsMFAVerified: true,
						}
					}).Return(nil).Once()
				mc.On("Delete", mock.Anything, mock.Anything).Return(nil).Once()
				mc.On("Set", mock.Anything, mock.Anything, true, 15*time.Minute).Return(nil).Once()
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache, mockLocator, mockAuthSvc, mw := freshAuthMiddleware(t)

			tt.setupMock(mockCache, mockLocator, mockAuthSvc, user.ID, sid, validToken)

			app := fiber.New()
			app.Get("/test", mw.JWTAuthenticator, mw.ZeroTrustValidator, func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			})

			req := tt.setupRequest(validToken)
			resp, err := app.Test(req, 1000)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockCache.AssertExpectations(t)
			mockLocator.AssertExpectations(t)
			mockAuthSvc.AssertExpectations(t)
		})
	}
}

func TestZeroTrustValidator_InvalidIP(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-jwt-secret-key-minimum-32-chars!")
	defer os.Unsetenv("JWT_SECRET")

	user := testUserEntity()
	jti := uuid.New()
	sid := uuid.New()

	validToken, err := utils.GenerateAccessToken(user, jti, sid)
	assert.NoError(t, err)

	mockCache, mockLocator, _, mw := freshAuthMiddleware(t)

	mockCache.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	mockLocator.On("Lookup", mock.Anything).Return(&geoip.Location{City: "Jakarta"}, nil).Once()
		mockCache.On("Get", mock.Anything, mock.Anything, mock.Anything).
			Run(func(args mock.Arguments) {
				ptr := args.Get(2).(*dto.SessionValue)
				*ptr = dto.SessionValue{
					UserID:        user.ID,
					SID:           sid,
					UserAgent:     "test-agent",
					IPAddress:     "0.0.0.0",
					GeoLocation:   "Jakarta",
					DeviceID:      "device-123",
					IsMFAVerified: true,
				}
			}).Return(nil).Once()

	app := fiber.New(fiber.Config{
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"0.0.0.0/0"},
	})
	app.Get("/test", mw.JWTAuthenticator, mw.ZeroTrustValidator, func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: validToken})
	req.Header.Set("X-Forwarded-For", "10.0.0.1")
	req.Header.Set("User-Agent", "test-agent")
	req.Header.Set("X-Device-Id", "device-123")

	resp, err := app.Test(req, 1000)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mockCache.AssertExpectations(t)
	mockLocator.AssertExpectations(t)
}

// ==========================================
// GetContext Tests
// ==========================================

func TestGetContext(t *testing.T) {
	t.Run("success - sets context in locals", func(t *testing.T) {
		_, mockLocator, _, mw := freshAuthMiddleware(t)

		mockLocator.On("Lookup", mock.Anything).Return(&geoip.Location{City: "Jakarta"}, nil).Once()

		app := fiber.New()
		var capturedContext any
		app.Get("/test", mw.GetContext, func(c *fiber.Ctx) error {
			capturedContext = c.Locals("context")
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("User-Agent", "test-agent")
		req.Header.Set("X-Device-Id", "device-456")

		resp, err := app.Test(req, 1000)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		sessionReq, ok := capturedContext.(*dto.SessionRequest)
		assert.True(t, ok)
		assert.NotNil(t, sessionReq)
		assert.Equal(t, "test-agent", sessionReq.UserAgent)
		assert.Equal(t, "device-456", sessionReq.DeviceID)
		assert.Equal(t, "Jakarta", sessionReq.GeoLocation)
		assert.NotEmpty(t, sessionReq.IPAddress)

		mockLocator.AssertExpectations(t)
	})

	t.Run("success - sets context with defaults when headers missing", func(t *testing.T) {
		_, mockLocator, _, mw := freshAuthMiddleware(t)

		mockLocator.On("Lookup", mock.Anything).Return(&geoip.Location{City: "Jakarta"}, nil).Once()

		app := fiber.New()
		var capturedContext any
		app.Get("/test", mw.GetContext, func(c *fiber.Ctx) error {
			capturedContext = c.Locals("context")
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		resp, err := app.Test(req, 1000)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		sessionReq, ok := capturedContext.(*dto.SessionRequest)
		assert.True(t, ok)
		assert.NotNil(t, sessionReq)

		mockLocator.AssertExpectations(t)
	})
}

// ==========================================
// CSRFDoubleSubmit Tests
// ==========================================

func TestCSRFDoubleSubmit(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-jwt-secret-key-minimum-32-chars!")
	os.Setenv("CSRF_SECRET", "test-csrf-secret-for-testing-12345")
	defer func() {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("CSRF_SECRET")
	}()

	user := testUserEntity()
	jti := uuid.New()
	sid := uuid.New()

	validToken, err := utils.GenerateAccessToken(user, jti, sid)
	assert.NoError(t, err)

	csrfToken := utils.GenerateCSRFToken(user.ID)

	tests := []struct {
		name           string
		method         string
		setupRequest   func() *http.Request
		expectedStatus int
	}{
		{
			name:   "success - GET request passes through",
			method: http.MethodGet,
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: validToken})
				return req
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "success - HEAD request passes through",
			method: http.MethodHead,
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodHead, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: validToken})
				return req
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "success - OPTIONS request passes through",
			method: http.MethodOptions,
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodOptions, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: validToken})
				return req
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "success - POST with matching CSRF tokens passes",
			method: http.MethodPost,
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: validToken})
				req.AddCookie(&http.Cookie{Name: "csrf_token", Value: csrfToken})
				req.Header.Set("X-CSRF-Token", csrfToken)
				return req
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "failure - POST missing cookie token returns 403",
			method: http.MethodPost,
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: validToken})
				req.Header.Set("X-CSRF-Token", csrfToken)
				return req
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:   "failure - POST missing header token returns 403",
			method: http.MethodPost,
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: validToken})
				req.AddCookie(&http.Cookie{Name: "csrf_token", Value: csrfToken})
				return req
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:   "failure - POST mismatched tokens returns 403",
			method: http.MethodPost,
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: validToken})
				req.AddCookie(&http.Cookie{Name: "csrf_token", Value: csrfToken})
				req.Header.Set("X-CSRF-Token", csrfToken+"tampered")
				return req
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:   "failure - POST invalid CSRF HMAC returns 403",
			method: http.MethodPost,
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/test", nil)
				req.AddCookie(&http.Cookie{Name: "access_token", Value: validToken})
				req.AddCookie(&http.Cookie{Name: "csrf_token", Value: "invalid-csrf-format"})
				req.Header.Set("X-CSRF-Token", "invalid-csrf-format")
				return req
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache, _, _, mw := freshAuthMiddleware(t)

			mockCache.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

			app := fiber.New()
			app.Add(tt.method, "/test", mw.JWTAuthenticator, mw.CSRFDoubleSubmit(), func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			})

			req := tt.setupRequest()
			resp, err := app.Test(req, 1000)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockCache.AssertExpectations(t)
		})
	}
}
