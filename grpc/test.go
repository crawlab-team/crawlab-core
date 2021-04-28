package grpc

import (
	"github.com/crawlab-team/crawlab-core/models"
	node2 "github.com/crawlab-team/crawlab-core/node"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path"
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
	masterNodeConfigName := "config-master.json"
	masterNodeConfigPath := path.Join(node2.DefaultConfigDirPath, masterNodeConfigName)
	err = ioutil.WriteFile(masterNodeConfigPath, []byte("{\"key\":\"master\",\"is_master\":true}"), os.ModePerm)
	require.Nil(t, err)
	masterNodeService, err := node2.NewService(&node2.ServiceOptions{
		ConfigPath: masterNodeConfigPath,
	})
	if TestMasterService, err = NewService(&ServiceOptions{
		NodeService: masterNodeService,
		Local: NewAddress(&AddressOptions{
			Host: "localhost",
			Port: TestMasterPort,
		}),
	}); err != nil {
		panic(err)
	}

	workerNodeConfigName := "config-worker.json"
	workerNodeConfigPath := path.Join(node2.DefaultConfigDirPath, workerNodeConfigName)
	err = ioutil.WriteFile(workerNodeConfigPath, []byte("{\"key\":\"worker\",\"is_worker\":false}"), os.ModePerm)
	require.Nil(t, err)
	workerNodeService, err := node2.NewService(&node2.ServiceOptions{
		ConfigPath: workerNodeConfigPath,
	})
	if TestWorkerService, err = NewService(&ServiceOptions{
		NodeService: workerNodeService,
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
