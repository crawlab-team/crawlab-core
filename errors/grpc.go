package errors

func NewGrpcError(msg string) (err error) {
	return NewError(ErrorPrefixGrpc, msg)
}

var ErrorGrpcClientFailedToStart = NewGrpcError("client failed to start")
var ErrorGrpcServerFailedToListen = NewGrpcError("server failed to listen")
var ErrorGrpcServerFailedToServe = NewGrpcError("server failed to serve")
var ErrorGrpcClientNotExists = NewGrpcError("client not exists")
var ErrorGrpcClientAlreadyExists = NewGrpcError("client already exists")
var ErrorGrpcInvalidType = NewGrpcError("invalid type")
var ErrorGrpcNotAllowed = NewGrpcError("not allowed")
var ErrorGrpcSubscribeNotExists = NewGrpcError("subscribe not exists")
