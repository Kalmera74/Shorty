package main

import (
	"log"
	"os"

	"github.com/Kalmera74/Shorty/db"
	"github.com/Kalmera74/Shorty/di"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
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

	di.SetupUser(app, dbConn)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}
