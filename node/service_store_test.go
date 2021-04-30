package node

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServiceStore_Set(t *testing.T) {
	var err error

	setupTest(t)

	err = TestServiceStore.Set("master", TestServiceMaster)
	require.Nil(t, err)

	_, ok := TestServiceStore.m.Load("master")
	require.True(t, ok)
}

func TestServiceStore_Get(t *testing.T) {
	var err error

	setupTest(t)

	err = TestServiceStore.Set("master", TestServiceMaster)
	require.Nil(t, err)
	err = TestServiceStore.Set("worker", TestServiceWorker)
	require.Nil(t, err)

	res, err := TestServiceStore.Get("master")
	require.Nil(t, err)
	svcMaster, ok := res.(interfaces.NodeService)
	require.True(t, ok)
	require.True(t, svcMaster.IsMaster())
	require.Equal(t, "master", svcMaster.GetNodeKey())

	res, err = TestServiceStore.Get("worker")
	require.Nil(t, err)
	svcWorker, ok := res.(interfaces.NodeService)
	require.True(t, ok)
	require.False(t, svcWorker.IsMaster())
	require.Equal(t, "worker", svcWorker.GetNodeKey())
}

func TestServiceStore_GetDefault(t *testing.T) {
	var err error

	setupTest(t)

	err = TestServiceStore.Set("default", TestService)
	require.Nil(t, err)

	res, err := TestServiceStore.GetDefault()
	require.Nil(t, err)
	svcDefault, ok := res.(interfaces.NodeService)
	require.True(t, ok)
	require.NotEmpty(t, svcDefault.GetNodeKey())
}
