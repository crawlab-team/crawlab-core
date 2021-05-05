package interfaces

type GrpcService interface {
	Injectable
	Init() error
	Stop() error
	GetServer() GrpcServer
	GetClient() GrpcClient
}
