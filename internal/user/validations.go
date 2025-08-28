package user

import (
	"net/mail"
)

func (u UserCreateRequest) Validate() error {

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return &InValidUserCreateRequestError{Msg: "User has no valid email address", Err: nil}
	}

	return nil
}

func (u UserUpdateRequest) Validate() error {

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return &InValidUserCreateRequestError{Msg: "User has no valid email address", Err: nil}
	}

	return nil
}
