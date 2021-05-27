package interfaces

import "time"

type Node interface {
	ModelWithTags
	GetKey() (key string)
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
