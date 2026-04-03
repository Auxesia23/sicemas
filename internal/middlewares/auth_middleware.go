package middlewares

import (
	"fmt"
	"log/slog"
	"net/netip"
	"situs-keagamaan/internal/cache"
	"situs-keagamaan/internal/dto"
	"situs-keagamaan/internal/geoip"
	"situs-keagamaan/internal/utils"
	"time"

	"github.com/casbin/casbin/v2"
	fcasbin "github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type AuthMiddleware interface {
	JWTAuthenticator(c *fiber.Ctx) error
	ZeroTrustValidator(c *fiber.Ctx) error
	CasbinAuthz() *fcasbin.Middleware
	GetContext(c *fiber.Ctx) error
	CSRFDoubleSubmit() fiber.Handler
}

type authMiddlewareImpl struct {
	enforcer *casbin.Enforcer
	locator  geoip.Locator
	cache    cache.Cache
	logger   *slog.Logger
}

func NewAuthMiddleware(
	enforcer *casbin.Enforcer,
	locator geoip.Locator,
	cache cache.Cache,
	logger *slog.Logger,
) AuthMiddleware {
	return &authMiddlewareImpl{
		enforcer: enforcer,
		locator:  locator,
		cache:    cache,
		logger:   logger,
	}
}

func (m *authMiddlewareImpl) JWTAuthenticator(c *fiber.Ctx) error {
	token := c.Cookies("access_token")
	if token == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var blocked bool
	m.cache.Get(c.Context(), fmt.Sprintf("blocked:%v", token), &blocked)
	if blocked {
		m.logger.Warn("blocked access token detected, request rejected")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claim, err := utils.ParseAccessToken(token)
	if err != nil {
		m.logger.Warn("invalid or expired access token", "error", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Locals("claim", claim)
	return c.Next()
}

func (m *authMiddlewareImpl) ZeroTrustValidator(c *fiber.Ctx) error {
	token := c.Cookies("access_token")
	jwtClaim := c.Locals("claim").(*dto.AccessToken)

	var ipStr string
	ipStr = c.Get("X-Forwarded-For")
	if ipStr == "" {
		ipStr = c.IP()
	}
	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
		m.logger.Error("failed to parse IP address", "error", err)
		return c.Status(500).SendString("Terjadi kesalahan!")
	}

	location, err := m.locator.Lookup(ip)
	if err != nil {
		m.logger.Error("failed to lookup GeoIP location", "error", err)
		return c.Status(500).SendString("Terjadi kesalahan!")
	}

	userAgent := string(c.Context().UserAgent())
	deviceId := c.Get("X-Device-Id")

	requestContect := &dto.SessionRequest{
		UserAgent:   userAgent,
		IPAddress:   ipStr,
		GeoLocation: location.City,
		DeviceID:    deviceId,
	}
	var curentSession dto.SessionValue
	sessionKey := fmt.Sprintf("rt:%v:%v", jwtClaim.Subject, jwtClaim.SID)
	err = m.cache.Get(c.Context(), sessionKey, &curentSession)
	if err != nil {
		if err == redis.Nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		m.logger.Error("failed to retrieve session from cache", "user_id", jwtClaim.Subject, "error", err)
		return c.Status(500).SendString("Terjadi kesalahan!")
	}

	trustScore := utils.CalculateTrustScore(requestContect, &curentSession)
	if trustScore <= 70 {
		m.logger.Warn("low trust score detected, session invalidated",
			"user_id", jwtClaim.Subject,
			"trust_score", trustScore,
			"ip_address", ipStr,
			"device_id", deviceId,
		)

		_ = m.cache.Delete(c.Context(), sessionKey)
		_ = m.cache.Set(c.Context(), fmt.Sprintf("blocked:%v", token), true, time.Minute*15)
		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			HTTPOnly: true,
		})

		c.Cookie(&fiber.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			HTTPOnly: true,
		})
		return c.Status(fiber.StatusUnauthorized).SendString("Terjadi perbuhan konteks dalam sesi!")
	}

	return c.Next()
}

func (m *authMiddlewareImpl) CasbinAuthz() *fcasbin.Middleware {
	authz := fcasbin.New(fcasbin.Config{
		Enforcer: m.enforcer,
		Lookup: func(c *fiber.Ctx) string {
			jwtClaim := c.Locals("claim").(*dto.AccessToken)
			return jwtClaim.Subject
		},
		Forbidden: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusForbidden).SendString("Anda tidak memiliki izin untuk melakukan aksi ini!")
		},
	})

	return authz
}

func (m *authMiddlewareImpl) GetContext(c *fiber.Ctx) error {
	var ipStr string
	ipStr = c.Get("X-True-Client-IP")
	if ipStr == "" {
		ipStr = c.Get("X-Forwarded-For")
	}
	if ipStr == "" {
		ipStr = c.IP()
	}
	if ipStr == "::1" || ipStr == "0:0:0:0:0:0:0:1" {
		ipStr = "127.0.0.1"
	}
	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
		m.logger.Error("failed to parse IP address in GetContext", "error", err)
		return c.Status(500).SendString("Terjadi kesalahan!")
	}

	deviceId := c.Get("X-Device-Id")
	locationn, _ := m.locator.Lookup(ip)
	userAgent := string(c.Context().UserAgent())

	loginContext := &dto.SessionRequest{
		UserAgent:   userAgent,
		IPAddress:   ipStr,
		GeoLocation: locationn.City,
		DeviceID:    deviceId,
	}

	c.Locals("context", loginContext)
	return c.Next()
}

func (m *authMiddlewareImpl) CSRFDoubleSubmit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		method := c.Method()
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			return c.Next()
		}

		claim := c.Locals("claim").(*dto.AccessToken)
		userId, err := uuid.Parse(claim.Subject)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Terjadi kesalahan!")
		}

		cookieToken := c.Cookies("csrf_token")
		headerToken := c.Get("X-CSRF-Token")

		if cookieToken == "" || headerToken == "" {
			m.logger.Warn("CSRF token mismatch detected", "user_id", userId, "reason", "missing_token")
			return c.Status(fiber.StatusForbidden).SendString("Akses ditolak. CSRF token tidak ditemukan")
		}

		if cookieToken != headerToken {
			m.logger.Warn("CSRF token mismatch detected", "user_id", userId, "reason", "token_mismatch")
			return c.Status(fiber.StatusForbidden).SendString("Akses ditolak. CSRF token tidak valid (Mismatch)")
		}

		isValid := utils.VerifyCSRFToken(cookieToken, userId)
		if !isValid {
			m.logger.Warn("CSRF token mismatch detected", "user_id", userId, "reason", "invalid_token")
			return c.Status(fiber.StatusForbidden).SendString("Akses ditolak. CSRF token tidak valid")
		}

		return c.Next()
	}
}
