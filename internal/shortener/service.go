package shortener

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
)

type URLService struct {
	store URLStore
}

func NewURLService(store URLStore) *URLService {
	return &URLService{store: store}
}

func (s *URLService) ShortenURL(userID uint, longURL string) (string, error) {
	if longURL == "" {
		return "", errors.New("long URL cannot be empty")
	}

	h := sha1.New()
	h.Write([]byte(longURL))
	shortID := hex.EncodeToString(h.Sum(nil))[:8]

	url := ShortenModel{
		UserID:  userID,
		LongURL: longURL,
		ShortID: shortID,
	}

	_, err := s.store.Create(url)
	if err != nil {
		return "", fmt.Errorf("failed to create short URL: %w", err)
	}

	return shortID, nil
}

func (s *URLService) GetURLByShortID(shortID string) (ShortenModel, error) {
	if shortID == "" {
		return ShortenModel{}, errors.New("short ID cannot be empty")
	}

	url, err := s.store.GetByShortID(shortID)
	if err != nil {
		return ShortenModel{}, fmt.Errorf("failed to get URL by short ID: %w", err)
	}

	return url, nil
}

func (s *URLService) GetAllURLsByUser(userID uint) ([]ShortenModel, error) {
	if userID == 0 {
		return nil, errors.New("invalid user ID")
	}

	urls, err := s.store.GetAllByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get URLs for user: %w", err)
	}

	return urls, nil
}

func (s *URLService) GetAllURLs() ([]ShortenModel, error) {
	urls, err := s.store.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get all URLs: %w", err)
	}
	return urls, nil
}

func (s *URLService) DeleteURL(shortID string) error {
	if shortID == "" {
		return errors.New("short ID cannot be empty")
	}

	if err := s.store.Delete(shortID); err != nil {
		return fmt.Errorf("failed to delete URL: %w", err)
	}

	return nil
}
