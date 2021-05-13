package errors

func NewFsError(msg string) (err error) {
	return NewError(ErrorPrefixFs, msg)
}

var ErrorFsForbidden = NewFsError("forbidden")
