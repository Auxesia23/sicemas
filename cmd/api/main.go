package main

import (
	"log"
	"situs-keagamaan/internal/app/handlers"
	"situs-keagamaan/internal/app/repositories"
	"situs-keagamaan/internal/app/services"
	"situs-keagamaan/internal/auth"
	"situs-keagamaan/internal/cache"
	"situs-keagamaan/internal/database"
	"situs-keagamaan/internal/geoip"
	"situs-keagamaan/internal/middlewares"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to load .env")
	}

	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err.Error())
	}
	rdb, err := cache.InitRedis()
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

	// Initiate service layer
	userService := services.NewUserService(userRepo, enforcer, cache)
	authService := services.NewAuthService(userRepo, cache)
	roleService := services.NewRoleService(roleRepo)
	policyService := services.NewPolicyService(enforcer)
	jenisSitusService := services.NewJenisSitusService(jenisSitusRepo)
	situsKeagamaanService := services.NewSitusKeagamaanService(situsKeagamaanRepo)

	// Initiate handler layer
	userHandler := handlers.NewUserHandler(userService, validate)
	authHandler := handlers.NewAuthHandler(authService, validate)
	roleHandler := handlers.NewRoleHandler(roleService, validate)
	policyHandler := handlers.NewPolicyHandler(policyService, validate)
	jenisSitusHandler := handlers.NewJenisSitusHandler(jenisSitusService, validate)
	situsKeagamaanHandler := handlers.NewSitusKeagamaanHandler(situsKeagamaanService, validate)
	// handler compositor
	handlers := handlers.NewHandlers(userHandler, authHandler, roleHandler, policyHandler, jenisSitusHandler, situsKeagamaanHandler)

	// Initiate middleware
	authMiddleware := middlewares.NewAuthMiddleware(enforcer, locator, cache)
	// middlewares compositor
	middlewares := middlewares.NewMiddlewares(authMiddleware)

	cfg := config{
		Addr:         ":8080",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  20 * time.Second,
	}

	server := newServer(cfg, middlewares, *handlers)

	server.run()
}
