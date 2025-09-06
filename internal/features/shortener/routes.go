package shortener

import (
	"github.com/Kalmera74/Shorty/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, handler *ShortHandler) {
	api := app.Group("/api/v1")
	shorts := api.Group("/shorts", middleware.Authenticate())

	shorts.Post("/", handler.Shorten)

	shorts.Get("/", middleware.Authorize("admin"), handler.GetAll)
	shorts.Post("/search", middleware.Authorize("admin"), handler.Search)
	shorts.Get("/user/:id/shorts", middleware.Authorize("admin"), handler.GetAllByUser)
	shorts.Get("/:id", middleware.Authorize("admin"), handler.GetById)
	shorts.Delete("/:id", middleware.Authorize("admin"), handler.Delete)

	app.Get("/:url<[a-zA-Z0-9_-]+>", handler.RedirectToOriginalUrl)
}
