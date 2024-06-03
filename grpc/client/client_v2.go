package client

import (
	"context"
	"encoding/json"
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/grpc/middlewares"
	"github.com/crawlab-team/crawlab-core/interfaces"
	nodeconfig "github.com/crawlab-team/crawlab-core/node/config"
	"github.com/crawlab-team/crawlab-core/utils"
	grpc2 "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"time"
)

type GrpcClientV2 struct {
	// dependencies
	nodeCfgSvc interfaces.NodeConfigService

	// settings
	address interfaces.Address
	timeout time.Duration

	// internals
	conn   *grpc.ClientConn
	stream grpc2.NodeService_SubscribeClient
	msgCh  chan *grpc2.StreamMessage
	err    error

	// clients
	NodeClient               grpc2.NodeServiceClient
	ModelBaseServiceV2Client grpc2.ModelBaseServiceV2Client
}

func (c *GrpcClientV2) Init() (err error) {
	return nil
}

func (c *GrpcClientV2) Start() (err error) {
	// connect
	if err := c.connect(); err != nil {
		return err
	}

	// register rpc services
	c.Register()

	// subscribe
	if err := c.subscribe(); err != nil {
		return err
	}

	// handle stream message
	go c.handleStreamMessage()

	return nil
}

func (c *GrpcClientV2) Stop() (err error) {
	// skip if connection is nil
	if c.conn == nil {
		return nil
	}

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

func (c *GrpcClientV2) Register() {
	// node
	c.NodeClient = grpc2.NewNodeServiceClient(c.conn)
	// model base service
	c.ModelBaseServiceV2Client = grpc2.NewModelBaseServiceV2Client(c.conn)

	// log
	log.Infof("[GrpcClient] grpc client registered client services")
	log.Debugf("[GrpcClient] NodeClient: %v", c.NodeClient)
	log.Debugf("[GrpcClient] ModelBaseServiceV2Client: %v", c.ModelBaseServiceV2Client)
}

func (c *GrpcClientV2) Context() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.timeout)
}

func (c *GrpcClientV2) NewRequest(d interface{}) (req *grpc2.Request) {
	return &grpc2.Request{
		NodeKey: c.nodeCfgSvc.GetNodeKey(),
		Data:    c.getRequestData(d),
	}
}

func (c *GrpcClientV2) IsStarted() (res bool) {
	return c.conn != nil
}

func (c *GrpcClientV2) IsClosed() (res bool) {
	if c.conn != nil {
		return c.conn.GetState() == connectivity.Shutdown
	}
	return false
}

func (c *GrpcClientV2) getRequestData(d interface{}) (data []byte) {
	if d == nil {
		return data
	}
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
	return data
}

func (c *GrpcClientV2) unsubscribe() (err error) {
	req := c.NewRequest(&entity.NodeInfo{
		Key:      c.nodeCfgSvc.GetNodeKey(),
		IsMaster: false,
	})
	if _, err = c.NodeClient.Unsubscribe(context.Background(), req); err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func (c *GrpcClientV2) connect() (err error) {
	op := func() error {
		// grpc server address
		address := c.address.String()

		// timeout context
		ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
		defer cancel()

		// connection
		// TODO: configure dial options
		var opts []grpc.DialOption
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		opts = append(opts, grpc.WithBlock())
		opts = append(opts, grpc.WithChainUnaryInterceptor(middlewares.GetAuthTokenUnaryChainInterceptor(c.nodeCfgSvc)))
		opts = append(opts, grpc.WithChainStreamInterceptor(middlewares.GetAuthTokenStreamChainInterceptor(c.nodeCfgSvc)))
		c.conn, err = grpc.DialContext(ctx, address, opts...)
		if err != nil {
			_ = trace.TraceError(err)
			return errors.ErrorGrpcClientFailedToStart
		}
		log.Infof("[GrpcClient] grpc client connected to %s", address)

		return nil
	}
	return backoff.RetryNotify(op, backoff.NewExponentialBackOff(), utils.BackoffErrorNotify("grpc client connect"))
}

func (c *GrpcClientV2) subscribe() (err error) {
	op := func() error {
		req := c.NewRequest(&entity.NodeInfo{
			Key:      c.nodeCfgSvc.GetNodeKey(),
			IsMaster: false,
		})
		c.stream, err = c.NodeClient.Subscribe(context.Background(), req)
		if err != nil {
			return trace.TraceError(err)
		}

		// log
		log.Infof("[GrpcClient] grpc client subscribed to remote server")

		return nil
	}
	return backoff.RetryNotify(op, backoff.NewExponentialBackOff(), utils.BackoffErrorNotify("grpc client subscribe"))
}

func (c *GrpcClientV2) handleStreamMessage() {
	log.Infof("[GrpcClient] start handling stream message...")
	for {
		// resubscribe if stream is set to nil
		if c.stream == nil {
			if err := backoff.RetryNotify(c.subscribe, backoff.NewExponentialBackOff(), utils.BackoffErrorNotify("grpc client subscribe")); err != nil {
				log.Errorf("subscribe")
				return
			}
		}

		// receive stream message
		msg, err := c.stream.Recv()
		log.Debugf("[GrpcClient] received message: %v", msg)
		if err != nil {
			// set error
			c.err = err

			// end
			if err == io.EOF {
				log.Infof("[GrpcClient] received EOF signal, disconnecting")
				return
			}

			// connection closed
			if c.IsClosed() {
				return
			}

			// error
			trace.PrintError(err)
			c.stream = nil
			time.Sleep(1 * time.Second)
			continue
		}

		// send stream message to channel
		c.msgCh <- msg

		// reset error
		c.err = nil
	}
}

func NewGrpcClientV2() (c *GrpcClientV2, err error) {
	client := &GrpcClientV2{
		address: entity.NewAddress(&entity.AddressOptions{
			Host: constants.DefaultGrpcClientRemoteHost,
			Port: constants.DefaultGrpcClientRemotePort,
		}),
		timeout: 10 * time.Second,
		msgCh:   make(chan *grpc2.StreamMessage),
	}
	client.nodeCfgSvc = nodeconfig.GetNodeConfigService()

	if viper.GetString("grpc.address") != "" {
		client.address, err = entity.NewAddressFromString(viper.GetString("grpc.address"))
		if err != nil {
			return nil, trace.TraceError(err)
		}
	}

	if err := client.Init(); err != nil {
		return nil, err
	}

	return client, nil
}

var _clientV2 *GrpcClientV2

func GetGrpcClientV2() (client *GrpcClientV2, err error) {
	if _clientV2 != nil {
		return _clientV2, nil
	}
	_clientV2, err = NewGrpcClientV2()
	if err != nil {
		return nil, err
	}
	return _clientV2, nil
}
