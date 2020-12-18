package services

import vcs "github.com/crawlab-team/crawlab-vcs"

type VcsServiceInterface interface {
	LocalClient() (c *vcs.Client, err error)
	RemoteClient() (c *vcs.Client, err error)
	Pull(target interface{}) (err error)
	Push(target interface{}) (err error)
}
