package shortener

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, handler *URLHandler) {
	short := app.Group("/shorten")

	short.Get("/", handler.GetAll)
	short.Post("/", handler.Shorten)
	short.Get("/:shortID", handler.GetByShortID)
	short.Delete("/:shortID", handler.Delete)
	short.Get("/user/:userID", handler.GetAllByUser)
}
