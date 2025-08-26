package di

import (
	"github.com/Kalmera74/Shorty/internal/user"
	"github.com/Kalmera74/Shorty/internal/user/stores"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupUser(app *fiber.App, dbConn *gorm.DB) {
	userStore := stores.NewUserPostgresStore(dbConn)
	userService := user.NewUserService(userStore)
	userHandler := user.NewUserHandler(userService)

	user.RegisterRoutes(app, userHandler)

}
