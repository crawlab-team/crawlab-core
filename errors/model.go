package errors

func NewModelError(msg string) (err error) {
	return NewError(ErrorPrefixModel, msg)
}

var ErrorModelInvalidType = NewModelError("invalid type")
var ErrorModelInvalidModelId = NewModelError("invalid model id")
var ErrorModelNotImplemented = NewModelError("not implemented")
var ErrorModelNotFound = NewModelError("not found")
