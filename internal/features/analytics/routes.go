package analytics

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, handler *analyticsHandler) {
	api := app.Group("/api/v1")
	analytics := api.Group("/analytics")

	analytics.Get("/", handler.GetAll)
	analytics.Get("/:shortUrl", handler.GetAllByShortUrl)
	analytics.Post("/", handler.Create)
}
