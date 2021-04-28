package node

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/go-trace"
	"sync"
)

var ServiceMap sync.Map

func GetServiceByKey(key string) (svc *Service, err error) {
	if key == "" {
		return GetDefaultService()
	}

	res, ok := ServiceMap.Load(key)
	if !ok {
		return nil, trace.TraceError(errors.ErrorNodeServiceNotExists)
	}
	svc, ok = res.(*Service)
	if !ok {
		return nil, trace.TraceError(errors.ErrorNodeInvalidType)
	}
	return svc, nil
}

func GetDefaultService() (svc *Service, err error) {
	return NewService(nil)
}
