package interfaces

import grpc "github.com/crawlab-team/crawlab-grpc"

type GrpcSubscribe interface {
	GetStream() grpc.NodeService_SubscribeServer
	GetFinished() chan bool
}
