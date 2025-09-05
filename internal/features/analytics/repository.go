package analytics

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IClickRepository interface {
	GetAll(ctx context.Context) ([]ClickModel, error)
	GetAllByShortUrl(ctx context.Context, shortUrl string) ([]ClickModel, error)
	Create(ctx context.Context, click ClickModel) (ClickModel, error)
}

type postgresClickRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) IClickRepository {
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
