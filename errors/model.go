package errors

func NewModelError(msg string) (err error) {
	return NewError(ErrorPrefixModel, msg)
}

var ErrorModelInvalidType = NewModelError("invalid type")
var ErrorModelInvalidModelId = NewModelError("invalid model id")
var ErrorModelNotImplemented = NewModelError("not implemented")
var ErrorModelNotFound = NewModelError("not found")
var ErrorModelAlreadyExists = NewModelError("already exists")
var ErrorModelMissingRequiredData = NewModelError("missing required data")
