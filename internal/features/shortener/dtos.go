package shortener

type ShortenRequest struct {
	UserID         uint    `json:"user_id" validate:"required,numeric,min=1"`
	Url            string  `json:"original_url" validate:"required,url"`
	CustomShortUrl *string `json:"custom_short_url,omitempty" validate:"max=8,omitempty"`
}
type ShortenResponse struct {
	Id          uint   `json:"id"`
	OriginalUrl string `json:"original_url"`
	ShortUrl    string `json:"short_url"`
}

type SearchRequest struct {
	OriginalUrl *string `json:"original_url,omitempty" validate:"omitempty,url"`
}
