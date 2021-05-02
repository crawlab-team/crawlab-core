package store

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
)

type NodeServiceStoreType struct {
	*ServiceStore
}

func (s *NodeServiceStoreType) GetNodeService(key interface{}) (svc interfaces.NodeService, err error) {
	res, err := s.Get(key)
	if err != nil {
		return nil, err
	}
	svc, ok := res.(interfaces.NodeService)
	if !ok {
		return nil, errors.ErrorStoreInvalidType
	}
	return svc, nil
}

func (s *NodeServiceStoreType) MustGetNodeService(key interface{}) (svc interfaces.NodeService) {
	svc, err := s.GetNodeService(key)
	if err != nil {
		panic(err)
	}
	return svc
}

func NewNodeServiceStore() (s *NodeServiceStoreType) {
	return &NodeServiceStoreType{NewServiceStore()}
}

var NodeServiceStore *NodeServiceStoreType
