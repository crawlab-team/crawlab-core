package interfaces

type TaskBaseService interface {
	WithConfigPath
	Module
	SaveTask(t Task, status string) (err error)
	IsStopped() (res bool)
}
