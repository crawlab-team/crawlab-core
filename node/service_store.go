package node

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"sync"
)

type ServiceStore struct {
	m sync.Map
}

func (s *ServiceStore) Set(key string, value interface{}) (err error) {
	if value == nil {
		return errors.ErrorStoreEmptyValue
	}
	s.m.Store(key, value)
	return nil
}

func (s *ServiceStore) MustSet(key string, value interface{}) {
	if err := s.Set(key, value); err != nil {
		panic(err)
	}
}

func (s *ServiceStore) Get(key string) (res interface{}, err error) {
	res, ok := s.m.Load(key)
	if !ok {
		return nil, errors.ErrorStoreNotExists
	}
	return res, nil
}

func (s *ServiceStore) MustGet(key string) (res interface{}) {
	res, err := s.Get(key)
	if err != nil {
		panic(err)
	}
	return res
}

func (s *ServiceStore) GetDefault() (res interface{}, err error) {
	ok := false
	s.m.Range(func(key, value interface{}) bool {
		res = value
		ok = true
		return false
	})
	if !ok {
		return nil, errors.ErrorStoreNotExists
	}
	return res, nil
}

func (s *ServiceStore) MustGetDefault() (res interface{}) {
	res, err := s.GetDefault()
	if err != nil {
		panic(err)
	}
	return res
}

func NewServiceStore() (s *ServiceStore) {
	return &ServiceStore{
		m: sync.Map{},
	}
}
