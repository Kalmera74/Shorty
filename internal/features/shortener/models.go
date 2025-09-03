package shortener

import (
	"time"

	"github.com/Kalmera74/Shorty/internal/types"
)

type ShortModel struct {
	ID          types.ShortId `gorm:"primaryKey"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	DeletedAt   time.Time     `json:"deleted_at,omitempty" gorm:"index"`
	UserID      types.UserId  `gorm:"not null" validate:"required,numeric,min=1"`
	OriginalUrl string        `gorm:"not null" validate:"required,url"`
	ShortUrl    string        `gorm:"unique;not null" validate:"required"`
}
