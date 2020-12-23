package services

type TaskGrpcServiceInterface interface {
	GetTask() (err error)
}

type TaskGrpcService struct {
}
