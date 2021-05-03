package grpc

import (
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/models"
	"github.com/crawlab-team/crawlab-core/node"
	"github.com/crawlab-team/crawlab-core/store"
	"testing"
)

var TestServiceMaster *Service
var TestServiceWorker *Service

var TestPortMaster = "9876"
var TestPortWorker = "9877"

func setupTest(t *testing.T) {
	var err error

	if err := models.InitModels(); err != nil {
		panic(err)
	}

	if err := node.ForceInitNode(); err != nil {
		panic(err)
	}

	if err := InitGrpc(); err != nil {
		panic(err)
	}

	node.SetupTest()
	store.NodeServiceStore = node.TestServiceStore

	if err := node.TestServiceStore.Set("master", node.TestServiceMaster); err != nil {
		panic(err)
	}
	TestServiceMaster, err = NewService(&ServiceOptions{
		NodeServiceKey: "master",
		Local:          entity.NewAddress(&entity.AddressOptions{Port: TestPortMaster}),
	})
	if err != nil {
		panic(err)
	}

	if err := node.TestServiceStore.Set("worker", node.TestServiceWorker); err != nil {
		panic(err)
	}
	TestServiceWorker, err = NewService(&ServiceOptions{
		NodeServiceKey: "worker",
		Local:          entity.NewAddress(&entity.AddressOptions{Port: TestPortWorker}),
		Remote:         entity.NewAddress(&entity.AddressOptions{Port: TestPortMaster}),
	})
	if err != nil {
		panic(err)
	}

	t.Cleanup(cleanupTest)
}

func cleanupTest() {
	_ = models.MustGetService(interfaces.ModelIdNode).Delete(nil)
	_ = store.GrpcService.Stop()
	_ = TestServiceMaster.Stop()
	_ = TestServiceWorker.Stop()
}
