package shortener

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Kalmera74/Shorty/internal/types"
	"github.com/Kalmera74/Shorty/pkg/caching"
)

type IShortService interface {
	ShortenURL(ctx context.Context, req ShortenRequest) (ShortModel, error)
	GetById(ctx context.Context, id types.ShortId) (ShortModel, error)
	GetByShortUrl(ctx context.Context, shortUrl string) (ShortModel, error)
	GetByLongUrl(ctx context.Context, originalUrl string) (ShortModel, error)
	Search(ctx context.Context, req SearchRequest) (ShortModel, error)
	GetAllByUser(ctx context.Context, userID types.UserId) ([]ShortModel, error)
	GetAllURLs(ctx context.Context) ([]ShortModel, error)
	DeleteURL(ctx context.Context, shortID types.ShortId) error
}
type shortService struct {
	store  ShortStore
	cacher caching.Cacher
}

func NewShortService(store ShortStore, cacher caching.Cacher) IShortService {
	return &shortService{
		store:  store,
		cacher: cacher,
	}
}

func (s *shortService) ShortenURL(ctx context.Context, req ShortenRequest) (ShortModel, error) {

	cachedShortID, err := s.cacher.Get(ctx, req.Url)
	if err == nil && cachedShortID != "" {
		id, _ := strconv.ParseUint(cachedShortID, 10, 64)
		existingShort, err := s.store.GetById(types.ShortId(id))
		if err == nil {
			return existingShort, nil
		}
	}

	existingShort, err := s.GetByLongUrl(ctx, req.Url)
	if err == nil {
		s.cacher.Set(ctx, req.Url, existingShort.ID, time.Minute*5)
		s.cacher.Set(ctx, existingShort.ShortUrl, existingShort.OriginalUrl, time.Minute*5)
		return existingShort, nil
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

	short, err := s.store.Create(url)

	if err != nil {
		return ShortModel{}, fmt.Errorf("%w: %v", ErrShortenFailed, err)
	}

	s.cacher.Set(ctx, req.Url, short.ID, time.Minute*5)
	marshalledShort, err := json.Marshal(short)
	if err == nil {
		s.cacher.Set(ctx, short.ShortUrl, marshalledShort, time.Minute*5)
	}

	return short, nil

}
func (s *shortService) GetById(ctx context.Context, id types.ShortId) (ShortModel, error) {

	short, err := s.store.GetById(id)
	if err != nil {
		return ShortModel{}, err
	}

	s.cacher.Set(ctx, short.OriginalUrl, short.ID, time.Minute*5)

	return short, nil

}
func (s *shortService) GetByShortUrl(ctx context.Context, shortUrl string) (ShortModel, error) {

	cachedShort, err := s.cacher.Get(ctx, shortUrl)
	if err == nil && cachedShort != "" {

		unMarshalledShort := ShortModel{}
		err := json.Unmarshal([]byte(cachedShort), &unMarshalledShort)
		if err == nil {
			return unMarshalledShort, nil
		}
	}

	short, err := s.store.GetByShortUrl(shortUrl)
	if err != nil {
		return ShortModel{}, fmt.Errorf("%w: %v", ErrShortNotFound, err)
	}

	marshalledShort, err := json.Marshal(short)
	if err == nil {

		s.cacher.Set(ctx, shortUrl, marshalledShort, time.Minute*5)
	}

	return short, nil
}
func (s *shortService) GetByLongUrl(ctx context.Context, originalUrl string) (ShortModel, error) {

	url, err := s.store.GetByLongUrl(originalUrl)
	if err != nil {
		return ShortModel{}, fmt.Errorf("%w: %v", ErrShortNotFound, err)
	}
	return url, nil
}
func (s *shortService) Search(ctx context.Context, req SearchRequest) (ShortModel, error) {
	if req.OriginalUrl != nil {
		url, err := s.store.GetByLongUrl(*req.OriginalUrl)
		if err != nil {
			return ShortModel{}, fmt.Errorf("%w: %v", ErrShortNotFound, err)
		}

		return url, nil
	}

	return ShortModel{}, fmt.Errorf("%w: no search parameters provided", ErrInvalidShortenRequest)
}
func (s *shortService) GetAllByUser(ctx context.Context, userID types.UserId) ([]ShortModel, error) {

	shorts, err := s.store.GetAllByUser(userID)
	if err != nil {
		return nil, err
	}

	return shorts, nil
}
func (s *shortService) GetAllURLs(ctx context.Context) ([]ShortModel, error) {
	allShorts, err := s.store.GetAll()
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
	if err := s.store.Delete(shortID); err != nil {
		return err
	}

	s.cacher.Delete(ctx, short.ShortUrl)
	s.cacher.Delete(ctx, short.OriginalUrl)
	return nil
}

func shortByShortUrlKey(shortUrl string) string {
	return fmt.Sprintf("short:byShortUrl:%s", shortUrl)
}

func shortByOriginalUrlKey(originalUrl string) string {
	return fmt.Sprintf("short:byOriginalUrl:%s", originalUrl)
}
