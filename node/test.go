package node

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var TestService *Service
var TestServiceMaster *Service
var TestServiceWorker *Service
var TestServiceStore *ServiceStore

func setupTest(t *testing.T) {
	var err error

	// default service
	TestService, err = NewService(nil)
	if err != nil {
		panic(err)
	}

	// master service
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

	// worker service
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

	// service store
	TestServiceStore = NewServiceStore()

	// cleanup
	t.Cleanup(cleanupTest)
}

func cleanupTest() {
	_ = os.RemoveAll(DefaultConfigDirPath)
}
