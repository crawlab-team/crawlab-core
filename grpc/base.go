package grpc

type ServiceInterface interface {
	Init() (err error)
	Start()
	Stop()
	Register()
}
