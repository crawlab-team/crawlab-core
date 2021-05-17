package client

import (
	"github.com/crawlab-team/crawlab-core/grpc/client"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/node/config"
	"go.uber.org/dig"
)

type ServiceDelegate struct {
	// settings variables
	cfgPath string

	// internals
	c interfaces.GrpcClient
}

func (d *ServiceDelegate) GetConfigPath() string {
	return d.cfgPath
}

func (d *ServiceDelegate) SetConfigPath(path string) {
	d.cfgPath = path
}

func (d *ServiceDelegate) NewBaseServiceDelegate(id interfaces.ModelId) (svc interfaces.GrpcClientModelBaseService, err error) {
	return NewBaseServiceDelegate(
		WithBaseServiceModelId(id),
		WithBaseServiceConfigPath(d.cfgPath),
	)
}

func NewServiceDelegate(opts ...ModelServiceDelegateOption) (svc2 interfaces.GrpcClientModelService, err error) {
	// service delegate
	svc := &ServiceDelegate{
		cfgPath: config.DefaultConfigPath,
	}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(client.ProvideGetClient(svc.cfgPath)); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(client interfaces.GrpcClient) {
		svc.c = client
	}); err != nil {
		return nil, err
	}

	return svc, nil
}

func ProvideServiceDelegate(path string) func() (svc interfaces.GrpcClientModelService, err error) {
	return func() (svc interfaces.GrpcClientModelService, err error) {
		return NewServiceDelegate(WithServiceConfigPath(path))
	}
}
