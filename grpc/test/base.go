package test

import (
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/grpc/client"
	"github.com/crawlab-team/crawlab-core/grpc/server"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/node/test"
	"testing"
)

func init() {
	var err error
	T, err = NewTest()
	if err != nil {
		panic(err)
	}
}

type Test struct {
	Server interfaces.GrpcServer
	Client interfaces.GrpcClient

	MasterNodeInfo *entity.NodeInfo
	WorkerNodeInfo *entity.NodeInfo
}

func (t *Test) Setup(t2 *testing.T) {
	test.T.Cleanup()
	t2.Cleanup(t.Cleanup)

	if err := t.Server.Start(); err != nil {
		panic(err)
	}
	if err := t.Client.Start(); err != nil {
		panic(err)
	}
}

func (t *Test) Cleanup() {
	if err := t.Client.Stop(); err != nil {
		panic(err)
	}
	if err := t.Server.Stop(); err != nil {
		panic(err)
	}
	test.T.Cleanup()
}

var T *Test

func NewTest() (res *Test, err error) {
	// test
	t := &Test{}

	// server
	t.Server, err = server.NewServer(server.WithConfigPath(test.T.MasterSvc.GetConfigPath()))
	if err != nil {
		return nil, err
	}

	// client
	t.Client, err = client.NewClient(client.WithConfigPath(test.T.WorkerSvc.GetConfigPath()))
	if err != nil {
		return nil, err
	}

	// master node info
	t.MasterNodeInfo = &entity.NodeInfo{
		Key:      "master",
		IsMaster: true,
	}

	// worker node info
	t.WorkerNodeInfo = &entity.NodeInfo{
		Key:      "worker",
		IsMaster: false,
	}

	return t, nil
}
