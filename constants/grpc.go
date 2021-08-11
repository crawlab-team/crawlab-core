package constants

const (
	DefaultGrpcServerHost       = "0.0.0.0"
	DefaultGrpcServerPort       = "9666"
	DefaultGrpcClientRemoteHost = "localhost"
	DefaultGrpcClientRemotePort = DefaultGrpcServerPort
	DefaultGrpcAuthKey          = "Crawlab2021!"
)

const (
	GrpcHeaderAuthorization = "authorization"
)

const (
	GrpcSubscribeTypeNode   = "node"
	GrpcSubscribeTypePlugin = "plugin"
)
