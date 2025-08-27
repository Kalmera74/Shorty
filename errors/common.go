package errors

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

