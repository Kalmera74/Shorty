package user

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, handler *UserHandler) {
	users := app.Group("/users")

	users.Get("/", handler.GetAllUsers)
	users.Get("/:id", handler.GetUser)
	users.Post("/",handler.CreateUser)
	users.Put("/:id", handler.UpdateUser)
	users.Delete("/:id", handler.DeleteUser)
}