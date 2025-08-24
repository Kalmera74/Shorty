package user

type UserStore interface {
	GetAll() ([]UserModel, error)
	Get(id uint) (UserModel, error)
	Update(id uint, user UserModel) error
	Delete(id uint) error
}
