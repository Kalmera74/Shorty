package stores

import (
	"errors"
	"fmt"

	"github.com/Kalmera74/Shorty/internal/shortener"
	"gorm.io/gorm"
)

type PostgresURLStore struct {
	db *gorm.DB
}

func NewURLPostgresStore(db *gorm.DB) *PostgresURLStore {
	return &PostgresURLStore{db: db}
}

func (s *PostgresURLStore) Create(url shortener.ShortenModel) (shortener.ShortenModel, error) {
	if url.LongURL == "" || url.ShortID == "" {
		return shortener.ShortenModel{}, errors.New("invalid URL data")
	}

	result := s.db.Create(&url)
	if result.Error != nil {
		return shortener.ShortenModel{}, fmt.Errorf("failed to create URL: %w", result.Error)
	}
	return url, nil
}

func (s *PostgresURLStore) GetByShortID(shortID string) (shortener.ShortenModel, error) {
	if shortID == "" {
		return shortener.ShortenModel{}, errors.New("short ID cannot be empty")
	}

	var url shortener.ShortenModel
	result := s.db.Where("short_id = ?", shortID).First(&url)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return shortener.ShortenModel{}, &shortener.URLNotFoundError{
			Msg: "URL not found",
			Err: result.Error,
		}
	}
	if result.Error != nil {
		return shortener.ShortenModel{}, fmt.Errorf("failed to get URL by short ID: %w", result.Error)
	}
	return url, nil
}

func (s *PostgresURLStore) GetAllByUser(userID uint) ([]shortener.ShortenModel, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	var urls []shortener.ShortenModel
	result := s.db.Where("user_id = ?", userID).Find(&urls)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get URLs for user: %w", result.Error)
	}
	if len(urls) == 0 {
		return nil, &shortener.URLNotFoundError{
			Msg: "No URLs found for this user",
			Err: nil,
		}
	}
	return urls, nil
}

func (s *PostgresURLStore) GetAll() ([]shortener.ShortenModel, error) {
	var urls []shortener.ShortenModel
	result := s.db.Find(&urls)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get all URLs: %w", result.Error)
	}
	if len(urls) == 0 {
		return nil, &shortener.URLNotFoundError{
			Msg: "No URLs found",
			Err: nil,
		}
	}
	return urls, nil
}

func (s *PostgresURLStore) Delete(shortID string) error {
	if shortID == "" {
		return errors.New("short ID cannot be empty")
	}

	result := s.db.Where("short_id = ?", shortID).Delete(&shortener.ShortenModel{})
	if result.RowsAffected == 0 {
		return &shortener.URLNotFoundError{
			Msg: "URL not found",
			Err: gorm.ErrRecordNotFound,
		}
	}
	if result.Error != nil {
		return fmt.Errorf("failed to delete URL: %w", result.Error)
	}
	return nil
}
