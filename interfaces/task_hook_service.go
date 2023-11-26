package interfaces

type TaskHookService interface {
	PreActions(Task, Spider, TaskFsService, TaskHandlerService) (err error)
	PostActions(Task, Spider, TaskFsService, TaskHandlerService) (err error)
}
