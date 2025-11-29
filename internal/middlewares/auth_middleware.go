package middlewares

import (
	"fmt"
	"net/netip"
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
}

type authMiddlewareImpl struct {
	enforcer *casbin.Enforcer
	locator  geoip.Locator
}

func NewAuthMiddleware(enforcer *casbin.Enforcer, locator geoip.Locator) AuthMiddleware {
	return &authMiddlewareImpl{
		enforcer,
		locator,
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
	claim, err := utils.ParseAccessToken(token[1])
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	fmt.Printf("User ID : %v", claim.Subject)

	c.Locals("user", "fajar")
	return c.Next()
}

func (m *authMiddlewareImpl) ZeroTrustValidator(c *fiber.Ctx) error {
	ipStr := c.Get("X-Forwarded-For")
	ip, _ := netip.ParseAddr(ipStr)
	locationn, _ := m.locator.Lookup(ip)
	fmt.Printf("Country : %v\nCity : %v\n", locationn.Country, locationn.City)
	return c.Next()
}

func (m *authMiddlewareImpl) CasbinAuthz() *fcasbin.Middleware {
	authz := fcasbin.New(fcasbin.Config{
		Enforcer: m.enforcer,
		Lookup: func(c *fiber.Ctx) string {
			return c.Locals("user").(string)
		},
	})

	return authz
}
