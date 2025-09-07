package shortener

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Kalmera74/Shorty/internal/types"
	"github.com/Kalmera74/Shorty/pkg/cache"
	"gorm.io/gorm"
)

type IShortService interface {
	ShortenURL(ctx context.Context, req ShortenRequest) (ShortModel, error)
	GetById(ctx context.Context, id types.ShortId) (ShortModel, error)
	GetByShortUrl(ctx context.Context, shortUrl string) (ShortModel, error)
	GetByLongUrl(ctx context.Context, originalUrl string) (ShortModel, error)
	Search(ctx context.Context, req SearchRequest) ([]ShortModel, error)
	GetAllByUser(ctx context.Context, userID types.UserId) ([]ShortModel, error)
	GetAll(ctx context.Context) ([]ShortModel, error)
	DeleteURL(ctx context.Context, shortID types.ShortId) error
}
type shortService struct {
	Repository IShortRepository
	Cacher     caching.ICacher
}

func NewShortService(store IShortRepository, cacher caching.ICacher) IShortService {
	return &shortService{
		Repository: store,
		Cacher:     cacher,
	}
}

func (s *shortService) ShortenURL(ctx context.Context, req ShortenRequest) (ShortModel, error) {

	search := SearchRequest{
		UserId:      &req.UserID,
		OriginalUrl: &req.Url,
	}

	searchResult, err := s.Search(ctx, search)
	if err == nil {
		short := searchResult[0]
		return short, nil
	}

	var shortID string

	h := sha1.New()
	h.Write([]byte(req.Url))
	shortID = hex.EncodeToString(h.Sum(nil))[:8]

	url := ShortModel{
		UserID:      types.UserId(req.UserID),
		OriginalUrl: req.Url,
		ShortUrl:    shortID,
	}

	short, err := s.Repository.Create(ctx, url)

	if err != nil {
		return ShortModel{}, fmt.Errorf("%w: %v", ErrShortenFailed, err)
	}

	return short, nil

}
func (s *shortService) GetById(ctx context.Context, id types.ShortId) (ShortModel, error) {

	short, err := s.Repository.GetById(ctx, id)
	if err != nil {
		return ShortModel{}, err
	}

	return short, nil

}
func (s *shortService) GetByShortUrl(ctx context.Context, shortUrl string) (ShortModel, error) {

	cachedShort, err := s.Cacher.Get(ctx, shortUrl)
	if err == nil && cachedShort != "" {
		unMarshalledShort := ShortModel{}
		err := json.Unmarshal([]byte(cachedShort), &unMarshalledShort)
		if err == nil {
			return unMarshalledShort, nil
		}
	}

	search := SearchRequest{
		ShortUrl: &shortUrl,
	}
	result, err := s.Repository.Search(ctx, search)
	if err != nil {
		return ShortModel{}, fmt.Errorf("%w: %v", ErrShortNotFound, err)
	}

	short := result[0]
	marshalledShort, err := json.Marshal(short)
	if err == nil {
		s.Cacher.Set(ctx, shortByShortUrlKey(shortUrl), marshalledShort, time.Minute*5)
	}

	return short, nil
}
func (s *shortService) GetByLongUrl(ctx context.Context, originalUrl string) (ShortModel, error) {

	search := SearchRequest{
		OriginalUrl: &originalUrl,
	}
	result, err := s.Repository.Search(ctx, search)
	if err != nil {
		return ShortModel{}, fmt.Errorf("%w: %v", ErrShortNotFound, err)
	}

	original := result[0]

	return original, nil
}

func (s *shortService) Search(ctx context.Context, req SearchRequest) ([]ShortModel, error) {
	result, err := s.Repository.Search(ctx, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w", ErrShortNotFound)
		}
		return nil, err
	}

	return result, nil
}

func (s *shortService) GetAllByUser(ctx context.Context, userID types.UserId) ([]ShortModel, error) {

	search := SearchRequest{
		UserId: &userID,
	}
	result, err := s.Repository.Search(ctx, search)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrShortNotFound, err)
	}

	return result, nil
}
func (s *shortService) GetAll(ctx context.Context) ([]ShortModel, error) {
	allShorts, err := s.Repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return allShorts, nil
}
func (s *shortService) DeleteURL(ctx context.Context, shortID types.ShortId) error {

	short, err := s.GetById(ctx, shortID)
	if err != nil {
		return err
	}
	if err := s.Repository.Delete(ctx, shortID); err != nil {
		return err
	}

	s.Cacher.Delete(ctx, short.ShortUrl)
	s.Cacher.Delete(ctx, short.OriginalUrl)
	return nil
}

func shortByShortUrlKey(shortUrl string) string {
	return fmt.Sprintf("short:byShortUrl:%s", shortUrl)
}
