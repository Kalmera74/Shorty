package shortener

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/Kalmera74/Shorty/internal/types"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Store //
type MockStore struct {
	mock.Mock
}

func (m *MockStore) Create(ctx context.Context, url ShortModel) (ShortModel, error) {
	args := m.Called(url)
	return args.Get(0).(ShortModel), args.Error(1)
}
func (m *MockStore) GetById(ctx context.Context, id types.ShortId) (ShortModel, error) {
	args := m.Called(id)
	return args.Get(0).(ShortModel), args.Error(1)
}

func (m *MockStore) Search(ctx context.Context, search SearchRequest) ([]ShortModel, error) {
	args := m.Called(search)
	return args.Get(0).([]ShortModel), args.Error(1)
}
func (m *MockStore) GetAll(ctx context.Context, offset, limit int) ([]ShortModel, int, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]ShortModel), args.Int(1), args.Error(2)
}
func (m *MockStore) Delete(ctx context.Context, id types.ShortId) error {
	args := m.Called(id)
	return args.Error(0)
}

// Mock Redis //
type MockRedis struct {
	mock.Mock
}

func (m *MockRedis) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	val, _ := args.Get(0).(string)
	err, _ := args.Get(1).(error)
	return val, err
}

func (m *MockRedis) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	m.Called(ctx, key, value, ttl)
	return nil
}

func (m *MockRedis) Delete(ctx context.Context, key string) error {
	m.Called(ctx, key)
	return nil
}

// ShortenURL Tests //
func TestShortenURL_Hit_DB(t *testing.T) {
	mockStore := new(MockStore)

	req := ShortenRequest{
		UserID: 1,
		Url:    "https://example.com",
	}

	expectedShort := []ShortModel{
		{
			ID:          1,
			UserID:      1,
			OriginalUrl: "https://example.com",
			ShortUrl:    "12345678",
		},
	}

	mockStore.On("Search", mock.Anything, mock.Anything).Return(expectedShort, nil)

	service := NewShortService(mockStore, nil)
	result, err := service.ShortenURL(nil, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedShort[0].ID, result.ID)

	mockStore.AssertExpectations(t)
}

func TestShortURL_Miss_DB(t *testing.T) {
	mockStore := new(MockStore)

	req := ShortenRequest{
		UserID: 1,
		Url:    "https://example.com",
	}

	expectedShort := ShortModel{
		ID:          1,
		UserID:      1,
		OriginalUrl: "https://example.com",
		ShortUrl:    "12345678",
	}

	mockStore.On("Search", mock.Anything, mock.Anything).Return([]ShortModel{}, ErrShortNotFound)
	mockStore.On("Create", mock.Anything, mock.Anything).Return(expectedShort, nil)

	service := NewShortService(mockStore, nil)
	result, err := service.ShortenURL(nil, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedShort.ID, result.ID)

	mockStore.AssertExpectations(t)
}

// GetById Tests
func TestGetById_ValidId(t *testing.T) {
	mockStore := new(MockStore)

	expectedShort := ShortModel{
		ID:          1,
		UserID:      1,
		OriginalUrl: "https://example.com",
		ShortUrl:    "12345678",
	}

	mockStore.On("GetById", types.ShortId(1)).Return(expectedShort, nil)

	service := NewShortService(mockStore, nil)
	result, err := service.GetById(nil, types.ShortId(1))

	assert.Equal(t, expectedShort, result)
	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
}

func TestGetById_StoreFails(t *testing.T) {
	mockStore := new(MockStore)
	mockRedis := new(MockRedis)

	mockStore.On("GetById", mock.Anything).Return(ShortModel{}, ErrShortNotFound)

	service := NewShortService(mockStore, mockRedis)
	_, err := service.GetById(nil, types.ShortId(999))

	assert.Error(t, err)
	assert.IsType(t, ErrShortNotFound, err)
	mockStore.AssertExpectations(t)
	mockRedis.AssertNotCalled(t, "Set")
}

// GetByShortUrl Tests
func TestGetByShortUrl_EmptyUrl(t *testing.T) {
	mockStore := new(MockStore)
	mockRedis := new(MockRedis)

	service := NewShortService(mockStore, mockRedis)

	mockStore.On("GetByShortUrl", mock.Anything).Return(ShortModel{}, ErrShortNotFound)
	mockRedis.On("Get", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("", ErrShortNotFound)
	mockStore.On("Search", mock.Anything, mock.Anything).Return([]ShortModel{}, ErrShortNotFound)

	_, err := service.GetByShortUrl(nil, "")

	assert.Error(t, err)
}

func TestGetByShortUrl_CacheHit(t *testing.T) {
	mockStore := new(MockStore)
	mockRedis := new(MockRedis)

	expectedShort := ShortModel{
		ID:          1,
		UserID:      1,
		OriginalUrl: "https://example.com",
		ShortUrl:    "12345678",
	}
	marshalledShort, _ := json.Marshal(expectedShort)

	mockRedis.On("Get", mock.Anything, expectedShort.ShortUrl).Return(string(marshalledShort), nil)
	mockStore.AssertNotCalled(t, "GetByShortUrl")

	service := NewShortService(mockStore, mockRedis)
	result, err := service.GetByShortUrl(nil, expectedShort.ShortUrl)

	assert.NoError(t, err)
	assert.Equal(t, expectedShort, result)
	mockRedis.AssertExpectations(t)
}

func TestGetByShortUrl_CacheMiss_DBHit(t *testing.T) {
	mockStore := new(MockStore)
	mockRedis := new(MockRedis)

	expectedShort := []ShortModel{
		{
			ID:          1,
			UserID:      1,
			OriginalUrl: "https://example.com",
			ShortUrl:    "12345678",
		},
	}

	mockRedis.On("Get", mock.Anything, mock.Anything).Return("", redis.Nil)
	mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockStore.On("Search", mock.Anything, mock.Anything).Return(expectedShort, nil)

	service := NewShortService(mockStore, mockRedis)
	result, err := service.GetByShortUrl(nil, "")

	assert.NoError(t, err)
	assert.Equal(t, expectedShort[0], result)
	mockStore.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
}

func TestGetByShortUrl_CacheMiss_DBMiss(t *testing.T) {
	mockStore := new(MockStore)
	mockRedis := new(MockRedis)

	mockRedis.On("Get", mock.Anything, mock.Anything).Return("", redis.Nil)
	mockStore.On("Search", mock.Anything, mock.Anything).Return([]ShortModel{}, ErrShortNotFound)

	service := NewShortService(mockStore, mockRedis)
	_, err := service.GetByShortUrl(nil, "nonexistent")

	assert.Error(t, err)
	mockStore.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
}

// GetByLongUrl Tests
func TestGetByLongUrl_ValidUrl_Hit(t *testing.T) {
	mockStore := new(MockStore)
	service := NewShortService(mockStore, nil)
	expectedShort := []ShortModel{
		{
			ID:          1,
			UserID:      1,
			OriginalUrl: "https://example.com",
			ShortUrl:    "12345678",
		},
	}
	mockStore.On("Search", mock.Anything, mock.Anything).Return(expectedShort, nil)

	result, err := service.GetByLongUrl(nil, "https://example.com")
	assert.NoError(t, err)
	assert.Equal(t, expectedShort[0], result)
	mockStore.AssertExpectations(t)
}

func TestGetByLongUrl_ValidUrl_Miss(t *testing.T) {
	mockStore := new(MockStore)
	service := NewShortService(mockStore, nil)

	mockStore.On("Search", mock.Anything, mock.Anything).Return([]ShortModel{}, ErrShortNotFound)

	_, err := service.GetByLongUrl(nil, "https://example.com")

	assert.Error(t, err)

	mockStore.AssertExpectations(t)
}

// Search Tests
func TestSearch_ByOriginalUrl_Found(t *testing.T) {
	mockStore := new(MockStore)
	service := NewShortService(mockStore, nil)

	expectedShort := []ShortModel{
		{
			OriginalUrl: "https://example.com",
		},
	}

	req := SearchRequest{
		OriginalUrl: &expectedShort[0].OriginalUrl,
	}

	mockStore.On("Search", req).Return(expectedShort, nil)

	result, err := service.Search(nil, req)
	assert.NoError(t, err)
	assert.Equal(t, expectedShort, result)
	mockStore.AssertExpectations(t)
}

func TestSearch_ByOriginalUrl_NotFound(t *testing.T) {
	mockStore := new(MockStore)
	service := NewShortService(mockStore, nil)

	url := "https://nonexistent.com"
	req := SearchRequest{OriginalUrl: &url}

	mockStore.On("Search", req).Return([]ShortModel{}, ErrShortNotFound)

	_, err := service.Search(nil, req)
	assert.Error(t, err)
	mockStore.AssertExpectations(t)
}

func TestSearch_NoCriteria(t *testing.T) {

	mockStore := new(MockStore)
	service := NewShortService(mockStore, nil)

	mockStore.On("Search", mock.Anything, mock.Anything).Return([]ShortModel{}, ErrShortNotFound)

	req := SearchRequest{}

	_, err := service.Search(nil, req)

	assert.Error(t, err)
}

// GetAllByUser Tests
func TestGetAllByUser_Found(t *testing.T) {
	mockStore := new(MockStore)
	service := NewShortService(mockStore, nil)

	userID := types.UserId(1)
	expectedShorts := []ShortModel{{ID: 1}, {ID: 2}}

	req := SearchRequest{
		UserId: &userID,
	}

	mockStore.On("Search", req).Return(expectedShorts, nil)

	result, err := service.GetAllByUser(nil, userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedShorts, result)
	mockStore.AssertExpectations(t)
}

func TestGetAllByUser_StoreError(t *testing.T) {
	mockStore := new(MockStore)
	service := NewShortService(mockStore, nil)
	userID := types.UserId(1)

	req := SearchRequest{
		UserId: &userID,
	}

	mockStore.On("Search", req).Return([]ShortModel{}, ErrShortNotFound)

	_, err := service.GetAllByUser(nil, userID)
	assert.Error(t, err)
	mockStore.AssertExpectations(t)
}

// GetAllURLs Tests
func TestGetAllURLs_Found(t *testing.T) {
	mockStore := new(MockStore)
	service := NewShortService(mockStore, nil)
	expectedShorts := []ShortModel{{ID: 1}, {ID: 2}}

	mockStore.On("GetAll", mock.Anything, mock.Anything, mock.Anything).Return(expectedShorts, 1, nil)
	result, _, err := service.GetAll(nil, 1, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedShorts, result)
	mockStore.AssertExpectations(t)
}

func TestGetAllURLs_StoreError(t *testing.T) {
	mockStore := new(MockStore)
	service := NewShortService(mockStore, nil)
	mockStore.On("GetAll", mock.Anything, mock.Anything, mock.Anything).Return([]ShortModel(nil), -1, errors.New("database error"))

	_, _, err := service.GetAll(nil, 1, 1)
	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	mockStore.AssertExpectations(t)
}

// DeleteURL Tests
func TestDeleteURL_Success(t *testing.T) {
	mockStore := new(MockStore)
	mockRedis := new(MockRedis)
	service := NewShortService(mockStore, mockRedis)

	shortID := types.ShortId(1)

	mockStore.On("Delete", shortID).Return(nil)
	mockStore.On("GetById", mock.Anything).Return(ShortModel{}, nil)

	mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	mockRedis.On("Delete", mock.Anything, mock.Anything).Return(nil)

	err := service.DeleteURL(context.Background(), shortID)
	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
}

func TestDeleteURL_StoreError(t *testing.T) {
	mockStore := new(MockStore)
	mockRedis := new(MockRedis)
	service := NewShortService(mockStore, mockRedis)

	shortID := types.ShortId(1)

	mockStore.On("Delete", shortID).Return(errors.New("database delete error"))
	mockStore.On("GetById", mock.Anything).Return(ShortModel{}, nil)

	mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	mockRedis.On("Delete", mock.Anything, mock.Anything).Return(nil)

	err := service.DeleteURL(nil, shortID)
	assert.Error(t, err)
	assert.EqualError(t, err, "database delete error")
	mockStore.AssertExpectations(t)
}
