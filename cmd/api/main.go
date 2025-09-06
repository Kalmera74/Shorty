package main

import (
	"os"
	"time"

	"github.com/Kalmera74/Shorty/internal/db"
	"github.com/Kalmera74/Shorty/internal/features/shortener"
	"github.com/Kalmera74/Shorty/internal/features/user"
	"github.com/Kalmera74/Shorty/pkg/auth"
	"github.com/Kalmera74/Shorty/pkg/caching"
	"github.com/Kalmera74/Shorty/pkg/messaging"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

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
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := godotenv.Load(".debug.env"); err != nil {
		log.Info().Msg("No .env file found, using system environment variables")
	}

	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the database")
	}

	if err := db.AutoMigrate(dbConn); err != nil {
		log.Fatal().Err(err).Msg("Failed to perform migrations")
	}

	app := fiber.New()

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &log.Logger,
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	auth.InitJwt()
	cacher := caching.NewCacher()

	mq, err := messaging.NewRabbitMQConnection()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to rabbitmq")
	}

	clickQue := os.Getenv("CLICK_QUEUE")
	if clickQue == "" {
		clickQue = "click_queue"
	}
	if err := mq.DeclareQueue(clickQue); err != nil {
		log.Fatal().Err(err).Msg("Failed to declare RabbitMQ queue")
	}
	defer mq.Close()

	userStore := user.NewUserRepository(dbConn)
	userService := user.NewUserService(userStore)
	userHandler := user.NewUserHandler(userService)
	user.RegisterRoutes(app, userHandler)

	shortStore := shortener.NewShortRepository(dbConn)
	shortService := shortener.NewShortService(shortStore, cacher)
	shortHandler := shortener.NewShortHandler(shortService, mq)
	shortener.RegisterRoutes(app, shortHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal().Err(app.Listen(":" + port)).Msg("Shorty app encounter a problem, quitting")
}
