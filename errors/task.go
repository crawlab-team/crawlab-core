package errors

func NewTaskError(msg string) (err error) {
	return NewError(ErrorPrefixTask, msg)
}

var (
	ErrorTaskNotExists          = NewTaskError("not exists")
	ErrorTaskAlreadyExists      = NewTaskError("already exists")
	ErrorTaskInvalidType        = NewTaskError("invalid type")
	ErrorTaskProcessStillExists = NewTaskError("process still exists")
	ErrorTaskUnableToCancel     = NewTaskError("unable to cancel")
	ErrorTaskForbidden          = NewTaskError("forbidden")
)
