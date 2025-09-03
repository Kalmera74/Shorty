package analytics

import "errors"

var (
	ErrClickNotFound   = errors.New("Click not found")
	ErrClicksNotFound  = errors.New("No clicks found")
	ErrClickCreateFail = errors.New("Failed to create click")
)
