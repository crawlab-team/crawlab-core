package constants

const (
	TaskStatusPending   = "pending"
	TaskStatusRunning   = "running"
	TaskStatusFinished  = "finished"
	TaskStatusError     = "error"
	TaskStatusCancelled = "cancelled"
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

const (
	TaskListQueuePrefixPublic = "tasks:public"
	TaskListQueuePrefixNodes  = "tasks:nodes"
)
