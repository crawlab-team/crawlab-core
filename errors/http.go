package errors

import "errors"

var ErrorHttpBadRequest = errors.New("bad request")
var ErrorHttpUnauthorized = errors.New("unauthorized")
var ErrorHttpNotFound = errors.New("not found")
