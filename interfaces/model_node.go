package interfaces

import "time"

type Node interface {
	BaseModelWithTagsInterface
	GetKey() (key string)
	UpdateStatus(active bool, activeTs time.Time, status string) (err error)
}
