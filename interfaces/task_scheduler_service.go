package interfaces

type TaskSchedulerService interface {
	TaskBaseService
	// Fetch continuously dequeue task from task queue and Assign to corresponding nodes to run tasks
	Fetch()
	// Assign tasks via sending grpc stream message to handler(s)
	Assign()
}
