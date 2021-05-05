package interfaces

type GrpcServer interface {
	GrpcBase
	SetAddress(Address)
}
