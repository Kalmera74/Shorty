package shortener

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
// @Success 200 {array} ShortenResponse
// @Failure 404 {object} map[string]string "No URLs found"
// @Router /shorten [get]
func (h *URLHandler) GetAll(c *fiber.Ctx) error {
	shortModels, err := h.service.GetAllURLs()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	shortResponses := []ShortenResponse{}

	for _, shortModel := range shortModels {
		shortResponses = append(shortResponses, ShortenResponse{
			Id:          shortModel.ID,
			OriginalUrl: shortModel.OriginalUrl,
			ShortUrl:    shortModel.ShortUrl,
		})
	}

	return c.JSON(shortResponses)
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

	short, err := h.service.ShortenURL(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	responseObj := ShortenResponse{
		Id:          short.ID,
		OriginalUrl: short.OriginalUrl,
		ShortUrl:    short.ShortUrl,
	}

	return c.JSON(responseObj)
}

// GetById godoc
// @Summary Get a URL by its Id
// @Tags Shortener
// @Produce json
// @Param shortID path string true "Short URL ID"
// @Success 200 {object} ShortenResponse
// @Failure 404 {object} map[string]string "URL not found"
// @Router /shorten/{id} [get]
func (h *URLHandler) GetById(c *fiber.Ctx) error {
	paramId := c.Params("id")
	id, err := strconv.ParseUint(paramId, 10, 64)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	shortModel, err := h.service.GetById(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err})
	}
	shortenResponse := ShortenResponse{
		Id:          shortModel.ID,
		OriginalUrl: shortModel.OriginalUrl,
		ShortUrl:    shortModel.ShortUrl,
	}
	return c.JSON(shortenResponse)
}

// GetByShortUrl godoc
// @Summary Gets a shortened URL by its short URL
// @Tags Shortener
// @Param shortID path string true "Short URL ID"
// @Success 301 {string} string "Redirects to the original URL"
// @Failure 404 {object} map[string]string "URL not found"
// @Router /short/{url} [get]
func (h *URLHandler) GetByShortUrl(c *fiber.Ctx) error {
	shortUrl := c.Params("url")

	shortModel, err := h.service.GetByShortUrl(shortUrl)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "URL not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	shortenResponse := ShortenResponse{
		Id:          shortModel.ID,
		OriginalUrl: shortModel.OriginalUrl,
		ShortUrl:    shortModel.ShortUrl,
	}

	return c.JSON(shortenResponse)
}

// GetByOriginalUrl godoc
// @Summary Get a shortened URL by its original URL
// @Description Looks up a short URL by providing its original, long URL.
// @Tags Shortener
// @Accept json
// @Produce json
// @Param body body ShortenRequest true "Original URL to look up"
// @Success 200 {object} ShortenResponse
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "URL not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /shorten/search [post]
func (h *URLHandler) GetByOriginalUrl(c *fiber.Ctx) error {
	var req SearchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	shortModel, err := h.service.GetByLongUrl(req.OriginalUrl)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	responseObj := ShortenResponse{
		Id:          shortModel.ID,
		OriginalUrl: shortModel.OriginalUrl,
		ShortUrl:    shortModel.ShortUrl,
	}

	return c.JSON(responseObj)
}

// RedirectToOriginalUrl godoc
// @Summary Redirect to the original URL
// @Tags Shortener
// @Param shortID path string true "Short URL ID"
// @Success 301 {string} string "Redirects to the original URL"
// @Failure 404 {object} map[string]string "URL not found"
// @Router /{url} [get]
func (h *URLHandler) RedirectToOriginalUrl(c *fiber.Ctx) error {
	short := c.Params("url")

	shortModel, err := h.service.GetByShortUrl(short)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "URL not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Redirect(shortModel.OriginalUrl, fiber.StatusMovedPermanently)
}

// GetAllByUser godoc
// @Summary Get all URLs for a user
// @Tags Shortener
// @Produce json
// @Param userID path int true "User ID"
// @Success 200 {array} ShortenResponse
// @Failure 404 {object} map[string]string "No URLs found"
// @Router /shorten/user/{id} [get]
func (h *URLHandler) GetAllByUser(c *fiber.Ctx) error {
	userIDParam := c.Params("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
	}

	shortModels, err := h.service.GetAllByUser(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	shortResponses := []ShortenResponse{}

	for _, shortModel := range shortModels {
		shortResponses = append(shortResponses, ShortenResponse{
			Id:          shortModel.ID,
			OriginalUrl: shortModel.OriginalUrl,
			ShortUrl:    shortModel.ShortUrl,
		})
	}

	return c.JSON(shortResponses)
}

// Delete godoc
// @Summary Delete a URL by its short ID
// @Tags Shortener
// @Param shortID path string true "Short URL ID"
// @Success 204
// @Failure 404 {object} map[string]string "URL not found"
// @Router /shorten/{id} [delete]
func (h *URLHandler) Delete(c *fiber.Ctx) error {
	shortIDParam := c.Params("id")

	id, err := strconv.ParseUint(shortIDParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}

	err = h.service.DeleteURL(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
