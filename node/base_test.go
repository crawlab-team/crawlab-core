package node

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var TestServiceDefault interfaces.NodeMasterService
var TestServiceMaster interfaces.NodeMasterService
var TestServiceWorker interfaces.NodeWorkerService

func initTest() {
	var err error

	// default service
	if err = ioutil.WriteFile(DefaultConfigPath, []byte("{\"key\":\"default\",\"is_master\":true}"), os.ModePerm); err != nil {
		panic(err)
	}
	if TestServiceDefault == nil {
		svc, err := NewService()
		if err != nil {
			panic(err)
		}
		TestServiceDefault = svc.(interfaces.NodeMasterService)
	}

	// master service
	if TestServiceMaster == nil {
		masterNodeConfigName := "config-master.json"
		masterNodeConfigPath := path.Join(DefaultConfigDirPath, masterNodeConfigName)
		if err = ioutil.WriteFile(masterNodeConfigPath, []byte("{\"key\":\"master\",\"is_master\":true}"), os.ModePerm); err != nil {
			panic(err)
		}
		svc, err := NewService()
		if err != nil {
			panic(err)
		}
		TestServiceMaster = svc.(interfaces.NodeMasterService)
	}

	// worker service
	if TestServiceWorker == nil {
		workerNodeConfigName := "config-worker.json"
		workerNodeConfigPath := path.Join(DefaultConfigDirPath, workerNodeConfigName)
		if err = ioutil.WriteFile(workerNodeConfigPath, []byte("{\"key\":\"worker\",\"is_worker\":false}"), os.ModePerm); err != nil {
			panic(err)
		}
		svc, err := NewService()
		if err != nil {
			panic(err)
		}
		TestServiceWorker = svc.(interfaces.NodeWorkerService)
	}
}

func setupTest(t *testing.T) {
	t.Cleanup(cleanupTest)
	initTest()
}

func cleanupTest() {
	_ = os.RemoveAll(DefaultConfigDirPath)
	TestServiceDefault = nil
	TestServiceMaster = nil
	TestServiceWorker = nil
}

type Test struct {
}

func (t *Test) Inject() error {
	return nil
}

func (t *Test) Setup() {
	initTest()
}

func (t *Test) Cleanup() {
	cleanupTest()
}

func NewServiceTest() (t *Test) {
	t = &Test{}
	_ = t.Inject()
	return t
}
