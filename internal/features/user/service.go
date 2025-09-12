package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Kalmera74/Shorty/internal/types"
	"github.com/Kalmera74/Shorty/pkg/security"
)

type IUserService interface {
	GetAllUsers(ctx context.Context, page, pageSize int) ([]UserModel, int, error)
	GetUser(ctx context.Context, id types.UserId) (UserModel, error)
	CreateUser(ctx context.Context, req UserRegisterRequest) (UserModel, error)
	UpdateUser(ctx context.Context, id types.UserId, req UserUpdateRequest) error
	DeleteUser(ctx context.Context, id types.UserId) error
	VerifyCredentials(ctx context.Context, email, password string) (*UserModel, error)
	GetByEmail(ctx context.Context, email string) (*UserModel, error)
}
type userService struct {
	Repository IUserRepository
}

func NewUserService(s IUserRepository) IUserService {
	return &userService{s}
}

func (s *userService) GetAllUsers(ctx context.Context, page, pageSize int) ([]UserModel, int, error) {
	offset := (page - 1) * pageSize
	users, total, err := s.Repository.GetAll(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("%w, %v", ErrUserNotFound, err)
	}

	return users, total, nil
}

func (s *userService) GetUser(ctx context.Context, id types.UserId) (UserModel, error) {

	userModel, err := s.Repository.Get(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return UserModel{}, err
		}
		return UserModel{}, fmt.Errorf("could not retrieve user %d: %w", id, err)
	}

	return userModel, nil
}

func (s *userService) CreateUser(ctx context.Context, req UserRegisterRequest) (UserModel, error) {

	existingUser, err := s.GetByEmail(ctx, req.Email)

	if existingUser != nil {
		return UserModel{}, errors.New("The email is already in use")
	}

	hasPss, err := security.HashPassword(req.Password)

	if err != nil {
		return UserModel{}, err
	}

	newUser := UserModel{
		UserName:     req.UserName,
		Email:        req.Email,
		PasswordHash: hasPss,
		Role:         "user",
	}

	createdUser, err := s.Repository.Add(ctx, newUser)
	if err != nil {
		return UserModel{}, fmt.Errorf("Could not create user. Reason: %v", err.Error())
	}

	return createdUser, nil
}
func (s *userService) UpdateUser(ctx context.Context, id types.UserId, req UserUpdateRequest) error {

	userModel, err := s.Repository.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("Could not retrieve user %d. Reason: %v", id, err.Error())
	}

	if req.UserName != nil {
		userModel.UserName = *req.UserName
	}
	if req.Email != nil {
		userModel.Email = *req.Email
	}

	err = s.Repository.Update(ctx, id, userModel)
	if err != nil {
		return fmt.Errorf("Could not update user %d. Reason: %v", id, err.Error())
	}

	return nil
}
func (s *userService) DeleteUser(ctx context.Context, id types.UserId) error {

	err := s.Repository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("Could not delete user %d. Reason: %v", id, err.Error())
	}
	return nil
}
func (s *userService) GetByEmail(ctx context.Context, email string) (*UserModel, error) {

	user, err := s.Repository.GetByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (s *userService) VerifyCredentials(ctx context.Context, email, password string) (*UserModel, error) {
	user, err := s.GetByEmail(ctx, email)

	if err != nil {
		return nil, fmt.Errorf("%w, %v", ErrInvalidCredentials, err)
	}

	if security.CheckPassword(password, user.PasswordHash) {
		return user, nil
	}

	return nil, fmt.Errorf("%w", ErrInvalidCredentials)
}
