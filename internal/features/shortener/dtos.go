package shortener

import "github.com/Kalmera74/Shorty/internal/types"

type ShortenRequest struct {
	UserID         types.UserId `json:"user_id" validate:"required,numeric,min=1"`
	Url            string       `json:"original_url" validate:"required,url"`
	CustomShortUrl *string      `json:"custom_short_url,omitempty"`
}
type ShortResponse struct {
	Id          types.ShortId   `json:"id"`
	OriginalUrl string `json:"original_url"`
	ShortUrl    string `json:"short_url"`
}

type SearchRequest struct {
	OriginalUrl *string       `json:"original_url,omitempty"`
	UserId      *types.UserId `json:"user_id,omitempty"`
	ShortUrl    *string       `json:"short_url,omitempty"`
}
