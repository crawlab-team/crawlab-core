package interfaces

type GrpcService interface {
	Init() error
	Stop() error
	GetServer() GrpcServer
	GetClient(Address) (GrpcClient, error)
	GetDefaultClient() (GrpcClient, error)
	MustGetDefaultClient() GrpcClient
	GetAllClients() ([]GrpcClient, error)
	AddClient(Options) error
	DeleteClient(Address) error
}
