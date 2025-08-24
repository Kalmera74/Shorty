package user

type UserError struct {
	Msg string
	Err error
}

func (e *UserError) Error() string {
	return e.Msg
}

func (e *UserError) Unwrap() error {
	return e.Err
}
