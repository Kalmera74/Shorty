package shortener

type ShortenRequest struct {
	UserID uint   `json:"user_id"`
	Url    string `json:"original_url"`
}

type ShortenResponse struct {
	Id          uint   `json:"id"`
	OriginalUrl string `json:"original_url"`
	ShortUrl    string `json:"short_url"`
}

type SearchRequest struct {
	OriginalUrl *string `json:"original_url,omitempty"`
}
