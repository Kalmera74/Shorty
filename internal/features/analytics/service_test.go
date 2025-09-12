package analytics

import (
	"context"
	"errors"
	"testing"

	"github.com/Kalmera74/Shorty/internal/features/shortener"
	"github.com/Kalmera74/Shorty/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAnalyticsRepository struct {
	mock.Mock
}

func (m *mockAnalyticsRepository) Create(ctx context.Context, click ClickModel) (ClickModel, error) {
	args := m.Called(ctx, click)
	return args.Get(0).(ClickModel), args.Error(1)
}

func (m *mockAnalyticsRepository) GetAll(ctx context.Context, offset, limit int) ([]ClickModel, int, error) {
	args := m.Called(ctx, offset, limit)
	var result []ClickModel
	if args.Get(0) != nil {
		result = args.Get(0).([]ClickModel)
	}
	return result, args.Int(1), args.Error(2)
}

func (m *mockAnalyticsRepository) GetAllByShortUrl(ctx context.Context, shortUrl string, offset, limit int) ([]ClickModel, int, error) {
	args := m.Called(ctx, shortUrl, offset, limit)
	var result []ClickModel
	if args.Get(0) != nil {
		result = args.Get(0).([]ClickModel)
	}
	return result, args.Int(1), args.Error(2)
}

func (m *mockAnalyticsRepository) GetByID(ctx context.Context, id types.ClickId) (ClickModel, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(ClickModel), args.Error(1)
}

// --- Create Tests ---
func TestCreate_Success(t *testing.T) {
	mockRepo := new(mockAnalyticsRepository)
	service := NewAnalyticService(mockRepo)

	clickToCreate := ClickModel{ShortID: 1}
	expectedClick := ClickModel{ID: 1, ShortID: types.ShortId(1)}
	mockRepo.On("Create", mock.Anything, clickToCreate).Return(expectedClick, nil).Once()

	result, err := service.Create(nil, clickToCreate)

	assert.NoError(t, err)
	assert.Equal(t, expectedClick, result)
	mockRepo.AssertExpectations(t)
}

func TestCreate_Failure(t *testing.T) {
	mockRepo := new(mockAnalyticsRepository)
	service := NewAnalyticService(mockRepo)

	clickToCreate := ClickModel{ShortID: types.ShortId(1)}
	repoError := errors.New("database error")
	mockRepo.On("Create", mock.Anything, clickToCreate).Return(ClickModel{}, repoError).Once()

	_, err := service.Create(nil, clickToCreate)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrClickCreateFail)
	mockRepo.AssertExpectations(t)
}

// --- GetAll Tests ---
func TestGetAll_Success(t *testing.T) {
	mockRepo := new(mockAnalyticsRepository)
	service := NewAnalyticService(mockRepo)

	expectedClicks := []ClickModel{{ID: 1}, {ID: 2}}

	mockRepo.On("GetAll", mock.Anything, 1, 1).Return(expectedClicks, 1, nil).Once()

	result, _, err := service.GetAll(nil, 1, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedClicks, result)
	mockRepo.AssertExpectations(t)
}

func TestGetAll_Failure_RepositoryError(t *testing.T) {
	mockRepo := new(mockAnalyticsRepository)
	service := NewAnalyticService(mockRepo)

	repoError := errors.New("repository error")
	mockRepo.On("GetAll", mock.Anything, 1, 1).Return(nil, -1, repoError)

	_, _, err := service.GetAll(nil, 1, 1)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetAll_Failure_NotFound(t *testing.T) {
	mockRepo := new(mockAnalyticsRepository)
	service := NewAnalyticService(mockRepo)

	mockRepo.On("GetAll", mock.Anything, 1, 1).Return([]ClickModel{}, 1, nil).Once()

	_, _, err := service.GetAll(nil, 1, 1)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

// --- GetAllByShortUrl Tests ---
func TestGetAllByShortUrl_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockAnalyticsRepository)
	service := NewAnalyticService(mockRepo)

	shortUrl := "testurl"
	expectedClicks := []ClickModel{
		{ID: 1, Short: shortener.ShortModel{
			ShortUrl: shortUrl,
		}}, {ID: 2, Short: shortener.ShortModel{
			ShortUrl: shortUrl,
		}}}

	mockRepo.On("GetAllByShortUrl", ctx, shortUrl, mock.Anything, mock.Anything).Return(expectedClicks, 1, nil).Once()

	result, _, err := service.GetAllByShortUrl(ctx, shortUrl, 1, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedClicks, result)
	mockRepo.AssertExpectations(t)
}

func TestGetAllByShortUrl_Failure_RepositoryError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockAnalyticsRepository)
	service := NewAnalyticService(mockRepo)
	shortUrl := "testurl"
	repoError := errors.New("repository error")
	mockRepo.On("GetAllByShortUrl", ctx, shortUrl, 1, 1).Return(nil, 1, repoError).Once()

	_, _, err := service.GetAllByShortUrl(ctx, shortUrl, 1, 1)

	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	mockRepo.AssertExpectations(t)
}

func TestGetAllByShortUrl_Failure_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockAnalyticsRepository)
	service := NewAnalyticService(mockRepo)
	shortUrl := "testurl"
	mockRepo.On("GetAllByShortUrl", ctx, shortUrl, 1, 1).Return([]ClickModel{}, 1, nil).Once()

	_, _, err := service.GetAllByShortUrl(ctx, shortUrl, 1, 1)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrClickNotFound)
	mockRepo.AssertExpectations(t)
}

// --- GetAllClicks Tests ---
func TestGetAllClicks_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockAnalyticsRepository)
	service := NewAnalyticService(mockRepo)
	expectedClicks := []ClickModel{{ID: 1}, {ID: 2}}
	mockRepo.On("GetAll", ctx, 1, 1).Return(expectedClicks, 1, nil).Once()

	result, _, err := service.GetAllClicks(ctx, 1, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedClicks, result)
	mockRepo.AssertExpectations(t)
}

func TestGetAllClicks_Failure_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockAnalyticsRepository)
	service := NewAnalyticService(mockRepo)
	mockRepo.On("GetAll", ctx, 1, 1).Return([]ClickModel{}, 1, nil).Once()

	_, _, err := service.GetAllClicks(ctx, 1, 1)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrClicksNotFound)
	mockRepo.AssertExpectations(t)
}

// --- GetByID Tests ---
func TestGetByID_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockAnalyticsRepository)
	service := NewAnalyticService(mockRepo)
	testID := types.ClickId(42)
	expectedClick := ClickModel{ID: testID}
	mockRepo.On("GetByID", ctx, testID).Return(expectedClick, nil).Once()

	result, err := service.GetByID(ctx, testID)

	assert.NoError(t, err)
	assert.Equal(t, expectedClick, result)
	mockRepo.AssertExpectations(t)
}

func TestGetByID_Failure_RepositoryError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockAnalyticsRepository)
	service := NewAnalyticService(mockRepo)
	testID := types.ClickId(42)
	repoError := errors.New("repository error")
	mockRepo.On("GetByID", ctx, testID).Return(ClickModel{}, repoError).Once()

	_, err := service.GetByID(ctx, testID)

	assert.Error(t, err)
	assert.Equal(t, repoError, err)
	mockRepo.AssertExpectations(t)
}

func TestGetByID_Failure_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockAnalyticsRepository)
	service := NewAnalyticService(mockRepo)
	testID := types.ClickId(42)
	mockRepo.On("GetByID", ctx, testID).Return(ClickModel{ID: 0}, nil).Once()

	_, err := service.GetByID(ctx, testID)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrClickNotFound)
	mockRepo.AssertExpectations(t)
}
