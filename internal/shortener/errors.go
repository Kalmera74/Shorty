package shortener

type URLNotFoundError struct {
	Msg string
	Err error
}

func (e *URLNotFoundError) Error() string {
	return e.Msg
}

func (e *URLNotFoundError) Unwrap() error {
	return e.Err
}
