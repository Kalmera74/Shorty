package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Kalmera74/Shorty/internal/types"
	"github.com/Kalmera74/Shorty/pkg/security"
)

type IUserService interface {
	GetAllUsers(ctx context.Context) ([]UserResponse, error)
	GetUser(ctx context.Context, id types.UserId) (UserResponse, error)
	CreateUser(ctx context.Context, req UserCreateRequest) (UserResponse, error)
	UpdateUser(ctx context.Context, id types.UserId, req UserUpdateRequest) error
	DeleteUser(ctx context.Context, id types.UserId) error
	VerifyCredentials(ctx context.Context, email, password string) (*UserModel, error)
	GetByEmail(ctx context.Context, email string) (*UserModel, error)
}
type userService struct {
	UserStore IUserRepository
}

func NewUserService(s IUserRepository) IUserService {
	return &userService{s}
}

func (s *userService) GetAllUsers(ctx context.Context) ([]UserResponse, error) {
	userModels, err := s.UserStore.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w, %v", ErrUserNotFound, err)
	}

	users := make([]UserResponse, 0, len(userModels))
	for _, user := range userModels {
		users = append(users, UserResponse{
			Id:       uint(user.ID),
			UserName: user.UserName,
			Email:    user.Email,
		})
	}

	return users, nil
}
func (s *userService) GetUser(ctx context.Context, id types.UserId) (UserResponse, error) {

	userModel, err := s.UserStore.Get(ctx,uint(id))
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return UserResponse{}, err
		}
		return UserResponse{}, fmt.Errorf("Could not retrieve user %d. Reason: %v", id, err.Error())
	}

	return UserResponse{
		Id:       uint(userModel.ID),
		UserName: userModel.UserName,
		Email:    userModel.Email,
	}, nil
}
func (s *userService) CreateUser(ctx context.Context, req UserCreateRequest) (UserResponse, error) {

	existingUser, err := s.GetByEmail(ctx,req.Email)

	if existingUser != nil {
		return UserResponse{}, errors.New("The email is already in use")
	}

	hasPss, err := security.HashPassword(req.Password)

	if err != nil {
		return UserResponse{}, err
	}

	newUser := UserModel{
		UserName:     req.UserName,
		Email:        req.Email,
		PasswordHash: hasPss,
	}

	createdUser, err := s.UserStore.Add(ctx,newUser)
	if err != nil {
		return UserResponse{}, fmt.Errorf("Could not create user. Reason: %v", err.Error())
	}

	return UserResponse{
		Id:       uint(createdUser.ID),
		UserName: createdUser.UserName,
		Email:    createdUser.Email,
	}, nil
}
func (s *userService) UpdateUser(ctx context.Context, id types.UserId, req UserUpdateRequest) error {

	userModel, err := s.UserStore.Get(ctx,uint(id))
	if err != nil {
		return fmt.Errorf("Could not retrieve user %d. Reason: %v", id, err.Error())
	}

	if req.UserName != nil {
		userModel.UserName = *req.UserName
	}
	if req.Email != nil {
		userModel.Email = *req.Email
	}

	err = s.UserStore.Update(ctx,uint(id), userModel)
	if err != nil {
		return fmt.Errorf("Could not update user %d. Reason: %v", id, err.Error())
	}

	return nil
}
func (s *userService) DeleteUser(ctx context.Context, id types.UserId) error {

	err := s.UserStore.Delete(ctx,uint(id))
	if err != nil {
		return fmt.Errorf("Could not delete user %d. Reason: %v", id, err.Error())
	}
	return nil
}
func (s *userService) GetByEmail(ctx context.Context, email string) (*UserModel, error) {

	user, err := s.UserStore.GetByEmail(ctx,email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (s *userService) VerifyCredentials(ctx context.Context, email, password string) (*UserModel, error) {
	user, err := s.GetByEmail(ctx,email)

	if err != nil {
		return nil, fmt.Errorf("%w, %v", ErrInvalidCredentials, err)
	}

	if security.CheckPassword(password, user.PasswordHash) {
		return user, nil
	}

	return nil, fmt.Errorf("%w", ErrInvalidCredentials)
}
