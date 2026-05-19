package middlewares

import (
	"log/slog"
	"sicemas/internal/dto"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/time/rate"
)

type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter interface {
	LimiterByDevice(r rate.Limit, burst int) fiber.Handler
	LimiterByIP(r rate.Limit, burst int) fiber.Handler
	LimiterByUser(r rate.Limit, burst int) fiber.Handler
}

type rateLimiterImpl struct {
	logger *slog.Logger
}

type limiterState struct {
	clients map[string]*Client
	mu      sync.Mutex

	rate rate.Limit
	burt int
}

func NewRateLimiter(logger *slog.Logger) RateLimiter {
	return &rateLimiterImpl{logger: logger}
}

func newLimiterState(rate rate.Limit, burst int) *limiterState {
	state := &limiterState{
		clients: make(map[string]*Client),
		rate:    rate,
		burt:    burst,
	}
	go state.cleanup()
	return state
}

func (rl *limiterState) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	client, exists := rl.clients[key]
	if !exists {
		limiter := rate.NewLimiter(rl.rate, rl.burt)
		rl.clients[key] = &Client{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}

	client.lastSeen = time.Now()
	return client.limiter
}

func (rl *limiterState) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for key, client := range rl.clients {
			if time.Since(client.lastSeen) > time.Minute {
				delete(rl.clients, key)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *rateLimiterImpl) LimiterByDevice(r rate.Limit, burst int) fiber.Handler {
	state := newLimiterState(r, burst)
	return func(c *fiber.Ctx) error {
		deviceID := c.Get("X-Device-Id")
		if deviceID == "" {
			return fiber.ErrBadRequest
		}
		limiter := state.getLimiter(deviceID)
		if !limiter.Allow() {
			rl.logger.Warn("rate limit exceeded", "device_id", deviceID, "ip_address", c.IP())
			return c.Status(fiber.StatusTooManyRequests).SendString("Terlalu banyak permintaan. Silakan tunggu beberapa saat.")
		}
		return c.Next()
	}
}

func (rl *rateLimiterImpl) LimiterByIP(r rate.Limit, burst int) fiber.Handler {
	state := newLimiterState(r, burst)
	return func(c *fiber.Ctx) error {
		ip := c.IP()
		if ip == "" {
			return fiber.ErrBadRequest
		}
		limiter := state.getLimiter(ip)
		if !limiter.Allow() {
			rl.logger.Warn("rate limit exceeded", "ip_address", c.IP())
			return c.Status(fiber.StatusTooManyRequests).SendString("Terlalu banyak permintaan. Silakan tunggu beberapa saat.")
		}
		return c.Next()
	}
}

// Rate limiter by user, make sure to use this middleware after JWTAuth middleware
func (rl *rateLimiterImpl) LimiterByUser(r rate.Limit, burst int) fiber.Handler {
	state := newLimiterState(r, burst)
	return func(c *fiber.Ctx) error {
		claim := c.Locals("claim").(*dto.AccessToken)
		if claim == nil {
			return fiber.ErrUnauthorized
		}
		limiter := state.getLimiter(claim.Subject)
		if !limiter.Allow() {
			rl.logger.Warn("rate limit exceeded", "user_id", claim.Subject, "device_id", c.Get("X-Device-Id"))
			return c.Status(fiber.StatusTooManyRequests).SendString("Terlalu banyak permintaan. Silakan tunggu beberapa saat.")
		}
		return c.Next()
	}

}
