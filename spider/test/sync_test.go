package test

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/fs"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestFsService_SyncToFs(t *testing.T) {
	var err error
	T.Setup(t)

	// save file to local
	filePath := path.Join(T.masterFsSvc.GetWorkspacePath(), T.scriptName)
	err = ioutil.WriteFile(filePath, []byte(T.script), os.ModePerm)
	require.Nil(t, err)

	// commit
	err = T.masterFsSvc.GetFsService().Commit("initial commit")
	require.Nil(t, err)

	// sync to fs
	err = T.masterSyncSvc.SyncToFs(T.s.Id)
	require.Nil(t, err)

	// validate
	remotePath := fmt.Sprintf("%s/%s/%s", fs.DefaultFsPath, T.s.Id.Hex(), T.scriptName)
	data, err := T.fsSvc.GetFs().GetFile(remotePath)
	require.Nil(t, err)
	require.Equal(t, T.script, string(data))
}

func TestFsService_SyncToWorkspace(t *testing.T) {
	var err error
	T.Setup(t)

	// save file to local
	require.Nil(t, err)
	err = T.masterFsSvc.GetFsService().Save(T.scriptName, []byte(T.script))

	// sync to fs
	err = T.workerSyncSvc.SyncToWorkspace(T.s.Id)
	require.Nil(t, err)

	// validate
	filePath := path.Join(T.workerFsSvc.GetWorkspacePath(), T.scriptName)
	data, err := ioutil.ReadFile(filePath)
	require.Nil(t, err)
	require.Equal(t, T.script, string(data))
}
