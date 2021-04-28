package errors

func NewNodeError(msg string) (err error) {
	return NewError(ErrorPrefixNode, msg)
}

var ErrorNodeUnregistered = NewNodeError("unregistered")
var ErrorNodeServiceNotExists = NewNodeError("service not exists")
var ErrorNodeInvalidType = NewNodeError("invalid type")
