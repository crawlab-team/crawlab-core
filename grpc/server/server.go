package server

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/node"
	. "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	nodeSvc interfaces.NodeService

	svr     *grpc.Server
	l       net.Listener
	address interfaces.Address
}

func (svr *Server) Init() (err error) {
	// register
	if err := svr.Register(); err != nil {
		return err
	}

	return nil
}

func (svr *Server) Start() (err error) {
	// grpc server binding address
	address := svr.address.String()

	// listener
	svr.l, err = net.Listen("tcp", address)
	if err != nil {
		_ = trace.TraceError(err)
		return errors.ErrorGrpcServerFailedToListen
	}
	log.Infof("grpc server listens to %s", address)

	// start grpc server
	go func() {
		if err := svr.svr.Serve(svr.l); err != nil {
			if err != grpc.ErrServerStopped {
				_ = trace.TraceError(err)
				log.Error(errors.ErrorGrpcServerFailedToServe.Error())
			}
		}
	}()

	return nil
}

func (svr *Server) Stop() (err error) {
	svr.svr.GracefulStop()
	_ = svr.l.Close()
	log.Infof("grpc server stopped")
	return nil
}

func (svr *Server) Register() (err error) {
	nodeSvc := svr.nodeSvc
	if !nodeSvc.IsMaster() {
		return
	}

	masterNodeSvc, ok := nodeSvc.(interfaces.NodeMasterService)
	if !ok {
		return errors.ErrorGrpcInvalidType
	}

	// model delegate
	RegisterModelDelegateServiceServer(svr.svr, *NewModelDelegateServer(masterNodeSvc))

	// node
	RegisterNodeServiceServer(svr.svr, *NewNodeServer(masterNodeSvc))

	// task
	//grpc2.RegisterTaskServiceServer(svr.svr, TaskService)

	return nil
}

func (svr *Server) SetAddress(address interfaces.Address) {
	svr.address = address
}

func NewServer(nodeSvc interfaces.NodeService, opts ...Option) (svr *Server, err error) {
	// recovery options
	var recoveryFunc grpc_recovery.RecoveryHandlerFunc
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(recoveryFunc),
	}

	// construct server
	svr = &Server{
		svr: grpc.NewServer(
			grpc_middleware.WithUnaryServerChain(
				grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
			),
			grpc_middleware.WithStreamServerChain(
				grpc_recovery.StreamServerInterceptor(recoveryOpts...),
			),
		),
	}

	// options
	for _, opt := range opts {
		opt(svr)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(node.NewService); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(nodeSvc interfaces.NodeService) {
		svr.nodeSvc = nodeSvc
	}); err != nil {
		return nil, err
	}

	// initialize
	if err := svr.Init(); err != nil {
		return nil, err
	}

	return svr, nil
}
