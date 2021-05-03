package grpc

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	. "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	. "github.com/grpc-ecosystem/go-grpc-middleware"
	. "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	svr  *grpc.Server
	l    net.Listener
	opts *ServerOptions
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
	address := svr.opts.Address.String()

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
	nodeSvc := svr.opts.NodeService
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

func NewServer(opts *ServerOptions) (server *Server, err error) {
	if opts == nil {
		opts = &ServerOptions{}
	}
	opts = opts.FillEmpty().(*ServerOptions)

	// recovery options
	var recoveryFunc RecoveryHandlerFunc
	recoveryOpts := []Option{
		WithRecoveryHandler(recoveryFunc),
	}

	// construct server
	server = &Server{
		svr: grpc.NewServer(
			WithUnaryServerChain(
				UnaryServerInterceptor(recoveryOpts...),
			),
			WithStreamServerChain(
				StreamServerInterceptor(recoveryOpts...),
			),
		),
		opts: opts,
	}

	// initialize
	if err := server.Init(); err != nil {
		return nil, err
	}

	return server, nil
}
