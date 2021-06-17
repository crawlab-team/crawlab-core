package interfaces

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	Init() (err error)
	SetJwtSecret(secret string)
	SetJwtSigningMethod(method jwt.SigningMethod)
	Create(opts *UserCreateOptions) (err error)
	Login(opts *UserLoginOptions) (token string, u User, err error)
	CheckToken(token string) (u User, err error)
	ChangePassword(id primitive.ObjectID, password string) (err error)
}
