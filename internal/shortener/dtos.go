package shortener

type ShortenRequest struct {
	UserID uint   `json:"user_id"`
	Url    string `json:"original_url"`
}

type ShortenResponse struct {
	ShortID     uint   `json:"id"`
	OriginalUrl string `json:"original_url"`
	ShortUrl    string `json:"short_url"`
}
