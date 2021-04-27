package grpc

import (
	"fmt"
	"github.com/apex/log"
	grpc2 "github.com/crawlab-team/crawlab-grpc"
	"google.golang.org/grpc"
)

type Client struct {
	conn *grpc.ClientConn
	ch   chan int
	opts *ClientOptions

	nodeClient grpc2.NodeServiceClient
	taskClient grpc2.TaskServiceClient
}

func (svc *Client) Init() (err error) {
	// register
	svc.Register()

	return nil
}

func (svc *Client) Start() {
	// grpc server address
	host := svc.opts.Address.Host
	port := svc.opts.Address.Port
	address := fmt.Sprintf("%s:%s", host, port)

	// connection
	var err error
	svc.conn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to start: %v", err)
		return
	}

	// wait for signal
	svc.waitForStop()
}

func (svc *Client) Stop() {
	svc.ch <- 1
	<-svc.ch
}

func (svc *Client) Register() {
	// node
	svc.nodeClient = grpc2.NewNodeServiceClient(svc.conn)

	// task
	svc.taskClient = grpc2.NewTaskServiceClient(svc.conn)
}

func (svc *Client) waitForStop() {
	for {
		sig := <-svc.ch
		if sig > 0 {
			_ = svc.conn.Close()
			svc.ch <- 1
			return
		}
	}
}

type ClientOptions struct {
	Address Address
}

func NewClient(opts *ClientOptions) (client *Client, err error) {
	if opts == nil {
		opts = &ClientOptions{
			Address: NewAddress(nil),
		}
	}
	client = &Client{
		conn: nil,
		ch:   make(chan int),
		opts: opts,
	}
	if err := client.Init(); err != nil {
		return nil, err
	}
	return client, nil
}
