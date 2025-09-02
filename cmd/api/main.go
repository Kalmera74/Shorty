package main

import (
	"log"
	"os"

	"github.com/Kalmera74/Shorty/internal/db"
	"github.com/Kalmera74/Shorty/internal/features/shortener"
	"github.com/Kalmera74/Shorty/internal/features/user"
	"github.com/Kalmera74/Shorty/pkg/auth"
	"github.com/Kalmera74/Shorty/pkg/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	_ "github.com/Kalmera74/Shorty/docs"
	"github.com/gofiber/swagger"
)

// @title Shorty API
// @version 1.0
// @description REST API for Shorty URL shortener
// @license.name MIT
// @host localhost:8080
// @BasePath /
func main() {

	if err := godotenv.Load(".debug.env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err := db.AutoMigrate(dbConn); err != nil {
		log.Fatalf("Failed to perform migrations: %v", err)
	}

	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)


	auth.InitJwt()
	redis.InitRedisClient()

	userStore := user.NewUserRepository(dbConn)
	userService := user.NewUserService(userStore)
	userHandler := user.NewUserHandler(userService)
	user.RegisterRoutes(app, userHandler)

	shortStore := shortener.NewShortRepository(dbConn)
	shortService := shortener.NewShortService(shortStore, redis.NewCacher(redis.Client))
	shortHandler := shortener.NewShortHandler(shortService)
	shortener.RegisterRoutes(app, shortHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}
