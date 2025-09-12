package analytics

import (
	"context"
	"errors"
	"fmt"

	"github.com/Kalmera74/Shorty/internal/types"
	"gorm.io/gorm"
)

type IAnalyticsRepository interface {
	GetAll(ctx context.Context) ([]ClickModel, error)
	GetAllByShortUrl(ctx context.Context, shortUrl string) ([]ClickModel, error)
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

func (p *postgresClickRepository) GetAll(ctx context.Context) ([]ClickModel, error) {
	var clicks []ClickModel
	if err := p.db.
		WithContext(ctx).
		Preload("Short").
		Find(&clicks).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch clicks: %w", err)
	}
	if len(clicks) == 0 {
		return nil, fmt.Errorf("no clicks found %w", ErrClicksNotFound)
	}
	return clicks, nil
}

func (p *postgresClickRepository) GetAllByShortUrl(ctx context.Context, shortUrl string) ([]ClickModel, error) {
	var clicks []ClickModel

	if err := p.db.
		WithContext(ctx).
		Joins("JOIN short_models ON short_models.id = click_models.short_id").
		Where("short_models.short_url = ?", shortUrl).
		Preload("Short").
		Find(&clicks).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no clicks found for short url %s: %w", shortUrl, ErrClicksNotFound)
		}
		return nil, fmt.Errorf("could not find clicks for short url %s: %w reason: %v",
			shortUrl, ErrClicksNotFound, err)
	}

	if len(clicks) == 0 {
		return nil, fmt.Errorf("no clicks found for short url %s: %w", shortUrl, ErrClicksNotFound)
	}

	return clicks, nil
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
