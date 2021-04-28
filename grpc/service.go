package grpc

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/node"
	"github.com/crawlab-team/go-trace"
	"sync"
)

type Service struct {
	server     *Server
	clientsMap sync.Map
	opts       *ServiceOptions
	nodeSvc    *node.Service
	remotes    []Address
}

func (svc *Service) Init() (err error) {
	// start server
	if err := svc.server.Start(); err != nil {
		return err
	}
	return nil
}

func (svc *Service) Stop() (err error) {
	// stop server
	if err := svc.server.Stop(); err != nil {
		return err
	}

	// stop and delete all clients
	clients, err := svc.GetAllClients()
	for _, client := range clients {
		if err := client.Stop(); err != nil {
			return err
		}
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

func (svc *Service) GetDefaultClient() (client *Client, err error) {
	err = errors.ErrorGrpcClientNotExists
	svc.clientsMap.Range(func(key, value interface{}) bool {
		var ok bool
		client, ok = value.(*Client)
		if !ok {
			err = errors.ErrorGrpcInvalidType
		} else {
			err = nil
		}
		return false
	})
	return client, err
}

func (svc *Service) GetAllClients() (clients []*Client, err error) {
	svc.clientsMap.Range(func(key, value interface{}) bool {
		client, ok := value.(*Client)
		if !ok {
			err = errors.ErrorGrpcInvalidType
			return false
		}
		clients = append(clients, client)
		return true
	})
	return clients, nil
}

func (svc *Service) AddClient(opts *ClientOptions) (err error) {
	if opts == nil {
		opts = DefaultClientOptions
	}
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

func (svc *Service) DeleteClient(address Address) (err error) {
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

var serviceMap = sync.Map{}

type ServiceOptions struct {
	NodeService *node.Service
	Local       Address
	Remotes     []Address
}

var DefaultServiceOptions = &ServiceOptions{
	Local:   NewAddress(nil),
	Remotes: []Address{},
}

func NewService(opts *ServiceOptions) (svc *Service, err error) {
	if opts == nil {
		opts = DefaultServiceOptions
	}
	if opts.NodeService == nil {
		opts.NodeService, _ = node.GetDefaultService()
	}

	// attempt to load existing service by node key
	res, ok := serviceMap.Load(opts.NodeService.GetNodeKey())
	if ok {
		svc, ok := res.(*Service)
		if !ok {
			return nil, errors.ErrorGrpcInvalidType
		}
		return svc, nil
	}

	// service
	svc = &Service{
		nodeSvc:    opts.NodeService,
		server:     nil,
		clientsMap: sync.Map{},
		opts:       opts,
	}

	// server
	svc.server, err = NewServer(&ServerOptions{
		NodeService: opts.NodeService,
		Address:     opts.Local,
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

var service *Service

func GetService() (svc *Service, err error) {
	if service == nil {
		service, err = NewService(nil)
		if err != nil {
			return nil, err
		}
	}
	return service, nil
}
