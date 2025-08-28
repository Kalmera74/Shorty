package di

import (
	"github.com/Kalmera74/Shorty/internal/shortener"
	"github.com/Kalmera74/Shorty/internal/shortener/stores"
	"github.com/Kalmera74/Shorty/pkgs/redis"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupShortener(app *fiber.App, dbConn *gorm.DB) {

	store := stores.NewURLPostgresStore(dbConn)
	service := shortener.NewURLService(store, redis.Client)
	handler := shortener.NewURLHandler(service)

	shortener.RegisterRoutes(app, handler)
}
