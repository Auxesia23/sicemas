package main

import (
	"log"
	"situs-keagamaan/internal/app/handlers"
	"situs-keagamaan/internal/dto"
	"situs-keagamaan/internal/middlewares"
	"time"

	"github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Requested-With",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

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
		users.Use(s.middlewares.Auth.JWTAuthenticator)
		users.Use(s.middlewares.Auth.ZeroTrustValidator)

		users.Get("/me", s.middlewares.Auth.JWTAuthenticator, func(c *fiber.Ctx) error {
			jwtClaim := c.Locals("claim").(*dto.AccessToken)

			return c.Status(200).JSON(fiber.Map{
				"nama": jwtClaim.NamaLengkap,
				"id":   jwtClaim.Subject,
			})
		})
		users.Get("/",
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"user:read"}, casbin.WithValidationRule(casbin.MatchAllRule)),
			s.handlers.User.GetAllUser,
		)
		users.Post("/",
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"user:create"}, casbin.WithValidationRule(casbin.MatchAllRule)),
			s.handlers.User.Register,
		)
		users.Put("/:userId",
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"user:update"}, casbin.WithValidationRule(casbin.MatchAllRule)),
			s.handlers.User.UpdateUser,
		)
		users.Delete("/:userId",
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"user:update"}, casbin.WithValidationRule(casbin.MatchAllRule)),
			s.handlers.User.DeleteUser,
		)
	}

	roles := app.Group("/roles")
	{
		roles.Use(s.middlewares.Auth.JWTAuthenticator)
		roles.Use(s.middlewares.Auth.ZeroTrustValidator)

		roles.Get("/",
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"role:read"}, casbin.WithValidationRule(casbin.MatchAllRule)),
			s.handlers.Role.GetAllRole,
		)
		roles.Post("/",
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"role:create"}, casbin.WithValidationRule(casbin.MatchAllRule)),
			s.handlers.Role.CreateRole,
		)
		roles.Delete("/:roleId",
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"role:create"}, casbin.WithValidationRule(casbin.MatchAllRule)),
			s.handlers.Role.DeleteRole,
		)
	}

	policies := app.Group("/policies")
	{
		policies.Use(s.middlewares.Auth.JWTAuthenticator)
		policies.Use(s.middlewares.Auth.ZeroTrustValidator)

		policies.Get("/",
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"policy:read"}, casbin.WithValidationRule(casbin.MatchAllRule)),
			s.handlers.Policy.GetFilteredPolicy,
		)
		policies.Post("/",
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"policy:create"}, casbin.WithValidationRule(casbin.MatchAllRule)),
			s.handlers.Policy.AddPolicy,
		)
		policies.Delete("/",
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"policy:delete"}, casbin.WithValidationRule(casbin.MatchAllRule)),
			s.handlers.Policy.RemovePolicy,
		)

		groups := policies.Group("/groups")
		{
			groups.Get("/",
				s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"policy:read"}, casbin.WithValidationRule(casbin.MatchAllRule)),
				s.handlers.Policy.GetFilteredGroupPolicy,
			)
			groups.Post("/",
				s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"policy:create"}, casbin.WithValidationRule(casbin.MatchAllRule)),
				s.handlers.Policy.AddGroupPolicy,
			)
			groups.Delete("/",
				s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"policy:delete"}, casbin.WithValidationRule(casbin.MatchAllRule)),
				s.handlers.Policy.RemoveGroupPolicy,
			)
		}
	}

	log.Fatal(app.Listen(s.cfg.Addr))
}
