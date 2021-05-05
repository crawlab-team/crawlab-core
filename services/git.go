package services

import (
	"fmt"
	"github.com/crawlab-team/crawlab-core/constants"
	models2 "github.com/crawlab-team/crawlab-core/models/models"
	"github.com/crawlab-team/crawlab-vcs"
	"github.com/spf13/viper"
)

func NewGitService(spider models2.Spider) (s *GitService, err error) {
	// base paths
	remoteBasePath := viper.GetString("git.basePath.remote")
	localBasePath := viper.GetString("git.basePath.local")

	// paths
	remotePath := fmt.Sprintf("%s/%s", remoteBasePath, spider.Id.Hex())
	localPath := fmt.Sprintf("%s/%s", localBasePath, spider.Id.Hex())

	// remote git client
	remoteClient, err := vcs.NewGitClient(&vcs.GitOptions{
		Path:   remotePath,
		IsBare: true,
	})
	if err != nil {
		return s, err
	}

	// local git client
	localClient, err := vcs.NewGitClient(&vcs.GitOptions{
		Path:      localPath,
		RemoteUrl: remotePath,
		IsBare:    false,
		IsMem:     true,
	})
	if err != nil {
		return s, err
	}

	// git service
	s = &GitService{
		local:  localClient,
		remote: remoteClient,
	}

	return
}

type GitService struct {
	local  *vcs.GitClient
	remote *vcs.GitClient
}

func (s *GitService) LocalClient() (c *vcs.GitClient, err error) {
	if s.local == nil {
		return c, constants.ErrNotExists
	}
	c = s.local
	return c, nil
}

func (s *GitService) RemoteClient() (c *vcs.GitClient, err error) {
	if s.remote == nil {
		return c, constants.ErrNotExists
	}
	c = s.remote
	return c, nil
}

func (s *GitService) Pull(target interface{}) (err error) {
	return s.local.Pull()
}

func (s *GitService) Push(target interface{}) (err error) {
	if err := s.local.CommitAll("commit"); err != nil {
		return err
	}
	return s.local.Push()
}
