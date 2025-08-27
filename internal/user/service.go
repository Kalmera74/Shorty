package user

import (
	"errors"
	"fmt"

	"github.com/Kalmera74/Shorty/validation"
)

type UserService struct {
	UserStore UserStore
}

func NewUserService(s UserStore) *UserService {
	return &UserService{s}
}

func (s *UserService) GetAllUsers() ([]UserResponse, error) {
	userModels, err := s.UserStore.GetAll()
	if err != nil {
		return nil, &UserError{Msg: fmt.Sprintf("Could not retrieve users. Reason: %v", err.Error()), Err: err}
	}

	users := make([]UserResponse, 0, len(userModels))
	for _, user := range userModels {
		users = append(users, UserResponse{
			Id:       user.ID,
			UserName: user.UserName,
			Email:    user.Email,
		})
	}

	return users, nil
}

func (s *UserService) GetUser(id uint) (UserResponse, error) {

	if err := validation.IsValidID(id); err != nil {
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
		Id:       userModel.ID,
		UserName: userModel.UserName,
		Email:    userModel.Email,
	}, nil
}

func (s *UserService) CreateUser(req UserCreateRequest) (UserResponse, error) {

	if err := req.Validate(); err != nil {
		return UserResponse{}, err
	}

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
		Id:       createdUser.ID,
		UserName: createdUser.UserName,
		Email:    createdUser.Email,
	}, nil
}

func (s *UserService) UpdateUser(id uint, req UserUpdateRequest) error {

	if err := req.Validate(); err != nil {
		return err
	}

	userModel, err := s.UserStore.Get(id)
	if err != nil {
		return &UserError{
			Msg: fmt.Sprintf("Could not retrieve user %d. Reason: %v", id, err.Error()),
			Err: err,
		}
	}

	if req.UserName != "" {
		userModel.UserName = req.UserName
	}
	if req.Email != "" {
		userModel.Email = req.Email
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

func (s *UserService) DeleteUser(id uint) error {

	if err := validation.IsValidID(id); err != nil {
		return err
	}
	err := s.UserStore.Delete(id)
	if err != nil {
		return &UserError{Msg: fmt.Sprintf("Could not delete user %d. Reason: %v", id, err.Error()), Err: err}
	}
	return nil
}
