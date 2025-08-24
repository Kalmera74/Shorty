package stores

import (
	"errors"

	"github.com/Kalmera74/Shorty/internal/user"
	"gorm.io/gorm"
)

type PostgresUserStore struct {
	db *gorm.DB
}

func NewPostgresUserStore(db *gorm.DB) *PostgresUserStore {
	return &PostgresUserStore{db}
}

func (s *PostgresUserStore) GetAll() ([]user.UserModel, error) {
	var users []user.UserModel
	result := s.db.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	if len(users) == 0 {
		return nil, &user.UserError{Msg: "No user was found", Err: nil}
	}

	return users, nil
}

func (s *PostgresUserStore) Get(id uint) (user.UserModel, error) {
	var u user.UserModel
	result := s.db.First(&u, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return user.UserModel{}, &user.UserError{Msg: "User not found", Err: result.Error}
	}
	if result.Error != nil {
		return user.UserModel{}, result.Error
	}

	return u, nil
}

func (s *PostgresUserStore) Update(id uint, req user.UserModel) error {
	var existingUser user.UserModel

	result := s.db.First(&existingUser, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &user.UserError{Msg: "User not found", Err: result.Error}
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
	result := s.db.Delete(&user.UserModel{}, id)

	if result.RowsAffected == 0 {
		return &user.UserError{Msg: "User not found", Err: gorm.ErrRecordNotFound}
	}

	return result.Error
}
