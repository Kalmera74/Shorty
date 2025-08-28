package stores

import (
	"errors"
	"fmt"

	"github.com/Kalmera74/Shorty/internal/shortener"
	"github.com/Kalmera74/Shorty/validation"
	"gorm.io/gorm"
)

type PostgresURLStore struct {
	db *gorm.DB
}

func NewURLPostgresStore(db *gorm.DB) *PostgresURLStore {
	return &PostgresURLStore{db: db}
}

func (s *PostgresURLStore) Create(short shortener.ShortModel) (shortener.ShortModel, error) {
	if err := short.Validate(); err != nil {
		return shortener.ShortModel{}, err

	}

	result := s.db.Create(&short)
	if result.Error != nil {
		return shortener.ShortModel{}, &shortener.ShortenError{Msg: fmt.Sprintf("Could not create shortened url. Reason: %v", result.Error.Error()), Err: result.Error}
	}
	return short, nil
}

func (s *PostgresURLStore) GetById(shortID uint) (shortener.ShortModel, error) {
	var url shortener.ShortModel

	result := s.db.First(&url, shortID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return shortener.ShortModel{}, &shortener.ShortNotFoundError{Msg: fmt.Sprintf("No short with the given Id: %v found", shortID)}
	}

	if result.Error != nil {
		return shortener.ShortModel{}, result.Error
	}

	return url, nil
}

func (s *PostgresURLStore) GetByShortUrl(shortUrl string) (shortener.ShortModel, error) {

	if err := validation.ValidateUrl(shortUrl); err != nil {
		return shortener.ShortModel{}, &shortener.ShortenError{Msg: fmt.Sprintf("Could not get the short with the short url: %v. Reason: %v", shortUrl, err.Error()), Err: err}
	}

	var short shortener.ShortModel
	result := s.db.Where("short_url = ?", shortUrl).First(&short)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return shortener.ShortModel{}, &shortener.ShortNotFoundError{
			Msg: "URL not found",
			Err: result.Error,
		}
	}
	if result.Error != nil {
		return shortener.ShortModel{}, &shortener.ShortenError{Msg: fmt.Sprintf("Could not create shortened url. Reason: %v", result.Error.Error()), Err: result.Error}
	}
	return short, nil
}
func (s *PostgresURLStore) GetByLongUrl(originalUrl string) (shortener.ShortModel, error) {
	if err := validation.ValidateUrl(originalUrl); err != nil {
		return shortener.ShortModel{}, &shortener.ShortenError{Msg: fmt.Sprintf("Could not get the short with the Url: %v. Reason: %v", originalUrl, err.Error()), Err: err}
	}

	var short shortener.ShortModel
	result := s.db.Where("original_url = ?", originalUrl).First(&short)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return shortener.ShortModel{}, &shortener.ShortNotFoundError{
			Msg: "URL not found",
			Err: result.Error,
		}
	}
	if result.Error != nil {
		return shortener.ShortModel{}, &shortener.ShortenError{Msg: fmt.Sprintf("Could not create shortened url. Reason: %v", result.Error.Error()), Err: result.Error}
	}
	return short, nil
}

func (s *PostgresURLStore) GetAllByUser(userID uint) ([]shortener.ShortModel, error) {
	if err := validation.ValidateID(userID); err != nil {
		return nil, &shortener.ShortenError{Msg: fmt.Sprintf("Could not get the shortened Urls by the user with the Id: %v. Reason: %v", userID, err.Error()), Err: err}
	}

	var urls []shortener.ShortModel
	result := s.db.Where("user_id = ?", userID).Find(&urls)

	if result.Error != nil {

		return nil, &shortener.ShortenError{Msg: fmt.Sprintf("Could not get the shortened Urls by the user with the Id: %v. Reason: %v", userID, result.Error.Error()), Err: result.Error}
	}
	if len(urls) == 0 {
		return nil, &shortener.ShortNotFoundError{
			Msg: "No URLs found for this user",
			Err: nil,
		}
	}
	return urls, nil
}

func (s *PostgresURLStore) GetAll() ([]shortener.ShortModel, error) {
	var urls []shortener.ShortModel
	result := s.db.Find(&urls)
	if result.Error != nil {
		return nil, &shortener.ShortenError{Msg: fmt.Sprintf("failed to get all URLs. Reason: %v", result.Error), Err: result.Error}
	}
	if len(urls) == 0 {
		return nil, &shortener.ShortNotFoundError{
			Msg: "No URLs found",
			Err: nil,
		}
	}
	return urls, nil
}

func (s *PostgresURLStore) Delete(shortId uint) error {
	if err := validation.ValidateID(shortId); err != nil {
		return err
	}

	result := s.db.Find(shortId).Delete(&shortener.ShortModel{})

	if result.Error != nil {
		return &shortener.ShortenError{Msg: fmt.Sprintf("Failed to delete short with the Id: %v. Reason: %v", shortId, result.Error.Error()), Err: result.Error}
	}

	if result.RowsAffected == 0 {
		return &shortener.ShortNotFoundError{Msg: fmt.Sprintf("No short is found with the given Id: %v", shortId), Err: gorm.ErrRecordNotFound}
	}

	return nil
}
