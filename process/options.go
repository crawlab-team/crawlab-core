package process

import "github.com/crawlab-team/crawlab-core/interfaces"

type DaemonOption func(d interfaces.ProcessDaemon)

func WithDaemonMaxErrors(maxErrors int) DaemonOption {
	return func(d interfaces.ProcessDaemon) {
		d.SetMaxErrors(maxErrors)
	}
}
