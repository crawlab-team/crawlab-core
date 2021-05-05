package grpc

import (
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/utils"
	"os"
	"testing"
)

var TestServiceMaster interfaces.GrpcService
var TestServiceWorker interfaces.GrpcService

var TestPortMaster = "9876"

func setupTest(t *testing.T) {
	var err error

	if err := os.Setenv("TEST_ENV", "master,worker"); err != nil {
		panic(err)
	}

	if TestServiceMaster, err = NewService(
		WithEnv("master"),
		WithLocal(entity.NewAddress(&entity.AddressOptions{Port: TestPortMaster})),
	); err != nil {
		panic(err)
	}

	if TestServiceWorker, err = NewService(
		WithEnv("worker"),
		WithRemote(entity.NewAddress(&entity.AddressOptions{Port: TestPortMaster})),
	); err != nil {
		panic(err)
	}

	t.Cleanup(cleanupTest)
}

func cleanupTest() {
	var nodeSvcTest interfaces.NodeServiceTest
	utils.MustResolveModule("", nodeSvcTest)
	nodeSvcTest.Cleanup()
}
