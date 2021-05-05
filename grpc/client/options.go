package client

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"time"
)

type Option func(client interfaces.GrpcClient)

func WithAddress(address interfaces.Address) Option {
	return func(c interfaces.GrpcClient) {
		c.SetAddress(address)
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c interfaces.GrpcClient) {
	}
}
