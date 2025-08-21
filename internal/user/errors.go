package user

type UserError struct {
	msg string
	err error
}

func (e *UserError) Error() string {
	return e.msg
}

func (e *UserError) Unwrap() error {
	return e.err
}
