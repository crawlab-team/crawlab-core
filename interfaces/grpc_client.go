package interfaces

import (
	"context"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"time"
)

type GrpcClient interface {
	GrpcBase
	WithConfigPath
	GetModelDelegateClient() grpc.ModelDelegateClient
	GetModelBaseServiceClient() grpc.ModelBaseServiceClient
	GetNodeClient() grpc.NodeServiceClient
	GetTaskClient() grpc.TaskServiceClient
	SetAddress(Address)
	SetTimeout(time.Duration)
	Context() (context.Context, context.CancelFunc)
	NewRequest(interface{}) *grpc.Request
	GetMessageChannel() chan *grpc.StreamMessage
	Restart() error
	NewModelBaseServiceRequest(ModelId, GrpcBaseServiceParams) (*grpc.Request, error)
	IsStarted() bool
	IsClosed() bool
}
