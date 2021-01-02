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
	Stop()
}

func NewCrawlabGrpcService() (s *CrawlabGrpcService, err error) {
	s = &CrawlabGrpcService{
		server: grpc.NewServer(), // grpc server
		ch:     make(chan int),   // signal channel
	}
	return s, nil
}

type CrawlabGrpcService struct {
	CrawlabGrpcServiceInterface
	server *grpc.Server
	ch     chan int
	rs     *TaskResultItemService
}

func (s *CrawlabGrpcService) Init() (err error) {
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

	// register services
	pb.RegisterTaskServiceServer(s.server, TaskService)

	// wait for signal
	go func() {
		for {
			sig := <-s.ch
			if sig > 0 {
				s.server.Stop()
				return
			}
		}
	}()

	// start task result item service
	s.rs = NewTaskResultItemService(&TaskResultItemServiceOptions{
		FlushWaitSeconds: viper.GetInt("grpc.task.flushWaitSeconds"),
	})

	// start grpc server
	if err := s.server.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}

	return nil
}

func (s *CrawlabGrpcService) Stop() {
	s.ch <- 1
}
