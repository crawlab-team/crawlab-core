package interfaces

type Spider interface {
	Model
	GetCmd() (cmd string)
	GetType() (ty string)
}
