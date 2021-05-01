package errors

func NewStoreError(msg string) (err error) {
	return NewError(ErrorPrefixStore, msg)
}

var ErrorStoreEmptyValue = NewStoreError("empty value")
var ErrorStoreNotExists = NewStoreError("not exists")
var ErrorStoreInvalidType = NewStoreError("invalid type")
