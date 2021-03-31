package models

type userService struct {
	*Service
}

func NewUserService() (svc *userService) {
	return &userService{NewService(ModelIdUser)}
}

var UserService = NewUserService()
