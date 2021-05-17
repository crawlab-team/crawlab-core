package errors

func NewTaskError(msg string) (err error) {
	return NewError(ErrorPrefixTask, msg)
}

var ErrorTaskNotExists = NewTaskError("not exists")
var ErrorTaskAlreadyExists = NewTaskError("already exists")
var ErrorTaskInvalidType = NewTaskError("invalid type")
var ErrorTaskProcessStillExists = NewTaskError("process still exists")
var ErrorTaskUnableToCancel = NewTaskError("unable to cancel")
var ErrorTaskForbidden = NewTaskError("forbidden")
