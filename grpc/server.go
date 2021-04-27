package grpc

import (
	"fmt"
	"github.com/apex/log"
	grpc2 "github.com/crawlab-team/crawlab-grpc"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	svr  *grpc.Server
	ch   chan int
	opts *ServerOptions
}

func (svc *Server) Init() (err error) {

	// register
	svc.Register()

	return nil
}

func (svc *Server) Start() {
	// grpc server binding address
	host := svc.opts.Address.Host
	port := svc.opts.Address.Port
	address := fmt.Sprintf("%s:%s", host, port)

	// listen
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}

	// start grpc server
	go func() {
		if err := svc.svr.Serve(listen); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// wait for signal
	svc.waitForStop()
}

func (svc *Server) Stop() {
	svc.ch <- 1
	<-svc.ch
}

func (svc *Server) Register() {
	// node
	grpc2.RegisterNodeServiceServer(svc.svr, NodeService)

	// task
	grpc2.RegisterTaskServiceServer(svc.svr, TaskService)
}

func (svc *Server) waitForStop() {
	for {
		sig := <-svc.ch
		if sig > 0 {
			svc.svr.GracefulStop()
			svc.ch <- 1
			return
		}
	}
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
		ch:   make(chan int),
		opts: opts,
	}
	if err := server.Init(); err != nil {
		return nil, err
	}
	return server, nil
}
