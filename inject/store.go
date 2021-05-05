package inject

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/goava/di"
	"sync"
)

type store struct {
	m sync.Map
}

func (s *store) GetMap() (m *sync.Map) {
	return &s.m
}

func (s *store) Set(key string, value *di.Container) (err error) {
	if value == nil {
		return errors.ErrorInjectEmptyValue
	}
	s.m.Store(key, value)
	return nil
}

func (s *store) Get(key string) (c *di.Container, err error) {
	res, ok := s.m.Load(key)
	if !ok {
		c, err = di.New()
		if err != nil {
			return nil, err
		}
		s.m.Store(key, c)
		return c, nil
	}
	c, ok = res.(*di.Container)
	if !ok {
		return nil, errors.ErrorInjectInvalidType
	}
	return c, nil
}

func NewStore() (s *store) {
	return &store{
		m: sync.Map{},
	}
}

var Store = NewStore()
