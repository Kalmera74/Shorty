package shortener

import (
	"errors"
	"strconv"

	"github.com/Kalmera74/Shorty/internal/types"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type URLHandler struct {
	service IShortService
}

func NewShortHandler(service IShortService) *URLHandler {
	return &URLHandler{service: service}
}

// GetAll godoc
// @Summary Get all Shorts in the system
// @Tags Shortener
// @Produce json
// @Success 200 {array} ShortenResponse
// @Failure 404 {object} map[string]string "No Shorts found"
// @Router /shorten [get]
func (h *URLHandler) GetAll(c *fiber.Ctx) error {
	shortModels, err := h.service.GetAllURLs()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	shortResponses := []ShortenResponse{}

	for _, shortModel := range shortModels {
		shortResponses = append(shortResponses, ShortenResponse{
			Id:          uint(shortModel.ID),
			OriginalUrl: shortModel.OriginalUrl,
			ShortUrl:    shortModel.ShortUrl,
		})
	}

	return c.JSON(shortResponses)
}

// Shorten godoc
// @Summary Shorten a URL
// @Description Creates a new Short for the given long URL.
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
		Id:          uint(short.ID),
		OriginalUrl: short.OriginalUrl,
		ShortUrl:    short.ShortUrl,
	}

	return c.JSON(responseObj)
}

// GetById godoc
// @Summary Get a Short by its numeric ID
// @Tags Shortener
// @Produce json
// @Param id path string true "Short numeric ID"
// @Success 200 {object} ShortenResponse
// @Failure 404 {object} map[string]string "Short not found"
// @Router /shorten/{id} [get]
func (h *URLHandler) GetById(c *fiber.Ctx) error {
	paramId := c.Params("id")
	id, err := strconv.ParseUint(paramId, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	shortModel, err := h.service.GetById(types.ShortId(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	shortenResponse := ShortenResponse{
		Id:          uint(shortModel.ID),
		OriginalUrl: shortModel.OriginalUrl,
		ShortUrl:    shortModel.ShortUrl,
	}
	return c.JSON(shortenResponse)
}

// GetByShortUrl godoc
// @Summary Gets a Short by its short URL
// @Tags Shortener
// @Produce json
// @Param url path string true "Short URL string"
// @Success 200 {object} ShortenResponse
// @Failure 404 {object} map[string]string "Short not found"
// @Router /shorten/short/{url} [get]
func (h *URLHandler) GetByShortUrl(c *fiber.Ctx) error {
	shortUrl := c.Params("url")

	shortModel, err := h.service.GetByShortUrl(shortUrl)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Short not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	shortenResponse := ShortenResponse{
		Id:          uint(shortModel.ID),
		OriginalUrl: shortModel.OriginalUrl,
		ShortUrl:    shortModel.ShortUrl,
	}

	return c.JSON(shortenResponse)
}

// SearchByOriginalUrl godoc
// @Summary Search for a Short by its original URL
// @Description Looks up a Short by providing its original, long URL.
// @Tags Shortener
// @Accept json
// @Produce json
// @Param body body SearchRequest true "Original URL to look up"
// @Success 200 {object} ShortenResponse
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Short not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /shorten/search [post]
func (h *URLHandler) Search(c *fiber.Ctx) error {
	var req SearchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	shortModel, err := h.service.Search(req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Short not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	responseObj := ShortenResponse{
		Id:          uint(shortModel.ID),
		OriginalUrl: shortModel.OriginalUrl,
		ShortUrl:    shortModel.ShortUrl,
	}

	return c.JSON(responseObj)
}

// RedirectToOriginalUrl godoc
// @Summary Redirect to the original URL
// @Tags Shortener
// @Param url path string true "Short string ID"
// @Success 301 {string} string "Redirects to the original URL"
// @Failure 404 {object} map[string]string "Short not found"
// @Router /{url} [get]
func (h *URLHandler) RedirectToOriginalUrl(c *fiber.Ctx) error {
	short := c.Params("url")

	shortModel, err := h.service.GetByShortUrl(short)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Short not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Redirect(shortModel.OriginalUrl, fiber.StatusMovedPermanently)
}

// GetAllByUser godoc
// @Summary Get all Shorts for a user
// @Tags Shortener
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {array} ShortenResponse
// @Failure 404 {object} map[string]string "No Shorts found"
// @Router /shorten/user/{id} [get]
func (h *URLHandler) GetAllByUser(c *fiber.Ctx) error {
	userIDParam := c.Params("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
	}

	shortModels, err := h.service.GetAllByUser(types.UserId(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	shortResponses := []ShortenResponse{}

	for _, shortModel := range shortModels {
		shortResponses = append(shortResponses, ShortenResponse{
			Id:          uint(shortModel.ID),
			OriginalUrl: shortModel.OriginalUrl,
			ShortUrl:    shortModel.ShortUrl,
		})
	}

	return c.JSON(shortResponses)
}

// Delete godoc
// @Summary Delete a Short by its numeric ID
// @Tags Shortener
// @Param id path string true "Short numeric ID"
// @Success 204
// @Failure 404 {object} map[string]string "Short not found"
// @Router /shorten/{id} [delete]
func (h *URLHandler) Delete(c *fiber.Ctx) error {
	shortIDParam := c.Params("id")

	id, err := strconv.ParseUint(shortIDParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = h.service.DeleteURL(types.ShortId(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
