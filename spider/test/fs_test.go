package test

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/fs"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFsService_SyncToFs(t *testing.T) {
	var err error
	T.Setup(t)

	// add file
	err = T.masterFsSvc.GetFsService().Save("main.go", []byte(T.script))
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
