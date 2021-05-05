package service

import (
	"github.com/crawlab-team/crawlab-core/color"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"go.uber.org/dig"
)

type Service struct {
	env      string
	colorSvc interfaces.ColorService
}

func (svc *Service) NewBaseService(id interfaces.ModelId) (svc2 interfaces.ModelBaseService) {
	return NewBaseService(id)
}

func NewService(opts ...Option) (svc2 ModelService, err error) {
	// service
	svc := &Service{}

	// options
	for _, opt := range opts {
		opt(svc)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(color.NewService); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(colorSvc interfaces.ColorService) {
		svc.colorSvc = colorSvc
	}); err != nil {
		return nil, err
	}

	return svc, nil
}
