package main

import (
	"log"
	"situs-keagamaan/internal/database"
	"situs-keagamaan/internal/handlers"
	"situs-keagamaan/internal/repositories"
	"situs-keagamaan/internal/services"
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
	rdb, err := database.InitRedis()
	if err != nil {
		log.Fatal(err.Error())
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	cache := repositories.NewCache(rdb)
	userRepo := repositories.NewUserRepo(db)

	userService := services.NewUserService(userRepo, cache)

	userHandler := handlers.NewUserHandler(userService, validate)

	cfg := config{
		Addr:         ":8080",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  20 * time.Second,
	}

	server := newServer(cfg, userHandler)

	server.run()
}
