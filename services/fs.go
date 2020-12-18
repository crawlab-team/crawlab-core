package services

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/errors"
	cfs "github.com/crawlab-team/crawlab-fs"
	"strings"
)

type FileSystemServiceInterface interface {
	GetFile(path string) (data []byte, err error)
	Save(path, content string) (err error)
	Rename(path, newPath string) (err error)
	Delete(path string) (err error)
}

type FileSystemServiceOptions struct {
	BasePath string
}

func NewFileSystemService(options *FileSystemServiceOptions) (s *FileSystemService, err error) {
	if options == nil {
		options = &FileSystemServiceOptions{
			BasePath: "/fs",
		}
	}
	m, err := cfs.NewSeaweedFSManager()
	if err != nil {
		return s, err
	}
	basePath := options.BasePath
	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}
	s = &FileSystemService{
		m:        m,
		basePath: basePath,
	}
	return
}

type FileSystemService struct {
	m        *cfs.SeaweedFSManager
	basePath string
}

func (f *FileSystemService) GetFile(path string) (data []byte, err error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	remotePath := fmt.Sprintf("%s%s", f.basePath, path)
	return f.m.GetFile(remotePath)
}

func (f *FileSystemService) Save(path string, data []byte) (err error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	remotePath := fmt.Sprintf("%s%s", f.basePath, path)
	return f.m.UpdateFile(remotePath, data)
}

func (f *FileSystemService) Rename(path, newPath string) (err error) {
	// normalize paths
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasPrefix(newPath, "/") {
		newPath = "/" + newPath
	}

	// remote paths
	remotePath := fmt.Sprintf("%s%s", f.basePath, path)
	newRemotePath := fmt.Sprintf("%s%s", f.basePath, newPath)

	// error if new remote path exists
	ok, err := f.m.Exists(newRemotePath)
	if err != nil {
		return err
	}
	if ok {
		return errors.ErrAlreadyExists
	}

	// get original file data
	data, err := f.m.GetFile(remotePath)
	if err != nil {
		return err
	}

	// save data to new file
	if err := f.m.UpdateFile(newRemotePath, data); err != nil {
		return err
	}

	// delete original file
	if err := f.m.DeleteFile(remotePath); err != nil {
		return err
	}

	return nil
}

func (f *FileSystemService) Delete(path string) (err error) {
	// normalize paths
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// delete remote file
	remotePath := fmt.Sprintf("%s%s", f.basePath, path)
	if err := f.m.DeleteFile(remotePath); err != nil {
		return err
	}
	return nil
}
