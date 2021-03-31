package models

type tokenService struct {
	*Service
}

func NewTokenService() (svc *tokenService) {
	return &tokenService{NewService(ModelIdToken)}
}

var TokenService = NewTokenService()
