package user

type UserCreateRequest struct {
	UserName string `json:"user_name" validate:"required,min=3,max=10" `
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=30" `
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=30"`
}

type UserLoginResponse struct {
	Token string
}
type UserUpdateRequest struct {
	UserName *string `json:"user_name,omitempty" validate:"min=3,max=30,omitempty"`
	Email    *string `json:"email,omitempty" validate:"email,omitempty"`
	Password *string `json:"password,omitempty" validate:"min=5,max=30"`
}

type UserResponse struct {
	Id       uint   `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}
