package shortener

import (
	"github.com/Kalmera74/Shorty/pkg/auth"
	"github.com/gofiber/fiber/v2"

	jwtware "github.com/gofiber/contrib/jwt"
)

func RegisterRoutes(app *fiber.App, handler *URLHandler) {

	api := app.Group("/api/v1")
	shorts := api.Group("/shorts")

	shorts.Post("/", jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(auth.JwtSecretKey)},
	}), handler.Shorten)

	shorts.Get("/", handler.GetAll)
	shorts.Post("/search", handler.Search)
	shorts.Get("/user/:id+/shorts", handler.GetAllByUser)
	shorts.Get("/:id+", handler.GetById)
	shorts.Delete("/:id+", handler.Delete)

	app.Get("/:url{[a-zA-Z0-9]+}", handler.RedirectToOriginalUrl)
}
