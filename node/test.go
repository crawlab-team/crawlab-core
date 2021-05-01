package node

import (
	"github.com/crawlab-team/crawlab-core/store"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var TestServiceDefault *Service
var TestServiceMaster *Service
var TestServiceWorker *Service
var TestServiceStore *store.NodeServiceStoreType

func initTest() {
	var err error

	// default service
	if err = ioutil.WriteFile(DefaultConfigPath, []byte("{\"key\":\"default\",\"is_master\":true}"), os.ModePerm); err != nil {
		panic(err)
	}
	if TestServiceDefault == nil {
		TestServiceDefault, err = NewService(nil)
		if err != nil {
			panic(err)
		}
	}

	// master service
	if TestServiceMaster == nil {
		masterNodeConfigName := "config-master.json"
		masterNodeConfigPath := path.Join(DefaultConfigDirPath, masterNodeConfigName)
		if err = ioutil.WriteFile(masterNodeConfigPath, []byte("{\"key\":\"master\",\"is_master\":true}"), os.ModePerm); err != nil {
			panic(err)
		}
		if TestServiceMaster, err = NewService(&ServiceOptions{
			ConfigPath: masterNodeConfigPath,
		}); err != nil {
			panic(err)
		}
	}

	// worker service
	if TestServiceWorker == nil {
		workerNodeConfigName := "config-worker.json"
		workerNodeConfigPath := path.Join(DefaultConfigDirPath, workerNodeConfigName)
		if err = ioutil.WriteFile(workerNodeConfigPath, []byte("{\"key\":\"worker\",\"is_worker\":false}"), os.ModePerm); err != nil {
			panic(err)
		}
		if TestServiceWorker, err = NewService(&ServiceOptions{
			ConfigPath: workerNodeConfigPath,
		}); err != nil {
			panic(err)
		}
	}

	// service store
	if TestServiceStore == nil {
		TestServiceStore = store.NewNodeServiceStore()
	}
}

func setupTest(t *testing.T) {
	initTest()

	// cleanup
	t.Cleanup(cleanupTest)
}

func cleanupTest() {
	_ = os.RemoveAll(DefaultConfigDirPath)
	TestServiceDefault = nil
	TestServiceMaster = nil
	TestServiceWorker = nil
	TestServiceStore = nil
}

func SetupTest() {
	initTest()
}

func CleanupTest() {
	cleanupTest()
}
