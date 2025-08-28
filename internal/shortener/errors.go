package shortener

type InvalidShortModelError struct {
	Msg string
	Err error
}

func (e *InvalidShortModelError) Error() string {
	return e.Msg
}
func (e *InvalidShortModelError) Unwrap() error {
	return e.Err
}

type ShortNotFoundError struct {
	Msg string
	Err error
}

func (e *ShortNotFoundError) Error() string {
	return e.Msg
}
func (e *ShortNotFoundError) Unwrap() error {
	return e.Err
}

type ShortenError struct {
	Msg string
	Err error
}

func (e *ShortenError) Error() string {
	return e.Msg
}
func (e *ShortenError) Unwrap() error {
	return e.Err
}

type InvalidShortenRequestError struct {
	Msg string
	Err error
}

func (e *InvalidShortenRequestError) Error() string {
	return e.Msg
}
func (e *InvalidShortenRequestError) Unwrap() error {
	return e.Err
}
