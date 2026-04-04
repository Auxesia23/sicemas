package main

import (
	"context"
	"log"
	"log/slog"
	"sicemas/internal/app/handlers"
	"sicemas/internal/app/repositories"
	"sicemas/internal/app/services"
	"sicemas/internal/auth"
	"sicemas/internal/cache"
	"sicemas/internal/database"
	"sicemas/internal/geoip"
	"sicemas/internal/logger"
	"sicemas/internal/middlewares"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to load .env")
	}

	appLog, writer := logger.NewLogger("./logs/app.log")
	slog.SetDefault(appLog)

	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err.Error())
	}
	rdb, err := cache.InitRedis()
	if err != nil {
		log.Fatal(err.Error())
	}
	cld, err := database.NewCloudinary(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	enforcer, err := auth.NewEnforcer()
	if err != nil {
		log.Fatal(err.Error())
	}
	locator, err := geoip.NewLocator()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Seeder data
	seeder := database.NewSeeder(db, enforcer)
	// User seeder
	if err := seeder.UserSeeder(); err != nil {
		log.Fatal(err.Error())
	}
	// Policiy seeder
	seeder.PolicySeeder()
	// Role seeder
	if err := seeder.RoleSeeder(); err != nil {
		log.Fatal(err.Error())
	}

	// Initiate validator
	validate := validator.New(validator.WithRequiredStructEnabled())

	// Initiate cache implementatiom
	cache := cache.NewCache(rdb)

	// Initiate repository layer
	userRepo := repositories.NewUserRepo(db)
	roleRepo := repositories.NewRoleRepo(db)
	jenisSitusRepo := repositories.NewJenisSitusRepo(db)
	situsKeagamaanRepo := repositories.NewSitusKeagamaanRepo(db)
	fotoSitusRepo := repositories.NewFotoSitusRepo(db)
	activityRepo := repositories.NewActivityRepo(db)
	dashboardRepo := repositories.NewDashboardRepo(db)

	// Initiate service layer
	userService := services.NewUserService(userRepo, activityRepo, enforcer, cache, appLog)
	authService := services.NewAuthService(userRepo, cache, appLog)
	roleService := services.NewRoleService(roleRepo, activityRepo, appLog)
	policyService := services.NewPolicyService(enforcer, activityRepo, appLog)
	jenisSitusService := services.NewJenisSitusService(jenisSitusRepo, activityRepo, appLog)
	situsKeagamaanService := services.NewSitusKeagamaanService(situsKeagamaanRepo, fotoSitusRepo, activityRepo, cld, enforcer, appLog)
	dashboardService := services.NewDashboardService(activityRepo, dashboardRepo, appLog)

	// Initiate handler layer
	userHandler := handlers.NewUserHandler(userService, validate)
	authHandler := handlers.NewAuthHandler(authService, validate)
	roleHandler := handlers.NewRoleHandler(roleService, validate)
	policyHandler := handlers.NewPolicyHandler(policyService, validate)
	jenisSitusHandler := handlers.NewJenisSitusHandler(jenisSitusService, validate)
	situsKeagamaanHandler := handlers.NewSitusKeagamaanHandler(situsKeagamaanService, validate)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	// handler compositor
	handlers := handlers.NewHandlers(userHandler, authHandler, roleHandler, policyHandler, jenisSitusHandler, situsKeagamaanHandler, dashboardHandler)

	// Initiate middleware
	authMiddleware := middlewares.NewAuthMiddleware(enforcer, locator, cache, appLog)
	rateLimiter := middlewares.NewRateLimiter(appLog)
	// middlewares compositor
	middlewares := middlewares.NewMiddlewares(authMiddleware, &rateLimiter)

	cfg := config{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  20 * time.Second,
	}

	server := newServer(cfg, middlewares, *handlers, writer, appLog)

	server.run()
}
