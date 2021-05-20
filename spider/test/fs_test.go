package test

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/fs"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestFsService_SyncToFs(t *testing.T) {
	var err error
	T.Setup(t)

	// write file
	filePath, _ := filepath.Abs(path.Join(T.masterFsSvc.GetWorkspacePath(), "main.go"))
	err = ioutil.WriteFile(filePath, []byte(T.script), os.ModePerm)
	require.Nil(t, err)

	// commit
	err = T.masterFsSvc.GetFsService().Commit("initial commit")
	require.Nil(t, err)

	// sync to fs
	err = T.masterSyncSvc.SyncToFs(T.s.Id)
	require.Nil(t, err)

	// validate
	remotePath := fmt.Sprintf("%s/%s/%s", fs.DefaultFsPath, T.s.Id.Hex(), "main.go")
	data, err := T.fsSvc.GetFs().GetFile(remotePath)
	require.Nil(t, err)
	require.Equal(t, T.script, string(data))
}
