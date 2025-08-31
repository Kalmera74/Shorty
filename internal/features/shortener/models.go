package shortener

import (
	"time"

	"github.com/Kalmera74/Shorty/internal/types"
	"gorm.io/gorm"
)

type ShortModel struct {
	ID          types.ShortId `gorm:"primaryKey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	UserID      types.UserId    `gorm:"not null" validate:"required,numeric,min=1"`
	OriginalUrl string         `gorm:"not null" validate:"required,url"`
	ShortUrl    string         `gorm:"unique;not null" validate:"required"`
}
