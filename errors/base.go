package errors

import (
	"errors"
	"fmt"
)

const (
	ErrorPrefixController = "controller"
	ErrorPrefixModel      = "model"
	ErrorPrefixFilter     = "filter"
	ErrorPrefixHttp       = "http"
	ErrorPrefixGrpc       = "grpc"
	ErrorPrefixNode       = "node"
)

type ErrorPrefix string

func NewError(prefix ErrorPrefix, msg string) (err error) {
	return errors.New(fmt.Sprintf("%s error: %s", prefix, msg))
}
