package analytics

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, handler *analyticsHandler) {
	api := app.Group("/api/v1")

	analytics := api.Group("/analytics")
	analytics.Get("/", handler.GetAllAnalytics)
	analytics.Get("/:shortUrl", handler.GetAllAnalyticsByShortUrl)

	clicks := api.Group("/clicks")
	clicks.Post("/", handler.CreateClick)
	clicks.Get("/", handler.GetAllClicks)
	clicks.Get("/:id", handler.GetClickById)
}
