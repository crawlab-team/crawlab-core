package fs

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/node/config"
	cfs "github.com/crawlab-team/crawlab-fs"
	vcs "github.com/crawlab-team/crawlab-vcs"
	"github.com/crawlab-team/go-trace"
	"github.com/google/uuid"
	"github.com/ztrue/tracerr"
	"go.uber.org/dig"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type Service struct {
	// settings
	cfgPath       string
	fsPath        string
	workspacePath string
	repoPath      string

	// dependencies
	nodeCfgSvc interfaces.NodeConfigService
	fs         cfs.Manager
	local      *vcs.GitClient
	repo       *vcs.GitClient
}

func (svc *Service) List(path string, opts ...interfaces.FsServiceCrudOption) (files []interfaces.FsFileInfo, err error) {
	// forbidden if not master
	if !svc.nodeCfgSvc.IsMaster() {
		return files, errors.ErrorFsForbidden
	}

	// apply options
	o := svc.newCrudOptions()
	for _, opt := range opts {
		opt(o)
	}

	// remote path
	remotePath := svc.getRemotePath(path, o)

	// list items of remote path recursively
	items, err := svc.fs.ListDir(remotePath, false)
	if err != nil {
		return files, err
	}
	for _, item := range items {
		itemPath := strings.Replace(item.FullPath, svc.fsPath, "", 1)
		if o.IsAbsolute {
			itemPath = item.FullPath
		}
		f := &entity.FsFileInfo{
			Name:      item.Name,
			Path:      itemPath,
			FullPath:  item.FullPath,
			Extension: item.Extension,
			Md5:       item.Md5,
			IsDir:     item.IsDir,
			FileSize:  item.FileSize,
		}
		if item.IsDir {
			relativePath := strings.Replace(item.FullPath, svc.fsPath, "", 1)
			if o.IsAbsolute {
				relativePath = item.FullPath
			}
			f.Children, err = svc.List(relativePath, opts...)
			if err != nil {
				return files, err
			}
		}
		files = append(files, f)
	}

	return files, nil
}

func (svc *Service) GetFile(path string, opts ...interfaces.FsServiceCrudOption) (data []byte, err error) {
	// forbidden if not master
	if !svc.nodeCfgSvc.IsMaster() {
		return data, trace.TraceError(errors.ErrorFsForbidden)
	}

	// apply options
	o := svc.newCrudOptions()
	for _, opt := range opts {
		opt(o)
	}

	// normalize remote path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	remotePath := fmt.Sprintf("%s%s", svc.fsPath, path)
	if o.IsAbsolute {
		remotePath = path
	}
	return svc.fs.GetFile(remotePath)
}

func (svc *Service) Save(path string, data []byte, opts ...interfaces.FsServiceCrudOption) (err error) {
	// forbidden if not master
	if !svc.nodeCfgSvc.IsMaster() {
		return trace.TraceError(errors.ErrorFsForbidden)
	}

	// apply options
	o := svc.newCrudOptions()
	for _, opt := range opts {
		opt(o)
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	remotePath := fmt.Sprintf("%s%s", svc.fsPath, path)
	if o.IsAbsolute {
		remotePath = path
	}
	return svc.fs.UpdateFile(remotePath, data)
}

func (svc *Service) Rename(path, newPath string, opts ...interfaces.FsServiceCrudOption) (err error) {
	// TODO: implement rename directory
	// forbidden if not master
	if !svc.nodeCfgSvc.IsMaster() {
		return trace.TraceError(errors.ErrorFsForbidden)
	}

	// apply options
	o := svc.newCrudOptions()
	for _, opt := range opts {
		opt(o)
	}

	// remote paths
	remotePath := svc.getRemotePath(path, o)
	newRemotePath := svc.getRemotePath(newPath, o)

	// error if new remote path exists
	ok, err := svc.fs.Exists(newRemotePath)
	if err != nil {
		return err
	}
	if ok {
		return constants.ErrAlreadyExists
	}

	// get original file data
	data, err := svc.fs.GetFile(remotePath)
	if err != nil {
		return err
	}

	// save data to new file
	if err := svc.fs.UpdateFile(newRemotePath, data); err != nil {
		return err
	}

	// delete original file
	if err := svc.fs.DeleteFile(remotePath); err != nil {
		return err
	}

	return nil
}

func (svc *Service) Delete(path string, opts ...interfaces.FsServiceCrudOption) (err error) {
	// forbidden if not master
	if !svc.nodeCfgSvc.IsMaster() {
		return trace.TraceError(errors.ErrorFsForbidden)
	}

	// apply options
	o := svc.newCrudOptions()
	for _, opt := range opts {
		opt(o)
	}

	// remote path
	remotePath := svc.getRemotePath(path, o)

	// delete remote file
	if err := svc.fs.DeleteFile(remotePath); err != nil {
		return err
	}
	return nil
}

func (svc *Service) Copy(path, newPath string, opts ...interfaces.FsServiceCrudOption) (err error) {
	// forbidden if not master
	if !svc.nodeCfgSvc.IsMaster() {
		return tracerr.Wrap(errors.ErrorFsForbidden)
	}

	// apply options
	o := svc.newCrudOptions()
	for _, opt := range opts {
		opt(o)
	}

	// normalize paths
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasPrefix(newPath, "/") {
		newPath = "/" + newPath
	}

	// iterate all files
	files, err := svc.List(path, opts...)
	if err != nil {
		return trace.TraceError(err)
	}
	for _, f := range files {
		if f.GetIsDir() {
			// directory
			dirPathNew := fmt.Sprintf("%s/%s", newPath, f.GetName())
			dirPath := f.GetPath()
			if o.IsAbsolute {
				dirPath = f.GetFullPath()
			}
			if err := svc.Copy(dirPath, dirPathNew, opts...); err != nil {
				return err
			}
		} else {
			// file
			filePath := f.GetPath()
			if o.IsAbsolute {
				filePath = f.GetFullPath()
			}
			data, err := svc.GetFile(filePath, opts...)
			if err != nil {
				return err
			}
			filePathNew := fmt.Sprintf("%s/%s", newPath, f.GetName())
			if err := svc.Save(filePathNew, data, opts...); err != nil {
				return err
			}
		}
	}

	return nil
}

func (svc *Service) Commit(msg string) (err error) {
	// forbidden if not master
	if !svc.nodeCfgSvc.IsMaster() {
		return trace.TraceError(errors.ErrorFsForbidden)
	}

	// local temp repo
	c, dirPath, err := svc.getLocalTempGitClient()
	if err != nil {
		return err
	}
	defer c.Dispose()

	// sync to local workspace from remote fs
	if err := svc.fs.SyncRemoteToLocal(svc.fsPath, dirPath); err != nil {
		return err
	}

	// commit
	if err := c.CommitAll(msg); err != nil {
		return trace.TraceError(err)
	}

	// push to repo
	if err := c.Push(vcs.GitDefaultRemoteName); err != nil {
		return trace.TraceError(err)
	}

	return nil
}

func (svc *Service) SyncToFs() (err error) {
	// forbidden if not master
	if !svc.nodeCfgSvc.IsMaster() {
		return trace.TraceError(errors.ErrorFsForbidden)
	}

	// local temp repo
	c, dirPath, err := svc.getLocalTempGitClient()
	if err != nil {
		return err
	}
	defer c.Dispose()

	// sync to fs
	if err := svc.fs.SyncLocalToRemote(dirPath, svc.fsPath); err != nil {
		return err
	}

	return nil
}

func (svc *Service) SyncToWorkspace() (err error) {
	// validate workspace path
	if svc.workspacePath == "" {
		return trace.TraceError(errors.ErrorFsEmptyWorkspacePath)
	}

	// create workspace directory if not exists
	if _, err := os.Stat(svc.workspacePath); err != nil {
		if err := os.MkdirAll(svc.workspacePath, os.ModePerm); err != nil {
			return trace.TraceError(err)
		}
	}

	// sync to local workspace from remote fs
	if err := svc.fs.SyncRemoteToLocal(svc.fsPath, svc.workspacePath); err != nil {
		return err
	}

	return nil
}

func (svc *Service) GetFsPath() (path string) {
	return svc.fsPath
}

func (svc *Service) SetFsPath(path string) {
	svc.fsPath = path
}

func (svc *Service) GetWorkspacePath() (path string) {
	return svc.workspacePath
}

func (svc *Service) SetWorkspacePath(path string) {
	svc.workspacePath = path
}

func (svc *Service) GetRepoPath() (path string) {
	return svc.repoPath
}

func (svc *Service) SetRepoPath(path string) {
	svc.repoPath = path
}

func (svc *Service) GetConfigPath() string {
	return svc.cfgPath
}

func (svc *Service) SetConfigPath(path string) {
	svc.cfgPath = path
}

func (svc *Service) GetFs() (fs cfs.Manager) {
	return svc.fs
}

func (svc *Service) getLocalTempGitClient() (c *vcs.GitClient, dirPath string, err error) {
	// create temp directory
	tmpPath := os.TempDir()
	dirPath = path.Join(tmpPath, uuid.New().String())
	if _, err := os.Stat(dirPath); err == nil {
		if err := os.RemoveAll(dirPath); err != nil {
			return c, dirPath, trace.TraceError(err)
		}
	}

	// absolute repo path
	matched, _ := regexp.MatchString("^http|^ssh", svc.repoPath)
	repoPath := svc.repoPath
	if !matched {
		repoPath, err = filepath.Abs(repoPath)
		if err != nil {
			return c, dirPath, trace.TraceError(err)
		}
	}

	// create git client
	c, err = vcs.NewGitClient(&vcs.GitOptions{
		Path:      dirPath,
		RemoteUrl: repoPath,
		IsBare:    false,
		IsMem:     false,
	})
	if err != nil {
		return c, dirPath, trace.TraceError(err)
	}

	return c, dirPath, nil
}

func (svc *Service) getRemotePath(path string, o *interfaces.FsServiceCrudOptions) (remotePath string) {
	// normalize path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if o.IsAbsolute {
		// absolute path
		return path
	} else {
		// relative path
		return fmt.Sprintf("%s%s", svc.fsPath, path)
	}
}

func (svc *Service) newCrudOptions() (opts *interfaces.FsServiceCrudOptions) {
	return &interfaces.FsServiceCrudOptions{
		IsAbsolute: false,
	}
}

func NewFsService(opts ...Option) (svc2 interfaces.FsService, err error) {
	// service
	svc := &Service{
		cfgPath:       config.DefaultConfigPath,
		fsPath:        DefaultFsPath,
		workspacePath: DefaultWorkspacePath,
		//repoPath:      DefaultRepoPath,
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// normalize fs base path
	if !strings.HasPrefix(svc.fsPath, "/") {
		svc.fsPath = "/" + svc.fsPath
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(config.ProvideConfigService(svc.cfgPath)); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Provide(cfs.NewSeaweedFsManager); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(nodeCfgSvc interfaces.NodeConfigService, fs cfs.Manager) {
		svc.nodeCfgSvc = nodeCfgSvc
		svc.fs = fs
	}); err != nil {
		return nil, trace.TraceError(err)
	}

	// local and remote repos
	if svc.repoPath != "" {
		// remote repo
		svc.repo, err = vcs.NewGitClient(&vcs.GitOptions{
			Path:   svc.repoPath,
			IsBare: true,
			IsMem:  false,
		})
		if err != nil {
			return nil, trace.TraceError(err)
		}

		// local temp repo (mem)
		svc.local, err = vcs.NewGitClient(&vcs.GitOptions{
			Path:      svc.repoPath,
			RemoteUrl: svc.repoPath,
			IsBare:    false,
			IsMem:     true,
		})
		if err != nil {
			return nil, trace.TraceError(err)
		}
	}

	return svc, nil
}
