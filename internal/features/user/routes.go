package user

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, handler *UserHandler) {

	api := app.Group("/api/v1")
	users := api.Group("/users")

	users.Get("/", handler.GetAllUsers)
	users.Post("/", handler.CreateUser)
	users.Post("/login", handler.Login)
	users.Get("/:id", handler.GetUser)
	users.Put("/:id", handler.UpdateUser)
	users.Delete("/:id", handler.DeleteUser)
}
