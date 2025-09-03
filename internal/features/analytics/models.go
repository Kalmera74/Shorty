package analytics

import (
	"time"

	"github.com/Kalmera74/Shorty/internal/features/shortener"
	"github.com/Kalmera74/Shorty/internal/types"
)

type ClickModel struct {
	ID        types.ClickId        `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	DeletedAt time.Time            `json:"deleted_at,omitempty" gorm:"index"`
	ShortID   types.ShortId        `json:"short_id" validate:"required"`
	Short     shortener.ShortModel `json:"short,omitempty" gorm:"constraint:OnDelete:CASCADE;"`
	IpAddress string               `json:"ip_address" validate:"required,ip"`
	UserAgent string               `json:"user_agent" validate:"required"`
}
