package shortener

import "errors"

var (
	ErrInvalidShortModel     = errors.New("Invalid short model")
	ErrShortNotFound         = errors.New("Short not found")
	ErrShortenFailed         = errors.New("Failed to shorten URL")
	ErrInvalidShortenRequest = errors.New("Invalid shorten request")
	ErrShortDeleteFail       = errors.New("Failed to delete short URL")
)
