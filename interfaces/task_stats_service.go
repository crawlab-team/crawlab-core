package interfaces

type TaskStatsService interface {
	TaskBaseService
	InsertData(t Task, records ...interface{}) (err error)
	InsertLogs(t Task, logs ...string) (err error)
}
