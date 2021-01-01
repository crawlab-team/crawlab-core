package grpc

import (
	"fmt"
	"github.com/apex/log"
	pb "github.com/crawlab-team/crawlab-grpc"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
)

type CrawlabGrpcServiceInterface interface {
	Init() (err error)
}

func NewCrawlabGrpcService() (s *CrawlabGrpcService, err error) {
	s = &CrawlabGrpcService{}
	return s, nil
}

type CrawlabGrpcService struct {
}

func (s *CrawlabGrpcService) Init() (err error) {
	host := viper.GetString("grpc.host")
	port := viper.GetString("grpc.port")
	address := fmt.Sprintf("%s:%s", host, port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	// construct grpc server
	server := grpc.NewServer()

	// register services
	pb.RegisterTaskServiceServer(server, TaskService)

	// start grpc server
	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}

	return nil
}
