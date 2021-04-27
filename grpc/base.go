package grpc

import (
	"fmt"
	"github.com/apex/log"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
)

type Service interface {
	Init() (err error)
	Stop()
	Register()
}

type BaseService struct {
	svr *grpc.Server
	ch  chan int
}

func (svc *BaseService) Init() (err error) {
	// grpc address
	host := viper.GetString("grpc.host")
	port := viper.GetString("grpc.port")
	address := fmt.Sprintf("%s:%s", host, port)

	// listen
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	// wait for signal
	go func() {
		for {
			sig := <-svc.ch
			if sig > 0 {
				svc.svr.Stop()
				return
			}
		}
	}()

	// start grpc server
	if err := svc.svr.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}

	return nil
}

func (svc *BaseService) Stop() {
	svc.ch <- 1
}

func (svc *BaseService) Register() {
	panic("implement me")
}
