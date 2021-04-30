package node

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestNewService(t *testing.T) {
	setupTest(t)

	svc := TestServiceDefault
	require.NotNil(t, svc.cfg)
	require.NotEmpty(t, svc.cfg.Key)
	require.False(t, svc.cfg.IsMaster)
	require.DirExists(t, path.Dir(svc.opts.ConfigPath))
	require.FileExists(t, svc.opts.ConfigPath)

	nodeKey := "test-node-key"
	cfg := NewConfig(&ConfigOptions{
		Key:      nodeKey,
		IsMaster: true,
	})
	data, err := json.Marshal(cfg)
	require.Nil(t, err)
	err = ioutil.WriteFile(svc.opts.ConfigPath, data, os.ModePerm)
	require.Nil(t, err)

	svc2, err := NewService(nil)
	require.Nil(t, err)
	require.Equal(t, nodeKey, svc2.cfg.Key)
	require.True(t, svc2.cfg.IsMaster)
}

func TestService_GetNodeKey(t *testing.T) {
	setupTest(t)

	key := TestServiceDefault.GetNodeKey()
	require.NotEmpty(t, key)
}

func TestService_IsMaster(t *testing.T) {
	setupTest(t)

	isMaster := TestServiceDefault.IsMaster()
	require.False(t, isMaster)
}
