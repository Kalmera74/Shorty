package shortener

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, handler *URLHandler) {

	short := app.Group("/shorten")

	short.Get("/", handler.GetAll)
	short.Post("/", handler.Shorten)
	short.Post("/search", handler.Search)

	short.Get("/short/:url{[a-zA-Z0-9]+}", handler.GetByShortUrl)

	short.Get("/user/:id+", handler.GetAllByUser)
	short.Get("/:id+", handler.GetById)
	short.Delete("/:id+", handler.Delete)

	app.Get("/:url{[a-zA-Z0-9]+}", handler.RedirectToOriginalUrl)

}
