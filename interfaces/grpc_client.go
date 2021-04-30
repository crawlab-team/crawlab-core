package interfaces

import grpc "github.com/crawlab-team/crawlab-grpc"

type GrpcClient interface {
	GrpcBase
	GetModelDelegateClient() grpc.ModelDelegateServiceClient
	GetNodeClient() grpc.NodeServiceClient
	GetTaskClient() grpc.TaskServiceClient
}
