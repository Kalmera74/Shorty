package user

type UserCreateRequest struct {
	UserName string
	Email    string
}
type UserUpdateRequest struct {
	UserName string
	Email    string
}
type UserResponse struct {
	Id       uint
	UserName string
	Email    string
}
