package apps

import (
	"github.com/apex/log"
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
	dck *Docker

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
	// log node info
	app.logNodeInfo()

	if utils.IsMaster() {

		// initialize controllers
		if err := controllers.InitControllers(); err != nil {
			panic(err)
		}
	}
}

func (app *Server) Start() {
	if utils.IsMaster() {
		// start docker app
		if utils.IsDocker() {
			go app.dck.Start()
		}

		// start api
		go app.api.Start()

		// import demo
		if utils.IsDemo() {
			go app.importDemo()
		}
	}

	// start node service
	go app.nodeSvc.Start()
}

func (app *Server) Wait() {
	<-app.quit
}

func (app *Server) Stop() {
	app.api.Stop()
	app.quit <- 1
}

func (app *Server) importDemo() {
	for {
		if app.api.Ready() {
			break
		}
		time.Sleep(1 * time.Second)
	}
	_ = utils.ImportDemo()
}

func (app *Server) logNodeInfo() {
	log.Infof("current node type: %s", utils.GetNodeType())
	if utils.IsDocker() {
		log.Infof("running in docker container")
	}
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

	// master modules
	if utils.IsMaster() {
		// api
		svr.api = GetApi()

		// docker
		if utils.IsDocker() {
			svr.dck = GetDocker()
		}
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
