package shortener

import (
	"context"
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
	cmd := redis.NewStringResult(val, err)
	return val, cmd.Err()
}

func (m *MockRedis) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	m.Called(ctx, key, value, ttl)
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

	expectedShort.ID = 1

	mockRedis.On("Get", mock.Anything, req.Url).Return("", redis.Nil)
	mockStore.On("GetByLongUrl", mock.Anything).Return(ShortModel{}, &ShortNotFoundError{})
	mockStore.On("Create", mock.Anything).Return(expectedShort, nil)
	mockRedis.On("Set", mock.Anything, req.Url, expectedShort.ID, time.Minute*5).Return(nil)
	mockRedis.On("Set", mock.Anything, expectedShort.ShortUrl, mock.Anything, time.Minute*5).Return(nil)

	service := NewShortService(mockStore, mockRedis)
	result, err := service.ShortenURL(req)

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

	mockStore.On("GetById", expectedShort.ID).Return(expectedShort, nil)

	service := NewShortService(mockStore, mockRedis)
	result, err := service.ShortenURL(req)

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
	result, err := service.ShortenURL(req)

	assert.NoError(t, err)
	assert.Equal(t, expectedShort.ID, result.ID)

	mockStore.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
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

	mockStore.On("GetById", mock.Anything).Return(expectedShort, nil)
	mockRedis.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	service := NewShortService(mockStore, mockRedis)
	result, err := service.GetById(expectedShort.ID)

	assert.Equal(t, expectedShort, result)
	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
	mockRedis.AssertExpectations(t)

}
