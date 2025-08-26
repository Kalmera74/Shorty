package user

import "gorm.io/gorm"

type UserModel struct{

	gorm.Model
	UserName string
	Email string
}
