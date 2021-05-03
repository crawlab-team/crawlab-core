package interfaces

import grpc "github.com/crawlab-team/crawlab-grpc"

type GrpcNodeServer interface {
	HandleSendStreamMessage(stream grpc.NodeService_StreamServer, chMsg chan *grpc.StreamMessage)
	AllowRegister(nodeKey string) (res bool)
}
