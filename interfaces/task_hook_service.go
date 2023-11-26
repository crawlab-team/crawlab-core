package interfaces

type TaskHookService interface {
	PreActions() (err error)
	PostActions() (err error)
}
