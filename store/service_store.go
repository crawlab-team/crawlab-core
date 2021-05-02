package store

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"sync"
)

type ServiceStore struct {
	m sync.Map
}

func (s *ServiceStore) GetMap() (m *sync.Map) {
	return &s.m
}

func (s *ServiceStore) Set(key interface{}, value interface{}) (err error) {
	if value == nil {
		return errors.ErrorStoreEmptyValue
	}
	s.m.Store(key, value)
	return nil
}

func (s *ServiceStore) Get(key interface{}) (res interface{}, err error) {
	if key == "" {
		return s.GetDefault()
	}
	res, ok := s.m.Load(key)
	if !ok {
		return nil, errors.ErrorStoreNotExists
	}
	return res, nil
}

func (s *ServiceStore) SetByString(key string, value interface{}) (err error) {
	return s.Set(key, value)
}

func (s *ServiceStore) GetByString(key string) (res interface{}, err error) {
	return s.Get(key)
}

func (s *ServiceStore) SetByInt(key int, value interface{}) (err error) {
	return s.Set(key, value)
}

func (s *ServiceStore) GetByInt(key int) (res interface{}, err error) {
	return s.Get(key)
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

func NewServiceStore() (s *ServiceStore) {
	return &ServiceStore{
		m: sync.Map{},
	}
}
