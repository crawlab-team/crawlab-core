package apps

type ServerOption func(app ServerApp)

type DockerOption func(dck DockerApp)

func WithDockerParent(parent NodeApp) DockerOption {
	return func(dck DockerApp) {
		dck.SetParent(parent)
	}
}
