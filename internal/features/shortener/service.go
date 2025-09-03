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
	"github.com/Kalmera74/Shorty/pkg/redis"
)

type IShortService interface {
	ShortenURL(req ShortenRequest) (ShortModel, error)
	GetById(id types.ShortId) (ShortModel, error)
	GetByShortUrl(shortUrl string) (ShortModel, error)
	GetByLongUrl(originalUrl string) (ShortModel, error)
	Search(req SearchRequest) (ShortModel, error)
	GetAllByUser(userID types.UserId) ([]ShortModel, error)
	GetAllURLs() ([]ShortModel, error)
	DeleteURL(shortID types.ShortId) error
}
type shortService struct {
	store  ShortStore
	cacher redis.Cacher
}

func NewShortService(store ShortStore, cacher redis.Cacher) IShortService {
	return &shortService{
		store:  store,
		cacher: cacher,
	}
}

func (s *shortService) ShortenURL(req ShortenRequest) (ShortModel, error) {

	//TODO: Update the user model to include the crated short
	ctx := context.Background()
	cachedShortID, err := s.cacher.Get(ctx, req.Url)
	if err == nil && cachedShortID != "" {
		id, _ := strconv.ParseUint(cachedShortID, 10, 64)
		existingShort, err := s.store.GetById(types.ShortId(id))
		if err == nil {
			return existingShort, nil
		}
	}

	existingShort, err := s.GetByLongUrl(req.Url)
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
func (s *shortService) GetById(id types.ShortId) (ShortModel, error) {

	short, err := s.store.GetById(id)
	if err != nil {
		return ShortModel{}, err
	}

	ctx := context.Background()
	s.cacher.Set(ctx, short.OriginalUrl, short.ID, time.Minute*5)

	return short, nil

}
func (s *shortService) GetByShortUrl(shortUrl string) (ShortModel, error) {

	ctx := context.Background()
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
func (s *shortService) GetByLongUrl(originalUrl string) (ShortModel, error) {

	url, err := s.store.GetByLongUrl(originalUrl)
	if err != nil {
		return ShortModel{}, fmt.Errorf("%w: %v", ErrShortNotFound, err)
	}
	return url, nil
}
func (s *shortService) Search(req SearchRequest) (ShortModel, error) {
	if req.OriginalUrl != nil {
		url, err := s.store.GetByLongUrl(*req.OriginalUrl)
		if err != nil {
			return ShortModel{}, fmt.Errorf("%w: %v", ErrShortNotFound, err)
		}

		return url, nil
	}

	return ShortModel{}, fmt.Errorf("%w: no search parameters provided", ErrInvalidShortenRequest)
}
func (s *shortService) GetAllByUser(userID types.UserId) ([]ShortModel, error) {

	shorts, err := s.store.GetAllByUser(userID)
	if err != nil {
		return nil, err
	}

	return shorts, nil
}
func (s *shortService) GetAllURLs() ([]ShortModel, error) {
	allShorts, err := s.store.GetAll()
	if err != nil {
		return nil, err
	}
	return allShorts, nil
}
func (s *shortService) DeleteURL(shortID types.ShortId) error {

	if err := s.store.Delete(shortID); err != nil {
		return err
	}

	return nil
}
