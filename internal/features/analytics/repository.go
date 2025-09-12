package analytics

import (
	"context"
	"errors"
	"fmt"

	"github.com/Kalmera74/Shorty/internal/types"
	"gorm.io/gorm"
)

type IAnalyticsRepository interface {
	GetAll(ctx context.Context, offset, limit int) ([]ClickModel, int, error)
	GetAllByShortUrl(ctx context.Context, shortUrl string, offset, limit int) ([]ClickModel, int, error)
	GetByID(ctx context.Context, id types.ClickId) (ClickModel, error)
	Create(ctx context.Context, click ClickModel) (ClickModel, error)
}

type postgresClickRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) IAnalyticsRepository {
	return &postgresClickRepository{db}
}

func (p *postgresClickRepository) Create(ctx context.Context, click ClickModel) (ClickModel, error) {
	if err := p.db.WithContext(ctx).Create(&click).Error; err != nil {
		return ClickModel{}, fmt.Errorf("%w: %v", ErrClickCreateFail, err)
	}
	return click, nil
}

func (p *postgresClickRepository) GetAll(ctx context.Context, offset, limit int) ([]ClickModel, int, error) {
	var clicks []ClickModel
	var total int64

	if err := p.db.WithContext(ctx).Model(&ClickModel{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count clicks: %w", err)
	}

	result := p.db.
		WithContext(ctx).
		Preload("Short").
		Limit(limit).
		Offset(offset).
		Order("id DESC").
		Find(&clicks)

	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to fetch clicks: %w", result.Error)
	}
	if len(clicks) == 0 {
		return nil, 0, ErrClicksNotFound
	}

	return clicks, int(total), nil
}

func (p *postgresClickRepository) GetAllByShortUrl(ctx context.Context, shortUrl string, offset, limit int) ([]ClickModel, int, error) {
	var clicks []ClickModel
	var total int64

	// Count total for this shortUrl
	if err := p.db.WithContext(ctx).
		Model(&ClickModel{}).
		Joins("JOIN short_models ON short_models.id = click_models.short_id").
		Where("short_models.short_url = ?", shortUrl).
		Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count clicks for short url %s: %w", shortUrl, err)
	}

	// Fetch paginated clicks
	if err := p.db.WithContext(ctx).
		Joins("JOIN short_models ON short_models.id = click_models.short_id").
		Where("short_models.short_url = ?", shortUrl).
		Preload("Short").
		Order("click_models.id DESC").
		Offset(offset).
		Limit(limit).
		Find(&clicks).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch clicks for short url %s: %w", shortUrl, err)
	}

	if len(clicks) == 0 {
		return nil, 0, fmt.Errorf("no clicks found for short url %s: %w", shortUrl, ErrClicksNotFound)
	}

	return clicks, int(total), nil
}

func (p *postgresClickRepository) GetByID(ctx context.Context, id types.ClickId) (ClickModel, error) {
	var click ClickModel
	if err := p.db.
		WithContext(ctx).
		Preload("Short").
		First(&click, id).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ClickModel{}, fmt.Errorf("click with id %d not found: %w", id, ErrClickNotFound)
		}
		return ClickModel{}, fmt.Errorf("failed to fetch click with id %d: %w", id, err)
	}
	return click, nil
}
