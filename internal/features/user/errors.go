package user

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")

	// Validation errors
	ErrInvalidUserCreateRequest = errors.New("invalid user create request")
	ErrInvalidUserUpdateRequest = errors.New("invalid user update request")
	ErrInvalidCredentials       = errors.New("invalid credentials")
)
