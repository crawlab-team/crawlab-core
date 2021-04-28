package grpc

import (
	"context"
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab-core/errors"
	grpc2 "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"google.golang.org/grpc"
	"time"
)

type Client struct {
	conn *grpc.ClientConn
	opts *ClientOptions

	NodeClient grpc2.NodeServiceClient
	TaskClient grpc2.TaskServiceClient
}

func (svc *Client) Init() (err error) {
	// register
	if err := svc.Register(); err != nil {
		return err
	}

	return nil
}

func (svc *Client) Start() (err error) {
	if err := backoff.Retry(svc.connect, backoff.NewExponentialBackOff()); err != nil {
		return err
	}

	return nil
}

func (svc *Client) Stop() (err error) {
	// grpc server address
	address := svc.opts.Address.String()

	// close connection
	if err := svc.conn.Close(); err != nil {
		return err
	}
	log.Infof("grpc client disconnected from %s", address)

	return nil
}

func (svc *Client) Register() (err error) {
	// node
	svc.NodeClient = grpc2.NewNodeServiceClient(svc.conn)

	// task
	svc.TaskClient = grpc2.NewTaskServiceClient(svc.conn)

	return nil
}

func (svc *Client) connect() (err error) {
	// grpc server address
	address := svc.opts.Address.String()

	// timeout context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(svc.opts.TimeoutSeconds)*time.Second)
	defer cancel()

	// connection
	svc.conn, err = grpc.DialContext(ctx, address, grpc.WithInsecure())
	if err != nil {
		_ = trace.TraceError(err)
		return errors.ErrorGrpcClientFailedToStart
	}
	log.Infof("grpc client connected to %s", address)

	return nil
}

type ClientOptions struct {
	Address        Address
	TimeoutSeconds int
}

var DefaultClientOptions = &ClientOptions{
	Address:        NewAddress(nil),
	TimeoutSeconds: 30,
}

func NewClient(opts *ClientOptions) (client *Client, err error) {
	if opts == nil {
		opts = &ClientOptions{
			Address: NewAddress(nil),
		}
	}
	if opts.TimeoutSeconds == 0 {
		opts.TimeoutSeconds = DefaultClientOptions.TimeoutSeconds
	}
	client = &Client{
		conn: nil,
		opts: opts,
	}
	if err := client.Init(); err != nil {
		return nil, err
	}
	return client, nil
}
