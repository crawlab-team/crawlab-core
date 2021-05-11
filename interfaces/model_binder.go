package interfaces

type modelBinder interface {
	Bind() (res Model, err error)
	MustBind() (res Model)
	Process(d Model) (res Model, err error)
}

type ModelBinder interface {
	modelBinder
	AssignFields(d Model, fieldIds ...ModelId) (res Model, err error)
}
