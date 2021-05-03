package interfaces

import grpc "github.com/crawlab-team/crawlab-grpc"

type NodeMasterService interface {
	NodeService
	GetInboundStreamMessageChannel(key string) (chan *grpc.StreamMessage, error)
	GetOutboundStreamMessageChannel(key string) (chan *grpc.StreamMessage, error)
	Monitor()
	GetStream(key string) (stream grpc.NodeService_StreamServer, err error)
	SetStream(key string, stream grpc.NodeService_StreamServer)
	DeleteStream(key string)
}
