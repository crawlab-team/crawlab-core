package interfaces

import "time"

type TaskSchedulerService interface {
	TaskBaseService
	// Fetch dequeue task from task queue
	Fetch() (t Task, err error)
	// Assign tasks via sending grpc stream message to handler(s)
	Assign()
	// SetFetchInterval set the interval or duration between two adjacent fetches
	SetFetchInterval(interval time.Duration)
	// GetTaskChannel internal channel of task
	GetTaskChannel() (ch chan []Task)
}
