package client

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"time"
)

type Option func(client interfaces.GrpcClient)

func WithConfigPath(path string) Option {
	return func(c interfaces.GrpcClient) {
		c.SetConfigPath(path)
	}
}

func WithAddress(address interfaces.Address) Option {
	return func(c interfaces.GrpcClient) {
		c.SetAddress(address)
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c interfaces.GrpcClient) {
	}
}

type PoolOption func(p interfaces.GrpcClientPool)

func WithPoolConfigPath(path string) PoolOption {
	return func(c interfaces.GrpcClientPool) {
		c.SetConfigPath(path)
	}
}

func WithPoolSize(size int) PoolOption {
	return func(c interfaces.GrpcClientPool) {
		c.SetSize(size)
	}
}
