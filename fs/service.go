package fs

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/node/config"
	cfs "github.com/crawlab-team/crawlab-fs"
	vcs "github.com/crawlab-team/crawlab-vcs"
	"github.com/crawlab-team/go-trace"
	"github.com/ztrue/tracerr"
	"go.uber.org/dig"
	"os"
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

	// internals
	w *vcs.GitClient // workspace git client (only for master)
}

func (svc *Service) List(path string, opts ...interfaces.ServiceCrudOption) (files []interfaces.FsFileInfo, err error) {
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

func (svc *Service) GetFile(path string, opts ...interfaces.ServiceCrudOption) (data []byte, err error) {
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

func (svc *Service) Save(path string, data []byte, opts ...interfaces.ServiceCrudOption) (err error) {
	// forbidden if not master
	if !svc.nodeCfgSvc.IsMaster() {
		return trace.TraceError(errors.ErrorFsForbidden)
	}

	// save fs
	if err := svc.saveFs(path, data, opts...); err != nil {
		return err
	}

	// sync to workspace
	return svc.SyncToWorkspace()
}

func (svc *Service) Rename(path, newPath string, opts ...interfaces.ServiceCrudOption) (err error) {
	// TODO: implement rename directory
	// forbidden if not master
	if !svc.nodeCfgSvc.IsMaster() {
		return trace.TraceError(errors.ErrorFsForbidden)
	}

	// rename fs
	if err := svc.renameFs(path, newPath, opts...); err != nil {
		return err
	}

	// sync to workspace
	return svc.SyncToWorkspace()
}

func (svc *Service) Delete(path string, opts ...interfaces.ServiceCrudOption) (err error) {
	// forbidden if not master
	if !svc.nodeCfgSvc.IsMaster() {
		return trace.TraceError(errors.ErrorFsForbidden)
	}

	// delete fs
	if err := svc.deleteFs(path, opts...); err != nil {
		return err
	}

	// sync to workspace
	return svc.SyncToWorkspace()
}

func (svc *Service) Copy(path, newPath string, opts ...interfaces.ServiceCrudOption) (err error) {
	// forbidden if not master
	if !svc.nodeCfgSvc.IsMaster() {
		return tracerr.Wrap(errors.ErrorFsForbidden)
	}

	if err := svc.copyFs(path, newPath, opts...); err != nil {
		return err
	}

	return nil
}

func (svc *Service) Commit(msg string) (err error) {
	// forbidden if not master
	if !svc.nodeCfgSvc.IsMaster() {
		return trace.TraceError(errors.ErrorFsForbidden)
	}

	// TODO: check whether need to sync remote fs to local (based on md5?)
	//if err := svc.fs.SyncRemoteToLocal(svc.fsPath, c.GetPath()); err != nil {
	//	return err
	//}

	// commit
	if err := svc.w.CommitAll(msg); err != nil {
		return trace.TraceError(err)
	}

	// push to repo
	if err := svc.w.Push(vcs.WithRemoteNamePush(vcs.GitDefaultRemoteName)); err != nil {
		return trace.TraceError(err)
	}

	return nil
}

// SyncToFs sync from repo/workspace to fs
func (svc *Service) SyncToFs(opts ...interfaces.ServiceCrudOption) (err error) {
	// forbidden if not master
	if !svc.nodeCfgSvc.IsMaster() {
		return trace.TraceError(errors.ErrorFsForbidden)
	}

	// apply options
	o := svc.newCrudOptions()
	for _, opt := range opts {
		opt(o)
	}

	// sync from repo to workspace
	if !o.OnlyFromWorkspace {
		if err := svc.syncFromRepoToWorkspace(); err != nil {
			return err
		}
	}

	// sync from workspace to fs
	if err := svc.fs.SyncLocalToRemote(svc.GetWorkspacePath(), svc.fsPath); err != nil {
		return err
	}

	return nil
}

// SyncToWorkspace sync from fs to workspace
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

func (svc *Service) saveFs(path string, data []byte, opts ...interfaces.ServiceCrudOption) (err error) {
	// apply options
	o := svc.newCrudOptions()
	for _, opt := range opts {
		opt(o)
	}

	// remote path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	remotePath := fmt.Sprintf("%s%s", svc.fsPath, path)
	if o.IsAbsolute {
		remotePath = path
	}

	return svc.fs.UpdateFile(remotePath, data)
}

func (svc *Service) renameFs(path, newPath string, opts ...interfaces.ServiceCrudOption) (err error) {
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
		return trace.TraceError(errors.ErrorFsAlreadyExists)
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

func (svc *Service) deleteFs(path string, opts ...interfaces.ServiceCrudOption) (err error) {
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

func (svc *Service) copyFs(path, newPath string, opts ...interfaces.ServiceCrudOption) (err error) {
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
			if err := svc.copyFs(dirPath, dirPathNew, opts...); err != nil {
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

func (svc *Service) syncFromRepoToWorkspace() (err error) {
	// TODO: specify remote
	return svc.w.Reset()
}

func (svc *Service) getRemotePath(path string, o *interfaces.ServiceCrudOptions) (remotePath string) {
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

func (svc *Service) newCrudOptions() (o *interfaces.ServiceCrudOptions) {
	return &interfaces.ServiceCrudOptions{
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

	// workspace directory
	if _, err := os.Stat(svc.GetWorkspacePath()); err != nil {
		if err := os.MkdirAll(svc.GetWorkspacePath(), os.FileMode(0766)); err != nil {
			return nil, trace.TraceError(err)
		}
	}

	// remote repo and local workspace git client
	// TODO: external repo
	if svc.repoPath != "" {
		// remote repo
		if !vcs.IsGitRepoExists(svc.GetRepoPath()) {
			if err := vcs.CreateBareGitRepo(svc.GetRepoPath()); err != nil {
				return nil, err
			}
		}

		// local workspace git client
		// TODO: auth
		gitOpts := []vcs.GitOption{
			vcs.WithPath(svc.GetWorkspacePath()),
			vcs.WithRemoteUrl(svc.GetRepoPath()),
		}
		svc.w, err = vcs.NewGitClient(gitOpts...)
		if err != nil {
			return nil, err
		}
	}

	return svc, nil
}
