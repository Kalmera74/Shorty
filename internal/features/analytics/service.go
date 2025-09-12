package analytics

import (
	"context"
	"fmt"

	"github.com/Kalmera74/Shorty/internal/types"
)

type IAnalyticsService interface {
	GetAll(ctx context.Context, offset, limit int) ([]ClickModel, int, error)
	GetAllByShortUrl(ctx context.Context, shortUrl string, offset, limit int) ([]ClickModel, int, error)
	Create(ctx context.Context, click ClickModel) (ClickModel, error)
	GetAllClicks(ctx context.Context, offset, limit int) ([]ClickModel, int, error)
	GetByID(ctx context.Context, id types.ClickId) (ClickModel, error)
}

type analyticsService struct {
	Repository IAnalyticsRepository
}

func NewAnalyticService(p IAnalyticsRepository) IAnalyticsService {
	return &analyticsService{p}
}

func (s *analyticsService) Create(ctx context.Context, click ClickModel) (ClickModel, error) {
	createdClick, err := s.Repository.Create(ctx, click)
	if err != nil {
		return ClickModel{}, fmt.Errorf("%w: %v", ErrClickCreateFail, err)
	}
	return createdClick, nil
}

func (s *analyticsService) GetAll(ctx context.Context, offset, limit int) ([]ClickModel, int, error) {
	clicks, total, err := s.Repository.GetAll(ctx, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("could not retrieve clicks: %w", err)
	}
	if len(clicks) == 0 {
		return nil, 0, ErrClicksNotFound
	}
	return clicks, total, nil
}

func (s *analyticsService) GetAllByShortUrl(ctx context.Context, shortUrl string, offset, limit int) ([]ClickModel, int, error) {
	clicks, total, err := s.Repository.GetAllByShortUrl(ctx, shortUrl, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	if len(clicks) == 0 {
		return nil, 0, fmt.Errorf("%w: no clicks found for short url %s", ErrClickNotFound, shortUrl)
	}
	return clicks, total, nil
}

func (s *analyticsService) GetAllClicks(ctx context.Context, offset, limit int) ([]ClickModel, int, error) {
	clicks, total, err := s.Repository.GetAll(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	if len(clicks) == 0 {
		return nil, 0, ErrClicksNotFound
	}
	return clicks, total, nil
}

func (s *analyticsService) GetByID(ctx context.Context, id types.ClickId) (ClickModel, error) {
	click, err := s.Repository.GetByID(ctx, id)
	if err != nil {
		return ClickModel{}, err
	}
	if click.ID == 0 {
		return ClickModel{}, ErrClickNotFound
	}
	return click, nil
}
