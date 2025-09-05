package analytics

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type analyticsHandler struct {
	service IAnalyticsService
}

func NewAnalyticsHandler(service IAnalyticsService) *analyticsHandler {
	return &analyticsHandler{service: service}
}

// GetAll godoc
// @Summary      Get all click analytics
// @Description  Returns all click analytics grouped by short URLs
// @Tags         analytics
// @Produce      json
// @Success      200 {array} Analysis
// @Failure      404 {object} map[string]string "No clicks found"
// @Failure      500 {object} map[string]string "Failed to fetch analytics"
// @Router       /api/analytics [get]
func (h *analyticsHandler) GetAll(c *fiber.Ctx) error {
	clickModels, err := h.service.GetAll(c.Context())
	if err != nil {
		if errors.Is(err, ErrClicksNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch analytics", "cause": err.Error()})
	}

	if len(clickModels) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no click analytics found"})
	}

	clickMap := make(map[string][]ClickModel)
	for _, model := range clickModels {
		clickMap[model.Short.ShortUrl] = append(clickMap[model.Short.ShortUrl], model)
	}

	analyticsList := make([]Analysis, 0, len(clickMap))
	for shortUrl, list := range clickMap {
		if len(list) == 0 {
			continue
		}

		analysis := Analysis{
			ShortUrl:     shortUrl,
			OriginalUrl:  list[0].Short.OriginalUrl,
			UsageDetails: make([]Usage, 0, len(list)),
		}

		for _, item := range list {
			analysis.UsageDetails = append(analysis.UsageDetails, Usage{
				ClickTimes: item.CreatedAt,
				IpAddress:  item.IpAddress,
				UserAgents: item.UserAgent,
			})
		}

		analyticsList = append(analyticsList, analysis)
	}

	return c.JSON(analyticsList)
}

// Create godoc
// @Summary      Create a new click record
// @Description  Records a click for a short URL
// @Tags         analytics
// @Accept       json
// @Produce      json
// @Param        click body ClickEvent true "Click information"
// @Success      201 {object} ClickModel
// @Failure      400 {object} map[string]string "Invalid request body"
// @Failure      500 {object} map[string]string "Failed to create click"
// @Router       /api/analytics [post]
func (h *analyticsHandler) Create(c *fiber.Ctx) error {
	var record ClickEvent
	if err := c.BodyParser(&record); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid request body", "cause": err.Error()})
	}

	click := ClickModel{
		IpAddress: record.Ip,
		UserAgent: record.UserAgent,
		CreatedAt: record.TimeStamp,
		ShortID:   record.ShortID,
	}

	createdClick, err := h.service.Create(c.Context(), click)
	if err != nil {
		if errors.Is(err, ErrClickCreateFail) {
			return c.Status(fiber.StatusInternalServerError).
				JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "failed to create click", "cause": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(createdClick)
}

// GetAllByShortUrl godoc
// @Summary      Get analytics for a specific short URL
// @Description  Returns click analytics for the given short URL
// @Tags         analytics
// @Produce      json
// @Param        shortUrl path string true "Short URL identifier"
// @Success      200 {object} Analysis
// @Failure      400 {object} map[string]string "Missing shortUrl parameter"
// @Failure      404 {object} map[string]string "No clicks found for this short URL"
// @Failure      500 {object} map[string]string "Failed to fetch clicks"
// @Router       /api/analytics/{shortUrl} [get]
func (h *analyticsHandler) GetAllByShortUrl(c *fiber.Ctx) error {
	shortUrl := c.Params("shortUrl")
	if shortUrl == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "shortUrl parameter is required"})
	}

	clickModels, err := h.service.GetAllByShortUrl(c.Context(), shortUrl)
	if err != nil {
		if errors.Is(err, ErrClickNotFound) {
			return c.Status(fiber.StatusNotFound).
				JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "failed to fetch clicks", "cause": err.Error()})
	}

	if len(clickModels) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no clicks found for this short URL"})
	}

	analysis := Analysis{
		ShortUrl:     clickModels[0].Short.ShortUrl,
		OriginalUrl:  clickModels[0].Short.OriginalUrl,
		UsageDetails: make([]Usage, 0, len(clickModels)),
	}

	for _, item := range clickModels {
		analysis.UsageDetails = append(analysis.UsageDetails, Usage{
			ClickTimes: item.CreatedAt,
			IpAddress:  item.IpAddress,
			UserAgents: item.UserAgent,
		})
	}

	return c.JSON(analysis)
}
