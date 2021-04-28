package grpc

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/errors"
	grpc2 "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	svr  *grpc.Server
	opts *ServerOptions
	l    net.Listener
}

func (svc *Server) Init() (err error) {
	// register
	if err := svc.Register(); err != nil {
		return err
	}

	return nil
}

func (svc *Server) Start() (err error) {
	// grpc server binding address
	address := svc.opts.Address.String()

	// listener
	svc.l, err = net.Listen("tcp", address)
	if err != nil {
		_ = trace.TraceError(err)
		return errors.ErrorGrpcServerFailedToListen
	}
	log.Infof("grpc server listens to %s", address)

	// start grpc server
	go func() {
		if err := svc.svr.Serve(svc.l); err != nil {
			if err != grpc.ErrServerStopped {
				_ = trace.TraceError(err)
				log.Error(errors.ErrorGrpcServerFailedToServe.Error())
			}
		}
	}()

	return nil
}

func (svc *Server) Stop() (err error) {
	svc.svr.GracefulStop()
	_ = svc.l.Close()
	log.Infof("grpc server stopped")
	return nil
}

func (svc *Server) Register() (err error) {
	// node
	grpc2.RegisterNodeServiceServer(svc.svr, NodeService)

	// task
	grpc2.RegisterTaskServiceServer(svc.svr, TaskService)

	return nil
}

type ServerOptions struct {
	Address Address
}

var DefaultServerOptions = &ServerOptions{
	Address: NewAddress(nil),
}

func NewServer(opts *ServerOptions) (server *Server, err error) {
	if opts == nil {
		opts = &ServerOptions{
			Address: NewAddress(nil),
		}
	}
	server = &Server{
		svr:  grpc.NewServer(),
		opts: opts,
	}
	if err := server.Init(); err != nil {
		return nil, err
	}
	return server, nil
}
