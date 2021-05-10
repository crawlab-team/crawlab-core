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

func WithTimeout(timeout time.Duration) Option {
	return func(c interfaces.GrpcClient) {
	}
}
