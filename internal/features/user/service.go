package user

import (
	"errors"
	"fmt"

	"github.com/Kalmera74/Shorty/internal/types"
	"github.com/Kalmera74/Shorty/pkg/security"
)

type IUserService interface {
	GetAllUsers() ([]UserResponse, error)
	GetUser(id types.UserId) (UserResponse, error)
	CreateUser(req UserCreateRequest) (UserResponse, error)
	UpdateUser(id types.UserId, req UserUpdateRequest) error
	DeleteUser(id types.UserId) error
	VerifyCredentials(email, password string) (*UserModel, error)
	GetByEmail(email string) (*UserModel, error)
}
type userService struct {
	UserStore UserStore
}

func NewUserService(s UserStore) IUserService {
	return &userService{s}
}

func (s *userService) GetAllUsers() ([]UserResponse, error) {
	userModels, err := s.UserStore.GetAll()
	if err != nil {
		return nil, &UserError{Msg: fmt.Sprintf("Could not retrieve users. Reason: %v", err.Error()), Err: err}
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
func (s *userService) GetUser(id types.UserId) (UserResponse, error) {

	userModel, err := s.UserStore.Get(uint(id))
	if err != nil {
		if errors.Is(err, &UserNotFoundError{}) {
			return UserResponse{}, err
		}
		return UserResponse{}, &UserError{Msg: fmt.Sprintf("Could not retrieve user %d. Reason: %v", id, err.Error()), Err: err}
	}

	return UserResponse{
		Id:       uint(userModel.ID),
		UserName: userModel.UserName,
		Email:    userModel.Email,
	}, nil
}
func (s *userService) CreateUser(req UserCreateRequest) (UserResponse, error) {

	existingUser, err := s.GetByEmail(req.Email)


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

	createdUser, err := s.UserStore.Add(newUser)
	if err != nil {
		return UserResponse{}, &UserError{
			Msg: fmt.Sprintf("Could not create user. Reason: %v", err.Error()),
			Err: err,
		}
	}

	return UserResponse{
		Id:       uint(createdUser.ID),
		UserName: createdUser.UserName,
		Email:    createdUser.Email,
	}, nil
}
func (s *userService) UpdateUser(id types.UserId, req UserUpdateRequest) error {

	userModel, err := s.UserStore.Get(uint(id))
	if err != nil {
		return &UserError{
			Msg: fmt.Sprintf("Could not retrieve user %d. Reason: %v", id, err.Error()),
			Err: err,
		}
	}

	if req.UserName != nil {
		userModel.UserName = *req.UserName
	}
	if req.Email != nil {
		userModel.Email = *req.Email
	}

	err = s.UserStore.Update(uint(id), userModel)
	if err != nil {
		return &UserError{
			Msg: fmt.Sprintf("Could not update user %d. Reason: %v", id, err.Error()),
			Err: err,
		}
	}

	return nil
}
func (s *userService) DeleteUser(id types.UserId) error {

	err := s.UserStore.Delete(uint(id))
	if err != nil {
		return &UserError{Msg: fmt.Sprintf("Could not delete user %d. Reason: %v", id, err.Error()), Err: err}
	}
	return nil
}

func (s *userService) GetByEmail(email string) (*UserModel, error) {

	user, err := s.UserStore.GetByEmail(email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userService) VerifyCredentials(email, password string) (*UserModel, error) {
	user, err := s.GetByEmail(email)

	if err != nil {
		return nil, err
	}

	if security.CheckPassword(password, user.PasswordHash) {
		return user, nil
	}

	return nil, errors.New("invalid credentials")
}
