package shortener

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
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

func (m *MockStore) Create(url ShortModel) (ShortModel, error) {
	args := m.Called(url)
	return args.Get(0).(ShortModel), args.Error(1)
}
func (m *MockStore) GetById(id types.ShortId) (ShortModel, error) {
	args := m.Called(id)
	return args.Get(0).(ShortModel), args.Error(1)
}
func (m *MockStore) GetByShortUrl(shortUrl string) (ShortModel, error) {
	args := m.Called(shortUrl)
	return args.Get(0).(ShortModel), args.Error(1)
}
func (m *MockStore) GetByLongUrl(longUrl string) (ShortModel, error) {
	args := m.Called(longUrl)
	return args.Get(0).(ShortModel), args.Error(1)
}
func (m *MockStore) GetAllByUser(userId types.UserId) ([]ShortModel, error) {
	args := m.Called(userId)
	return args.Get(0).([]ShortModel), args.Error(1)
}
func (m *MockStore) GetAll() ([]ShortModel, error) {
	args := m.Called()
	return args.Get(0).([]ShortModel), args.Error(1)
}
func (m *MockStore) Delete(id types.ShortId) error {
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
func TestShortenURL_Miss_Cache_Miss_DB(t *testing.T) {
	mockStore := new(MockStore)
	mockRedis := new(MockRedis)

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

	mockRedis.On("Get", mock.Anything, req.Url).Return("", redis.Nil)
	mockStore.On("GetByLongUrl", mock.Anything).Return(ShortModel{}, ErrShortNotFound)
	mockStore.On("Create", mock.Anything).Return(expectedShort, nil)
	mockRedis.On("Set", mock.Anything, req.Url, expectedShort.ID, time.Minute*5).Return(nil)
	mockRedis.On("Set", mock.Anything, expectedShort.ShortUrl, mock.Anything, time.Minute*5).Return(nil)

	service := NewShortService(mockStore, mockRedis)
	result, err := service.ShortenURL(nil, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedShort.ID, result.ID)

	mockStore.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
}

func TestShortenURL_Hit_Cache_Hit_DB(t *testing.T) {
	mockStore := new(MockStore)
	mockRedis := new(MockRedis)

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

	mockRedis.On("Get", mock.Anything, req.Url).Return(strconv.Itoa(int(expectedShort.ID)), nil)
	mockStore.On("GetById", types.ShortId(expectedShort.ID)).Return(expectedShort, nil)

	service := NewShortService(mockStore, mockRedis)
	result, err := service.ShortenURL(nil, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedShort.ID, result.ID)

	mockStore.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
}

func TestShortURL_Miss_Cache_Hit_DB(t *testing.T) {
	mockStore := new(MockStore)
	mockRedis := new(MockRedis)

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

	mockRedis.On("Get", mock.Anything, req.Url).Return("", redis.Nil)
	mockStore.On("GetByLongUrl", req.Url).Return(expectedShort, nil)
	mockRedis.On("Set", mock.Anything, req.Url, expectedShort.ID, time.Minute*5).Return(nil)
	mockRedis.On("Set", mock.Anything, expectedShort.ShortUrl, expectedShort.OriginalUrl, time.Minute*5).Return(nil)

	service := NewShortService(mockStore, mockRedis)
	result, err := service.ShortenURL(nil, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedShort.ID, result.ID)

	mockStore.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
}

func TestShortenURL_New_StoreCreateFails(t *testing.T) {
	mockStore := new(MockStore)
	mockRedis := new(MockRedis)

	req := ShortenRequest{
		UserID: 1,
		Url:    "https://example.com",
	}

	mockRedis.On("Get", mock.Anything, req.Url).Return("", redis.Nil)
	mockStore.On("GetByLongUrl", mock.Anything).Return(ShortModel{}, ErrShortNotFound)
	mockStore.On("Create", mock.Anything).Return(ShortModel{}, errors.New("database create error"))

	service := NewShortService(mockStore, mockRedis)
	_, err := service.ShortenURL(nil, req)

	assert.Error(t, err)

	mockStore.AssertExpectations(t)
	mockRedis.AssertNotCalled(t, "Set")
}

// GetById Tests
func TestGetById_ValidId(t *testing.T) {
	mockStore := new(MockStore)
	mockRedis := new(MockRedis)

	expectedShort := ShortModel{
		ID:          1,
		UserID:      1,
		OriginalUrl: "https://example.com",
		ShortUrl:    "12345678",
	}

	mockStore.On("GetById", types.ShortId(1)).Return(expectedShort, nil)
	mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	service := NewShortService(mockStore, mockRedis)
	result, err := service.GetById(nil, types.ShortId(1))

	assert.Equal(t, expectedShort, result)
	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
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

	expectedShort := ShortModel{
		ID:          1,
		UserID:      1,
		OriginalUrl: "https://example.com",
		ShortUrl:    "12345678",
	}

	mockRedis.On("Get", mock.Anything, expectedShort.ShortUrl).Return("", redis.Nil)
	mockStore.On("GetByShortUrl", expectedShort.ShortUrl).Return(expectedShort, nil)
	mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	service := NewShortService(mockStore, mockRedis)
	result, err := service.GetByShortUrl(nil, expectedShort.ShortUrl)

	assert.NoError(t, err)
	assert.Equal(t, expectedShort, result)
	mockStore.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
}

func TestGetByShortUrl_CacheMiss_DBMiss(t *testing.T) {
	mockStore := new(MockStore)
	mockRedis := new(MockRedis)

	mockRedis.On("Get", mock.Anything, mock.Anything).Return("", redis.Nil)
	mockStore.On("GetByShortUrl", mock.Anything).Return(ShortModel{}, ErrShortNotFound)

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
	expectedShort := ShortModel{OriginalUrl: "https://example.com"}

	mockStore.On("GetByLongUrl", mock.Anything).Return(expectedShort, nil)
	result, err := service.GetByLongUrl(nil, "https://example.com")
	assert.NoError(t, err)
	assert.Equal(t, expectedShort, result)
	mockStore.AssertExpectations(t)
}

func TestGetByLongUrl_ValidUrl_Miss(t *testing.T) {
	mockStore := new(MockStore)
	service := NewShortService(mockStore, nil)

	mockStore.On("GetByLongUrl", mock.Anything).Return(ShortModel{}, ErrShortNotFound)

	_, err := service.GetByLongUrl(nil, "https://example.com")

	assert.Error(t, err)

	mockStore.AssertExpectations(t)
}

// Search Tests
func TestSearch_ByOriginalUrl_Found(t *testing.T) {
	mockStore := new(MockStore)
	service := NewShortService(mockStore, nil)
	expectedShort := ShortModel{OriginalUrl: "https://example.com"}
	req := SearchRequest{OriginalUrl: &expectedShort.OriginalUrl}

	mockStore.On("GetByLongUrl", *req.OriginalUrl).Return(expectedShort, nil)
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

	mockStore.On("GetByLongUrl", url).Return(ShortModel{}, ErrShortNotFound)
	_, err := service.Search(nil, req)
	assert.Error(t, err)
	mockStore.AssertExpectations(t)
}

func TestSearch_NoCriteria(t *testing.T) {
	service := NewShortService(nil, nil)
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

	mockStore.On("GetAllByUser", userID).Return(expectedShorts, nil)
	result, err := service.GetAllByUser(nil, userID)
	assert.NoError(t, err)
	assert.Equal(t, expectedShorts, result)
	mockStore.AssertExpectations(t)
}

func TestGetAllByUser_StoreError(t *testing.T) {
	mockStore := new(MockStore)
	service := NewShortService(mockStore, nil)
	userID := types.UserId(1)
	mockStore.On("GetAllByUser", userID).Return([]ShortModel(nil), errors.New("database error"))

	_, err := service.GetAllByUser(nil, userID)
	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
	mockStore.AssertExpectations(t)
}

// GetAllURLs Tests
func TestGetAllURLs_Found(t *testing.T) {
	mockStore := new(MockStore)
	service := NewShortService(mockStore, nil)
	expectedShorts := []ShortModel{{ID: 1}, {ID: 2}}

	mockStore.On("GetAll").Return(expectedShorts, nil)
	result, err := service.GetAllURLs(nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedShorts, result)
	mockStore.AssertExpectations(t)
}

func TestGetAllURLs_StoreError(t *testing.T) {
	mockStore := new(MockStore)
	service := NewShortService(mockStore, nil)
	mockStore.On("GetAll").Return([]ShortModel(nil), errors.New("database error"))

	_, err := service.GetAllURLs(nil)
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
