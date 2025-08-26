package shortener

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type URLHandler struct {
	service *URLService
}

func NewURLHandler(service *URLService) *URLHandler {
	return &URLHandler{service: service}
}

// GetAll godoc
// @Summary Get all URLs in the system
// @Tags Shortener
// @Produce json
// @Success 200 {array} ShortenModel
// @Failure 404 {object} map[string]string "No URLs found"
// @Router /shorten [get]
func (h *URLHandler) GetAll(c *fiber.Ctx) error {
	urls, err := h.service.GetAllURLs()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(urls)
}

// Shorten godoc
// @Summary Shorten a URL for a logged-in user
// @Description Shortens the given long URL. User must be logged in.
// @Tags Shortener
// @Accept json
// @Produce json
// @Param body body ShortenRequest true "Shorten URL request"
// @Success 200 {object} ShortenResponse
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /shorten [post]
func (h *URLHandler) Shorten(c *fiber.Ctx) error {
	var req ShortenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	shortID, err := h.service.ShortenURL(req.UserID, req.LongURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(ShortenResponse{ShortURL: shortID})
}

// GetByShortID godoc
// @Summary Get a URL by its short ID
// @Tags Shortener
// @Produce json
// @Param shortID path string true "Short URL ID"
// @Success 200 {object} ShortenModel
// @Failure 404 {object} map[string]string "URL not found"
// @Router /shorten/{shortID} [get]
func (h *URLHandler) GetByShortID(c *fiber.Ctx) error {
	shortID := c.Params("shortID")
	url, err := h.service.GetURLByShortID(shortID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(url)
}

// GetAllByUser godoc
// @Summary Get all URLs for a user
// @Tags Shortener
// @Produce json
// @Param userID path int true "User ID"
// @Success 200 {array} ShortenModel
// @Failure 404 {object} map[string]string "No URLs found"
// @Router /shorten/user/{userID} [get]
func (h *URLHandler) GetAllByUser(c *fiber.Ctx) error {
	userIDParam := c.Params("userID")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
	}

	urls, err := h.service.GetAllURLsByUser(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(urls)
}

// Delete godoc
// @Summary Delete a URL by its short ID
// @Tags Shortener
// @Param shortID path string true "Short URL ID"
// @Success 204
// @Failure 404 {object} map[string]string "URL not found"
// @Router /shorten/{shortID} [delete]
func (h *URLHandler) Delete(c *fiber.Ctx) error {
	shortID := c.Params("shortID")
	err := h.service.DeleteURL(shortID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
