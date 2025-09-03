package user

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type IUserRepository interface {
	GetAll() ([]UserModel, error)
	Get(id uint) (UserModel, error)
	Add(user UserModel) (UserModel, error)
	Update(id uint, user UserModel) error
	Delete(id uint) error
	GetByEmail(email string) (UserModel, error)
}

type postgresUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &postgresUserRepository{db}
}

func (s *postgresUserRepository) GetAll() ([]UserModel, error) {
	var users []UserModel
	result := s.db.Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("could not retrieve users. Reason: %v", result.Error)
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("%w, %v", ErrUserNotFound, "no users in database")
	}

	return users, nil
}
func (s *postgresUserRepository) Get(id uint) (UserModel, error) {
	var u UserModel
	result := s.db.First(&u, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return UserModel{}, fmt.Errorf("%w, %v", ErrUserNotFound, id)
	}
	if result.Error != nil {
		return UserModel{}, fmt.Errorf("could not retrieve user %d. Reason: %v", id, result.Error)
	}

	return u, nil
}
func (s *postgresUserRepository) GetByEmail(email string) (UserModel, error) {

	var user UserModel
	result := s.db.Where("email =?", email).First(&user)

	if result.Error != nil {
		return UserModel{}, result.Error
	}

	return user, nil

}
func (s *postgresUserRepository) Add(u UserModel) (UserModel, error) {
	result := s.db.Create(&u)

	if result.Error != nil {
		return UserModel{}, fmt.Errorf("could not create user. Reason: %v", result.Error)
	}

	return u, nil
}
func (s *postgresUserRepository) Update(id uint, req UserModel) error {
	var existingUser UserModel

	result := s.db.First(&existingUser, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("%w, %v", ErrUserNotFound, id)
	}
	if result.Error != nil {
		return fmt.Errorf("could not retrieve user %d. Reason: %v", id, result.Error)
	}

	updateResult := s.db.Model(&existingUser).Updates(req)

	if updateResult.Error != nil {
		return fmt.Errorf("could not update user %d. Reason: %v", id, updateResult.Error)
	}

	return nil
}
func (s *postgresUserRepository) Delete(id uint) error {
	result := s.db.Delete(&UserModel{}, id)

	if result.RowsAffected == 0 {
		return fmt.Errorf("%w, %v", ErrUserNotFound, id)
	}

	return nil
}
