package interfaces

type FsServiceCrudOptions struct {
	IsAbsolute bool
}

type FsServiceCrudOption func(*FsServiceCrudOptions)
