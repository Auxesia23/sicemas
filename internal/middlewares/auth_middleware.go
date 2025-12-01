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
	fmt.Printf("JTI : %v\nUSER ID : %v\n", jwtClaim.ID, jwtClaim.Subject)

	var ipStr string
	ipStr = c.Get("X-Forwarded-For")
	if ipStr == "" {
		ipStr = c.IP()
	}
	ip, _ := netip.ParseAddr(ipStr)
	locationn, _ := m.locator.Lookup(ip)

	userAgent := string(c.Context().UserAgent())
	fmt.Printf("IP Addr : %v\nUser Agent : %v\nCountry : %v\nCity : %v\n", ipStr, userAgent, locationn.Country, locationn.City)
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
	ip, _ := netip.ParseAddr(ipStr)

	deviceId := c.Get("X-Device-Id")
	fmt.Println("DEVICE ID : ", deviceId)
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
