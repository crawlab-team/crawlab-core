package interfaces

import "time"

type Node interface {
	ModelWithTags
	GetName() (name string)
	GetKey() (key string)
	GetIsMaster() (ok bool)
	SetActive(active bool)
	SetActiveTs(activeTs time.Time)
	GetStatus() (status string)
	SetStatus(status string)
	GetEnabled() (enabled bool)
	SetEnabled(enabled bool)
	GetAvailableRunners() (runners int)
	SetAvailableRunners(runners int)
	GetMaxRunners() (runners int)
	SetMaxRunners(runners int)
	IncrementAvailableRunners()
	DecrementAvailableRunners()
}
