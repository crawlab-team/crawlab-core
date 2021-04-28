package grpc

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/errors"
	node2 "github.com/crawlab-team/crawlab-core/node"
	grpc2 "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
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
	if nodeSvc == nil {
		nodeSvc, err = node2.GetDefaultService()
		if err != nil {
			return err
		}
	}

	// node
	grpc2.RegisterNodeServiceServer(svr.svr, NewNodeServer(nodeSvc))

	// task
	//grpc2.RegisterTaskServiceServer(svr.svr, TaskService)

	return nil
}

type ServerOptions struct {
	NodeService *node2.Service
	Address     Address
}

var DefaultServerOptions = &ServerOptions{
	Address: NewAddress(nil),
}

func NewServer(opts *ServerOptions) (server *Server, err error) {
	if opts == nil {
		opts = DefaultServerOptions
	}

	// recovery options
	var recoveryFunc grpc_recovery.RecoveryHandlerFunc
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(recoveryFunc),
	}

	// construct server
	server = &Server{
		svr: grpc.NewServer(
			grpc_middleware.WithUnaryServerChain(
				grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
			),
			grpc_middleware.WithStreamServerChain(
				grpc_recovery.StreamServerInterceptor(recoveryOpts...),
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
