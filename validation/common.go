package validation

import (
	"github.com/Kalmera74/Shorty/errors"
)

func IsValidID(id uint) error {
	if id > 0 {
		return nil
	}

	return &errors.InvalidIdError{Msg: "ID cannot be less then or equal to 0"}
}
