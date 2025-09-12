package analytics

import (
	"errors"

	"github.com/Kalmera74/Shorty/internal/types"
	"github.com/gofiber/fiber/v2"
)

type analyticsHandler struct {
	service IAnalyticsService
}

func NewAnalyticsHandler(service IAnalyticsService) *analyticsHandler {
	return &analyticsHandler{service: service}
}

// GetAllAnalytics godoc
// @Summary      Get paginated click analytics
// @Description  Returns paginated click analytics grouped by short URLs
// @Tags         analytics
// @Produce      json
// @Param        page     query int false "Page number" default(1)
// @Param        pageSize query int false "Items per page" default(10)
// @Success      200 {object} PaginatedAnalytics
// @Failure      404 {object} map[string]string "No clicks found"
// @Failure      500 {object} map[string]string "Failed to fetch analytics"
// @Router       /api/v1/analytics [get]
func (h *analyticsHandler) GetAllAnalytics(c *fiber.Ctx) error {
	// Parse query params
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	clickModels, total, err := h.service.GetAll(c.Context(), offset, pageSize)
	if err != nil {
		if errors.Is(err, ErrClicksNotFound) {
			return c.Status(fiber.StatusNotFound).
				JSON(fiber.Map{"error": "no click analytics found"})
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "failed to fetch analytics", "cause": err.Error()})
	}

	// Group clicks by short URL
	clickMap := make(map[string][]ClickModel)
	for _, model := range clickModels {
		clickMap[model.Short.ShortUrl] = append(clickMap[model.Short.ShortUrl], model)
	}

	// Build analytics response
	analyticsList := make([]Analysis, 0, len(clickMap))
	for shortUrl, list := range clickMap {
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

	// Paginated response
	response := PaginatedAnalytics{
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		Analytics: analyticsList,
	}

	return c.JSON(response)
}



// CreateClick godoc
// @Summary      CreateClick a new click record
// @Description  Records a click for a short URL
// @Tags         clicks
// @Accept       json
// @Produce      json
// @Param        click body ClickEvent true "Click information"
// @Success      201 {object} ClickModel
// @Failure      400 {object} map[string]string "Invalid request body"
// @Failure      500 {object} map[string]string "Failed to create click"
// @Router       /api/v1/analytics [post]
func (h *analyticsHandler) CreateClick(c *fiber.Ctx) error {
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

// GetAllAnalyticsByShortUrl godoc
// @Summary      Get paginated analytics for a specific short URL
// @Description  Returns click analytics for the given short URL
// @Tags         analytics
// @Produce      json
// @Param        shortUrl path string true "Short URL identifier"
// @Param        page query int false "Page number" default(1)
// @Param        pageSize query int false "Items per page" default(10)
// @Success      200 {object} PaginatedAnalysis
// @Failure      400 {object} map[string]string "Missing shortUrl parameter"
// @Failure      404 {object} map[string]string "No clicks found for this short URL"
// @Failure      500 {object} map[string]string "Failed to fetch clicks"
// @Router       /api/v1/analytics/{shortUrl} [get]
func (h *analyticsHandler) GetAllAnalyticsByShortUrl(c *fiber.Ctx) error {
	shortUrl := c.Params("shortUrl")
	if shortUrl == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "shortUrl parameter is required"})
	}

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	clickModels, total, err := h.service.GetAllByShortUrl(c.Context(), shortUrl, offset, pageSize)
	if err != nil {
		if errors.Is(err, ErrClickNotFound) {
			return c.Status(fiber.StatusNotFound).
				JSON(fiber.Map{"error": "no clicks found for this short URL"})
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "failed to fetch clicks", "cause": err.Error()})
	}

	if len(clickModels) == 0 {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": "no clicks found for this short URL"})
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

	response := PaginatedAnalysis{
		Total:   total,
		Page:    page,
		Limit:   pageSize,
		Results: analysis,
	}

	return c.JSON(response)
}


// GetAllClicks godoc
// @Summary      Get paginated click records
// @Description  Returns all individual click records (not grouped)
// @Tags         clicks
// @Produce      json
// @Param        page     query int false "Page number" default(1)
// @Param        pageSize query int false "Items per page" default(10)
// @Success      200 {object} PaginatedClicks
// @Failure      404 {object} map[string]string "No clicks found"
// @Failure      500 {object} map[string]string "Failed to fetch clicks"
// @Router       /api/v1/clicks [get]
func (h *analyticsHandler) GetAllClicks(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	clicks, total, err := h.service.GetAllClicks(c.Context(), offset, pageSize)
	if err != nil {
		if errors.Is(err, ErrClicksNotFound) {
			return c.Status(fiber.StatusNotFound).
				JSON(fiber.Map{"error": "no clicks found"})
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "failed to fetch clicks", "cause": err.Error()})
	}

	clickEvents := make([]ClickEvent, 0, len(clicks))
	for _, click := range clicks {
		clickEvents = append(clickEvents, ClickEvent{
			ShortID:   click.ShortID,
			Ip:        click.IpAddress,
			UserAgent: click.UserAgent,
			TimeStamp: click.CreatedAt,
		})
	}

	response := PaginatedClicks{
		Total:  total,
		Page:   page,
		Limit:  pageSize,
		Clicks: clickEvents,
	}

	return c.JSON(response)
}


// GetClickById godoc
// @Summary      Get a click record by ID
// @Description  Returns a single click record based on its ID
// @Tags         clicks
// @Produce      json
// @Param        id path int true "Click ID"
// @Success      200 {object} ClickEvent
// @Failure      400 {object} map[string]string "Invalid ID parameter"
// @Failure      404 {object} map[string]string "Click not found"
// @Failure      500 {object} map[string]string "Failed to fetch click"
// @Router       /api/v1/clicks/{id} [get]
func (h *analyticsHandler) GetClickById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid ID parameter"})
	}

	click, err := h.service.GetByID(c.Context(), types.ClickId(id))
	if err != nil {
		if errors.Is(err, ErrClickNotFound) {
			return c.Status(fiber.StatusNotFound).
				JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": "failed to fetch click", "cause": err.Error()})
	}

	clickEvent := ClickEvent{
		ShortID:   click.ShortID,
		Ip:        click.IpAddress,
		UserAgent: click.UserAgent,
		TimeStamp: click.CreatedAt,
	}

	return c.JSON(clickEvent)
}
