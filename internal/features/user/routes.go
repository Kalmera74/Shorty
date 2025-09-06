package user

import (
	"github.com/Kalmera74/Shorty/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, handler *UserHandler) {

	api := app.Group("/api/v1")
	api.Post("/register", handler.CreateUser)
	api.Post("/login", handler.Login)

	users := api.Group("/users", middleware.Authenticate(), middleware.Authorize("admin"))
	users.Get("/", handler.GetAllUsers)
	users.Get("/:id", handler.GetUser)
	users.Put("/:id", handler.UpdateUser)
	users.Delete("/:id", handler.DeleteUser)
}
