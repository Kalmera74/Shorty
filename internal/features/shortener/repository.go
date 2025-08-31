package shortener

import (
	"errors"
	"fmt"

	"github.com/Kalmera74/Shorty/internal/types"
	"github.com/Kalmera74/Shorty/internal/validation"
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
		return ShortModel{}, &ShortenError{Msg: fmt.Sprintf("Could not create shortened url. Reason: %v", result.Error.Error()), Err: result.Error}
	}
	return short, nil
}
func (s *PostgresURLStore) GetById(shortID types.ShortId) (ShortModel, error) {
	var url ShortModel

	result := s.db.First(&url, shortID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ShortModel{}, &ShortNotFoundError{Msg: fmt.Sprintf("No short with the given Id: %v found", shortID)}
	}

	if result.Error != nil {
		return ShortModel{}, result.Error
	}

	return url, nil
}
func (s *PostgresURLStore) GetByShortUrl(shortUrl string) (ShortModel, error) {

	if err := validation.ValidateUrl(shortUrl); err != nil {
		return ShortModel{}, &ShortenError{Msg: fmt.Sprintf("Could not get the short with the short url: %v. Reason: %v", shortUrl, err.Error()), Err: err}
	}

	var short ShortModel
	result := s.db.Where("short_url = ?", shortUrl).First(&short)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ShortModel{}, &ShortNotFoundError{
			Msg: "URL not found",
			Err: result.Error,
		}
	}
	if result.Error != nil {
		return ShortModel{}, &ShortenError{Msg: fmt.Sprintf("Could not create shortened url. Reason: %v", result.Error.Error()), Err: result.Error}
	}
	return short, nil
}
func (s *PostgresURLStore) GetByLongUrl(originalUrl string) (ShortModel, error) {


	var short ShortModel
	result := s.db.Where("original_url = ?", originalUrl).First(&short)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ShortModel{}, &ShortNotFoundError{
			Msg: "URL not found",
			Err: result.Error,
		}
	}
	if result.Error != nil {
		return ShortModel{}, &ShortenError{Msg: fmt.Sprintf("Could not create shortened url. Reason: %v", result.Error.Error()), Err: result.Error}
	}
	return short, nil
}
func (s *PostgresURLStore) GetAllByUser(userID types.UserId) ([]ShortModel, error) {
	var urls []ShortModel
	result := s.db.Where("user_id = ?", userID).Find(&urls)

	if result.Error != nil {

		return nil, &ShortenError{Msg: fmt.Sprintf("Could not get the shortened Urls by the user with the Id: %v. Reason: %v", userID, result.Error.Error()), Err: result.Error}
	}
	if len(urls) == 0 {
		return nil, &ShortNotFoundError{
			Msg: "No URLs found for this user",
			Err: nil,
		}
	}
	return urls, nil
}
func (s *PostgresURLStore) GetAll() ([]ShortModel, error) {
	var urls []ShortModel
	result := s.db.Find(&urls)
	if result.Error != nil {
		return nil, &ShortenError{Msg: fmt.Sprintf("failed to get all URLs. Reason: %v", result.Error), Err: result.Error}
	}
	if len(urls) == 0 {
		return nil, &ShortNotFoundError{
			Msg: "No URLs found",
			Err: nil,
		}
	}
	return urls, nil
}
func (s *PostgresURLStore) Delete(shortId types.ShortId) error {
	result := s.db.Find(shortId).Delete(&ShortModel{})

	if result.Error != nil {
		return &ShortenError{Msg: fmt.Sprintf("Failed to delete short with the Id: %v. Reason: %v", shortId, result.Error.Error()), Err: result.Error}
	}

	if result.RowsAffected == 0 {
		return &ShortNotFoundError{Msg: fmt.Sprintf("No short is found with the given Id: %v", shortId), Err: gorm.ErrRecordNotFound}
	}

	return nil
}
