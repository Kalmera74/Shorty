package shortener

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/Kalmera74/Shorty/internal/types"
	"github.com/Kalmera74/Shorty/pkg/messaging"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type ShortHandler struct {
	service   IShortService
	messaging messaging.IMessaging
}

func NewShortHandler(service IShortService, messaging messaging.IMessaging) *ShortHandler {
	return &ShortHandler{service: service, messaging: messaging}
}

// GetAll godoc
// @Summary Get all shortened URLs
// @Description Retrieve all shortened URLs
// @Tags shorts
// @Produce json
// @Success 200 {array} ShortenResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/shorts [get]
func (h *ShortHandler) GetAll(c *fiber.Ctx) error {
	shortModels, err := h.service.GetAllURLs(c.Context(), )
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
// @Summary Create a new shortened URL
// @Description Shorten a given URL
// @Tags shorts
// @Accept json
// @Produce json
// @Param request body ShortenRequest true "Shorten Request"
// @Success 200 {object} ShortenResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/shorts [post]
func (h *ShortHandler) Shorten(c *fiber.Ctx) error {
	var req ShortenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	short, err := h.service.ShortenURL(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(ShortenResponse{
		Id:          uint(short.ID),
		OriginalUrl: short.OriginalUrl,
		ShortUrl:    short.ShortUrl,
	})
}

// GetById godoc
// @Summary Get a shortened URL by ID
// @Tags shorts
// @Produce json
// @Param id path int true "Short ID"
// @Success 200 {object} ShortenResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/shorts/{id} [get]
func (h *ShortHandler) GetById(c *fiber.Ctx) error {
	paramId := c.Params("id")
	id, err := strconv.ParseUint(paramId, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	shortModel, err := h.service.GetById(c.Context(), types.ShortId(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(ShortenResponse{
		Id:          uint(shortModel.ID),
		OriginalUrl: shortModel.OriginalUrl,
		ShortUrl:    shortModel.ShortUrl,
	})
}

// GetByShortUrl godoc
// @Summary Get a shortened URL by short code
// @Tags shorts
// @Produce json
// @Param url path string true "Short URL"
// @Success 200 {object} ShortenResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/shorts/short/{url} [get]
func (h *ShortHandler) GetByShortUrl(c *fiber.Ctx) error {
	shortUrl := c.Params("url")
	shortModel, err := h.service.GetByShortUrl(c.Context(), shortUrl)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Short not found"})
	}
	return c.JSON(ShortenResponse{
		Id:          uint(shortModel.ID),
		OriginalUrl: shortModel.OriginalUrl,
		ShortUrl:    shortModel.ShortUrl,
	})
}

// Search godoc
// @Summary Search for a shortened URL
// @Description Search using filters (URL, etc.)
// @Tags shorts
// @Accept json
// @Produce json
// @Param request body SearchRequest true "Search criteria"
// @Success 200 {object} ShortenResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/shorts/search [post]
func (h *ShortHandler) Search(c *fiber.Ctx) error {
	var req SearchRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	shortModel, err := h.service.Search(c.Context(), req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Short not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(ShortenResponse{
		Id:          uint(shortModel.ID),
		OriginalUrl: shortModel.OriginalUrl,
		ShortUrl:    shortModel.ShortUrl,
	})
}

// RedirectToOriginalUrl godoc
// @Summary Redirect to the original URL
// @Tags shorts
// @Produce json
// @Param url path string true "Short URL"
// @Success 301 {string} string "Redirects to the original URL"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /{url} [get]
func (h *ShortHandler) RedirectToOriginalUrl(c *fiber.Ctx) error {

	short := c.Params("url")
	shortModel, err := h.service.GetByShortUrl(c.Context(), short)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Short not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	go func() {
		event := map[string]any{
			"short_id":   shortModel.ID,
			"ip":         c.IP(),
			"user_agent": c.Get("User-Agent"),
			"time_stamp": time.Now(),
		}
		payload, _ := json.Marshal(event)
		_ = h.messaging.Publish("clicks", payload)
	}()

	return c.Redirect(shortModel.OriginalUrl, fiber.StatusMovedPermanently)
}

// GetAllByUser godoc
// @Summary Get all shorts for a specific user
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} ShortenResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/users/{id}/shorts [get]
func (h *ShortHandler) GetAllByUser(c *fiber.Ctx) error {
	userIDParam := c.Params("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user ID"})
	}
	shortModels, err := h.service.GetAllByUser(c.Context(), types.UserId(userID))
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
// @Summary Delete a shortened URL
// @Tags shorts
// @Produce json
// @Param id path int true "Short ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/shorts/{id} [delete]
func (h *ShortHandler) Delete(c *fiber.Ctx) error {
	shortIDParam := c.Params("id")
	id, err := strconv.ParseUint(shortIDParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err = h.service.DeleteURL(c.Context(), types.ShortId(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
