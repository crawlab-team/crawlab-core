package models

type Service struct {
}

func NewService() (svc *Service) {
	return &Service{}
}

var RootModelService *Service
