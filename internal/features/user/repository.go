package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/Kalmera74/Shorty/internal/types"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetAll(ctx context.Context, offset, limit int) ([]UserModel, int, error)
	Get(ctx context.Context, id types.UserId) (UserModel, error)
	Add(ctx context.Context, user UserModel) (UserModel, error)
	Update(ctx context.Context, id types.UserId, user UserModel) error
	Delete(ctx context.Context, id types.UserId) error
	GetByEmail(ctx context.Context, email string) (UserModel, error)
}

type postgresUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &postgresUserRepository{db}
}

func (s *postgresUserRepository) GetAll(ctx context.Context, offset, limit int) ([]UserModel, int, error) {
	var users []UserModel
	var total int64

	if err := s.db.WithContext(ctx).Model(&UserModel{}).Count(&total).Error; err != nil {
		return nil, -1, fmt.Errorf("%w: %v", ErrUserNotFound, err)
	}

	if err := s.db.WithContext(ctx).
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, -1, fmt.Errorf("%w: %v", ErrUserNotFound, err)
	}

	if len(users) == 0 {
		return nil, -1, fmt.Errorf("%w", ErrUserNotFound)
	}

	return users, int(total), nil
}

func (s *postgresUserRepository) Get(ctx context.Context, id types.UserId) (UserModel, error) {
	var u UserModel
	result := s.db.WithContext(ctx).Preload("Shorts").First(&u, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return UserModel{}, fmt.Errorf("%w, %v", ErrUserNotFound, id)
	}
	if result.Error != nil {
		return UserModel{}, fmt.Errorf("could not retrieve user %d. Reason: %v", id, result.Error)
	}

	return u, nil
}
func (s *postgresUserRepository) GetByEmail(ctx context.Context, email string) (UserModel, error) {

	var user UserModel
	result := s.db.WithContext(ctx).Where("email =?", email).First(&user)

	if result.Error != nil {
		return UserModel{}, fmt.Errorf("%w: %v", ErrUserNotFound, result.Error)
	}

	return user, nil

}
func (s *postgresUserRepository) Add(ctx context.Context, u UserModel) (UserModel, error) {
	result := s.db.WithContext(ctx).Create(&u)

	if result.Error != nil {
		return UserModel{}, fmt.Errorf("could not create user. Reason: %v", result.Error)
	}

	return u, nil
}
func (s *postgresUserRepository) Update(ctx context.Context, id types.UserId, req UserModel) error {
	var existingUser UserModel

	result := s.db.WithContext(ctx).First(&existingUser, id)

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
func (s *postgresUserRepository) Delete(ctx context.Context, id types.UserId) error {
	result := s.db.WithContext(ctx).Delete(&UserModel{}, id)

	if result.RowsAffected == 0 {
		return fmt.Errorf("%w, %v", ErrUserNotFound, id)
	}

	return nil
}
