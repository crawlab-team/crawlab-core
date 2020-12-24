package services

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	cfs "github.com/crawlab-team/crawlab-fs"
	vcs "github.com/crawlab-team/crawlab-vcs"
	"github.com/google/uuid"
	"github.com/linxGnu/goseaweedfs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type FileSystemServiceInterface interface {
	// CRUD actions on fs
	List(path string) (files []goseaweedfs.FileInfo, err error)
	GetFile(path string) (data []byte, err error)
	Save(path string, data []byte) (err error)
	Rename(path, newPath string) (err error)
	Delete(path string) (err error)

	// commit all files from fs and push to git repo
	Commit(msg string) (err error)

	// pull from git repo and sync to fs
	SyncToFs() (err error)

	// sync to local workspace from fs
	SyncToWorkspace() (err error)
}

type FileSystemServiceOptions struct {
	IsMaster      bool
	FsPath        string
	WorkspacePath string
	RepoPath      string
}

func NewFileSystemService(options *FileSystemServiceOptions) (s *FileSystemService, err error) {
	// options
	if options == nil {
		options = &FileSystemServiceOptions{
			IsMaster:      false,
			WorkspacePath: "/tmp/test_workspace",
			RepoPath:      "/repo",
			FsPath:        "/fs",
		}
	}

	// normalize fs base path
	fsPath := options.FsPath
	if !strings.HasPrefix(fsPath, "/") {
		fsPath = "/" + fsPath
	}

	// file system
	fs, err := cfs.NewSeaweedFSManager()
	if err != nil {
		return s, err
	}

	// local and remote repos
	var repo *vcs.GitClient
	var local *vcs.GitClient
	if options.RepoPath != "" {
		// remote repo
		repo, err = vcs.NewGitClient(&vcs.GitOptions{
			Path:   options.RepoPath,
			IsBare: true,
			IsMem:  false,
		})
		if err != nil {
			return s, err
		}

		// local temp repo (mem)
		local, err = vcs.NewGitClient(&vcs.GitOptions{
			Path:      options.RepoPath,
			RemoteUrl: options.RepoPath,
			IsBare:    false,
			IsMem:     true,
		})
		if err != nil {
			return s, err
		}
	}

	// file system service
	s = &FileSystemService{
		fs:    fs,
		local: local,
		repo:  repo,
		opts:  options,
	}

	return s, nil
}

type FileSystemService struct {
	fs    *cfs.SeaweedFSManager
	local *vcs.GitClient
	repo  *vcs.GitClient
	opts  *FileSystemServiceOptions
}

func (s *FileSystemService) List(path string) (files []entity.FsFileInfo, err error) {
	// forbidden if not master
	if !s.opts.IsMaster {
		return files, constants.ErrForbidden
	}

	// normalize path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	remotePath := fmt.Sprintf("%s%s", s.opts.FsPath, path)

	// list items of remote path recursively
	items, err := s.fs.ListDir(remotePath, false)
	if err != nil {
		return files, err
	}
	for _, item := range items {
		info := entity.FsFileInfo{
			Name:      item.Name,
			Path:      item.FullPath,
			Extension: item.Extension,
			Md5:       item.Md5,
			IsDir:     item.IsDir,
			FileSize:  item.FileSize,
		}
		if item.IsDir {
			relativePath := strings.Replace(item.FullPath, s.opts.FsPath, "", 1)
			info.Children, err = s.List(relativePath)
			if err != nil {
				return files, err
			}
		}
		files = append(files, info)
	}

	return files, nil
}

func (s *FileSystemService) GetFile(path string) (data []byte, err error) {
	// forbidden if not master
	if !s.opts.IsMaster {
		return data, constants.ErrForbidden
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	remotePath := fmt.Sprintf("%s%s", s.opts.FsPath, path)
	return s.fs.GetFile(remotePath)
}

func (s *FileSystemService) Save(path string, data []byte) (err error) {
	// forbidden if not master
	if !s.opts.IsMaster {
		return constants.ErrForbidden
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	remotePath := fmt.Sprintf("%s%s", s.opts.FsPath, path)
	return s.fs.UpdateFile(remotePath, data)
}

func (s *FileSystemService) Rename(path, newPath string) (err error) {
	// forbidden if not master
	if !s.opts.IsMaster {
		return constants.ErrForbidden
	}

	// normalize paths
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if !strings.HasPrefix(newPath, "/") {
		newPath = "/" + newPath
	}

	// remote paths
	remotePath := fmt.Sprintf("%s%s", s.opts.FsPath, path)
	newRemotePath := fmt.Sprintf("%s%s", s.opts.FsPath, newPath)

	// error if new remote path exists
	ok, err := s.fs.Exists(newRemotePath)
	if err != nil {
		return err
	}
	if ok {
		return constants.ErrAlreadyExists
	}

	// get original file data
	data, err := s.fs.GetFile(remotePath)
	if err != nil {
		return err
	}

	// save data to new file
	if err := s.fs.UpdateFile(newRemotePath, data); err != nil {
		return err
	}

	// delete original file
	if err := s.fs.DeleteFile(remotePath); err != nil {
		return err
	}

	return nil
}

func (s *FileSystemService) Delete(path string) (err error) {
	// forbidden if not master
	if !s.opts.IsMaster {
		return constants.ErrForbidden
	}

	// normalize paths
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// delete remote file
	remotePath := fmt.Sprintf("%s%s", s.opts.FsPath, path)
	if err := s.fs.DeleteFile(remotePath); err != nil {
		return err
	}
	return nil
}

func (s *FileSystemService) Commit(msg string) (err error) {
	// forbidden if not master
	if !s.opts.IsMaster {
		return constants.ErrForbidden
	}

	// validate options
	if s.opts.RepoPath == "" || s.opts.FsPath == "" {
		return constants.ErrInvalidOptions
	}

	// local temp repo
	c, dirPath, err := s.GetLocalTempGitClient()
	if err != nil {
		return err
	}
	defer c.Dispose()

	// sync to local workspace from remote fs
	if err := s.fs.SyncRemoteToLocal(s.opts.FsPath, dirPath); err != nil {
		return err
	}

	// commit
	if err := c.CommitAll(msg); err != nil {
		return err
	}

	// push to repo
	if err := c.Push(vcs.GitDefaultRemoteName); err != nil {
		return err
	}

	return nil
}

func (s *FileSystemService) SyncToFs() (err error) {
	// forbidden if not master
	if !s.opts.IsMaster {
		return constants.ErrForbidden
	}

	// validate options
	if s.opts.RepoPath == "" || s.opts.FsPath == "" {
		return constants.ErrInvalidOptions
	}

	// local temp repo
	c, dirPath, err := s.GetLocalTempGitClient()
	if err != nil {
		return err
	}
	defer c.Dispose()

	// sync to fs
	if err := s.fs.SyncLocalToRemote(dirPath, s.opts.FsPath); err != nil {
		return err
	}

	return nil
}

func (s *FileSystemService) SyncToWorkspace() (err error) {
	// validate options
	if s.opts.WorkspacePath == "" {
		return constants.ErrInvalidOptions
	}

	// create workspace directory if not exists
	if _, err := os.Stat(s.opts.WorkspacePath); err != nil {
		if err := os.MkdirAll(s.opts.WorkspacePath, os.ModePerm); err != nil {
			return err
		}
	}

	// sync to local workspace from remote fs
	if err := s.fs.SyncRemoteToLocal(s.opts.FsPath, s.opts.WorkspacePath); err != nil {
		return err
	}

	return nil
}

func (s *FileSystemService) GetLocalTempGitClient() (c *vcs.GitClient, dirPath string, err error) {
	// validate options
	if s.opts.RepoPath == "" {
		return c, dirPath, constants.ErrInvalidOptions
	}

	// create temp directory
	tmpPath := os.TempDir()
	dirPath = path.Join(tmpPath, uuid.New().String())
	if _, err := os.Stat(dirPath); err == nil {
		if err := os.RemoveAll(dirPath); err != nil {
			return c, dirPath, err
		}
	}

	// absolute repo path
	matched, _ := regexp.MatchString("^http|^ssh", s.opts.RepoPath)
	repoPath := s.opts.RepoPath
	if !matched {
		repoPath, err = filepath.Abs(repoPath)
		if err != nil {
			return c, dirPath, err
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
		return c, dirPath, err
	}

	return c, dirPath, nil
}
