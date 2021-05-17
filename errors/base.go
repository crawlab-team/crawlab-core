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
	ErrorPrefixInject     = "inject"
	ErrorPrefixSpider     = "spider"
	ErrorPrefixFs         = "fs"
	ErrorPrefixTask       = "task"
)

type ErrorPrefix string

func NewError(prefix ErrorPrefix, msg string) (err error) {
	return errors.New(fmt.Sprintf("%s error: %s", prefix, msg))
}
