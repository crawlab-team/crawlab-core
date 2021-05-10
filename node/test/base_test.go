package test

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNodeServices(t *testing.T) {
	T.Setup(t)

	go T.MasterSvc.Start()
	time.Sleep(1 * time.Second)
	go T.WorkerSvc.Start()
	time.Sleep(5 * time.Second)

	// validate master
	masterNodeKey := T.WorkerSvc.GetConfigService().GetNodeKey()
	masterNode, err := T.ModelSvc.GetNodeByKey(masterNodeKey, nil)
	require.Nil(t, err)
	require.Equal(t, constants.NodeStatusOnline, masterNode.Status)

	// validate worker
	workerNodeKey := T.WorkerSvc.GetConfigService().GetNodeKey()
	workerNode, err := T.ModelSvc.GetNodeByKey(workerNodeKey, nil)
	require.Nil(t, err)
	require.Equal(t, constants.NodeStatusOnline, workerNode.Status)
}
