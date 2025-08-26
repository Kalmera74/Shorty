package shortener

import "gorm.io/gorm"

type ShortenModel struct {
	gorm.Model
	UserID  uint   `gorm:"not null"`
	LongURL string `gorm:"not null"`
	ShortID string `gorm:"unique;not null"`
}
