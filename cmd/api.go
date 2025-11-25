package main

import (
	"log"
	"situs-keagamaan/internal/handlers"
	"time"

	"github.com/gofiber/fiber/v2"
)

type server struct {
	cfg         config
	userHandler handlers.UserHandler
}

func newServer(cfg config, userHandler handlers.UserHandler) *server {
	return &server{
		cfg,
		userHandler,
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
		auth.Post("/register", s.userHandler.RegisterUser)
		auth.Post("/login", s.userHandler.LoginUser)
		auth.Post("/verify-otp", s.userHandler.VerifyOTP)
	}

	log.Fatal(app.Listen(s.cfg.Addr))
}
