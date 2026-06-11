package middlewares

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"sicemas/internal/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func setupRateLimiterTest() RateLimiter {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewRateLimiter(logger)
}

// ==========================================
// LimiterByDevice Tests
// ==========================================

func TestLimiterByDevice(t *testing.T) {
	rl := setupRateLimiterTest()

	t.Run("success - within rate limit", func(t *testing.T) {
		app := fiber.New()
		app.Get("/test", rl.LimiterByDevice(rate.Limit(10), 5), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Device-Id", "device-test-001")

		resp, err := app.Test(req, 1000)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("failure - missing device ID returns 400", func(t *testing.T) {
		app := fiber.New()
		app.Get("/test", rl.LimiterByDevice(rate.Limit(10), 5), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		resp, err := app.Test(req, 1000)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("failure - empty device ID returns 400", func(t *testing.T) {
		app := fiber.New()
		app.Get("/test", rl.LimiterByDevice(rate.Limit(10), 5), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Device-Id", "")

		resp, err := app.Test(req, 1000)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("failure - rate exceeded returns 429", func(t *testing.T) {
		app := fiber.New()
		hdl := rl.LimiterByDevice(rate.Limit(1), 1)
		app.Get("/test", hdl, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Device-Id", "device-rate-test")

		resp1, err := app.Test(req, 1000)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp1.StatusCode)

		req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
		req2.Header.Set("X-Device-Id", "device-rate-test")

		resp2, err := app.Test(req2, 1000)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTooManyRequests, resp2.StatusCode)
	})

	t.Run("different devices have separate rate limits", func(t *testing.T) {
		app := fiber.New()
		hdl := rl.LimiterByDevice(rate.Limit(1), 1)
		app.Get("/test", hdl, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
		req1.Header.Set("X-Device-Id", "device-a")

		resp1, err := app.Test(req1, 1000)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp1.StatusCode)

		req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
		req2.Header.Set("X-Device-Id", "device-b")

		resp2, err := app.Test(req2, 1000)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp2.StatusCode)
	})
}

// ==========================================
// LimiterByIP Tests
// ==========================================

func TestLimiterByIP(t *testing.T) {
	rl := setupRateLimiterTest()

	t.Run("success - within rate limit", func(t *testing.T) {
		app := fiber.New()
		app.Get("/test", rl.LimiterByIP(rate.Limit(10), 5), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		resp, err := app.Test(req, 1000)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("failure - rate exceeded returns 429", func(t *testing.T) {
		app := fiber.New()
		hdl := rl.LimiterByIP(rate.Limit(1), 1)
		app.Get("/test", hdl, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req1 := httptest.NewRequest(http.MethodGet, "/test", nil)

		resp1, err := app.Test(req1, 1000)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp1.StatusCode)

		req2 := httptest.NewRequest(http.MethodGet, "/test", nil)

		resp2, err := app.Test(req2, 1000)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTooManyRequests, resp2.StatusCode)
	})

	t.Run("same IP respects rate limit", func(t *testing.T) {
		app := fiber.New()
		hdl := rl.LimiterByIP(rate.Limit(1), 1)
		app.Get("/test", hdl, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req1 := httptest.NewRequest(http.MethodGet, "/test", nil)
		resp1, err := app.Test(req1, 1000)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp1.StatusCode)

		req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
		resp2, err := app.Test(req2, 1000)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTooManyRequests, resp2.StatusCode)
	})
}

// ==========================================
// LimiterByUser Tests
// ==========================================

func TestLimiterByUser(t *testing.T) {
	rl := setupRateLimiterTest()

	userID := uuid.New()

	t.Run("success - within rate limit", func(t *testing.T) {
		app := fiber.New()
		app.Get("/test", withClaim(userID), rl.LimiterByUser(rate.Limit(10), 5), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		resp, err := app.Test(req, 1000)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("failure - no claim returns 401", func(t *testing.T) {
		app := fiber.New()
		app.Get("/test", rl.LimiterByUser(rate.Limit(10), 5), func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		resp, err := app.Test(req, 1000)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("failure - rate exceeded returns 429", func(t *testing.T) {
		app := fiber.New()
		hdl := rl.LimiterByUser(rate.Limit(1), 1)
		app.Get("/test", withClaim(userID), hdl, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req1 := httptest.NewRequest(http.MethodGet, "/test", nil)

		resp1, err := app.Test(req1, 1000)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp1.StatusCode)

		req2 := httptest.NewRequest(http.MethodGet, "/test", nil)

		resp2, err := app.Test(req2, 1000)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusTooManyRequests, resp2.StatusCode)
	})

	t.Run("different users have separate rate limits", func(t *testing.T) {
		app := fiber.New()
		hdl := rl.LimiterByUser(rate.Limit(1), 1)
		app.Get("/test", withClaim(uuid.New()), hdl, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})
		app.Get("/test2", withClaim(uuid.New()), hdl, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req1 := httptest.NewRequest(http.MethodGet, "/test", nil)

		resp1, err := app.Test(req1, 1000)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp1.StatusCode)

		req2 := httptest.NewRequest(http.MethodGet, "/test2", nil)

		resp2, err := app.Test(req2, 1000)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp2.StatusCode)
	})
}

// Note: LimiterByUser should be used after JWTAuthenticator middleware.
// For testing, we wrap the claim-setting in a test helper middleware.
func withClaim(uid uuid.UUID) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("claim", &dto.AccessToken{
			NamaLengkap: "Test User",
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: uid.String(),
			},
		})
		return c.Next()
	}
}
