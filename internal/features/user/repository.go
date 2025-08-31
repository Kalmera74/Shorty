package user

import (
	"errors"

	"gorm.io/gorm"
)

type UserStore interface {
	GetAll() ([]UserModel, error)
	Get(id uint) (UserModel, error)
	Add(user UserModel) (UserModel, error)
	Update(id uint, user UserModel) error
	Delete(id uint) error
}

type PostgresUserStore struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *PostgresUserStore {
	return &PostgresUserStore{db}
}

func (s *PostgresUserStore) GetAll() ([]UserModel, error) {
	var users []UserModel
	result := s.db.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	if len(users) == 0 {
		return nil, &UserNotFoundError{Msg: "No user was found", Err: nil}
	}

	return users, nil
}
func (s *PostgresUserStore) Get(id uint) (UserModel, error) {
	var u UserModel
	result := s.db.First(&u, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return UserModel{}, &UserNotFoundError{Msg: "User not found", Err: result.Error}
	}
	if result.Error != nil {
		return UserModel{}, result.Error
	}

	return u, nil
}
func (s *PostgresUserStore) Add(u UserModel) (UserModel, error) {
	result := s.db.Create(&u)

	if result.Error != nil {
		return UserModel{}, result.Error
	}

	return u, nil
}
func (s *PostgresUserStore) Update(id uint, req UserModel) error {
	var existingUser UserModel

	result := s.db.First(&existingUser, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &UserNotFoundError{Msg: "User not found", Err: result.Error}
	}
	if result.Error != nil {
		return result.Error
	}

	updateResult := s.db.Model(&existingUser).Updates(req)

	if updateResult.Error != nil {
		return updateResult.Error
	}

	return nil
}
func (s *PostgresUserStore) Delete(id uint) error {
	result := s.db.Delete(&UserModel{}, id)

	if result.RowsAffected == 0 {
		return &UserNotFoundError{Msg: "User not found", Err: gorm.ErrRecordNotFound}
	}

	return result.Error
}
