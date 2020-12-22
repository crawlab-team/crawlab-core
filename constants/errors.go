package constants

import (
	"errors"
	e "github.com/crawlab-team/crawlab-core/errors"
	"net/http"
)

var (
	ErrorMongoError                = e.NewSystemOPError(1001, "system error:[mongo]%s", http.StatusInternalServerError)
	ErrorUserNotFound              = e.NewBusinessError(10001, "user not found.", http.StatusUnauthorized)
	ErrorUsernameOrPasswordInvalid = e.NewBusinessError(11001, "username or password invalid", http.StatusUnauthorized)
	ErrAlreadyExists               = errors.New("already exists")
	ErrNotExists                   = errors.New("not exists")
	ErrForbidden                   = errors.New("forbidden")
	ErrInvalidOptions              = errors.New("invalid options")
)
