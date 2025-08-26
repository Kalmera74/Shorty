package user

type UserCreateRequest struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

type UserUpdateRequest struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

type UserResponse struct {
	Id       uint   `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}
