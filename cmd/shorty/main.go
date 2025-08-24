package main

import (
	"log"
	"os"

	"github.com/Kalmera74/Shorty/db"
	"github.com/Kalmera74/Shorty/internal/user"
	"github.com/Kalmera74/Shorty/internal/user/stores"
	"github.com/gofiber/fiber/v2"
)

func main() {
	dbConn, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err := db.AutoMigrate(dbConn); err != nil {
		log.Fatalf("Failed to perform migrations: %v", err)
	}

	userStore := stores.NewPostgresUserStore(dbConn)
	userService := user.NewUserService(userStore)
	userHandler := user.NewUserHandler(userService)

	app := fiber.New()

	user.RegisterRoutes(app, userHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}
