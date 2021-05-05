package inject

import "github.com/goava/di"

func init() {
	di.SetTracer(&di.StdTracer{})

	var err error
	DefaultContainer, err = di.New()
	if err != nil {
		panic(err)
	}

	if err := Store.Set("", DefaultContainer); err != nil {
		panic(err)
	}
}

var DefaultContainer *di.Container
