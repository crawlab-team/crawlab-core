package test

import (
	"context"
	gtest "github.com/crawlab-team/crawlab-core/grpc/test"
	ntest "github.com/crawlab-team/crawlab-core/node/test"
	"github.com/crawlab-team/crawlab-db/mongo"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"testing"
)

func TestMain(m *testing.M) {
	if err := mongo.InitMongo(); err != nil {
		panic(err)
	}
	T.TestNode.GetKey()
	grpcT, _ := gtest.NewTest()
	_ = grpcT.Client.Start()
	defer grpcT.Client.Stop()
	req := &grpc.Request{
		NodeKey: T.TestNode.GetKey(),
	}
	_, _ = grpcT.Client.GetNodeClient().Subscribe(context.Background(), req)

	m.Run()

	_ = ntest.T.ModelSvc.DropAll()
}
