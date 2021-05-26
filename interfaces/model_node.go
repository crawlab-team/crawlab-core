package interfaces

import "time"

type Node interface {
	ModelWithTags
	GetKey() (key string)
	SetActive(active bool)
	SetActiveTs(activeTs time.Time)
	SetStatus(status string)
	SetAvailable(available bool)
}
