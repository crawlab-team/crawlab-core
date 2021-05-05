package interfaces

type ModelBinder interface {
	Bind() (res Model, err error)
	MustBind() (res Model)
	Process(d Model, fieldIds ...ModelId) (res Model, err error)
	AssignFields(d Model, fieldIds ...ModelId) (res Model, err error)
}
