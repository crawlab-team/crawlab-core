package interfaces

type TaskHookService interface {
	PreActions(Task, Spider, TaskFsService) (err error)
	PostActions(Task, Spider, TaskFsService) (err error)
}
