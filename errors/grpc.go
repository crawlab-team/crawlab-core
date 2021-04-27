package errors

func NewGrpcError(msg string) (err error) {
	return NewError(ErrorPrefixGrpc, msg)
}

var ErrorGrpcClientNotExists = NewGrpcError("client not exists")
var ErrorGrpcClientAlreadyExists = NewGrpcError("client already exists")
