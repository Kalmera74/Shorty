package user

import (
	"time"

	"github.com/Kalmera74/Shorty/internal/features/shortener"
	"github.com/Kalmera74/Shorty/internal/types"
)

type UserModel struct {
	ID           types.UserId           `gorm:"primaryKey"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	DeletedAt    time.Time              `json:"deleted_at,omitempty" gorm:"index"`
	UserName     string                 `validate:"required,min=3,max=30"`
	Email        string                 `validate:"required,email"`
	PasswordHash string                 `validate:"required"`
	Role         string                 `validate:"required"`
	Shorts       []shortener.ShortModel `gorm:"foreignKey:UserID" validate:"dive"`
}
