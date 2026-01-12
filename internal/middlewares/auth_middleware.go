package middlewares

import (
	"fmt"
	"net/netip"
	"situs-keagamaan/internal/cache"
	"situs-keagamaan/internal/dto"
	"situs-keagamaan/internal/geoip"
	"situs-keagamaan/internal/utils"
	"strings"

	"github.com/casbin/casbin/v2"
	fcasbin "github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type AuthMiddleware interface {
	JWTAuthenticator(c *fiber.Ctx) error
	ZeroTrustValidator(c *fiber.Ctx) error
	CasbinAuthz() *fcasbin.Middleware
	GetContext(c *fiber.Ctx) error
}

type authMiddlewareImpl struct {
	enforcer *casbin.Enforcer
	locator  geoip.Locator
	cache    cache.Cache
}

func NewAuthMiddleware(enforcer *casbin.Enforcer, locator geoip.Locator, cache cache.Cache) AuthMiddleware {
	return &authMiddlewareImpl{
		enforcer,
		locator,
		cache,
	}
}

func (m *authMiddlewareImpl) JWTAuthenticator(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	token := strings.Split(authHeader, " ")
	if token[0] != "Bearer" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	var blocked bool
	m.cache.Get(c.Context(), fmt.Sprintf("blocked:%v", token[1]), &blocked)
	if blocked {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	claim, err := utils.ParseAccessToken(token[1])
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Locals("claim", claim)
	return c.Next()
}

func (m *authMiddlewareImpl) ZeroTrustValidator(c *fiber.Ctx) error {
	jwtClaim := c.Locals("claim").(*dto.AccessToken)

	var ipStr string
	ipStr = c.Get("X-Forwarded-For")
	if ipStr == "" {
		ipStr = c.IP()
	}
	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
		return c.Status(500).SendString("Terjadi kesalahan!")
	}

	location, err := m.locator.Lookup(ip)
	if err != nil {
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
	err = m.cache.Get(c.Context(), fmt.Sprintf("rt:%v:%v", jwtClaim.Subject, jwtClaim.SID), &curentSession)
	if err != nil {
		if err == redis.Nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Status(500).SendString("Terjadi kesalahan!")
	}

	trustScore := utils.CalculateTrustScore(requestContect, &curentSession)
	if trustScore <= 70 {
		return c.Status(fiber.StatusForbidden).SendString("Terjadi perbuhan konteks dalam sesi!")
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
	})

	return authz
}

func (m *authMiddlewareImpl) GetContext(c *fiber.Ctx) error {
	var ipStr string
	ipStr = c.Get("X-Forwarded-For")
	if ipStr == "" {
		ipStr = c.IP()
	}
	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
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
	fmt.Printf("IP Addr : %v\nLocation : %v\nUser Agent : %v\nDevice ID : %v\n", ipStr, locationn.City, userAgent, deviceId)
	return c.Next()
}
