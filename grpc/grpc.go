package grpc

func InitServices() {
	NodeService = nodeServer{}
	TaskService = taskService{}
}
