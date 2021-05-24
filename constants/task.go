package constants

const (
	TaskStatusPending   = "pending"
	TaskStatusAssigned  = "assigned"
	TaskStatusRunning   = "running"
	TaskStatusFinished  = "finished"
	TaskStatusError     = "error"
	TaskStatusCancelled = "cancelled"
	TaskStatusAbnormal  = "abnormal"
)

const (
	TaskFinish = "finish"
	TaskCancel = "cancel"
)

const (
	RunTypeAllNodes      = "all-nodes"
	RunTypeRandom        = "random"
	RunTypeSelectedNodes = "selected-nodes"
)

const (
	TaskTypeSpider = "spider"
	TaskTypeSystem = "system"
)

type TaskSignal int

const (
	TaskSignalFinish TaskSignal = iota
	TaskSignalCancel
	TaskSignalError
	TaskSignalLost
)
