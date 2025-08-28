package shortener

import (
	"github.com/Kalmera74/Shorty/validation"
)

func (s *ShortModel) Validate() error {
	if s.OriginalUrl == "" {
		return &InvalidShortModelError{Msg: "Original URL cannot be nil or empty"}
	}
	if s.ShortUrl == "" {

		return &InvalidShortModelError{Msg: "Short URL cannot be nil or empty"}
	}

	if err := validation.ValidateID(s.UserID); err != nil {

		return &InvalidShortModelError{Msg: "User ID cannot be nil or empty", Err: err}
	}

	return nil
}

func (s *ShortenRequest) Validate() error {
	if s.Url == "" {
		return &InvalidShortenRequestError{Msg: "Given URL cannot be nil or empty"}
	}

	if err := validation.ValidateUrl(s.Url); err != nil {
		return &InvalidShortenRequestError{Msg: "Given URL is not a valid url", Err: err}
	}

	return nil
}
