package shortener

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, handler *URLHandler) {

	short := app.Group("/shorten")

	short.Get("/", handler.GetAll)
	short.Post("/", handler.Shorten)
	short.Post("/search", handler.GetByOriginalUrl)
	short.Get("/short/:url", handler.GetByShortUrl)
	short.Get("/user/:id", handler.GetAllByUser)
	short.Get("/:id", handler.GetById)
	short.Delete("/:id", handler.Delete)

	app.Get("/:url", handler.RedirectToOriginalUrl)

}
