package analytics

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, handler *analyticsHandler) {
	api := app.Group("/api/v1")
	analytics := api.Group("/analytics")
	clicks := app.Group("/clicks")

	clicks.Post("/clicks", handler.CreateClick)
	clicks.Get("/", handler.GetAllClicks)
	clicks.Get("/:id", handler.GetClickById)

	analytics.Get("/", handler.GetAllAnalytics)
	analytics.Get("/:shortUrl", handler.GetAllAnalyticsByShortUrl)
}
