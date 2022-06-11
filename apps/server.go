package apps

import (
	"github.com/crawlab-team/crawlab-core/config"
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/node/service"
	"github.com/crawlab-team/crawlab-core/utils"
	"time"
)

type Server struct {
	// settings
	grpcAddress interfaces.Address

	// dependencies
	interfaces.WithConfigPath
	nodeSvc interfaces.NodeService

	// modules
	api *Api

	// internals
	quit chan int
}

func (app *Server) SetGrpcAddress(address interfaces.Address) {
	app.grpcAddress = address
}

func (app *Server) GetApi() (api *Api) {
	return app.api
}

func (app *Server) GetNodeService() (svc interfaces.NodeService) {
	return app.nodeSvc
}

func (app *Server) Init() {
	// initialize controllers
	if utils.IsMaster() {
		if err := controllers.InitControllers(); err != nil {
			panic(err)
		}
	}
}

func (app *Server) Start() {
	if utils.IsMaster() {
		go app.api.Start()
	}
	go app.nodeSvc.Start()

	// import demo
	if utils.IsMaster() && utils.IsDemo() {
		go func() {
			for {
				if app.api.Ready() {
					break
				}
				time.Sleep(1 * time.Second)
			}
			_ = utils.ImportDemo()
		}()
	}
}

func (app *Server) Wait() {
	<-app.quit
}

func (app *Server) Stop() {
	app.api.Stop()
	app.quit <- 1
}

func NewServer(opts ...ServerOption) (app ServerApp) {
	// server
	svr := &Server{
		WithConfigPath: config.NewConfigPathService(),
		quit:           make(chan int, 1),
	}

	// apply options
	for _, opt := range opts {
		opt(svr)
	}

	// service options
	var svcOpts []service.Option
	if svr.grpcAddress != nil {
		svcOpts = append(svcOpts, service.WithAddress(svr.grpcAddress))
	}

	// api
	if utils.IsMaster() {
		svr.api = GetApi()
	}

	// node service
	var err error
	if utils.IsMaster() {
		svr.nodeSvc, err = service.NewMasterService(svcOpts...)
	} else {
		svr.nodeSvc, err = service.NewWorkerService(svcOpts...)
	}
	if err != nil {
		panic(err)
	}

	return svr
}

var server ServerApp

func GetServer(opts ...ServerOption) ServerApp {
	if server != nil {
		return server
	}
	server = NewServer(opts...)
	return server
}
