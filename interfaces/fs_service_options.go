package interfaces

type FsServiceCrudOptions struct {
	IsAbsolute        bool // whether the path is absolute
	OnlyFromWorkspace bool // true if only sync from workspace
}

type FsServiceCrudOption func(o *FsServiceCrudOptions)

func WithOnlyFromWorkspace() FsServiceCrudOption {
	return func(o *FsServiceCrudOptions) {
		o.OnlyFromWorkspace = true
	}
}
