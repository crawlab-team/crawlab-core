package entity

import grpc "github.com/crawlab-team/crawlab-grpc"

type GrpcSubscribe struct {
	Stream   grpc.NodeService_SubscribeServer
	Finished chan<- bool
}
