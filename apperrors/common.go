package apperrors

type InvalidIdError struct {
	Msg string
	Err error
}

func (e *InvalidIdError) Error() string {
	return e.Msg
}

func (e *InvalidIdError) Unwrap() error {
	return e.Err
}

type InvalidUrlError struct {
	Msg string
	Err error
}

func (e *InvalidUrlError) Error() string {
	return e.Msg
}

func (e *InvalidUrlError) Unwrap() error {
	return e.Err
}
