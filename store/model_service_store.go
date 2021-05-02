package store

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
)

type ModelServiceStoreType struct {
	*ServiceStore
}

func (s *ModelServiceStoreType) GetModelService(key interfaces.ModelId) (svc interfaces.ModelService, err error) {
	res, err := s.Get(key)
	if err != nil {
		return nil, err
	}
	svc, ok := res.(interfaces.ModelService)
	if !ok {
		return nil, errors.ErrorStoreInvalidType
	}
	return svc, nil
}

func (s *ModelServiceStoreType) MustGetModelService(key interfaces.ModelId) (svc interfaces.ModelService) {
	svc, err := s.GetModelService(key)
	if err != nil {
		panic(err)
	}
	return svc
}

func NewModelServiceStore() (s *ModelServiceStoreType) {
	return &ModelServiceStoreType{NewServiceStore()}
}

var ModelServiceStore *ModelServiceStoreType
