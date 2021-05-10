package interfaces

import (
	"github.com/crawlab-team/crawlab-core/entity"
	grpc "github.com/crawlab-team/crawlab-grpc"
)

type GrpcServer interface {
	GrpcBase
	SetAddress(Address)
	GetInboundStreamMessageChannel(key string) (chan *grpc.StreamMessage, error)
	GetOutboundStreamMessageChannel(key string) (chan *grpc.StreamMessage, error)
	GetSubscribe(key string) (sub *entity.GrpcSubscribe, err error)
	SetSubscribe(key string, sub *entity.GrpcSubscribe)
	DeleteSubscribe(key string)
}
