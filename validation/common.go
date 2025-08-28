package validation

import (
	"net/url"

	"github.com/Kalmera74/Shorty/apperrors"
)

func ValidateID(id uint) error {
	if id > 0 {
		return nil
	}

	return &apperrors.InvalidIdError{Msg: "ID cannot be less then or equal to 0"}
}

func ValidateUrl(URL string) error {
	if _, err := url.Parse(URL); err != nil {
		return &apperrors.InvalidUrlError{Msg: "Given Url is not a valid Url"}
	}

	return nil
}
