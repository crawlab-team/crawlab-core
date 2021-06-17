package interfaces

import "github.com/dgrijalva/jwt-go"

type UserService interface {
	Init() (err error)
	SetJwtSecret(secret string)
	SetJwtSigningMethod(method jwt.SigningMethod)
	Create(opts *UserCreateOptions) (err error)
	Login(opts *UserLoginOptions) (token string, u User, err error)
	CheckToken(token string) (u User, err error)
}
