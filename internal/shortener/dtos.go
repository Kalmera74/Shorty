package shortener

type ShortenRequest struct {
	UserID  uint   `json:"user_id"`
	LongURL string `json:"long_url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}
