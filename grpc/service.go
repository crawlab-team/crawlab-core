package grpc

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/go-trace"
	"sync"
)

type Service struct {
	server     *Server
	clientsMap sync.Map
	opts       *ServiceOptions
}

func (svc *Service) Init() (err error) {
	// start server
	if err := svc.server.Start(); err != nil {
		return err
	}
	return nil
}

func (svc *Service) GetClient(address Address) (client *Client, err error) {
	res, ok := svc.clientsMap.Load(address)
	if !ok {
		return nil, errors.ErrorGrpcClientNotExists
	}
	client, ok = res.(*Client)
	if !ok {
		return nil, errors.ErrorGrpcInvalidType
	}
	return client, nil
}

func (svc *Service) AddClient(opts *ClientOptions) (err error) {
	_, ok := svc.clientsMap.LoadAndDelete(opts.Address)
	if ok {
		return errors.ErrorGrpcClientAlreadyExists
	}
	client, err := NewClient(opts)
	if err != nil {
		return err
	}
	if err := client.Start(); err != nil {
		return err
	}
	svc.clientsMap.Store(client.opts.Address, client)
	return nil
}

func (svc *Service) RemoveClient(address Address) (err error) {
	client, err := svc.GetClient(address)
	if err != nil {
		return err
	}
	if err := client.Stop(); err != nil {
		_ = trace.TraceError(err)
	}
	svc.clientsMap.Delete(address)
	return nil
}

type ServiceOptions struct {
	Local   Address
	Remotes []Address
}

func NewService(opts *ServiceOptions) (svc *Service, err error) {
	if opts == nil {
		opts = &ServiceOptions{}
	}

	// service
	svc = &Service{
		server:     nil,
		clientsMap: sync.Map{},
		opts:       opts,
	}

	// server
	svc.server, err = NewServer(&ServerOptions{
		Address: opts.Local,
	})
	if err != nil {
		return nil, err
	}

	// clients
	for _, remote := range opts.Remotes {
		if err := svc.AddClient(&ClientOptions{Address: remote}); err != nil {
			return nil, err
		}
	}

	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}
