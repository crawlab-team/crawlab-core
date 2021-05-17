package task

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"time"
)

type Option func(svc interfaces.TaskService)

func WithMaxRunners(maxRunners int) Option {
	return func(svc interfaces.TaskService) {
		svc.SetMaxRunners(maxRunners)
	}
}

func WithPollWaitDuration(duration time.Duration) Option {
	return func(svc interfaces.TaskService) {
		svc.SetPollWaitDuration(duration)
	}
}
