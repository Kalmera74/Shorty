package validation

import (
	"fmt"
	"net/url"

	"github.com/Kalmera74/Shorty/internal/apperrors"
)

func ValidateID(id uint) error {
	if id > 0 {
		return nil
	}

	return &apperrors.InvalidIdError{Msg: "ID cannot be less then or equal to 0"}
}

func ValidateUrl(URL string) error {
	if URL == "" {
		return &apperrors.InvalidUrlError{Msg: "Given Url cannot be empty."}
	}

	parsedUrl, err := url.Parse(URL)
	if err != nil {
		return &apperrors.InvalidUrlError{
			Msg: fmt.Sprintf("Given Url is not a valid Url. Reason: %v", err.Error()),
			Err: err,
		}
	}

	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		return &apperrors.InvalidUrlError{
			Msg: "Url must have a valid scheme (http or https).",
		}
	}
	
	if parsedUrl.Host == "" {
		return &apperrors.InvalidUrlError{
			Msg: "Url is missing a host.",
		}
	}

	return nil
}