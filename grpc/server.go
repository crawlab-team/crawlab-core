package grpc

import (
	"fmt"
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
	host := svc.opts.Address.Host
	port := svc.opts.Address.Port
	address := fmt.Sprintf("%s:%s", host, port)

	// listen
	listen, err := net.Listen("tcp", address)
	if err != nil {
		_ = trace.TraceError(err)
		return errors.ErrorGrpcServerFailedToListen
	}

	// start grpc server
	go func() {
		if err := svc.svr.Serve(listen); err != nil {
			_ = trace.TraceError(err)
			log.Error(errors.ErrorGrpcServerFailedToServe.Error())
		}
	}()

	return nil
}

func (svc *Server) Stop() (err error) {
	svc.svr.GracefulStop()
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
