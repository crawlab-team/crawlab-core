package interfaces

import (
	vcs "github.com/crawlab-team/crawlab-vcs"
	"os"
)

type FsServiceV2 interface {
	List(path string) (files []os.FileInfo, err error)
	GetFile(path string) (data []byte, err error)
	GetFileInfo(path string) (file os.FileInfo, err error)
	Save(path string, data []byte) (err error)
	Rename(path, newPath string) (err error)
	Delete(path string) (err error)
	Copy(path, newPath string) (err error)
	GetGitClient() (c *vcs.GitClient)
}
