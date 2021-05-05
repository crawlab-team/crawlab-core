package client

import (
	"context"
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	grpc2 "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"google.golang.org/grpc"
	"time"
)

type Client struct {
	conn    *grpc.ClientConn
	address interfaces.Address
	timeout time.Duration

	ModelDelegateClient grpc2.ModelDelegateServiceClient
	NodeClient          grpc2.NodeServiceClient
	TaskClient          grpc2.TaskServiceClient
}

func (c *Client) Init() (err error) {
	// do nothing
	return nil
}

func (c *Client) Start() (err error) {
	// connect
	if err := backoff.Retry(c.connect, backoff.NewExponentialBackOff()); err != nil {
		return err
	}

	// register
	if err := c.Register(); err != nil {
		return err
	}

	return nil
}

func (c *Client) Stop() (err error) {
	// grpc server address
	address := c.address.String()

	// close connection
	if err := c.conn.Close(); err != nil {
		return err
	}
	log.Infof("grpc client disconnected from %s", address)

	return nil
}

func (c *Client) Register() (err error) {
	// model delegate
	c.ModelDelegateClient = grpc2.NewModelDelegateServiceClient(c.conn)

	// node
	c.NodeClient = grpc2.NewNodeServiceClient(c.conn)

	// task
	c.TaskClient = grpc2.NewTaskServiceClient(c.conn)

	return nil
}

func (c *Client) GetModelDelegateClient() (res grpc2.ModelDelegateServiceClient) {
	return c.ModelDelegateClient
}

func (c *Client) GetNodeClient() grpc2.NodeServiceClient {
	return c.NodeClient
}

func (c *Client) GetTaskClient() grpc2.TaskServiceClient {
	return c.TaskClient
}

func (c *Client) connect() (err error) {
	// grpc server address
	address := c.address.String()

	// timeout context
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	// connection
	c.conn, err = grpc.DialContext(ctx, address, grpc.WithInsecure())
	if err != nil {
		_ = trace.TraceError(err)
		return errors.ErrorGrpcClientFailedToStart
	}
	log.Infof("grpc client connected to %s", address)

	return nil
}

func (c *Client) SetAddress(address interfaces.Address) {
	c.address = address
}

func (c *Client) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

func NewClient(opts ...Option) (c interfaces.GrpcClient, err error) {
	c = &Client{
		address: entity.NewAddress(&entity.AddressOptions{
			Host: "localhost",
			Port: "9666",
		}),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}
