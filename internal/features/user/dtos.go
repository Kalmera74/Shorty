package user

type UserCreateRequest struct {
	UserName string `json:"user_name" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
}

type UserUpdateRequest struct {
	UserName *string `json:"user_name,omitempty" validate:"min=3,max=30,omitempty"`
	Email    *string `json:"email,omitempty" validate:"email,omitempty"`
}

type UserResponse struct {
	Id       uint   `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}
