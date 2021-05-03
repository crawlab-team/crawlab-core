package models

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/store"
	"github.com/crawlab-team/go-trace"
)

func GetRootService() (svc *Service, err error) {
	if RootModelService == nil {
		return nil, trace.TraceError(errors.ErrorModelNotFound)
	}
	return RootModelService, nil
}

func MustGetRootService() (svc *Service) {
	svc, err := GetRootService()
	if err != nil {
		panic(err)
	}
	return svc
}

func GetService(id interfaces.ModelId) (svc interfaces.ModelService, err error) {
	svc, err = store.ModelServiceStore.GetModelService(id)
	if err != nil {
		return nil, err
	}
	return svc, nil
}

func MustGetService(id interfaces.ModelId) (svc interfaces.ModelService) {
	svc, err := GetService(id)
	if err != nil {
		panic(err)
	}
	return svc
}
