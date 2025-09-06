package analytics

import (
	"context"
	"fmt"
)

type IAnalyticsService interface {
	GetAll(ctx context.Context) ([]ClickModel, error)
	GetAllByShortUrl(ctx context.Context, shortUrl string) ([]ClickModel, error)
	Create(ctx context.Context, click ClickModel) (ClickModel, error)
	GetAllClicks(ctx context.Context) ([]ClickModel, error)
	GetByID(ctx context.Context, id uint) (ClickModel, error)
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

func (s *analyticsService) GetAllClicks(ctx context.Context) ([]ClickModel, error) {
	clicks, err := s.Repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	if len(clicks) == 0 {
		return nil, ErrClicksNotFound
	}
	return clicks, nil
}

func (s *analyticsService) GetByID(ctx context.Context, id uint) (ClickModel, error) {
	click, err := s.Repository.GetByID(ctx, id)
	if err != nil {
		return ClickModel{}, err
	}
	if click.ID == 0 {
		return ClickModel{}, ErrClickNotFound
	}
	return click, nil
}
