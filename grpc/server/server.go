package server

import (
	"encoding/json"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/node/config"
	. "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"go.uber.org/dig"
	"go/types"
	"google.golang.org/grpc"
	"net"
	"sync"
)

type Server struct {
	// dependencies
	nodeCfgSvc       interfaces.NodeConfigService
	nodeSvr          *NodeServer
	modelDelegateSvr *ModelDelegateServer

	// settings variables
	cfgPath string
	address interfaces.Address

	// internals
	svr  *grpc.Server
	l    net.Listener
	subs sync.Map
}

func (svr *Server) Init() (err error) {
	// register
	if err := svr.Register(); err != nil {
		return err
	}

	return nil
}

func (svr *Server) Start() (err error) {
	// grpc server binding address
	address := svr.address.String()

	// listener
	svr.l, err = net.Listen("tcp", address)
	if err != nil {
		_ = trace.TraceError(err)
		return errors.ErrorGrpcServerFailedToListen
	}
	log.Infof("grpc server listens to %s", address)

	// start grpc server
	go func() {
		if err := svr.svr.Serve(svr.l); err != nil {
			if err != grpc.ErrServerStopped {
				_ = trace.TraceError(err)
				log.Error(errors.ErrorGrpcServerFailedToServe.Error())
			}
		}
	}()

	return nil
}

func (svr *Server) Stop() (err error) {
	svr.svr.GracefulStop()
	_ = svr.l.Close()
	log.Infof("grpc server stopped")
	return nil
}

func (svr *Server) Register() (err error) {
	// model delegate
	RegisterModelDelegateServiceServer(svr.svr, *svr.modelDelegateSvr)

	// node
	RegisterNodeServiceServer(svr.svr, *svr.nodeSvr)

	// task
	//grpc2.RegisterTaskServiceServer(svr.svr, TaskService)

	return nil
}

func (svr *Server) SetAddress(address interfaces.Address) {
	svr.address = address
}

func (svr *Server) GetConfigPath() (path string) {
	return svr.cfgPath
}

func (svr *Server) SetConfigPath(path string) {
	svr.cfgPath = path
}

func (svr *Server) GetSubscribe(key string) (sub interfaces.GrpcSubscribe, err error) {
	res, ok := svr.subs.Load(key)
	if !ok {
		return nil, errors.ErrorNodeStreamNotFound
	}
	sub, ok = res.(interfaces.GrpcSubscribe)
	if !ok {
		return nil, errors.ErrorNodeInvalidType
	}
	return sub, nil
}

func (svr *Server) SetSubscribe(key string, sub interfaces.GrpcSubscribe) {
	svr.subs.Store(key, sub)
}

func (svr *Server) DeleteSubscribe(key string) {
	svr.subs.Delete(key)
}

func (svr *Server) SendStreamMessage(nodeKey string, code StreamMessageCode, d interface{}) (err error) {
	var data []byte
	switch d.(type) {
	case types.Nil:
		// do nothing
	case []byte:
		data = d.([]byte)
	default:
		var err error
		data, err = json.Marshal(d)
		if err != nil {
			panic(err)
		}
	}
	sub, err := svr.GetSubscribe(nodeKey)
	if err != nil {
		return err
	}
	msg := &StreamMessage{
		Code:    code,
		NodeKey: svr.nodeCfgSvc.GetNodeKey(),
		Data:    data,
	}
	return sub.GetStream().Send(msg)
}

func NewServer(opts ...Option) (svr2 interfaces.GrpcServer, err error) {
	// recovery options
	var recoveryFunc grpc_recovery.RecoveryHandlerFunc
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(recoveryFunc),
	}

	// server
	svr := &Server{
		cfgPath: config.DefaultConfigPath,
		address: entity.NewAddress(&entity.AddressOptions{
			Host: constants.DefaultGrpcServerHost,
			Port: constants.DefaultGrpcServerPort,
		}),
		svr: grpc.NewServer(
			grpc_middleware.WithUnaryServerChain(
				grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
			),
			grpc_middleware.WithStreamServerChain(
				grpc_recovery.StreamServerInterceptor(recoveryOpts...),
			),
		),
		subs: sync.Map{},
	}

	// options
	for _, opt := range opts {
		opt(svr)
	}

	// dependency injection
	c := dig.New()
	if err := c.Provide(config.ProvideConfigService(svr.GetConfigPath())); err != nil {
		return nil, err
	}
	if err := c.Provide(NewModelDelegateServer); err != nil {
		return nil, err
	}
	if err := c.Provide(ProvideNodeServer(svr)); err != nil {
		return nil, err
	}
	if err := c.Invoke(func(nodeCfgSvc interfaces.NodeConfigService, modelDelegateSvr *ModelDelegateServer, nodeSvr *NodeServer) {
		svr.nodeCfgSvc = nodeCfgSvc
		svr.modelDelegateSvr = modelDelegateSvr
		svr.nodeSvr = nodeSvr
	}); err != nil {
		return nil, err
	}

	// initialize
	if err := svr.Init(); err != nil {
		return nil, err
	}

	return svr, nil
}

func ProvideServer(path string) func() (res interfaces.GrpcServer, err error) {
	return func() (res interfaces.GrpcServer, err error) {
		return NewServer(WithConfigPath(path))
	}
}
