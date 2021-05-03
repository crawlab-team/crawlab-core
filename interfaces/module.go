package interfaces

type ModuleId int

const (
	ModuleIdNode = iota
	ModuleIdModels
	ModuleIdGrpc
)

type Module interface {
	Init() error
	Start()
	Wait()
	Stop()
}
