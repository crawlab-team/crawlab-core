package interfaces

import (
	grpc "github.com/crawlab-team/crawlab-grpc"
	"time"
)

type GrpcClient interface {
	GrpcBase
	GetModelDelegateClient() grpc.ModelDelegateServiceClient
	GetNodeClient() grpc.NodeServiceClient
	GetTaskClient() grpc.TaskServiceClient
	SetAddress(Address)
	SetTimeout(time.Duration)
}
