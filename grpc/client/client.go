package client

import (
	"context"
	"encoding/json"
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/node/config"
	grpc2 "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"go.uber.org/dig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"
	"io"
	"sync"
	"time"
)

type Client struct {
	// dependencies
	nodeCfgSvc interfaces.NodeConfigService

	// settings variables
	cfgPath string
	address interfaces.Address
	timeout time.Duration

	// internals
	conn   *grpc.ClientConn
	stream grpc2.NodeService_SubscribeClient
	msgCh  chan *grpc2.StreamMessage

	// grpc clients
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
	if err := c.connect(); err != nil {
		return err
	}

	// register rpc services
	if err := c.Register(); err != nil {
		return err
	}

	// subscribe
	if err := c.subscribe(); err != nil {
		return err
	}

	// handle stream message
	go c.handleStreamMessage()

	return nil
}

func (c *Client) Stop() (err error) {
	// grpc server address
	address := c.address.String()

	// unsubscribe
	if err := c.unsubscribe(); err != nil {
		return err
	}
	log.Infof("grpc client unsubscribed from %s", address)

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

	// log
	log.Infof("grpc client registered client services")

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

func (c *Client) SetAddress(address interfaces.Address) {
	c.address = address
}

func (c *Client) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

func (c *Client) Context() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.timeout)
}

func (c *Client) NewRequest(d interface{}) (req *grpc2.Request) {
	var data []byte
	switch d.(type) {
	case []byte:
		data = d.([]byte)
	default:
		var err error
		data, err = json.Marshal(d)
		if err != nil {
			panic(err)
		}
	}
	return &grpc2.Request{
		NodeKey: c.nodeCfgSvc.GetNodeKey(),
		Data:    data,
	}
}

func (c *Client) GetConfigPath() (path string) {
	return c.cfgPath
}

func (c *Client) SetConfigPath(path string) {
	c.cfgPath = path
}

func (c *Client) GetMessageChannel() (msgCh chan *grpc2.StreamMessage) {
	return c.msgCh
}

func (c *Client) Restart() (err error) {
	if c.needRestart() {
		return c.Start()
	}
	return nil
}

func (c *Client) connect() (err error) {
	return backoff.RetryNotify(c._connect, backoff.NewExponentialBackOff(), func(err error, duration time.Duration) {
		log.Errorf("grpc client connect error: %v. reattempt in %.1f seconds...", err, duration.Seconds())
		trace.PrintError(err)
	})
}

func (c *Client) _connect() (err error) {
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

func (c *Client) subscribe() (err error) {
	return backoff.RetryNotify(c._subscribe, backoff.NewExponentialBackOff(), func(err error, duration time.Duration) {
		log.Errorf("grpc client subscribe error: %v. reattempt in %.1f seconds...", err, duration.Seconds())
		trace.PrintError(err)
	})
}

func (c *Client) _subscribe() (err error) {
	req := c.NewRequest(&entity.NodeInfo{
		Key:      c.nodeCfgSvc.GetNodeKey(),
		IsMaster: false,
	})
	c.stream, err = c.GetNodeClient().Subscribe(context.Background(), req)
	if err != nil {
		return trace.TraceError(err)
	}

	// log
	log.Infof("grpc client subscribed to remote server")

	return nil
}

func (c *Client) unsubscribe() (err error) {
	req := c.NewRequest(&entity.NodeInfo{
		Key:      c.nodeCfgSvc.GetNodeKey(),
		IsMaster: false,
	})
	if _, err = c.GetNodeClient().Unsubscribe(context.Background(), req); err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (c *Client) handleStreamMessage() {
	for {
		// resubscribe if stream is set to nil
		if c.stream == nil {
			if err := backoff.RetryNotify(c.subscribe, backoff.NewExponentialBackOff(), func(err error, duration time.Duration) {
				_ = trace.TraceError(err)
				return
			}); err != nil {
				log.Errorf("subscribe")
				return
			}
		}

		// receive stream message
		msg, err := c.stream.Recv()
		if err != nil {
			// end
			if err == io.EOF {
				log.Infof("received EOF signal, disconnecting")
				return
			}

			// connection closing
			s, ok := status.FromError(err)
			if ok && s.Code() == codes.Canceled {
				log.Infof("received Canceled signal, disconnecting")
				return
			}

			// error
			c.stream = nil
			time.Sleep(1 * time.Second)
			continue
		}

		// send stream message to channel
		c.msgCh <- msg
	}
}

func (c *Client) needRestart() bool {
	switch c.conn.GetState() {
	case connectivity.Shutdown, connectivity.TransientFailure:
		return true
	case connectivity.Idle, connectivity.Connecting, connectivity.Ready:
		return false
	default:
		return false
	}
}

func NewClient(opts ...Option) (res interfaces.GrpcClient, err error) {
	// client
	client := &Client{
		cfgPath: config.DefaultConfigPath,
		address: entity.NewAddress(&entity.AddressOptions{
			Host: constants.DefaultGrpcClientRemoteHost,
			Port: constants.DefaultGrpcClientRemotePort,
		}),
		timeout: 10 * time.Second,
		msgCh:   make(chan *grpc2.StreamMessage),
	}

	// apply options
	for _, opt := range opts {
		opt(client)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(config.ProvideConfigService(client.GetConfigPath())); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(nodeCfgSvc interfaces.NodeConfigService) {
		client.nodeCfgSvc = nodeCfgSvc
	}); err != nil {
		return nil, err
	}

	// init
	if err := client.Init(); err != nil {
		return nil, err
	}

	return client, nil
}

func ProvideClient(path string) func() (res interfaces.GrpcClient, err error) {
	return func() (res interfaces.GrpcClient, err error) {
		return NewClient(WithConfigPath(path))
	}
}

var clientStore = sync.Map{}

func GetClient(path string) (c interfaces.GrpcClient, err error) {
	res, ok := clientStore.Load(path)
	if !ok {
		return createClient(path)
	}
	c, ok = res.(interfaces.GrpcClient)
	if !ok {
		return createClient(path)
	}
	return c, nil
}

func ForceGetClient(path string) (p interfaces.GrpcClient, err error) {
	return createClient(path)
}

func createClient(path string) (client2 interfaces.GrpcClient, err error) {
	c := dig.New()
	if err := c.Provide(ProvideClient(path)); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := c.Invoke(func(client interfaces.GrpcClient) {
		client2 = client
	}); err != nil {
		return nil, trace.TraceError(err)
	}
	if err := client2.Start(); err != nil {
		return nil, err
	}
	clientStore.Store(path, client2)
	return client2, nil
}
