package user

import (
	"time"

	"github.com/Kalmera74/Shorty/internal/features/shortener"
	"github.com/Kalmera74/Shorty/internal/types"
	"gorm.io/gorm"
)

type UserModel struct {
	ID           types.UserId `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt         `gorm:"index"`
	UserName     string                 `validate:"required,min=3,max=30"`
	Email        string                 `validate:"required,email"`
	PasswordHash string                 `validate:"required"`
	Shorts       []shortener.ShortModel `gorm:"foreignKey:UserID" validate:"dive"`
}
