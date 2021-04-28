package grpc

type ServiceInterface interface {
	Init() (err error)
	Start() (err error)
	Stop() (err error)
	Register() (err error)
}
