package interfaces

type ServiceCrudOptions struct {
	IsAbsolute        bool // whether the path is absolute
	OnlyFromWorkspace bool // true if only sync from workspace
}

type ServiceCrudOption func(o *ServiceCrudOptions)

func WithOnlyFromWorkspace() ServiceCrudOption {
	return func(o *ServiceCrudOptions) {
		o.OnlyFromWorkspace = true
	}
}
