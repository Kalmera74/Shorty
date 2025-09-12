package user

import (
	"context"
	"errors"
	"testing"

	"github.com/Kalmera74/Shorty/internal/types"
	"github.com/Kalmera74/Shorty/pkg/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock Repository ---
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAll(ctx context.Context, offset, limit int) ([]UserModel, int, error) {
	args := m.Called()
	return args.Get(0).([]UserModel), args.Int(1), args.Error(2)
}

func (m *MockUserRepository) Get(ctx context.Context, id types.UserId) (UserModel, error) {
	args := m.Called(id)
	return args.Get(0).(UserModel), args.Error(1)
}

func (m *MockUserRepository) Add(ctx context.Context, u UserModel) (UserModel, error) {
	args := m.Called(u)
	return args.Get(0).(UserModel), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, id types.UserId, u UserModel) error {
	args := m.Called(id, u)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id types.UserId) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (UserModel, error) {
	args := m.Called(email)
	return args.Get(0).(UserModel), args.Error(1)
}

// --- Tests ---
func TestGetAllUsers_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockRepo.On("GetAll").Return([]UserModel{
		{ID: 1, UserName: "alice", Email: "alice@test.com"},
	}, 1, nil)

	svc := NewUserService(mockRepo)

	users, _, err := svc.GetAllUsers(nil, 1, 1)

	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "alice", users[0].UserName)
	mockRepo.AssertExpectations(t)
}

func TestGetUser_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockRepo.On("Get", types.UserId(99)).Return(UserModel{}, ErrUserNotFound)

	svc := NewUserService(mockRepo)

	_, err := svc.GetUser(context.Background(), types.UserId(99))

	assert.Error(t, err)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	req := UserRegisterRequest{
		UserName: "bob",
		Email:    "bob@test.com",
		Password: "secure123",
	}

	mockRepo.On("GetByEmail", "bob@test.com").
		Return(UserModel{}, errors.New("not found")) // email not taken
	mockRepo.On("Add", mock.AnythingOfType("UserModel")).
		Return(UserModel{ID: 1, UserName: "bob", Email: "bob@test.com"}, nil)

	svc := NewUserService(mockRepo)

	resp, err := svc.CreateUser(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, "bob", resp.UserName)
	assert.Equal(t, "bob@test.com", resp.Email)
	mockRepo.AssertExpectations(t)
}

func TestVerifyCredentials_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	hashed, _ := security.HashPassword("mypassword")

	mockRepo.On("GetByEmail", "eve@test.com").
		Return(UserModel{ID: 2, Email: "eve@test.com", PasswordHash: hashed}, nil)

	svc := NewUserService(mockRepo)

	usr, err := svc.VerifyCredentials(context.Background(), "eve@test.com", "mypassword")

	assert.NoError(t, err)
	assert.NotNil(t, usr)
	assert.Equal(t, "eve@test.com", usr.Email)
}

func TestVerifyCredentials_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	hashed, _ := security.HashPassword("rightpassword")

	mockRepo.On("GetByEmail", "john@test.com").
		Return(UserModel{ID: 3, Email: "john@test.com", PasswordHash: hashed}, nil)

	svc := NewUserService(mockRepo)

	usr, err := svc.VerifyCredentials(context.Background(), "john@test.com", "wrongpassword")

	assert.Error(t, err)
	assert.Nil(t, usr)
	assert.EqualError(t, err, "invalid credentials")
}
