package errors

func NewControllerError(msg string) (err error) {
	return NewError(ErrorPrefixController, msg)
}

var ErrorControllerInvalidControllerId = NewControllerError("invalid controller id")
var ErrorControllerInvalidType = NewControllerError("invalid type")
var ErrorControllerAddError = NewControllerError("add error")
var ErrorControllerUpdateError = NewControllerError("update error")
var ErrorControllerDeleteError = NewControllerError("delete error")
var ErrorControllerNotImplemented = NewControllerError("not implemented")
var ErrorControllerNoModelService = NewControllerError("no model service")
var ErrorControllerRequestPayloadInvalid = NewControllerError("request payload invalid")
var ErrorControllerMissingInCache = NewControllerError("missing in cache")
var ErrorControllerNotCancellable = NewControllerError("not cancellable")
