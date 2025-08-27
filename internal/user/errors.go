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

type UserNotFoundError struct {
	Msg string
	Err error
}
func (e *UserNotFoundError) Error() string {
	return e.Msg
}
func (e *UserNotFoundError) Unwrap() error {
	return e.Err
}

type InValidUserCreateRequestError struct {
	Msg string
	Err error
}

func (e *InValidUserCreateRequestError) Error() string {
	return e.Msg
}
func (e *InValidUserCreateRequestError) Unwrap() error {
	return e.Err
}

type InValidUserUpdateRequestError struct {
	Msg string
	Err error
}

func (e *InValidUserUpdateRequestError) Error() string {
	return e.Msg
}
func (e *InValidUserUpdateRequestError) Unwrap() error {
	return e.Err
}
