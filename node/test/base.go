package test

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	service2 "github.com/crawlab-team/crawlab-core/models/service"
	"github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/crawlab-core/node/service"
	"github.com/crawlab-team/crawlab-core/utils"
	"go.uber.org/dig"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

func init() {
	var err error
	T, err = NewTest()
	if err != nil {
		panic(err)
	}
}

var T *Test

type Test struct {
	MasterSvc interfaces.NodeMasterService
	WorkerSvc interfaces.NodeWorkerService
	ModelSvc  service2.ModelService
}

func NewTest() (res *Test, err error) {
	// test
	t := &Test{}

	// master config
	masterNodeConfigName := "config-master.json"
	masterNodeConfigPath := path.Join(config.DefaultConfigDirPath, masterNodeConfigName)
	if err := ioutil.WriteFile(masterNodeConfigPath, []byte("{\"key\":\"master\",\"is_master\":true}"), os.ModePerm); err != nil {
		return nil, err
	}

	// worker config
	workerNodeConfigName := "config-worker.json"
	workerNodeConfigPath := path.Join(config.DefaultConfigDirPath, workerNodeConfigName)
	if err = ioutil.WriteFile(workerNodeConfigPath, []byte("{\"key\":\"worker\",\"is_worker\":false}"), os.ModePerm); err != nil {
		return nil, err
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(service.ProvideMasterService(masterNodeConfigPath, service.WithMonitorInterval(3*time.Second))); err != nil {
		return nil, err
	}
	if err := c.Provide(service.ProvideWorkerService(workerNodeConfigPath, service.WithHeartbeatInterval(1*time.Second))); err != nil {
		return nil, err
	}
	if err := c.Provide(service2.NewService); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(masterSvc interfaces.NodeMasterService, workerSvc interfaces.NodeWorkerService, modelSvc service2.ModelService) {
		t.MasterSvc = masterSvc
		t.WorkerSvc = workerSvc
		t.ModelSvc = modelSvc
	}); err != nil {
		return nil, err
	}

	// visualize dependencies
	if err := utils.VisualizeContainer(c); err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Test) Setup(t2 *testing.T) {
	t.Cleanup()
	t2.Cleanup(t.Cleanup)
}

func (t *Test) Cleanup() {
	if err := t.ModelSvc.DropAll(); err != nil {
		panic(err)
	}
}
