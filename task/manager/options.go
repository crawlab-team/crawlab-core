package manager

import "github.com/crawlab-team/crawlab-core/interfaces"

type Option func(svc interfaces.TaskManagerService)

func WithConfigPath(path string) Option {
	return func(svc interfaces.TaskManagerService) {
		svc.SetConfigPath(path)
	}
}
