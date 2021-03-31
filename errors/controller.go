package errors

func NewControllerError(msg string) (err error) {
	return NewError(ErrorPrefixController, msg)
}

var ErrorControllerInvalidControllerId = NewControllerError("invalid controller id")
var ErrorControllerAddError = NewControllerError("add error")
var ErrorControllerUpdateError = NewControllerError("update error")
var ErrorControllerDeleteError = NewControllerError("delete error")
var ErrorControllerNotImplemented = NewControllerError("not implemented")
