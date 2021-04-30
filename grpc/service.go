package grpc

import (
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/node"
	"github.com/crawlab-team/crawlab-core/store"
	"github.com/crawlab-team/go-trace"
	"sync"
)

type Service struct {
	server     *Server
	clientsMap sync.Map
	opts       *ServiceOptions
	nodeSvc    *node.Service
	remotes    []entity.Address
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

func (svc *Service) GetServer() (svr interfaces.GrpcServer) {
	return svc.server
}

func (svc *Service) GetClient(address interfaces.Address) (c interfaces.GrpcClient, err error) {
	addr, ok := address.(*entity.Address)
	if !ok {
		return c, errors.ErrorGrpcInvalidType
	}
	res, ok := svc.clientsMap.Load(addr.String())
	if !ok {
		return nil, errors.ErrorGrpcClientNotExists
	}
	c, ok = res.(interfaces.GrpcClient)
	if !ok {
		return nil, errors.ErrorGrpcInvalidType
	}
	return c, nil
}

func (svc *Service) GetDefaultClient() (c interfaces.GrpcClient, err error) {
	err = errors.ErrorGrpcClientNotExists
	svc.clientsMap.Range(func(key, value interface{}) bool {
		var ok bool
		c, ok = value.(interfaces.GrpcClient)
		if !ok {
			err = errors.ErrorGrpcInvalidType
		} else {
			err = nil
		}
		return false
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (svc *Service) MustGetDefaultClient() (client interfaces.GrpcClient) {
	client, err := svc.GetDefaultClient()
	if err != nil {
		panic(err)
	}
	return client
}

func (svc *Service) GetAllClients() (clients []interfaces.GrpcClient, err error) {
	svc.clientsMap.Range(func(key, value interface{}) bool {
		client, ok := value.(interfaces.GrpcClient)
		if !ok {
			err = errors.ErrorGrpcInvalidType
			return false
		}
		clients = append(clients, client)
		return true
	})
	return clients, nil
}

func (svc *Service) AddClient(opts interfaces.Options) (err error) {
	opts_ := &ClientOptions{}
	if opts != nil {
		opts_ = opts.FillEmpty().(*ClientOptions)
	}
	_, ok := svc.clientsMap.LoadAndDelete(opts_.Address)
	if ok {
		return errors.ErrorGrpcClientAlreadyExists
	}
	client, err := NewClient(opts_)
	if err != nil {
		return err
	}
	if err := client.Start(); err != nil {
		return err
	}
	svc.clientsMap.Store(client.opts.Address.String(), client)
	return nil
}

func (svc *Service) DeleteClient(address interfaces.Address) (err error) {
	client, err := svc.GetClient(address)
	if err != nil {
		return err
	}
	if err := client.Stop(); err != nil {
		_ = trace.TraceError(err)
	}
	svc.clientsMap.Delete(address.String())
	return nil
}

func NewService(opts *ServiceOptions) (res2 *Service, err error) {
	if opts == nil {
		opts = &ServiceOptions{}
	}
	opts = opts.FillEmpty().(*ServiceOptions)

	// attempt to load existing service by node key
	res, err := store.NodeServiceStore.Get(opts.NodeServiceKey)
	if err != nil {
		return nil, err
	}
	nodeSvc, ok := res.(*node.Service)
	if !ok {
		return nil, errors.ErrorGrpcInvalidType
	}

	// service
	svc := &Service{
		nodeSvc:    nodeSvc,
		server:     nil,
		clientsMap: sync.Map{},
		opts:       opts,
	}

	// server
	svc.server, err = NewServer(&ServerOptions{
		NodeService: nodeSvc,
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
