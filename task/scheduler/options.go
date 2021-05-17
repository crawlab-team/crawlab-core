package scheduler

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
)

type Option func(svc interfaces.TaskSchedulerService)

func WithConfigPath(path string) Option {
	return func(svc interfaces.TaskSchedulerService) {
		svc.SetConfigPath(path)
	}
}
