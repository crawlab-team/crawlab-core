package grpc

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/store"
)

type Service struct {
	server  *Server
	client  interfaces.GrpcClient
	opts    *ServiceOptions
	nodeSvc interfaces.NodeService
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
	if svc.server != nil {
		if err := svc.server.Stop(); err != nil {
			return err
		}
	}

	// stop client
	if svc.client != nil {
		if err := svc.client.Stop(); err != nil {
			return err
		}
	}

	return nil
}

func (svc *Service) GetServer() (svr interfaces.GrpcServer) {
	return svc.server
}

func (svc *Service) GetClient() (c interfaces.GrpcClient) {
	return svc.client
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
	nodeSvc, ok := res.(interfaces.NodeService)
	if !ok {
		return nil, errors.ErrorGrpcInvalidType
	}

	// service
	svc := &Service{
		nodeSvc: nodeSvc,
		server:  nil,
		client:  nil,
		opts:    opts,
	}

	if svc.nodeSvc.IsMaster() {
		// master server
		svc.server, err = NewServer(&ServerOptions{
			NodeService: nodeSvc,
			Address:     opts.Local,
		})
		if err != nil {
			return nil, err
		}
	} else {
		// worker client
		svc.client, err = NewClient(&ClientOptions{
			Address: opts.Remote,
		})
		if err != nil {
			return nil, err
		}
	}

	if err := svc.Init(); err != nil {
		return nil, err
	}

	return svc, nil
}
