package store

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
)

type GrpcServiceStoreType struct {
	*ServiceStore
}

func (s *GrpcServiceStoreType) GetGrpcService(key interface{}) (svc interfaces.GrpcService, err error) {
	res, err := s.Get(key)
	if err != nil {
		return nil, err
	}
	svc, ok := res.(interfaces.GrpcService)
	if !ok {
		return nil, errors.ErrorStoreInvalidType
	}
	return svc, nil
}

func (s *GrpcServiceStoreType) MustGetGrpcService(key interface{}) (svc interfaces.GrpcService) {
	svc, err := s.GetGrpcService(key)
	if err != nil {
		panic(err)
	}
	return svc
}

func NewGrpcServiceStore() (s *GrpcServiceStoreType) {
	return &GrpcServiceStoreType{NewServiceStore()}
}

var GrpcServiceStore *GrpcServiceStoreType
