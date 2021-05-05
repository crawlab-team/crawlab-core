package interfaces

import "time"

type ModelNodeDelegate interface {
	ModelDelegate
	UpdateStatus(active bool, activeTs time.Time, status string) (err error)
}
