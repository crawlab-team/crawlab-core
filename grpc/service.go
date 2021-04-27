package grpc

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"sync"
)

type Service struct {
	server     *Server
	clientsMap sync.Map
}

func (svc *Service) Init() (err error) {
	if err := svc.server.Init(); err != nil {
		return err
	}
	return nil
}

func (svc *Service) GetClient(address Address) (client *Client, err error) {
	res, ok := svc.clientsMap.Load(address)
	if !ok {
		return nil, errors.ErrorGrpcClientNotExists
	}
	client = res.(*Client)
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
	svc.clientsMap.Store(client.opts.Address, client)
	return nil
}

func (svc *Service) RemoveClient(address Address) (err error) {
	_, ok := svc.clientsMap.LoadAndDelete(address)
	if !ok {
		return errors.ErrorGrpcClientNotExists
	}
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
	server, err := NewServer(&ServerOptions{
		Address: opts.Local,
	})
	if err != nil {
		return nil, err
	}
	svc = &Service{
		server:     server,
		clientsMap: sync.Map{},
	}
	return svc, nil
}
