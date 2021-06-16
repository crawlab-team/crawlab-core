package interfaces

import (
	"time"
)

type ScheduleService interface {
	WithConfigPath
	Module
	GetLocation() (loc *time.Location)
	SetLocation(loc *time.Location)
	GetDelay() (delay bool)
	SetDelay(delay bool)
	GetSkip() (skip bool)
	SetSkip(skip bool)
	GetUpdateInterval() (interval time.Duration)
	SetUpdateInterval(interval time.Duration)
	Enable(s Schedule) (err error)
	Disable(s Schedule) (err error)
	Update()
}
