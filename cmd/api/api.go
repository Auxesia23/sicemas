package main

import (
	"io"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"runtime/debug"
	"sicemas/internal/app/handlers"
	"sicemas/internal/middlewares"
	"sicemas/ui"
	"time"

	"github.com/gofiber/contrib/casbin"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"
	"golang.org/x/time/rate"
)

type server struct {
	cfg         config
	middlewares *middlewares.Middlewares
	handlers    handlers.Handlers
	logWriter   io.Writer
	logger      *slog.Logger
}

func newServer(
	cfg config,
	middlewares *middlewares.Middlewares,
	handlers handlers.Handlers,
	logWriter io.Writer,
	logger *slog.Logger,
) *server {
	return &server{
		cfg:         cfg,
		middlewares: middlewares,
		handlers:    handlers,
		logWriter:   logWriter,
		logger:      logger,
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
		// ProxyHeader:  "X-Forwarded-For",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			s.logError("API_ERROR", code, c, err.Error(), nil)
			return c.Status(code).SendString(err.Error())
		},
	})

	app.Use(fiberRecover.New(fiberRecover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e any) {
			s.logError("PANIC_RECOVERED", 500, c, e, debug.Stack())
		},
	}))

	app.Use(logger.New(logger.Config{
		Output:     s.logWriter,
		TimeFormat: "2006-01-02T15:04:05.000Z07:00",
		Format:     `{"time":"${time}", "level":"INFO", "msg":"access_log", "status":${status}, "latency":"${latency}", "method":"${method}", "path":"${path}", "ip":"${ip}"}` + "\n",
	}))

	api := app.Group("/api")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hallo, dari server KUA Ciemas")
	})

	auth := api.Group("/auth")
	{
		auth.Use(s.middlewares.Limiter.LimiterByDevice(rate.Limit(1), 3))
		auth.Use(s.middlewares.Limiter.LimiterByIP(rate.Limit(1), 3))
		auth.Post("/login", s.handlers.Auth.Login)
		auth.Post("/verify-otp", s.middlewares.Auth.GetContext, s.handlers.Auth.VerifyOTP)
		auth.Post("/refresh", s.middlewares.Auth.GetContext, s.handlers.Auth.Refresh)
		auth.Post("/logout", s.handlers.Auth.Logout)
		auth.Post("/resend-stepup", s.middlewares.Auth.JWTAuthenticator,
			s.handlers.Auth.ResendStepUpOTP)
		auth.Post("/verify-stepup",
			s.middlewares.Auth.GetContext, s.middlewares.Auth.JWTAuthenticator,
			s.handlers.Auth.VerifyStepUpOTP)
	}

	dashboard := api.Group("/dashboard")
	{
		dashboard.Use(s.middlewares.Auth.JWTAuthenticator)
		dashboard.Use(s.middlewares.Auth.ZeroTrustValidator)
		dashboard.Use(s.middlewares.Limiter.LimiterByDevice(rate.Limit(5), 20))
		dashboard.Use(s.middlewares.Limiter.LimiterByUser(rate.Limit(5), 20))
		dashboard.Get("/",
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"dashboard:read"}, casbin.WithValidationRule(casbin.AtLeastOneRule)),
			s.handlers.DashboardHandler.GetDashboardData)
	}

	users := api.Group("/users")
	{
		users.Use(s.middlewares.Auth.JWTAuthenticator)
		users.Use(s.middlewares.Auth.CSRFDoubleSubmit())
		users.Use(s.middlewares.Auth.ZeroTrustValidator)

		users.Use(s.middlewares.Limiter.LimiterByDevice(rate.Limit(5), 20))
		users.Use(s.middlewares.Limiter.LimiterByUser(rate.Limit(5), 20))

		users.Get("/me", s.handlers.User.Me)
		users.Get("/profile", s.handlers.User.Profile)
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

	roles := api.Group("/roles")
	{
		roles.Use(s.middlewares.Auth.JWTAuthenticator)
		roles.Use(s.middlewares.Auth.CSRFDoubleSubmit())
		roles.Use(s.middlewares.Auth.ZeroTrustValidator)

		roles.Use(s.middlewares.Limiter.LimiterByDevice(rate.Limit(5), 20))
		roles.Use(s.middlewares.Limiter.LimiterByUser(rate.Limit(5), 20))

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

	policies := api.Group("/policies")
	{
		policies.Use(s.middlewares.Auth.JWTAuthenticator)
		policies.Use(s.middlewares.Auth.CSRFDoubleSubmit())
		policies.Use(s.middlewares.Auth.ZeroTrustValidator)

		policies.Use(s.middlewares.Limiter.LimiterByDevice(rate.Limit(5), 20))
		policies.Use(s.middlewares.Limiter.LimiterByUser(rate.Limit(5), 20))

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
	}

	jenisSitus := api.Group("/jenis-situs")
	{
		jenisSitus.Use(s.middlewares.Auth.JWTAuthenticator)
		jenisSitus.Use(s.middlewares.Auth.CSRFDoubleSubmit())
		jenisSitus.Use(s.middlewares.Auth.ZeroTrustValidator)

		jenisSitus.Use(s.middlewares.Limiter.LimiterByDevice(rate.Limit(5), 20))
		jenisSitus.Use(s.middlewares.Limiter.LimiterByUser(rate.Limit(5), 20))

		jenisSitus.Get("/",
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"jenis-situs:read"}, casbin.WithValidationRule(casbin.AtLeastOneRule)),
			s.handlers.JenisSitus.GetAllJenisSitus,
		)
		jenisSitus.Post("/",
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"jenis-situs:create"}, casbin.WithValidationRule(casbin.MatchAllRule)),
			s.handlers.JenisSitus.CreateJenisSitus,
		)
		jenisSitus.Put("/:id",
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"jenis-situs:update"}, casbin.WithValidationRule(casbin.MatchAllRule)),
			s.handlers.JenisSitus.UpdateJenisSitus,
		)
		jenisSitus.Delete("/:id",
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"jenis-situs:delete"}, casbin.WithValidationRule(casbin.MatchAllRule)),
			s.handlers.JenisSitus.DeleteJenisSitus,
		)
	}

	situs := api.Group("/situs")
	{
		situs.Use(s.middlewares.Auth.JWTAuthenticator)
		situs.Use(s.middlewares.Auth.CSRFDoubleSubmit())
		situs.Use(s.middlewares.Auth.ZeroTrustValidator)

		exportLimiterDevice := s.middlewares.Limiter.LimiterByDevice(rate.Every(3*time.Second), 1)
		exportLimiterUser := s.middlewares.Limiter.LimiterByUser(rate.Every(3*time.Second), 1)

		situs.Get("/export",
			exportLimiterDevice,
			exportLimiterUser,
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"situs:export"}, casbin.WithValidationRule(casbin.MatchAllRule)),
			s.handlers.SitusKeagamaan.ExportSitusToExcel,
		)

		stdLimiterDevice := s.middlewares.Limiter.LimiterByDevice(rate.Limit(5), 20)
		stdLimiterUser := s.middlewares.Limiter.LimiterByUser(rate.Limit(5), 20)

		situs.Post("/",
			stdLimiterDevice,
			stdLimiterUser,
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"situs:create"}, casbin.WithValidationRule(casbin.MatchAllRule)),
			s.handlers.SitusKeagamaan.CreateSitus,
		)
		situs.Get("/",
			stdLimiterDevice,
			stdLimiterUser,
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"situs:read_all", "situs:read_own"}, casbin.WithValidationRule(casbin.AtLeastOneRule)),
			s.handlers.SitusKeagamaan.GetAllSitus,
		)
		situs.Get("/:id",
			stdLimiterDevice,
			stdLimiterUser,
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"situs:read_all", "situs:read_own"}, casbin.WithValidationRule(casbin.AtLeastOneRule)),
			s.handlers.SitusKeagamaan.GetDetailSitus,
		)
		situs.Put("/:id",
			stdLimiterDevice,
			stdLimiterUser,
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"situs:update"}, casbin.WithValidationRule(casbin.MatchAllRule)),
			s.handlers.SitusKeagamaan.UpdateSitus,
		)
		situs.Patch("/:id/verify",
			stdLimiterDevice,
			stdLimiterUser,
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"situs:verify"}, casbin.WithValidationRule(casbin.MatchAllRule)),
			s.handlers.SitusKeagamaan.VerifySitus,
		)
		situs.Delete("/:id",
			stdLimiterDevice,
			stdLimiterUser,
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"situs:delete"}, casbin.WithValidationRule(casbin.MatchAllRule)),
			s.handlers.SitusKeagamaan.DeleteSitus,
		)
		situs.Post("/:id/foto",
			stdLimiterDevice,
			stdLimiterUser,
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"situs:create", "situs:update"}, casbin.WithValidationRule(casbin.AtLeastOneRule)),
			s.handlers.SitusKeagamaan.UploadFotoSitus,
		)
		situs.Delete("/:id/foto",
			stdLimiterDevice,
			stdLimiterUser,
			s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"situs:delete", "situs:update"}, casbin.WithValidationRule(casbin.AtLeastOneRule)),
			s.handlers.SitusKeagamaan.DeleteFotoSitus,
		)
	}

	aktivitas := api.Group("/activities")
	{
		aktivitas.Use(s.middlewares.Auth.JWTAuthenticator)
		aktivitas.Use(s.middlewares.Auth.ZeroTrustValidator)

		aktivitas.Use(s.middlewares.Limiter.LimiterByDevice(rate.Limit(5), 20))
		aktivitas.Use(s.middlewares.Limiter.LimiterByUser(rate.Limit(5), 20))

		aktivitas.Get("/", s.middlewares.Auth.CasbinAuthz().RequiresPermissions([]string{"activity:read"}, casbin.WithValidationRule(casbin.MatchAllRule)),
			s.handlers.ActivityHandler.GetActivityByDate,
		)
	}

	api.Use(func(c *fiber.Ctx) error {
		return c.Status(404).SendString("API endpoint not found")
	})

	buildFS, err := fs.Sub(ui.FrontendFS, "build")
	if err != nil {
		log.Fatal("Gagal memuat embedded filesystem: ", err)
	}

	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(buildFS),
		Index:        "index.html",
		NotFoundFile: "index.html",
	}))

	log.Fatal(app.Listen(s.cfg.Addr))
}

func (s *server) logError(msg string, code int, c *fiber.Ctx, err any, stack []byte) {
	fields := []any{
		"status", code,
		"method", c.Method(),
		"path", c.Path(),
		"ip", c.IP(),
		"error", err,
	}

	if stack != nil {
		fields = append(fields, "stack", string(stack))
	}

	s.logger.Error(msg, fields...)
}
