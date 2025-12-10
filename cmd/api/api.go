package main

import (
	"log"
	"situs-keagamaan/internal/app/handlers"
	"situs-keagamaan/internal/middlewares"
	"time"

	"github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
)

type server struct {
	cfg         config
	middlewares *middlewares.Middlewares
	handlers    handlers.Handlers
}

func newServer(cfg config, middlewares *middlewares.Middlewares, handlers handlers.Handlers) *server {
	return &server{
		cfg,
		middlewares,
		handlers,
	}
}

type config struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func (s *server) run() {
	app := fiber.New(fiber.Config{
		ServerHeader: "KUA-Ciemas",
		ReadTimeout:  s.cfg.ReadTimeout,
		WriteTimeout: s.cfg.WriteTimeout,
		IdleTimeout:  s.cfg.IdleTimeout,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hallo, dari server KUA Ciemas")
	})

	auth := app.Group("/auth")
	{
		auth.Post("/login", s.handlers.Auth.Login)
		auth.Post("/verify-otp", s.middlewares.Auth.GetContext, s.handlers.Auth.VerifyOTP)
		auth.Post("/refresh", s.middlewares.Auth.GetContext, s.handlers.Auth.Refresh)
		auth.Post("/logout", s.handlers.Auth.Logout)
	}

	users := app.Group("/users")
	{
		users.Get("/me", func(c *fiber.Ctx) error {
			return c.SendStatus(200)
		})
		users.Use(s.middlewares.Auth.JWTAuthenticator)
		users.Use(s.middlewares.Auth.ZeroTrustValidator)

		users.Post("/",
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"user:create"}, casbin.WithValidationRule(casbin.MatchAllRule)),
			s.handlers.User.Register,
		)
	}

	log.Fatal(app.Listen(s.cfg.Addr))
}
