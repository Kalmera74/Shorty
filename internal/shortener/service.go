package shortener

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Kalmera74/Shorty/validation"

	"github.com/go-redis/redis/v8"
)

type URLService struct {
	store URLStore
	redis *redis.Client
}

func NewURLService(store URLStore, redis *redis.Client) *URLService {
	return &URLService{
		store: store,
		redis: redis,
	}
}

func (s *URLService) ShortenURL(req ShortenRequest) (ShortModel, error) {
	if err := req.Validate(); err != nil {
		return ShortModel{}, err
	}

	ctx := context.Background()
	cachedShortID, err := s.redis.Get(ctx, req.Url).Result()
	if err == nil && cachedShortID != "" {
		id, _ := strconv.ParseUint(cachedShortID, 10, 64)
		existingShort, err := s.store.GetById(uint(id))
		if err == nil {
			return existingShort, nil
		}
	}

	existingShort, err := s.GetByLongUrl(req.Url)
	if err == nil {
		s.redis.Set(ctx, req.Url, existingShort.ID, time.Minute*5)
		s.redis.Set(ctx, existingShort.ShortUrl, existingShort.OriginalUrl, time.Minute*5)
		return existingShort, nil
	}

	h := sha1.New()
	h.Write([]byte(req.Url))
	shortID := hex.EncodeToString(h.Sum(nil))[:8]

	url := ShortModel{
		UserID:      req.UserID,
		OriginalUrl: req.Url,
		ShortUrl:    shortID,
	}

	short, err := s.store.Create(url)

	if err != nil {
		return ShortModel{}, &ShortenError{Msg: fmt.Sprintf("Could not create the shortened url. Reason: %v", err.Error()), Err: err}
	}

	s.redis.Set(ctx, req.Url, short.ID, 0)
	marshalledShort, err := json.Marshal(short)
	if err == nil {
		s.redis.Set(ctx, short.ShortUrl, marshalledShort, time.Minute*5)
	}

	return short, nil

}

func (s *URLService) GetById(id uint) (ShortModel, error) {
	if err := validation.ValidateID(id); err != nil {
		return ShortModel{}, err
	}

	short, err := s.store.GetById(id)
	if err != nil {
		return ShortModel{}, err
	}

	ctx := context.Background()
	s.redis.Set(ctx, short.OriginalUrl, short.ID, time.Minute*5)

	return short, nil

}

func (s *URLService) GetByShortUrl(shortUrl string) (ShortModel, error) {
	if shortUrl == "" {
		return ShortModel{}, &ShortenError{Msg: "Short url cannot be nil or empty"}
	}

	ctx := context.Background()
	cachedShort, err := s.redis.Get(ctx, shortUrl).Result()
	if err == nil && cachedShort != "" {

		unMarshalledShort := ShortModel{}
		err := json.Unmarshal([]byte(cachedShort), &unMarshalledShort)
		if err == nil {
			return unMarshalledShort, nil
		}
	}

	short, err := s.store.GetByShortUrl(shortUrl)
	if err != nil {
		if errors.Is(err, &ShortNotFoundError{}) {
			return ShortModel{}, err
		}
		return ShortModel{}, &ShortenError{Msg: fmt.Sprintf("Could not get the short with the Id: %v Reason: %v", shortUrl, err.Error()), Err: err}
	}

	marshalledShort, err := json.Marshal(short)
	if err == nil {

		s.redis.Set(ctx, shortUrl, marshalledShort, time.Minute*5)
	}
	return short, nil
}
func (s *URLService) GetByLongUrl(originalUrl string) (ShortModel, error) {
	if err := validation.ValidateUrl(originalUrl); err != nil {
		return ShortModel{}, err
	}

	url, err := s.store.GetByLongUrl(originalUrl)
	if err != nil {
		if errors.Is(err, &ShortNotFoundError{}) {
			return ShortModel{}, err
		}
		return ShortModel{}, &ShortenError{Msg: fmt.Sprintf("Could not get the short with the original Url: %v Reason: %v", originalUrl, err.Error()), Err: err}
	}

	return url, nil
}

func (s *URLService) GetAllByUser(userID uint) ([]ShortModel, error) {
	if err := validation.ValidateID(userID); err != nil {
		return nil, err
	}

	shorts, err := s.store.GetAllByUser(userID)
	if err != nil {
		return nil, err
	}

	return shorts, nil
}

func (s *URLService) GetAllURLs() ([]ShortModel, error) {
	allShorts, err := s.store.GetAll()
	if err != nil {
		return nil, err
	}
	return allShorts, nil
}

func (s *URLService) DeleteURL(shortID uint) error {
	if err := validation.ValidateID(shortID); err != nil {
		return err
	}

	if err := s.store.Delete(shortID); err != nil {
		return err
	}

	return nil
}
