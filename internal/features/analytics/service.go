package analytics

import (
	"context"
	"fmt"
)

type IAnalyticsService interface {
	GetAll(ctx context.Context) ([]ClickModel, error)
	GetAllByShortUrl(ctx context.Context, shortUrl string) ([]ClickModel, error)
	Create(ctx context.Context, click ClickModel) (ClickModel, error)
}

type analyticsService struct {
	Repository IClickRepository
}

func NewAnalyticService(p IClickRepository) IAnalyticsService {
	return &analyticsService{p}
}

func (s *analyticsService) Create(ctx context.Context, click ClickModel) (ClickModel, error) {
	createdClick, err := s.Repository.Create(ctx, click)
	if err != nil {
		return ClickModel{}, fmt.Errorf("%w: %v", ErrClickCreateFail, err)
	}
	return createdClick, nil
}

func (s *analyticsService) GetAll(ctx context.Context) ([]ClickModel, error) {
	clicks, err := s.Repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if len(clicks) == 0 {
		return nil, ErrClickNotFound
	}
	return clicks, nil
}

func (s *analyticsService) GetAllByShortUrl(ctx context.Context, shortUrl string) ([]ClickModel, error) {
	clicks, err := s.Repository.GetAllByShortUrl(ctx, shortUrl)
	if err != nil {
		return nil, err
	}
	if len(clicks) == 0 {
		return nil, fmt.Errorf("%w: no clicks found for short url %s", ErrClickNotFound, shortUrl)
	}
	return clicks, nil
}
