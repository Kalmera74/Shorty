package shortener

import (
	"context"
	"errors"
	"fmt"

	"github.com/Kalmera74/Shorty/internal/types"
	"gorm.io/gorm"
)

type IShortRepository interface {
	Create(ctx context.Context, short ShortModel) (ShortModel, error)
	GetById(ctx context.Context, id types.ShortId) (ShortModel, error)
	Search(ctx context.Context, req SearchRequest) ([]ShortModel, error)
	GetAll(ctx context.Context, offset, limit int) ([]ShortModel, int, error)
	Delete(ctx context.Context, shortenID types.ShortId) error
}

type postgresURLStore struct {
	db *gorm.DB
}

func NewShortRepository(db *gorm.DB) IShortRepository {
	return &postgresURLStore{db: db}
}

func (s *postgresURLStore) Create(ctx context.Context, short ShortModel) (ShortModel, error) {

	result := s.db.Create(&short)
	if result.Error != nil {
		return ShortModel{}, fmt.Errorf("%w: %v", ErrShortenFailed, result.Error)
	}
	return short, nil
}
func (s *postgresURLStore) GetById(ctx context.Context, shortID types.ShortId) (ShortModel, error) {
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
func (r *postgresURLStore) Search(ctx context.Context, req SearchRequest) ([]ShortModel, error) {
	var shorts []ShortModel
	query := r.db.WithContext(ctx).Model(&ShortModel{})

	if req.OriginalUrl != nil {
		query = query.Where("original_url = ?", *req.OriginalUrl)
	}
	if req.UserId != nil {
		query = query.Where("user_id = ?", *req.UserId)
	}
	if req.ShortUrl != nil {
		query = query.Where("short_url =?", *req.ShortUrl)
	}

	if err := query.Find(&shorts).Error; err != nil {
		return nil, err
	}

	if len(shorts) == 0 {
		return nil, fmt.Errorf("%w", ErrShortNotFound)
	}

	return shorts, nil
}
func (r *postgresURLStore) GetAll(ctx context.Context, offset, limit int) ([]ShortModel, int, error) {
	var shorts []ShortModel
	var total int64

	if err := r.db.WithContext(ctx).Model(&ShortModel{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrShortNotFound, err)
	}

	if err := r.db.WithContext(ctx).
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&shorts).Error; err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrShortNotFound, err)
	}

	if len(shorts) == 0 {
		return nil, -1, fmt.Errorf("%w", ErrShortNotFound)
	}

	return shorts, int(total), nil
}

func (s *postgresURLStore) Delete(ctx context.Context, shortId types.ShortId) error {
	result := s.db.Find(shortId).Delete(&ShortModel{})

	if result.Error != nil {
		return fmt.Errorf("%w: %v", ErrShortDeleteFail, result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%w: no URL found with ID %d", ErrShortNotFound, shortId)
	}

	return nil
}
