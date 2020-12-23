package services

import (
	"github.com/crawlab-team/crawlab-core/constants"
	vcs "github.com/crawlab-team/crawlab-vcs"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func setupFs() (err error) {
	return cleanupFs()
}

func cleanupFs() (err error) {
	s, err := NewFileSystemService(&FileSystemServiceOptions{FsPath: "/test", IsMaster: true})
	if err != nil {
		return err
	}
	ok, err := s.fs.Exists("/test")
	if err != nil {
		return err
	}
	if ok {
		if err := s.fs.DeleteDir("/test"); err != nil {
			return err
		}
	}
	if _, err := os.Stat("./tmp"); err == nil {
		if err := os.RemoveAll("./tmp"); err != nil {
			return err
		}
	}
	if err := os.MkdirAll("./tmp", os.ModePerm); err != nil {
		return err
	}
	return nil
}

func TestNewFileSystemService(t *testing.T) {
	// setupFs
	err := setupFs()
	require.Nil(t, err)

	s, err := NewFileSystemService(&FileSystemServiceOptions{FsPath: "/test", IsMaster: true})
	require.Nil(t, err)

	require.NotNil(t, s)
	require.Equal(t, "/test", s.opts.FsPath)

	// cleanupFs
	err = cleanupFs()
	require.Nil(t, err)
}

func TestFileSystemService_Save(t *testing.T) {
	// setupFs
	err := setupFs()
	require.Nil(t, err)

	s, err := NewFileSystemService(&FileSystemServiceOptions{FsPath: "/test", IsMaster: true})
	require.Nil(t, err)

	// save new file to remote
	content := "it works"
	err = s.Save("test_file.txt", []byte(content))
	require.Nil(t, err)

	// get file
	data, err := s.GetFile("test_file.txt")
	require.Nil(t, err)
	require.Equal(t, content, string(data))

	// cleanupFs
	err = cleanupFs()
	require.Nil(t, err)
}

func TestFileSystemService_Rename(t *testing.T) {
	// setupFs
	err := setupFs()
	require.Nil(t, err)

	s, err := NewFileSystemService(&FileSystemServiceOptions{FsPath: "/test", IsMaster: true})
	require.Nil(t, err)

	// save new file to remote
	content := "it works"
	err = s.Save("test_file.txt", []byte(content))
	require.Nil(t, err)
	ok, err := s.fs.Exists("/test/test_file.txt")
	require.Nil(t, err)
	require.True(t, ok)

	// rename file
	err = s.Rename("test_file.txt", "test_file2.txt")
	require.Nil(t, err)
	ok, err = s.fs.Exists("/test/test_file.txt")
	require.Nil(t, err)
	require.False(t, ok)
	ok, err = s.fs.Exists("/test/test_file2.txt")
	require.Nil(t, err)
	require.True(t, ok)

	// rename to existing
	err = s.Save("test_file.txt", []byte(content))
	require.Nil(t, err)
	err = s.Rename("test_file.txt", "test_file2.txt")
	require.Equal(t, constants.ErrAlreadyExists, err)

	// cleanupFs
	err = cleanupFs()
	require.Nil(t, err)
}

func TestFileSystemService_Delete(t *testing.T) {
	// setupFs
	err := setupFs()
	require.Nil(t, err)

	s, err := NewFileSystemService(&FileSystemServiceOptions{FsPath: "/test", IsMaster: true})
	require.Nil(t, err)

	// save new file to remote
	content := "it works"
	err = s.Save("test_file.txt", []byte(content))
	require.Nil(t, err)

	// delete remote file
	err = s.Delete("test_file.txt")
	require.Nil(t, err)
	ok, err := s.fs.Exists("/test/test_file.txt")
	require.Nil(t, err)
	require.False(t, ok)

	// cleanupFs
	err = cleanupFs()
	require.Nil(t, err)
}

func TestFileSystemService_Commit(t *testing.T) {
	// setupFs
	err := setupFs()
	require.Nil(t, err)

	// create a master fs service
	s, err := NewFileSystemService(&FileSystemServiceOptions{
		IsMaster: true,
		FsPath:   "/test",
		RepoPath: "./tmp/test_repo",
	})
	require.Nil(t, err)

	// save new file to remote
	content := "it works"
	err = s.Save("test_file.txt", []byte(content))
	require.Nil(t, err)

	// commit to repo
	err = s.Commit("test commit")
	require.Nil(t, err)

	// new git client from remote repo
	c, err := vcs.NewGitClient(&vcs.GitOptions{
		Path:      "./tmp/test_local",
		RemoteUrl: "./tmp/test_repo",
		IsBare:    false,
		IsMem:     false,
	})
	require.Nil(t, err)
	require.NotNil(t, c)
	require.FileExists(t, "./tmp/test_local/test_file.txt")
	data, err := ioutil.ReadFile("./tmp/test_local/test_file.txt")
	require.Nil(t, err)
	require.Equal(t, content, string(data))

	// cleanupFs
	err = cleanupFs()
	require.Nil(t, err)
}

func TestFileSystemService_SyncToFs(t *testing.T) {
	// setupFs
	err := setupFs()
	require.Nil(t, err)

	// create a master fs service
	s, err := NewFileSystemService(&FileSystemServiceOptions{
		IsMaster: true,
		FsPath:   "/test",
		RepoPath: "./tmp/test_repo",
	})
	require.Nil(t, err)

	// save new file to remote
	content := "it works"
	err = s.Save("test_file.txt", []byte(content))
	require.Nil(t, err)

	// commit to repo
	err = s.Commit("test commit")
	require.Nil(t, err)

	// edit the file
	content2 := "hello world"
	err = s.Save("test_file.txt", []byte(content2))
	require.Nil(t, err)

	// test file content
	data, err := s.GetFile("test_file.txt")
	require.Nil(t, err)
	require.Equal(t, content2, string(data))

	// sync to fs
	err = s.SyncToFs()
	require.Nil(t, err)

	// test file content
	data, err = s.GetFile("test_file.txt")
	require.Nil(t, err)
	require.Equal(t, content, string(data))

	// cleanupFs
	err = cleanupFs()
	require.Nil(t, err)
}

func TestFileSystemService_SyncToWorkspace(t *testing.T) {
	// setupFs
	err := setupFs()
	require.Nil(t, err)

	// create a master fs service
	s, err := NewFileSystemService(&FileSystemServiceOptions{
		IsMaster: true,
		FsPath:   "/test",
		RepoPath: "./tmp/test_repo",
	})
	require.Nil(t, err)

	// save new file to remote
	content := "it works"
	err = s.Save("test_file.txt", []byte(content))
	require.Nil(t, err)

	// create a worker fs service
	s2, err := NewFileSystemService(&FileSystemServiceOptions{
		IsMaster:      false,
		FsPath:        "/test",
		WorkspacePath: "./tmp/test_workspace",
	})
	require.Nil(t, err)

	// sync to workspace
	err = s2.SyncToWorkspace()
	require.Nil(t, err)
	require.FileExists(t, "./tmp/test_workspace/test_file.txt")
	data, err := ioutil.ReadFile("./tmp/test_workspace/test_file.txt")
	require.Nil(t, err)
	require.Equal(t, content, string(data))

	// cleanupFs
	err = cleanupFs()
	require.Nil(t, err)
}

func TestFileSystemService_WorkerFsService(t *testing.T) {
	// setupFs
	err := setupFs()
	require.Nil(t, err)

	// create a master fs service
	s, err := NewFileSystemService(&FileSystemServiceOptions{
		IsMaster: true,
		FsPath:   "/test",
		RepoPath: "./tmp/test_repo",
	})
	require.Nil(t, err)

	// save new file to remote
	content := "it works"
	err = s.Save("test_file.txt", []byte(content))
	require.Nil(t, err)

	// create a worker fs service
	s2, err := NewFileSystemService(&FileSystemServiceOptions{
		IsMaster:      false,
		FsPath:        "/test",
		WorkspacePath: "./tmp/test_workspace",
	})
	require.Nil(t, err)

	// test methods
	_, err = s2.GetFile("test_file.txt")
	require.Equal(t, constants.ErrForbidden, err)
	err = s2.Save("test_file.txt", []byte("it works"))
	require.Equal(t, constants.ErrForbidden, err)
	err = s2.Rename("test_file.txt", "new_test_file.txt")
	require.Equal(t, constants.ErrForbidden, err)
	err = s2.Delete("test_file.txt")
	require.Equal(t, constants.ErrForbidden, err)
	err = s2.Commit("test commit")
	require.Equal(t, constants.ErrForbidden, err)
	err = s2.SyncToFs()
	require.Equal(t, constants.ErrForbidden, err)
	err = s2.SyncToWorkspace()
	require.Nil(t, err)
	data, err := ioutil.ReadFile("./tmp/test_workspace/test_file.txt")
	require.Nil(t, err)
	require.Equal(t, content, string(data))

	// cleanupFs
	err = cleanupFs()
	require.Nil(t, err)
}
