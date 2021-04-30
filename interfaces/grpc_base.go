package interfaces

type GrpcBase interface {
	Init() (err error)
	Start() (err error)
	Stop() (err error)
	Register() (err error)
}
