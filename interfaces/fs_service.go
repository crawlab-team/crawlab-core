package interfaces

import (
	cfs "github.com/crawlab-team/crawlab-fs"
)

type FsService interface {
	WithConfigPath
	List(path string, opts ...FsServiceCrudOption) (files []FsFileInfo, err error)
	GetFile(path string, opts ...FsServiceCrudOption) (data []byte, err error)
	Save(path string, data []byte, opts ...FsServiceCrudOption) (err error)
	Rename(path, newPath string, opts ...FsServiceCrudOption) (err error)
	Delete(path string, opts ...FsServiceCrudOption) (err error)
	Copy(path, newPath string, opts ...FsServiceCrudOption) (err error)
	Commit(msg string) (err error)
	SyncToFs(opts ...FsServiceCrudOption) (err error)
	SyncToWorkspace() (err error)
	GetFsPath() (path string)
	SetFsPath(path string)
	GetWorkspacePath() (path string)
	SetWorkspacePath(path string)
	GetRepoPath() (path string)
	SetRepoPath(path string)
	GetFs() (fs cfs.Manager)
}
