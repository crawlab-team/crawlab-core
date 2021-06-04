package plugin

import "github.com/crawlab-team/crawlab-core/interfaces"

type Option func(svc interfaces.PluginService)

func WithDirPath(path string) Option {
	return func(svc interfaces.PluginService) {
		svc.SetDirPath(path)
	}
}

func WithFsPathBase(path string) Option {
	return func(svc interfaces.PluginService) {
		svc.SetFsPathBase(path)
	}
}
