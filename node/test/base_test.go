package test

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNodeServices_Master_Worker(t *testing.T) {
	T, _ = NewTest()
	T.Setup(t)
	startMasterWorker(t)

	// validate master
	masterNodeKey := T.MasterSvc.GetConfigService().GetNodeKey()
	masterNode, err := T.ModelSvc.GetNodeByKey(masterNodeKey, nil)
	require.Nil(t, err)
	require.Equal(t, constants.NodeStatusOnline, masterNode.Status)
	require.Equal(t, masterNodeKey, masterNode.Key)
	require.True(t, masterNode.IsMaster)

	// validate worker
	workerNodeKey := T.WorkerSvc.GetConfigService().GetNodeKey()
	workerNode, err := T.ModelSvc.GetNodeByKey(workerNodeKey, nil)
	require.Nil(t, err)
	require.Equal(t, constants.NodeStatusOnline, workerNode.Status)
	require.Equal(t, workerNodeKey, workerNode.Key)
	require.False(t, workerNode.IsMaster)

	stopMasterWorker(t)
}

func TestNodeServices_Default(t *testing.T) {
	T, _ = NewTest()
	T.Setup(t)

	go T.DefaultSvc.Start()
	time.Sleep(1 * time.Second)

	// validate default
	defaultNodeKey := T.DefaultSvc.GetConfigService().GetNodeKey()
	defaultNode, err := T.ModelSvc.GetNodeByKey(defaultNodeKey, nil)
	require.Nil(t, err)
	require.Equal(t, constants.NodeStatusOnline, defaultNode.Status)
	require.Equal(t, defaultNodeKey, defaultNode.Key)
	require.True(t, defaultNode.IsMaster)

	T.DefaultSvc.Stop()
	time.Sleep(1 * time.Second)
}

func TestNodeServices_Monitor(t *testing.T) {
	T, _ = NewTest()
	T.Setup(t)
	startMasterWorkerMonitor(t)
	time.Sleep(3 * time.Second)

	// stop worker
	T.WorkerSvcMonitor.Stop()
	time.Sleep(5 * time.Second)

	// validate
	require.True(t, T.MasterSvcMonitor.GetServer().IsStopped())

	stopMasterWorkerMonitor(t)
}

func startMasterWorker(t *testing.T) {
	go T.MasterSvc.Start()
	time.Sleep(1 * time.Second)
	go T.WorkerSvc.Start()
	time.Sleep(1 * time.Second)
}

func stopMasterWorker(t *testing.T) {
	go T.MasterSvc.Stop()
	time.Sleep(1 * time.Second)
	go T.WorkerSvc.Stop()
	time.Sleep(1 * time.Second)
}

func startMasterWorkerMonitor(t *testing.T) {
	go T.MasterSvcMonitor.Start()
	time.Sleep(1 * time.Second)
	go T.WorkerSvcMonitor.Start()
	time.Sleep(1 * time.Second)
}

func stopMasterWorkerMonitor(t *testing.T) {
	go T.MasterSvcMonitor.Stop()
	time.Sleep(1 * time.Second)
	go T.WorkerSvcMonitor.Stop()
	time.Sleep(1 * time.Second)
}
