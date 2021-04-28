package grpc

import (
	"github.com/crawlab-team/crawlab-core/models"
	"testing"
)

var TestMasterService *Service
var TestWorkerService *Service

var TestMasterPort = "9876"
var TestWorkerPort = "9877"

func setupTest(t *testing.T) {
	if err := models.InitModelServices(); err != nil {
		panic(err)
	}

	var err error
	if TestMasterService, err = NewService(&ServiceOptions{
		Local: NewAddress(&AddressOptions{
			Host: "localhost",
			Port: TestMasterPort,
		}),
	}); err != nil {
		panic(err)
	}

	if TestWorkerService, err = NewService(&ServiceOptions{
		Local: NewAddress(&AddressOptions{
			Host: "localhost",
			Port: TestWorkerPort,
		}),
		Remotes: []Address{
			NewAddress(&AddressOptions{
				Host: "localhost",
				Port: TestMasterPort,
			}),
		},
	}); err != nil {
		panic(err)
	}

	if err = TestMasterService.AddClient(&ClientOptions{
		Address: NewAddress(&AddressOptions{
			Host: "localhost",
			Port: TestWorkerPort,
		}),
	}); err != nil {
		panic(err)
	}

	t.Cleanup(cleanupTest)
}

func cleanupTest() {
	_ = models.NodeService.Delete(nil)
	_ = TestMasterService.Stop()
	_ = TestWorkerService.Stop()
}
