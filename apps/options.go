package apps

import "github.com/crawlab-team/crawlab-core/interfaces"

type MasterOption func(app MasterApp)

func WithMasterConfigPath(path string) MasterOption {
	return func(app MasterApp) {
		app.SetConfigPath(path)
	}
}

func WithMasterGrpcAddress(address interfaces.Address) MasterOption {
	return func(app MasterApp) {
		app.SetGrpcAddress(address)
	}
}

func WithRunOnMaster(ok bool) MasterOption {
	return func(app MasterApp) {
		app.SetRunOnMaster(ok)
	}
}

type WorkerOption func(app WorkerApp)

func WithWorkerConfigPath(path string) WorkerOption {
	return func(app WorkerApp) {
		app.SetConfigPath(path)
	}
}

func WithWorkerGrpcAddress(address interfaces.Address) WorkerOption {
	return func(app WorkerApp) {
		app.SetGrpcAddress(address)
	}
}

type ServerOption func(app ServerApp)

func WithServerConfigPath(path string) ServerOption {
	return func(app ServerApp) {
		app.SetConfigPath(path)
	}
}

func WithServerGrpcAddress(address interfaces.Address) ServerOption {
	return func(app ServerApp) {
		app.SetGrpcAddress(address)
	}
}
