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
	address := c.opts.Address.String()

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
	address := c.opts.Address.String()

	// timeout context
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.opts.TimeoutSeconds)*time.Second)
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

func NewClient(opts *ClientOptions) (client *Client, err error) {
	if opts == nil {
		opts = &ClientOptions{}
	}
	client = &Client{
		conn: nil,
		opts: opts.FillEmpty().(*ClientOptions),
	}
	if err := client.Init(); err != nil {
		return nil, err
	}
	return client, nil
}
