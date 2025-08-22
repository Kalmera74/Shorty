package user

type UserCreate struct {
	UserName string
	Email    string
}

type User struct {
	Id       uint
	UserName string
	Email string
}

