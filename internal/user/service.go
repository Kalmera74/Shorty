package user

import "fmt"

type UserService struct {
	UserStore UserStore
}

func (s *UserService) GetAllUsers() ([]User, error) {
	userModels, err := s.UserStore.GetAll()
	if err != nil {
		return nil, &UserError{msg: fmt.Sprintf("Could not retrive users. Reason: %v", err.Error()), err: err}
	}

	users := make([]User, len(userModels))

	for _, user := range userModels {
		convertedUser := User{Id: user.ID, UserName: user.UserName, Email: user.Email}
		users = append(users, convertedUser)

	}
	return users, nil
}
