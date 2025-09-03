package analytics

import "fmt"

type IAnalyticsService interface {
	GetAll() ([]ClickModel, error)
	GetAllByShortUrl(shortUrl string) ([]ClickModel, error)
	Create(click ClickModel) (ClickModel, error)
}

type analyticsService struct {
	Repository IClickRepository
}

func NewAnalyticService(p IClickRepository) IAnalyticsService {
	return &analyticsService{p}
}

func (s *analyticsService) Create(click ClickModel) (ClickModel, error) {
	//TODO: Publish a rabbit message to handle the clicks data
	createdClick, err := s.Repository.Create(click)
	if err != nil {
		return ClickModel{}, fmt.Errorf("%w: %v", ErrClickCreateFail, err)
	}
	return createdClick, nil
}

func (s *analyticsService) GetAll() ([]ClickModel, error) {
	clicks, err := s.Repository.GetAll()
	if err != nil {
		return nil, err
	}
	if len(clicks) == 0 {
		return nil, ErrClickNotFound
	}
	return clicks, nil
}

func (s *analyticsService) GetAllByShortUrl(shortUrl string) ([]ClickModel, error) {
	clicks, err := s.Repository.GetAllByShortUrl(shortUrl)
	if err != nil {
		return nil, err
	}
	if len(clicks) == 0 {
		return nil, fmt.Errorf("%w: no clicks found for short url %s", ErrClickNotFound, shortUrl)
	}
	return clicks, nil
}
