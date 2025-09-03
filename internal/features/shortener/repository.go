package shortener

import (
	"errors"
	"fmt"

	"github.com/Kalmera74/Shorty/internal/types"
	"gorm.io/gorm"
)

type ShortStore interface {
	Create(short ShortModel) (ShortModel, error)
	GetById(id types.ShortId) (ShortModel, error)
	GetByShortUrl(shortUrl string) (ShortModel, error)
	GetByLongUrl(originalUrl string) (ShortModel, error)
	GetAllByUser(userID types.UserId) ([]ShortModel, error)
	GetAll() ([]ShortModel, error)
	Delete(shortenID types.ShortId) error
}

type PostgresURLStore struct {
	db *gorm.DB
}

func NewShortRepository(db *gorm.DB) *PostgresURLStore {
	return &PostgresURLStore{db: db}
}

func (s *PostgresURLStore) Create(short ShortModel) (ShortModel, error) {

	result := s.db.Create(&short)
	if result.Error != nil {
		return ShortModel{}, fmt.Errorf("%w: %v", ErrShortenFailed, result.Error)
	}
	return short, nil
}
func (s *PostgresURLStore) GetById(shortID types.ShortId) (ShortModel, error) {
	var url ShortModel

	result := s.db.First(&url, shortID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ShortModel{}, fmt.Errorf("%w: %v", ErrShortNotFound, result.Error)
	}

	if result.Error != nil {
		return ShortModel{}, result.Error
	}

	return url, nil
}
func (s *PostgresURLStore) GetByShortUrl(shortUrl string) (ShortModel, error) {

	var short ShortModel
	result := s.db.Where("short_url = ?", shortUrl).First(&short)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ShortModel{}, fmt.Errorf("%w: %v", ErrShortNotFound, result.Error)
	}
	if result.Error != nil {
		return ShortModel{}, fmt.Errorf("%w: %v", ErrShortenFailed, result.Error)
	}
	return short, nil
}
func (s *PostgresURLStore) GetByLongUrl(originalUrl string) (ShortModel, error) {

	var short ShortModel
	result := s.db.Where("original_url = ?", originalUrl).First(&short)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ShortModel{}, fmt.Errorf("%w: %v", ErrShortNotFound, result.Error)
	}
	if result.Error != nil {
		return ShortModel{}, fmt.Errorf("%w: %v", ErrShortenFailed, result.Error)
	}
	return short, nil
}
func (s *PostgresURLStore) GetAllByUser(userID types.UserId) ([]ShortModel, error) {
	var urls []ShortModel
	result := s.db.Where("user_id = ?", userID).Find(&urls)

	if result.Error != nil {

		return nil, fmt.Errorf("%w: %v", ErrShortenFailed, result.Error)
	}
	if len(urls) == 0 {
		return nil, fmt.Errorf("%w: no URLs found for user %d", ErrShortNotFound, userID)
	}
	return urls, nil
}
func (s *PostgresURLStore) GetAll() ([]ShortModel, error) {
	var urls []ShortModel
	result := s.db.Find(&urls)
	if result.Error != nil {
		return nil, fmt.Errorf("%w: %v", ErrShortenFailed, result.Error)
	}
	if len(urls) == 0 {
		return nil, fmt.Errorf("%w: no URLs found", ErrShortNotFound)
	}
	return urls, nil
}
func (s *PostgresURLStore) Delete(shortId types.ShortId) error {
	result := s.db.Find(shortId).Delete(&ShortModel{})

	if result.Error != nil {
		return fmt.Errorf("%w: %v", ErrShortDeleteFail, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: no URL found with ID %d", ErrShortNotFound, shortId)
	}

	return nil
}
