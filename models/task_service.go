package models

type taskService struct {
	*Service
}

func NewTaskService() (svc *taskService) {
	return &taskService{NewService(ModelIdTask)}
}

var TaskService = NewTaskService()
