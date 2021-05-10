package entity

import grpc "github.com/crawlab-team/crawlab-grpc"

type GrpcSubscribe struct {
	Stream   grpc.NodeService_SubscribeServer
	Finished chan bool
}

func (sub *GrpcSubscribe) GetStream() grpc.NodeService_SubscribeServer {
	return sub.Stream
}

func (sub *GrpcSubscribe) GetFinished() chan bool {
	return sub.Finished
}
