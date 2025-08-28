package shortener

import "gorm.io/gorm"

type ShortModel struct {
	gorm.Model
	UserID      uint   `gorm:"not null"`
	OriginalUrl string `gorm:"not null"`
	ShortUrl    string `gorm:"unique;not null"`
}
