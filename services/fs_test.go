package services

import (
	cerr "github.com/crawlab-team/crawlab-core/errors"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func setup() (err error) {
	return cleanup()
}

func cleanup() (err error) {
	s, err := NewFileSystemService(&FileSystemServiceOptions{BasePath: "/test"})
	if err != nil {
		return err
	}
	ok, err := s.m.Exists("/test")
	if err != nil {
		return err
	}
	if ok {
		if err := s.m.DeleteDir("/test"); err != nil {
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
	// setup
	err := setup()
	require.Nil(t, err)

	s, err := NewFileSystemService(&FileSystemServiceOptions{BasePath: "/test"})
	require.Nil(t, err)

	require.NotNil(t, s)
	require.Equal(t, "/test", s.basePath)

	// cleanup
	err = cleanup()
	require.Nil(t, err)
}

func TestFileSystemService_Save(t *testing.T) {
	// setup
	err := setup()
	require.Nil(t, err)

	s, err := NewFileSystemService(&FileSystemServiceOptions{BasePath: "/test"})
	require.Nil(t, err)

	// save new file to remote
	content := "it works"
	err = s.Save("test_file.txt", []byte(content))
	require.Nil(t, err)

	// get file
	data, err := s.GetFile("test_file.txt")
	require.Nil(t, err)
	require.Equal(t, content, string(data))

	// cleanup
	err = cleanup()
	require.Nil(t, err)
}

func TestFileSystemService_Rename(t *testing.T) {
	// setup
	err := setup()
	require.Nil(t, err)

	s, err := NewFileSystemService(&FileSystemServiceOptions{BasePath: "/test"})
	require.Nil(t, err)

	// save new file to remote
	content := "it works"
	err = s.Save("test_file.txt", []byte(content))
	require.Nil(t, err)
	ok, err := s.m.Exists("/test/test_file.txt")
	require.Nil(t, err)
	require.True(t, ok)

	// rename file
	err = s.Rename("test_file.txt", "test_file2.txt")
	require.Nil(t, err)
	ok, err = s.m.Exists("/test/test_file.txt")
	require.Nil(t, err)
	require.False(t, ok)
	ok, err = s.m.Exists("/test/test_file2.txt")
	require.Nil(t, err)
	require.True(t, ok)

	// rename to existing
	err = s.Save("test_file.txt", []byte(content))
	require.Nil(t, err)
	err = s.Rename("test_file.txt", "test_file2.txt")
	require.Equal(t, cerr.ErrAlreadyExists, err)

	// cleanup
	err = cleanup()
	require.Nil(t, err)
}

func TestFileSystemService_Delete(t *testing.T) {
	// setup
	err := setup()
	require.Nil(t, err)

	s, err := NewFileSystemService(&FileSystemServiceOptions{BasePath: "/test"})
	require.Nil(t, err)

	// save new file to remote
	content := "it works"
	err = s.Save("test_file.txt", []byte(content))
	require.Nil(t, err)

	// delete remote file
	err = s.Delete("test_file.txt")
	require.Nil(t, err)
	ok, err := s.m.Exists("/test/test_file.txt")
	require.Nil(t, err)
	require.False(t, ok)

	// cleanup
	err = cleanup()
	require.Nil(t, err)
}
