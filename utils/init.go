package utils

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"sync"
)

var moduleInitializedMap = sync.Map{}

func InitModule(id interfaces.ModuleId, fn func() error) (err error) {
	res, ok := moduleInitializedMap.Load(id)
	if ok {
		return nil
	}
	initialized, ok := res.(bool)
	if !ok || initialized {
		moduleInitializedMap.Store(id, true)
		return nil
	}

	return fn()
}
