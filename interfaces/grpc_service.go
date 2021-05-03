package interfaces

type GrpcService interface {
	Init() error
	Stop() error
	GetServer() GrpcServer
	GetClient() GrpcClient
}
