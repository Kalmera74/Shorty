package user

import (
	"errors"
	"fmt"

	"github.com/Kalmera74/Shorty/internal/validation"
)

type IUserService interface {
	GetAllUsers() ([]UserResponse, error)
	GetUser(id uint) (UserResponse, error)
	CreateUser(req UserCreateRequest) (UserResponse, error)
	UpdateUser(id uint, req UserUpdateRequest) error
	DeleteUser(id uint) error
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
func (s *userService) GetUser(id uint) (UserResponse, error) {

	if err := validation.ValidateID(id); err != nil {
		return UserResponse{}, err
	}

	userModel, err := s.UserStore.Get(id)
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

	newUser := UserModel{
		UserName: req.UserName,
		Email:    req.Email,
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
func (s *userService) UpdateUser(id uint, req UserUpdateRequest) error {

	userModel, err := s.UserStore.Get(id)
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

	err = s.UserStore.Update(id, userModel)
	if err != nil {
		return &UserError{
			Msg: fmt.Sprintf("Could not update user %d. Reason: %v", id, err.Error()),
			Err: err,
		}
	}

	return nil
}
func (s *userService) DeleteUser(id uint) error {

	if err := validation.ValidateID(id); err != nil {
		return err
	}
	err := s.UserStore.Delete(id)
	if err != nil {
		return &UserError{Msg: fmt.Sprintf("Could not delete user %d. Reason: %v", id, err.Error()), Err: err}
	}
	return nil
}
